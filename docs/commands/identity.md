# identity

Manage Identity AD Service configuration, connectors, and ISPM

## identity ack-exposures

Acknowledge or unacknowledge ISPM exposures

```text
s1ctl identity ack-exposures [flags]
```

Set the acknowledged status on exposures. Requires --detection and --domain.
Dry-run by default — pass --yes to apply. Use --unack to reverse.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | filter by account ID |
| `--detection` | stringSlice | - | detection name(s) (required) |
| `--domain` | stringSlice | - | domain name(s) (required) |
| `--site-id` | stringSlice | - | filter by site ID |
| `--unack` | bool | false | reverse acknowledgement |
| `--yes` | bool | false | apply changes (default: dry-run) |

## identity config

Manage AD configurations

```text
s1ctl identity config
```

## identity connector

Manage AD connectors (Cloudlink agents)

```text
s1ctl identity connector
```

## identity domains

List AD domains

```text
s1ctl identity domains [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | filter by account ID |
| `--site-id` | stringSlice | - | filter by site ID |

## identity features

List available AD features

```text
s1ctl identity features [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | filter by account ID |
| `--site-id` | stringSlice | - | filter by site ID |

## identity onboard

Show AD service onboarding status

```text
s1ctl identity onboard [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | filter by account ID |
| `--site-id` | stringSlice | - | filter by site ID |

## identity skip-exposures

Skip or unskip ISPM exposures

```text
s1ctl identity skip-exposures [flags]
```

Set exposures as skipped (accepted risk) or unskip previously skipped exposures.
Requires --detection and --domain. Dry-run by default — pass --yes to apply.
Use --unskip to reverse a previous skip.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | filter by account ID |
| `--detection` | stringSlice | - | detection name(s) (required) |
| `--domain` | stringSlice | - | domain name(s) (required) |
| `--reason` | string | - | reason for skipping |
| `--site-id` | stringSlice | - | filter by site ID |
| `--unskip` | bool | false | reverse a previous skip |
| `--yes` | bool | false | apply changes (default: dry-run) |

## identity timezones

List available timezones for AD configuration

```text
s1ctl identity timezones [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | filter by account ID |
| `--site-id` | stringSlice | - | filter by site ID |
