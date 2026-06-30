# port-scan-detector

A lightweight network port scan detector written in Go. Listens on a network interface, tracks incoming TCP connections per source IP using a sliding time window, and fires an alert when a single IP hits enough unique ports to indicate a scan.

Part of the [hartguard](../) — blue team tooling for Linux servers.

## Features

- Detects SYN, FIN, NULL, and XMAS scan patterns
- Sliding window algorithm — configurable threshold and time window
- IP whitelist support
- Structured JSON log output
- CLI flags override config file
- Requires root (raw packet capture)

## Detected Scan Types

| Scan | TCP Flags | Description |
|------|-----------|-------------|
| SYN  | SYN       | Standard nmap default scan |
| FIN  | FIN       | Stealthy, bypasses some firewalls |
| NULL | (none)    | No flags set |
| XMAS | FIN+PSH+URG | All "Christmas tree" flags |

## Requirements

- Go 1.22+
- libpcap (`sudo apt install libpcap-dev`)
- Root or `CAP_NET_RAW` capability

## Installation

```bash
git clone https://github.com/chrishartserver/hartguard.git
cd hartguard/port-scan-detector
go build -o port-scan-detector .
```

## Configuration

Edit `config.yaml`:

```yaml
interface: eth0      # network interface to listen on
threshold: 10        # unique ports within window to trigger alert
window: 5            # time window in seconds
whitelist:
  - 127.0.0.1        # IPs that will never trigger alerts
log_file: logs/detections.log
```

## Usage

```bash
# Run with config file
sudo ./port-scan-detector

# Override interface and threshold via flags
sudo ./port-scan-detector --interface ens18 --threshold 5 --window 3

# Custom config path
sudo ./port-scan-detector --config /etc/psd/config.yaml
```

## Example Output

```
2024/01/15 14:23:01 Starting port-scan-detector | interface=ens18 threshold=10 window=5s
[ALERT] SYN scan from 192.168.1.50 | 15 ports hit: [22 23 80 443 3306 5432 6379 8080 8443 27017 ...]
```

Log file entry (JSON):
```json
{"timestamp":"2024-01-15T14:23:01Z","src_ip":"192.168.1.50","scan_type":"SYN","ports_hit":[22,23,80,443,3306],"port_count":15}
```

## Running Tests

Tests cover the detection engine and do not require root:

```bash
go test ./internal/detector/...
```

## Roadmap

- [ ] IPv6 support
- [ ] Email / webhook alerting
- [ ] Systemd service unit file
- [ ] Rewrite in Go v2 → already Go; Python port as learning exercise

## License

MIT
