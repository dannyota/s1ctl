# status

Show environment health summary

## status capabilities

Show s1ctl version, config, and API reachability

```text
s1ctl status capabilities
```

## status enums

List known enum values used across the CLI

```text
s1ctl status enums [flags]
```

Show all known enum values grouped by domain. Use --group to filter
to a specific domain (e.g. alerts, threats, agents, rules, exclusions, policies).

**Flags**

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--group` | string | - | filter to a specific enum group |

## status surfaces

List all API surfaces and supported operations

```text
s1ctl status surfaces
```
