# network

Manage network quarantine rules

## network configuration

Get or set network quarantine control configuration

```text
s1ctl network configuration
```

## network copy

Copy network quarantine rules between scopes

```text
s1ctl network copy [flags]
```

Copy network quarantine rules from a source scope to a target scope.

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
| `--target-group-id` | string | - | target group ID |
| `--target-site-id` | string | - | target site ID |
| `--yes` | bool | false | apply changes (default: dry-run) |

## network delete

Delete network quarantine rules

```text
s1ctl network delete <rule-id>... [flags]
```

Delete one or more network quarantine rules by ID.
Dry-run by default — pass --yes to apply.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply changes (default: dry-run) |

## network disable

Disable network quarantine rules

```text
s1ctl network disable <rule-id>... [flags]
```

Disable one or more network quarantine rules by ID.
Dry-run by default — pass --yes to apply.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply changes (default: dry-run) |

## network enable

Enable network quarantine rules

```text
s1ctl network enable <rule-id>... [flags]
```

Enable one or more network quarantine rules by ID.
Dry-run by default — pass --yes to apply.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply changes (default: dry-run) |

## network export

Export network quarantine rules to a JSON file

```text
s1ctl network export [flags]
```

Export network quarantine rules from a scope to a JSON file.
The exported file can be imported into another scope with "network import".

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--out` | string | network-quarantine-rules.json | output file (use - for stdout) |
| `--site-id` | stringSlice | - | scope: site IDs |

## network get

Get a network quarantine rule

```text
s1ctl network get <rule-id>
```

## network import

Import network quarantine rules from a JSON file

```text
s1ctl network import <file> [flags]
```

Import network quarantine rules from a previously exported JSON file into a scope.
Dry-run by default — pass --yes to apply.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | target account IDs |
| `--group-id` | stringSlice | - | target group IDs |
| `--site-id` | stringSlice | - | target site IDs |
| `--yes` | bool | false | apply changes (default: dry-run) |

## network list

List network quarantine rules

```text
s1ctl network list [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--all` | bool | false | fetch all pages |
| `--cursor` | string | - | pagination cursor |
| `--limit` | int | 0 | max results per page (default 50) |
| `--query` | string | - | free text search |
| `--site-id` | stringSlice | - | filter by site ID |

## network move

Move network quarantine rules to another scope

```text
s1ctl network move <rule-id>... [flags]
```

Move one or more network quarantine rules to a target scope.

Use --target-site-id, --target-account-id, or --target-group-id for the
destination. At least one target flag is required.
Dry-run by default — pass --yes to apply.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--target-account-id` | string | - | target account ID |
| `--target-group-id` | string | - | target group ID |
| `--target-site-id` | string | - | target site ID |
| `--yes` | bool | false | apply changes (default: dry-run) |

## network protocols

List available network quarantine protocols

```text
s1ctl network protocols [flags]
```

Show protocols that can be used in network quarantine rules.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--query` | string | - | search protocols |

## network pull

Pull network quarantine rules to local YAML files

```text
s1ctl network pull [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--out` | string | network-quarantine | output directory |
| `--site-id` | stringSlice | - | filter by site ID |

## network push

Push network quarantine rules from local YAML files

```text
s1ctl network push [flags]
```

Read network quarantine rule YAML files from a directory and sync them to SentinelOne.
Rules are matched by name: existing rules are updated, new rules are created,
and unchanged rules are skipped. Dry-run by default — pass --yes to apply changes.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--dir` | string | network-quarantine | directory containing network quarantine rule YAML files |
| `--site-id` | stringSlice | - | target site IDs |
| `--yes` | bool | false | apply changes (default: dry-run) |

## network reorder

Reorder network quarantine rules

```text
s1ctl network reorder <id:order>... [flags]
```

Change the evaluation order of network quarantine rules.

Each argument is an id:order pair, for example:

  s1ctl network reorder 123:1 456:2 789:3

The order determines rule evaluation priority (1 = first).
Dry-run by default — pass --yes to apply.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | scope: account IDs |
| `--group-id` | stringSlice | - | scope: group IDs |
| `--site-id` | stringSlice | - | scope: site IDs |
| `--yes` | bool | false | apply changes (default: dry-run) |

## network set-location

Set the location assignment of network quarantine rules

```text
s1ctl network set-location <rule-id>... [flags]
```

Assign a location matcher to one or more network quarantine rules.

--type is one of all, specific, or fallback. For "specific", pass one or more
--location-id values.
Dry-run by default — pass --yes to apply.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--location-id` | stringSlice | - | location IDs (for --type specific) |
| `--type` | string | all | location type: all, specific, or fallback |
| `--yes` | bool | false | apply changes (default: dry-run) |

## network tags

Add or remove tags on network quarantine rules

```text
s1ctl network tags
```
