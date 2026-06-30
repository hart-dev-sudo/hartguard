# vpn-watch

A lightweight VPN leak detector for containerized VPN setups. Verifies that protected containers are routing traffic through the VPN container — not leaking their real IP.

Part of [hartguard](../) — blue team security tooling for Linux servers.

## What it does

On each check cycle, vpn-watch:
1. Confirms the VPN container is running
2. Resolves the VPN container's external IP
3. Resolves each protected container's external IP
4. Alerts if any protected container's IP doesn't match the VPN IP

## Features

- Works with any containerized VPN (gluetun, wireguard, openvpn, etc.)
- Configurable container list, check URL, and interval
- Continuous monitoring or one-shot mode (`--once`)
- Structured JSON log output
- No vendor lock-in — uses any IP-echo endpoint

## Requirements

- Go 1.22+
- Docker with containers that have `wget` available

## Installation

```bash
git clone https://github.com/hart-dev-sudo/hartguard.git
cd hartguard/vpn-watch
go build -o vpn-watch .
```

## Configuration

Edit `config.yaml`:

```yaml
vpn_container: gluetun       # your VPN container name
check_containers:
  - qbittorrent              # containers that must route through VPN
check_url: https://api.ipify.org  # any IP-echo endpoint
interval: 60                 # seconds between checks (0 = one-shot)
log_file: logs/vpn-watch.log
```

## Usage

```bash
# Run continuously (uses interval from config)
./vpn-watch

# Single check and exit
./vpn-watch --once

# Custom config path
./vpn-watch --config /etc/hartguard/vpn-watch.yaml
```

## Example Output

```
[2024-01-15 14:30:00] Running VPN check...
[OK]    VPN IP: 185.213.154.20 (gluetun)
[OK]    qbittorrent IP: 185.213.154.20 — matches VPN

[2024-01-15 14:31:00] Running VPN check...
[ALERT] LEAK: qbittorrent is using IP 93.184.216.34 (expected 185.213.154.20)
```

Log entry (JSON):
```json
{"timestamp":"2024-01-15T14:31:00Z","event":"LEAK_DETECTED","vpn_container":"gluetun","vpn_ip":"185.213.154.20","container":"qbittorrent","container_ip":"93.184.216.34"}
```

## Running Tests

```bash
go test ./internal/checker/...
```

## License

MIT
