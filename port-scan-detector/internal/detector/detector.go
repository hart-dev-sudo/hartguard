package detector

import (
	"sync"
	"time"
)

// ScanType maps TCP flag combinations to scan names
var ScanType = map[uint16]string{
	0x002: "SYN",
	0x001: "FIN",
	0x000: "NULL",
	0x029: "XMAS",
}

type event struct {
	timestamp time.Time
	port      uint16
}

// Alert is passed to the alerter when a scan is detected
type Alert struct {
	SrcIP    string
	ScanType string
	Ports    []uint16
}

// Detector tracks per-IP port access history and fires alerts via a channel
type Detector struct {
	threshold int
	window    time.Duration
	whitelist map[string]struct{}
	Alerts    chan Alert
	mu        sync.Mutex
	tracker   map[string][]event
}

func New(threshold int, windowSecs int, whitelist []string) *Detector {
	wl := make(map[string]struct{}, len(whitelist))
	for _, ip := range whitelist {
		wl[ip] = struct{}{}
	}
	return &Detector{
		threshold: threshold,
		window:    time.Duration(windowSecs) * time.Second,
		whitelist: wl,
		Alerts:    make(chan Alert, 100),
		tracker:   make(map[string][]event),
	}
}

func (d *Detector) Process(srcIP string, port uint16, flags uint16) {
	if _, ok := d.whitelist[srcIP]; ok {
		return
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-d.window)

	// Prune events outside the window
	events := d.tracker[srcIP]
	pruned := events[:0]
	for _, e := range events {
		if e.timestamp.After(cutoff) {
			pruned = append(pruned, e)
		}
	}
	pruned = append(pruned, event{timestamp: now, port: port})
	d.tracker[srcIP] = pruned

	// Count unique ports
	seen := make(map[uint16]struct{}, len(pruned))
	for _, e := range pruned {
		seen[e.port] = struct{}{}
	}

	if len(seen) >= d.threshold {
		ports := make([]uint16, 0, len(seen))
		for p := range seen {
			ports = append(ports, p)
		}
		scanName := ScanType[flags]
		if scanName == "" {
			scanName = "UNKNOWN"
		}
		d.Alerts <- Alert{SrcIP: srcIP, ScanType: scanName, Ports: ports}
		// Clear tracker for this IP to avoid alert spam
		delete(d.tracker, srcIP)
	}
}
