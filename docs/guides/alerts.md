# Alerts

Query SentinelOne unified alerts with `s1ctl alerts list`. This command
uses the **GraphQL UAM API** (Unified Alert Management), which provides
richer filtering and faster pagination than the REST alerts endpoint.

## Prerequisites

- s1ctl [installed](guides/install.md) and [configured](guides/configure.md)
- `S1_CONSOLE_URL` and `S1_TOKEN` set (env, config file, or flags)

## Command reference

```text
s1ctl alerts list [flags]
```

| Flag | Type | Default | Description |
| --- | --- | --- | --- |
| `--severity` | `[]string` | all | Filter by severity (`LOW`, `MEDIUM`, `HIGH`, `CRITICAL`) |
| `--verdict` | `[]string` | all | Filter by analyst verdict |
| `--limit` | `int` | 50 | Max results per page |
| `--all` | `bool` | false | Fetch all pages (auto-paginate) |
| `--after` | `string` | | GraphQL cursor for manual pagination |
| `--output` | `string` | `table` | Output format (`table`, `json`, `csv`) |
| `--json` | `bool` | false | Shorthand for `--output json` |

## Output columns

| Column | Description |
| --- | --- |
| ID | Alert identifier |
| Name | Alert rule or detection name (truncated to 40 chars) |
| Severity | `LOW`, `MEDIUM`, `HIGH`, or `CRITICAL` |
| Status | Alert status |
| Verdict | Analyst verdict (empty if unreviewed) |
| Detected | Detection timestamp |

## Filtering

### By severity

Pass one or more severity levels as a comma-separated list:

```bash
# Critical only
s1ctl alerts list --severity CRITICAL

# High and critical
s1ctl alerts list --severity HIGH,CRITICAL
```

### By analyst verdict

Filter alerts that analysts have already triaged:

```bash
s1ctl alerts list --verdict TRUE_POSITIVE
s1ctl alerts list --verdict FALSE_POSITIVE --limit 100
```

### Combined filters

Flags compose with AND logic:

```bash
s1ctl alerts list --severity CRITICAL --verdict TRUE_POSITIVE --limit 20
```

## Pagination

The alerts command uses **GraphQL cursor-based pagination**, not offset
pagination. Two approaches:

### Automatic (--all)

Fetch every matching alert across all pages:

```bash
s1ctl alerts list --severity HIGH,CRITICAL --all
```

The command auto-paginates through all results and prints a running count.

### Manual (--after)

Use the cursor from a previous response to fetch the next page:

```bash
# First page
s1ctl alerts list --limit 25

# Next page (use the cursor from the previous output)
s1ctl alerts list --limit 25 --after "YWxlcnQ6MDAwMDAw"
```

Note: alerts use `--after` (GraphQL cursor), not `--cursor` which is used
by REST-based commands.

## Workflows

### List critical alerts

Quick triage view of high-priority alerts:

```bash
s1ctl alerts list --severity CRITICAL --limit 20
```

### Export alerts to JSON

Export for SIEM integration, scripting, or archival:

```bash
# All critical and high alerts as JSON
s1ctl alerts list --severity HIGH,CRITICAL --all --json > alerts.json

# Pipe to jq for field extraction
s1ctl alerts list --all --json | jq '[.[] | {id: .ID, name: .Name, severity: .Severity}]'
```

### Export as CSV

```bash
s1ctl alerts list --severity HIGH,CRITICAL --all --output csv > alerts.csv
```

### Filter unreviewed alerts

Find alerts that have not been triaged:

```bash
s1ctl alerts list --severity HIGH,CRITICAL --limit 50
```

Review the Verdict column -- empty values indicate untriaged alerts.

### Count alerts by severity

```bash
s1ctl alerts list --severity CRITICAL --all --json | jq 'length'
```

## GraphQL vs REST

The `alerts list` command uses GraphQL UAM exclusively. GraphQL provides:

- Cursor-based pagination (stable across concurrent changes)
- Field-level filtering (severity, verdict)
- Lower latency for filtered queries

There is no `--protocol` flag on `alerts list` -- it always uses GraphQL.

## Go SDK

Query alerts programmatically with the `graphql` package:

```go
import "danny.vn/s1/graphql"

client := graphql.NewClient("https://your-console.sentinelone.net", token)

params := &graphql.ListParams{
    First: 25,
    Filters: []graphql.Filter{
        {
            FieldID:  "severity",
            StringIn: &graphql.InStr{Values: []string{"HIGH", "CRITICAL"}},
        },
    },
}

conn, err := client.AlertsList(ctx, params)
for _, edge := range conn.Edges {
    fmt.Println(edge.Node.ID, edge.Node.Severity, edge.Node.Name)
}
```
