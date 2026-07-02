# Config-as-code

Manage SentinelOne configuration through local files and git.

## Core loop

1. **Pull** live state to local files
2. **Review** changes in `git diff`
3. **Push** desired state back (dry-run by default)

## Surfaces

| Surface | Pull | Push | Format |
|---------|------|------|--------|
| Exclusions | `exclusions pull` | `exclusions push` | JSON |
| Policies | `policies pull` | `policies push` | YAML (per site) |
| Rules | `rules pull` | `rules push` | YAML (per rule) |
| Device control | `devicecontrol pull` | `devicecontrol push` | JSON |
| Firewall | `firewall pull` | `firewall push` | JSON |
| Sites | `sites pull` | `sites push` | JSON |
| Groups | `groups pull` | `groups push` | JSON |
| Tags | `tags pull` | `tags push` | JSON |
| Cloud policies | `cloud-policies pull` | `cloud-policies push` | JSON |

## Exclusions

Pull all exclusions (or filter by site):

```bash
s1ctl exclusions pull --out samples/
s1ctl exclusions pull --site-id 000000 --out samples/
```

This creates `samples/exclusions.json`. Edit the file, then push:

```bash
s1ctl exclusions push --file samples/exclusions.json --site-id 000000
# dry-run output shows what would be created
s1ctl exclusions push --file samples/exclusions.json --site-id 000000 --yes
```

## Policies

Pull fetches one YAML file per site:

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

Same pull/push pattern:

```bash
s1ctl devicecontrol pull --site-id 000000 --out samples/
s1ctl devicecontrol push --file samples/device-control.json --site-id 000000 --yes

s1ctl firewall pull --site-id 000000 --out samples/
s1ctl firewall push --file samples/firewall-rules.json --site-id 000000 --yes
```

## Sites, groups, and tags

Pull the hierarchy objects to JSON, review, and push new ones back:

```bash
s1ctl sites pull --out samples/
s1ctl sites push --file samples/sites.json --yes

s1ctl groups pull --site-id 000000 --out samples/
s1ctl groups push --file samples/groups.json --yes

s1ctl tags pull --site-id 000000 --out samples/
s1ctl tags push --file samples/tags.json --yes
```

For these surfaces, `push` creates the objects listed in the file. See the
[Sites and groups](guides/sites-groups.md) guide for per-command flags.

## Cloud policies

Pull cloud security policies, flip their enabled/disabled state in the file,
then push the reconciled status back:

```bash
s1ctl cloud-policies pull --out samples/
# edit samples/cloud-policies.json: set "status" to "enabled" or "disabled" per policy
s1ctl cloud-policies push --file samples/cloud-policies.json --yes
```

`cloud-policies push` reconciles only the enabled/disabled status of each
policy in the file — it enables or disables policies to match, and never
creates or deletes them.

## Settings

Settings follow a read-edit-apply round-trip rather than a bulk pull. Read a
category, edit the JSON, then apply it:

```bash
s1ctl settings get syslog > syslog.json
# edit syslog.json
s1ctl settings update syslog --from-file syslog.json --yes
```

Updatable categories: `notifications`, `sso`, `smtp`, `syslog`. Scope with
`--site-id` or `--account-id`. Secrets are never echoed: `settings get`
redacts sensitive fields, and `settings update` reports status only without
printing the payload back. Because the pulled file has secrets redacted,
re-enter any secret fields (passwords, tokens, certificate contents) before
pushing — otherwise the update writes them back empty.

## Tips

- All push commands are **dry-run by default**. Pass `--yes` to apply.
- Commit pulled files to git for audit trail and diff review.
- Use `--site-id` to scope pulls to specific sites.
- The push diff output shows exactly what will change before applying.
