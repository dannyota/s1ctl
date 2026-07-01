# groups

Manage groups

## groups count

Count groups

```text
s1ctl groups count [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--site-id` | stringSlice | - | filter by site ID |

## groups create

Create a group

```text
s1ctl groups create [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--description` | string | - | group description |
| `--name` | string | - | group name (required) |
| `--site-id` | string | - | site ID (required) |
| `--yes` | bool | false | apply the action (default: dry-run) |

## groups delete

Delete a group

```text
s1ctl groups delete <group-id> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply the action (default: dry-run) |

## groups get

Get group details

```text
s1ctl groups get <group-id>
```

## groups list

List groups

```text
s1ctl groups list [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--all` | bool | false | fetch all pages |
| `--cursor` | string | - | pagination cursor |
| `--limit` | int | 0 | max results per page (default 50) |
| `--query` | string | - | free text search |
| `--site-id` | stringSlice | - | filter by site ID |
| `--sort-by` | string | - | sort field (e.g. name, type) |
| `--sort-order` | string | - | sort direction (asc, desc) |
