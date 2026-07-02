# unified-exclusions

Manage unified exclusions

## unified-exclusions create

Create a unified exclusion

```text
s1ctl unified-exclusions create [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--description` | string | - | exclusion description |
| `--engines` | string | - | engines |
| `--interaction-level` | string | - | interaction level |
| `--mode-type` | string | - | mode type (required) |
| `--name` | string | - | exclusion name (required) |
| `--os-type` | string | - | target OS type (required) |
| `--path-type` | string | - | path exclusion type |
| `--reason` | string | - | exclusion reason (required) |
| `--scope-id` | string | - | scope level ID |
| `--scope-level` | string | - | scope level (required) |
| `--source` | string | - | exclusion source |
| `--threat-type` | string | - | threat type (required) |
| `--type` | string | - | exclusion type |
| `--value` | string | - | exclusion value |
| `--yes` | bool | false | apply the action (default: dry-run) |

## unified-exclusions export

Export unified exclusions

```text
s1ctl unified-exclusions export [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | filter by account ID |
| `--group-id` | stringSlice | - | filter by group ID |
| `--mode-type` | stringSlice | - | filter by mode type |
| `--os-type` | stringSlice | - | filter by OS type |
| `--out` | string | - | write export to file (default: stdout) |
| `--site-id` | stringSlice | - | filter by site ID |
| `--source` | stringSlice | - | filter by source |
| `--threat-type` | stringSlice | - | filter by threat type |

## unified-exclusions list

List unified exclusions

```text
s1ctl unified-exclusions list [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | filter by account ID |
| `--all` | bool | false | fetch all pages |
| `--cursor` | string | - | pagination cursor |
| `--engines` | stringSlice | - | filter by engines |
| `--group-id` | stringSlice | - | filter by group ID |
| `--limit` | int | 0 | max results per page (default 50) |
| `--mode-type` | stringSlice | - | filter by mode type |
| `--name` | stringSlice | - | filter by name (contains) |
| `--os-type` | stringSlice | - | filter by OS type |
| `--site-id` | stringSlice | - | filter by site ID |
| `--sort-by` | string | - | sort field |
| `--sort-order` | string | - | sort direction (asc, desc) |
| `--source` | stringSlice | - | filter by source |
| `--threat-type` | stringSlice | - | filter by threat type |
| `--value` | stringSlice | - | filter by value (contains) |
