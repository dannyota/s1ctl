# updates

Manage agent update packages

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
