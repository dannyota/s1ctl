# Remote operations

Manage remote scripts: list, inspect, execute on agents, review results, and
handle pending approvals and guardrails.

> Prerequisites: `s1ctl` installed and configured (`S1_CONSOLE_URL`, `S1_TOKEN`).
> Remote ops requires the appropriate license and RBAC permissions.

## List scripts

```bash
s1ctl remoteops list
s1ctl remoteops list --site-id 000000
s1ctl remoteops list --query "Collect Logs" --all --json
```

| Flag | Description |
|------|-------------|
| `--site-id` | Filter by site ID (repeatable) |
| `--query` | Free text search |
| `--limit` | Max results per page (default 50) |
| `--all` | Fetch all pages |
| `--cursor` | Pagination cursor |

## Get script details

```bash
s1ctl remoteops get 000000
s1ctl remoteops get 000000 --json
```

## View script content

Print the script body to stdout or save to a file.

```bash
s1ctl remoteops content 000000
s1ctl remoteops content 000000 --out script.ps1
```

| Flag | Description |
|------|-------------|
| `--out` | Write content to file (default: stdout) |

## Run a script

Execute a remote script from the Script Library on targeted agents. Requires
at least one targeting flag. Dry-run by default.

```bash
# Target specific agents
s1ctl remoteops run 000000 --agent-id 000001 --yes

# Target a site
s1ctl remoteops run 000000 --site-id 000000 --yes

# Target a group
s1ctl remoteops run 000000 --group-id 000000 --yes

# With parameters and custom timeout
s1ctl remoteops run 000000 \
  --agent-id 000001 \
  --input-params "-Path C:\Logs" \
  --timeout 3600 \
  --output-dest SentinelCloud \
  --yes
```

| Flag | Description |
|------|-------------|
| `--agent-id` | Target agent IDs (repeatable) |
| `--site-id` | Target site IDs (repeatable) |
| `--group-id` | Target group IDs (repeatable) |
| `--input-params` | Script input parameters |
| `--timeout` | Runtime timeout in seconds (60--172800) |
| `--output-dest` | Output destination: `SentinelCloud`, `Local`, `None`, `SingularityXDR` (default `SentinelCloud`) |
| `--description` | Task description |
| `--yes` | Apply (default: dry-run) |

## View execution results

```bash
s1ctl remoteops results 000000
s1ctl remoteops results 000000 --status completed
s1ctl remoteops results 000000 --all --json
```

| Flag | Description |
|------|-------------|
| `--status` | Filter by status: `created`, `pending`, `in_progress`, `completed`, `failed`, `canceled`, `expired` (repeatable) |
| `--limit` | Max results per page (default 50) |
| `--all` | Fetch all pages |
| `--cursor` | Pagination cursor |

## Update script metadata

Update a script's properties (name, type, OS types, timeout) from a JSON file.
Dry-run by default.

```bash
s1ctl remoteops update 000000 --from-file script-meta.json --yes
```

Example `script-meta.json`:

```json
{
  "scriptName": "Collect Logs",
  "scriptType": "dataCollection",
  "osTypes": ["linux", "macos"],
  "inputRequired": false,
  "inputExample": "-",
  "inputInstructions": "-",
  "scriptRuntimeTimeoutSeconds": 3600
}
```

| Flag | Description |
|------|-------------|
| `--from-file` | JSON file with the update data (required) |
| `--yes` | Apply (default: dry-run) |

## Upload limits

Check package upload size limits.

```bash
s1ctl remoteops upload-limits
```

## Pending approvals

When guardrails are configured, script executions on large numbers of
endpoints require approval before running.

### List pending executions

```bash
s1ctl remoteops pending list
s1ctl remoteops pending list --site-id 000000 --all --json
```

| Flag | Description |
|------|-------------|
| `--site-id` | Filter by site ID (repeatable) |
| `--account-id` | Filter by account ID (repeatable) |
| `--group-id` | Filter by group ID (repeatable) |
| `--sort-by` | Sort field: `id`, `createdAt`, `state` |
| `--sort-order` | Sort direction (`asc`, `desc`) |
| `--limit` | Max results per page (default 50) |
| `--all` | Fetch all pages |
| `--cursor` | Pagination cursor |

### Approve or decline

Dry-run by default.

```bash
s1ctl remoteops pending approve 000000 --yes
s1ctl remoteops pending decline 000000 --yes
```

## Guardrails

Guardrails require approval before scripts run on large numbers of endpoints.
Each guardrail is configured per scope (account, site, or group).

### Get guardrail configuration

```bash
s1ctl remoteops guardrails get --scope-id 000000 --scope-level site
```

| Flag | Description |
|------|-------------|
| `--scope-id` | Scope ID (required) |
| `--scope-level` | Scope level: `account`, `site`, or `group` (required) |

### Set a guardrail

Create or update a guardrail from a JSON file. Dry-run by default.

```bash
s1ctl remoteops guardrails set --from-file guardrail.json --yes
```

Example `guardrail.json`:

```json
{
  "scopeId": "000000000000000000",
  "scopeLevel": "site",
  "endpointsQuantity": 100,
  "scriptTypes": ["action"],
  "enabled": true
}
```

### Check a guardrail

Read-only pre-check: would running a script on given agents trip a guardrail?

```bash
s1ctl remoteops guardrails check --from-file check.json
```

Example `check.json`:

```json
{
  "scriptId": "000000000000000001",
  "agentIds": ["000000000000000002"]
}
```

### Delete a guardrail

```bash
s1ctl remoteops guardrails delete \
  --scope-id 000000 --scope-level site --yes
```

## Workflows

### Run a script and monitor results

```bash
# 1. Find the script
s1ctl remoteops list --query "Collect Logs"

# 2. Dry-run to preview
s1ctl remoteops run 000000 --agent-id 000001

# 3. Execute
s1ctl remoteops run 000000 --agent-id 000001 --yes

# 4. Check results (use the parent task ID from step 3)
s1ctl remoteops results 000000
```

### Review and approve pending executions

```bash
s1ctl remoteops pending list --site-id 000000
s1ctl remoteops pending approve 000000 --yes
```

### Save a script locally for review

```bash
s1ctl remoteops content 000000 --out review-script.ps1
```

## See also

- [Agents](agents.md) -- agent management and targeting
- [`remoteops` command reference](../commands/remoteops.md)
