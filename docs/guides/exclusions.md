# Exclusions

Manage SentinelOne exclusions (allowlist entries) as code: list, inspect, pull
to local files, review in git, and push back.

## The config-as-code loop

```mermaid
flowchart LR
    A[Console] -->|pull| B[Local JSON]
    B -->|git diff| C[Review]
    C -->|push --yes| A
    C -->|git commit| D[Git history]
```

Pull the live exclusions to a local file, review the diff, commit to git for
audit history, then push changes back. Every mutation is dry-run by default --
pass `--yes` to apply.

## List exclusions

```bash
s1ctl exclusions list
```

Filter by site, type, or OS:

```bash
s1ctl exclusions list --site-id 000000
s1ctl exclusions list --type path --os-type windows
s1ctl exclusions list --query "chrome" --limit 20
s1ctl exclusions list --all --json
```

| Flag | Description |
|------|-------------|
| `--site-id` | Filter by site ID (repeatable) |
| `--type` | Filter by exclusion type (repeatable): `path`, `white_hash`, `certificate`, `browser`, `file_type` |
| `--os-type` | Filter by OS type (repeatable): `windows`, `macos`, `linux` |
| `--query` | Free text search |
| `--sort-by` | Sort field (e.g. `type`, `osType`) |
| `--sort-order` | Sort direction (`asc`, `desc`) |
| `--limit` | Max results per page (default 50) |
| `--all` | Fetch all pages |
| `--cursor` | Pagination cursor for manual paging |
| `--json` | Machine-readable output |

Table columns: ID, Type, Value, OS, Mode.

## Get exclusion details

```bash
s1ctl exclusions get 000000
```

Returns: ID, type, value, OS, mode, description, scope name, user, and
creation date.

## Pull exclusions

Download all exclusions matching the given filters to a local JSON file:

```bash
s1ctl exclusions pull --site-id 000000
```

| Flag | Description |
|------|-------------|
| `--site-id` | Filter by site ID (repeatable) |
| `--out` | Output directory (default `samples`) |

The command auto-paginates (fetches all pages) and writes to
`<out>/exclusions.json`.

```bash
# Pull to the default samples/ directory
s1ctl exclusions pull --site-id 000000

# Pull to a custom directory
s1ctl exclusions pull --site-id 000000 --out snapshots/prod
```

## Push exclusions

Create exclusions from a local JSON file. Dry-run by default:

```bash
# Preview what would be created
s1ctl exclusions push --site-id 000000

# Apply
s1ctl exclusions push --site-id 000000 --yes
```

| Flag | Description |
|------|-------------|
| `--file` | Input file (default `samples/exclusions.json`) |
| `--site-id` | Target site IDs (repeatable) |
| `--yes` | Apply changes (default: dry-run) |

Failed entries log a warning and continue -- the command does not abort on the
first error.

## Pull, diff, push workflow

The core loop for managing exclusions as code:

```mermaid
flowchart TD
    A[s1ctl exclusions pull] --> B[samples/exclusions.json]
    B --> C{git diff}
    C -->|Changes look right| D[git commit]
    C -->|Edit needed| E[Edit exclusions.json]
    E --> C
    D --> F[s1ctl exclusions push --yes]
    F --> G[Console updated]
```

### Step by step

1. Pull the current exclusions from the console:

   ```bash
   s1ctl exclusions pull --site-id 000000
   ```

2. Review the diff against the last committed version:

   ```bash
   git diff samples/exclusions.json
   ```

3. Edit the file if changes are needed (add, remove, or modify entries).

4. Commit the snapshot to git:

   ```bash
   git add samples/exclusions.json
   git commit -m "exclusions: update path exclusions for site 000000"
   ```

5. Push the updated exclusions back to the console:

   ```bash
   s1ctl exclusions push --site-id 000000 --yes
   ```

### Why this matters

- **Audit trail.** Every exclusion change is a git commit with author, date, and
  diff.
- **Review gate.** The dry-run default prevents accidental mutations.
- **Rollback.** Revert a commit and push again to undo a change.

## Workflows

### Audit exclusions across sites

Pull exclusions from multiple sites and compare:

```bash
s1ctl exclusions pull --site-id 111111 --out samples/site-a
s1ctl exclusions pull --site-id 222222 --out samples/site-b
diff samples/site-a/exclusions.json samples/site-b/exclusions.json
```

### Export all exclusions as JSON

```bash
s1ctl exclusions list --all --json > all-exclusions.json
```

### Count exclusions by type

```bash
s1ctl exclusions list --all --json \
  | jq 'group_by(.type) | map({type: .[0].type, count: length})'
```

### Find path exclusions on a specific OS

```bash
s1ctl exclusions list --type path --os-type windows --all --json \
  | jq '.[] | {id, value, mode}'
```

### Copy exclusions between sites

Pull from one site and push to another:

```bash
s1ctl exclusions pull --site-id 111111
s1ctl exclusions push --site-id 222222          # dry-run first
s1ctl exclusions push --site-id 222222 --yes    # apply
```
