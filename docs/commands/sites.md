# sites

Manage sites

## sites count

Count sites

```text
s1ctl sites count [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | filter by account ID |

## sites create

Create a site

```text
s1ctl sites create [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | string | - | account ID (required) |
| `--description` | string | - | site description |
| `--expiration` | string | - | expiration timestamp (RFC 3339) |
| `--name` | string | - | site name (required) |
| `--site-type` | string | - | site type |
| `--total-licenses` | int | 0 | total licenses |
| `--unlimited-licenses` | bool | false | unlimited licenses |
| `--yes` | bool | false | apply the action (default: dry-run) |

## sites delete

Delete a site

```text
s1ctl sites delete <site-id> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply the action (default: dry-run) |

## sites get

Get site details

```text
s1ctl sites get <site-id>
```

## sites licenses

Show license utilization across sites

```text
s1ctl sites licenses
```

Aggregate license health view across all sites.
Each site shows active vs total licenses, utilization percentage,
expiration date, and a status indicator (OK, WARNING, CRITICAL).

## sites list

List sites

```text
s1ctl sites list [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | filter by account ID |
| `--all` | bool | false | fetch all pages |
| `--cursor` | string | - | pagination cursor |
| `--limit` | int | 0 | max results per page (default 50) |
| `--query` | string | - | free text search |
| `--sort-by` | string | - | sort field (e.g. name, state) |
| `--sort-order` | string | - | sort direction (asc, desc) |
| `--state` | stringSlice | - | filter by state |

## sites pull

Pull sites to local YAML files

```text
s1ctl sites pull [flags]
```

Fetch all sites and write them as YAML files.

Each site produces one file named by its sanitized name. Server-only metadata
(ID, state, licenses in use, timestamps) is omitted so the files contain only
the declarative site definition.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | filter by account ID |
| `--out` | string | sites | output directory |

## sites push

Push sites from local YAML files

```text
s1ctl sites push [flags]
```

Read site YAML files from a directory and sync them to SentinelOne.

Sites are matched by name: existing sites are updated, new sites are created,
and unchanged sites are skipped. Dry-run by default — pass --yes to apply changes.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--dir` | string | sites | directory containing site YAML files |
| `--yes` | bool | false | apply changes (default: dry-run) |

## sites update

Update a site

```text
s1ctl sites update <site-id> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--description` | string | - | new description |
| `--expiration` | string | - | new expiration timestamp (RFC 3339) |
| `--name` | string | - | new site name |
| `--total-licenses` | int | 0 | new total licenses |
| `--unlimited-licenses` | bool | false | unlimited licenses |
| `--yes` | bool | false | apply the action (default: dry-run) |
