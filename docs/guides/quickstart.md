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

Pull exclusions to a local directory, review, then push back:

```bash
s1ctl exclusions pull --site-id 000000
git diff exclusions/
s1ctl exclusions push --site-id 000000 --yes
```

## Cloud and vulnerability management

Query xSPM findings and cloud policies:

```bash
s1ctl vulnerabilities list --limit 10
s1ctl misconfigurations list --limit 10
s1ctl cloud-policies list --limit 10
```

## Data lake

Run a powerquery against the Singularity Data Lake:

```bash
s1ctl datalake powerquery --query "endpoint.name contains 'srv'"
s1ctl datalake powerquery --query "endpoint.name contains 'srv'" --protocol rest
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

## Go SDK

The SDK packages are independently importable:

```go
import "danny.vn/s1/mgmt"

client := mgmt.NewClient("https://your-console.sentinelone.net", token)
agents, _, err := client.AgentsList(ctx, nil)
```

```go
import "danny.vn/s1/graphql"

client := graphql.NewClient("https://your-console.sentinelone.net", token)
alerts, err := client.AlertsList(ctx, &graphql.ListParams{First: 10})
```

```go
import "danny.vn/s1/sdl"

client := sdl.NewClient("https://your-console.sentinelone.net", token)
resp, err := client.PowerQueryGraphQL(ctx, &sdl.PowerQueryRequest{
    Query:     "endpoint.name contains 'srv'",
    StartTime: "24h",
})
```
