package checker

import (
	"fmt"
	"os/exec"
	"strings"
)

// Executor abstracts docker exec calls so the checker is testable without Docker
type Executor interface {
	ContainerIP(container, checkURL string) (string, error)
	IsRunning(container string) bool
}

type DockerExecutor struct{}

func (d DockerExecutor) IsRunning(container string) bool {
	out, err := exec.Command("docker", "ps", "--format", "{{.Names}}").Output()
	if err != nil {
		return false
	}
	for _, name := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if name == container {
			return true
		}
	}
	return false
}

func (d DockerExecutor) ContainerIP(container, checkURL string) (string, error) {
	out, err := exec.Command(
		"docker", "exec", container,
		"wget", "-qO-", "--timeout=10", checkURL,
	).Output()
	if err != nil {
		return "", fmt.Errorf("exec in %s failed: %w", container, err)
	}
	ip := strings.TrimSpace(string(out))
	if ip == "" {
		return "", fmt.Errorf("empty response from %s", container)
	}
	return ip, nil
}

// ContainerResult holds the IP check result for a single protected container
type ContainerResult struct {
	Name string
	IP   string
	Leak bool
}

// Result is the full outcome of a single VPN check cycle
type Result struct {
	VPNContainer string
	VPNIP        string
	Protected    []ContainerResult
	VPNDown      bool
	LeakDetected bool
}

type Checker struct {
	vpnContainer      string
	checkContainers   []string
	checkURL          string
	exec              Executor
}

func New(vpnContainer string, checkContainers []string, checkURL string, exec Executor) *Checker {
	return &Checker{
		vpnContainer:    vpnContainer,
		checkContainers: checkContainers,
		checkURL:        checkURL,
		exec:            exec,
	}
}

func (c *Checker) Run() Result {
	result := Result{VPNContainer: c.vpnContainer}

	if !c.exec.IsRunning(c.vpnContainer) {
		result.VPNDown = true
		return result
	}

	vpnIP, err := c.exec.ContainerIP(c.vpnContainer, c.checkURL)
	if err != nil {
		result.VPNDown = true
		return result
	}
	result.VPNIP = vpnIP

	for _, container := range c.checkContainers {
		cr := ContainerResult{Name: container}
		ip, err := c.exec.ContainerIP(container, c.checkURL)
		if err != nil {
			cr.IP = "UNKNOWN"
			cr.Leak = true
		} else {
			cr.IP = ip
			cr.Leak = ip != vpnIP
		}
		if cr.Leak {
			result.LeakDetected = true
		}
		result.Protected = append(result.Protected, cr)
	}

	return result
}
