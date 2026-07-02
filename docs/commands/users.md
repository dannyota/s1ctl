# users

Manage users

## users delete

Delete a user

```text
s1ctl users delete <user-id> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
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
