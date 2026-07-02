# blocklist

Manage the blocklist (blocked file hashes)

## blocklist create

Add a hash to the blocklist

```text
s1ctl blocklist create [flags]
```

Add a SHA1 (--value) and/or SHA256 (--sha256) hash to the blocklist.

OS types: windows, linux, macos, windows_legacy
Type must be black_hash (any other value creates an exclusion instead).

New items are added to the scope given by --site-id/--group-id/--account-id, or
to the global (tenant) blocklist when no scope flag is set.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | target account IDs |
| `--description` | string | - | blocklist item description |
| `--group-id` | stringSlice | - | target group IDs |
| `--os-type` | string | - | target OS (windows, linux, macos, windows_legacy) (required) |
| `--sha256` | string | - | SHA256 hash to block |
| `--site-id` | stringSlice | - | target site IDs |
| `--source` | string | - | blocklist item source |
| `--type` | string | black_hash | restriction type (black_hash) |
| `--value` | string | - | SHA1 hash to block (required) |
| `--yes` | bool | false | apply the action (default: dry-run) |

## blocklist delete

Delete a blocklist item

```text
s1ctl blocklist delete <blocklist-id> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply the action (default: dry-run) |

## blocklist export

Export blocklist items as CSV

```text
s1ctl blocklist export [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | filter by account ID |
| `--group-id` | stringSlice | - | filter by group ID |
| `--os-type` | stringSlice | - | filter by OS type |
| `--out` | string | - | write export to file (default: stdout) |
| `--site-id` | stringSlice | - | filter by site ID |
| `--tenant` | bool | false | export the global (tenant) blocklist |

## blocklist list

List blocklist items

```text
s1ctl blocklist list [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | filter by account ID |
| `--all` | bool | false | fetch all pages |
| `--cursor` | string | - | pagination cursor |
| `--group-id` | stringSlice | - | filter by group ID |
| `--limit` | int | 0 | max results per page (default 50) |
| `--os-type` | stringSlice | - | filter by OS type (windows, linux, macos, windows_legacy) |
| `--query` | string | - | free text search |
| `--site-id` | stringSlice | - | filter by site ID |
| `--sort-by` | string | - | sort field (e.g. createdAt, osType) |
| `--sort-order` | string | - | sort direction (asc, desc) |
| `--value` | string | - | filter by hash value |

## blocklist pull

Pull blocklist items to local YAML files

```text
s1ctl blocklist pull [flags]
```

Fetch all blocklist items and write them as YAML files.

Each item produces one file. Server-only metadata (ID, scope, source,
timestamps) is omitted so the files contain only the declarative definition.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--out` | string | blocklist | output directory |
| `--site-id` | stringSlice | - | filter by site ID |

## blocklist push

Push blocklist items from local YAML files

```text
s1ctl blocklist push [flags]
```

Read blocklist YAML files from a directory and sync them to SentinelOne.

Items are matched by type + OS + value: existing items are updated, new items
are created, and unchanged items are skipped. Dry-run by default — pass --yes
to apply changes.

New items are created at the scope specified by --site-id. If no scope flag is
given, they are created at the global (tenant) scope.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--dir` | string | blocklist | directory containing blocklist YAML files |
| `--site-id` | stringSlice | - | scope for new items (default: global/tenant) |
| `--yes` | bool | false | apply changes (default: dry-run) |

## blocklist update

Update a blocklist item (full replacement)

```text
s1ctl blocklist update <blocklist-id> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--description` | string | - | blocklist item description |
| `--os-type` | string | - | target OS (required) |
| `--sha256` | string | - | SHA256 hash |
| `--source` | string | - | blocklist item source |
| `--type` | string | black_hash | restriction type (black_hash) |
| `--value` | string | - | SHA1 hash (required) |
| `--yes` | bool | false | apply the action (default: dry-run) |

## blocklist validate

Check whether a hash is Not Allowed or Not Recommended

```text
s1ctl blocklist validate [flags]
```

Check whether a hash is on SentinelOne's "Not Allowed" or "Not Recommended"
list before adding it to the blocklist. This is a read-only check.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | scope account IDs |
| `--group-id` | stringSlice | - | scope group IDs |
| `--os-type` | string | - | target OS (windows, linux, macos, windows_legacy) |
| `--sha256` | string | - | SHA256 hash to validate |
| `--site-id` | stringSlice | - | scope site IDs |
| `--value` | string | - | SHA1 hash to validate |
