# visibility

Deep Visibility threat hunting

## visibility query

Run a Deep Visibility query

```text
s1ctl visibility query [flags]
```

Run a Deep Visibility query to hunt for endpoint events.

Initiates a query, polls until complete, then fetches and displays results.
The query uses SentinelOne's Deep Visibility query language.

Examples:
  s1ctl visibility query --query "EventType = \"Process Creation\""
  s1ctl visibility query --query "ProcessName contains \"cmd.exe\"" --from 7d
  s1ctl visibility query --query "SHA256 = \"abc123...\"" --json

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--from` | string | 24h | start time (duration like 24h/7d, or RFC3339) |
| `--limit` | int | 100 | max events per page (1-1000) |
| `--max-results` | int | 0 | stop after fetching this many events (0 = all) |
| `--poll-interval` | duration | 2s | interval between status polls |
| `--query` | string | - | Deep Visibility query expression (required) |
| `--site-id` | stringSlice | - | filter by site ID |
| `--sort-by` | string | createdAt | sort field (e.g. createdAt, pid) |
| `--sort-order` | string | desc | sort direction (asc, desc) |
| `--to` | string | - | end time (default: now) |
