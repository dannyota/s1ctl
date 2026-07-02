# dlp

Manage Data Loss Prevention (DLP) rules and classifications

## dlp classifications

Manage DLP classifications

```text
s1ctl dlp classifications
```

## dlp rules

Manage data protection rules

```text
s1ctl dlp rules
```

## dlp settings

Show DLP engine settings for a scope

```text
s1ctl dlp settings [flags]
```

Show DLP engine settings. A scope is required by the API: pass both
--scope-level and --scope-id.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--scope-id` | string | - | account, site, or group ID |
| `--scope-level` | string | - | scope level (account, site, group) |
