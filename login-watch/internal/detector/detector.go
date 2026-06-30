package detector

import (
	"sync"
	"time"
)

type Alert struct {
	SrcIP      string
	Username   string
	EventType  string
	Count      int
	WindowSecs int
}

type entry struct {
	timestamp time.Time
	eventType string
	username  string
}

type Detector struct {
	threshold  int
	window     time.Duration
	whitelist  map[string]struct{}
	Alerts     chan Alert
	mu         sync.Mutex
	tracker    map[string][]entry
}

func New(threshold, windowSecs int, whitelist []string) *Detector {
	wl := make(map[string]struct{}, len(whitelist))
	for _, ip := range whitelist {
		wl[ip] = struct{}{}
	}
	return &Detector{
		threshold: threshold,
		window:    time.Duration(windowSecs) * time.Second,
		whitelist: wl,
		Alerts:    make(chan Alert, 100),
		tracker:   make(map[string][]entry),
	}
}

func (d *Detector) Process(srcIP, eventType, username string) {
	if _, ok := d.whitelist[srcIP]; ok {
		return
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-d.window)

	entries := d.tracker[srcIP]
	pruned := entries[:0]
	for _, e := range entries {
		if e.timestamp.After(cutoff) {
			pruned = append(pruned, e)
		}
	}
	pruned = append(pruned, entry{timestamp: now, eventType: eventType, username: username})
	d.tracker[srcIP] = pruned

	if len(pruned) >= d.threshold {
		d.Alerts <- Alert{
			SrcIP:      srcIP,
			Username:   username,
			EventType:  eventType,
			Count:      len(pruned),
			WindowSecs: int(d.window.Seconds()),
		}
		delete(d.tracker, srcIP)
	}
}
