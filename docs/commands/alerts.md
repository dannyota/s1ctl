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
| `--status` | stringSlice | - | filter by status (NEW, IN_PROGRESS, RESOLVED) |
| `--verdict` | stringSlice | - | filter by analyst verdict (e.g. FALSE_POSITIVE_BENIGN, TRUE_POSITIVE_MALWARE; see 'enums' cmd) |

## alerts counts

Count alert values per field (filter counts or group-by counts)

```text
s1ctl alerts counts --field <fieldId> [--field ...] [flags]
```

Return per-field value counts for the current alert selection.

By default uses alertFiltersCount (distinct filterable values and their
cardinality). Pass --group-by to use the deprecated alertGroupByCount query
instead; for grouped alert volume prefer "alerts stats" (alertGroups).

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--field` | stringSlice | - | field ID to count (repeatable, required) |
| `--group-by` | bool | false | use the deprecated alertGroupByCount query |
| `--scope-id` | string | - | account, site, or group ID |
| `--scope-level` | string | - | scope level (account, site, group) |
| `--severity` | stringSlice | - | filter by severity (HIGH, CRITICAL, etc.) |
| `--status` | stringSlice | - | filter by status (NEW, IN_PROGRESS, RESOLVED) |

## alerts delete-note

Delete an alert note

```text
s1ctl alerts delete-note <note-id> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply the action (default: dry-run) |

## alerts export

Export alerts to a CSV file

```text
s1ctl alerts export [flags]
```

Export alerts matching the filters as CSV via alertsCsvExport.

The API returns the full CSV inline. It is written to --out, or to stdout when
--out is omitted.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--out` | string | - | output file (default: stdout) |
| `--scope-id` | string | - | account, site, or group ID |
| `--scope-level` | string | - | scope level (account, site, group) |
| `--severity` | stringSlice | - | filter by severity (HIGH, CRITICAL, etc.) |
| `--status` | stringSlice | - | filter by status (NEW, IN_PROGRESS, RESOLVED) |
| `--verdict` | stringSlice | - | filter by analyst verdict |
| `--view` | string | - | predefined view (ALL, CLOUD, ENDPOINT, IDENTITY, CUSTOM_ALERTS, THIRD_PARTY) |

## alerts get

Get alert details

```text
s1ctl alerts get <id>
```

## alerts history

Show audit trail for an alert

```text
s1ctl alerts history <alert-id> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--all` | bool | false | fetch all pages |
| `--limit` | int | 0 | max results per page (default 50) |

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
| `--status` | stringSlice | - | filter by status (NEW, IN_PROGRESS, RESOLVED) |
| `--verdict` | stringSlice | - | filter by analyst verdict (e.g. FALSE_POSITIVE_BENIGN, TRUE_POSITIVE_MALWARE; see 'enums' cmd) |

## alerts notes

List investigation notes on an alert

```text
s1ctl alerts notes <alert-id>
```

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

## alerts stats

Show alert volume grouped by field

```text
s1ctl alerts stats [flags]
```

Show alert counts grouped by a specified field using the GraphQL alertGroups query.

Common group-by fields: severity, status, analystVerdict, classification,
detectionSource.product, assets.name.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--group-by` | string | severity | field to group by (e.g. severity, status, analystVerdict) |
| `--severity` | stringSlice | - | filter by severity (HIGH, CRITICAL, etc.) |
| `--status` | stringSlice | - | filter by status (NEW, RESOLVED, etc.) |

## alerts status

Update alert status (NEW, IN_PROGRESS, RESOLVED)

```text
s1ctl alerts status <id> <status> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply the action (default: dry-run) |

## alerts timeline

Show the timeline for an alert

```text
s1ctl alerts timeline <alert-id> [flags]
```

Show the alert timeline: notes, activities, enrichments, indicators,
asset operations, mitigation actions, and related alerts, newest first.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--all` | bool | false | fetch all pages |
| `--limit` | int | 0 | max results per page (default 50) |

## alerts update-note

Update the text of an alert note

```text
s1ctl alerts update-note <note-id> --text <text> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--text` | string | - | new note text (required) |
| `--yes` | bool | false | apply the action (default: dry-run) |

## alerts verdict

Update alert analyst verdict (e.g. FALSE_POSITIVE_BENIGN, TRUE_POSITIVE_MALWARE)

```text
s1ctl alerts verdict <id> <verdict> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply the action (default: dry-run) |
