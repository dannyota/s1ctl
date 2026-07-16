# Automation

Manage SentinelOne Hyperautomation workflows: list, inspect, import, run, and
control lifecycle (activate/deactivate).

Hyperautomation workflows automate response actions, enrichment, and
notification across the SentinelOne platform. Each workflow has versioned
definitions that can be exported, imported, and promoted between environments.

## Prerequisites

- s1ctl [installed](install.md) and [configured](configure.md)
- `S1_CONSOLE_URL` and `S1_TOKEN` set
- Hyperautomation module enabled on your console

## List workflows

```bash
s1ctl automation list
s1ctl automation list --json
```

## Get a workflow version

```bash
s1ctl automation get 000000
s1ctl automation get 000000 --json
```

Returns the workflow in export format (suitable for re-import).

## Versions

List all versions of a workflow:

```bash
s1ctl automation versions 000000
```

## Export and import

Export a workflow version as JSON for backup or promotion to another
environment:

```bash
s1ctl automation export 000000 > workflow.json
```

Import (create) a workflow from an exported JSON file:

```bash
s1ctl automation create --from-file workflow.json          # dry-run
s1ctl automation create --from-file workflow.json --yes    # apply
```

The create command is **dry-run by default**; pass `--yes` to apply.

## Lifecycle

Activate or deactivate a workflow version:

```bash
s1ctl automation activate 000000 --yes
s1ctl automation deactivate 000000 --yes
```

Both are **dry-run by default**; pass `--yes` to apply.

## Run a workflow

Trigger a manual workflow execution:

```bash
s1ctl automation run 000000 --yes
```

The run command is **dry-run by default**; pass `--yes` to apply.

## Executions

List workflow executions and inspect results:

```bash
s1ctl automation executions 000000
s1ctl automation execution-get 000000
s1ctl automation execution-output 000000
```

`executions` lists all runs for a workflow. `execution-get` returns details
for a specific execution. `execution-output` shows the output payload.

## Workflows (patterns)

### Promote a workflow between environments

Export from one console, import to another:

```bash
# On source console
s1ctl automation export 000000 > workflow.json

# On target console (different S1_CONSOLE_URL / S1_TOKEN)
s1ctl automation create --from-file workflow.json --yes
```

### Audit workflow health

List all workflows and check which are active:

```bash
s1ctl automation list --json > workflows.json
```

### Monitor execution results

```bash
s1ctl automation executions 000000 --json
s1ctl automation execution-output 000000 --json
```

## See also

- [The loop](the-loop.md) — core pull/review/push mental model
- [Examples](examples.md) — common workflows
