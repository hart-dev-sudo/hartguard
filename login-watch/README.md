# login-watch

A real-time SSH brute force detector for Linux servers. Tails `/var/log/auth.log`, tracks failed login attempts per source IP using a sliding window, and alerts when a threshold is hit.

Part of [hartguard](../) — blue team security tooling for Linux servers. Pairs with [port-scan-detector](../port-scan-detector/) to cover the full recon → brute force attack chain.

## What it detects

- `Failed password` — wrong password attempts
- `Invalid user` — attempts with unknown usernames
- `Connection closed by invalid user` — scan-pattern connections
- `authentication failure` — PAM-level failures
- `BREAK-IN ATTEMPT` — kernel-flagged intrusion attempts

## Features

- Real-time log tailing — no polling delay
- Sliding window algorithm — configurable threshold and time window
- Username extraction — alerts show who was targeted, not just the source IP
- Heartbeat output in watch mode so you know it's alive
- IP whitelist support
- `--scan` mode to analyze existing log history
- Structured JSON log output
- Systemd service unit file for production deployment

## Requirements

- Go 1.22+
- Read access to `/var/log/auth.log` (requires root or `adm` group)

## Installation

```bash
git clone https://github.com/hart-dev-sudo/hartguard.git
cd hartguard/login-watch
go build -o login-watch .
```

### Production install (systemd)

```bash
sudo make install
sudo systemctl enable --now login-watch
sudo journalctl -u login-watch -f
```

## Configuration

Copy the example config and edit for your system:

```bash
cp config.yaml.example config.yaml
```

```yaml
auth_log: /var/log/auth.log   # Ubuntu/Debian; Fedora/RHEL use /var/log/secure
threshold: 5                  # failures within window to trigger alert
window: 60                    # sliding window size in seconds
whitelist:
  - 127.0.0.1
log_file: logs/login-watch.log
```

**Auth log path by distro:**

| Distro | Path |
|--------|------|
| Ubuntu / Debian | `/var/log/auth.log` |
| Fedora / RHEL / CentOS | `/var/log/secure` |
| Arch Linux | `/var/log/auth.log` or use `journalctl` |

After `make install`, config lives at `/etc/hartguard/login-watch.yaml`.

## Usage

```bash
# Watch live (requires log read access)
sudo ./login-watch

# Scan existing log history and exit
sudo ./login-watch --scan

# Custom config
sudo ./login-watch --config /etc/hartguard/login-watch.yaml
```

## Example Output

**Scan mode:**
```
========================================
  login-watch — scan mode
  log:       /var/log/auth.log
  threshold: 5 failures in 60s
========================================
[*]     18 failure events processed
[ALERT] Brute force from 203.0.113.42 (user: root) — 7 failures in 60s (Failed password)
[ALERT] Brute force from 198.51.100.9 (user: admin) — 5 failures in 60s (Invalid user)
========================================
  scan complete
========================================
```

**Watch mode (live):**
```
========================================
  login-watch — live mode
  log:       /var/log/auth.log
  threshold: 5 failures in 60s
========================================
[*]     still watching — 14:32:01
[*]     still watching — 14:33:01
[ALERT] Brute force from 203.0.113.42 (user: root) — 5 failures in 60s (Failed password)
```

**Log entry (JSON):**
```json
{"timestamp":"2024-01-15T14:23:01Z","src_ip":"203.0.113.42","username":"root","event_type":"Failed password","count":5,"window_secs":60}
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
