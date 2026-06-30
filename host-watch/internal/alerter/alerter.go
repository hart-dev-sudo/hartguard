package alerter

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/chrishartserver/hartguard/host-watch/internal/checker"
)

type logEntry struct {
	Timestamp string `json:"timestamp"`
	Event     string `json:"event"`
	Detail    string `json:"detail"`
}

type Alerter struct {
	logger *log.Logger
}

func New(logFile string) (*Alerter, error) {
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("opening log file: %w", err)
	}
	return &Alerter{logger: log.New(f, "", 0)}, nil
}

func (a *Alerter) Report(r checker.Result, diskWarnPct int) {
	fmt.Println("-- Disk --")
	for _, d := range r.Disk {
		status := "OK"
		if d.UsedPercent >= diskWarnPct {
			status = "WARN"
			a.write("DISK_WARN", checker.DiskSummary(d))
		}
		fmt.Printf("[%-4s] %s\n", status, checker.DiskSummary(d))
	}

	fmt.Printf("\n-- Memory --\n")
	fmt.Printf("[INFO] %d MB used / %d MB total (%d%%)\n",
		r.Memory.UsedMB, r.Memory.TotalMB, r.Memory.UsedPercent)

	fmt.Println("\n-- Containers --")
	for _, c := range r.Containers {
		if c.Running {
			fmt.Printf("[OK  ] %s\n", c.Name)
		} else {
			fmt.Printf("[CRIT] %s — not running\n", c.Name)
			a.write("CONTAINER_DOWN", c.Name)
		}
	}

	if len(r.Services) > 0 {
		fmt.Println("\n-- Services --")
		for _, s := range r.Services {
			if s.Reachable {
				fmt.Printf("[OK  ] %s\n", s.Name)
			} else {
				fmt.Printf("[CRIT] %s unreachable (%s)\n", s.Name, s.URL)
				a.write("SERVICE_DOWN", fmt.Sprintf("%s %s", s.Name, s.URL))
			}
		}
	}
}

func (a *Alerter) write(event, detail string) {
	entry := logEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Event:     event,
		Detail:    detail,
	}
	b, _ := json.Marshal(entry)
	a.logger.Println(string(b))
}
