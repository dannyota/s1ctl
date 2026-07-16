# firewall

Manage firewall control rules

## firewall copy

Copy firewall rules between scopes

```text
s1ctl firewall copy [flags]
```

Copy firewall rules from a source scope to a target scope.

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

## firewall delete

Delete firewall rules

```text
s1ctl firewall delete <rule-id>... [flags]
```

Delete one or more firewall rules by ID.
Dry-run by default — pass --yes to apply.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply changes (default: dry-run) |

## firewall disable

Disable firewall rules

```text
s1ctl firewall disable <rule-id>... [flags]
```

Disable one or more firewall rules by ID.
Dry-run by default — pass --yes to apply.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply changes (default: dry-run) |

## firewall enable

Enable firewall rules

```text
s1ctl firewall enable <rule-id>... [flags]
```

Enable one or more firewall rules by ID.
Dry-run by default — pass --yes to apply.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply changes (default: dry-run) |

## firewall export

Export firewall rules to a JSON file

```text
s1ctl firewall export [flags]
```

Export firewall rules from a scope to a JSON file.
The exported file can be imported into another scope with "firewall import".

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--out` | string | firewall-rules.json | output file (use - for stdout) |
| `--site-id` | stringSlice | - | scope: site IDs |

## firewall get

Get a firewall rule

```text
s1ctl firewall get <rule-id>
```

## firewall import

Import firewall rules from a JSON file

```text
s1ctl firewall import <file> [flags]
```

Import firewall rules from a previously exported JSON file into a scope.
Dry-run by default — pass --yes to apply.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | target account IDs |
| `--group-id` | stringSlice | - | target group IDs |
| `--site-id` | stringSlice | - | target site IDs |
| `--yes` | bool | false | apply changes (default: dry-run) |

## firewall list

List firewall rules

```text
s1ctl firewall list [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--all` | bool | false | fetch all pages |
| `--cursor` | string | - | pagination cursor |
| `--limit` | int | 0 | max results per page (default 50) |
| `--query` | string | - | free text search |
| `--site-id` | stringSlice | - | filter by site ID |

## firewall protocols

List available firewall protocols

```text
s1ctl firewall protocols [flags]
```

Show protocols that can be used in firewall rules.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--query` | string | - | search protocols |

## firewall pull

Pull firewall rules to local YAML files

```text
s1ctl firewall pull [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--out` | string | firewall | output directory |
| `--site-id` | stringSlice | - | filter by site ID |

## firewall push

Push firewall rules from local YAML files

```text
s1ctl firewall push [flags]
```

Read firewall rule YAML files from a directory and sync them to SentinelOne.
Rules are matched by name: existing rules are updated, new rules are created,
and unchanged rules are skipped. Dry-run by default — pass --yes to apply changes.

Note: the plan is built against the list scoped by --site-id. Without --site-id
the list is unscoped, which may match rules from other sites and produce an
incorrect plan. Always pass --site-id when pushing to a specific site.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--dir` | string | firewall | directory containing firewall rule YAML files |
| `--site-id` | stringSlice | - | target site IDs |
| `--yes` | bool | false | apply changes (default: dry-run) |

## firewall reorder

Reorder firewall rules

```text
s1ctl firewall reorder <id:order>... [flags]
```

Change the evaluation order of firewall rules.

Each argument is an id:order pair, for example:

  s1ctl firewall reorder 123:1 456:2 789:3

The order determines rule evaluation priority (1 = first).
Dry-run by default — pass --yes to apply.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | scope: account IDs |
| `--group-id` | stringSlice | - | scope: group IDs |
| `--site-id` | stringSlice | - | scope: site IDs |
| `--yes` | bool | false | apply changes (default: dry-run) |
