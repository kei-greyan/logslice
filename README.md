# logslice

Stream and filter structured JSON logs from multiple sources with a unified query syntax.

## Installation

```bash
go install github.com/yourusername/logslice@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/logslice.git && cd logslice && go build ./...
```

## Usage

```bash
# Filter logs by level from a local file
logslice --source ./app.log --filter "level=error"

# Stream logs from multiple sources and filter by field
logslice --source ./api.log --source ./worker.log --filter "service=auth"

# Pipe from stdin
cat app.log | logslice --filter "status>=500"

# Stream from a remote source with a time range
logslice --source http://logs.example.com/stream --filter "level=warn" --since 1h
```

### Query Syntax

| Operator | Example | Description |
|----------|---------|-------------|
| `=` | `level=error` | Exact match |
| `!=` | `level!=debug` | Not equal |
| `>=` | `status>=500` | Greater than or equal |
| `~` | `msg~timeout` | Contains substring |

Output is newline-delimited JSON, making it easy to pipe into tools like `jq`:

```bash
logslice --source ./app.log --filter "level=error" | jq '.message'
```

## Configuration

logslice can be configured via a `.logslice.yaml` file in the project root or your home directory:

```yaml
sources:
  - ./logs/api.log
  - ./logs/worker.log
default_filter: "level!=debug"
```

## License

MIT © [yourusername](https://github.com/yourusername)