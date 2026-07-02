# Data lake

Query the Singularity Data Lake (SDL) with `s1ctl datalake powerquery`.
Supports two protocols: GraphQL (default) and REST.

## Prerequisites

- s1ctl [installed](guides/install.md) and [configured](guides/configure.md)
- `S1_CONSOLE_URL` and `S1_TOKEN` set (env, config file, or flags)
- For REST protocol only: `S1_SDL_URL` set to the XDR data lake host

## Command reference

```text
s1ctl datalake powerquery [flags]
```

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--query` | `string` | | PowerQuery expression (required) |
| `--start` | `string` | `24h` | Start time (`24h`, `7d`, `30d`, etc.) |
| `--end` | `string` | | End time (defaults to now) |
| `--protocol` | `string` | `graphql` | API protocol (`graphql`, `rest`) |
| `--priority` | `string` | `low` | Query priority (`low`, `high`) [REST only] |
| `--output` | `string` | `table` | Output format (`table`, `json`, `csv`) |
| `--json` | `bool` | false | Shorthand for `--output json` |
| `--no-progress` | `bool` | false | Disable spinner (for scripting) |

## Protocols

### GraphQL (default)

Connects through the management console URL (`S1_CONSOLE_URL`). No
additional configuration needed.

```bash
s1ctl datalake powerquery --query "endpoint.name contains 'srv'"
```

Use GraphQL when:

- You only have a console URL and token
- You want the simplest setup
- You are running interactive queries

### REST

Connects directly to the XDR data lake host (`S1_SDL_URL`). Requires
separate SDL credentials.

```bash
s1ctl datalake powerquery --query "endpoint.name contains 'srv'" --protocol rest
```

Use REST when:

- You need query priority control (`--priority high`)
- Your environment requires direct SDL access
- GraphQL is rate-limited or unavailable

## PowerQuery syntax

PowerQuery is SentinelOne's query language for the data lake. It is
not YARA-L or KQL.

### Basic structure

```text
<source> | <filter> | <aggregation>
```

### Common patterns

| Pattern | Example |
| --- | --- |
| Filter by field | `endpoint.name contains 'srv'` |
| Exact match | `event.type = 'Process Creation'` |
| Wildcard | `src.process.name matches '*chrome*'` |
| Multiple conditions | `endpoint.name contains 'srv' AND event.type = 'DNS'` |
| Aggregation | `endpoint.name contains 'srv' \| count by endpoint.name` |

## Time ranges

The `--start` flag accepts relative durations:

| Value | Meaning |
| --- | --- |
| `1h` | Last hour |
| `24h` | Last 24 hours (default) |
| `7d` | Last 7 days |
| `30d` | Last 30 days |

Combine with `--end` to query a specific window:

```bash
s1ctl datalake powerquery --query "endpoint.name contains 'srv'" --start 7d --end 1d
```

## Progress spinner

Long-running queries display a spinner on stderr. The spinner is
automatically suppressed when output is not a TTY (e.g., piped to a file).

Force disable for scripting:

```bash
s1ctl datalake powerquery --query "event.type = 'DNS'" --no-progress --json
```

## Workflows

### Search by endpoint name

```bash
s1ctl datalake powerquery --query "endpoint.name contains 'prod-web'"
```

### Search process events

```bash
s1ctl datalake powerquery --query "event.type = 'Process Creation' AND src.process.name = 'powershell.exe'" --start 7d
```

### Search DNS events

```bash
s1ctl datalake powerquery --query "event.type = 'DNS' AND event.dns.request contains 'example.com'" --start 24h
```

### Export results to JSON

```bash
s1ctl datalake powerquery --query "endpoint.name contains 'srv'" --json > results.json
```

### Export as CSV

```bash
s1ctl datalake powerquery --query "endpoint.name contains 'srv'" --output csv > results.csv
```

### Pipe to jq for processing

```bash
s1ctl datalake powerquery --query "event.type = 'DNS'" --json \
  | jq '[.Values[] | {endpoint: .[0], query: .[1]}]'
```

### Scripting and automation

Disable the spinner and use JSON output for reliable parsing:

```bash
#!/bin/bash
results=$(s1ctl datalake powerquery \
  --query "event.type = 'Process Creation'" \
  --start 24h \
  --no-progress \
  --json)

count=$(echo "$results" | jq '.Values | length')
echo "Found $count events"
```

### High-priority query (REST)

Use `--priority high` for time-sensitive investigations:

```bash
s1ctl datalake powerquery \
  --query "event.type = 'Process Creation' AND src.process.name = 'mimikatz.exe'" \
  --protocol rest \
  --priority high \
  --start 7d
```

## Facet aggregation

`datalake facet` returns the most common values of a single field (SDL REST
protocol, so `S1_SDL_URL` must be set). Pass `--output csv` for a
spreadsheet-friendly table.

```bash
s1ctl datalake facet --field event.type --start 24h
s1ctl datalake facet --field src.process.name \
  --filter "event.type = 'Process Creation'" --start 7d --max-count 20
```

| Flag | Description |
| --- | --- |
| `--field` | Field to aggregate (required) |
| `--start` | Start time, e.g. `24h` or a timestamp (required) |
| `--end` | End time (defaults to now) |
| `--filter` | Query filter expression |
| `--max-count` | Max distinct values to return |

## Time-series aggregation

`datalake timeseries` buckets an aggregation over time (SDL REST protocol).
Pass `--output csv` to emit one row per bucket for spreadsheets.

```bash
s1ctl datalake timeseries --filter "event.type = 'DNS'" \
  --function count --start 24h --buckets 24
s1ctl datalake timeseries --filter "event.type = 'Process Creation'" \
  --function "mean(duration)" --start 7d
```

| Flag | Description |
| --- | --- |
| `--filter` | Query filter expression (required) |
| `--start` | Start time, e.g. `24h` or a timestamp (required) |
| `--end` | End time (defaults to now) |
| `--function` | Aggregation function, e.g. `count`, `mean(field)` |
| `--buckets` | Number of time buckets |

## Ingest

Send data into the data lake. Both subcommands are **dry-run by default**;
pass `--yes` to apply.

### Structured events

`datalake ingest events` uploads a JSON array of events (addEvents):

```bash
s1ctl datalake ingest events --file events.json --session my-uploader-1 --yes
```

| Flag | Description |
| --- | --- |
| `--file` | JSON file containing an array of events (required) |
| `--session` | Unique session ID for this uploader (required) |

### Plain-text logs

`datalake ingest logs` uploads a raw log file (uploadLogs):

```bash
s1ctl datalake ingest logs --file app.log --parser myparser --yes
```

| Flag | Description |
| --- | --- |
| `--file` | Log file to upload (required) |
| `--parser` | Parser to apply on ingest |
| `--logfile` | `logfile` attribute for the events |
| `--server-host` | `serverHost` attribute for the events |

## Configuration files

`datalake files` manages data lake configuration files such as parsers and
lookups.

List and fetch (`files list` honors `--output csv`):

```bash
s1ctl datalake files list
s1ctl datalake files list --output csv
s1ctl datalake files get path/to/file.conf
s1ctl datalake files get path/to/file.conf --out local-copy.conf
```

Create, update, or delete (dry-run by default; pass `--yes` to apply):

```bash
s1ctl datalake files put path/to/file.conf --from-file local.conf --yes
s1ctl datalake files put path/to/file.conf --delete --yes
```

| Flag | Description |
| --- | --- |
| `--from-file` | Local file with the new content |
| `--delete` | Delete the remote file |
| `--out` | (get) Write content to a local file instead of stdout |
| `--expected-version` | (put) Fail if the remote version differs |

## Go SDK

Query the data lake programmatically:

```go
import "danny.vn/s1/sdl"

client := sdl.NewClient("https://your-console.sentinelone.net", token)

resp, err := client.PowerQueryGraphQL(ctx, &sdl.PowerQueryRequest{
    Query:     "endpoint.name contains 'srv'",
    StartTime: "24h",
})

for _, row := range resp.Values {
    fmt.Println(row)
}
```

REST protocol:

```go
client := sdl.NewClient("https://your-xdr-host.sentinelone.net", token)

resp, err := client.PowerQuery(ctx, &sdl.PowerQueryRequest{
    Query:     "endpoint.name contains 'srv'",
    StartTime: "24h",
    Priority:  "high",
})
```
