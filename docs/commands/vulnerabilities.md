# vulnerabilities

Manage xSPM vulnerabilities

## vulnerabilities assign

Assign a vulnerability to a user

```text
s1ctl vulnerabilities assign <id> --user-id <user-id> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--user-id` | string | - | assignee user ID (required) |
| `--yes` | bool | false | apply the action (default: dry-run) |

## vulnerabilities cve

Get CVE details

```text
s1ctl vulnerabilities cve <id>
```

## vulnerabilities cves

List CVEs

```text
s1ctl vulnerabilities cves [flags]
```

List CVEs via the cves query.

The cves server-side filter (CveFilterInput) supports only datetime-range
filtering, so --min-cvss is applied client-side against each CVE's NVD base
score after fetching. It only sees the fetched page, so pair it with --all to
filter the full result set rather than a single page.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--after` | string | - | pagination cursor |
| `--all` | bool | false | fetch all pages |
| `--limit` | int | 0 | max results per page (default 50) |
| `--min-cvss` | float64 | 0 | only show CVEs with NVD base score >= this value (client-side) |

## vulnerabilities export

Export vulnerabilities to a CSV file

```text
s1ctl vulnerabilities export [flags]
```

Export vulnerabilities matching the filters as CSV via
vulnerabilitiesExportToCsv. The API returns the full CSV inline; it is written
to --out, or to stdout when --out is omitted.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--out` | string | - | output file (default: stdout) |
| `--scope-id` | string | - | account, site, or group ID |
| `--scope-level` | string | - | scope level (account, site, group) |
| `--severity` | stringSlice | - | filter by severity (HIGH, CRITICAL, etc.) |
| `--status` | stringSlice | - | filter by status |

## vulnerabilities get

Get vulnerability details

```text
s1ctl vulnerabilities get <id>
```

## vulnerabilities health

Summarize vulnerabilities by severity and status

```text
s1ctl vulnerabilities health
```

Show a breakdown of vulnerability counts by severity and open/resolved status.
Uses count queries — no bulk data fetch needed.

## vulnerabilities history

Show the history of a vulnerability

```text
s1ctl vulnerabilities history <id>
```

## vulnerabilities list

List vulnerabilities

```text
s1ctl vulnerabilities list [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--after` | string | - | pagination cursor |
| `--all` | bool | false | fetch all pages |
| `--limit` | int | 0 | max results per page (default 50) |
| `--severity` | stringSlice | - | filter by severity (HIGH, CRITICAL, etc.) |
| `--status` | stringSlice | - | filter by status |

## vulnerabilities note-add

Add an investigation note to a vulnerability

```text
s1ctl vulnerabilities note-add <id> --text <text> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--text` | string | - | note text (required) |
| `--yes` | bool | false | apply the action (default: dry-run) |

## vulnerabilities note-delete

Delete a vulnerability note

```text
s1ctl vulnerabilities note-delete <note-id> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply the action (default: dry-run) |

## vulnerabilities note-update

Update the text of a vulnerability note

```text
s1ctl vulnerabilities note-update <note-id> --text <text> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--text` | string | - | new note text (required) |
| `--yes` | bool | false | apply the action (default: dry-run) |

## vulnerabilities notes

List investigation notes on a vulnerability

```text
s1ctl vulnerabilities notes <id>
```

## vulnerabilities related-assets

List assets related to a vulnerability

```text
s1ctl vulnerabilities related-assets <id>
```

## vulnerabilities stats

Summarize vulnerability posture (unique CVEs + top vulnerable applications/assets/OS)

```text
s1ctl vulnerabilities stats [flags]
```

Summarize vulnerability posture.

By default reports the unique CVE count plus the top vulnerable applications,
assets, and OS types. Pass --top applications|assets|os to show only one list.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--limit` | int | 0 | number of top entries per list (default 10) |
| `--scope-id` | string | - | account, site, or group ID |
| `--scope-level` | string | - | scope level (account, site, group) |
| `--severity` | stringSlice | - | filter by severity (HIGH, CRITICAL, etc.) |
| `--top` | string | - | show only one list: applications, assets, or os |

## vulnerabilities status

Update vulnerability status

```text
s1ctl vulnerabilities status <id> <status> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply the action (default: dry-run) |

## vulnerabilities verdict

Update vulnerability analyst verdict (TRUE_POSITIVE, FALSE_POSITIVE, SUSPICIOUS, UNDEFINED)

```text
s1ctl vulnerabilities verdict <id> <verdict> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply the action (default: dry-run) |
