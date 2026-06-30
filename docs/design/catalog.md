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
| agents | list, get, count | actions (isolate, scan, ...) | -- | designed |
| threats | list, get | mitigate, notes | -- | designed |
| alerts | list (REST), query (GraphQL) | triage, notes | -- | designed |
| sites | list, get | create, update, delete | -- | designed |
| groups | list, get | create, update, delete | -- | designed |
| accounts | list, get | -- | -- | designed |
| policies | list, get | update | pull/push | designed |
| exclusions | list, get | create, update, delete | pull/push | designed |

## Detection and response

| Surface | Read | Write | Config-as-code | Status |
|---------|------|-------|----------------|--------|
| rules | list, get | create, update, delete | pull/push | designed |
| visibility | query | -- | -- | designed |
| remoteops | list scripts, results | run, upload | -- | designed |

## Application and device control

| Surface | Read | Write | Config-as-code | Status |
|---------|------|-------|----------------|--------|
| applications | list, get | -- | -- | -- |
| devices | list, get | rules | -- | -- |
| firewall | list, get | create, update, delete | pull/push | -- |
| network | list, get | quarantine | -- | -- |

## Cloud and vulnerability management

| Surface | Read | Write | Config-as-code | Status |
|---------|------|-------|----------------|--------|
| vulnerabilities | list, query, export | triage, notes | -- | designed |
| misconfigurations | list, query, export | triage, notes | -- | designed |
| cloud policies | list, get | create, update, delete | pull/push | designed |
| cloud onboarding | list, get | onboard, delete | -- | designed |
| cloud compliance | -- | -- | -- | blocked |

## Data lake

| Surface | Read | Write | Config-as-code | Status |
|---------|------|-------|----------------|--------|
| datalake | query, powerquery, files | addEvents, uploadLogs | -- | designed |

## Platform administration

| Surface | Read | Write | Config-as-code | Status |
|---------|------|-------|----------------|--------|
| users | list, get | create, update, delete | -- | -- |
| settings | list, get | update | pull/push | -- |
| updates | list, get | upgrade | -- | -- |
| tags | list, get | create, update, delete | -- | -- |
| activities | list | -- | -- | -- |

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
