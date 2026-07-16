# marketplace

Manage Singularity Marketplace applications

## marketplace catalog

List marketplace catalog applications

```text
s1ctl marketplace catalog [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--all` | bool | false | fetch all pages |
| `--category` | string | - | filter by category (contains) |
| `--category-id` | stringSlice | - | filter by category ID |
| `--limit` | int | 0 | max results per page |
| `--name` | string | - | filter by name (contains) |
| `--query` | string | - | free-text search |
| `--sort-by` | string | - | sort field |
| `--sort-order` | string | - | sort order (asc, desc) |

## marketplace catalog-config

Show configuration fields for a catalog application

```text
s1ctl marketplace catalog-config CATALOG_ID
```

## marketplace config

Show configuration for an installed application

```text
s1ctl marketplace config APP_ID
```

## marketplace delete

Delete an installed marketplace application

```text
s1ctl marketplace delete [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | scope to account ID |
| `--group-id` | stringSlice | - | scope to group ID |
| `--id` | string | - | application ID (required) |
| `--site-id` | stringSlice | - | scope to site ID |
| `--tenant` | bool | false | scope to tenant |
| `--yes` | bool | false | apply the change (default: dry-run) |

## marketplace disable

Disable an installed marketplace application

```text
s1ctl marketplace disable APP_ID [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | scope to account ID |
| `--group-id` | stringSlice | - | scope to group ID |
| `--site-id` | stringSlice | - | scope to site ID |
| `--yes` | bool | false | apply the change (default: dry-run) |

## marketplace enable

Enable an installed marketplace application

```text
s1ctl marketplace enable APP_ID [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | scope to account ID |
| `--group-id` | stringSlice | - | scope to group ID |
| `--site-id` | stringSlice | - | scope to site ID |
| `--yes` | bool | false | apply the change (default: dry-run) |

## marketplace install

Install a marketplace application

```text
s1ctl marketplace install [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | scope to account ID |
| `--catalog-id` | string | - | catalog application ID (required) |
| `--config` | stringSlice | - | configuration (id=value, repeatable) |
| `--group-id` | stringSlice | - | scope to group ID |
| `--name` | string | - | instance name (required) |
| `--site-id` | stringSlice | - | scope to site ID |
| `--tenant` | bool | false | scope to tenant |
| `--yes` | bool | false | apply the change (default: dry-run) |

## marketplace list

List installed marketplace applications

```text
s1ctl marketplace list [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | filter by account ID |
| `--all` | bool | false | fetch all pages |
| `--catalog-id` | string | - | filter by catalog application ID |
| `--creator` | string | - | filter by creator (contains) |
| `--limit` | int | 0 | max results per page |
| `--name` | string | - | filter by name (contains) |
| `--query` | string | - | free-text search |
| `--site-id` | stringSlice | - | filter by site ID |
| `--sort-by` | string | - | sort field |
| `--sort-order` | string | - | sort order (asc, desc) |

## marketplace log

Show log entries for an installed application

```text
s1ctl marketplace log APP_ID [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--only-errors` | bool | false | show only error entries |

## marketplace update

Update an installed marketplace application

```text
s1ctl marketplace update [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | scope to account ID |
| `--config` | stringSlice | - | configuration (id=value, repeatable) |
| `--group-id` | stringSlice | - | scope to group ID |
| `--id` | string | - | application ID (required) |
| `--name` | string | - | new instance name |
| `--site-id` | stringSlice | - | scope to site ID |
| `--yes` | bool | false | apply the change (default: dry-run) |
