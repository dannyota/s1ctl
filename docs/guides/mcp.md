# MCP server

Use s1ctl as a [Model Context Protocol](https://modelcontextprotocol.io) server
to give AI agents — Claude Code, Claude Desktop, Cursor, Windsurf, or any MCP
client — direct access to your SentinelOne console.

The server uses **dynamic tool loading**: instead of exposing all 370+ commands
upfront, it starts with 4 meta-tools and loads group-specific typed tools on
demand. This keeps the agent's context window small while covering every
command.

## Install

```bash
go install danny.vn/s1/cmd/s1ctl@latest
```

Or download a pre-built binary from the
[releases page](https://github.com/dannyota/s1ctl/releases).

## Configure credentials

The MCP server reads the same config as the CLI. Set two environment variables:

```bash
export S1_CONSOLE_URL=https://your-console.sentinelone.net
export S1_TOKEN=your-api-token
```

Or run `s1ctl config` to write `~/.s1ctl/config.yaml`. See
[Configure](guides/configure.md) for details.

Verify connectivity:

```bash
s1ctl doctor
```

## Register with your MCP client

### Claude Code

From your project directory:

```bash
s1ctl mcp install
```

This writes an entry to `.mcp.json` in the current directory. Restart Claude
Code to pick up the new server.

### Claude Desktop

Add to your `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "s1ctl": {
      "command": "s1ctl",
      "args": ["mcp", "serve"],
      "env": {
        "S1_CONSOLE_URL": "https://your-console.sentinelone.net",
        "S1_TOKEN": "your-api-token"
      }
    }
  }
}
```

### Cursor / other MCP clients

Point your client at the stdio command:

```text
s1ctl mcp serve
```

Pass `S1_CONSOLE_URL` and `S1_TOKEN` as environment variables, or ensure
`~/.s1ctl/config.yaml` exists.

## How it works

The server exposes:

- **Meta-tools** — `help`, `run`, `focus`, `unfocus` (always loaded)
- **Group tools** — loaded on demand via `focus` (e.g. `focus group="agents"`)
- **Resources** — one per guide (`guide://{name}`)

### Workflow

1. Call `help` to discover available command groups
2. Call `help group="agents"` to list subcommands in a group
3. Call `focus group="agents"` to load typed tools (`agents_list`, `agents_get`, etc.)
4. Use the typed tools with full parameter schemas
5. Call `unfocus group="agents"` when done to free context
6. Use `run` anytime for quick one-off commands without focusing

### Example conversation

> **You:** How many Windows agents are online?
>
> The agent calls `run command="agents count --query windows --status active"`
> and returns the count.

> **You:** Show me unresolved threats from the last 24 hours.
>
> The agent calls `focus group="threats"`, then `threats_list` with the
> appropriate filters.

All tool output is JSON. Mutations are dry-run by default — the agent must
pass `--yes` to apply, same as the CLI.

## Security

- The MCP server inherits the permissions of your API token. Use a
  least-privilege token scoped to what the agent needs.
- All mutations require `--yes` — the agent cannot accidentally modify your
  environment without explicit confirmation.
- The server runs locally on stdio. No network listener is opened.
