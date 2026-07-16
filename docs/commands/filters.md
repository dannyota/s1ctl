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

## filters pull

Pull saved filters to local YAML files

```text
s1ctl filters pull [flags]
```

Fetch all saved filters and write them as YAML files.

Each filter produces one file named by its sanitized name. Server-only metadata
(ID, scope, timestamps) is omitted so the files contain only the declarative
definition: the filter name and its filterFields criteria set.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | filter by account ID |
| `--out` | string | filters | output directory |
| `--site-id` | stringSlice | - | filter by site ID |

## filters push

Push saved filters from local YAML files

```text
s1ctl filters push [flags]
```

Read filter YAML files from a directory and sync them to SentinelOne.

Filters are matched by name: existing filters are updated, new ones are created,
and unchanged ones are skipped. Dry-run by default — pass --yes to apply changes.
New filters are created at the scope given by --site-id (default: global/tenant).

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | scope for new filters |
| `--dir` | string | filters | directory containing filter YAML files |
| `--site-id` | stringSlice | - | scope for new filters (default: global/tenant) |
| `--yes` | bool | false | apply changes (default: dry-run) |

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
