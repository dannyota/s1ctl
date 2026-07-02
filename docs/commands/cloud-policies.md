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

## cloud-policies pull

Pull cloud security policies to local YAML files

```text
s1ctl cloud-policies pull [flags]
```

Fetch all cloud security policies and write them as YAML files.

Each policy produces one file carrying its ID, name, and status. Cloud policies
cannot be created through this surface, so push only reconciles status.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--out` | string | cloud-policies | output directory |

## cloud-policies push

Push cloud security policy status from local YAML files

```text
s1ctl cloud-policies push [flags]
```

Read cloud policy YAML files from a directory and reconcile their status.

Policies are matched by ID: a status change (enabled/disabled) is applied,
unchanged policies are skipped, and a local file whose ID has no live match
fails per-item since policies cannot be created through this surface. Dry-run by
default — pass --yes to apply changes.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--dir` | string | cloud-policies | directory containing cloud policy YAML files |
| `--yes` | bool | false | apply changes (default: dry-run) |
