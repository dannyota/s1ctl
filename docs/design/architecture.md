# Architecture

s1ctl is an **SDK-first** project: the Go SDK packages are the product, the CLI
is a thin consumer. Both ship together.

## SDK packages

```text
danny.vn/s1/mgmt      REST MGMT client (generated from swagger v2.1)
danny.vn/s1/sdl       SDL Data Lake client (hand-written, 12 endpoints)
danny.vn/s1/graphql   GraphQL client (generated via genqlient)
danny.vn/s1/auth      Token management (shared by all three)
danny.vn/s1/config    Instance config resolution
```

Each package is independently importable:

```go
import "danny.vn/s1/mgmt"

client := mgmt.NewClient(mgmt.WithToken(token), mgmt.WithConsole(url))
agents, err := client.Agents.List(ctx, &mgmt.AgentListParams{})
```

The SDK packages are **pure** — HTTP calls and typed structs, no disk I/O. All
on-disk layout (pull/push file trees, config files) lives in `internal/`.

## Three protocols, one CLI

SentinelOne exposes three API protocols. The CLI unifies them under one command
tree — users never think about which protocol is underneath.

```mermaid
flowchart LR
  CLI["s1ctl CLI"]

  subgraph SDK["Go SDK packages"]
    MGMT["mgmt/<br/>REST MGMT v2.1<br/>680 endpoints"]
    SDL["sdl/<br/>Singularity Data Lake<br/>12 endpoints"]
    GQL["graphql/<br/>GraphQL<br/>6 domains"]
  end

  subgraph S1["SentinelOne Console"]
    REST_API["REST API"]
    SDL_API["SDL API"]
    GQL_API["GraphQL API"]
  end

  CLI --> MGMT
  CLI --> SDL
  CLI --> GQL
  MGMT --> REST_API
  SDL --> SDL_API
  GQL --> GQL_API
```

| Protocol | Package | Auth header | Scope |
|----------|---------|-------------|-------|
| REST MGMT v2.1 | `mgmt/` | `ApiToken <token>` | Agents, threats, sites, groups, exclusions, policies, remote ops |
| SDL | `sdl/` | `Bearer <token>` | PowerQuery, log ingest/query, file ops |
| GraphQL | `graphql/` | `Bearer <token>` | UAM alerts, xSPM vulns/misconfigs, cloud security |

## Protocol selection

When a surface is available via multiple protocols (e.g. alerts via both REST
and GraphQL), s1ctl defaults to the protocol with the best performance and
richest filtering. Users can override with `--protocol rest|graphql|sdl`.

Each command's `--help` documents which protocol is the default and why.

## Two planes

| Plane | Loop | Source of truth | Surfaces |
|-------|------|-----------------|----------|
| **Control** | pull &rarr; `git diff` &rarr; push | Git | Exclusions, custom rules, policies, cloud policies |
| **Operational** | query &rarr; review &rarr; act | Live instance | Agents, threats, alerts, vulns, misconfigs, data lake, remote ops, inventory |

Control plane is narrow — most of SentinelOne is operational. Config-as-code
surfaces use the same reconcile model as secopsctl: identity by server ID,
canonical diff, dry-run by default, `--yes` to apply.

## Codegen strategy

The REST MGMT API has 680 endpoints across 119 tags — hand-writing is not
viable. The GraphQL API spans 6 schema domains with ~236 operations.

| Surface | Generator | Source | Output |
|---------|-----------|--------|--------|
| REST MGMT | `oapi-codegen` | `references/swagger_2_1.json` | `mgmt/` |
| GraphQL | `genqlient` | `references/graphql/*.graphql` | `graphql/` |
| SDL | Hand-written | `references/sdl-api/*.md` | `sdl/` |

Generated code is never edited directly. Hand-written wrappers live alongside
in separate files, providing ergonomic service-oriented access.

## CLI structure

Commands follow SentinelOne's official terminology. Plural nouns at the top
level, verbs nested underneath.

```text
s1ctl agents list|get|isolate|scan|...
s1ctl threats list|get|mitigate|...
s1ctl alerts query|get|...
s1ctl exclusions list|get|pull|push|...
s1ctl datalake query|powerquery|...
s1ctl config
s1ctl doctor
```

### Cross-cutting flags

| Flag | Scope | Default |
|------|-------|---------|
| `--json` | All read commands | false (table output) |
| `--yes` | All mutations | false (dry-run) |
| `--site-id` | Most commands | from config |
| `--limit` | List commands | API default |
| `--config` | All | `~/.s1ctl/config.yaml` |

Full command naming conventions in [CLI naming](cli-naming.md). Domain-to-API
mapping in [Surfaces](surfaces.md). Implementation status in
[Catalog](catalog.md).
