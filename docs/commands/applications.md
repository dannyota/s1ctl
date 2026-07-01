# applications

Application inventory

## applications list

List installed applications

```text
s1ctl applications list [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--agent-id` | stringSlice | - | filter by agent ID |
| `--all` | bool | false | fetch all pages |
| `--cursor` | string | - | pagination cursor |
| `--limit` | int | 0 | max results per page (default 50) |
| `--name` | string | - | filter by application name (contains) |
| `--publisher` | string | - | filter by publisher (contains) |
| `--site-id` | stringSlice | - | filter by site ID |
| `--version` | string | - | filter by version (contains) |
