# mcp

Run Model Context Protocol server

## mcp install

Register s1ctl in the project .mcp.json

```text
s1ctl mcp install
```

Add s1ctl as an MCP server in the project-level .mcp.json so every
Claude Code session in this directory gets s1ctl tools automatically.
Idempotent — updates the entry if it already exists.

## mcp serve

Start the MCP server on stdio

```text
s1ctl mcp serve [flags]
```

Start a Model Context Protocol (MCP) server that exposes every s1ctl
command as an MCP tool and every docs guide as an MCP resource.

Tools are auto-generated from the command tree — adding a command
automatically creates a tool. Resources are embedded from docs/guides/.

Use --read-only to hide mutation tools and block mutations via run.

Configure Claude Code to use this server:

  s1ctl mcp install

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--read-only` | bool | false | expose only read-only tools and block mutations |
