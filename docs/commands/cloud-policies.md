# cloud-policies

Manage cloud security policies (CNS rules)

## cloud-policies delete

Delete cloud security policies

```text
s1ctl cloud-policies delete <id> [id...] [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply the action (default: dry-run) |

## cloud-policies disable

Disable cloud security policies

```text
s1ctl cloud-policies disable <id> [id...] [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply the action (default: dry-run) |

## cloud-policies enable

Enable cloud security policies

```text
s1ctl cloud-policies enable <id> [id...] [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply the action (default: dry-run) |

## cloud-policies get

Get cloud policy details

```text
s1ctl cloud-policies get <id>
```

## cloud-policies list

List cloud security policies

```text
s1ctl cloud-policies list [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--after` | string | - | pagination cursor |
| `--all` | bool | false | fetch all pages |
| `--limit` | int | 0 | max results per page (default 50) |
| `--severity` | stringSlice | - | filter by severity (HIGH, CRITICAL, etc.) |
| `--status` | stringSlice | - | filter by status |
