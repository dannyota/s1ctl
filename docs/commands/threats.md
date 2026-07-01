# threats

Manage threats

## threats add-note

Add a note to a threat

```text
s1ctl threats add-note <threat-id> <text> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply the action (default: dry-run) |

## threats count

Count threats

```text
s1ctl threats count [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--site-id` | stringSlice | - | filter by site ID |

## threats get

Get threat details

```text
s1ctl threats get <threat-id>
```

## threats list

List threats

```text
s1ctl threats list [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--all` | bool | false | fetch all pages |
| `--classification` | stringSlice | - | filter by classification |
| `--cursor` | string | - | pagination cursor |
| `--limit` | int | 0 | max results per page (default 50) |
| `--mitigation-status` | stringSlice | - | filter by mitigation status (not_mitigated, mitigated, etc.) |
| `--query` | string | - | free text search |
| `--site-id` | stringSlice | - | filter by site ID |
| `--sort-by` | string | - | sort field (e.g. createdAt, classification) |
| `--sort-order` | string | - | sort direction (asc, desc) |
| `--status` | stringSlice | - | filter by incident status |
| `--verdict` | stringSlice | - | filter by analyst verdict |

## threats mitigate

Apply mitigation action to a threat

```text
s1ctl threats mitigate <threat-id> [flags]
```

Actions: kill, quarantine, remediate, rollback-remediation

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--action` | string | - | mitigation action (kill, quarantine, remediate, rollback-remediation) |
| `--yes` | bool | false | apply the action (default: dry-run) |

## threats notes

List notes for a threat

```text
s1ctl threats notes <threat-id> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--all` | bool | false | fetch all pages |
| `--cursor` | string | - | pagination cursor |
| `--limit` | int | 0 | max results per page (default 50) |
| `--sort-by` | string | - | sort field |
| `--sort-order` | string | - | sort direction (asc, desc) |

## threats resolve

Resolve threats (bulk)

```text
s1ctl threats resolve [threat-id...] [flags]
```

Set incident status to "resolved" on one or more threats.

Specify threat IDs as arguments, or use filter flags to match threats.
Filter flags only match unresolved threats. Dry-run by default.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--classification` | stringSlice | - | filter by classification (e.g. Malware, PUP) |
| `--mitigation-status` | stringSlice | - | filter by mitigation status |
| `--name` | string | - | match threats by name (contains, case-insensitive) |
| `--query` | string | - | free text search filter |
| `--site-id` | stringSlice | - | filter by site ID |
| `--verdict` | stringSlice | - | filter by analyst verdict |
| `--yes` | bool | false | apply the action (default: dry-run) |

## threats status

Update incident status on a threat

```text
s1ctl threats status <threat-id> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--status` | string | - | incident status (unresolved, in_progress, resolved) |
| `--yes` | bool | false | apply the action (default: dry-run) |

## threats verdict

Update analyst verdict on a threat

```text
s1ctl threats verdict <threat-id> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--verdict` | string | - | analyst verdict (true_positive, false_positive, suspicious, undefined) |
| `--yes` | bool | false | apply the action (default: dry-run) |
