# datalake

Query Singularity Data Lake (SDL)

## datalake powerquery

Execute a PowerQuery

```text
s1ctl datalake powerquery [flags]
```

Execute a PowerQuery against the Singularity Data Lake.

By default, uses the GraphQL protocol which connects through the management
console and does not require a separate SDL URL. Use --protocol rest to use
the REST API, which requires S1_SDL_URL to be configured.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--col-width` | int | 120 | max column width in table output |
| `--end` | string | - | end time |
| `--priority` | string | low | query priority (low, high) [REST only] |
| `--protocol` | string | graphql | API protocol (graphql, rest) |
| `--query` | string | - | PowerQuery expression (required) |
| `--start` | string | 24h | start time (e.g. 24h, 7d) |

## datalake query

Execute a basic log query

```text
s1ctl datalake query [flags]
```

Execute a basic log query against the Singularity Data Lake.

Uses the SDL REST API (/api/query) to search log events. Requires
S1_SDL_URL to be configured.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--all` | bool | false | fetch all pages via continuation token |
| `--end` | string | - | end time |
| `--max-count` | int | 0 | max events to return |
| `--protocol` | string | rest | API protocol (rest) |
| `--query` | string | - | query expression (required) |
| `--start` | string | 24h | start time (e.g. 24h, 7d) |

## datalake saved-queries

List saved PowerQueries

```text
s1ctl datalake saved-queries
```

List saved searches from the Singularity Data Lake console.
Shows both private and shared saved queries.
