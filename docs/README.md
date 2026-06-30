# s1ctl

Operate **SentinelOne Singularity Platform** as code — one Go CLI and one
importable Go SDK covering the REST Management API, Singularity Data Lake, and
GraphQL surfaces. The core loop is **pull live state, review in `git diff`, push
back**, with git history as the source of truth.

> Mutating commands default to `--dry-run` and print a banner — nothing changes
> until you pass `--yes`. Always dry-run, read it, then apply.

New here? **[Install](guides/install.md) &rarr;
[Configure](guides/configure.md).**
Building it? **[Architecture](design/architecture.md).**

## API surfaces

| Surface | Protocol | Scope |
|---------|----------|-------|
| **REST MGMT** (v2.1) | REST | Agents, threats, sites, groups, exclusions, policies, remote ops |
| **SDL** | REST | PowerQuery, log ingest/query, file ops |
| **GraphQL** | GraphQL | UAM alerts, xSPM vulnerabilities/misconfigurations, cloud security |

## Quick start

```bash
go install danny.vn/s1/cmd/s1ctl@latest

s1ctl config          # one-screen wizard
s1ctl doctor          # verify auth + API reach
```
