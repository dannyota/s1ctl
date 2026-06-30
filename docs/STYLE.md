# Docs style guide

The contract every doc in `docs/` follows. Keep it short; keep it true to the code.

## Where a doc goes

| Folder | Audience | Answers |
|--------|----------|---------|
| `design/` | devs building **s1ctl** | *how is it built?* (architecture, surfaces, catalog) |
| `guides/` | operators **using** s1ctl | *how do I do X?* (install, auth, per-area how-tos, SDK) |

Root holds only the map (`README.md`), this guide (`STYLE.md`), and the sidebar.
One concept per file. If a file passes **450 lines**, split it.

## Voice

- **Short, dense, technical.** State what's true; cut filler, hedging, and history.
- **Active, imperative** in guides ("Pull the agents, filter, export.").
- **Tenant-neutral always** — placeholders only (`your-console`, `000000`).
  Never a real account/site/host/IP/rule name.

## Format

- **Tables and lists over prose** for any set of things (surfaces, flags, steps).
- **Fenced code** for every command/snippet; show the command, then the why.
- **Mermaid for every flow or structure** — a diagram beats a paragraph.
- **One H1** per file (the title). Sentence-case headings.

## Mermaid

Use a fenced ` ```mermaid ` block (renders on GitHub **and** the docsify site).

Diagrams are **part of the doc-to-code contract**: a diagram must match what the
code does. When the code changes, update the diagram in the same change.

Keep `<br/>` for line breaks (not `\n`) and escape literal angle brackets as
`&lt;`/`&gt;`. The theme toggle re-renders diagrams automatically.

## Formatting rules

- A **blank line before every table, list, and fenced block**, and **after every
  heading** — consistent whitespace keeps the source readable and the linter happy.
- New page under `docs/`? Add it to `docs/_sidebar.md` or it's unreachable.

## Length and lint

- **A doc is capped at 450 lines.** Over it: split into a focused page or trim.
- **Fenced code blocks must declare a language** (` ```bash `, ` ```go `, ` ```text `).
- Run `npx markdownlint-cli2 "docs/**/*.md"` before committing docs changes.
