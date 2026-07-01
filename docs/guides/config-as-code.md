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

Push syncs rules by name — existing rules are updated, new ones created:

```bash
s1ctl rules push --dir rules/ --yes
```

Review which rules have fired vs dormant:

```bash
s1ctl rules diff
```

## Device control and firewall

Same pull/push pattern:

```bash
s1ctl devicecontrol pull --site-id 000000 --out samples/
s1ctl devicecontrol push --file samples/device-control.json --site-id 000000 --yes

s1ctl firewall pull --site-id 000000 --out samples/
s1ctl firewall push --file samples/firewall-rules.json --site-id 000000 --yes
```

## Tips

- All push commands are **dry-run by default**. Pass `--yes` to apply.
- Commit pulled files to git for audit trail and diff review.
- Use `--site-id` to scope pulls to specific sites.
- The push diff output shows exactly what will change before applying.
