# Catalog

Implementation status for every surface. Updated in the same commit that moves
a surface forward.

## Status legend

| Status | Meaning |
|--------|---------|
| **designed** | API mapped, CLI shape decided, not yet built |
| **built** | Code exists, passes tests |
| **verified** | Tested against a live console |
| **blocked** | API limitation or missing access |
| **--** | Not yet scoped |

## Endpoint security

| Surface | Read | Write | Config-as-code | Status |
|---------|------|-------|----------------|--------|
| agents | list, get, count | isolate, connect, scan, decommission, update-software, move-to-site, fetch-logs, restart, enable, disable, reset-config, approve/reject-uninstall, mark-up-to-date, set-external-id, randomize-uuid, firewall-logging | -- | built |
| threats | list, get | mitigate, verdict, status, add-to-blacklist, fetch-file | -- | built |
| alerts | list, get (GraphQL) | status, verdict | -- | built |
| sites | list, get | create, update, delete | -- | built |
| groups | list, get | create, update, delete | -- | built |
| accounts | list, get | -- | -- | built |
| policies | get (site, account, group) | update (site, account, group) | -- | built |
| exclusions | list, get | create, update, delete | pull/push | built |

## Detection and response

| Surface | Read | Write | Config-as-code | Status |
|---------|------|-------|----------------|--------|
| rules | list, get | create, update, delete | pull/push | designed |
| visibility | query | -- | -- | designed |
| remoteops | list, get | -- | -- | built |

## Application and device control

| Surface | Read | Write | Config-as-code | Status |
|---------|------|-------|----------------|--------|
| applications | list | -- | -- | built |
| devices | list, get | -- | -- | built |
| firewall | list, get | -- | -- | built |
| network | list, get | quarantine | -- | -- |

## Cloud and vulnerability management

| Surface | Read | Write | Config-as-code | Status |
|---------|------|-------|----------------|--------|
| vulnerabilities | list, get (GraphQL) | status, verdict | -- | built |
| misconfigurations | list, get (GraphQL) | status, verdict | -- | built |
| cloud policies | list, get (GraphQL) | -- | pull/push | built |
| cloud onboarding | list, get | onboard, delete | -- | designed |
| cloud compliance | -- | -- | -- | blocked |

## Data lake

| Surface | Read | Write | Config-as-code | Status |
|---------|------|-------|----------------|--------|
| datalake | powerquery (GraphQL + REST), query, facet, timeseries | addEvents, uploadLogs | -- | built |
| files | getFile, listFiles | putFile | -- | built |

## Platform administration

| Surface | Read | Write | Config-as-code | Status |
|---------|------|-------|----------------|--------|
| users | list, get | delete | -- | built |
| settings | list, get | update | pull/push | -- |
| updates | list, get | -- | -- | built |
| tags | list, get | create, update, delete | -- | built |
| activities | list | -- | -- | built |

## Automation and integration

| Surface | Read | Write | Config-as-code | Status |
|---------|------|-------|----------------|--------|
| automation | list, get | create, run | -- | -- |
| marketplace | list, get | install | -- | -- |

## Asset inventory

| Surface | Read | Write | Config-as-code | Status |
|---------|------|-------|----------------|--------|
| inventory | list, get (all types) | tags, actions | -- | -- |

## Identity

| Surface | Read | Write | Config-as-code | Status |
|---------|------|-------|----------------|--------|
| identity | list, get | configure | -- | -- |
