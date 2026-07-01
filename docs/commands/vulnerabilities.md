# vulnerabilities

Manage xSPM vulnerabilities

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

Update vulnerability analyst verdict

```text
s1ctl vulnerabilities verdict <id> <verdict> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply the action (default: dry-run) |
