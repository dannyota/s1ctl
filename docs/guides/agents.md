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

All actions are **dry-run by default**: they print what would happen and
change nothing. Pass `--yes` to apply. Most actions take one or more agent
IDs as arguments.

### Isolation and scanning

| Command | Description |
|---------|-------------|
| `agents isolate <id...>` | Network-isolate agents |
| `agents reconnect <id...>` | Reconnect network-isolated agents |
| `agents scan <id>` | Start a full disk scan |
| `agents abort-scan <id>` | Abort a running disk scan |

`isolate` and `reconnect` also accept `--filter key=value` to target agents
by API query (repeatable, and combinable with explicit IDs):

```bash
s1ctl agents isolate 000000                        # dry-run
s1ctl agents isolate 000000 --yes                  # apply
s1ctl agents isolate --filter infected=true --yes  # all infected agents

s1ctl agents reconnect 000000 --yes
s1ctl agents scan 000000 --yes
s1ctl agents abort-scan 000000 --yes
```

### Lifecycle

| Command | Description |
|---------|-------------|
| `agents restart <id>` | Restart the endpoint |
| `agents shutdown <id>` | Shut down the endpoint |
| `agents decommission <id>` | Remove the agent from the console |
| `agents uninstall <id>` | Uninstall the agent |
| `agents approve-uninstall <id>` | Approve a pending uninstall request |
| `agents reject-uninstall <id>` | Reject a pending uninstall request |
| `agents reset-passphrase <id>` | Reset the agent maintenance passphrase |

```bash
s1ctl agents restart 000000 --yes
s1ctl agents shutdown 000000 --yes
s1ctl agents decommission 000000 --yes
s1ctl agents reset-passphrase 000000 --yes
```

> **Warning:** Decommission and uninstall remove the agent from management.
> Neither can be undone from the CLI.

### Passphrases

```bash
s1ctl agents passphrases --site-id 000000
s1ctl agents passphrases --site-id 000000 --json
```

List agent maintenance passphrases for a site.

> **Note:** This command outputs sensitive data (maintenance passphrases).
> The output is printed to stdout for scripting, and a notice is printed to
> stderr reminding the user that the output is sensitive.

### State and configuration

| Command | Description |
|---------|-------------|
| `agents enable <id>` | Enable a disabled agent |
| `agents disable <id>` | Disable an agent |
| `agents reset-config <id>` | Reset agent local configuration |
| `agents mark-up-to-date <id>` | Mark the agent as up to date |
| `agents randomize-uuid <id>` | Randomize the agent UUID |
| `agents ranger <id>` | Enable or disable Ranger network discovery |
| `agents local-upgrade <id>` | Authorize or revoke local upgrade/downgrade |
| `agents local-upgrade-status <id>` | Show local upgrade/downgrade authorization status |

```bash
s1ctl agents enable 000000 --yes
s1ctl agents disable 000000 --yes
s1ctl agents reset-config 000000 --yes
```

`ranger` uses `--state on|off`:

```bash
s1ctl agents ranger 000000 --state on --yes
s1ctl agents ranger 000000 --state off --yes
```

`local-upgrade` uses `--authorize` or `--revoke`:

```bash
s1ctl agents local-upgrade 000000 --authorize --yes
s1ctl agents local-upgrade 000000 --revoke --yes
```

`local-upgrade-status` is a read command (no `--yes` needed):

```bash
s1ctl agents local-upgrade-status 000000
s1ctl agents local-upgrade-status 000000 --json
```

### Organization

Move an agent between groups or sites, or set its external ID:

| Command | Required flag | Description |
|---------|---------------|-------------|
| `agents move <id>` | `--group-id` | Move to a different group |
| `agents move-to-site <id>` | `--site-id` | Move to a different site |
| `agents set-external-id <id>` | `--external-id` | Set the external ID |

```bash
s1ctl agents move 000000 --group-id 000000 --yes
s1ctl agents move-to-site 000000 --site-id 000000 --yes
s1ctl agents set-external-id 000000 --external-id my-asset-tag --yes
```

### Data collection and broadcast

| Command | Description |
|---------|-------------|
| `agents broadcast <id...>` | Display a message on the agent's endpoint |
| `agents fetch-files <id>` | Fetch specific files from an agent to the console |
| `agents fetch-installed-apps <id>` | Fetch installed-applications inventory |
| `agents fetch-firewall-rules <id>` | Fetch current firewall-rules inventory |

`broadcast` requires `--message`:

```bash
s1ctl agents broadcast 000000 --message "Scheduled maintenance at 22:00 UTC" --yes
```

`fetch-files` accepts `--password` to encrypt the fetched archive. Using
`--password` is recommended for security:

```bash
s1ctl agents fetch-files 000000 --yes
s1ctl agents fetch-files 000000 --password "s3cret" --yes
```

`fetch-installed-apps` and `fetch-firewall-rules` trigger an inventory
refresh from the agent:

```bash
s1ctl agents fetch-installed-apps 000000 --yes
s1ctl agents fetch-firewall-rules 000000 --yes
```

### Firewall logging

Toggle firewall logging on an agent with `--state on|off`:

```bash
s1ctl agents firewall-logging 000000 --state on --yes
s1ctl agents firewall-logging 000000 --state off --yes
```

### Upgrade

Trigger an agent software upgrade. Identify the package with exactly one of
`--package-id`, `--file-name` (which also needs `--os-type`), or `--path`.
Target agents by ID, or by `--site-id` / `--group-id` / `--query` filter:

```bash
s1ctl agents upgrade 000000 --package-id 000000 --yes
s1ctl agents upgrade --group-id 000000 --package-id 000000 --yes
s1ctl agents upgrade 000000 --file-name AgentSetup.exe --os-type windows --yes
```

| Flag | Description |
|------|-------------|
| `--package-id` | Upgrade package ID |
| `--file-name` | Package file name (requires `--os-type`) |
| `--path` | Local path to the package on the endpoint |
| `--os-type` | Target OS (`linux`, `macos`, `windows`) |
| `--package-type` | `Agent`, `Ranger`, or `AgentAndRanger` |
| `--allow-downgrade` | Allow downgrading the agent version |
| `--scheduled` | Upgrade per the agent upgrade schedule |

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

### Retrieve maintenance passphrases

```bash
s1ctl agents passphrases --site-id 000000 --json
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
