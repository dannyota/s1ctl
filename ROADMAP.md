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

## Wave 4 — Mutation commands (complete)

Agent and threat actions, exclusion sync, and core mutations — all wired to
the CLI.

| Surface | CLI wired |
|---------|-----------|
| Agent actions | isolate, reconnect, scan, abort-scan, decommission, uninstall, approve/reject-uninstall, upgrade, move, move-to-site, fetch-logs, restart, shutdown, enable, disable, reset-config, mark-up-to-date, set-external-id, randomize-uuid, firewall-logging |
| Threat actions | mitigate, verdict, status, resolve, add-note, blacklist, fetch-file |
| Exclusions | pull, push, create, update, delete |
| Sites | create, update, delete |
| Groups | create, update, delete |
| Tags | get, create, update, delete |
| Policies | update, revert (site, account, group) |
| Users | delete |

## Wave 5 — Cloud and xSPM CLI (complete)

GraphQL domain SDKs wired into CLI commands: misconfigurations and
vulnerabilities (list/get/status/verdict), cloud-policies (list/get/enable/
disable/delete), and expanded alerts (get/status/verdict/history/stats).

## Wave 6 — Data Lake CLI (complete)

Remaining SDL operations wired into CLI commands: `datalake query`,
`datalake facet`, `datalake timeseries`, `datalake ingest events|logs`
(addEvents / uploadLogs), and `datalake files list|get|put`
(getFile / listFiles / putFile).

## Wave 7 — Config-as-code (complete)

Pull/push loop across mutable resources: exclusions, policies (site/account/
group), rules, firewall, device control, sites, groups, tags, and cloud
policies. Every push is dry-run by default with a diff, and applies only
with `--yes`.

## Wave 8 — Extended mutations CLI (complete)

Remaining SDK mutations wired into CLI commands: the full agent action set,
threat blacklist and fetch-file, site/group/tag CRUD, policy updates, user
deletion, settings updates (notifications/sso/smtp/syslog), and the
upgrade-policies lifecycle (create/update/delete/activate/deactivate).

## Wave 9 — Reconcile engine and drift gate (complete)

A shared reconcile engine (`internal/reconcile`) backs every sync surface: one
on-disk model, one set of create-or-update semantics, and a `drift` command
that reports per-surface divergence for CI gating.

- **Reconcile engine (D1)** — one engine drives pull/push across rules,
  firewall, device control, sites, groups, tags, exclusions, and cloud
  policies. Push is create-or-update, matched by a stable per-surface identity.
- **Drift CI gate (D2)** — `s1ctl drift` plans every committed surface and
  exits non-zero when any surface diverges.
- **Breaking (v0.6.0):** sites, groups, tags, exclusions, and cloud policies
  move from single JSON array files to per-object YAML directories; re-pull to
  regenerate local state.

## Wave 10 — Polish and release

- Progress indicators (bubbletea spinners for long operations)
- Rate limiting (x/time/rate)
- Release automation (goreleaser — Win/Mac/Linux x amd64/arm64)
- Docs site updates for all new commands
- **Release v1.0.0**

## Backlog

Surfaces scoped but not yet built:

- **Network control** — list/get/quarantine endpoints
- **Cloud onboarding** — list/get/onboard/delete cloud accounts
- **Automation** — list/get/create/run automation rules
- **Marketplace** — list/get/install integrations
- **Inventory** — unified asset inventory across all types
- **Identity** — identity posture list/get/configure
