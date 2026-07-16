# MCP server

Use s1ctl as a [Model Context Protocol](https://modelcontextprotocol.io) server
to give AI agents — Claude Code, Claude Desktop, Cursor, Windsurf, or any MCP
client — direct access to your SentinelOne console.

The server uses **dynamic tool loading**: instead of exposing all 370+ commands
upfront, it starts with 5 meta-tools and loads group-specific typed tools on
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

Or run `s1ctl config init` to write `~/.s1ctl/config.yaml`. See
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

- **Meta-tools** — `help`, `run`, `usage`, `focus`, `unfocus` (always loaded)
- **Group tools** — loaded on demand via `focus` (e.g. `focus group="agents"`)
- **Resources** — one per guide (`guide://{name}`)

### Workflow

1. Call `help` to discover available command groups
2. Call `help group="agents"` to list subcommands in a group
3. Call `usage command="agents list"` to see one command's flags and args
4. Call `focus group="agents"` to load typed tools (`agents_list`, `agents_get`, etc.)
5. Use the typed tools with full parameter schemas
6. Call `unfocus group="agents"` when done to free context
7. Use `run` anytime for quick one-off commands without focusing

### Example conversation

> **You:** Which Windows agents are online right now?
>
> The agent calls `run command="agents list --os-type windows --active"`
> and summarizes the results.

> **You:** Show me unresolved threats from the last 24 hours.
>
> The agent calls `focus group="threats"`, then `threats_list` with the
> appropriate filters.

All tool output is JSON. Mutations are dry-run by default — the agent must
pass `--yes` to apply, same as the CLI.

### Tool annotations

Every tool carries MCP `annotations` that classify it for safety:

- **Read-only commands**: `readOnlyHint: true`, `destructiveHint: false`
- **Mutation commands**: `readOnlyHint: false`, `destructiveHint: true`
- **Meta-tools** (help, usage, focus, unfocus): `readOnlyHint: true`
- **run**: no annotations in normal mode (it can invoke either kind)

Clients that support annotations can use these to gate or warn before
executing destructive tools.

### Structured errors

When a tool call fails, the content text is a JSON envelope:

```json
{"error":{"message":"<error details>"}}
```

Parse `error.message` for the cause. Partial output (e.g. API error bodies)
is included in the message string.

### Large output and concurrency

Each tool call runs in its own subprocess, so calls execute concurrently and
one slow query never blocks the rest. A single result is capped at 4 MiB:
larger output is written to a temporary file and the tool returns a JSON
pointer with `file`, `bytes`, `preview`, and `message`. The `preview` field
contains the first 2 KiB (rune-safe) for quick inspection without reading the
full file. Spill files are removed after 24 hours.

### Read-only mode

Start the server with `--read-only` to restrict it to read-only operations:

```bash
s1ctl mcp serve --read-only
```

Or in your MCP client config:

```json
{
  "mcpServers": {
    "s1ctl": {
      "command": "s1ctl",
      "args": ["mcp", "serve", "--read-only"]
    }
  }
}
```

When read-only mode is active:

- Mutation tools are hidden from `tools/list`
- Focused groups load only read-only tools
- The `run` meta-tool sets `S1_READONLY=1` so mutations are blocked
- Server instructions note the mode

This is useful for monitoring agents that should observe but never modify.

## Security

- The MCP server inherits the permissions of your API token. Use a
  least-privilege token scoped to what the agent needs.
- All mutations require `--yes` — the agent cannot accidentally modify your
  environment without explicit confirmation.
- Use `--read-only` for agents that should only observe.
- The server runs locally on stdio. No network listener is opened.
