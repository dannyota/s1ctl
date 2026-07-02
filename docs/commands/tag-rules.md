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
