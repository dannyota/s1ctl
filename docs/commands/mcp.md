# mcp

Model Context Protocol server for AI agent integration.

Every CLI command is automatically exposed as an MCP tool, and every
guide page as an MCP resource — zero maintenance when commands are added.

## mcp serve

Start the MCP server on stdio.

```text
s1ctl mcp serve
```

The server exposes:

- **Tools** — one per CLI leaf command (auto-generated from the command tree)
- **Resources** — one per `docs/guides/*.md` file (`guide://{name}`)

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
