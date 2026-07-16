# MCP workflows

Agent-oriented cheat sheet for working with s1ctl as an MCP server.

## Discovery flow

1. `help` -- list command groups with read/mutation counts
2. `help {group}` -- list subcommands in a group
3. `usage {command}` -- flags, args, and full schema for one command
4. `focus {group}` -- load typed tool schemas (structured calls)
5. `unfocus {group}` -- unload when done to free context

## Running commands

- `run {command}` -- execute any command as a string
- Typed tools (after `focus`) -- structured parameters, no quoting needed

## Read vs mutation semantics

- Read commands execute immediately and return JSON
- Mutation commands are dry-run by default; pass `--yes` to apply
- Tool descriptions include `[mutation]` for mutation commands
- Tool annotations: `readOnlyHint` and `destructiveHint` classify each tool

## Read-only mode

When the server starts with `--read-only`:

- Mutation tools are hidden from `tools/list`
- Focused groups only load read-only tools
- The `run` meta-tool sets `S1_READONLY=1` in subprocesses
- Even if a mutation command is invoked via `run`, the guard blocks it

## JSON output contract

All tool output is JSON. Successful calls return the command's JSON output
directly. The `--json` and `--no-progress` flags are added automatically.

## Error envelope

When a tool call fails, the content text is a JSON envelope:

```json
{"error":{"message":"<error details>"}}
```

Parse `error.message` for the cause. Partial output (e.g. API error bodies)
is included in the message.

## Large output and spill files

Output over 4 MiB is written to a temporary file. The tool returns a JSON
pointer instead:

```json
{
  "file": "/tmp/s1ctl-mcp/s1ctl-mcp-123.json",
  "bytes": 5242880,
  "preview": "<first 2 KiB of output>",
  "message": "Output exceeded 4 MiB limit. ..."
}
```

- `preview` contains the first 2 KiB (rune-safe) for quick inspection
- Read the `file` path for full results
- Spill files are removed after 24 hours
- Use `--max-results` or narrower filters to reduce output

## Tips

- Always scope queries with `--site-id`
- Use `focus` for repeated structured calls within a group
- Use `run` for quick one-off commands
- Prefer `focus` for commands with filter expressions
- In `run`, use shell-style quoting: `--filter 'event.type = "Login"'`
