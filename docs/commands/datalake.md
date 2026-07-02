# datalake

Query Singularity Data Lake (SDL)

## datalake dashboards

Manage Data Lake dashboards

```text
s1ctl datalake dashboards
```

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

## datalake numeric

Run a numeric aggregation query (SDL REST)

```text
s1ctl datalake numeric [flags]
```

Run a numeric aggregation query against the Singularity Data Lake.

Counts events, computes event rate, or applies an aggregation function
(e.g. mean, min, max) to a numeric field across one or more time buckets.

Note: numericQuery is effectively deprecated in favour of timeseries with
createSummaries=false, but remains useful for sub-30-second bucket
granularity and users with limited query permissions.

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--buckets` | int | 0 | number of buckets (1-5000) |
| `--end` | string | - | end time |
| `--filter` | string | - | query filter expression |
| `--function` | string | - | aggregation function (e.g. rate, count, mean(field)) |
| `--priority` | string | - | query priority (low, high) |
| `--start` | string | - | start time, e.g. 1h or timestamp (required) |

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

Manage saved PowerQueries

```text
s1ctl datalake saved-queries
```

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
