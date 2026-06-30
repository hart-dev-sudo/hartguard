package checker

import "testing"

type mockExecutor struct {
	running map[string]bool
	ips     map[string]string
}

func (m mockExecutor) IsRunning(container string) bool {
	return m.running[container]
}

func (m mockExecutor) ContainerIP(container, _ string) (string, error) {
	return m.ips[container], nil
}

func TestVPNDown(t *testing.T) {
	exec := mockExecutor{running: map[string]bool{"vpn": false}}
	c := New("vpn", []string{"app"}, "https://api.ipify.org", exec)
	r := c.Run()
	if !r.VPNDown {
		t.Fatal("expected VPNDown when VPN container not running")
	}
}

func TestNoLeak(t *testing.T) {
	exec := mockExecutor{
		running: map[string]bool{"vpn": true},
		ips:     map[string]string{"vpn": "1.2.3.4", "app": "1.2.3.4"},
	}
	c := New("vpn", []string{"app"}, "https://api.ipify.org", exec)
	r := c.Run()
	if r.LeakDetected {
		t.Fatal("expected no leak when IPs match")
	}
}

func TestLeakDetected(t *testing.T) {
	exec := mockExecutor{
		running: map[string]bool{"vpn": true},
		ips:     map[string]string{"vpn": "1.2.3.4", "app": "9.9.9.9"},
	}
	c := New("vpn", []string{"app"}, "https://api.ipify.org", exec)
	r := c.Run()
	if !r.LeakDetected {
		t.Fatal("expected leak when IPs differ")
	}
}

func TestMultipleContainers(t *testing.T) {
	exec := mockExecutor{
		running: map[string]bool{"vpn": true},
		ips:     map[string]string{"vpn": "1.2.3.4", "app1": "1.2.3.4", "app2": "9.9.9.9"},
	}
	c := New("vpn", []string{"app1", "app2"}, "https://api.ipify.org", exec)
	r := c.Run()
	if !r.LeakDetected {
		t.Fatal("expected leak when one container IP differs")
	}
	if r.Protected[0].Leak {
		t.Fatal("app1 should not be leaking")
	}
	if !r.Protected[1].Leak {
		t.Fatal("app2 should be leaking")
	}
}
