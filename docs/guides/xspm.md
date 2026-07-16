# xSPM: vulnerabilities and misconfigurations

Manage xSPM findings: list, triage, investigate, export, and track
vulnerabilities and misconfigurations.

> Prerequisites: `s1ctl` installed and configured (`S1_CONSOLE_URL`, `S1_TOKEN`).
> xSPM requires the Singularity Cloud Security or xSPM license.

## Vulnerabilities

### List vulnerabilities

```bash
s1ctl vulnerabilities list
s1ctl vulns list --severity HIGH,CRITICAL
s1ctl vulns list --status open --all --json
```

`vulns` is an alias for `vulnerabilities`.

| Flag | Description |
|------|-------------|
| `--severity` | Filter by severity: `HIGH`, `CRITICAL`, etc. (repeatable) |
| `--status` | Filter by status (repeatable) |
| `--limit` | Max results per page (default 50) |
| `--all` | Fetch all pages |
| `--after` | Pagination cursor |

### Get vulnerability details

```bash
s1ctl vulns get 000000
s1ctl vulns get 000000 --json
```

### Health summary

Show a breakdown of vulnerability counts by severity and open/resolved status.
Uses count queries -- no bulk data fetch needed.

```bash
s1ctl vulns health
s1ctl vulns health --json
```

### Posture stats

Summarize vulnerability posture: unique CVE count plus the top vulnerable
applications, assets, and OS types.

```bash
s1ctl vulns stats
s1ctl vulns stats --severity CRITICAL --limit 5
s1ctl vulns stats --top applications
s1ctl vulns stats --scope-level site --scope-id 000000
```

| Flag | Description |
|------|-------------|
| `--severity` | Filter by severity (repeatable) |
| `--top` | Show only one list: `applications`, `assets`, or `os` |
| `--limit` | Number of top entries per list (default 10) |
| `--scope-level` | Scope level: `account`, `site`, `group` |
| `--scope-id` | Scope ID |

### List CVEs

```bash
s1ctl vulns cves
s1ctl vulns cves --min-cvss 9.0 --all
s1ctl vulns cves --all --json
```

| Flag | Description |
|------|-------------|
| `--min-cvss` | Only show CVEs with NVD base score >= this value (client-side filter) |
| `--limit` | Max results per page (default 50) |
| `--all` | Fetch all pages |
| `--after` | Pagination cursor |

### Get CVE details

```bash
s1ctl vulns cve CVE-2024-00000
s1ctl vulns cve CVE-2024-00000 --json
```

### Investigation notes

```bash
# List notes
s1ctl vulns notes 000000

# Add a note
s1ctl vulns add-note 000000 --text "Investigating patch availability" --yes

# Update a note
s1ctl vulns update-note 000001 --text "Patch available in next release" --yes

# Delete a note
s1ctl vulns delete-note 000001 --yes
```

### Assign to a user

```bash
s1ctl vulns assign 000000 --user-id 000001 --yes
```

### History

```bash
s1ctl vulns history 000000
s1ctl vulns history 000000 --json
```

### Related assets

```bash
s1ctl vulns related-assets 000000
s1ctl vulns related-assets 000000 --json
```

### Update status

```bash
s1ctl vulns status 000000 resolved --yes
```

### Update verdict

```bash
s1ctl vulns verdict 000000 TRUE_POSITIVE --yes
s1ctl vulns verdict 000000 FALSE_POSITIVE --yes
```

### Export as CSV

```bash
s1ctl vulns export --out vulns.csv
s1ctl vulns export --severity CRITICAL --out critical.csv
s1ctl vulns export --scope-level site --scope-id 000000 --out site-vulns.csv
```

| Flag | Description |
|------|-------------|
| `--out` | Output file (default: stdout) |
| `--severity` | Filter by severity (repeatable) |
| `--status` | Filter by status (repeatable) |
| `--scope-level` | Scope level: `account`, `site`, `group` |
| `--scope-id` | Scope ID |

---

## Misconfigurations

### List misconfigurations

```bash
s1ctl misconfigurations list
s1ctl misconfigs list --severity HIGH,CRITICAL --all --json
```

`misconfigs` is an alias for `misconfigurations`.

| Flag | Description |
|------|-------------|
| `--severity` | Filter by severity (repeatable) |
| `--status` | Filter by status (repeatable) |
| `--limit` | Max results per page (default 50) |
| `--all` | Fetch all pages |
| `--after` | Pagination cursor |

### Get misconfiguration details

```bash
s1ctl misconfigs get 000000
s1ctl misconfigs get 000000 --json
```

### Investigation notes

```bash
s1ctl misconfigs notes 000000
s1ctl misconfigs add-note 000000 --text "Reviewing remediation steps" --yes
s1ctl misconfigs update-note 000001 --text "Remediation applied" --yes
s1ctl misconfigs delete-note 000001 --yes
```

### Assign to a user

```bash
s1ctl misconfigs assign 000000 --user-id 000001 --yes
```

### History

```bash
s1ctl misconfigs history 000000
```

### Related assets

```bash
s1ctl misconfigs related-assets 000000
```

### Update status

```bash
s1ctl misconfigs status 000000 resolved --yes
```

### Update verdict

```bash
s1ctl misconfigs verdict 000000 TRUE_POSITIVE --yes
s1ctl misconfigs verdict 000000 FALSE_POSITIVE --yes
```

### Export as CSV

```bash
s1ctl misconfigs export --out misconfigs.csv
s1ctl misconfigs export --severity CRITICAL --out critical.csv
s1ctl misconfigs export --scope-level site --scope-id 000000
```

| Flag | Description |
|------|-------------|
| `--out` | Output file (default: stdout) |
| `--severity` | Filter by severity (repeatable) |
| `--status` | Filter by status (repeatable) |
| `--scope-level` | Scope level: `account`, `site`, `group` |
| `--scope-id` | Scope ID |

---

## Workflows

### Triage critical vulnerabilities

```bash
s1ctl vulns list --severity CRITICAL --all --json \
  | jq '.[].id' -r \
  | while read id; do s1ctl vulns get "$id"; done
```

### Weekly posture report

```bash
s1ctl vulns health --json > vuln-health.json
s1ctl vulns stats --json > vuln-stats.json
```

### Find high-CVSS CVEs

```bash
s1ctl vulns cves --min-cvss 9.0 --all --json
```

### Export all findings for compliance

```bash
s1ctl vulns export --all --out vulnerabilities.csv
s1ctl misconfigs export --all --out misconfigurations.csv
```

## See also

- [`vulnerabilities` command reference](../commands/vulnerabilities.md)
- [`misconfigurations` command reference](../commands/misconfigurations.md)
