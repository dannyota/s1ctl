# filters

Manage saved endpoint filters

## filters create

Create a saved filter from a JSON file

```text
s1ctl filters create --from-file <filter.json> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | create in these account IDs |
| `--from-file` | string | - | filter definition JSON file (required) |
| `--site-id` | stringSlice | - | create in these site IDs (default: global/tenant) |
| `--yes` | bool | false | apply the action (default: dry-run) |

## filters delete

Delete a saved filter

```text
s1ctl filters delete <filter-id> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply the action (default: dry-run) |

## filters list

List saved filters

```text
s1ctl filters list [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | filter by account ID |
| `--all` | bool | false | fetch all pages |
| `--cursor` | string | - | pagination cursor |
| `--limit` | int | 0 | max results per page (default 50) |
| `--query` | string | - | free text search on filter name |
| `--site-id` | stringSlice | - | filter by site ID |

## filters update

Update a saved filter from a JSON file

```text
s1ctl filters update <filter-id> --from-file <filter.json> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--from-file` | string | - | filter definition JSON file (required) |
| `--yes` | bool | false | apply the action (default: dry-run) |
