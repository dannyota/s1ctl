# Threats

List, inspect, triage, and mitigate threats.

## List threats

```bash
s1ctl threats list
s1ctl threats list --status unresolved --limit 20
s1ctl threats list --classification malware --verdict suspicious
s1ctl threats list --site-id 000000 --sort-by createdAt --sort-order desc
```

### Flags

| Flag | Type | Description |
|------|------|-------------|
| `--query` | string | Free text search across threat fields |
| `--site-id` | string[] | Filter by site ID (repeatable) |
| `--classification` | string[] | Filter by classification (repeatable) |
| `--status` | string[] | Filter by incident status (repeatable) |
| `--verdict` | string[] | Filter by analyst verdict (repeatable) |
| `--sort-by` | string | Sort field (e.g. `createdAt`, `classification`) |
| `--sort-order` | string | Sort direction (`asc`, `desc`) |
| `--limit` | int | Max results per page (default 50) |
| `--all` | bool | Fetch all pages automatically |
| `--cursor` | string | Pagination cursor for manual paging |

The `--output` and `--json` flags are global and work on all read commands.

## Get threat details

```bash
s1ctl threats get 000000
s1ctl threats get 000000 --json
```

Returns a detail view: ID, name, classification, confidence level,
mitigation status, analyst verdict, incident status, agent ID, and
creation date.

### Notes

List investigation notes attached to a threat:

```bash
s1ctl threats notes 000000
s1ctl threats notes 000000 --json
```

### Timeline

Show the activity timeline for a threat:

```bash
s1ctl threats timeline 000000
s1ctl threats timeline 000000 --json
```

## Actions

All actions are **dry-run by default**. Pass `--yes` to apply.

### Mitigate

Apply a mitigation action to a threat.

```bash
s1ctl threats mitigate 000000 --action kill          # dry-run
s1ctl threats mitigate 000000 --action kill --yes     # apply
```

| Action | Description |
|--------|-------------|
| `kill` | Kill the threat process |
| `quarantine` | Quarantine the threat file |
| `remediate` | Remediate (kill + quarantine + undo changes) |
| `rollback-remediation` | Undo a previous remediation |

### Update verdict

Set the analyst verdict on a threat.

```bash
s1ctl threats verdict 000000 --verdict true_positive        # dry-run
s1ctl threats verdict 000000 --verdict true_positive --yes   # apply
```

| Verdict | Description |
|---------|-------------|
| `true_positive` | Confirmed threat |
| `false_positive` | Not a threat |
| `suspicious` | Requires further investigation |
| `undefined` | Reset to no verdict |

### Update status

Set the incident status on a threat.

```bash
s1ctl threats status 000000 --status resolved              # dry-run
s1ctl threats status 000000 --status resolved --yes         # apply
```

| Status | Description |
|--------|-------------|
| `unresolved` | Not yet addressed |
| `in_progress` | Under investigation |
| `resolved` | Fully handled |

### Blacklist

Add the threat's file hash to the blacklist so the same file is blocked
across the tenant on next encounter.

```bash
s1ctl threats blacklist 000000        # dry-run
s1ctl threats blacklist 000000 --yes   # apply
```

### Add to exclusions

Create an exclusion directly from a threat. The exclusion type and hash
are auto-populated from the threat.

```bash
s1ctl threats add-to-exclusions 000000        # dry-run
s1ctl threats add-to-exclusions 000000 --yes   # apply
```

### Mitigate alerts

Mark an alert as a threat and run a mitigation action in one step.

```bash
s1ctl threats mitigate-alerts 000000 --action kill          # dry-run
s1ctl threats mitigate-alerts 000000 --action kill --yes     # apply
```

| Action | Description |
|--------|-------------|
| `kill` | Kill the threat process |
| `quarantine` | Quarantine the threat file |
| `remediate` | Remediate (kill + quarantine + undo changes) |
| `rollback-remediation` | Undo a previous remediation |

### Set ticket

Set the external ticket ID and/or URL on a threat for tracking in an
external ticketing system.

```bash
s1ctl threats set-ticket 000000 --ticket-id "INC-1234"                    # dry-run
s1ctl threats set-ticket 000000 --ticket-id "INC-1234" --yes              # apply
s1ctl threats set-ticket 000000 --ticket-url "https://jira.example.com/INC-1234" --yes
s1ctl threats set-ticket 000000 --ticket-id "INC-1234" --ticket-url "https://jira.example.com/INC-1234" --yes
```

### Fetch file

Retrieve the threat file from the endpoint to the console for offline
analysis (for example, to download it later for sandbox detonation).

```bash
s1ctl threats fetch-file 000000        # dry-run
s1ctl threats fetch-file 000000 --yes   # apply
```

## Read-only inspection

### Quarantined files

List files quarantined for a threat:

```bash
s1ctl threats quarantined-files 000000
s1ctl threats quarantined-files 000000 --json
```

### Exclusion options

Show the exclusion (whitening) options available for a threat. Useful
to review what exclusion types can be created before running
`add-to-exclusions`.

```bash
s1ctl threats exclusion-options 000000
s1ctl threats exclusion-options 000000 --json
```

## Export

Export threats to a CSV file. Use `--all` to fetch all pages before
writing the file.

```bash
s1ctl threats export --output threats.csv
s1ctl threats export --all --output threats.csv
s1ctl threats export --status unresolved --all --output unresolved.csv
```

## Workflows

### Triage unresolved threats

List all unresolved threats, newest first:

```bash
s1ctl threats list --status unresolved --sort-by createdAt --sort-order desc
```

Inspect a specific threat, then mitigate and mark:

```bash
s1ctl threats get 000000
s1ctl threats mitigate 000000 --action kill --yes
s1ctl threats verdict 000000 --verdict true_positive --yes
s1ctl threats status 000000 --status resolved --yes
```

### Bulk export threats

Export all threats as JSON for offline analysis:

```bash
s1ctl threats list --all --json > threats.json
```

Export as CSV for spreadsheets:

```bash
s1ctl threats list --all --output csv > threats.csv
```

### Filter by classification and verdict

```bash
s1ctl threats list --classification malware --verdict suspicious --json \
  | jq '.[].threatName'
```

### Scope to a site

```bash
s1ctl threats list --site-id 000000 --status unresolved
```

## Output formats

| Flag | Format | Use case |
|------|--------|----------|
| (default) | table | Human-readable terminal output |
| `--json` | JSON | Pipe to jq, scripts, automation |
| `--output csv` | CSV | Spreadsheets, bulk analysis |

## FAQ

### How do I see all unresolved threats?

```bash
s1ctl threats list --status unresolved --all
```

### What is the difference between mitigate actions?

`kill` stops the process. `quarantine` isolates the file. `remediate`
does all three: kill, quarantine, and undo changes. Use `remediate` for
the most thorough response.

### How do I mark a threat as a false positive?

```bash
s1ctl threats verdict 000000 --verdict false_positive --yes
s1ctl threats status 000000 --status resolved --yes
```

### Can I link threats to my ticketing system?

Yes. Use `set-ticket` with `--ticket-id` and/or `--ticket-url`:

```bash
s1ctl threats set-ticket 000000 --ticket-id "INC-1234" --yes
```
