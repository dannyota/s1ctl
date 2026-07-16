# Custom detection rules (STAR)

Manage STAR custom detection rules: list, inspect, assess health, review
detection trends, and manage rules as code with pull/push.

> Prerequisites: `s1ctl` installed and configured (`S1_CONSOLE_URL`, `S1_TOKEN`).

## List rules

```bash
s1ctl rules list
s1ctl rules list --site-id 000000
s1ctl rules list --status Active --severity High,Critical
s1ctl rules list --query-type events --scope site
s1ctl rules list --name "Lateral Movement" --all --json
```

| Flag | Description |
|------|-------------|
| `--site-id` | Filter by site ID (repeatable) |
| `--status` | Filter by status: `Draft`, `Active`, `Disabled`, etc. (repeatable) |
| `--severity` | Filter by severity: `Info`, `Low`, `Medium`, `High`, `Critical` (repeatable) |
| `--query-type` | Filter by query type: `events`, `correlation`, `scheduled` (repeatable) |
| `--scope` | Filter by scope: `global`, `account`, `site`, `group` (repeatable) |
| `--name` | Substring match on rule name |
| `--query` | Free text search on S1QL |
| `--sort-by` | Sort field (e.g. `name`, `severity`, `createdAt`) |
| `--sort-order` | Sort direction (`asc`, `desc`) |
| `--limit` | Max results per page (default 50) |
| `--all` | Fetch all pages |
| `--cursor` | Pagination cursor |

## Get rule details

```bash
s1ctl rules get 000000
s1ctl rules get 000000 --json
```

## Rule health

Classify all rules by operational state: firing (active with alerts), silent
(active with zero alerts), disabled, or erroring (reached alert limit).

```bash
s1ctl rules health
s1ctl rules health --site-id 000000
```

## Detection trends

Sort rules by generated alert count (descending) to identify alert fatigue
candidates.

```bash
s1ctl rules trends
s1ctl rules trends --site-id 000000 --top 10
```

| Flag | Description |
|------|-------------|
| `--site-id` | Filter by site ID (repeatable) |
| `--top` | Show only top N rules (default: all) |

## View detections for a rule

List cloud detection alerts (STAR alerts) filtered by rule name.

```bash
s1ctl rules detections "Lateral Movement"
s1ctl rules detections "Lateral Movement" --since 2025-01-01T00:00:00Z
s1ctl rules detections "Lateral Movement" --severity High --status unresolved
s1ctl rules detections "Lateral Movement" --group-by agent --all --json
```

| Flag | Description |
|------|-------------|
| `--site-id` | Filter by site ID (repeatable) |
| `--severity` | Filter by severity (repeatable) |
| `--status` | Filter by incident status (repeatable) |
| `--since` | Show detections after this time (RFC 3339) |
| `--group-by` | Group results (e.g. `agent`) |
| `--sort-by` | Sort field (default: `id`) |
| `--sort-order` | Sort direction (default: `desc`) |
| `--limit` | Max results per page (default 50) |
| `--all` | Fetch all pages |
| `--cursor` | Pagination cursor |

## Enable and disable rules

All mutations are dry-run by default. Pass `--yes` to apply.

```bash
s1ctl rules enable 000000                  # dry-run
s1ctl rules enable 000000 --yes            # apply
s1ctl rules disable 000000 000001 --yes    # disable multiple
```

## Validate rules locally

Check rule YAML files for missing fields, invalid enums, and empty queries.
No API calls are made.

```bash
s1ctl rules validate
s1ctl rules validate --dir my-rules
```

## Diff local vs live

Compare local rule YAML files against the live console to preview what would
change before pushing.

```bash
s1ctl rules diff
s1ctl rules diff --dir my-rules
```

## Config-as-code

Pull rules to local YAML files, review in git, and push back. See
[Config-as-code](config-as-code.md) for the general pattern.

### Pull

```bash
s1ctl rules pull --site-id 000000
s1ctl rules pull --site-id 000000 --out snapshots/rules
```

| Flag | Description |
|------|-------------|
| `--site-id` | Filter by site ID (repeatable) |
| `--out` | Output directory (default `rules`) |

### Push

Rules are matched by name: existing rules are updated, new rules are created,
unchanged rules are skipped. Dry-run by default.

```bash
s1ctl rules push                    # dry-run
s1ctl rules push --yes              # apply
s1ctl rules push --dir my-rules     # custom directory
```

| Flag | Description |
|------|-------------|
| `--dir` | Input directory (default `rules`) |
| `--yes` | Apply changes (default: dry-run) |

### Pull, diff, push workflow

```bash
# 1. Pull current rules
s1ctl rules pull --site-id 000000

# 2. Review
git diff rules/

# 3. Validate locally
s1ctl rules validate

# 4. Preview changes against live
s1ctl rules diff

# 5. Commit snapshot
git add rules/ && git commit -m "rules: snapshot site 000000"

# 6. Push
s1ctl rules push --yes
```

## Workflows

### Find noisy rules causing alert fatigue

```bash
s1ctl rules trends --top 5 --site-id 000000
```

Inspect the top offender and its detections:

```bash
s1ctl rules detections "Noisy Rule Name" --all --json \
  | jq 'group_by(.agentId) | map({agent: .[0].agentId, count: length})'
```

### Audit silent rules

Find active rules that have never fired:

```bash
s1ctl rules health --site-id 000000 --json \
  | jq '.[] | select(.state == "silent")'
```

### Copy rules between sites

```bash
s1ctl rules pull --site-id 111111
s1ctl rules push --yes    # pushes to the configured scope
```

## See also

- [Config-as-code](config-as-code.md) -- the pull/review/push loop
- [Alerts](alerts.md) -- STAR rules generate cloud detection alerts
- [`rules` command reference](../commands/rules.md)
