# logdrift

A lightweight CLI tool for tailing and filtering structured JSON logs across multiple services simultaneously.

---

## Installation

```bash
go install github.com/yourusername/logdrift@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/logdrift.git
cd logdrift
go build -o logdrift .
```

---

## Usage

Tail logs from multiple services and filter by log level or field values:

```bash
# Tail logs from multiple files
logdrift tail -f /var/log/api.log /var/log/worker.log

# Filter by log level
logdrift tail -f /var/log/api.log --level error

# Filter by a specific JSON field value
logdrift tail -f /var/log/*.log --filter service=payments

# Output pretty-printed JSON
logdrift tail -f /var/log/api.log --pretty
```

### Example Output

```
[api]     2024-01-15T10:23:01Z  ERROR  "message": "connection timeout"  service=api
[worker]  2024-01-15T10:23:02Z  INFO   "message": "job completed"       service=worker
```

---

## Flags

| Flag | Description |
|------|-------------|
| `-f` | Files to tail (supports glob patterns) |
| `--level` | Filter by log level (debug, info, warn, error) |
| `--filter` | Filter by field value (key=value) |
| `--pretty` | Pretty-print JSON output |

---

## License

MIT © [yourusername](https://github.com/yourusername)