# upgrade-policies

Manage agent auto-upgrade policies

## upgrade-policies activate

Activate an upgrade policy

```text
s1ctl upgrade-policies activate <policy-id> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply the action (default: dry-run) |

## upgrade-policies create

Create an upgrade policy

```text
s1ctl upgrade-policies create [flags]
```

Create a new agent auto-upgrade policy.

Scope levels: account, group, site, tenant
OS types: linux, macos, windows

Use "upgrade-policies packages" to find available package versions and file IDs.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--active` | bool | false | activate the policy immediately |
| `--all-endpoints` | bool | true | apply to all endpoints (set false with tags) |
| `--build` | string | - | package build version |
| `--description` | string | - | policy description |
| `--file-id` | string | - | package file ID (required; see 'upgrade-policies packages') |
| `--major` | string | - | package major version |
| `--max-retries` | int | 5 | max upgrade retries on failure |
| `--minor` | string | - | package minor version |
| `--name` | string | - | policy name (required) |
| `--os-type` | string | - | OS type: linux, macos, windows (required) |
| `--scope-id` | string | - | scope ID |
| `--scope-level` | string | - | scope level: account, group, site, tenant (required) |
| `--tag` | stringSlice | - | endpoint tags (when --all-endpoints=false) |
| `--yes` | bool | false | apply the action (default: dry-run) |

## upgrade-policies deactivate

Deactivate an upgrade policy

```text
s1ctl upgrade-policies deactivate <policy-id> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply the action (default: dry-run) |

## upgrade-policies delete

Delete an upgrade policy

```text
s1ctl upgrade-policies delete <policy-id> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply the action (default: dry-run) |

## upgrade-policies get

Get upgrade policy details

```text
s1ctl upgrade-policies get <policy-id> [flags]
```

Get details for a single upgrade policy by ID.

The API requires scope and OS filters even for a single lookup.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--os-type` | string | - | OS type (linux, macos, windows) [required] |
| `--scope-id` | string | - | scope ID |
| `--scope-level` | string | - | scope level (account, group, site, tenant) [required] |

## upgrade-policies list

List upgrade policies

```text
s1ctl upgrade-policies list [flags]
```

List agent auto-upgrade policies for a given scope and OS type.

Scope levels: account, group, site, tenant
OS types: linux, macos, windows

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--all` | bool | false | fetch all pages |
| `--limit` | int | 0 | max results per page (default 50) |
| `--os-type` | string | - | OS type (linux, macos, windows) [required] |
| `--scope-id` | string | - | scope ID |
| `--scope-level` | string | - | scope level (account, group, site, tenant) [required] |
| `--skip` | int | 0 | skip first N results |
| `--sort-by` | string | - | sort field (default: priority) |
| `--sort-order` | string | - | sort direction (asc, desc; default: asc) |

## upgrade-policies packages

List available upgrade packages

```text
s1ctl upgrade-policies packages [flags]
```

List agent packages available for upgrade policies.

Scope levels: account, group, site, tenant
OS types: linux, macos, windows

Each package may include multiple file variants. Use the file ID
when creating an upgrade policy (--file-id).

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--os-type` | string | - | OS type (linux, macos, windows) [required] |
| `--query` | string | - | filter by display name (partial match) |
| `--scope-id` | string | - | scope ID |
| `--scope-level` | string | - | scope level (account, group, site, tenant) [required] |

## upgrade-policies pull

Pull upgrade policies to local YAML files

```text
s1ctl upgrade-policies pull [flags]
```

Fetch upgrade policies and write them as YAML files.

Each policy produces one file named by its sanitized name. Server-only metadata
(ID, priority, timestamps) is omitted so the files contain only the declarative
definition.

The API requires --scope-level and --os-type to list policies. Pull fetches the
specified partition; use multiple invocations for different OS types or scopes.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--os-type` | string | - | OS type (linux, macos, windows) [required] |
| `--out` | string | upgrade-policies | output directory |
| `--scope-id` | string | - | scope ID (site/account/group ID) |
| `--scope-level` | string | - | scope level (account, group, site, tenant) [required] |

## upgrade-policies push

Push upgrade policies from local YAML files

```text
s1ctl upgrade-policies push [flags]
```

Read upgrade policy YAML files from a directory and sync them to SentinelOne.

Policies are matched by name: existing policies are updated, new ones are
created, and unchanged ones are skipped. Dry-run by default — pass --yes to
apply changes.

The API requires --scope-level and --os-type to list live policies for matching.
These must match the scope and OS in the local files.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--dir` | string | upgrade-policies | directory containing upgrade policy YAML files |
| `--os-type` | string | - | OS type (linux, macos, windows) [required] |
| `--scope-id` | string | - | scope ID (site/account/group ID) |
| `--scope-level` | string | - | scope level (account, group, site, tenant) [required] |
| `--yes` | bool | false | apply changes (default: dry-run) |

## upgrade-policies update

Update an upgrade policy

```text
s1ctl upgrade-policies update <policy-id> [flags]
```

Update an existing agent auto-upgrade policy.

The full policy body is sent, so provide every flag as with "create".

Scope levels: account, group, site, tenant
OS types: linux, macos, windows

Use "upgrade-policies packages" to find available package versions and file IDs.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--active` | bool | false | activate the policy immediately |
| `--all-endpoints` | bool | true | apply to all endpoints (set false with tags) |
| `--build` | string | - | package build version |
| `--description` | string | - | policy description |
| `--file-id` | string | - | package file ID (required; see 'upgrade-policies packages') |
| `--major` | string | - | package major version |
| `--max-retries` | int | 5 | max upgrade retries on failure |
| `--minor` | string | - | package minor version |
| `--name` | string | - | policy name (required) |
| `--os-type` | string | - | OS type: linux, macos, windows (required) |
| `--scope-id` | string | - | scope ID |
| `--scope-level` | string | - | scope level: account, group, site, tenant (required) |
| `--tag` | stringSlice | - | endpoint tags (when --all-endpoints=false) |
| `--yes` | bool | false | apply the action (default: dry-run) |
