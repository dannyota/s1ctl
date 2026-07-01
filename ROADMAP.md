# Roadmap

Wave-based delivery. Updated each time a wave completes.

## Wave 1 — SDK (complete)

Build all three SDK packages as independently importable Go libraries.

| Package | Surface | Methods |
|---------|---------|---------|
| `mgmt/` | REST MGMT API v2.1 | 72 (CRUD + actions across 15 resources) |
| `sdl/` | SDL REST + GraphQL | 13 (query, ingest, files, powerquery, GraphQL ops) |
| `graphql/` | GraphQL APIs (4 domains) | 15 (alerts, misconfigs, vulns, cloud policies) |
| `auth/` | Token management | 2 credential types |
| `config/` | Instance config | Load, validate, resolve |

## Wave 2 — Foundation (complete)

CLI skeleton and core infrastructure.

- Module init, entry point, root command
- Interactive config wizard (SDL URL support)
- Doctor command (verify connectivity to all three APIs)
- Table and JSON output formatting
- Version and commands catalog
- Shell completion (bash, zsh, fish, powershell)

## Wave 3 — Read commands (complete)

Read-only CLI commands across all surfaces.

| Domain | Commands |
|--------|----------|
| Endpoint security | agents, threats, alerts, sites, groups, accounts, policies, exclusions |
| Operations | users, tags, remoteops, applications, device-control, firewall, updates |
| Data lake | powerquery (GraphQL + REST protocols) |
| Platform | activities |

## Wave 4 — Mutation commands (complete: SDK; partial: CLI)

Agent and threat actions, exclusion sync, and core mutations.

| Surface | CLI wired | SDK ready |
|---------|-----------|-----------|
| Agent actions | isolate, connect, scan, decommission | + update-software, move-to-site, fetch-logs, restart, enable, disable, reset-config, approve/reject-uninstall, mark-up-to-date, set-external-id, randomize-uuid, firewall-logging |
| Threat actions | mitigate, verdict, status | + add-to-blacklist, fetch-file |
| Exclusions | pull, push, create, delete | + update |
| Sites | -- | create, update, delete |
| Groups | -- | create, update, delete |
| Tags | -- | create, update, delete |
| Policies | -- | update (site, account, group) |
| Users | -- | delete |

## Wave 5 — Cloud and xSPM CLI

Wire the GraphQL domain SDKs into CLI commands.

- `s1ctl misconfigurations list|get|status|verdict`
- `s1ctl vulnerabilities list|get|status|verdict`
- `s1ctl cloud-policies list|get`
- Expanded alerts: `s1ctl alerts get|status|verdict`

## Wave 6 — Data Lake CLI

Wire remaining SDL operations into CLI commands.

- `s1ctl datalake query` — raw log search
- `s1ctl datalake facet` — facet aggregation
- `s1ctl datalake timeseries` — time-series aggregation
- `s1ctl datalake ingest` — addEvents / uploadLogs
- `s1ctl datalake files` — getFile / putFile / listFiles

## Wave 7 — Config-as-code

Pull/push loop for all mutable resources.

- Sites, groups, tags pull/push
- Policies pull/push (site, account, group)
- Cloud policies pull/push
- Firewall rules pull/push
- Device control rules pull/push
- Dry-run diff before every push

## Wave 8 — Extended mutations CLI

Wire remaining SDK mutations into CLI commands.

- Agent actions: all 20 actions accessible as subcommands
- Threat actions: blacklist, fetch-file
- Site/group/tag CRUD commands
- Policy update commands
- User management commands

## Wave 9 — Polish and release

- Progress indicators (bubbletea spinners for long operations)
- Rate limiting (x/time/rate)
- Release automation (goreleaser — Win/Mac/Linux x amd64/arm64)
- Docs site updates for all new commands
- **Release v1.0.0**
