# mcp

Model Context Protocol server for AI agent integration.

The server uses **dynamic tool loading** to stay within context limits.
Instead of exposing all 300+ commands upfront, it starts with 4 meta-tools
and loads group-specific tools on demand via MCP `listChanged`.

## mcp serve

Start the MCP server on stdio.

```text
s1ctl mcp serve
```

The server exposes:

- **Meta-tools** — `run`, `help`, `focus`, `unfocus` (always loaded)
- **Group tools** — loaded on demand via `focus` (e.g. `focus group="agents"`)
- **Resources** — one per `docs/guides/*.md` file (`guide://{name}`)

### Meta-tools

| Tool | Purpose |
|------|---------|
| `run` | Run any s1ctl command by string (universal, always works) |
| `help` | List groups or subcommands within a group |
| `focus` | Load typed tools for a group (triggers `listChanged`) |
| `unfocus` | Unload a group's tools to free context |

### Workflow

1. Call `help` to discover available command groups
2. Call `focus group="agents"` to load typed tools for agents
3. Use the loaded `agents_list`, `agents_get`, etc. with full schemas
4. Call `unfocus group="agents"` when done to free context
5. Use `run` anytime for quick one-off commands without focusing

All tool output is JSON. Mutations are dry-run by default (the agent must
pass `--yes` to apply, same as the CLI).

## mcp install

Register s1ctl in the project `.mcp.json`.

```text
s1ctl mcp install
```

Writes to `.mcp.json` in the current directory, merging with any existing
entries. Idempotent — updates the entry if it already exists. Restart
Claude Code to pick up the new server.
