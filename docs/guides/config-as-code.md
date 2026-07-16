# Config-as-code

Manage SentinelOne configuration through local files and git.

## Core loop

1. **Pull** live state to local files
2. **Review** changes in `git diff`
3. **Push** desired state back (dry-run by default)

## Reconcile model

Most sync surfaces share one on-disk model and one engine (see
[Reconcile engine](../design/reconcile.md)). Each object is a single YAML file
in a per-surface directory named after the surface (`sites/`, `firewall/`,
`tags/`, ...). Files hold only the declarative definition — server-assigned
IDs, scopes, and timestamps are omitted so diffs never churn on server-managed
fields.

`pull` renders live objects to files (overwrite; never deletes). `push` matches
files to live objects by a stable identity, then **creates** objects with no
live match and **updates** objects whose body differs — dry-run by default,
`--yes` to apply. Live-only objects (live, no local file) are reported, never
deleted. A push exits non-zero if any per-item apply fails.

### Immutable fields

A file carries some fields that identify or place the object but that an update
cannot change — a site's `accountId` or `siteType`, a tag's `scope`. Editing
one of these does not move the object; the push keeps planning an `update`
every run because the live object never matches. Revert the field to its pulled
value to clear the drift.

### Stale local files

`pull` never deletes local files, but it warns about **stale** ones — a `.yaml`
file with no live counterpart. A stale file (its object was deleted or renamed
away in the console) plans as a `create` on the next push and would re-create
the object. Delete stale files before pushing.

## Surfaces

| Surface | Pull | Push | Layout |
|---------|------|------|--------|
| Exclusions | `exclusions pull` | `exclusions push` | YAML dir (per exclusion) |
| Policies | `policies pull` | `policies push` | YAML (per scope) |
| Rules | `rules pull` | `rules push` | YAML dir (per rule) |
| Device control | `devicecontrol pull` | `devicecontrol push` | YAML dir (per rule) |
| Firewall | `firewall pull` | `firewall push` | YAML dir (per rule) |
| Sites | `sites pull` | `sites push` | YAML dir (per site) |
| Groups | `groups pull` | `groups push` | YAML dir (per group) |
| Tags | `tags pull` | `tags push` | YAML dir (per tag) |
| Blocklist | `blocklist pull` | `blocklist push` | YAML dir (per item) |
| Network | `network pull` | `network push` | YAML dir (per rule) |
| Locations | `locations pull` | `locations push` | YAML dir (per location) |
| Cloud policies | `cloud-policies pull` | `cloud-policies push` | YAML dir (per policy) |
| Upgrade policies | `upgrade-policies pull` | `upgrade-policies push` | YAML dir (per policy) |
| Application control rules | `applications rules pull` | `applications rules push` | YAML dir (per rule) |

Policies are the one exception: they are scope-singletons (account/site/group)
with their own pull/push/diff/revert lane, not a per-object collection.

## Exclusions

Pull all exclusions (or filter by site) into a directory of per-exclusion
files:

```bash
s1ctl exclusions pull
s1ctl exclusions pull --site-id 000000
```

This writes one YAML file per exclusion under `exclusions/`:

```yaml
type: path
value: /opt/app/cache
osType: linux
mode: suppress
```

Edit or delete files, then push. Exclusions are matched by type + OS + value:
matching files update, files with no live match are created, live-only entries
are reported. New exclusions are created at the scope named by `--site-id`
(global/tenant scope if omitted):

```bash
s1ctl exclusions push --site-id 000000            # dry-run
s1ctl exclusions push --site-id 000000 --yes      # apply
```

## Policies

Pull fetches one YAML file per scope:

```bash
s1ctl policies pull --out policies/
```

Each file contains the key policy fields:

```yaml
siteId: "000000"
siteName: Production
mitigationMode: protect
mitigationModeSuspicious: detect
antiTamperingOn: true
networkQuarantineOn: false
```

Edit and push:

```bash
s1ctl policies push --dir policies/ --yes
```

Compare policies across sites to find inconsistencies:

```bash
s1ctl policies diff
```

## Rules

Pull creates one YAML file per custom detection rule:

```bash
s1ctl rules pull --out rules/
```

Each file contains the rule definition:

```yaml
name: Suspicious SSH Login
s1ql: "EventType = 'Login' AND src.process.name = 'sshd'"
severity: Medium
status: Active
queryType: events
expirationMode: Permanent
treatAsThreat: UNDEFINED
```

Validate rule files before pushing:

```bash
s1ctl rules validate --dir rules/
```

Compare local files against live rules to see what changed:

```bash
s1ctl rules diff --dir rules/
```

Push syncs rules by name — existing rules are updated, new ones created:

```bash
s1ctl rules push --dir rules/ --yes
```

Check operational health of deployed rules:

```bash
s1ctl rules health
s1ctl rules trends --top 10
```

## Device control and firewall

Same per-object pattern — pull to a surface directory, edit files, push:

```bash
s1ctl devicecontrol pull --site-id 000000
s1ctl devicecontrol push --site-id 000000 --yes

s1ctl firewall pull --site-id 000000
s1ctl firewall push --site-id 000000 --yes
```

Rules are matched by name. New rules are created at the scope named by
`--site-id`.

## Sites, groups, and tags

Pull the hierarchy objects to per-object directories, review in git, then push:

```bash
s1ctl sites pull
s1ctl sites push --yes

s1ctl groups pull --site-id 000000
s1ctl groups push --yes

s1ctl tags pull --site-id 000000
s1ctl tags push --yes
```

Each object is one file. Sites are matched by name, groups by site ID + name,
and tags by key: matching files update, files with no live match are created.
A site file looks like:

```yaml
name: Production
accountId: "000000"
siteType: Paid
totalLicenses: 500
```

See the [Sites and groups](sites-groups.md) guide for per-command flags.

## Cloud policies

Pull cloud security policies to a directory, flip the `status` field in the
per-policy files, then push the reconciled status back:

```bash
s1ctl cloud-policies pull
# edit files in cloud-policies/: set status to enabled or disabled
s1ctl cloud-policies push --yes
```

Each file carries the policy identity and status:

```yaml
id: "000000"
name: Example CNS policy
status: enabled
```

`cloud-policies push` reconciles only the enabled/disabled status of each
policy — matched by ID, it enables or disables policies to match, and never
creates or deletes them. A local file whose ID has no live match fails per-item
since policies cannot be created through this surface.

## Upgrade policies

Pull auto-upgrade policies to local YAML files, review changes in git, then
push back. The API requires `--scope-level` and `--os-type` on every pull and
push:

```bash
s1ctl upgrade-policies pull --scope-level site --scope-id 000000 --os-type linux
# edit files in upgrade-policies/
s1ctl upgrade-policies push --scope-level site --scope-id 000000 --os-type linux          # dry-run
s1ctl upgrade-policies push --scope-level site --scope-id 000000 --os-type linux --yes    # apply
```

Policies are scope-partitioned by level (tenant/account/site/group) and OS
type (linux/macos/windows). Pull one partition at a time; use multiple
invocations for different OS types or scopes.

## Application control rules

Pull and push application control rules through the reconcile engine:

```bash
s1ctl applications rules pull --scope-id 000000
# edit files in appcontrol-rules/
s1ctl applications rules push --scope-id 000000          # dry-run
s1ctl applications rules push --scope-id 000000 --yes    # apply
```

Rules are matched by name. See the [Applications](applications.md) guide for
the full command reference.

## Drift detection

`s1ctl drift` reports the difference between committed files and live state for
every sync surface, without applying anything. For each surface that has a local
directory, it loads the files, lists live objects, and prints a per-surface
summary of creates, updates, live-only, and unchanged counts. It is read-only —
there is no apply path.

```bash
s1ctl drift
s1ctl drift --surface firewall --surface sites
```

Exit code is 0 when every checked surface is clean and 1 when any surface has
drift, so a CI job can fail the build on a non-zero exit:

```bash
# in CI: check the committed config against the console
s1ctl drift
```

Two caveats:

- **Surfaces without a local directory are skipped.** Drift checks only what is
  committed. An empty-but-present surface directory is not skipped: every live
  object reports as live-only, so the surface shows drift. Pull the surface
  first.
- **Live-only counts as drift.** An object that exists in the console but has no
  committed file is drift, the same as a stale local file that would re-create a
  deleted object on push. Reconcile either the console or the files to clear it.
- **Per-surface SKIPPED.** Some surfaces (e.g. upgrade-policies) are skipped by
  drift when they require scope parameters that drift cannot infer from the
  local directory alone. Drift reports these as `SKIPPED` rather than clean or
  dirty. Pull the surface with explicit scope flags to enable drift checking.

## Settings

Settings follow a read-edit-apply round-trip rather than a bulk pull. Read a
category, edit the JSON, then apply it:

```bash
s1ctl settings get syslog > syslog.json
# edit syslog.json
s1ctl settings update syslog --from-file syslog.json --yes
```

Updatable categories: `notifications`, `sso`, `smtp`, `syslog`, `sms`,
`recipients`, `ad`, `ad-scope-mapping`. Scope with `--site-id` or
`--account-id`.

> **Note:** `settings list` shows `active-directory` as the category name,
> but the CLI token for `settings get` and `settings update` is `ad`
> (abbreviated). Similarly, `ad-scope-mapping` is the CLI token for the
> Active Directory scope mapping category.

Secrets are never echoed: `settings get` redacts sensitive fields, and
`settings update` reports status only without printing the payload back.
Because the pulled file has secrets redacted, re-enter any secret fields
(passwords, tokens, certificate contents) before pushing — otherwise the
update writes them back empty.

## Tips

- All push commands are **dry-run by default**. Pass `--yes` to apply.
- Commit pulled directories to git for audit trail and diff review.
- Use `--site-id` to scope pulls to specific sites.
- The push diff output shows exactly what will change before applying.
- Run `s1ctl drift` in CI to catch console changes that bypass the loop.
