# settings

Manage platform settings

## settings get

Get settings configuration

```text
s1ctl settings get <type> [flags]
```

Get configuration for a specific settings type.

Types: notifications, sso, smtp, syslog

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | filter by account ID |
| `--site-id` | stringSlice | - | filter by site ID |

## settings list

List settings categories

```text
s1ctl settings list
```

## settings test

Test settings connectivity

```text
s1ctl settings test <type> [flags]
```

Test connectivity for SMTP or syslog settings.

Types: smtp, syslog

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | filter by account ID |
| `--site-id` | stringSlice | - | filter by site ID |
| `--yes` | bool | false | apply the action (default: dry-run) |
