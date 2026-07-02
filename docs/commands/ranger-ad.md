# ranger-ad

Ranger AD exposure assessment (ISPM)

## ranger-ad affected-objects

List objects affected by an AD exposure

```text
s1ctl ranger-ad affected-objects [flags]
```

List Active Directory objects affected by a specific detection.
Requires --detection and --domain flags to identify the exposure.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | filter by account ID |
| `--all` | bool | false | fetch all pages |
| `--detection` | stringSlice | - | detection name (required) |
| `--domain` | stringSlice | - | domain name (required) |
| `--limit` | int | 0 | max results per page (default 50) |
| `--object-type` | stringSlice | - | filter by object type (Computer, User, Group, ...) |
| `--site-id` | stringSlice | - | filter by site ID |

## ranger-ad assess

Trigger a new AD assessment

```text
s1ctl ranger-ad assess [flags]
```

Trigger a Ranger AD assessment scan.
Use --full-scan for a complete scan, or omit for a targeted reassessment.
Dry-run by default — pass --yes to apply.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | filter by account ID |
| `--domain` | stringSlice | - | domain names to scan |
| `--full-scan` | bool | false | perform a full scan (default: targeted) |
| `--scan-source` | string | - | scan source (AD, Azure) |
| `--site-id` | stringSlice | - | filter by site ID |
| `--yes` | bool | false | apply changes (default: dry-run) |

## ranger-ad exposures

List AD exposures

```text
s1ctl ranger-ad exposures [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | filter by account ID |
| `--all` | bool | false | fetch all pages |
| `--detection` | stringSlice | - | filter by detection name |
| `--domain` | stringSlice | - | filter by domain name |
| `--limit` | int | 0 | max results per page (default 50) |
| `--severity` | stringSlice | - | filter by severity (Critical, High, Medium, Low) |
| `--site-id` | stringSlice | - | filter by site ID |
| `--source` | stringSlice | - | filter by source (OnPremAD, AzureAD) |
| `--status` | stringSlice | - | filter by detection status (Vulnerable, Not_Vulnerable, Skipped, ...) |

## ranger-ad status

Show AD assessment status

```text
s1ctl ranger-ad status [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | filter by account ID |
| `--site-id` | stringSlice | - | filter by site ID |
