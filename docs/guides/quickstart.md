# Quickstart

Common workflows with s1ctl. Assumes you have already
[installed](guides/install.md) and [configured](guides/configure.md) the CLI.

## Read operations

List agents, filter and paginate:

```bash
s1ctl agents list --limit 10 --query "win"
s1ctl agents list --json | jq '.[].computerName'
s1ctl agents count
```

List threats and alerts:

```bash
s1ctl threats list --limit 5
s1ctl alerts list --limit 10
```

Browse sites, groups, accounts:

```bash
s1ctl sites list
s1ctl groups list
s1ctl accounts list
```

## Mutations

All mutations are dry-run by default. Pass `--yes` to apply.

Isolate an agent:

```bash
s1ctl agents isolate --id 000000          # dry-run
s1ctl agents isolate --id 000000 --yes    # apply
```

Mitigate a threat:

```bash
s1ctl threats mitigate --id 000000 --action kill --yes
```

## Config-as-code

Pull exclusions to a local file, review, then push back:

```bash
s1ctl exclusions pull --site-id 000000
git diff samples/exclusions.json
s1ctl exclusions push --site-id 000000 --yes
```

## Data lake

Run a powerquery against the Singularity Data Lake:

```bash
s1ctl datalake powerquery --query "endpoint.name contains 'srv'"
```

## Shell completion

```bash
source <(s1ctl completion bash)
s1ctl completion zsh > "${fpath[1]}/_s1ctl"
s1ctl completion fish | source
```

## JSON output

Every read command supports `--json` for machine-readable output:

```bash
s1ctl sites list --json | jq '.[] | {id, name, state}'
s1ctl agents list --json | jq 'length'
```
