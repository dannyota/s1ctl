# Roadmap

Wave-based delivery. Updated each time a wave completes.

## Wave 1 — SDK

Build all three SDK packages as independently importable Go libraries.

| Package | Source | Generator |
|---------|--------|-----------|
| `mgmt/` | REST MGMT swagger v2.1 | oapi-codegen |
| `graphql/` | GraphQL schemas (6 domains) | genqlient |
| `sdl/` | SDL API docs | Hand-written |
| `auth/` | — | Hand-written |
| `config/` | — | Hand-written |

## Wave 2 — Foundation

CLI skeleton and core infrastructure.

- Module init, entry point, root command
- Interactive config wizard
- Doctor command (verify connectivity to all three APIs)
- Table and JSON output formatting
- Version and commands catalog

## Wave 3 — Read commands

Read-only commands across all surfaces — every CLI group gets `list`, `get`,
and/or `query`.

| Domain | Surfaces |
|--------|----------|
| Endpoint security | agents, threats, alerts, sites, groups, accounts, policies, exclusions |
| Detection & response | rules, visibility, remoteops |
| Cloud & vuln mgmt | vulnerabilities, misconfigurations, cloud policies, cloud onboarding |
| Data lake | query, powerquery |
| Platform admin | users, settings, updates, tags, activities |
| App & device control | applications, devices, firewall, network |
| Other | automation, marketplace, inventory, identity |

**Release v0.1.0** after Wave 3.

## Wave 4 — Mutations and config-as-code

Write operations and the `pull`/`push` loop.

- Agent actions (isolate, connect, scan, decommission)
- Threat mitigation
- Config-as-code: exclusions, custom rules, policies, firewall rules, cloud policies, settings
- Dry-run guard on all mutations

## Wave 5 — Extended operations

- Remote ops (scripts, forensics)
- Data lake ingest
- Application and device management
- User and tag management
- Inventory actions

## Wave 6 — Polish

- Shell completion (bash, zsh, fish, powershell)
- Progress indicators
- Rate limiting
- Release automation (goreleaser — Win/Mac/Linux × amd64/arm64)
