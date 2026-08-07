# daml-escrow-commons

Shared, dependency-light Go utilities used by both `../daml-escrow` and `../daml-escrow-cms` (and any future
sibling repo). A Go module, imported via `go.mod` (local `replace` directive during development, tagged versions
once either consumer repo needs stability across independent releases).

## What's here

| Package | What it does | Why it's shared, not local |
|---|---|---|
| `schema` | Loads a directory of JSON Schema files, validates arbitrary JSON payloads against them by type name. | Both repos validate against the *same* schema authority (`daml-escrow/architecture/schemas/*.json`) — daml-escrow for `EscrowMetadata`, daml-escrow-cms for structurally-extracted contract terms before draft creation. One implementation, one behavior. |
| `hmacsig` | HMAC-SHA256 sign/verify with constant-time comparison. | Every webhook-style integration on the platform needs the same primitive — daml-escrow's oracle/fiat-settlement webhooks today, daml-escrow-cms's import/OCR callbacks and any third-party-CLM substitution point tomorrow. |
| `validate` | Small composable field-level checks (`RequireNonEmpty`, `RequirePositive`, `RequireOneOf`, `RequireValidEmail`) plus an `Errors` aggregator, for the `.Validate() error` DTO convention both repos use. | The primitives are identical everywhere; only the domain-specific composition (which fields, which rules) differs per DTO and stays local. |

## What's NOT here, and won't be

See `.claude/skills/commons-contribution/SKILL.md` for the full promotion checklist. In short: no ledger client
code, no business/domain logic (escrow lifecycle, milestone rules, custody thresholds), no identity resolution
(T1/T2/T3 — daml-escrow-cms explicitly should not touch T2/T3, so nothing here should tempt it to). If it's only
used by one repo today, it stays in that repo until a second consumer actually needs it.

## Using this module

Module path: `github.com/vdatacloud/daml-escrow-commons` (matches the repo location so tagged versions are
resolvable via `go get`, not just local `replace` — see `RELEASING.md`).

During local development, from a consumer repo's `go.mod`:

```
require github.com/vdatacloud/daml-escrow-commons v0.0.0
replace github.com/vdatacloud/daml-escrow-commons => ../daml-escrow-commons
```

Once a tagged release exists, drop the `replace` line and pin the `require` to a real version
(`go get github.com/vdatacloud/daml-escrow-commons@vX.Y.Z`) — see `RELEASING.md` for cutting one, including the
private-repo module-fetch setup a consumer needs one time.

## Status

Scaffolded 2026-08-07 with three packages, each with full unit test coverage (`go test ./...`). Not yet wired up
as a dependency of `daml-escrow` or `daml-escrow-cms` — both repos currently carry their own equivalent code
(`daml-escrow/internal/services/schema_service.go`, `internal/services/compliance.go`'s HMAC verification). See
`../daml-escrow/plans/CMS_SEPARATION_PLAN.md` for the migration note.
