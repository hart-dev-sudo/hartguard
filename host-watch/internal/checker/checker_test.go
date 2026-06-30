package checker

import "testing"

func TestHasWarningsDiskThreshold(t *testing.T) {
	r := Result{Disk: []DiskStat{{Path: "/", UsedPercent: 85}}}
	if !HasWarnings(r, 80) {
		t.Fatal("expected warning when disk exceeds threshold")
	}
}

func TestNoWarningsBelowThreshold(t *testing.T) {
	r := Result{Disk: []DiskStat{{Path: "/", UsedPercent: 60}}}
	if HasWarnings(r, 80) {
		t.Fatal("expected no warning when disk is below threshold")
	}
}

func TestHasWarningsContainerDown(t *testing.T) {
	r := Result{Containers: []ContainerStat{{Name: "nginx", Running: false}}}
	if !HasWarnings(r, 80) {
		t.Fatal("expected warning when container is down")
	}
}

func TestHasWarningsServiceUnreachable(t *testing.T) {
	r := Result{Services: []ServiceStat{{Name: "Plex", URL: "http://localhost:32400", Reachable: false}}}
	if !HasWarnings(r, 80) {
		t.Fatal("expected warning when service is unreachable")
	}
}

func TestNoWarningsAllHealthy(t *testing.T) {
	r := Result{
		Disk:       []DiskStat{{Path: "/", UsedPercent: 50}},
		Containers: []ContainerStat{{Name: "nginx", Running: true}},
		Services:   []ServiceStat{{Name: "App", Reachable: true}},
	}
	if HasWarnings(r, 80) {
		t.Fatal("expected no warnings when everything is healthy")
	}
}
