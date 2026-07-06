---
name: s1ctl
description: >
  Operating guide for AI agents driving the s1ctl CLI (v0.7.3) against a
  SentinelOne Singularity Platform instance. Covers the three API surfaces,
  the full 370-command map across 48 groups, the mutation ritual, the
  config-as-code reconcile engine with drift detection, secret-output
  conventions, self-discovery commands, output contracts, end-to-end recipes,
  and gotchas the per-command --help can't express. Read this before issuing
  any s1ctl command.
---

# s1ctl agent operating guide

s1ctl operates a SentinelOne Singularity Platform instance as code — agents,
threats, alerts, policies, exclusions, detections, device control, firewall,
network quarantine, cloud security (CNS rules, DLP), and more. This guide
makes you productive without the repo docs; the live commands (`commands
--json`, `status surfaces`, `<cmd> --help`) are the source of truth when
something here looks out of date.

## Session bootstrap — do these first

```bash
s1ctl doctor                      # config + auth + API reachability (read-only)
s1ctl status capabilities --json  # version, auth health, surface status
s1ctl commands --json             # every verb: name, kind (read/guarded-mutation), flags
```

`doctor` is the gate: if it reports unhealthy, fix config/auth before
proceeding. `commands --json` is the **live source of truth** for what this
binary supports (370 commands across 48 groups at v0.7.3).

## The three API surfaces

All authenticated via `S1_TOKEN` (API token from the console):

| Surface | Auth | Transport | Key groups |
|---|---|---|---|
| **REST MGMT** (v2.1) | `ApiToken <token>` | REST | agents, threats, sites, groups, accounts, policies, exclusions, blocklist, rules, settings, reports, tags, users, service-users, roles, activities, devicecontrol, firewall, network, upgrade-policies, IOCs, remoteops, filters, locations, tag-rules, maintenance |
| **SDL / Data Lake** | `Bearer <token>` | GraphQL + REST | datalake (powerquery, query, numeric, facet, timeseries, ingest, files, dashboards, saved-queries) |
| **GraphQL** | `Bearer <token>` | GraphQL | alerts (UAM), misconfigurations, vulnerabilities (xSPM), cloud-policies, cloud-rules (CNS), dlp, assets (XDR) |

### Config resolution (highest priority first)

`S1_*` env vars → `--config <path>` → `~/.s1ctl/config.yaml` →
`./config/config.yaml`.

Required: `S1_CONSOLE_URL` (e.g. `https://your-console.sentinelone.net`) and
`S1_TOKEN`. Persist in `.env` (gitignored) and `source .env` before running.

## Command map

370 commands across 48 groups. Run `commands --json` for the full live catalog.
Kind: **read** is always safe; **guarded-mutation** needs the dry-run → `--yes`
ritual.

| Group | Verbs |
|---|---|
| **accounts** | count, expire, get, list, reactivate, uninstall-password (generate/revoke/show) |
| **activities** | count, export, list, types |
| **agents** | abort-scan, approve-uninstall, broadcast, count, decommission, disable, enable, fetch-files, fetch-firewall-rules, fetch-installed-apps, fetch-logs, firewall-logging, get, health, isolate, list, local-upgrade, local-upgrade-status, mark-up-to-date, move, move-to-site, outdated, passphrases, randomize-uuid, ranger, reconnect, reject-uninstall, reset-config, reset-passphrase, restart, scan, set-external-id, shutdown, uninstall, upgrade, versions |
| **alerts** | add-note, count, counts, delete-note, export, get, history, list, notes, resolve, stats, status, timeline, update-note, verdict |
| **applications** | cves, list, risks |
| **assets** | categories, overview |
| **blocklist** | create, delete, export, list, pull, push, update, validate |
| **cloud-policies** | delete, disable, enable, get, list, pull, push |
| **cloud-rules** | create, delete, disable, enable, evaluate, get, list, types, update |
| **datalake** | dashboards (get/list), facet, files (get/list/put), ingest (events/logs), numeric, powerquery, query, saved-queries (delete/list), timeseries |
| **detection-library** | data-sources, disable, enable, list, surfaces |
| **devicecontrol** | copy, delete, disable, enable, events, get, list, pull, push, reorder |
| **dlp** | classifications (delete/get/list), rules (delete/disable/enable/get/list), settings |
| **drift** | (top-level — per-surface drift detection) |
| **exclusions** | create, delete, get, list, pull, push, update |
| **filters** | create, delete, list, update |
| **firewall** | copy, delete, disable, enable, export, get, import, list, protocols, pull, push, reorder |
| **groups** | count, create, delete, get, list, pull, push, update |
| **iocs** | config, create, delete, list |
| **locations** | create, delete, list, pull, push, update |
| **maintenance** | export, get, get-flexible, set, set-flexible |
| **misconfigurations** | add-note, assign, delete-note, export, get, history, list, notes, related-assets, status, update-note, verdict |
| **network** | configuration (get/set), copy, delete, disable, enable, export, get, import, list, move, protocols, pull, push, reorder, set-location, tags (add/remove) |
| **policies** | diff, get, list, pull, push, revert |
| **remoteops** | content, get, guardrails (check/delete/get/set), list, pending (approve/decline/list), results, run, update, upload-limits |
| **reports** | create, download, list, tasks, types |
| **roles** | create, delete, get, list, template, update |
| **rules** | detections, diff, disable, enable, get, health, list, pull, push, trends, validate |
| **service-users** | bulk-delete, create, delete, export, generate-token, get, list, update |
| **settings** | cancel-pending-emails, delete-recipient, get, list, sso-cert, test, update (ad/ad-scope-mapping/notifications/recipients/sms/smtp/sso/syslog) |
| **sites** | count, create, delete, duplicate, expire, get, licenses, list, pull, push, reactivate, regenerate-key, token, update |
| **tag-rules** | create, delete, list, test, update |
| **tags** | create, delete, get, list, pull, push, update |
| **threats** | add-note, add-to-exclusions, blacklist, count, exclusion-options, export, fetch-file, get, list, mitigate, mitigate-alerts, notes, quarantined-files, resolve, set-ticket, status, timeline, verdict |
| **unified-exclusions** | create, export, list |
| **upgrade-policies** | activate, create, deactivate, delete, get, list, packages, update |
| **users** | 2fa (disable/enable), delete, generate-token, get, list, revoke-token, token-details, update |
| **vulnerabilities** | add-note, assign, cve, cves, delete-note, export, get, health, history, list, notes, related-assets, stats, status, update-note, verdict |
| **status** | capabilities, enums, surfaces |
| **system** | info |
| **visibility** | query |

Config-as-code surfaces (pull/push via the reconcile engine): `blocklist`,
`cloud-policies`, `devicecontrol`, `exclusions`, `firewall`, `groups`,
`locations`, `network`, `rules`, `sites`, `tags`.

## The mutation ritual — every guarded verb

Every verb that changes live state is **guarded**: it defaults to a dry-run
preview; pass `--yes` to apply.

```bash
s1ctl agents upgrade --package-id PKG --site-id SITE   # 1. dry-run preview
s1ctl agents upgrade --package-id PKG --site-id SITE --yes  # 2. apply
```

**Never skip the preview.** A mutation is a production action against a live
console.

### Read-only mode

Launch with `S1_READONLY=1` (or `--read-only`): every guarded mutation
degrades to a dry-run even with `--yes`. Every mutation attempt is logged
to `~/.s1ctl/audit.jsonl`.

## Output contracts

- `--json` on any read command for parseable output.
- `--output csv` on list commands for CSV export.
- `--all` fetches every page (not just the first).
- `--out <file>` writes output to a file (on export/download commands).
- Under `--json`, dry-run emits `{"dryRun": true, "command", "action", "target"}`.
- Non-2xx errors return typed `APIError` with status, title, detail.

## The config-as-code reconcile loop

```text
pull live state  →  review in git diff  →  push back  →  drift to verify
```

All 11 sync surfaces use the same reconcile engine: one YAML file per object,
name-matched plans (create/update/unchanged/live-only), no deletes. Push is
dry-run by default; `--yes` applies creates and updates, skipping unchanged
objects.

```bash
s1ctl firewall pull --out firewall/           # snapshot live rules to YAML
git diff firewall/                             # the review surface
s1ctl firewall push --dir firewall/            # dry-run preview
s1ctl firewall push --dir firewall/ --yes      # deploy creates + updates
s1ctl drift                                    # verify: exit 0 = clean
```

Same loop for all sync surfaces: `blocklist`, `cloud-policies`,
`devicecontrol`, `exclusions`, `firewall`, `groups`, `locations`, `network`,
`rules`, `sites`, `tags`.

### Drift detection (CI gate)

```bash
s1ctl drift                                    # all surfaces with local dirs
s1ctl drift --surface firewall --surface rules  # subset
s1ctl drift --dir-root /path/to/config         # custom root
```

Exit 0 = clean; exit 1 = drift detected (creates, updates, or live-only).
Read-only — lists and plans, never applies. Surfaces without a local directory
are skipped.

## Secret-output conventions

Commands that print sensitive values (tokens, passwords, passphrases):

| Command | Prints to stdout | Stderr notice |
|---|---|---|
| `sites token`, `sites regenerate-key` | registration token | yes |
| `accounts uninstall-password show/generate` | uninstall password | yes |
| `agents passphrases` | agent passphrases | yes |
| `users generate-token`, `service-users generate-token` | API token (once) | yes |
| `service-users create` (if token auto-generated) | API token | yes |

These values are NEVER written to the audit log. The `settings get` commands
redact secrets (SMTP password, syslog token/certs, AD bind password) on the
`--json` path.

## Common recipes

### Agent fleet management

```bash
s1ctl agents list --site-id SITE --json
s1ctl agents outdated --site-id SITE
s1ctl agents health --site-id SITE
s1ctl agents upgrade --package-id PKG --site-id SITE --yes
s1ctl agents isolate AGENT-ID --yes
s1ctl agents reconnect AGENT-ID --yes
s1ctl agents scan AGENT-ID --yes
s1ctl agents broadcast AGENT-ID --message "Reboot at 2am" --yes
s1ctl agents passphrases --site-id SITE         # SECRET
s1ctl agents ranger AGENT-ID --state on --yes   # toggle network discovery
```

### Threat triage

```bash
s1ctl threats list --site-id SITE --json
s1ctl threats get THREAT-ID --json
s1ctl threats timeline THREAT-ID
s1ctl threats notes THREAT-ID
s1ctl threats mitigate THREAT-ID --action remediate --yes
s1ctl threats verdict THREAT-ID --verdict true_positive --yes
s1ctl threats export --out threats.csv           # CSV export
s1ctl threats set-ticket THREAT-ID --ticket-id JIRA-123 --yes
```

### Alert triage (GraphQL)

```bash
s1ctl alerts list --json
s1ctl alerts history ALERT-ID
s1ctl alerts stats
s1ctl alerts timeline ALERT-ID
s1ctl alerts notes ALERT-ID
s1ctl alerts add-note ALERT-ID "investigation note" --yes
s1ctl alerts export --out alerts.csv
s1ctl alerts status ALERT-ID --status resolved --yes
```

### Cloud security

```bash
s1ctl cloud-policies list --json
s1ctl cloud-policies pull --out cloud-policies/
s1ctl cloud-rules list --json                    # CNS custom rules
s1ctl cloud-rules create --from-file rule.json --yes
s1ctl cloud-rules evaluate --rule rule.json --resource resource.json
s1ctl dlp rules list --json                      # DLP data-protection rules
s1ctl dlp classifications list
s1ctl dlp settings --scope-level site --scope-id SITE
s1ctl misconfigurations list --json
s1ctl vulnerabilities list --json
s1ctl vulnerabilities cves --json                # CVE inventory
s1ctl vulnerabilities stats                      # top vulns summary
```

### Policy and exclusion management

```bash
s1ctl policies list
s1ctl policies diff --site-id S1 --site-id S2
s1ctl policies pull --out policies/
s1ctl policies push --dir policies/ --yes
s1ctl exclusions pull --out exclusions/
s1ctl exclusions push --dir exclusions/ --yes
s1ctl blocklist list --json
s1ctl blocklist pull --out blocklist/
s1ctl blocklist push --dir blocklist/ --yes
```

### Network and device control

```bash
s1ctl firewall pull --out firewall/
s1ctl firewall push --dir firewall/ --yes
s1ctl network list --site-id SITE               # network quarantine rules
s1ctl network pull --out network-quarantine/
s1ctl network push --dir network-quarantine/ --yes
s1ctl devicecontrol pull --out devicecontrol/
s1ctl devicecontrol push --dir devicecontrol/ --yes
```

### Platform administration

```bash
s1ctl sites list --json
s1ctl sites reactivate SITE-ID --unlimited --yes
s1ctl sites token SITE-ID                        # SECRET: registration token
s1ctl accounts uninstall-password show ACCT-ID   # SECRET
s1ctl service-users list --json
s1ctl service-users create --name bot --scope tenant --yes
s1ctl roles list --json
s1ctl roles create --from-file role.json --yes
s1ctl users update USER-ID --full-name "New Name" --yes
s1ctl users 2fa enable USER-ID --yes
s1ctl settings get ad --json                     # AD settings (password redacted)
s1ctl settings update ad --from-file ad.json --yes
s1ctl settings sso-cert --out cert.pem
s1ctl maintenance get --task-type AgentSoftwareUpdate --json
```

### Remote operations

```bash
s1ctl remoteops list --json
s1ctl remoteops run SCRIPT-ID --site-id SITE --yes
s1ctl remoteops results TASK-ID
s1ctl remoteops content SCRIPT-ID               # view script source
s1ctl remoteops pending list                     # pending approvals
s1ctl remoteops pending approve EXEC-ID --yes
s1ctl remoteops guardrails get --scope-id SITE --scope-level site
```

### Data Lake queries

```bash
s1ctl datalake powerquery --query '* | limit 10' --from 24h
s1ctl datalake query --filter 'severity >= 4' --start 24h --max 1000
s1ctl datalake facet --field serverHost --start 24h
s1ctl datalake timeseries --filter '*' --start 24h --buckets 12
s1ctl datalake numeric --filter '*' --function rate --start 24h
s1ctl datalake ingest events --file events.json --session s1 --yes
s1ctl datalake dashboards list --json
s1ctl datalake saved-queries list
s1ctl datalake saved-queries delete --name Q --type saved --index 0 --yes
```

## Self-discovery commands

Do not guess command names, flags, or enums — read the live catalog:

- `s1ctl commands --json` — every verb with kind, flags, and descriptions.
- `s1ctl status capabilities --json` — version, auth health, surface status.
- `s1ctl status surfaces` — every surface with protocol and status.
- `s1ctl status enums` — typed enum values used across the API.
- `<cmd> --help` — per-command flags and usage.

## Gotchas

### Agent upgrade requires package identification

`agents upgrade` needs exactly one of `--package-id`, `--file-name`, or
`--path`. When using `--file-name`, you must also pass `--os-type`. Get
valid package IDs from `upgrade-policies packages`.

### Sites/accounts reactivate requires explicit expiration choice

`sites reactivate` and `accounts reactivate` require exactly one of
`--unlimited` or `--expiration <RFC3339>` — neither defaults to the other.

### Cloud policy/CNS rules empty IDs means "all"

The cloud-policies and CNS-rules APIs interpret an empty ID list as "act on
ALL rules". The CLI rejects empty id lists to prevent accidental bulk
operations. Same guard on DLP bulk operations.

### Roles are CRUD-only, not config-as-code

`roles` has list/get/template/create/update/delete but NO pull/push — the API
list endpoint lacks permission data, so the reconcile round-trip is unsound.
Use `roles create --from-file` / `roles update --from-file` directly.

### Settings get redacts secrets

`settings get smtp/syslog/ad` blanks password/token/cert fields in `--json`
output. When doing a get→edit→update round-trip, re-enter secrets before
pushing, or the update writes them back empty.

### SDL QueryAll pagination

QueryAll terminates on empty matches OR a non-advancing continuation token.
Sessions are merged from every page.

### Rate limiting

The console enforces rate limits. Use `WithRateLimit(rps, burst)` in SDK
code. The CLI handles retries automatically.

## Safety rules

- **No mutation without the dry-run review.** Show the preview, then `--yes`.
- **Never commit real identifiers.** Use placeholders in code and docs.
- **No secrets in the repo.** `.env` and tokens never committed.
- **All live testing targets UAT only.** Never run mutations against
  production without explicit approval.

## Quick reference

| I want to… | Command |
|---|---|
| Verify setup | `s1ctl doctor` |
| Discover commands | `s1ctl commands --json` |
| Check drift | `s1ctl drift` |
| List agents | `s1ctl agents list --site-id SITE` |
| Upgrade agents | `s1ctl agents upgrade --package-id PKG --site-id SITE --yes` |
| Triage threats | `threats list` → `threats get ID` → `threats mitigate ID --yes` |
| Triage alerts | `alerts list` → `alerts get ID` → `alerts status ID --yes` |
| Sync firewall rules | `firewall pull` → edit → `firewall push --yes` → `drift` |
| Manage blocklist | `blocklist pull` → edit → `blocklist push --yes` |
| Query data lake | `datalake powerquery --query '...' --from 24h` |
| Cloud posture | `misconfigurations list` / `vulnerabilities list` / `cloud-rules list` |
| DLP rules | `dlp rules list` / `dlp classifications list` |
| Manage IOCs | `iocs list` / `iocs create` / `iocs delete` |
| Platform admin | `sites list` / `roles list` / `service-users list` / `settings list` |
| Hard read-only | `S1_READONLY=1 s1ctl ...` |
