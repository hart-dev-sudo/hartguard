package alerter

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hart-dev-sudo/hartguard/login-watch/internal/detector"
)

type logEntry struct {
	Timestamp  string `json:"timestamp"`
	SrcIP      string `json:"src_ip"`
	Username   string `json:"username"`
	EventType  string `json:"event_type"`
	Count      int    `json:"count"`
	WindowSecs int    `json:"window_secs"`
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

const (
	colorRed   = "\033[31m"
	colorGreen = "\033[32m"
	colorReset = "\033[0m"
)

func (a *Alerter) Run(alerts <-chan detector.Alert, done func()) {
	defer done()
	for alert := range alerts {
		entry := logEntry{
			Timestamp:  time.Now().UTC().Format(time.RFC3339),
			SrcIP:      alert.SrcIP,
			Username:   alert.Username,
			EventType:  alert.EventType,
			Count:      alert.Count,
			WindowSecs: alert.WindowSecs,
		}
		b, _ := json.Marshal(entry)
		a.logger.Println(string(b))

		fmt.Printf("%s[ALERT]%s Brute force from %s (user: %s) — %d failures in %ds (%s)\n",
			colorRed, colorReset, alert.SrcIP, alert.Username, alert.Count, alert.WindowSecs, alert.EventType)
	}
}
