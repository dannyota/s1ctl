# Reconcile engine

The reconcile engine is the shared core behind every config-as-code surface:
`pull` renders live objects to files, `push` plans and applies the difference
between files and live state, and `drift` reports that difference without
applying anything. One engine, one on-disk model, one set of semantics.

## Why

Before the engine, the sync surfaces had grown two divergent models:

| Model | Surfaces | Layout | Push semantics |
|-------|----------|--------|----------------|
| Per-object YAML | rules, firewall, devicecontrol | one file per object | create-or-update, matched by name |
| Single JSON array | sites, groups, tags, exclusions, cloud-policies | one array file per surface | create-only (cloud-policies: status toggle) |

Each surface carried its own ~200-line pull/push implementation. The
duplication made every improvement an eight-file change, the two models
behaved differently under `git diff`, and nothing could answer "does my
committed config match the console?" without pushing.

The engine unifies all of it on the per-object model — the one that makes the
core loop work: **pull live state → review in `git diff` → push back**. A
change to one object is a one-file diff.

## Core types

The engine lives in `internal/reconcile` and imports no SDK. Each surface
supplies closures; the engine owns matching, planning, file I/O, and apply
semantics.

```go
// Object is one config item in canonical file form.
type Object struct {
    Name string // stable identity used for matching (surface-defined)
    ID   string // server ID; "" for local objects not yet created
    Body []byte // canonical declarative body (YAML of the surface's file shape)
}

// Capabilities declares what a surface supports.
type Capabilities struct {
    NoCreate bool // push never creates (e.g. cloud-policies)
}

// Surface adapts one resource to the engine.
type Surface struct {
    Name    string // singular resource noun, e.g. "device rule"
    Command string // CLI group name for guard/audit strings
    Caps    Capabilities
    Decode  func(data []byte) (Object, error) // file bytes → identity + canonical body
    List    func(ctx context.Context) ([]Object, error)
    Create  func(ctx context.Context, local Object) error // nil if NoCreate
    Update  func(ctx context.Context, id string, local Object) error
}
```

### Canonicalization

Bodies are comparable only because both sides pass through the same encoder.
Every surface defines a typed *file shape* struct (its declarative fields,
nothing server-managed). `Decode` unmarshals a local file into that struct and
re-marshals it; `List` converts each live SDK object into the same struct and
marshals it the same way. Key order, indentation, and omitted optional fields
are therefore identical on both sides; comments and unknown fields in local
files are dropped by the round-trip. A local file that fails to decode, or
that lacks its identity field, is a hard error — not a skipped item.

`BuildPlan(local, live)` matches objects by `Name` and classifies every item:

| Kind | Meaning |
|------|---------|
| `create` | local file with no live match |
| `update` | name matches, canonical bodies differ |
| `unchanged` | name matches, bodies byte-equal |
| `live-only` | live object with no local file — reported, never touched |

Duplicate live names resolve to the first object listed; the plan carries a
warning for each duplicate. Two local files declaring the same identity are a
hard error — the input is ambiguous and nothing is applied. Note the
interaction: pulling a console that holds duplicate names writes suffixed
files that share one identity, so resolve the duplicates (delete the extra
file, or fix the console) before pushing.

Identity is by name, so renaming an object locally plans as a `create` of the
new name plus a `live-only` report of the old one. The engine never deletes,
so the old object must be removed (or renamed) in the console; until then,
`drift` keeps reporting it.

`Apply` executes creates and updates with warn-and-continue semantics: a
failing item is reported and skipped, the rest proceed, and the command exits
non-zero with a failure summary when any item failed. On a `NoCreate` surface,
a planned create counts as a failure.

## On-disk layout

One YAML file per object. The filename stem is the sanitized object name;
duplicate stems get `-1`, `-2`… suffixes. Files contain only the declarative
definition — server-assigned IDs, scopes, and timestamps are omitted so diffs
never churn on server-managed fields.

`pull --out` and `push --dir` default to a directory named after the surface
(`sites/`, `firewall/`, …). Pull overwrites; it is the render of live truth.
Pull never deletes local files, but it warns about stale ones — `.yaml` files
in the output directory with no live counterpart. Heed the warning: a stale
file (its live object was deleted, or renamed away) plans as a `create` on the
next push and would re-create the object. Delete stale files before pushing.

## Surfaces

| Surface | Identity | File shape | Caps |
|---------|----------|-----------|------|
| rules | rule name | STAR rule definition | |
| firewall | rule name | firewall rule definition | |
| network | rule name | network quarantine rule definition | |
| devicecontrol | rule name | device rule definition | |
| sites | site name | name, account, type, licenses | |
| groups | site ID + group name | name, site ID, description | |
| tags | tag key | key, value, description, scope | |
| exclusions | type + OS + value | type, value, OS, mode | |
| blocklist | type + OS + value | type, value, OS, hash, description | |
| locations | location name | name, operator, detection conditions | |
| cloud-policies | policy ID | id, name, status | NoCreate |

cloud-policies is status-reconcile only: an `update` toggles the policy
between enabled and disabled; local files whose ID has no live match fail
per-item since policies cannot be created through this surface.

Policies (site/account/group protection policies) stay outside the engine:
they are scope-singletons with their own pull/push/diff/revert lane, not
per-object collections.

## Push semantics

1. Load local files, list live objects, build the plan.
2. Dry-run (default) prints the plan: creates, updates, unchanged count,
   live-only count. Nothing is sent.
3. `--yes` applies creates and updates through the same guarded path as every
   other mutation. Live-only objects are never deleted — there is no delete
   in v1, with or without flags.
4. Per-item failures warn and continue; the command exits non-zero when any
   item failed, and the audit log records the failure.

Because planning needs live state, `push` contacts the API even in dry-run.

## Drift

`s1ctl drift` runs the plan for every surface that has a local directory and
prints a per-surface summary:

```text
SURFACE        CREATE  UPDATE  LIVE-ONLY  UNCHANGED
firewall       0       2       1          14
sites          0       0       0          3
```

Exit code 0 means every checked surface is clean (no creates, no updates, no
live-only objects); 1 means drift. Surfaces without a local directory are
skipped — drift checks only what is committed. The command is read-only by
construction: it lists, plans, and reports, and has no apply path.

Typical CI use: check out the config repo, run `s1ctl drift`, fail the job on
a non-zero exit.

## Migration from the array layout

Sites, groups, tags, exclusions, and cloud-policies previously pulled to a
single JSON array file (`sites.json`, …). From v0.6.0 they use the per-object
layout above; array files are no longer read. Re-pull each surface to
regenerate the local state in the new layout. Push semantics change too:
sites, groups, tags, and exclusions pushed create-only before and now update
matched objects; cloud-policies keeps its status-toggle behavior, expressed
as engine updates.

## Out of scope (v1)

- **Deletes / `--prune`.** Live-only objects are reported, never removed.
- **Optimistic concurrency.** No etag/version tokens on writes.
- **Secret references.** File shapes contain no secret material.
