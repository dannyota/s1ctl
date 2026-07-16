# Cloud onboarding

Manage CNAPP cloud account onboarding: list onboarded entities, inspect their
details, onboard new cloud accounts, and offboard (delete) existing ones.

Cloud onboarding uses the GraphQL API. The `cloud-onboarding` command group
covers AWS, GCP, Azure, OCI, and Alibaba accounts and organizations.

## Prerequisites

- s1ctl [installed](install.md) and [configured](configure.md)
- `S1_CONSOLE_URL` and `S1_TOKEN` set
- CNAPP module enabled on your console

## List onboarded entities

```bash
s1ctl cloud-onboarding list
s1ctl cloud-onboarding list --json
```

The listing shows each entity's cloud provider, account ID, status, and
onboarding type.

## Get entity details

```bash
s1ctl cloud-onboarding get 000000
s1ctl cloud-onboarding get 000000 --json
```

Returns the full entity record including configuration, features, and
connected status.

## Onboard a new cloud entity

Create a JSON file with the onboarding payload — provider, account ID,
role ARN (AWS), project ID (GCP), or equivalent fields for other providers:

```bash
s1ctl cloud-onboarding onboard --from-file entity.json          # dry-run
s1ctl cloud-onboarding onboard --from-file entity.json --yes    # apply
```

The onboard command is **dry-run by default**; pass `--yes` to apply.

## Delete (offboard) an entity

```bash
s1ctl cloud-onboarding delete 000000          # dry-run
s1ctl cloud-onboarding delete 000000 --yes    # apply
```

Deleting is **dry-run by default**; pass `--yes` to apply.

## Workflows

### Audit onboarded accounts

List all onboarded entities and export to JSON for review:

```bash
s1ctl cloud-onboarding list --json > onboarded.json
```

### Onboard and verify

```bash
s1ctl cloud-onboarding onboard --from-file aws-account.json --yes
s1ctl cloud-onboarding list --json | grep '"accountId"'
```

## See also

- [Cloud policies](../commands/cloud-policies.md) — manage cloud security
  policies
- [Misconfigurations](../commands/misconfigurations.md) — xSPM findings
