# users

Manage users

## users 2fa

Enable or disable two-factor authentication for a user

```text
s1ctl users 2fa
```

## users delete

Delete a user

```text
s1ctl users delete <user-id> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply the action (default: dry-run) |

## users generate-token

Generate an API token for the current user

```text
s1ctl users generate-token [flags]
```

Generate an API token for the authenticated user. The token is shown
once and replaces any existing token for that user.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--force-legacy` | bool | false | request a legacy token even when auth-tokens is enabled |
| `--yes` | bool | false | apply the action (default: dry-run) |

## users get

Get user details

```text
s1ctl users get <user-id>
```

## users list

List users

```text
s1ctl users list [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--all` | bool | false | fetch all pages |
| `--cursor` | string | - | pagination cursor |
| `--limit` | int | 0 | max results per page (default 50) |
| `--query` | string | - | free text search |
| `--sort-by` | string | - | sort field (e.g. fullName, email) |
| `--sort-order` | string | - | sort direction (asc, desc) |

## users revoke-token

Revoke a user's API token

```text
s1ctl users revoke-token <user-id> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply the action (default: dry-run) |

## users token-details

Show API-token metadata (created/expires) for a user

```text
s1ctl users token-details [<user-id>]
```

Show API-token metadata for a user. With no argument, reports the
authenticated user's token; with a user ID, reports that user's token. Only
timestamps are shown — any secret value is redacted.

## users update

Update a user

```text
s1ctl users update <user-id> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--email` | string | - | new email address |
| `--full-name` | string | - | new full name |
| `--scope` | string | - | new scope |
| `--yes` | bool | false | apply the action (default: dry-run) |
