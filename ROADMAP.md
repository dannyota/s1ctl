# Roadmap

Wave-based delivery. Updated each time a wave completes.

## Wave 1 — SDK (complete)

Build all three SDK packages as independently importable Go libraries.

| Package | Surface | Methods |
|---------|---------|---------|
| `mgmt/` | REST MGMT API v2.1 | 297 (CRUD + actions across 40+ resources) |
| `sdl/` | SDL REST + GraphQL | 19 (query, ingest, files, powerquery, GraphQL ops) |
| `graphql/` | GraphQL APIs (6 domains) | 74 (alerts, misconfigs, vulns, cloud policies, CNS rules, DLP) |
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

## Wave 10 — Polish and release (in progress)

- Progress indicators (charmbracelet spinners for long operations) — done
- Rate limiting (x/time/rate) — done
- Release automation (goreleaser — Win/Mac/Linux x amd64/arm64) — done
- Docs site updates for all new commands — done
- **Release v1.0.0** — remaining

## Wave 11 — API coverage expansion (complete)

Broad surface coverage across the REST, GraphQL, and SDL APIs: new command
groups and enrichment of existing surfaces.

- **New command groups:** blocklist (CRUD + pull/push), service-users (CRUD +
  token generation), roles (CRUD), network quarantine (full mirror of firewall),
  DLP (rules + classifications), cloud-rules (CNS custom), filters, locations
  (CRUD + pull/push), tag-rules, maintenance windows
- **Enriched surfaces:** agents (broadcast, fetch-files, fetch-installed-apps,
  fetch-firewall-rules, reset-passphrase, ranger, local-upgrade, passphrases),
  threats (add-to-exclusions, mitigate-alerts, set-ticket, quarantined-files,
  exclusion-options, export), alerts (notes, timeline, counts, export),
  misconfigurations + vulnerabilities (notes, assign, history, related-assets,
  export, CVE queries), remoteops (update, content, upload-limits, pending
  approve/decline, guardrails), settings (sms, recipients, AD, AD scope
  mapping, SSO cert, cancel-pending-emails), sites (reactivate, expire,
  duplicate, regenerate-key, token), accounts (reactivate, expire,
  uninstall-password), users (update, generate-token, revoke-token,
  token-details, 2FA), datalake (numeric, dashboards, saved-queries delete)
- **Breaking:** `datalake saved-queries` is now a subcommand group
  (`saved-queries list`, `saved-queries delete`) — the previous
  `saved-queries` flat command is removed
- **Roles note:** RBAC roles surface is CRUD-only (list/get/create/update/
  delete); pull/push reconcile was planned but review found the round-trip
  unsafe due to permission-tree normalization

## Wave 12 — Cloud, config-as-code, identity, and automation (complete)

New command groups and config-as-code expansion across cloud, identity,
data lake, and automation surfaces.

- **Cloud onboarding (GraphQL)** — list/get/onboard/delete cloud entities
  (AWS, GCP, Azure, OCI, Alibaba)
- **Application control + config-as-code** — rules CRUD + pull/push, settings,
  management settings, labels
- **Config overrides + Sentinel Deploy** — settings overrides CRUD,
  updates deploy credential groups/details
- **Config-as-code expansion** — upgrade-policies (scope-partitioned pull/push),
  filters, tag-rules; maintenance evaluated and excluded as singleton config
- **SDL parsers + notebooks** — datalake parsers and notebooks subgroups
  (list/get/delete)
- **Identity AD Service + ISPM** — config, connector, onboard, domains,
  features, timezones, skip-exposures, ack-exposures
- **Automation (Hyperautomation)** — list/get/versions/export/create/run/
  activate/deactivate/executions
- **Deferred:** Unified Incidents and Unified Tags (endpoints inaccessible on
  the console)

## Backlog

Surfaces scoped but not yet built:

- **Marketplace** — list/get/install integrations
- **Inventory** — unified asset inventory across all types (the `assets` overview surface covers a first slice)
