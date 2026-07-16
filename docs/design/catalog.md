# Catalog

Implementation status for every surface. Updated in the same commit that moves
a surface forward. Read/Write columns list the union of CLI verbs and SDK
methods for the surface; config-as-code marks surfaces with `pull`/`push`.

## Status legend

| Status | Meaning |
|--------|---------|
| **designed** | API mapped, CLI shape decided, not yet built |
| **built** | Code exists, passes tests |
| **verified** | Tested against a live console |
| **blocked** | API limitation or missing access |
| **--** | Not yet scoped |

## Foundation

| Surface | Commands | Status |
|---------|----------|--------|
| status | health summary dashboard; capabilities, enums, surfaces | built |
| version | version info | built |
| doctor | config diagnostics | built |
| config | init wizard, show | built |
| commands | list all commands | built |
| completion | shell completions | built |
| docs | generate (command reference; hidden) | built |
| mcp | serve (stdio MCP server), install (Claude Code config) | built |

## Endpoint security

| Surface | Read | Write | Config-as-code | Status |
|---------|------|-------|----------------|--------|
| agents | list, get, count, outdated, versions, health, passphrases, local-upgrade-status | isolate, reconnect (by ID or filter), scan, abort-scan, decommission, uninstall, approve/reject-uninstall, upgrade (by package ID, file name, or path), move, fetch-logs, restart, shutdown, enable, disable, reset-config, mark-up-to-date, set-external-id, randomize-uuid, firewall-logging, broadcast, reset-passphrase, ranger, fetch-installed-apps, fetch-firewall-rules, fetch-files, local-upgrade | -- | built |
| threats | list, get, count, notes, timeline, quarantined-files, exclusion-options, export | mitigate, verdict, status, resolve, add-note, blacklist, fetch-file, add-to-exclusions, mitigate-alerts, set-ticket | -- | built |
| alerts | list, get, count, history, stats, notes, timeline, counts, export (GraphQL) | status, verdict, resolve, add-note, note-update, note-delete | -- | built |
| sites | list, get, count, licenses, token | create, update, delete, reactivate, expire, duplicate, regenerate-key | pull/push | built |
| groups | list, get, count | create, update, delete | pull/push | built |
| accounts | list, get, count, uninstall-password | reactivate, expire, uninstall-password generate/revoke | -- | built |
| policies | list, get, diff (site, account, group scopes) | update, revert (per scope) | pull/push | built |
| exclusions | list, get | create, update, delete | pull/push | built |
| unified-exclusions | list, count, export | create | -- | built |
| blocklist | list, export | create, update, delete, validate | pull/push | built |

## Detection and response

| Surface | Read | Write | Config-as-code | Status |
|---------|------|-------|----------------|--------|
| rules | list, get, health, trends, detections, diff, validate | create, update, enable, disable | pull/push | built |
| detection-library | list, surfaces, data-sources | enable, disable | -- | built |
| iocs | list, config | create, delete | -- | built |
| visibility | query | -- | -- | built |
| remoteops | list, get, results, content, upload-limits, pending, guardrails | run, update, pending approve/decline, guardrails set/delete/check | -- | built |

## Application and device control

| Surface | Read | Write | Config-as-code | Status |
|---------|------|-------|----------------|--------|
| applications | list, cves, risks, rules list/get, settings get, labels list, mgmt-settings get | rules create/update/delete, settings update, mgmt-settings update | rules pull/push | built |
| devicecontrol | list, get, events | create, update, delete, enable, disable, reorder, copy | pull/push | built |
| firewall | list, get, protocols, export | create, update, delete, enable, disable, reorder, copy, import | pull/push | built |
| network | list, get, protocols, configuration get, export | create, update, delete, enable, disable, reorder, copy, import, set-location, move, tags, configuration set | pull/push | built |
| locations | list | create, update, delete | pull/push | built |

## Cloud and vulnerability management

| Surface | Read | Write | Config-as-code | Status |
|---------|------|-------|----------------|--------|
| vulnerabilities | list, get, health, notes, history, related-assets, export, cves, cve, stats (GraphQL) | status, verdict, note-add, note-update, note-delete, assign | -- | built |
| misconfigurations | list, get, notes, history, related-assets, export (GraphQL) | status, verdict, note-add, note-update, note-delete, assign | -- | built |
| cloud policies | list, get (GraphQL) | enable, disable, delete | pull/push | built |
| cloud rules (CNS) | list, get, types (GraphQL) | create, update, enable, disable, delete, evaluate | -- | built |
| dlp | rules list, rules get, classifications list, classifications get, settings (GraphQL) | rules enable, rules disable, rules delete, classifications delete | -- | built |
| cloud onboarding | list, get (GraphQL) | onboard, delete | -- | built |
| cloud compliance | -- | -- | -- | blocked |

## Data lake

| Surface | Read | Write | Config-as-code | Status |
|---------|------|-------|----------------|--------|
| datalake | powerquery (GraphQL + REST), query, facet, timeseries, numeric, saved-queries list, dashboards (list, get), parsers (list, get), notebooks (list, get) | addEvents, uploadLogs, saved-queries delete, parsers delete, notebooks delete | -- | built |
| files | getFile, listFiles | putFile | -- | built |

## Platform administration

| Surface | Read | Write | Config-as-code | Status |
|---------|------|-------|----------------|--------|
| users | list, get, token-details | update, delete, generate-token, revoke-token, 2fa | -- | built |
| service-users | list, get, export | create, update, delete, bulk-delete, generate-token | -- | built |
| roles | list, get, template | create, update, delete | -- | built |
| settings | list, get (notifications/sso/smtp/syslog/sms/recipients/ad/ad-scope-mapping), sso-cert | update (same set), test (smtp/syslog/ad), delete-recipient, cancel-pending-emails | -- | built |
| config-overrides | list, get | create, update, delete | -- | built |
| updates | list, get | -- | -- | built |
| deploy | list-groups, list-details | create-group, delete-group, add-detail, update-detail, delete-detail | -- | built |
| upgrade-policies | list, get, packages | create, update, delete, activate, deactivate | pull/push | built |
| tags | list, get | create, update, delete | pull/push | built |
| filters | list | create, update, delete | pull/push | built |
| maintenance | get, get-flexible, export | set, set-flexible | -- | built |
| activities | list, count, export, types | -- | -- | built |
| audit | list (local mutation log) | -- | -- | built |
| drift | drift summary (all sync surfaces) | -- | -- | built |
| reports | list, tasks, types, download | create | -- | built |
| system | info | -- | -- | built |

## Automation and integration

| Surface | Read | Write | Config-as-code | Status |
|---------|------|-------|----------------|--------|
| automation | list, get | create, run | -- | designed |
| marketplace | list, get | install | -- | designed |

## Asset inventory

| Surface | Read | Write | Config-as-code | Status |
|---------|------|-------|----------------|--------|
| assets (XDR) | overview, categories | -- | -- | built |
| tag-rules | list, test | create, update, delete | pull/push | built |
| inventory | list, get (all types) | tags, actions | -- | designed |

## Identity

| Surface | Read | Write | Config-as-code | Status |
|---------|------|-------|----------------|--------|
| ranger-ad | status, exposures, affected-objects | assess | -- | built |
| identity | list, get | configure | -- | designed |
