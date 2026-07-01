# exclusions

Manage exclusions and blocklist

## exclusions create

Create an exclusion

```text
s1ctl exclusions create [flags]
```

Create a new exclusion entry.

Types: path, file_type, white_hash, browser, certificate, document_type
OS types: windows, linux, macos, windows_legacy
Modes: suppress, suppress_dynamic_only, suppress_app_control

For path exclusions, --path-type specifies the match type:
  subfolders (default), file, glob

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--description` | string | - | exclusion description |
| `--group-id` | stringSlice | - | target group IDs |
| `--mode` | string | suppress | exclusion mode (suppress, suppress_dynamic_only, suppress_app_control) |
| `--os-type` | string | - | target OS (windows, linux, macos) |
| `--path-type` | string | - | path exclusion type (subfolders, file, glob) |
| `--site-id` | stringSlice | - | target site IDs |
| `--type` | string | - | exclusion type (path, file_type, white_hash, browser, certificate, document_type) |
| `--value` | string | - | exclusion value (path, hash, extension, etc.) |
| `--yes` | bool | false | apply the action (default: dry-run) |

## exclusions get

Get exclusion details

```text
s1ctl exclusions get <exclusion-id>
```

## exclusions list

List exclusions

```text
s1ctl exclusions list [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--all` | bool | false | fetch all pages |
| `--cursor` | string | - | pagination cursor |
| `--limit` | int | 0 | max results per page (default 50) |
| `--os-type` | stringSlice | - | filter by OS type |
| `--query` | string | - | free text search |
| `--site-id` | stringSlice | - | filter by site ID |
| `--sort-by` | string | - | sort field (e.g. type, osType) |
| `--sort-order` | string | - | sort direction (asc, desc) |
| `--type` | stringSlice | - | filter by exclusion type |

## exclusions pull

Pull exclusions to local files

```text
s1ctl exclusions pull [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--out` | string | samples | output directory |
| `--site-id` | stringSlice | - | filter by site ID |

## exclusions push

Push exclusions from local file

```text
s1ctl exclusions push [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--file` | string | samples/exclusions.json | input file |
| `--site-id` | stringSlice | - | target site IDs |
| `--yes` | bool | false | apply changes (default: dry-run) |
