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
| agents | list, get, count | isolate, connect, scan, decommission | -- | built |
| threats | list, get | mitigate, verdict, status | -- | built |
| alerts | list (GraphQL) | -- | -- | built |
| sites | list, get | -- | -- | built |
| groups | list, get | -- | -- | built |
| accounts | list, get | -- | -- | built |
| policies | get | -- | -- | built |
| exclusions | list, get | create, delete | pull/push | built |

## Detection and response

| Surface | Read | Write | Config-as-code | Status |
|---------|------|-------|----------------|--------|
| rules | list, get | create, update, delete | pull/push | designed |
| visibility | query | -- | -- | designed |
| remoteops | list | -- | -- | built |

## Application and device control

| Surface | Read | Write | Config-as-code | Status |
|---------|------|-------|----------------|--------|
| applications | list | -- | -- | built |
| devices | list | -- | -- | built |
| firewall | list | -- | -- | built |
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
| datalake | powerquery | -- | -- | built |

## Platform administration

| Surface | Read | Write | Config-as-code | Status |
|---------|------|-------|----------------|--------|
| users | list, get | -- | -- | built |
| settings | list, get | update | pull/push | -- |
| updates | list | -- | -- | built |
| tags | list | -- | -- | built |
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
