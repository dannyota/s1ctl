# iocs

Manage threat intelligence IOCs

## iocs config

Show threat intelligence configuration

```text
s1ctl iocs config
```

## iocs create

Create a threat intelligence IOC

```text
s1ctl iocs create [flags]
```

Create a new threat intelligence indicator of compromise.

Types: DNS, IPV4, IPV6, MD5, SHA1, SHA256, URL
Severities: Low, Medium, High

Dry-run by default; pass --yes to apply.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--description` | string | - | IOC description |
| `--external-id` | string | - | external reference ID |
| `--method` | string | - | detection method |
| `--name` | string | - | IOC name |
| `--severity` | string | - | severity (Low, Medium, High) |
| `--source` | string | - | intelligence source |
| `--type` | string | - | IOC type (DNS, IPV4, IPV6, MD5, SHA1, SHA256, URL) |
| `--valid-until` | string | - | expiration date (ISO 8601) |
| `--value` | string | - | indicator value |
| `--yes` | bool | false | apply the action (default: dry-run) |

## iocs delete

Delete threat intelligence IOCs

```text
s1ctl iocs delete <ioc-id...> [flags]
```

Delete one or more threat intelligence IOCs by ID.

Dry-run by default; pass --yes to apply.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply the action (default: dry-run) |

## iocs list

List threat intelligence IOCs

```text
s1ctl iocs list [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--all` | bool | false | fetch all pages |
| `--cursor` | string | - | pagination cursor |
| `--limit` | int | 0 | max results per page (default 50) |
| `--severity` | stringSlice | - | filter by severity (Low, Medium, High) |
| `--sort-by` | string | - | sort field |
| `--sort-order` | string | - | sort direction (asc, desc) |
| `--source` | stringSlice | - | filter by source |
| `--type` | stringSlice | - | filter by IOC type (DNS, IPV4, IPV6, MD5, SHA1, SHA256, URL) |
| `--value` | string | - | filter by IOC value |
