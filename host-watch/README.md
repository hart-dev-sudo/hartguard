# host-watch

A configurable host health monitor for Linux servers. Checks disk usage, memory, container status, and service reachability in one pass — continuous or on-demand.

Part of [hartguard](../) — blue team security tooling for Linux servers.

## What it checks

- **Disk** — usage percentage per configured path, warns at threshold
- **Memory** — used / total in MB
- **Containers** — running state for each configured container
- **Services** — HTTP reachability for configured URLs

## Features

- Fully config-driven — no hardcoded container names or paths
- Continuous monitoring or one-shot mode (`--once`)
- Structured JSON log output for warnings and failures
- Clean color-coded terminal output

## Requirements

- Go 1.22+
- Docker (for container checks)

## Installation

```bash
git clone https://github.com/hart-dev-sudo/hartguard.git
cd hartguard/host-watch
go build -o host-watch .
```

## Configuration

Copy the example config and edit for your system:

```bash
cp config.yaml.example config.yaml
```

```yaml
disk_paths:
  - /                         # always check root
  - /mnt/data                 # add your mount points
disk_warn_percent: 80
containers:
  - my-container              # Docker container names to check
service_urls:
  - name: My Service
    url: http://localhost:8080/health
interval: 300
log_file: logs/host-watch.log
```

**Finding your mount points:** Run `df -h` to list mounted filesystems. **Finding container names:** Run `docker ps --format '{{.Names}}'`.

## Usage

```bash
# Run continuously
./host-watch

# Single check and exit
./host-watch --once

# Custom config path
./host-watch --config /etc/hartguard/host-watch.yaml
```

## Example Output

```
[2024-01-15 14:30:00] Host check
========================================
-- Disk --
[OK  ] /: 54% used
[WARN] /mnt/media: 83% used

-- Memory --
[INFO] 6210 MB used / 15987 MB total (38%)

-- Containers --
[OK  ] nginx
[CRIT] postgres — not running

-- Services --
[OK  ] My App
========================================
```

### Production install (systemd)

```bash
sudo make install
sudo systemctl enable --now host-watch
sudo journalctl -u host-watch -f
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
