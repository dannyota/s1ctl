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

## threats add-to-exclusions

Create an exclusion from a threat

```text
s1ctl threats add-to-exclusions <threat-id> [flags]
```

Create an exclusion from a threat, overriding the malicious verdict.
Scopes: group, site, account, tenant. Types: hash, path, certificate, browser, file_type.
Mode applies to path exclusions only (e.g. suppress, disable_all_monitors).

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--description` | string | - | exclusion description |
| `--mode` | string | - | exclusion mode (path exclusions only, e.g. suppress) |
| `--note` | string | - | note to add to the threat |
| `--path-exclusion-type` | string | - | excluded path type (path exclusions only) |
| `--scope` | string | - | exclusion scope (group, site, account, tenant) |
| `--ticket-id` | string | - | external ticket ID to set on the threat |
| `--type` | string | - | exclusion type (hash, path, certificate, browser, file_type) |
| `--value` | string | - | exclusion value (defaults to the threat's value) |
| `--yes` | bool | false | apply the action (default: dry-run) |

## threats blacklist

Add the threat file hash to the blacklist

```text
s1ctl threats blacklist <threat-id> [flags]
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

## threats exclusion-options

Show the exclusion (whitening) options available for a threat

```text
s1ctl threats exclusion-options <threat-id>
```

## threats export

Export threats to a CSV file

```text
s1ctl threats export [flags]
```

Export threats matching the filters as CSV. Writes to --out, or stdout when --out is omitted.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--classification` | stringSlice | - | filter by classification |
| `--mitigation-status` | stringSlice | - | filter by mitigation status |
| `--out` | string | - | output file (default: stdout) |
| `--query` | string | - | free text search |
| `--site-id` | stringSlice | - | filter by site ID |
| `--status` | stringSlice | - | filter by incident status |
| `--verdict` | stringSlice | - | filter by analyst verdict |

## threats fetch-file

Fetch the threat file from the endpoint to the console

```text
s1ctl threats fetch-file <threat-id> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply the action (default: dry-run) |

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
| `--status` | stringSlice | - | filter by incident status (unresolved, in_progress, resolved) |
| `--verdict` | stringSlice | - | filter by analyst verdict (true_positive, false_positive, suspicious, undefined) |

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

## threats mitigate-alerts

Mark an alert as a threat and run a mitigation action

```text
s1ctl threats mitigate-alerts [flags]
```

Mark a Deep Visibility alert (identified by agent ID and storyline) as a
threat and run a mitigation action.
Actions: kill, remediate, rollback-remediation, quarantine, un-quarantine, remove_macros, restore_macros.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--action` | string | - | mitigation action (kill, remediate, quarantine, etc.) |
| `--agent-id` | string | - | agent ID that reported the alert (required) |
| `--storyline` | string | - | storyline of the alert (required) |
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

## threats quarantined-files

List files quarantined for a threat

```text
s1ctl threats quarantined-files <threat-id>
```

## threats resolve

Resolve threats (bulk)

```text
s1ctl threats resolve [threat-id...] [flags]
```

Set incident status to "resolved" on one or more threats.

Specify threat IDs as arguments, or use filter flags to match threats.
Use typed flags (--classification, --verdict) or the generic --filter
flag with key=value pairs (e.g. --filter classifications=Malware).
Filter flags only match unresolved threats. Dry-run by default.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--classification` | stringSlice | - | filter by classification (e.g. Malware, PUP) |
| `--filter` | stringArray | - | key=value filter (e.g. --filter classifications=Malware) |
| `--mitigation-status` | stringSlice | - | filter by mitigation status |
| `--name` | string | - | match threats by name (contains, case-insensitive) |
| `--query` | string | - | free text search filter |
| `--site-id` | stringSlice | - | filter by site ID |
| `--verdict` | stringSlice | - | filter by analyst verdict (true_positive, false_positive, suspicious, undefined) |
| `--yes` | bool | false | apply the action (default: dry-run) |

## threats set-ticket

Set the external ticket ID on a threat

```text
s1ctl threats set-ticket <threat-id> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--ticket-id` | string | - | external ticket ID (required) |
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

## threats timeline

Show activity timeline for a threat

```text
s1ctl threats timeline <threat-id> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--all` | bool | false | fetch all pages |
| `--cursor` | string | - | pagination cursor |
| `--limit` | int | 0 | max results per page (default 50) |
| `--query` | string | - | free text search |
| `--sort-by` | string | - | sort field |
| `--sort-order` | string | - | sort direction (asc, desc) |

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
