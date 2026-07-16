# Blocklist

Manage the SentinelOne blocklist (restrictions): add, update, validate,
export, and sync blocked hashes as code.

The blocklist holds SHA1/SHA256 hashes that agents block from executing.
Items are scoped globally (tenant) or to accounts, sites, or groups.

> Prerequisites: `s1ctl` installed and configured (`S1_CONSOLE_URL`, `S1_TOKEN`).

## List blocklist items

```bash
s1ctl blocklist list
s1ctl blocklist list --site-id 000000
s1ctl blocklist list --os-type windows --query "mimikatz"
s1ctl blocklist list --all --json
```

| Flag | Description |
|------|-------------|
| `--site-id` | Filter by site ID (repeatable) |
| `--account-id` | Filter by account ID (repeatable) |
| `--group-id` | Filter by group ID (repeatable) |
| `--os-type` | Filter by OS: `windows`, `linux`, `macos`, `windows_legacy` (repeatable) |
| `--query` | Free text search |
| `--value` | Filter by hash value |
| `--sort-by` | Sort field (e.g. `createdAt`, `osType`) |
| `--sort-order` | Sort direction (`asc`, `desc`) |
| `--limit` | Max results per page (default 50) |
| `--all` | Fetch all pages |
| `--cursor` | Pagination cursor |

## Create a blocklist entry

Add a hash to the blocklist. Dry-run by default.

```bash
s1ctl blocklist create \
  --value da39a3ee5e6b4b0d3255bfef95601890afd80709 \
  --os-type windows \
  --description "Known malware sample"         # dry-run

s1ctl blocklist create \
  --value da39a3ee5e6b4b0d3255bfef95601890afd80709 \
  --sha256 e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855 \
  --os-type linux \
  --site-id 000000 \
  --yes                                        # apply
```

| Flag | Description |
|------|-------------|
| `--value` | SHA1 hash to block (required) |
| `--sha256` | SHA256 hash to block |
| `--os-type` | Target OS (required): `windows`, `linux`, `macos`, `windows_legacy` |
| `--description` | Item description |
| `--source` | Item source |
| `--type` | Restriction type (default `black_hash`) |
| `--site-id` | Target site IDs (repeatable) |
| `--account-id` | Target account IDs (repeatable) |
| `--group-id` | Target group IDs (repeatable) |
| `--yes` | Apply (default: dry-run) |

## Update a blocklist entry

Full replacement of an existing item. Dry-run by default.

```bash
s1ctl blocklist update 000000 \
  --value da39a3ee5e6b4b0d3255bfef95601890afd80709 \
  --os-type windows \
  --description "Updated description" --yes
```

## Delete a blocklist entry

```bash
s1ctl blocklist delete 000000          # dry-run
s1ctl blocklist delete 000000 --yes    # apply
```

## Validate a hash

Check whether a hash is on SentinelOne's "Not Allowed" or "Not Recommended"
list before adding it. Read-only.

```bash
s1ctl blocklist validate --value da39a3ee5e6b4b0d3255bfef95601890afd80709
s1ctl blocklist validate --sha256 e3b0c44298fc1c149afbf4c8996fb924... --os-type windows
s1ctl blocklist validate --value da39a3ee5e6b... --site-id 000000
```

| Flag | Description |
|------|-------------|
| `--value` | SHA1 hash to validate |
| `--sha256` | SHA256 hash to validate |
| `--os-type` | Target OS |
| `--site-id` | Scope site IDs (repeatable) |
| `--account-id` | Scope account IDs (repeatable) |
| `--group-id` | Scope group IDs (repeatable) |

## Export as CSV

```bash
s1ctl blocklist export --site-id 000000 --out blocklist.csv
s1ctl blocklist export --tenant --out global.csv
s1ctl blocklist export --os-type windows
```

| Flag | Description |
|------|-------------|
| `--site-id` | Filter by site ID (repeatable) |
| `--account-id` | Filter by account ID (repeatable) |
| `--group-id` | Filter by group ID (repeatable) |
| `--os-type` | Filter by OS type (repeatable) |
| `--tenant` | Export the global (tenant) blocklist |
| `--out` | Output file (default: stdout) |

## Config-as-code

Pull blocklist items to local YAML files, review in git, and push back. See
[Config-as-code](config-as-code.md) for the general pattern.

### Pull

```bash
s1ctl blocklist pull --site-id 000000
s1ctl blocklist pull --out snapshots/blocklist
```

| Flag | Description |
|------|-------------|
| `--site-id` | Filter by site ID (repeatable) |
| `--out` | Output directory (default `blocklist`) |

### Push

Items are matched by type + OS + value: existing items are updated, new items
are created, unchanged items are skipped. Dry-run by default.

```bash
s1ctl blocklist push --site-id 000000          # dry-run
s1ctl blocklist push --site-id 000000 --yes    # apply
s1ctl blocklist push --dir my-blocklist --yes
```

| Flag | Description |
|------|-------------|
| `--dir` | Input directory (default `blocklist`) |
| `--site-id` | Scope for new items (default: global/tenant) |
| `--yes` | Apply changes (default: dry-run) |

### Pull, diff, push workflow

```bash
# 1. Pull
s1ctl blocklist pull --site-id 000000

# 2. Review / edit
git diff blocklist/

# 3. Commit
git add blocklist/ && git commit -m "blocklist: add hash for site 000000"

# 4. Push
s1ctl blocklist push --site-id 000000 --yes
```

## Workflows

### Validate then add

```bash
s1ctl blocklist validate --value da39a3ee5e6b... --os-type windows
s1ctl blocklist create --value da39a3ee5e6b... --os-type windows --yes
```

### Copy blocklist between sites

```bash
s1ctl blocklist pull --site-id 111111
s1ctl blocklist push --site-id 222222          # dry-run
s1ctl blocklist push --site-id 222222 --yes    # apply
```

### Count entries by OS

```bash
s1ctl blocklist list --all --json \
  | jq 'group_by(.osType) | map({os: .[0].osType, count: length})'
```

## See also

- [Config-as-code](config-as-code.md) -- the pull/review/push loop
- [Exclusions](exclusions.md) -- allowlist entries (the inverse of blocklist)
- [`blocklist` command reference](../commands/blocklist.md)
