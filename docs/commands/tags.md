# tags

Manage tags

## tags list

List tags

```text
s1ctl tags list [flags]
```

List tags by type: firewall, network-quarantine, device-inventory.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--all` | bool | false | fetch all pages |
| `--cursor` | string | - | pagination cursor |
| `--limit` | int | 0 | max results per page (default 50) |
| `--query` | string | - | free text search |
| `--site-id` | stringSlice | - | filter by site ID |
| `--type` | string | - | tag type (firewall, network-quarantine, device-inventory) |
