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

Edit `config.yaml`:

```yaml
disk_paths:
  - /
  - /mnt/data
disk_warn_percent: 80
containers:
  - nginx
  - postgres
service_urls:
  - name: My App
    url: http://localhost:8080/health
interval: 300        # seconds between checks (0 = one-shot)
log_file: logs/host-watch.log
```

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

## Running Tests

```bash
go test ./internal/checker/...
```

## License

MIT
