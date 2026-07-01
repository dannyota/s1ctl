# devicecontrol

Device control rules

## devicecontrol list

List device control rules

```text
s1ctl devicecontrol list [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--all` | bool | false | fetch all pages |
| `--cursor` | string | - | pagination cursor |
| `--limit` | int | 0 | max results per page (default 50) |
| `--query` | string | - | free text search |
| `--site-id` | stringSlice | - | filter by site ID |

## devicecontrol pull

Pull device control rules to local YAML files

```text
s1ctl devicecontrol pull [flags]
```

Fetch all device control rules and write them as YAML files.

Each rule produces one file named by its sanitized rule name (e.g. block-usb-storage.yaml).
Server-only metadata (ID, scope, timestamps) is omitted from the YAML so the files
contain only the declarative rule definition.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | filter by account ID |
| `--out` | string | devicecontrol | output directory |
| `--site-id` | stringSlice | - | filter by site ID |

## devicecontrol push

Push device control rules from local YAML files

```text
s1ctl devicecontrol push [flags]
```

Read device control rule YAML files from a directory and sync them to SentinelOne.

Rules are matched by name: existing rules are updated, new rules are created.
Dry-run by default — pass --yes to apply changes.

New rules are created at the scope specified by --site-id. If no scope flag
is given, new rules are created at the global (tenant) scope.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--dir` | string | devicecontrol | directory containing device rule YAML files |
| `--site-id` | stringSlice | - | scope for new rules (default: global/tenant) |
| `--yes` | bool | false | apply changes (default: dry-run) |
