package detector

import (
	"testing"
	"time"
)

func newTestDetector(threshold, window int) *Detector {
	return New(threshold, window, []string{"192.168.1.1"})
}

func TestNoAlertBelowThreshold(t *testing.T) {
	d := newTestDetector(5, 10)
	for i := uint16(0); i < 4; i++ {
		d.Process("10.0.0.1", i, 0x002)
	}
	if len(d.Alerts) > 0 {
		t.Fatal("expected no alert below threshold")
	}
}

func TestAlertAtThreshold(t *testing.T) {
	d := newTestDetector(5, 10)
	for i := uint16(0); i < 5; i++ {
		d.Process("10.0.0.1", i, 0x002)
	}
	select {
	case alert := <-d.Alerts:
		if alert.SrcIP != "10.0.0.1" {
			t.Fatalf("unexpected src ip: %s", alert.SrcIP)
		}
		if alert.ScanType != "SYN" {
			t.Fatalf("unexpected scan type: %s", alert.ScanType)
		}
	default:
		t.Fatal("expected alert at threshold")
	}
}

func TestWhitelistIgnored(t *testing.T) {
	d := newTestDetector(5, 10)
	for i := uint16(0); i < 10; i++ {
		d.Process("192.168.1.1", i, 0x002)
	}
	if len(d.Alerts) > 0 {
		t.Fatal("whitelisted IP should not trigger alert")
	}
}

func TestWindowExpiry(t *testing.T) {
	d := newTestDetector(5, 1)
	for i := uint16(0); i < 4; i++ {
		d.Process("10.0.0.2", i, 0x002)
	}
	time.Sleep(1100 * time.Millisecond)
	d.Process("10.0.0.2", 99, 0x002)
	if len(d.Alerts) > 0 {
		t.Fatal("events outside window should not trigger alert")
	}
}

func TestScanTypeDetection(t *testing.T) {
	tests := []struct {
		flags    uint16
		expected string
	}{
		{0x002, "SYN"},
		{0x001, "FIN"},
		{0x000, "NULL"},
		{0x029, "XMAS"},
	}
	for _, tt := range tests {
		d := newTestDetector(1, 10)
		d.Process("10.0.0.3", 80, tt.flags)
		select {
		case alert := <-d.Alerts:
			if alert.ScanType != tt.expected {
				t.Errorf("flags %x: got %s, want %s", tt.flags, alert.ScanType, tt.expected)
			}
		default:
			t.Errorf("flags %x: expected alert", tt.flags)
		}
	}
}
