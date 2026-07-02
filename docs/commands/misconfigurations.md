# misconfigurations

Manage xSPM misconfigurations

## misconfigurations get

Get misconfiguration details

```text
s1ctl misconfigurations get <id>
```

## misconfigurations list

List misconfigurations

```text
s1ctl misconfigurations list [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--after` | string | - | pagination cursor |
| `--all` | bool | false | fetch all pages |
| `--limit` | int | 0 | max results per page (default 50) |
| `--severity` | stringSlice | - | filter by severity (HIGH, CRITICAL, etc.) |
| `--status` | stringSlice | - | filter by status |

## misconfigurations status

Update misconfiguration status

```text
s1ctl misconfigurations status <id> <status> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply the action (default: dry-run) |

## misconfigurations verdict

Update misconfiguration analyst verdict (TRUE_POSITIVE, FALSE_POSITIVE, SUSPICIOUS, UNDEFINED)

```text
s1ctl misconfigurations verdict <id> <verdict> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply the action (default: dry-run) |
