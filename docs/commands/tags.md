# tags

Manage tags

## tags create

Create a tag

```text
s1ctl tags create [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--description` | string | - | tag description |
| `--key` | string | - | tag key (required) |
| `--scope` | string | - | tag scope |
| `--scope-id` | string | - | tag scope ID |
| `--value` | string | - | tag value (required) |
| `--yes` | bool | false | apply the action (default: dry-run) |

## tags delete

Delete a tag

```text
s1ctl tags delete <tag-id> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply the action (default: dry-run) |

## tags get

Get a tag

```text
s1ctl tags get <tag-id>
```

## tags list

List tags

```text
s1ctl tags list [flags]
```

List tags by type: firewall, network-quarantine, device-inventory.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--all` | bool | false | fetch all pages |
| `--cursor` | string | - | pagination cursor |
| `--limit` | int | 0 | max results per page (default 50) |
| `--query` | string | - | free text search |
| `--site-id` | stringSlice | - | filter by site ID |
| `--type` | string | - | tag type (firewall, network-quarantine, device-inventory) |

## tags pull

Pull tags to local YAML files

```text
s1ctl tags pull [flags]
```

Fetch all tags and write them as YAML files.

Each tag produces one file named by its sanitized key. Server-only metadata
(ID, timestamps) is omitted so the files contain only the declarative
definition. Tags are matched by key: duplicate keys across scopes resolve to
the first one listed with a warning.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--out` | string | tags | output directory |
| `--site-id` | stringSlice | - | filter by site ID |

## tags push

Push tags from local YAML files

```text
s1ctl tags push [flags]
```

Read tag YAML files from a directory and sync them to SentinelOne.

Tags are matched by key: existing tags are updated, new tags are created,
and unchanged tags are skipped. Dry-run by default — pass --yes to apply changes.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--dir` | string | tags | directory containing tag YAML files |
| `--yes` | bool | false | apply changes (default: dry-run) |

## tags update

Update a tag

```text
s1ctl tags update <tag-id> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--description` | string | - | new description |
| `--key` | string | - | new tag key |
| `--value` | string | - | new tag value |
| `--yes` | bool | false | apply the action (default: dry-run) |
