# hartguard

A personal blue team security suite built incrementally — one layover at a time.

## What is this?

hartguard is a growing collection of defensive security tools built from scratch in Go. Each tool is written to understand the underlying detection techniques, not just wrap existing utilities. Counterpart to [hartkit](https://github.com/hart-dev-sudo/hartkit) — if hartkit scans, hartguard listens.

## Tools

| Tool | Description | Status |
|------|-------------|--------|
| [port-scan-detector](./port-scan-detector/) | Detects incoming TCP port scans using a sliding window algorithm — SYN, FIN, NULL, and XMAS | In progress |
| [vpn-watch](./vpn-watch/) | Monitors containerized VPN setups for IP leaks — verified against any IP-echo endpoint | Available |

## Structure

```
hartguard/
├── port-scan-detector/    # detects incoming port scans
│   ├── internal/          # detection, sniffing, alerting packages
│   ├── main.go
│   └── config.yaml
```

## Usage

Each tool builds and runs standalone:

```bash
cd <tool>
go build -o <tool> .
sudo ./<tool>
```

See each tool's README for full usage and configuration.

## Requirements

- Go 1.22+
- libpcap (`sudo apt install libpcap-dev`)
- Root or `CAP_NET_RAW` capability for packet capture

## Intended Use

These tools are built for defensive monitoring of systems you own and operate. Do not deploy them on networks or systems you do not have explicit permission to monitor. The author assumes no liability for misuse.

For testing, pair with [hartkit](https://github.com/hart-dev-sudo/hartkit) against your own lab environment.

## License

MIT
