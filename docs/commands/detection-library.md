# detection-library

Manage platform detection rules (detection library)

## detection-library data-sources

List available detection data sources

```text
s1ctl detection-library data-sources
```

## detection-library disable

Disable platform detection rules

```text
s1ctl detection-library disable <rule-id>... [flags]
```

Disable one or more platform detection rules by ID.

If --scope-level is not specified, the rule's inherited scope is
auto-detected. Platform rules inherited from a higher scope (e.g.
account) can only be toggled at that scope — not at site or group.

Dry-run by default — pass --yes to apply.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--scope-id` | string | - | account, site, or group ID (auto-detected if omitted) |
| `--scope-level` | string | - | scope level (auto-detected from rule if omitted) |
| `--yes` | bool | false | apply changes (default: dry-run) |

## detection-library enable

Enable platform detection rules

```text
s1ctl detection-library enable <rule-id>... [flags]
```

Enable one or more platform detection rules by ID.

If --scope-level is not specified, the rule's inherited scope is
auto-detected. Platform rules inherited from a higher scope (e.g.
account) can only be toggled at that scope — not at site or group.

Dry-run by default — pass --yes to apply.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--scope-id` | string | - | account, site, or group ID (auto-detected if omitted) |
| `--scope-level` | string | - | scope level (auto-detected from rule if omitted) |
| `--yes` | bool | false | apply changes (default: dry-run) |

## detection-library list

List platform detection rules

```text
s1ctl detection-library list [flags]
```

List platform detection rules from the detection library.

Requires --scope (global, account, site, group) and --scope-id for non-global scopes.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--all` | bool | false | fetch all pages |
| `--category` | stringSlice | - | filter by category (Events, Correlation, UEBAFirstSeen, Scheduled) |
| `--cursor` | string | - | pagination cursor |
| `--limit` | int | 0 | max results per page (default 50) |
| `--name` | string | - | filter by rule name (substring match) |
| `--scope` | string | - | scope level: global, account, site, group (required) |
| `--scope-id` | string | - | account, site, or group ID for scoped listing |
| `--severity` | stringSlice | - | filter by severity (Info, Low, Medium, High, Critical) |
| `--source` | stringSlice | - | filter by data source |
| `--status` | stringSlice | - | filter by status (Active, Disabled, Activating, Disabling) |
| `--surface` | stringSlice | - | filter by attack surface |
| `--tag` | stringSlice | - | filter by tag |

## detection-library surfaces

List available detection surfaces

```text
s1ctl detection-library surfaces
```
