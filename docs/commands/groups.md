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

## groups pull

Pull groups to local YAML files

```text
s1ctl groups pull [flags]
```

Fetch all groups and write them as YAML files.

Each group produces one file. Server-only metadata (ID, rank, agent counts,
timestamps) is omitted so the files contain only the declarative definition.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--out` | string | groups | output directory |
| `--site-id` | stringSlice | - | filter by site ID |

## groups push

Push groups from local YAML files

```text
s1ctl groups push [flags]
```

Read group YAML files from a directory and sync them to SentinelOne.

Groups are matched by site ID + name: existing groups are updated, new groups
are created, and unchanged groups are skipped. A group file without a siteId
fails at create time. Dry-run by default — pass --yes to apply changes.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--dir` | string | groups | directory containing group YAML files |
| `--yes` | bool | false | apply changes (default: dry-run) |

## groups update

Update a group

```text
s1ctl groups update <group-id> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--description` | string | - | new description |
| `--name` | string | - | new group name |
| `--yes` | bool | false | apply the action (default: dry-run) |
