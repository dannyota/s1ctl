# Firewall and network quarantine

Manage firewall control rules and network quarantine rules: CRUD, reorder,
copy between scopes, import/export, and sync as code with pull/push.

> Prerequisites: `s1ctl` installed and configured (`S1_CONSOLE_URL`, `S1_TOKEN`).

## Firewall control

### List firewall rules

```bash
s1ctl firewall list
s1ctl firewall list --site-id 000000
s1ctl firewall list --query "RDP" --all --json
```

| Flag | Description |
|------|-------------|
| `--site-id` | Filter by site ID (repeatable) |
| `--query` | Free text search |
| `--limit` | Max results per page (default 50) |
| `--all` | Fetch all pages |
| `--cursor` | Pagination cursor |

### Get a firewall rule

```bash
s1ctl firewall get 000000
s1ctl firewall get 000000 --json
```

### Enable and disable

Dry-run by default. Pass `--yes` to apply. Accepts multiple IDs.

```bash
s1ctl firewall enable 000000 --yes
s1ctl firewall disable 000000 000001 --yes
```

### Delete

```bash
s1ctl firewall delete 000000 --yes
```

### Reorder

Change rule evaluation order. Each argument is an `id:order` pair.

```bash
s1ctl firewall reorder 000000:1 000001:2 000002:3 --site-id 000000 --yes
```

| Flag | Description |
|------|-------------|
| `--site-id` | Scope: site IDs (repeatable) |
| `--account-id` | Scope: account IDs (repeatable) |
| `--group-id` | Scope: group IDs (repeatable) |
| `--yes` | Apply changes (default: dry-run) |

### Copy between scopes

```bash
s1ctl firewall copy \
  --source-site-id 111111 \
  --target-site-id 222222 --yes
```

| Flag | Description |
|------|-------------|
| `--source-site-id` | Source site IDs (repeatable) |
| `--source-account-id` | Source account IDs (repeatable) |
| `--target-site-id` | Target site ID |
| `--target-account-id` | Target account ID |
| `--target-group-id` | Target group ID |
| `--yes` | Apply changes (default: dry-run) |

### List available protocols

```bash
s1ctl firewall protocols
s1ctl firewall protocols --query "TCP"
```

### Export and import

Export rules to a JSON file, then import into another scope.

```bash
# Export
s1ctl firewall export --site-id 000000
s1ctl firewall export --site-id 000000 --out my-rules.json

# Import
s1ctl firewall import my-rules.json --site-id 222222          # dry-run
s1ctl firewall import my-rules.json --site-id 222222 --yes    # apply
```

### Config-as-code

Pull rules to YAML, review in git, push back. See
[Config-as-code](config-as-code.md) for the general pattern.

```bash
# Pull
s1ctl firewall pull --site-id 000000
s1ctl firewall pull --site-id 000000 --out snapshots/firewall

# Push
s1ctl firewall push --site-id 000000          # dry-run
s1ctl firewall push --site-id 000000 --yes    # apply
s1ctl firewall push --dir my-rules --yes
```

| Flag (pull) | Description |
|-------------|-------------|
| `--site-id` | Filter by site ID (repeatable) |
| `--out` | Output directory (default `firewall`) |

| Flag (push) | Description |
|-------------|-------------|
| `--dir` | Input directory (default `firewall`) |
| `--site-id` | Target site IDs (repeatable) |
| `--yes` | Apply changes (default: dry-run) |

---

## Network quarantine

The `network` command group manages network quarantine rules. These control
which network traffic is allowed when an agent is quarantined (isolated from
the network). The interface mirrors `firewall` with additional features for
location awareness, rule tagging, and scope movement.

### List network quarantine rules

```bash
s1ctl network list
s1ctl network list --site-id 000000 --all --json
```

| Flag | Description |
|------|-------------|
| `--site-id` | Filter by site ID (repeatable) |
| `--query` | Free text search |
| `--limit` | Max results per page (default 50) |
| `--all` | Fetch all pages |
| `--cursor` | Pagination cursor |

### Get a rule

```bash
s1ctl network get 000000
```

### Enable, disable, and delete

```bash
s1ctl network enable 000000 --yes
s1ctl network disable 000000 000001 --yes
s1ctl network delete 000000 --yes
```

### Reorder

```bash
s1ctl network reorder 000000:1 000001:2 --site-id 000000 --yes
```

### Copy and move

Copy duplicates rules to a new scope. Move reassigns existing rules.

```bash
# Copy
s1ctl network copy --source-site-id 111111 --target-site-id 222222 --yes

# Move
s1ctl network move 000000 000001 --target-site-id 222222 --yes
```

### Set location

Assign a location matcher to rules.

```bash
s1ctl network set-location 000000 --type all --yes
s1ctl network set-location 000000 --type specific --location-id 000001 --yes
s1ctl network set-location 000000 --type fallback --yes
```

| Flag | Description |
|------|-------------|
| `--type` | Location type: `all`, `specific`, or `fallback` (default `all`) |
| `--location-id` | Location IDs for `--type specific` (repeatable) |
| `--yes` | Apply changes (default: dry-run) |

### Tags

```bash
s1ctl network tags add 000000 --tag-id 000001 --yes
s1ctl network tags remove 000000 --tag-id 000001 --yes
```

### Configuration

Get or set the network quarantine control configuration for a scope.

```bash
s1ctl network configuration get --site-id 000000
s1ctl network configuration set --site-id 000000 --enabled --yes
s1ctl network configuration set --site-id 000000 \
  --location-aware --report-blocked --yes
```

| Flag (set) | Description |
|------------|-------------|
| `--site-id` | Scope: site IDs (repeatable) |
| `--account-id` | Scope: account IDs (repeatable) |
| `--group-id` | Scope: group IDs (repeatable) |
| `--enabled` | Enable network quarantine for the scope |
| `--location-aware` | Enable location awareness |
| `--report-blocked` | Report blocked events |
| `--selected-tag` | Selected tag IDs (repeatable) |
| `--yes` | Apply changes (default: dry-run) |

### Protocols

```bash
s1ctl network protocols
s1ctl network protocols --query "DNS"
```

### Export and import

```bash
s1ctl network export --site-id 000000
s1ctl network export --site-id 000000 --out nq-rules.json

s1ctl network import nq-rules.json --site-id 222222 --yes
```

### Config-as-code

```bash
# Pull
s1ctl network pull --site-id 000000
s1ctl network pull --site-id 000000 --out snapshots/nq

# Push
s1ctl network push --site-id 000000          # dry-run
s1ctl network push --site-id 000000 --yes    # apply
```

| Flag (pull) | Description |
|-------------|-------------|
| `--site-id` | Filter by site ID (repeatable) |
| `--out` | Output directory (default `network-quarantine`) |

| Flag (push) | Description |
|-------------|-------------|
| `--dir` | Input directory (default `network-quarantine`) |
| `--site-id` | Target site IDs (repeatable) |
| `--yes` | Apply changes (default: dry-run) |

## See also

- [Config-as-code](config-as-code.md) -- the pull/review/push loop
- [`firewall` command reference](../commands/firewall.md)
- [`network` command reference](../commands/network.md)
