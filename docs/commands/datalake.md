# datalake

Query Singularity Data Lake (SDL)

## datalake facet

Aggregate the most common values of a field (SDL REST)

```text
s1ctl datalake facet [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--end` | string | - | end time |
| `--field` | string | - | field to aggregate (required) |
| `--filter` | string | - | query filter expression |
| `--max-count` | int | 0 | max distinct values to return |
| `--start` | string | - | start time, e.g. 24h or timestamp (required) |

## datalake files

Manage data lake configuration files

```text
s1ctl datalake files
```

## datalake ingest

Ingest events or raw logs into the data lake

```text
s1ctl datalake ingest
```

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

By default, returns only the first page. Use --all to fetch all pages via
continuation token, or --max-events to cap the total number of events.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--all` | bool | false | fetch all pages via continuation token |
| `--end` | string | - | end time |
| `--max-count` | int | 0 | max events per page (1-5000) |
| `--max-events` | int | 0 | max total events across all pages (0 = no limit) |
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

## datalake timeseries

Run a time-series aggregation (SDL REST)

```text
s1ctl datalake timeseries [flags]
```

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--buckets` | int | 0 | number of time buckets |
| `--end` | string | - | end time |
| `--filter` | string | - | query filter expression (required) |
| `--function` | string | - | aggregation function (e.g. count, mean(field)) |
| `--start` | string | - | start time, e.g. 24h or timestamp (required) |
