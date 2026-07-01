# s1ctl

Operate **SentinelOne Singularity Platform** as code — one Go CLI and one
importable Go SDK covering the REST Management API, Singularity Data Lake, and
GraphQL surfaces.

The core loop is **pull live state, review in `git diff`, push back**, with git
history as the source of truth.

## Install

```bash
go install danny.vn/s1/cmd/s1ctl@latest
```

Or download a pre-built binary from the
[releases page](https://github.com/dannyota/s1ctl/releases).

## Configure

```bash
s1ctl config          # interactive wizard
s1ctl doctor          # verify auth + API reach
```

Set `S1_CONSOLE_URL` and `S1_TOKEN` as environment variables, or let the wizard
write `~/.s1ctl/config.yaml`.

## API surfaces

| Surface | Protocol | Package | Methods |
|---------|----------|---------|---------|
| **REST MGMT** (v2.1) | REST | `danny.vn/s1/mgmt` | 72 |
| **SDL** | REST + GraphQL | `danny.vn/s1/sdl` | 13 |
| **GraphQL** | GraphQL | `danny.vn/s1/graphql` | 15 |

## CLI usage

```bash
s1ctl agents list --limit 10
s1ctl threats list --status active
s1ctl alerts list --limit 10
s1ctl vulnerabilities list --limit 10
s1ctl datalake powerquery --query "endpoint.name contains 'srv'"
```

All mutations are dry-run by default — pass `--yes` to apply:

```bash
s1ctl agents isolate --id 000000 --yes
s1ctl threats mitigate --id 000000 --action kill --yes
```

Config-as-code:

```bash
s1ctl exclusions pull --site-id 000000
git diff samples/exclusions.json
s1ctl exclusions push --site-id 000000 --yes
```

Every read command supports `--json` for machine-readable output.

## Go SDK

Each package is independently importable:

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

## Documentation

Full docs at [s1ctl.danny.vn](https://s1ctl.danny.vn).

## License

MIT
