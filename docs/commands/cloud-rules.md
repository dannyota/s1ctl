# cloud-rules

Manage CNS custom cloud rules (Cloud Native Security)

## cloud-rules create

Create a CNS custom cloud rule from a JSON file

```text
s1ctl cloud-rules create --from-file <rule.json> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--from-file` | string | - | path to rule JSON file (required) |
| `--scope-id` | string | - | account, site, or group ID |
| `--scope-level` | string | - | scope level (account, site, group) |
| `--yes` | bool | false | apply (default: dry-run) |

## cloud-rules delete

Delete CNS custom cloud rules

```text
s1ctl cloud-rules delete <id> [id...] [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--scope-id` | string | - | account, site, or group ID |
| `--scope-level` | string | - | scope level (account, site, group) |
| `--yes` | bool | false | apply the action (default: dry-run) |

## cloud-rules disable

Disable CNS custom cloud rules

```text
s1ctl cloud-rules disable <id> [id...] [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--scope-id` | string | - | account, site, or group ID |
| `--scope-level` | string | - | scope level (account, site, group) |
| `--yes` | bool | false | apply the action (default: dry-run) |

## cloud-rules enable

Enable CNS custom cloud rules

```text
s1ctl cloud-rules enable <id> [id...] [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--scope-id` | string | - | account, site, or group ID |
| `--scope-level` | string | - | scope level (account, site, group) |
| `--yes` | bool | false | apply the action (default: dry-run) |

## cloud-rules evaluate

Evaluate a Rego query against asset JSON (dry-check)

```text
s1ctl cloud-rules evaluate --rule <rule.json> --resource <resource.json> [flags]
```

Evaluate a raw Rego query against an asset's JSON before creating or
updating a CNS rule. This is a read-only dry-check: it evaluates only and
mutates nothing.

The Rego query comes from the --rule file's rawQuery field, or from --query
when set. The asset JSON to test against is read from --resource.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--config` | string | - | inline rule config parameters JSON string |
| `--policy-id` | string | - | policy ID to source mandatory parameters |
| `--query` | string | - | inline Rego query (overrides rule file rawQuery) |
| `--resource` | string | - | asset JSON file to evaluate against (required) |
| `--rule` | string | - | rule JSON file (rawQuery extracted from it) |
| `--scope-id` | string | - | account, site, or group ID |
| `--scope-level` | string | - | scope level (account, site, group) |

## cloud-rules get

Get CNS custom cloud rule details

```text
s1ctl cloud-rules get <id> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--scope-id` | string | - | account, site, or group ID |
| `--scope-level` | string | - | scope level (account, site, group) |

## cloud-rules list

List CNS custom cloud rules

```text
s1ctl cloud-rules list [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--after` | string | - | pagination cursor |
| `--all` | bool | false | fetch all pages |
| `--limit` | int | 0 | max results per page (default 50) |
| `--scope-id` | string | - | account, site, or group ID |
| `--scope-level` | string | - | scope level (account, site, group) |
| `--severity` | stringSlice | - | filter by severity (LOW, MEDIUM, HIGH, CRITICAL) |
| `--status` | stringSlice | - | filter by status |

## cloud-rules types

List supported CNS rule types

```text
s1ctl cloud-rules types [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--scope-id` | string | - | account, site, or group ID |
| `--scope-level` | string | - | scope level (account, site, group) |

## cloud-rules update

Replace a CNS custom cloud rule from a JSON file

```text
s1ctl cloud-rules update <id> --from-file <rule.json> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--from-file` | string | - | path to rule JSON file (required) |
| `--scope-id` | string | - | account, site, or group ID |
| `--scope-level` | string | - | scope level (account, site, group) |
| `--yes` | bool | false | apply (default: dry-run) |
