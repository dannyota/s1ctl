# Policies

Manage endpoint protection policies: list, compare across sites, detect drift,
and sync as code with pull/push.

> Prerequisites: `s1ctl` installed and configured (`S1_CONSOLE_URL`, `S1_TOKEN`).

## List policies

Fetch and compare policies across all sites (or a filtered subset). Each site
has one policy; this command presents them side by side.

```bash
s1ctl policies list
s1ctl policies list --site-id 000000
s1ctl policies list --account-id 000000 --json
```

| Flag | Description |
|------|-------------|
| `--site-id` | Filter by site ID (repeatable) |
| `--account-id` | Filter by account ID (repeatable) |

## Get a policy

Retrieve the policy for a specific scope: site, account, or group.

```bash
s1ctl policies get --site-id 000000
s1ctl policies get --account-id 000000
s1ctl policies get --group-id 000000
s1ctl policies get --site-id 000000 --json
```

| Flag | Description |
|------|-------------|
| `--site-id` | Site ID |
| `--account-id` | Account ID |
| `--group-id` | Group ID |

## Diff policies

Compare policies across sites and highlight fields that differ. Useful for
spotting inconsistencies -- for example, one site in detect mode while others
are in protect mode.

```bash
s1ctl policies diff
s1ctl policies diff --site-id 000000 --site-id 000001
s1ctl policies diff --account-id 000000 --json
```

| Flag | Description |
|------|-------------|
| `--site-id` | Filter by site ID (repeatable) |
| `--account-id` | Filter by account ID (repeatable) |

## Revert a policy

Reset a policy to the values inherited from its parent scope. Site policies
revert to the account policy, group policies revert to the site policy, and
account policies revert to global defaults.

Dry-run by default.

```bash
s1ctl policies revert --scope site --id 000000          # dry-run
s1ctl policies revert --scope site --id 000000 --yes    # apply
s1ctl policies revert --scope group --id 000000 --yes
```

| Flag | Description |
|------|-------------|
| `--scope` | Policy scope: `site`, `account`, or `group` (default `site`) |
| `--id` | Target scope ID (site, account, or group ID) |
| `--yes` | Apply the revert (default: dry-run) |

## Config-as-code

Pull policies to local YAML files, review in git, and push back. See
[Config-as-code](config-as-code.md) for the general pattern.

### Pull

```bash
s1ctl policies pull
s1ctl policies pull --site-id 000000
s1ctl policies pull --scope account --account-id 000000
s1ctl policies pull --scope group --site-id 000000
s1ctl policies pull --out snapshots/policies
```

| Flag | Description |
|------|-------------|
| `--site-id` | Filter by site ID (repeatable) |
| `--account-id` | Filter by account ID (repeatable) |
| `--scope` | Policy scope: `site`, `account`, or `group` (default `site`) |
| `--out` | Output directory (default `policies`) |

### Push

Each YAML file must contain a scope field and matching scope ID. The command
fetches the current policy, diffs against the desired state, and applies
changes. Dry-run by default.

```bash
s1ctl policies push                    # dry-run
s1ctl policies push --yes              # apply
s1ctl policies push --dir my-policies --yes
```

| Flag | Description |
|------|-------------|
| `--dir` | Input directory (default `policies`) |
| `--yes` | Apply changes (default: dry-run) |

### Pull, diff, push workflow

```bash
# 1. Pull current policies
s1ctl policies pull --site-id 000000

# 2. Review
git diff policies/

# 3. Edit policy YAML if needed (e.g. change mitigationMode)

# 4. Commit
git add policies/ && git commit -m "policies: enable protect mode on site 000000"

# 5. Push
s1ctl policies push --yes
```

## Workflows

### Detect policy drift across sites

```bash
s1ctl policies diff --account-id 000000
```

Fields that differ between sites are highlighted. Use `--json` to pipe to
automation.

### Enforce a standard policy

Pull the baseline site's policy and push it to others:

```bash
s1ctl policies pull --site-id 111111 --out baseline
cp -r baseline/* policies/
# Edit scope IDs in each file to target other sites
s1ctl policies push --yes
```

### Audit policy history

Commit each pull to git. The git log becomes an audit trail of policy changes
over time.

```bash
s1ctl policies pull --site-id 000000
git add policies/ && git commit -m "policies: weekly snapshot"
```

## See also

- [Config-as-code](config-as-code.md) -- the pull/review/push loop
- [Sites and groups](sites-groups.md) -- site and group management
- [`policies` command reference](../commands/policies.md)
