# roles

Manage RBAC roles

## roles create

Create a role from a declarative role file

```text
s1ctl roles create --from-file <role.yaml> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | create in these account IDs |
| `--from-file` | string | - | role definition file, YAML or JSON (required) |
| `--group-id` | stringSlice | - | create in these group IDs |
| `--site-id` | stringSlice | - | create in these site IDs |
| `--tenant` | bool | false | create at the global (tenant) scope |
| `--yes` | bool | false | apply the action (default: dry-run) |

## roles delete

Delete a role

```text
s1ctl roles delete <role-id> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply the action (default: dry-run) |

## roles get

Get a role definition, including its permission tree

```text
s1ctl roles get <role-id>
```

## roles list

List RBAC roles

```text
s1ctl roles list [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--account-id` | stringSlice | - | filter by account ID |
| `--all` | bool | false | fetch all pages |
| `--cursor` | string | - | pagination cursor |
| `--group-id` | stringSlice | - | filter by group ID |
| `--limit` | int | 0 | max results per page (default 50) |
| `--predefined` | bool | false | filter by predefined roles (true) or custom roles (false) |
| `--query` | string | - | free text search (name, description) |
| `--site-id` | stringSlice | - | filter by site ID |
| `--sort-by` | string | - | sort field (e.g. name, createdAt) |
| `--sort-order` | string | - | sort direction (asc, desc) |

## roles template

Print the blank role template for editing

```text
s1ctl roles template
```

Fetch the blank role template (description and the full permission tree with
default values) and print it as JSON. Use it as a starting point for a new
role: edit the values, then create the role with 'roles create --from-file'.

## roles update

Update a role from a declarative role file

```text
s1ctl roles update <role-id> --from-file <role.yaml> [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--from-file` | string | - | role definition file, YAML or JSON (required) |
| `--yes` | bool | false | apply the action (default: dry-run) |
