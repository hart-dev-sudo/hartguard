package alerter

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	"github.com/hart-dev-sudo/hartguard/port-scan-detector/internal/detector"
)

type logEntry struct {
	Timestamp string   `json:"timestamp"`
	SrcIP     string   `json:"src_ip"`
	ScanType  string   `json:"scan_type"`
	Ports     []uint16 `json:"ports_hit"`
	PortCount int      `json:"port_count"`
}

// Alerter consumes alerts from the detector and writes them to stdout and a log file
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

func (a *Alerter) Run(alerts <-chan detector.Alert) {
	for alert := range alerts {
		sort.Slice(alert.Ports, func(i, j int) bool {
			return alert.Ports[i] < alert.Ports[j]
		})

		entry := logEntry{
			Timestamp: time.Now().UTC().Format(time.RFC3339),
			SrcIP:     alert.SrcIP,
			ScanType:  alert.ScanType,
			Ports:     alert.Ports,
			PortCount: len(alert.Ports),
		}

		b, _ := json.Marshal(entry)
		a.logger.Println(string(b))

		fmt.Printf("[ALERT] %s scan from %s | %d ports hit: %v\n",
			alert.ScanType, alert.SrcIP, len(alert.Ports), alert.Ports)
	}
}
