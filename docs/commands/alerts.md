# alerts

Manage unified alerts (GraphQL UAM)

## alerts add-note

Add an investigation note to an alert

```text
s1ctl alerts add-note <alert-id> <text> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply the action (default: dry-run) |

## alerts count

Count alerts

```text
s1ctl alerts count [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--severity` | stringSlice | - | filter by severity (HIGH, CRITICAL, etc.) |
| `--status` | stringSlice | - | filter by status (NEW, RESOLVED, etc.) |
| `--verdict` | stringSlice | - | filter by analyst verdict |

## alerts get

Get alert details

```text
s1ctl alerts get <id>
```

## alerts list

List alerts

```text
s1ctl alerts list [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--after` | string | - | pagination cursor |
| `--all` | bool | false | fetch all pages |
| `--limit` | int | 0 | max results per page (default 50) |
| `--severity` | stringSlice | - | filter by severity (HIGH, CRITICAL, etc.) |
| `--sort-by` | string | - | sort field (e.g. detectedAt, severity) |
| `--sort-order` | string | - | sort direction (ASC, DESC) |
| `--source` | stringSlice | - | filter by detection source (STAR, EDR, CWS) |
| `--status` | stringSlice | - | filter by status (NEW, RESOLVED, etc.) |
| `--verdict` | stringSlice | - | filter by analyst verdict |

## alerts resolve

Resolve alerts by ID or filter

```text
s1ctl alerts resolve [id...] [flags]
```

Set status to "RESOLVED" on one or more alerts.

Specify alert IDs directly, or use --name/--severity/--source to match alerts.
Filter flags only match alerts with status NEW. Dry-run by default.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--name` | string | - | match alerts by name (contains, case-insensitive) |
| `--severity` | stringSlice | - | filter by severity (HIGH, CRITICAL) |
| `--source` | stringSlice | - | filter by detection source (STAR, EDR, CWS) |
| `--yes` | bool | false | apply the action (default: dry-run) |

## alerts status

Update alert status

```text
s1ctl alerts status <id> <status> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply the action (default: dry-run) |

## alerts verdict

Update alert analyst verdict

```text
s1ctl alerts verdict <id> <verdict> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply the action (default: dry-run) |
