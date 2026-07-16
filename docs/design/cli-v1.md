# CLI v1 — read, export, page, sort, progress, errors

Output, pagination, sorting, progress, and error conventions shared by every
read command. Designed in the v1 wave; all of it is implemented.

## Output formats

A global `--output` flag selects the format:

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
| `--no-progress` | Global | Disable spinners and progress (for scripts/AI agents) |

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

A TTY-aware spinner (charmbracelet):

| Operation | Display |
|-----------|---------|
| SDL PowerQuery (poll loop) | `Running query...` with elapsed time |
| `--all` pagination | `Fetching <resource>... N/total` (or `N fetched` if total unknown) |

Spinners render on stderr so stdout stays clean for piping. Suppressed
entirely when stderr is not a TTY or when `--no-progress` is passed.

`--no-progress` is useful when AI agents or scripts consume s1ctl output and
need clean, predictable output without ANSI escape sequences.

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

## Implementation map

| File | Role |
|------|------|
| `internal/cli/paging.go` | `listAll()` auto-paginator with progress |
| `internal/cli/progress.go` | TTY-aware spinner |
| `internal/cli/root.go` | `--verbose`, `--output` flags; error formatting in `Execute()` |
| `internal/cli/output.go` | `printCSV()`, error formatting, `--output` handling |
| `internal/cli/*_list.go` | `--all`, `--cursor`/`--after`, `--sort-by`, `--sort-order` where applicable |

Mutations, config-as-code, and the wider command surface were added in later
waves — see the [Roadmap](https://github.com/dannyota/s1ctl/blob/master/ROADMAP.md)
and [Catalog](catalog.md).
