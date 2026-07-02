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
