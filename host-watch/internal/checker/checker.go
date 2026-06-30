package checker

import (
	"fmt"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type DiskStat struct {
	Path        string
	UsedPercent int
}

type ContainerStat struct {
	Name    string
	Running bool
}

type ServiceStat struct {
	Name      string
	URL       string
	Reachable bool
}

type MemStat struct {
	TotalMB     uint64
	UsedMB      uint64
	UsedPercent int
}

type Result struct {
	Disk       []DiskStat
	Memory     MemStat
	Containers []ContainerStat
	Services   []ServiceStat
}

type ServiceURL struct {
	Name string
	URL  string
}

type Config struct {
	DiskPaths   []string
	DiskWarnPct int
	Containers  []string
	ServiceURLs []ServiceURL
}

func Run(cfg Config) Result {
	return Result{
		Disk:       checkDisk(cfg.DiskPaths),
		Memory:     checkMemory(),
		Containers: checkContainers(cfg.Containers),
		Services:   checkServices(cfg.ServiceURLs),
	}
}

func checkDisk(paths []string) []DiskStat {
	stats := make([]DiskStat, 0, len(paths))
	for _, path := range paths {
		out, err := exec.Command("df", "--output=pcent", path).Output()
		if err != nil {
			continue
		}
		lines := strings.Split(strings.TrimSpace(string(out)), "\n")
		if len(lines) < 2 {
			continue
		}
		pctStr := strings.TrimSpace(strings.TrimSuffix(lines[1], "%"))
		pct, _ := strconv.Atoi(pctStr)
		stats = append(stats, DiskStat{Path: path, UsedPercent: pct})
	}
	return stats
}

func checkMemory() MemStat {
	out, err := exec.Command("free", "-m").Output()
	if err != nil {
		return MemStat{}
	}
	for _, line := range strings.Split(string(out), "\n") {
		if !strings.HasPrefix(line, "Mem:") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 3 {
			break
		}
		total, _ := strconv.ParseUint(fields[1], 10, 64)
		used, _ := strconv.ParseUint(fields[2], 10, 64)
		pct := 0
		if total > 0 {
			pct = int(used * 100 / total)
		}
		return MemStat{TotalMB: total, UsedMB: used, UsedPercent: pct}
	}
	return MemStat{}
}

func checkContainers(names []string) []ContainerStat {
	out, _ := exec.Command("docker", "ps", "--format", "{{.Names}}").Output()
	running := make(map[string]bool)
	for _, name := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		running[strings.TrimSpace(name)] = true
	}
	stats := make([]ContainerStat, 0, len(names))
	for _, name := range names {
		stats = append(stats, ContainerStat{Name: name, Running: running[name]})
	}
	return stats
}

func checkServices(urls []ServiceURL) []ServiceStat {
	client := &http.Client{Timeout: 5 * time.Second}
	stats := make([]ServiceStat, 0, len(urls))
	for _, svc := range urls {
		_, err := client.Get(svc.URL)
		stats = append(stats, ServiceStat{
			Name:      svc.Name,
			URL:       svc.URL,
			Reachable: err == nil,
		})
	}
	return stats
}

func HasWarnings(r Result, diskWarnPct int) bool {
	for _, d := range r.Disk {
		if d.UsedPercent >= diskWarnPct {
			return true
		}
	}
	for _, c := range r.Containers {
		if !c.Running {
			return true
		}
	}
	for _, s := range r.Services {
		if !s.Reachable {
			return true
		}
	}
	return false
}

func DiskSummary(d DiskStat) string {
	return fmt.Sprintf("%s: %d%% used", d.Path, d.UsedPercent)
}
