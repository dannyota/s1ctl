# locations

Manage firewall locations

## locations create

Create a firewall location

```text
s1ctl locations create --name <name> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | create in these account IDs |
| `--description` | string | - | location description |
| `--name` | string | - | location name (required) |
| `--operator` | string | any | match operator: all, any, none |
| `--site-id` | stringSlice | - | create in these site IDs |
| `--yes` | bool | false | apply the action (default: dry-run) |

## locations delete

Delete a firewall location

```text
s1ctl locations delete <location-id> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply the action (default: dry-run) |

## locations list

List firewall locations

```text
s1ctl locations list [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | filter by account ID |
| `--all` | bool | false | fetch all pages |
| `--cursor` | string | - | pagination cursor |
| `--group-id` | stringSlice | - | filter by group ID |
| `--limit` | int | 0 | max results per page (default 50) |
| `--site-id` | stringSlice | - | filter by site ID |

## locations pull

Pull locations to local YAML files

```text
s1ctl locations pull [flags]
```

Fetch all locations and write them as YAML files.

Each location produces one file named by its sanitized name. Server-only
metadata (ID, scope, counters, timestamps) is omitted so the files contain only
the declarative definition, including the detection parameters.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--out` | string | locations | output directory |
| `--site-id` | stringSlice | - | filter by site ID |

## locations push

Push locations from local YAML files

```text
s1ctl locations push [flags]
```

Read location YAML files from a directory and sync them to SentinelOne.

Locations are matched by name: existing locations are updated, new ones are
created, and unchanged ones are skipped. Dry-run by default — pass --yes to
apply. New locations are created at the scope given by --site-id (default:
global/tenant).

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--dir` | string | locations | directory containing location YAML files |
| `--site-id` | stringSlice | - | scope for new locations (default: global/tenant) |
| `--yes` | bool | false | apply changes (default: dry-run) |

## locations update

Update a firewall location's name, description, or operator

```text
s1ctl locations update <location-id> --name <name> [flags]
```

Update a location's name, description, and match operator.

Note: this replaces the location definition with the supplied fields; detection
parameters not expressible as flags are managed through 'locations push'.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--description` | string | - | location description |
| `--name` | string | - | location name (required) |
| `--operator` | string | any | match operator: all, any, none |
| `--yes` | bool | false | apply the action (default: dry-run) |
