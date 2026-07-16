# assets

Manage XDR asset inventory

## assets action

Perform an action on assets

```text
s1ctl assets action [flags]
```

Perform an action on one or more assets.

Action names are passed through to the API (e.g. mark_asset_criticality_high,
mark_asset_criticality_medium). Dry-run by default; pass --yes to apply.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--action` | string | - | action name (required) |
| `--id` | stringSlice | - | asset ID(s) to act on (required, repeatable) |
| `--type` | string | - | asset type slug (omit for cross-type action) |
| `--yes` | bool | false | apply the action (default: dry-run) |

## assets categories

List asset categories with counts

```text
s1ctl assets categories [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | filter by account ID |
| `--group-id` | stringSlice | - | filter by group ID |
| `--site-id` | stringSlice | - | filter by site ID |

## assets export

Export assets from the XDR inventory

```text
s1ctl assets export [flags]
```

Export assets as raw CSV or JSON from the API.

Streams the raw export response to a file or stdout.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | filter by account ID |
| `--filter` | stringArray | - | key=value filter |
| `--group-id` | stringSlice | - | filter by group ID |
| `--output-file` | string | - | write export to file instead of stdout |
| `--site-id` | stringSlice | - | filter by site ID |
| `--type` | string | - | asset type slug (e.g. device, server) |

## assets filter-options

Show available filter fields for an asset type

```text
s1ctl assets filter-options [flags]
```

Show available filters for the given asset type.

Combines autocomplete and free-text filter information.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | filter by account ID |
| `--group-id` | stringSlice | - | filter by group ID |
| `--site-id` | stringSlice | - | filter by site ID |
| `--type` | string | - | asset type slug (required) |

## assets list

List assets from the XDR inventory

```text
s1ctl assets list [flags]
```

List assets from the XDR asset inventory.

When --type is omitted, lists assets across all types.
Use --filter key=value to pass type-specific API query parameters.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | filter by account ID |
| `--all` | bool | false | fetch all pages |
| `--cursor` | string | - | pagination cursor |
| `--filter` | stringArray | - | key=value filter (e.g. --filter osTypes=windows) |
| `--group-id` | stringSlice | - | filter by group ID |
| `--limit` | int | 0 | max results per page |
| `--site-id` | stringSlice | - | filter by site ID |
| `--skip` | int | 0 | number of results to skip |
| `--sort-by` | string | - | sort field |
| `--sort-order` | string | - | sort direction (asc, desc) |
| `--type` | string | - | asset type slug (e.g. device, server, surface/cloud) |

## assets notes

Manage asset notes

```text
s1ctl assets notes
```

## assets overview

Show asset counts by category and surface

```text
s1ctl assets overview [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | filter by account ID |
| `--group-id` | stringSlice | - | filter by group ID |
| `--site-id` | stringSlice | - | filter by site ID |

## assets sub-categories

List asset sub-categories

```text
s1ctl assets sub-categories [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | filter by account ID |
| `--group-id` | stringSlice | - | filter by group ID |
| `--site-id` | stringSlice | - | filter by site ID |
