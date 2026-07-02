# accounts

Manage accounts

## accounts count

Count accounts

```text
s1ctl accounts count
```

## accounts expire

Expire an account immediately

```text
s1ctl accounts expire <account-id> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply the action (default: dry-run) |

## accounts get

Get account details

```text
s1ctl accounts get <account-id>
```

## accounts list

List accounts

```text
s1ctl accounts list [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--all` | bool | false | fetch all pages |
| `--cursor` | string | - | pagination cursor |
| `--limit` | int | 0 | max results per page (default 50) |
| `--query` | string | - | free text search |
| `--state` | stringSlice | - | filter by state |

## accounts reactivate

Reactivate an expired account

```text
s1ctl accounts reactivate <account-id> [flags]
```

Reactivate an expired account. Specify exactly one of --unlimited (no
expiration) or --expiration (an RFC3339 timestamp) to set the new license
window.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--expiration` | string | - | new expiration as an RFC3339 timestamp |
| `--unlimited` | bool | false | reactivate with no expiration |
| `--yes` | bool | false | apply the action (default: dry-run) |

## accounts uninstall-password

Manage an account's agent uninstall password

```text
s1ctl accounts uninstall-password
```
