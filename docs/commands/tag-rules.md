# tag-rules

Manage dynamic asset tag rules

## tag-rules create

Create a dynamic tag rule from a JSON file

```text
s1ctl tag-rules create --from-file <rule.json> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--from-file` | string | - | tag rule definition JSON file (required) |
| `--yes` | bool | false | apply the action (default: dry-run) |

## tag-rules delete

Delete a dynamic tag rule

```text
s1ctl tag-rules delete <rule-id> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply the action (default: dry-run) |

## tag-rules list

List dynamic tag rules

```text
s1ctl tag-rules list [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | filter by account ID |
| `--all` | bool | false | fetch all pages |
| `--cursor` | string | - | pagination cursor |
| `--limit` | int | 0 | max results per page (default 50) |
| `--name` | string | - | filter by rule name |
| `--site-id` | stringSlice | - | filter by site ID |
| `--status` | string | - | filter by status (enabled, disabled) |

## tag-rules pull

Pull dynamic tag rules to local YAML files

```text
s1ctl tag-rules pull [flags]
```

Fetch all dynamic tag rules and write them as YAML files.

Each rule produces one file named by its sanitized name. Server-only metadata
(ID, scope IDs, audit fields, timestamps) is omitted so the files contain only
the declarative definition: name, status, conditions, scopes, tags, and excluded
assets.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | filter by account ID |
| `--out` | string | tag-rules | output directory |
| `--site-id` | stringSlice | - | filter by site ID |

## tag-rules push

Push dynamic tag rules from local YAML files

```text
s1ctl tag-rules push [flags]
```

Read tag rule YAML files from a directory and sync them to SentinelOne.

Rules are matched by name: existing rules are updated, new ones are created, and
unchanged ones are skipped. Dry-run by default — pass --yes to apply changes.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | filter by account ID |
| `--dir` | string | tag-rules | directory containing tag rule YAML files |
| `--site-id` | stringSlice | - | filter by site ID |
| `--yes` | bool | false | apply changes (default: dry-run) |

## tag-rules test

Report how many assets a candidate tag rule matches

```text
s1ctl tag-rules test --from-file <rule.json> [flags]
```

Report how many inventory assets a candidate tag rule would match, without
saving it. This is a read-only dry-check against live inventory.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--from-file` | string | - | tag rule definition JSON file (required) |

## tag-rules update

Update a dynamic tag rule from a JSON file

```text
s1ctl tag-rules update <rule-id> --from-file <rule.json> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--from-file` | string | - | tag rule definition JSON file (required) |
| `--yes` | bool | false | apply the action (default: dry-run) |
