# Usage examples

Real-world workflows using s1ctl.

## Triage alerts

List NEW alerts sorted by detection time:

```bash
s1ctl alerts list --status NEW --sort-by detectedAt --sort-order DESC
```

Count CRITICAL alerts:

```bash
s1ctl alerts count --severity CRITICAL
```

Resolve alerts in bulk:

```bash
s1ctl alerts resolve <id1> <id2> <id3> --yes
```

## Investigate threats

List unresolved threats with agent name and creation date:

```bash
s1ctl threats list --status unresolved
```

Filter by mitigation status:

```bash
s1ctl threats list --mitigation-status not_mitigated --all
```

Resolve benign threats in bulk:

```bash
s1ctl threats resolve <id1> <id2> --yes
```

## Agent management

Find outdated agents:

```bash
s1ctl agents outdated --all
```

Show version distribution:

```bash
s1ctl agents versions
```

Trigger upgrade on outdated agents in a site:

```bash
s1ctl agents upgrade --site-id 000000 --yes
```

## Environment health

Get a one-shot dashboard:

```bash
s1ctl status
```

JSON output for scripting:

```bash
s1ctl status --json | jq '.unresolved'
```

## Config-as-code

Pull exclusions scoped to a site:

```bash
s1ctl exclusions pull --site-id 000000
```

Create a quick exclusion:

```bash
s1ctl exclusions create --type path --value "/opt/app/logs/" \
  --os-type linux --site-id 000000 --yes
```

Pull and compare policies across sites:

```bash
s1ctl policies diff
```

Pull policies to YAML, edit, then push:

```bash
s1ctl policies pull --out policies/
# edit policies/production.yaml
s1ctl policies push --dir policies/ --yes
```

Review rule effectiveness:

```bash
s1ctl rules diff
```

## Data lake queries

Run a PowerQuery with wider columns:

```bash
s1ctl datalake powerquery --query "src.process.name = 'sshd'" \
  --start 7d --col-width 120
```

List saved queries from the console:

```bash
s1ctl datalake saved-queries
```

## Group management

Create a group:

```bash
s1ctl groups create --site-id 000000 --name "Staging" --yes
```

Delete a group:

```bash
s1ctl groups delete <group-id> --yes
```

## Output formats

All read commands support `--output` (table, json, csv):

```bash
s1ctl agents list --output csv > agents.csv
s1ctl threats list --json | jq '.[] | .threatName'
```
