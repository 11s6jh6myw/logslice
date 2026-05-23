# logslice

Fast log file segmenter that splits large logs by time range or pattern.

## Installation

```bash
go install github.com/yourusername/logslice@latest
```

## Usage

Split a log file by time range:

```bash
logslice --input app.log --from "2024-01-15 08:00:00" --to "2024-01-15 12:00:00" --output morning.log
```

Split by pattern:

```bash
logslice --input app.log --pattern "ERROR|FATAL" --output errors.log
```

Slice into fixed time windows:

```bash
logslice --input app.log --window 1h --output-dir ./slices/
```

### Flags

| Flag | Description |
|------|-------------|
| `--input` | Path to the source log file |
| `--output` | Path to the output file |
| `--output-dir` | Directory for windowed output files |
| `--from` | Start timestamp (RFC3339 or common log formats) |
| `--to` | End timestamp |
| `--pattern` | Regex pattern to match log lines |
| `--window` | Split into time windows (e.g. `1h`, `30m`) |
| `--format` | Log timestamp format (default: auto-detect) |

## Install from Source

```bash
git clone https://github.com/yourusername/logslice.git
cd logslice
go build -o logslice .
```

## License

MIT