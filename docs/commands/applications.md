# applications

Application inventory and risk management

## applications cves

List CVEs across applications

```text
s1ctl applications cves [flags]
```

List CVEs across applications.

Requires both --app-name and --vendor unless querying by application IDs.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | filter by account ID |
| `--all` | bool | false | fetch all pages |
| `--app-name` | string | - | filter by application name |
| `--cursor` | string | - | pagination cursor |
| `--cve-id` | string | - | filter by CVE ID (contains) |
| `--limit` | int | 0 | max results per page (default 50) |
| `--severity` | stringSlice | - | filter by severity (CRITICAL, HIGH, MEDIUM, LOW) |
| `--site-id` | stringSlice | - | filter by site ID |
| `--vendor` | string | - | filter by vendor |

## applications list

List installed applications

```text
s1ctl applications list [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--agent-id` | stringSlice | - | filter by agent ID |
| `--all` | bool | false | fetch all pages |
| `--cursor` | string | - | pagination cursor |
| `--limit` | int | 0 | max results per page (default 50) |
| `--name` | string | - | filter by application name (contains) |
| `--publisher` | string | - | filter by publisher (contains) |
| `--site-id` | stringSlice | - | filter by site ID |
| `--version` | string | - | filter by version (contains) |

## applications risks

List application risks (CVE vulnerabilities per endpoint)

```text
s1ctl applications risks [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | filter by account ID |
| `--all` | bool | false | fetch all pages |
| `--cursor` | string | - | pagination cursor |
| `--limit` | int | 0 | max results per page (default 50) |
| `--severity` | stringSlice | - | filter by severity (CRITICAL, HIGH, MEDIUM, LOW) |
| `--site-id` | stringSlice | - | filter by site ID |
| `--vendor` | string | - | filter by vendor (contains) |
