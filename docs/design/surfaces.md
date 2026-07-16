# Surfaces

Every CLI command group maps to one or more API surfaces. This page documents
the mapping from SentinelOne's official API taxonomy to s1ctl's command tree.

## Domain groups

### Endpoint security

| CLI group | API tags | Protocol | Plane |
|-----------|----------|----------|-------|
| `agents` | Agents, Agent Actions | REST | Operational |
| `threats` | Threats, Threat Notes | REST | Operational |
| `alerts` | UAM unified alerts | GraphQL | Operational |
| `sites` | Sites | REST | Operational |
| `groups` | Groups | REST | Operational |
| `accounts` | Accounts | REST | Operational |
| `policies` | Policies | REST | Control |
| `exclusions` | Exclusions and Blocklist | REST | Control |
| `unified-exclusions` | Exclusions v2.1 | REST | Control |
| `blocklist` | Exclusions and Blocklist (restrictions) | REST | Control |

### Detection and response

| CLI group | API tags | Protocol | Plane |
|-----------|----------|----------|-------|
| `rules` | Custom Detection Rule | REST | Control |
| `detection-library` | Platform Detection Rules | REST | Control |
| `iocs` | Threat Intelligence | REST | Operational |
| `visibility` | Deep Visibility | REST | Operational |
| `remoteops` | RemoteOps Scripts, RemoteOps Forensics, Remote Ops MMS | REST | Operational |

### Application and device control

| CLI group | API tags | Protocol | Plane |
|-----------|----------|----------|-------|
| `applications` | Application, Application Management, Application Risk, Application Control - Rules, Application Control - Settings and Labels, Application Management Settings | REST | Operational + Control |
| `devicecontrol` | Device Control | REST | Control |
| `firewall` | Firewall Control | REST | Control |
| `network` | Network Quarantine Control | REST | Control |
| `locations` | Locations | REST | Control |

### Cloud and vulnerability management

| CLI group | API | Protocol | Plane |
|-----------|-----|----------|-------|
| `vulnerabilities` | xSPM Vulnerabilities | GraphQL | Operational |
| `misconfigurations` | xSPM Misconfigurations | GraphQL | Operational |
| `cloud-policies` | Cloud Security Policies | GraphQL | Control |
| `cloud-rules` | CNS Custom Rules | GraphQL | Control |
| `dlp` | Data Protection Rules, DLP Classifications | GraphQL | Control |

### Data lake

| CLI group | API | Protocol | Plane |
|-----------|-----|----------|-------|
| `datalake` | SDL (query, powerQuery, addEvents, files, dashboards) | SDL (REST + GraphQL) | Operational |

### Asset inventory and identity

| CLI group | API tags | Protocol | Plane |
|-----------|----------|----------|-------|
| `assets` | XDR assets | REST | Operational |
| `tag-rules` | Dynamic tag rules | REST | Control |
| `ranger-ad` | Ranger AD (`/ranger-ad`) | REST | Operational |
| `identity` | Identity AD Service - Configuration/Connector/Onboarding, ISPM | REST | Control |

### Platform administration

| CLI group | API tags | Protocol | Plane |
|-----------|----------|----------|-------|
| `users` | Users | REST | Operational |
| `service-users` | Service Users | REST | Operational |
| `roles` | RBAC | REST | Control |
| `settings` | Settings | REST | Control |
| `updates` | Updates | REST | Operational |
| `upgrade-policies` | Auto Upgrade Policy | REST | Control |
| `tags` | Tags, Tag Manager | REST | Operational |
| `filters` | Filters | REST | Control |
| `maintenance` | Tasks (`/tasks-configuration`) | REST | Control |
| `activities` | Activities | REST | Operational |
| `reports` | Default Reports | REST | Operational |
| `system` | System | REST | Operational |

### Local and meta commands

These groups run locally and call no API surface (except `doctor` and
`status`, which probe connectivity):

| CLI group | Purpose |
|-----------|---------|
| `status` | Health summary dashboard; capabilities, enums, surfaces |
| `doctor` | Config and connectivity diagnostics |
| `config` | Init wizard, show |
| `drift` | Plan every committed sync surface, exit non-zero on divergence |
| `audit` | Local mutation log |
| `mcp` | MCP server (serve, install) |
| `commands`, `completion`, `version` | Catalog, shell completion, version |
| `docs` | Command-reference generation (hidden) |

### Automation and integration

| CLI group | API tags | Protocol | Plane |
|-----------|----------|----------|-------|
| `automation` | Hyperautomation | REST | Operational |

### Planned groups

Scoped in the roadmap backlog, not yet built: `marketplace`,
`inventory` (Inventory suite).

## Protocol routing

When a surface is available via multiple protocols, the CLI defaults to the
protocol with the best performance and richest filtering, and exposes
`--protocol` on commands that support more than one.

| Command | Default protocol | Reason | Override available |
|---------|------------------|--------|--------------------|
| `s1ctl agents list` | REST | Only available via REST | No |
| `s1ctl alerts list` | GraphQL UAM | Richer filtering, grouping | No |
| `s1ctl vulnerabilities list` | GraphQL xSPM | Only available via GraphQL | No |
| `s1ctl datalake query` | SDL REST | Only available via SDL REST | No |
| `s1ctl datalake powerquery` | GraphQL | Lower latency, no separate SDL URL needed | Yes (`--protocol rest`) |
| `s1ctl cloud-policies list` | GraphQL | Only available via GraphQL | No |

## Config-as-code surfaces

These surfaces support the `pull` / `push` loop through the shared reconcile
engine (see [Reconcile engine](reconcile.md)):

| Surface | Identity key | Reconcile model |
|---------|-------------|-----------------|
| exclusions | type + OS + value | Create or update |
| blocklist | type + OS + value | Create or update |
| rules | rule name | Create or update |
| firewall | rule name | Create or update |
| devicecontrol | rule name | Create or update |
| network | rule name | Create or update |
| sites | site name | Create or update |
| groups | site ID + group name | Create or update |
| tags | tag key | Create or update |
| locations | location name | Create or update |
| filters | filter name | Create or update |
| tag-rules | rule name | Create or update |
| upgrade-policies | policy name | Create or update |
| cloud-policies | policy ID | Update only (no create) |
| applications rules | rule name | Create or update |

Protection policies (`policies pull/push/diff/revert`) have their own lane
outside the engine: they are scope-singletons, not per-object collections.

All mutations are dry-run by default. Pass `--yes` to apply. The `drift`
command plans every surface above and exits non-zero on divergence.
