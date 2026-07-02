# reports

Manage reports and report tasks

## reports create

Create a report task

```text
s1ctl reports create [flags]
```

Create a new report task or schedule.

Schedule types: manually, scheduled
Frequencies: manually, weekly, monthly
Days (for weekly): sunday, monday, tuesday, wednesday, thursday, friday, saturday

Use "reports types" to list available insight types, then pass them
as a JSON array via --insight-types.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | target account IDs |
| `--attachment-type` | stringSlice | - | attachment types (pdf, html) |
| `--day` | string | - | day of week for weekly schedules |
| `--frequency` | string | - | frequency: manually, weekly, monthly |
| `--from-date` | string | - | report date range start (ISO timestamp) |
| `--insight-types` | string | - | insight types as JSON array (required; see 'reports types') |
| `--name` | string | - | report task name (required) |
| `--recipient` | stringSlice | - | email recipients |
| `--schedule-type` | string | - | schedule type: manually, scheduled (required) |
| `--scope` | string | - | scope filter |
| `--site-id` | stringSlice | - | target site IDs |
| `--to-date` | string | - | report date range end (ISO timestamp) |
| `--trend` | bool | false | trend report (period = last month) |
| `--yes` | bool | false | apply the action (default: dry-run) |

## reports download

Download a generated report

```text
s1ctl reports download <report-id> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--format` | string | pdf | report format (pdf, html) |
| `--output` | string | - | output file path (default: report-<id>.<format>) |

## reports list

List generated reports

```text
s1ctl reports list [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--all` | bool | false | fetch all pages |
| `--cursor` | string | - | pagination cursor |
| `--frequency` | string | - | filter by frequency (manually, weekly, monthly) |
| `--limit` | int | 0 | max results per page (default 50) |
| `--query` | string | - | free text search |
| `--schedule-type` | string | - | filter by schedule type (manually, scheduled) |
| `--scope` | string | - | filter by scope (group, site, account, tenant) |
| `--site-id` | stringSlice | - | filter by site ID |
| `--sort-by` | string | - | sort field (e.g. name, createdAt, status) |
| `--sort-order` | string | - | sort direction (asc, desc) |

## reports tasks

List report tasks and schedules

```text
s1ctl reports tasks [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--all` | bool | false | fetch all pages |
| `--cursor` | string | - | pagination cursor |
| `--frequency` | string | - | filter by frequency (manually, weekly, monthly) |
| `--limit` | int | 0 | max results per page (default 50) |
| `--query` | string | - | free text search |
| `--schedule-type` | string | - | filter by schedule type (manually, scheduled) |
| `--scope` | string | - | filter by scope (group, site, account, tenant) |
| `--site-id` | stringSlice | - | filter by site ID |
| `--sort-by` | string | - | sort field (e.g. name, frequency, scope) |
| `--sort-order` | string | - | sort direction (asc, desc) |

## reports types

List available report insight types

```text
s1ctl reports types [flags]
```

List available report insight types. Output is always JSON because the schema is opaque.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | filter by account ID |
| `--group-id` | stringSlice | - | filter by group ID |
| `--site-id` | stringSlice | - | filter by site ID |
