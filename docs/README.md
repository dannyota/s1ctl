# s1ctl

Open-source **CLI**, **Go SDK**, and **MCP server** for SentinelOne Singularity
Platform. One tool, three interfaces — operate your SentinelOne environment from
the terminal, from Go code, or from any AI agent that speaks
[Model Context Protocol](https://modelcontextprotocol.io).

> Community project. Not affiliated with or endorsed by SentinelOne, Inc.

## Three ways in

| Interface | What it does | Get started |
|-----------|-------------|-------------|
| **MCP Server** | Give Claude, Cursor, or any MCP client full access to your SentinelOne console — 370+ tools with dynamic loading | [MCP guide](guides/mcp.md) |
| **CLI** | `s1ctl agents list`, `s1ctl threats mitigate` — pull live state, review in `git diff`, push back | [Install](guides/install.md) |
| **Go SDK** | `import "danny.vn/s1/mgmt"` — typed clients for REST, SDL, and GraphQL | [SDK guide](guides/sdk.md) |

## MCP server — quick start

```bash
go install danny.vn/s1/cmd/s1ctl@latest

export S1_CONSOLE_URL=https://your-console.sentinelone.net
export S1_TOKEN=your-api-token

s1ctl mcp install   # writes .mcp.json in the current project
```

Restart your MCP client. The server exposes four meta-tools (`help`, `run`,
`focus`, `unfocus`) and loads group-specific typed tools on demand — staying
within context limits while covering every command.

[Full MCP setup guide &rarr;](guides/mcp.md)

## CLI — quick start

```bash
go install danny.vn/s1/cmd/s1ctl@latest

s1ctl config init     # one-screen wizard
s1ctl doctor          # verify auth + API reach

s1ctl agents list --limit 10
s1ctl threats list --limit 5
s1ctl datalake powerquery --query "endpoint.name contains 'srv'"
```

Mutating commands default to `--dry-run` — nothing changes until you pass
`--yes`.

[Full CLI quickstart &rarr;](guides/quickstart.md)

## API surfaces

| Surface | Protocol | Scope |
|---------|----------|-------|
| **REST MGMT** (v2.1) | REST | Agents, threats, sites, groups, exclusions, policies, remote ops |
| **SDL** | REST + GraphQL | PowerQuery, log ingest/query, file ops |
| **GraphQL** | GraphQL | UAM alerts, xSPM vulnerabilities/misconfigurations, cloud security |

## License

MIT. See [GitHub](https://github.com/dannyota/s1ctl) for source.
