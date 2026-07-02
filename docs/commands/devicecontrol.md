# devicecontrol

Device control rules

## devicecontrol copy

Copy device control rules between scopes

```text
s1ctl devicecontrol copy [flags]
```

Copy device control rules from a source scope to a target scope.

Use --source-site-id or --source-account-id to define the source, and
--target-site-id, --target-account-id, or --target-group-id for the destination.
At least one target flag is required.
Dry-run by default — pass --yes to apply.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--source-account-id` | stringSlice | - | source account IDs |
| `--source-site-id` | stringSlice | - | source site IDs |
| `--target-account-id` | string | - | target account ID |
| `--target-group-id` | stringSlice | - | target group IDs |
| `--target-site-id` | string | - | target site ID |
| `--yes` | bool | false | apply changes (default: dry-run) |

## devicecontrol delete

Delete device control rules

```text
s1ctl devicecontrol delete <rule-id>... [flags]
```

Delete one or more device control rules by ID.
Dry-run by default — pass --yes to apply.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply changes (default: dry-run) |

## devicecontrol disable

Disable device control rules

```text
s1ctl devicecontrol disable <rule-id>... [flags]
```

Disable one or more device control rules by ID.
Dry-run by default — pass --yes to apply.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply changes (default: dry-run) |

## devicecontrol enable

Enable device control rules

```text
s1ctl devicecontrol enable <rule-id>... [flags]
```

Enable one or more device control rules by ID.
Dry-run by default — pass --yes to apply.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply changes (default: dry-run) |

## devicecontrol events

List device control events

```text
s1ctl devicecontrol events [flags]
```

Show device control events from endpoints with Device Control-enabled Agents.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--all` | bool | false | fetch all pages |
| `--cursor` | string | - | pagination cursor |
| `--interface` | stringSlice | - | filter by interface (USB, Bluetooth, Thunderbolt, SDCard) |
| `--limit` | int | 0 | max results per page (default 50) |
| `--query` | string | - | free text search |
| `--site-id` | stringSlice | - | filter by site ID |

## devicecontrol list

List device control rules

```text
s1ctl devicecontrol list [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--all` | bool | false | fetch all pages |
| `--cursor` | string | - | pagination cursor |
| `--limit` | int | 0 | max results per page (default 50) |
| `--query` | string | - | free text search |
| `--site-id` | stringSlice | - | filter by site ID |

## devicecontrol pull

Pull device control rules to local YAML files

```text
s1ctl devicecontrol pull [flags]
```

Fetch all device control rules and write them as YAML files.

Each rule produces one file named by its sanitized rule name (e.g. block-usb-storage.yaml).
Server-only metadata (ID, scope, timestamps) is omitted from the YAML so the files
contain only the declarative rule definition.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | filter by account ID |
| `--out` | string | devicecontrol | output directory |
| `--site-id` | stringSlice | - | filter by site ID |

## devicecontrol push

Push device control rules from local YAML files

```text
s1ctl devicecontrol push [flags]
```

Read device control rule YAML files from a directory and sync them to SentinelOne.

Rules are matched by name: existing rules are updated, new rules are created.
Dry-run by default — pass --yes to apply changes.

New rules are created at the scope specified by --site-id. If no scope flag
is given, new rules are created at the global (tenant) scope.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--dir` | string | devicecontrol | directory containing device rule YAML files |
| `--site-id` | stringSlice | - | scope for new rules (default: global/tenant) |
| `--yes` | bool | false | apply changes (default: dry-run) |

## devicecontrol reorder

Reorder device control rules

```text
s1ctl devicecontrol reorder <id:order>... [flags]
```

Change the evaluation order of device control rules.

Each argument is an id:order pair, for example:

  s1ctl devicecontrol reorder 123:1 456:2 789:3

The order determines rule evaluation priority (1 = first).
Dry-run by default — pass --yes to apply.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | scope: account IDs |
| `--group-id` | stringSlice | - | scope: group IDs |
| `--site-id` | stringSlice | - | scope: site IDs |
| `--yes` | bool | false | apply changes (default: dry-run) |
