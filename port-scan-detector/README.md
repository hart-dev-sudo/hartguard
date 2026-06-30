# port-scan-detector

A lightweight network port scan detector written in Go. Listens on a network interface, tracks incoming TCP connections per source IP using a sliding time window, and fires an alert when a single IP hits enough unique ports to indicate a scan.

Part of [hartguard](../) — blue team security tooling for Linux servers. Pairs with [login-watch](../login-watch/) to cover the full recon → brute force attack chain.

## What it detects

| Scan | TCP Flags | Description |
|------|-----------|-------------|
| SYN  | SYN | Standard nmap default scan |
| FIN  | FIN | Stealthy, bypasses some firewalls |
| NULL | (none) | No flags set |
| XMAS | FIN+PSH+URG | All "Christmas tree" flags |

## Features

- Detects SYN, FIN, NULL, and XMAS scan patterns
- Sliding window algorithm — configurable threshold and time window
- IP whitelist support
- Structured JSON log output
- CLI flags override config file
- Systemd service unit file for production deployment
- Requires root or `CAP_NET_RAW` (raw packet capture)

## Requirements

- Go 1.22+
- libpcap (`sudo apt install libpcap-dev`)
- Root or `CAP_NET_RAW` capability

## Installation

```bash
git clone https://github.com/hart-dev-sudo/hartguard.git
cd hartguard/port-scan-detector
go build -o port-scan-detector .
```

### Production install (systemd)

```bash
sudo make install
sudo systemctl enable --now port-scan-detector
sudo journalctl -u port-scan-detector -f
```

## Configuration

Copy the example config and edit for your system:

```bash
cp config.yaml.example config.yaml
```

```yaml
interface: eth0      # network interface to listen on
threshold: 10        # unique ports within window to trigger alert
window: 5            # time window in seconds
whitelist:
  - 127.0.0.1
log_file: logs/detections.log
```

**Finding your interface:** Run `ip addr show` and look for the interface with `state UP` and your local IP. Common names: `eth0`, `ens18`, `enp3s0`, `wlan0`. Example:

```
2: eth0: <BROADCAST,MULTICAST,UP,LOWER_UP> ...
    inet 192.168.1.10/24 ...
```

After `make install`, config lives at `/etc/hartguard/port-scan-detector.yaml`.

## Usage

```bash
# Run with config file
sudo ./port-scan-detector

# Override interface and threshold via flags
sudo ./port-scan-detector --interface ens18 --threshold 5 --window 3

# Custom config path
sudo ./port-scan-detector --config /etc/hartguard/port-scan-detector.yaml
```

## Example Output

```
========================================
  port-scan-detector
  interface: ens18
  threshold: 10 unique ports in 5s
========================================
[ALERT] SYN scan from 192.168.1.50 | 15 ports hit: [22 23 80 443 3306 5432 6379 8080 8443 27017]
```

Log entry (JSON):
```json
{"timestamp":"2024-01-15T14:23:01Z","src_ip":"192.168.1.50","scan_type":"SYN","ports_hit":[22,23,80,443,3306],"port_count":15}
```

## Running Tests

```bash
make test
```

## Makefile targets

| Target | Description |
|--------|-------------|
| `make build` | Compile binary |
| `make test` | Run all tests |
| `make install` | Install binary, service file, and config |
| `make uninstall` | Remove all installed files |
| `make clean` | Remove compiled binary |

## License

MIT
