# rules

Manage custom detection rules (STAR)

## rules detections

List recent detections for a rule

```text
s1ctl rules detections <rule-name> [flags]
```

Fetch cloud detection alerts (STAR alerts) filtered by rule name.
Shows what a specific rule is catching.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--all` | bool | false | fetch all pages |
| `--cursor` | string | - | pagination cursor |
| `--group-by` | string | - | group results (agent) |
| `--limit` | int | 0 | max results per page (default 50) |
| `--severity` | stringSlice | - | filter by severity |
| `--since` | string | - | show detections after this time (RFC3339) |
| `--site-id` | stringSlice | - | filter by site ID |
| `--sort-by` | string | - | sort field (default: id) |
| `--sort-order` | string | - | sort direction (default: desc) |
| `--status` | stringSlice | - | filter by incident status |

## rules diff

Compare local rule YAML files against live rules

```text
s1ctl rules diff [flags]
```

Read rule YAML files from a directory, fetch corresponding live rules
by name, and show what differs. Helps review changes before pushing.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--dir` | string | rules | directory containing rule YAML files |

## rules disable

Disable custom detection rules

```text
s1ctl rules disable <rule-id>... [flags]
```

Deactivate one or more custom detection rules by ID.
Dry-run by default — pass --yes to apply.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply changes (default: dry-run) |

## rules enable

Enable custom detection rules

```text
s1ctl rules enable <rule-id>... [flags]
```

Activate one or more custom detection rules by ID.
Dry-run by default — pass --yes to apply.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply changes (default: dry-run) |

## rules get

Get custom detection rule details

```text
s1ctl rules get <rule-id>
```

## rules health

Classify rules by operational state

```text
s1ctl rules health [flags]
```

Fetch all custom detection rules and classify them as firing (active
with alerts), silent (active with zero alerts), disabled, or erroring
(reached alert limit). Helps identify rules that need attention.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--site-id` | stringSlice | - | filter by site ID |

## rules list

List custom detection rules

```text
s1ctl rules list [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--all` | bool | false | fetch all pages |
| `--cursor` | string | - | pagination cursor |
| `--limit` | int | 0 | max results per page (default 50) |
| `--name` | string | - | filter by rule name (substring match) |
| `--query` | string | - | free text search on S1QL |
| `--query-type` | stringSlice | - | filter by query type (events, correlation, scheduled) |
| `--scope` | stringSlice | - | filter by scope (global, account, site, group) |
| `--severity` | stringSlice | - | filter by severity (Info, Low, Medium, High, Critical) |
| `--site-id` | stringSlice | - | filter by site ID |
| `--sort-by` | string | - | sort field (e.g. name, severity, createdAt) |
| `--sort-order` | string | - | sort direction (asc, desc) |
| `--status` | stringSlice | - | filter by status (Draft, Active, Disabled, ...) |

## rules pull

Pull custom detection rules to local YAML files

```text
s1ctl rules pull [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--out` | string | rules | output directory |
| `--site-id` | stringSlice | - | filter by site ID |

## rules push

Push custom detection rules from local YAML files

```text
s1ctl rules push [flags]
```

Read rule YAML files from a directory and sync them to SentinelOne.
Rules are matched by name: existing rules are updated, new rules are created,
and unchanged rules are skipped. Dry-run by default — pass --yes to apply changes.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--dir` | string | rules | directory containing rule YAML files |
| `--yes` | bool | false | apply changes (default: dry-run) |

## rules trends

Show noisiest rules by detection count

```text
s1ctl rules trends [flags]
```

Fetch all custom detection rules and sort by generated alert count
(descending). Helps identify alert fatigue candidates for tuning.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--site-id` | stringSlice | - | filter by site ID |
| `--top` | int | 0 | show only top N rules (default: all) |

## rules validate

Validate rule YAML files without deploying

```text
s1ctl rules validate [flags]
```

Read rule YAML files from a directory and check for errors:
missing required fields, invalid enum values, and empty queries.
No API calls are made — this is a local-only check.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--dir` | string | rules | directory containing rule YAML files |
