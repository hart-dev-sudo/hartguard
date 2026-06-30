# hartguard

A personal blue team security suite built incrementally — one layover at a time.

## What is this?

hartguard is a growing collection of defensive security tools built from scratch in Go. Each tool is written to understand the underlying detection techniques, not just wrap existing utilities. Counterpart to [hartkit](https://github.com/hart-dev-sudo/hartkit) — if hartkit scans, hartguard listens.

## Tools

| Tool | Description | Status |
|------|-------------|--------|
| [port-scan-detector](./port-scan-detector/) | Detects incoming TCP port scans using a sliding window algorithm — SYN, FIN, NULL, and XMAS | Available |
| [vpn-watch](./vpn-watch/) | Monitors containerized VPN setups for IP leaks — verified against any IP-echo endpoint | Available |
| [host-watch](./host-watch/) | Checks disk, memory, container status, and service reachability in one pass | Available |
| [login-watch](./login-watch/) | Detects SSH brute force attempts in real time — username extraction, sliding window, systemd-ready | Available |

## Structure

```
hartguard/
├── port-scan-detector/    # detects incoming port scans
├── vpn-watch/             # VPN leak detection
├── host-watch/            # host health monitoring
└── login-watch/           # SSH brute force detection
```

## Usage

Each tool builds and runs standalone:

```bash
cd <tool>
cp config.yaml.example config.yaml   # edit for your system
go build -o <tool> .
sudo ./<tool>
```

See each tool's README for full usage and configuration.

## Requirements

- Go 1.22+
- See each tool's README for additional dependencies

## Intended Use

These tools are built for defensive monitoring of systems you own and operate. Do not deploy them on networks or systems you do not have explicit permission to monitor. Unauthorized monitoring may violate the Computer Fraud and Abuse Act (CFAA) and equivalent laws in your jurisdiction.

For testing, pair with [hartkit](https://github.com/hart-dev-sudo/hartkit) against your own lab environment.

The author assumes no liability for misuse of these tools.

## License

MIT
