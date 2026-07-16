# Applications

Manage application inventory, risk assessment, and application control rules.
The `applications` command group (alias: `apps`) covers installed application
listing, CVE and risk queries, and the application control subsystem (rules,
settings, labels, and management settings).

## Prerequisites

- s1ctl [installed](install.md) and [configured](configure.md)
- `S1_CONSOLE_URL` and `S1_TOKEN` set

## Application inventory

List installed applications across endpoints:

```bash
s1ctl applications list --site-id 000000
s1ctl applications list --site-id 000000 --json
```

## CVEs and risks

Query known CVEs across your application inventory:

```bash
s1ctl applications cves --site-id 000000
```

List application risks (CVE vulnerabilities per endpoint):

```bash
s1ctl applications risks --site-id 000000
```

## Application control rules

Application control rules define allow/block policies for applications.
Rules support the full config-as-code loop: pull, review in `git diff`, push.

### List and inspect

```bash
s1ctl applications rules list --site-id 000000
s1ctl applications rules get 000000
```

### Create and update

```bash
s1ctl applications rules create --name "My Rule" --behavior allow --scope-type site --scope-id 000000 --yes
s1ctl applications rules update 000000 --behavior block --yes
```

### Delete

```bash
s1ctl applications rules delete 000000 --yes
```

### Pull and push (config-as-code)

Pull application control rules to local YAML files, review changes in git,
then push back:

```bash
s1ctl applications rules pull --scope-id 000000
# review and edit files in appcontrol-rules/
s1ctl applications rules push --scope-id 000000          # dry-run
s1ctl applications rules push --scope-id 000000 --yes    # apply
```

Rules are matched by name. New rules create, changed rules update. See
[Config-as-code](config-as-code.md) for the full reconcile model.

## Application control settings

Read and update the application control settings (enable/disable, mode):

```bash
s1ctl applications settings get
s1ctl applications settings update --scope-type site --scope-id 000000 --enable true --yes
```

## Management settings

Manage application management settings (scan schedule, extensive scan):

```bash
s1ctl applications mgmt-settings get --site-id 000000
s1ctl applications mgmt-settings update --site-id 000000 --extensive-scan true --yes
```

## Labels

List application control labels:

```bash
s1ctl applications labels list
```

## Workflows

### Audit application control posture

```bash
s1ctl applications settings get --site-id 000000 --json > app-settings.json
s1ctl applications rules list --site-id 000000 --json > app-rules.json
```

### Sync rules across sites

Pull from one site and push to another:

```bash
s1ctl applications rules pull --scope-id 111111
s1ctl applications rules push --scope-id 222222          # dry-run
s1ctl applications rules push --scope-id 222222 --yes    # apply
```

## See also

- [Config-as-code](config-as-code.md) — reconcile model and drift detection
- [Settings](settings.md) — platform settings management
