# CLI naming

Command naming conventions for s1ctl. Names follow SentinelOne's official API
terminology — never invented abstractions.

## Rules

1. **Top-level groups are plural nouns** naming the subject: `agents`, `threats`,
   `alerts`, `exclusions`, `vulnerabilities`.
2. **Verbs nest under nouns**: `s1ctl agents list`, `s1ctl agents isolate`.
3. **Protocol is invisible** — users never type `graphql` or `rest`. The CLI
   routes to the right API.
4. **Multi-word names are hyphenated**: `cloud-policies`, `service-users`,
   `tag-rules`.
5. **No aliases** — one surface, one runnable name.
6. **`pull`/`push` nest under each surface**: `s1ctl exclusions pull`,
   `s1ctl rules push`.
7. **`commands --json`** is the machine-readable command catalog.

## Standard verbs

| Verb | Meaning | Example |
|------|---------|---------|
| `list` | Paginated listing with filters | `s1ctl agents list --os windows` |
| `get` | Single resource by ID | `s1ctl threats get <id>` |
| `count` | Count matching resources | `s1ctl agents count --infected true` |
| `query` | Rich query (SDL/Deep Visibility) | `s1ctl visibility query --query ...` |
| `create` | Create a resource | `s1ctl exclusions create --type path ...` |
| `update` | Update a resource | `s1ctl exclusions update <id> ...` |
| `delete` | Delete a resource | `s1ctl exclusions delete <id>` |
| `export` | Export to CSV/JSON | `s1ctl vulnerabilities export --format csv` |

## Action verbs (operational plane)

| Verb | Domain | Example |
|------|--------|---------|
| `isolate` | agents | `s1ctl agents isolate <id> --yes` |
| `reconnect` | agents | `s1ctl agents reconnect <id> --yes` |
| `scan` | agents | `s1ctl agents scan <id> --yes` |
| `mitigate` | threats | `s1ctl threats mitigate <id> --yes` |
| `run` | remoteops | `s1ctl remoteops run <script-id> --yes` |

## Config-as-code verbs

```text
s1ctl <surface> pull [--site-id ID] [--out DIR]
s1ctl <surface> push [--dir DIR] [--yes]
```

`pull` writes one YAML file per object to a directory named after the surface
(`exclusions/`, `sites/`, …). `push` reads the same directory, plans the
difference against live state, and applies creates and updates with `--yes`.
See [Reconcile engine](reconcile.md) for the full semantics.

## Examples

```bash
s1ctl agents list --os linux --json
s1ctl agents isolate abc123 --yes
s1ctl threats list --status active --limit 50
s1ctl alerts list --severity HIGH --json
s1ctl vulnerabilities list --limit 100
s1ctl datalake powerquery --query "EventType = 'Process Creation'"
s1ctl exclusions pull --site-id 000000
s1ctl exclusions push --yes
s1ctl config init
s1ctl doctor
```
