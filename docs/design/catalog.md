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

## Endpoint security

| Surface | Read | Write | Config-as-code | Status |
|---------|------|-------|----------------|--------|
| agents | list, get, count, outdated, versions, health | isolate, reconnect (by ID or filter), scan, abort-scan, decommission, uninstall, approve/reject-uninstall, upgrade (by package ID, file name, or path), move, fetch-logs, restart, shutdown, enable, disable, reset-config, mark-up-to-date, set-external-id, randomize-uuid, firewall-logging | -- | built |
| threats | list, get, count, notes, timeline | mitigate, verdict, status, resolve, add-note, blacklist, fetch-file | -- | built |
| alerts | list, get, count, history, stats (GraphQL) | status, verdict, resolve, add-note | -- | built |
| sites | list, get, count, licenses | create, update, delete | pull/push | built |
| groups | list, get, count | create, update, delete | pull/push | built |
| accounts | list, get, count | -- | -- | built |
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
| remoteops | list, get, results | run | -- | built |

## Application and device control

| Surface | Read | Write | Config-as-code | Status |
|---------|------|-------|----------------|--------|
| applications | list, cves, risks | -- | -- | built |
| devicecontrol | list, get, events | create, update, delete, enable, disable, reorder, copy | pull/push | built |
| firewall | list, get, protocols, export | create, update, delete, enable, disable, reorder, copy, import | pull/push | built |
| network | list, get | quarantine | -- | -- |

## Cloud and vulnerability management

| Surface | Read | Write | Config-as-code | Status |
|---------|------|-------|----------------|--------|
| vulnerabilities | list, get, health (GraphQL) | status, verdict | -- | built |
| misconfigurations | list, get (GraphQL) | status, verdict | -- | built |
| cloud policies | list, get (GraphQL) | enable, disable, delete | pull/push | built |
| cloud onboarding | list, get | onboard, delete | -- | designed |
| cloud compliance | -- | -- | -- | blocked |

## Data lake

| Surface | Read | Write | Config-as-code | Status |
|---------|------|-------|----------------|--------|
| datalake | powerquery (GraphQL + REST), query, facet, timeseries, saved-queries | addEvents, uploadLogs | -- | built |
| files | getFile, listFiles | putFile | -- | built |

## Platform administration

| Surface | Read | Write | Config-as-code | Status |
|---------|------|-------|----------------|--------|
| users | list, get, token-details | update, delete, generate-token, revoke-token, 2fa | -- | built |
| service-users | list, get, export | create, update, delete, bulk-delete, generate-token | -- | built |
| roles | list, get, template | create, update, delete | -- | built |
| settings | list, get | update, test | -- | built |
| updates | list, get | -- | -- | built |
| upgrade-policies | list, get, packages | create, update, delete, activate, deactivate | -- | built |
| tags | list, get | create, update, delete | pull/push | built |
| activities | list, count, export, types | -- | -- | built |
| audit | list (local mutation log) | -- | -- | built |
| drift | drift summary (all sync surfaces) | -- | -- | built |
| reports | list, tasks, types, download | create | -- | built |
| system | info | -- | -- | built |

## Automation and integration

| Surface | Read | Write | Config-as-code | Status |
|---------|------|-------|----------------|--------|
| automation | list, get | create, run | -- | -- |
| marketplace | list, get | install | -- | -- |

## Asset inventory

| Surface | Read | Write | Config-as-code | Status |
|---------|------|-------|----------------|--------|
| assets (XDR) | overview, categories | -- | -- | built |
| inventory | list, get (all types) | tags, actions | -- | -- |

## Identity

| Surface | Read | Write | Config-as-code | Status |
|---------|------|-------|----------------|--------|
| ranger-ad | status, exposures, affected-objects | assess | -- | built |
| identity | list, get | configure | -- | -- |
