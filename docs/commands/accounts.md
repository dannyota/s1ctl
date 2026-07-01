# accounts

Manage accounts

## accounts count

Count accounts

```text
s1ctl accounts count
```

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
