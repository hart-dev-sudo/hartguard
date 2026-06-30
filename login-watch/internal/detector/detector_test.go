package detector

import (
	"testing"
	"time"
)

func TestNoAlertBelowThreshold(t *testing.T) {
	d := New(5, 60, nil)
	for i := 0; i < 4; i++ {
		d.Process("10.0.0.1", "Failed password", "root")
	}
	if len(d.Alerts) > 0 {
		t.Fatal("expected no alert below threshold")
	}
}

func TestAlertAtThreshold(t *testing.T) {
	d := New(5, 60, nil)
	for i := 0; i < 5; i++ {
		d.Process("10.0.0.1", "Failed password", "root")
	}
	select {
	case alert := <-d.Alerts:
		if alert.SrcIP != "10.0.0.1" {
			t.Errorf("unexpected src ip: %s", alert.SrcIP)
		}
		if alert.Username != "root" {
			t.Errorf("expected username root, got %s", alert.Username)
		}
		if alert.Count != 5 {
			t.Errorf("expected count 5, got %d", alert.Count)
		}
	default:
		t.Fatal("expected alert at threshold")
	}
}

func TestWhitelistIgnored(t *testing.T) {
	d := New(3, 60, []string{"192.168.1.1"})
	for i := 0; i < 10; i++ {
		d.Process("192.168.1.1", "Failed password", "admin")
	}
	if len(d.Alerts) > 0 {
		t.Fatal("whitelisted IP should not trigger alert")
	}
}

func TestWindowExpiry(t *testing.T) {
	d := New(3, 1, nil)
	d.Process("10.0.0.2", "Failed password", "root")
	d.Process("10.0.0.2", "Failed password", "root")
	time.Sleep(1100 * time.Millisecond)
	d.Process("10.0.0.2", "Failed password", "root")
	if len(d.Alerts) > 0 {
		t.Fatal("events outside window should not trigger alert")
	}
}
