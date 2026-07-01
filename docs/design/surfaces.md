# Surfaces

Every CLI command group maps to one or more API surfaces. This page documents
the mapping from SentinelOne's official API taxonomy to s1ctl's command tree.

## Domain groups

### Endpoint security

| CLI group | API tags | Protocol | Plane |
|-----------|----------|----------|-------|
| `agents` | Agents, Agent Actions | REST | Operational |
| `threats` | Threats, Threat Notes, Threat Intelligence | REST | Operational |
| `alerts` | alerts (REST) + UAM (GraphQL) | REST + GraphQL | Operational |
| `sites` | Sites | REST | Operational |
| `groups` | Groups | REST | Operational |
| `accounts` | Accounts | REST | Operational |
| `policies` | Policies | REST | Control |
| `exclusions` | Exclusions and Blocklist, Exclusions v2.1 | REST | Control |

### Detection and response

| CLI group | API tags | Protocol | Plane |
|-----------|----------|----------|-------|
| `rules` | Custom Detection Rule, Platform Detection Rules | REST | Control |
| `visibility` | Deep Visibility | REST | Operational |
| `remoteops` | RemoteOps Scripts, RemoteOps Forensics, Remote Ops MMS | REST | Operational |

### Application and device control

| CLI group | API tags | Protocol | Plane |
|-----------|----------|----------|-------|
| `applications` | Application, Application Management, Application Control | REST | Operational |
| `devices` | Device Control | REST | Operational |
| `firewall` | Firewall Control | REST | Control |
| `network` | Network Quarantine Control, Network Discovery | REST | Operational |

### Cloud and vulnerability management

| CLI group | API | Protocol | Plane |
|-----------|-----|----------|-------|
| `vulnerabilities` | xSPM Vulnerabilities | GraphQL | Operational |
| `misconfigurations` | xSPM Misconfigurations | GraphQL | Operational |
| `cloud` | Cloud Policies, Cloud Onboarding, Cloud Compliance | GraphQL | Mixed |

### Data lake

| CLI group | API | Protocol | Plane |
|-----------|-----|----------|-------|
| `datalake` | SDL (query, powerQuery, addEvents, files) | SDL (REST + GraphQL) | Operational |

### Platform administration

| CLI group | API tags | Protocol | Plane |
|-----------|----------|----------|-------|
| `users` | Users, RBAC, Service Users | REST | Operational |
| `settings` | Settings, Config Overrides | REST | Control |
| `updates` | Updates, Auto Upgrade Policy, Sentinel Deploy | REST | Operational |
| `tags` | Tags, Tag Manager, Dynamic tag rules | REST | Operational |
| `activities` | Activities | REST | Operational |

### Automation and integration

| CLI group | API tags | Protocol | Plane |
|-----------|----------|----------|-------|
| `automation` | Hyperautomation | REST | Operational |
| `marketplace` | marketplace | REST | Operational |

### Asset inventory

| CLI group | API tags | Protocol | Plane |
|-----------|----------|----------|-------|
| `inventory` | Inventory (all subtypes) | REST | Operational |

### Identity

| CLI group | API tags | Protocol | Plane |
|-----------|----------|----------|-------|
| `identity` | Identity AD Service, ISPM | REST | Operational |

## Protocol routing

The CLI defaults to the protocol with the best performance and richest
filtering. Users can override with `--protocol rest|graphql|sdl` on any
command that supports multiple protocols.

| Command | Default protocol | Reason | Override available |
|---------|------------------|--------|--------------------|
| `s1ctl agents list` | REST | Only available via REST | No |
| `s1ctl alerts query` | GraphQL UAM | Richer filtering, grouping | Yes (`--protocol rest`) |
| `s1ctl alerts list` | REST | Simpler, lower latency for basic lists | Yes (`--protocol graphql`) |
| `s1ctl vulnerabilities list` | GraphQL xSPM | Only available via GraphQL | No |
| `s1ctl datalake query` | SDL REST | Only available via SDL REST | No |
| `s1ctl datalake powerquery` | GraphQL | Lower latency, no separate SDL URL needed | Yes (`--protocol rest`) |
| `s1ctl cloud policies list` | GraphQL | Only available via GraphQL | No |

## Config-as-code surfaces

These surfaces support the `pull` / `push` loop:

| Surface | Identity key | Reconcile model |
|---------|-------------|-----------------|
| Exclusions | Server ID | Full CRUD |
| Custom rules | Rule ID | Full CRUD |
| Policies | Policy ID | Update only (policies are site-scoped) |
| Firewall rules | Rule ID | Full CRUD |
| Cloud policies | Policy ID | Full CRUD |
| Settings | Key name | Update only |

All mutations are dry-run by default. Pass `--yes` to apply.
