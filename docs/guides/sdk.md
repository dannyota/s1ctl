# Go SDK

The SDK packages are independently importable Go libraries. Each covers one
SentinelOne API protocol — use one, two, or all three.

```text
danny.vn/s1/mgmt      REST MGMT v2.1 — agents, threats, sites, policies, …
danny.vn/s1/sdl       Singularity Data Lake — queries, ingest, files
danny.vn/s1/graphql   GraphQL — alerts, vulnerabilities, misconfigurations, cloud
```

## Client creation

All three packages use the same constructor pattern:

```go
client := mgmt.NewClient("https://your-console.sentinelone.net", token)
client := sdl.NewClient("https://your-console.sentinelone.net", token)
client := graphql.NewClient("https://your-console.sentinelone.net", token)
```

Optional `WithHTTPClient` overrides the default `http.Client`:

```go
client := mgmt.NewClient(url, token, mgmt.WithHTTPClient(customHTTP))
```

## Error handling

Every package returns typed errors for non-2xx responses.

| Package | Error type | Fields |
|---------|-----------|--------|
| `mgmt` | `*mgmt.APIError` | `Status`, `Title`, `Detail` |
| `sdl` | `*sdl.APIError` | `Status`, `Body` |
| `graphql` | `*graphql.HTTPError` | `Status`, `Body` |
| `graphql` | `*graphql.QueryError` | `Errors []GQLError` (GraphQL-level errors) |

```go
agents, _, err := client.AgentsList(ctx, nil)
if err != nil {
    var apiErr *mgmt.APIError
    if errors.As(err, &apiErr) {
        log.Printf("API error %d: %s", apiErr.Status, apiErr.Title)
    }
}
```

## Raw JSON access

Every struct carries its original JSON in the `Raw` field:

```go
agent, err := client.AgentsGet(ctx, id)
fmt.Println(string(agent.Raw))  // full API response
```

This lets callers access fields not yet mapped to struct fields without
losing data through the typed layer.

## mgmt — REST Management API

72 methods across 15 resource types. Covers agents, threats, sites, groups,
accounts, exclusions, policies, tags, users, and more.

### Listing and pagination

List methods return `(items, *Pagination, error)`. Pass `nil` for defaults:

```go
agents, pag, err := client.AgentsList(ctx, nil)
```

Filter with typed params:

```go
agents, pag, err := client.AgentsList(ctx, &mgmt.AgentListParams{
    Limit:      10,
    OSTypes:    []string{"windows"},
    IsInfected: ptr(true),
})
```

### Single resource

```go
agent, err := client.AgentsGet(ctx, "000000")
site, err := client.SitesGet(ctx, "000000")
threat, err := client.ThreatsGet(ctx, "000000")
```

### CRUD

Create, update, and delete follow the same pattern:

```go
site, err := client.SitesCreate(ctx, mgmt.SiteCreate{
    Name:      "Production",
    SiteType:  "Paid",
    AccountID: "000000",
})
site, err = client.SitesUpdate(ctx, site.ID, mgmt.SiteUpdate{Name: "Prod"})
err = client.SitesDelete(ctx, site.ID)
```

Same pattern for groups, tags, and exclusions.

### Actions

Agent and threat actions use filter-based targeting:

```go
affected, err := client.AgentsIsolate(ctx, mgmt.ActionFilter{IDs: []string{"000000"}})
affected, err = client.ThreatsUpdateVerdict(ctx, mgmt.ActionFilter{
    IDs: []string{"000000"},
}, "true_positive")
```

### Policies

Policies are scoped to site, account, or group:

```go
policy, err := client.PolicyGetSite(ctx, siteID)
policy, err = client.PolicyUpdateSite(ctx, siteID, updatedJSON)
```

### Method catalog

| Resource | Methods |
|----------|---------|
| Agents | `List`, `Get`, `Count`, `Isolate`, `Connect`, `Disconnect`, `InitiateScan`, `AbortScan`, `Shutdown`, `Uninstall`, `Decommission`, `UpdateSoftware`, `MoveToSite`, `FetchLogs`, `RestartMachine`, `EnableAgent`, `DisableAgent`, `ResetLocalConfig`, `ApproveUninstall`, `RejectUninstall`, `MarkUpToDate`, `SetExternalID`, `RandomizeUUID`, `FirewallLogging` |
| Threats | `List`, `Get`, `Mitigate`, `UpdateStatus`, `UpdateVerdict`, `AddToBlacklist`, `FetchFile` |
| Sites | `List`, `Get`, `Create`, `Update`, `Delete` |
| Groups | `List`, `Get`, `Create`, `Update`, `Delete` |
| Tags | `List`, `Get`, `Create`, `Update`, `Delete` |
| Exclusions | `List`, `Get`, `Create`, `Update`, `Delete` |
| Policies | `GetSite`, `GetAccount`, `GetGroup`, `UpdateSite`, `UpdateAccount`, `UpdateGroup` |
| Accounts | `List`, `Get` |
| Users | `List`, `Get`, `Delete` |
| Applications | `List` |
| Device control | `List`, `Get` |
| Firewall | `List`, `Get` |
| Remote scripts | `List`, `Get` |
| Updates | `List`, `Get` |
| Activities | `List` |

## sdl — Singularity Data Lake

13 methods covering queries, log ingest, and file operations. Two protocols:
GraphQL (default, through the management console) and REST (requires a
separate SDL URL).

### PowerQuery

GraphQL (default) — no separate SDL URL needed:

```go
resp, err := client.PowerQueryGraphQL(ctx, &sdl.PowerQueryRequest{
    Query:     "endpoint.name contains 'srv'",
    StartTime: "24h",
})
for _, row := range resp.Values {
    fmt.Println(row)
}
```

REST — requires `S1_SDL_URL`:

```go
sdlClient := sdl.NewClient("https://xdr.us1.sentinelone.net", token)
resp, err := sdlClient.PowerQuery(ctx, &sdl.PowerQueryRequest{
    Query:     "endpoint.name contains 'srv'",
    StartTime: "24h",
})
```

Both return the same `*PowerQueryResponse` — callers can switch transparently.

### GraphQL query lifecycle

For advanced use, control the launch-poll-cleanup cycle directly:

```go
result, err := client.LaunchQuery(ctx, group)
result, err = client.PingQuery(ctx, result.IDs, result.StepsCompleted, result.Token)
err = client.RemoveQuery(ctx, result.Token)
```

### Other queries

```go
resp, err := client.Query(ctx, &sdl.LogQueryRequest{...})
resp, err := client.FacetQuery(ctx, &sdl.FacetQueryRequest{...})
resp, err := client.TimeseriesQuery(ctx, &sdl.TimeseriesQueryRequest{...})
```

### Ingest

```go
err := client.AddEvents(ctx, &sdl.AddEventsRequest{Events: events})
err = client.UploadLogs(ctx, &sdl.UploadLogsRequest{
    Token:  token,
    Logs:   logData,
    Parser: "json",
})
```

### Files

```go
data, err := client.GetFile(ctx, path)
err = client.PutFile(ctx, path, content)
files, err := client.ListFiles(ctx, prefix)
```

### Method catalog

| Category | Methods |
|----------|---------|
| PowerQuery | `PowerQueryGraphQL`, `PowerQuery` |
| GraphQL ops | `LaunchQuery`, `PingQuery`, `RemoveQuery` |
| Queries | `Query`, `FacetQuery`, `TimeseriesQuery` |
| Ingest | `AddEvents`, `UploadLogs` |
| Files | `GetFile`, `PutFile`, `ListFiles` |

## graphql — GraphQL API

15 methods across 4 domains: UAM alerts, xSPM vulnerabilities, xSPM
misconfigurations, and cloud security policies. All use Relay-style pagination.

### Listing with pagination

List methods accept `*ListParams` and return `*Connection[T]`:

```go
alerts, err := client.AlertsList(ctx, &graphql.ListParams{First: 10})
for _, edge := range alerts.Edges {
    fmt.Println(edge.Node.Name, edge.Node.Severity)
}
if alerts.PageInfo.HasNextPage {
    next, err := client.AlertsList(ctx, &graphql.ListParams{
        First: 10,
        After: alerts.PageInfo.EndCursor,
    })
}
```

### Filtering and scoping

```go
params := &graphql.ListParams{
    First: 50,
    Filters: []graphql.Filter{{
        FieldID:  "severity",
        StringIn: &graphql.InStr{Values: []string{"Critical", "High"}},
    }},
    Scope: &graphql.Scope{
        ScopeIDs:  []string{"000000"},
        ScopeType: "site",
    },
}
vulns, err := client.VulnerabilitiesList(ctx, params)
```

### Single resource and mutations

```go
alert, err := client.AlertsGet(ctx, []string{"000000"})
err = client.AlertsUpdateStatus(ctx, []string{"000000"}, "resolved")
err = client.AlertsUpdateVerdict(ctx, []string{"000000"}, "true_positive")
```

Same pattern for vulnerabilities and misconfigurations:

```go
err = client.VulnerabilitiesUpdateStatus(ctx, []string{id}, "resolved")
err = client.MisconfigurationsUpdateVerdict(ctx, []string{id}, "true_positive")
```

### Cloud policies (read-only)

```go
policies, err := client.CloudPoliciesList(ctx, &graphql.ListParams{First: 50})
policy, err := client.CloudPoliciesGet(ctx, []string{"000000"})
```

### Method catalog

| Domain | Methods |
|--------|---------|
| Alerts | `List`, `Get`, `UpdateStatus`, `UpdateVerdict` |
| Vulnerabilities | `List`, `Get`, `UpdateStatus`, `UpdateVerdict` |
| Misconfigurations | `List`, `Get`, `UpdateStatus`, `UpdateVerdict` |
| Cloud policies | `List`, `Get` |
| Low-level | `Do` (send any GraphQL query to any endpoint) |
