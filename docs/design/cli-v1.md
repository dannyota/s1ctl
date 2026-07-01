# CLI v1 — read, export, page, sort, progress, errors

First usable release of the CLI. Every read command gets consistent output
formats, pagination, sorting, progress feedback, and actionable error messages.

## Output formats

A global `--output` flag replaces the current `--json` boolean:

| Format | Description |
|--------|-------------|
| `table` | Default. Lipgloss table with headers, styled for TTY |
| `json` | Pretty-printed JSON array (or object for `get`) |
| `csv` | RFC 4180 CSV with header row, pipeable |

`--json` is kept as shorthand for `--output json`.

## Pagination

### Defaults

| Protocol | Default page size |
|----------|-------------------|
| REST | 50 |
| GraphQL | 50 |

### Flags

| Flag | Scope | Description |
|------|-------|-------------|
| `--limit N` | All list commands | Page size (overrides default) |
| `--all` | All list commands | Auto-paginate until exhausted |
| `--cursor` | REST list commands | Resume from a cursor value |
| `--after` | GraphQL list commands | Resume from an end cursor |

### Footer

Table output prints a footer after results:

```text
Showing 50 of 1,234 agents (use --all to fetch all)
```

When `--all` is used, the footer shows the total fetched:

```text
1,234 agents
```

JSON and CSV output suppress the footer — the data speaks for itself.

### Auto-pagination (`--all`)

Loops using cursor (REST) or endCursor (GraphQL) until no more pages. Shows
a progress line on TTY:

```text
Fetching agents... 200/1,234
```

Progress is suppressed when stdout is not a TTY (piped to file or another
command).

## Sorting

`--sort-by` and `--sort-order` flags on commands whose SDK params support
sorting:

| Command | Sortable fields (examples) |
|---------|---------------------------|
| `agents list` | computerName, lastActiveDate, osType |
| `threats list` | createdAt, mitigationStatus |
| `sites list` | name, state, expiration |
| `groups list` | name, type |
| `exclusions list` | type, value |
| `users list` | fullName, email, dateJoined |

`--sort-order` accepts `asc` (default) or `desc`.

Commands without server-side sort support (activities, applications, device
control, firewall, remote ops, updates, tags) omit these flags.

GraphQL commands (alerts, vulnerabilities, misconfigurations, cloud policies)
use server-side filtering but do not expose sort flags — the GraphQL schemas
do not support sort parameters.

## Progress indicators

A TTY-aware spinner using bubbletea (already a dependency):

| Operation | Display |
|-----------|---------|
| SDL PowerQuery (poll loop) | `Running query...` with elapsed time |
| `--all` pagination | `Fetching <resource>... N/total` (or `N fetched` if total unknown) |

Spinners render on stderr so stdout stays clean for piping. Suppressed
entirely when stderr is not a TTY.

## Error handling

### Default

Short, actionable message:

```text
Error: HTTP 401: Unauthorized
```

### `--verbose`

Full API error with response body:

```text
Error: HTTP 401: Unauthorized

  Title:  Unauthorized
  Detail: The API token is invalid or expired.
  Body:   {"errors":[{"code":4010010,"detail":"...","title":"Unauthorized"}]}
```

### `--json` errors

When `--output json` (or `--json`) is active, errors are also JSON:

```json
{
  "error": {
    "status": 401,
    "title": "Unauthorized",
    "detail": "The API token is invalid or expired.",
    "body": "{\"errors\":[...]}"
  }
}
```

This lets scripts parse errors programmatically and lets users paste the full
error when reporting issues.

### `--verbose` on non-API errors

For non-API errors (network timeouts, config issues, etc.), `--verbose` shows
the full Go error chain. Default shows only the leaf message.

## Implementation scope

### New files

| File | Purpose |
|------|---------|
| `internal/cli/paging.go` | `listAll()` auto-paginator with progress |
| `internal/cli/progress.go` | TTY-aware spinner wrapping bubbletea |

### Modified files

| File | Changes |
|------|---------|
| `root.go` | Add `--verbose`, `--output` flags; update error formatting in `Execute()` |
| `output.go` | Add `printCSV()`, `formatError()`, respect `--output` flag |
| All `*_list.go` | Add `--all`, `--cursor`/`--after`, `--sort-by`, `--sort-order` where applicable; update footer |
| `datalake_query.go` | Add spinner to powerquery poll loop |

### Not in scope

- Mutation commands (create, update, delete) — separate wave
- Config-as-code (pull/push) — separate wave
- New resource commands — only enhancing existing ones
