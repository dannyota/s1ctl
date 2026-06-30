# Configure

s1ctl needs two things: your console URL and an API token.

## Quick start

```bash
s1ctl config
```

The interactive wizard walks through setup and writes `~/.s1ctl/config.yaml`.

## Environment variables

| Variable | Description |
|----------|-------------|
| `S1_CONSOLE_URL` | Console URL (e.g. `https://your-console.sentinelone.net`) |
| `S1_TOKEN` | API token |

Environment variables override the config file.

## Config file

Default location: `~/.s1ctl/config.yaml`.

```yaml
console_url: https://your-console.sentinelone.net
token: your-api-token
```

Override the path with `--config`:

```bash
s1ctl --config /path/to/config.yaml doctor
```

## Resolution order

Highest priority first:

1. `S1_*` environment variables
2. `--config` flag
3. `~/.s1ctl/config.yaml`
4. `./config/config.yaml`

## Verify

```bash
s1ctl doctor
```

Checks connectivity to all three API surfaces (REST MGMT, SDL, GraphQL).
