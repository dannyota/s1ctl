# Agents

Query, inspect, and act on endpoint agents.

## List agents

```bash
s1ctl agents list
s1ctl agents list --query "web-server" --os-type linux --limit 20
s1ctl agents list --site-id 000000 --sort-by lastActiveDate --sort-order desc
```

### Flags

| Flag | Type | Description |
|------|------|-------------|
| `--query` | string | Free text search across agent fields |
| `--site-id` | string[] | Filter by site ID (repeatable) |
| `--group-id` | string[] | Filter by group ID (repeatable) |
| `--os-type` | string[] | Filter by OS type (repeatable) |
| `--sort-by` | string | Sort field (e.g. `computerName`, `lastActiveDate`) |
| `--sort-order` | string | Sort direction (`asc`, `desc`) |
| `--limit` | int | Max results per page (default 50) |
| `--all` | bool | Fetch all pages automatically |
| `--cursor` | string | Pagination cursor for manual paging |

The `--output` and `--json` flags are global and work on all read commands.

## Get agent details

```bash
s1ctl agents get 000000
s1ctl agents get 000000 --json
```

Returns a detail view: ID, name, OS, version, network status, infection
state, site, group, external IP, last active date, and registration date.

## Count agents

```bash
s1ctl agents count
s1ctl agents count --site-id 000000
```

Returns the total agent count. Accepts `--site-id` to scope the count to
one or more sites.

## Actions

All actions are **dry-run by default**. Pass `--yes` to apply.

| Command | Description |
|---------|-------------|
| `agents isolate <id>` | Network-isolate an agent |
| `agents connect <id>` | Reconnect an isolated agent |
| `agents scan <id>` | Start a full disk scan |
| `agents decommission <id>` | Decommission an agent |

### Isolate and reconnect

```bash
s1ctl agents isolate 000000            # dry-run: prints what would happen
s1ctl agents isolate 000000 --yes      # applies the isolation

s1ctl agents connect 000000 --yes      # reconnect the agent
```

### Scan

```bash
s1ctl agents scan 000000 --yes
```

### Decommission

```bash
s1ctl agents decommission 000000       # dry-run
s1ctl agents decommission 000000 --yes
```

> **Warning:** Decommission removes the agent from the console. This cannot
> be undone from the CLI.

## Workflows

### Find inactive agents

List agents sorted by last active date, oldest first:

```bash
s1ctl agents list --sort-by lastActiveDate --sort-order asc --limit 25
```

### Export all agents as CSV

```bash
s1ctl agents list --all --output csv > agents.csv
```

### Count agents per site

```bash
s1ctl agents count --site-id 000000
s1ctl agents count --site-id 111111
```

Or count all agents across the account:

```bash
s1ctl agents count
```

### Filter by OS and pipe to jq

```bash
s1ctl agents list --os-type linux --json | jq '.[].computerName'
```

### Isolate a compromised agent

```bash
s1ctl agents get 000000                    # confirm the agent
s1ctl agents isolate 000000                # dry-run
s1ctl agents isolate 000000 --yes          # apply
```

## Output formats

| Flag | Format | Use case |
|------|--------|----------|
| (default) | table | Human-readable terminal output |
| `--json` | JSON | Pipe to jq, scripts, automation |
| `--output csv` | CSV | Spreadsheets, bulk analysis |
