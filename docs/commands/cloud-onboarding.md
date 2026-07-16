# cloud-onboarding

Manage CNAPP cloud account onboarding

## cloud-onboarding delete

Delete (offboard) cloud entities

```text
s1ctl cloud-onboarding delete <account-id> [account-id...] [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--yes` | bool | false | apply (default: dry-run) |

## cloud-onboarding get

Get onboarded cloud entity details

```text
s1ctl cloud-onboarding get <account-id>
```

## cloud-onboarding list

List onboarded cloud entities

```text
s1ctl cloud-onboarding list [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--after` | string | - | pagination cursor |
| `--all` | bool | false | fetch all pages |
| `--limit` | int | 0 | max results per page (default 50) |
| `--provider` | stringSlice | - | filter by cloud provider (AWS, GCP, AZURE, OCI, ALIBABA) |
| `--status` | stringSlice | - | filter by operational status |

## cloud-onboarding onboard

Onboard a cloud entity from a JSON file

```text
s1ctl cloud-onboarding onboard --from-file <request.json> [flags]
```

Onboard a new cloud entity (AWS account, GCP project, Azure
subscription, etc.) using a CnappCloudOnboardingRequest JSON file. The file
must contain the full onboarding payload including cloudProvider, products,
and credentials.

Dry-run by default; pass --yes to apply.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--from-file` | string | - | path to onboarding request JSON file (required) |
| `--yes` | bool | false | apply (default: dry-run) |
