package alerter

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hart-dev-sudo/hartguard/vpn-watch/internal/checker"
)

type logEntry struct {
	Timestamp    string `json:"timestamp"`
	Event        string `json:"event"`
	VPNContainer string `json:"vpn_container"`
	VPNIP        string `json:"vpn_ip,omitempty"`
	Container    string `json:"container,omitempty"`
	ContainerIP  string `json:"container_ip,omitempty"`
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

func (a *Alerter) Report(r checker.Result) {
	ts := time.Now().UTC().Format(time.RFC3339)

	if r.VPNDown {
		entry := logEntry{Timestamp: ts, Event: "VPN_DOWN", VPNContainer: r.VPNContainer}
		a.write(entry)
		fmt.Printf("[ALERT] VPN container %q is down or unreachable\n", r.VPNContainer)
		return
	}

	fmt.Printf("[OK]    VPN IP: %s (%s)\n", r.VPNIP, r.VPNContainer)

	for _, cr := range r.Protected {
		if cr.Leak {
			entry := logEntry{
				Timestamp:    ts,
				Event:        "LEAK_DETECTED",
				VPNContainer: r.VPNContainer,
				VPNIP:        r.VPNIP,
				Container:    cr.Name,
				ContainerIP:  cr.IP,
			}
			a.write(entry)
			fmt.Printf("[ALERT] LEAK: %s is using IP %s (expected %s)\n", cr.Name, cr.IP, r.VPNIP)
		} else {
			fmt.Printf("[OK]    %s IP: %s — matches VPN\n", cr.Name, cr.IP)
		}
	}
}

func (a *Alerter) write(entry logEntry) {
	b, _ := json.Marshal(entry)
	a.logger.Println(string(b))
}
