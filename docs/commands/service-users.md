# service-users

Manage service users (API-token identities)

## service-users bulk-delete

Delete multiple service users by ID

```text
s1ctl service-users bulk-delete <service-user-id>... [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply the action (default: dry-run) |

## service-users create

Create a service user (generates an API token)

```text
s1ctl service-users create [flags]
```

Create a service user. Creation mints an API token that is shown once.

Scope must be one of: tenant, account, site. For account/site scope, pass
--scope-id (the account/site ID) and a role via --role-id or --role-name.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--description` | string | - | description |
| `--expiration` | string | - | token expiration, RFC3339 (required) |
| `--name` | string | - | service user name (required) |
| `--role-id` | string | - | RBAC role ID to assign |
| `--role-name` | string | - | predefined role name to assign |
| `--scope` | string | - | scope: tenant, account, site (required) |
| `--scope-id` | string | - | account/site ID for account/site scope |
| `--yes` | bool | false | apply the action (default: dry-run) |

## service-users delete

Delete a service user

```text
s1ctl service-users delete <service-user-id> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply the action (default: dry-run) |

## service-users export

Export service users

```text
s1ctl service-users export [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | filter by account ID |
| `--out` | string | - | write export to file (default: stdout) |
| `--query` | string | - | free text search (name, description) |
| `--role-id` | stringSlice | - | filter by RBAC role ID |
| `--site-id` | stringSlice | - | filter by site ID |

## service-users generate-token

Regenerate a service user's API token

```text
s1ctl service-users generate-token <service-user-id> [flags]
```

Regenerate the API token for a service user. The new token is shown once
and replaces any existing token.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--expiration` | string | - | token expiration, RFC3339 (required) |
| `--yes` | bool | false | apply the action (default: dry-run) |

## service-users get

Get service user details

```text
s1ctl service-users get <service-user-id>
```

## service-users list

List service users

```text
s1ctl service-users list [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | filter by account ID |
| `--all` | bool | false | fetch all pages |
| `--cursor` | string | - | pagination cursor |
| `--limit` | int | 0 | max results per page (default 50) |
| `--query` | string | - | free text search (name, description) |
| `--role-id` | stringSlice | - | filter by RBAC role ID |
| `--site-id` | stringSlice | - | filter by site ID |
| `--sort-by` | string | - | sort field (e.g. id, name) |
| `--sort-order` | string | - | sort direction (asc, desc) |

## service-users update

Update a service user

```text
s1ctl service-users update <service-user-id> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--description` | string | - | new description |
| `--role-id` | string | - | RBAC role ID to assign |
| `--role-name` | string | - | predefined role name to assign |
| `--scope` | string | - | new scope: tenant, account, site |
| `--scope-id` | string | - | account/site ID for account/site scope |
| `--yes` | bool | false | apply the action (default: dry-run) |
