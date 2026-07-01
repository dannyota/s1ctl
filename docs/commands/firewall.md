# firewall

Firewall control rules

## firewall list

List firewall rules

```text
s1ctl firewall list [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--all` | bool | false | fetch all pages |
| `--cursor` | string | - | pagination cursor |
| `--limit` | int | 0 | max results per page (default 50) |
| `--query` | string | - | free text search |
| `--site-id` | stringSlice | - | filter by site ID |

## firewall pull

Pull firewall rules to local YAML files

```text
s1ctl firewall pull [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--out` | string | firewall | output directory |
| `--site-id` | stringSlice | - | filter by site ID |

## firewall push

Push firewall rules from local YAML files

```text
s1ctl firewall push [flags]
```

Read firewall rule YAML files from a directory and sync them to SentinelOne.
Rules are matched by name: existing rules are updated, new rules are created.
Dry-run by default — pass --yes to apply changes.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--dir` | string | firewall | directory containing firewall rule YAML files |
| `--site-id` | stringSlice | - | target site IDs |
| `--yes` | bool | false | apply changes (default: dry-run) |
