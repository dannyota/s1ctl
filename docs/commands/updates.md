# updates

Manage agent update packages

## updates deploy

Manage Sentinel Deploy credential groups

```text
s1ctl updates deploy
```

Manage credential groups used by Sentinel Deploy (Ranger auto-deploy)
to install agents on unprotected endpoints.

Credential groups contain encrypted credentials that agents use to
authenticate when deploying to new endpoints.

## updates get

Get an update package

```text
s1ctl updates get <package-id>
```

## updates list

List update packages

```text
s1ctl updates list [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--all` | bool | false | fetch all pages |
| `--cursor` | string | - | pagination cursor |
| `--limit` | int | 0 | max results per page (default 50) |
| `--query` | string | - | free text search |
| `--site-id` | stringSlice | - | filter by site ID |
