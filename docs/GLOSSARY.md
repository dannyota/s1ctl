# Glossary

Plain-language definitions of the terms used across these docs and the CLI.
Start at the [quickstart](guides/quickstart.md) for everyday work — this page
is here when a doc or `--help` string uses a word you don't recognize.

## Core concepts

| Term | What it means for you |
|------|----------------------|
| **console** | Your SentinelOne management web UI. Every API call goes through it. Set `S1_CONSOLE_URL` to point at yours (e.g. `https://your-console.sentinelone.net`). |
| **site** | A top-level tenant partition in the console. Most objects belong to a site. Use `--site-id` to scope commands. |
| **group** | A subdivision of a site. Agents belong to a group; policies can be set per-group. |
| **account** | The highest organizational level, above sites. Multi-site deployments share one account. |
| **agent** | The SentinelOne endpoint agent installed on a machine. `s1ctl agents list` shows them. |
| **threat** | A malicious or suspicious event detected by an agent. `s1ctl threats list` shows them. |
| **alert** | A detection raised for triage — from behavioral AI, STAR rules, or threat intelligence. `s1ctl alerts list` queries them via GraphQL (UAM). |
| **exclusion** | A path, hash, certificate, or browser rule that tells agents to skip scanning. The first surface with full config-as-code support (pull/push). |
| **policy** | The agent configuration applied at the site, group, or account level — scan modes, engine toggles, network controls. |

## API surfaces

| Term | What it means for you |
|------|----------------------|
| **REST MGMT API** | The main REST API (v2.1). Covers agents, threats, sites, groups, exclusions, policies, remote ops. Auth: `ApiToken <token>`. |
| **SDL (Singularity Data Lake)** | SentinelOne's centralized event store. Ingests endpoint, cloud, and identity data. Queried via powerQuery (REST or GraphQL). Auth: `Bearer <token>`. |
| **GraphQL API** | Newer API for alerts (UAM), xSPM findings, and cloud security. Auth: `Bearer <token>`. s1ctl defaults to GraphQL where it offers richer filtering. |
| **protocol** | Which API protocol a command uses underneath. When a surface is available via multiple protocols, `--protocol rest\|graphql\|sdl` overrides the default. |
| **`S1_TOKEN`** | Your API token. Authenticates all three surfaces. Resolved from env vars or config file — never hardcoded. |
| **`S1_CONSOLE_URL`** | Your console base URL (e.g. `https://your-console.sentinelone.net`). Set in `.env` or config. |

## CLI concepts

| Term | What it means for you |
|------|----------------------|
| **the loop** | Core mental model: pull live state &rarr; review in `git diff` &rarr; push back. Git history is the source of truth. See [the loop](guides/the-loop.md). |
| **pull** | Read-only. Downloads live config to local files. Never changes anything on the console. |
| **push** | Deploys local file state to the live console. A mutation — defaults to dry-run. |
| **dry run** | A preview of what a push or action would do, without doing it. The default for every mutating command. |
| **`--yes`** | The flag that makes a mutation real. Without it, every mutation is a dry-run preview. |
| **`--json`** | Machine-readable output flag. Every read command supports it. Pipe to `jq` for filtering. |
| **`--site-id`** | Scopes a command to a specific site. Falls back to the value in your config file. |
| **`--limit`** | Controls pagination — how many items a list command returns per page. |
| **doctor** | `s1ctl doctor` — connectivity check. Validates config, token, and console reachability. |
| **config** | `s1ctl config` — manage the local config file (`~/.s1ctl/config.yaml`). |

## Data lake

| Term | What it means for you |
|------|----------------------|
| **powerQuery** | SentinelOne's query language for the data lake. `s1ctl datalake powerquery` runs one. Available via both GraphQL and REST — GraphQL is the default. |
| **deep visibility** | Legacy name for endpoint telemetry in the data lake. Same data, accessed via powerQuery. |

## Security operations

| Term | What it means for you |
|------|----------------------|
| **isolation** | Disconnects an agent from the network (keeps the console tunnel). `s1ctl agents isolate 000000 --yes`. Reversed with `reconnect`. |
| **mitigation** | An action taken on a threat. Four actions: **kill** (terminate the process), **quarantine** (isolate the malicious file), **remediate** (undo changes made by the threat), **rollback** (restore to a pre-threat VSS snapshot). |
| **verdict** | The analyst's classification: true positive, false positive, suspicious, or undefined. Applies to threats and alerts. |
| **incident status** | The workflow state of a threat: unresolved, in-progress, or resolved. |
| **scan** | A full disk scan triggered on an agent. `s1ctl agents scan 000000 --yes`. |
| **decommission** | Permanently removes an agent from the console. Irreversible — use with care. |

## Cloud and vulnerability management

| Term | What it means for you |
|------|----------------------|
| **xSPM** | Extended Security Posture Management. Covers vulnerabilities and misconfigurations across endpoints and cloud workloads. |
| **vulnerability** | A known CVE found on an endpoint by the SentinelOne agent. `s1ctl vulnerabilities list` queries via GraphQL. |
| **misconfiguration** | A security posture finding — weak OS setting, missing patch, insecure config. `s1ctl misconfigurations list` queries via GraphQL. |
| **cloud policy** | A security policy governing cloud accounts (AWS, Azure, GCP). `s1ctl cloud-policies list` queries via GraphQL. |

## See also

- [Quickstart](guides/quickstart.md) — common workflows
- [The loop](guides/the-loop.md) — core mental model
- [Architecture](design/architecture.md) — how s1ctl is built
- [Catalog](design/catalog.md) — per-surface implementation status
