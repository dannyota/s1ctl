# policies

View endpoint policies

## policies diff

Compare policies across sites

```text
s1ctl policies diff [flags]
```

Fetch policies for all sites (or a filtered subset) and highlight
fields that differ between them. Useful for spotting inconsistencies
like one site in detect mode while others are in protect mode.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | filter by account ID |
| `--site-id` | stringSlice | - | filter by site ID |

## policies get

Get policy for a scope (site, account, or group)

```text
s1ctl policies get [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | string | - | account ID |
| `--group-id` | string | - | group ID (requires --site-id) |
| `--site-id` | string | - | site ID |

## policies list

List policies across sites

```text
s1ctl policies list [flags]
```

List endpoint policies across all sites (or a filtered subset).

The SentinelOne API returns one policy per scope. This command fetches sites
and retrieves each site's policy, presenting them side by side for comparison.

Use --account-id or --site-id to narrow the scope.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | filter by account ID |
| `--site-id` | stringSlice | - | filter by site ID |

## policies pull

Pull site policies to local YAML files

```text
s1ctl policies pull [flags]
```

Fetch endpoint policies for each site and write them as YAML files.

Each site produces one file named by its sanitized site name (e.g. production.yaml).
The YAML includes site metadata (siteId, siteName) for push matching, plus the key
policy fields: mitigationMode, antiTamperingOn, networkQuarantineOn, etc.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | filter by account ID |
| `--out` | string | policies | output directory |
| `--site-id` | stringSlice | - | filter by site ID |

## policies push

Push site policies from local YAML files

```text
s1ctl policies push [flags]
```

Read policy YAML files from a directory and update the corresponding site policies.

Each file must contain a siteId field to identify the target site. The command
fetches the current policy, diffs it against the desired state, and applies changes.
Dry-run by default — pass --yes to apply changes.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--dir` | string | policies | directory containing policy YAML files |
| `--yes` | bool | false | apply changes (default: dry-run) |
