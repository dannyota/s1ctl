# Configure

s1ctl needs two things: your console URL and an API token.

## Quick start

```bash
s1ctl config init
```

The interactive wizard walks through setup and writes `~/.s1ctl/config.yaml`.

## Environment variables

| Variable | Description |
|----------|-------------|
| `S1_CONSOLE_URL` | Console URL (e.g. `https://your-console.sentinelone.net`) |
| `S1_TOKEN` | API token |
| `S1_SDL_URL` | SDL REST console URL (e.g. `https://xdr.us1.sentinelone.net`) — only needed with `--protocol rest` |

Environment variables override the config file.

> **Note:** Data Lake queries use GraphQL by default, which connects through
> your management console — no extra URL needed. If you prefer the REST
> protocol (`--protocol rest`), set `S1_SDL_URL` to your Data Lake console URL.

## Config file

Default location: `~/.s1ctl/config.yaml`.

```yaml
console_url: https://your-console.sentinelone.net
token: your-api-token
sdl_url: https://xdr.us1.sentinelone.net  # optional, only for --protocol rest
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

## FAQ

### How do I create an API token?

In the SentinelOne console: Settings > Users > API Token > Generate.
Grant it the minimum scope needed for your use case.

### Can I use multiple consoles?

Yes. Use `--config` to point to a different config file per console, or
switch `S1_CONSOLE_URL` and `S1_TOKEN` environment variables.

### What is `S1_SDL_URL` and when do I need it?

Only needed if you use `--protocol rest` for data lake queries. GraphQL
(the default) routes through your management console URL.
