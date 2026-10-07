# daml-escrow-commons

Shared, dependency-light Go utilities used by both `../daml-escrow` and `../daml-escrow-cms` (and any future
sibling repo). A Go module, imported via `go.mod` (local `replace` directive during development, tagged versions
once either consumer repo needs stability across independent releases).

## What's here

| Package | What it does | Why it's shared, not local |
|---|---|---|
| `schema` | Loads a directory of JSON Schema files, validates arbitrary JSON payloads against them by type name. | Both repos validate against the *same* schema authority (`daml-escrow/architecture/schemas/*.json`) — daml-escrow for `EscrowMetadata`, daml-escrow-cms for structurally-extracted contract terms before draft creation. One implementation, one behavior. |
| `hmacsig` | HMAC-SHA256 sign/verify with constant-time comparison. | Every webhook-style integration on the platform needs the same primitive — daml-escrow's oracle/fiat-settlement webhooks today, daml-escrow-cms's import/OCR callbacks and any third-party-CLM substitution point tomorrow. |
| `pseudoref` | Short (12-hex), stable, non-reversible keyed references (HMAC-SHA256) to sensitive ids. | Operational telemetry (traces, timing records) needs enough of a tenant/transaction reference to cross-reference, never the id itself -- daml-escrow PLAN.md Phase 68. |
| `validate` | Small composable field-level checks (`RequireNonEmpty`, `RequirePositive`, `RequireOneOf`, `RequireValidEmail`) plus an `Errors` aggregator, for the `.Validate() error` DTO convention both repos use. | The primitives are identical everywhere; only the domain-specific composition (which fields, which rules) differs per DTO and stays local. |
| `metering` | `LedgerCommandEvent`/`SettlementEvent` — price-free usage-fact shapes for `daml-escrow-platform`'s rating/billing layer (Phase 34), plus their `Validate()`. | Both daml-escrow and daml-escrow-cms submit ledger commands and must emit identically-shaped usage events, or `daml-escrow-platform`'s per-tenant aggregation needs repo-specific special-casing. One shape, one behavior, same as `schema`'s rationale. |
| `cantonid` | Typed Canton identifiers -- `Fingerprint` (a key: `1220`+64 hex, `FingerprintOf(ed25519)`), `PartyID` (`hint::fingerprint`; `ControlledBy(fp)` for external parties), `Hash` (transaction/topology) -- each with one parser, full-value JSON, and one `Short()` display form; `Short(string)` for untyped values. Golden `testdata/parse.json`. | Every service and UI that shows or checks a party, key or hash must parse and abbreviate it identically (`relaytest::1220ebb7…4288`); `@vdatacloud/cx-commons` (`sdk/canton-id`) mirrors it against the same fixture. Pure parsing -- no resolution of who a party is. |
| `apierror` | The canonical API error envelope (`{"error","code","status","stage","hint","requestId","upstream","details"}`): `Error` with Full/Summary detail levels and a byte-exact canonical JSON form, `Parse` for any error body (canonical, legacy `{"error"}`, plain text), an HTTP `Writer` + `RequestID` middleware, `ShortID` (one short form for long ids in messages and logs) and `LogFields` (logger-agnostic, short forms only). JSON Schema embedded (`Schema()`); golden fixtures in `testdata/canonical`. | Every service (daml-escrow, -cms, -identity, -platform) answers errors, and ~540 handlers still send plain-text `http.Error` -- one envelope means one parser and one error UI. `@vdatacloud/cx-commons` (`sdk/api-error`) mirrors it in TypeScript and tests against the same fixtures, so Go and the UX produce identical bytes. |
| `apierror/canton` | Turns a Canton JSON Ledger API refusal into an `apierror.Error`: parses Canton's error body (`LedgerError`), maps its gRPC code to status/code/hint, lifts the expected transaction hash out of a signature refusal, and defines `SigningDetails` for signing/authorization failures. | Any service submitting ledger commands (daml-escrow today; daml-escrow-cms emits ledger-command metering too) needs the same classification, and the cx-commons mirror (`fromCantonError`) gives a UI that talks to a participant or wallet gateway directly the same treatment -- fixtures in `testdata/classify` are shared. Pure wire-format parsing: party ids are opaque strings, no T2/T3 resolution, so it doesn't pull CMS across its T1 boundary. |
| `agreementsig` | How a party signs one version of an agreement: a canonical JSON form and SHA-256 version hash, a fixed-English signed message naming the version, its hash and the signer, and `Verify` for signatures that name their algorithm (Ed25519 for wallets, ECDSA P-256/SHA-256 low-s DER for KMS keys such as the custodian's; `NormalizeECDSA` for KMS output). Two frozen schemas: `tripart.draft-version/1` (moved from daml-escrow `internal/draftsig`, Phase 77) and `tripart.agreement-version/1` (Phase 81: id, version, parent hash, amendment chain, rendered-document and canonical-source hashes -- never the storage location -- and terms as an opaque JSON object named by its schema, integers only, money as decimal strings). Golden `testdata/*.json`. | daml-escrow-cms produces and collects approvals for agreement versions, daml-escrow verifies and records them on the ledger, and a browser recomputes the hash of what it shows before signing -- three implementations that must agree byte for byte. `@vdatacloud/cx-commons` (`sdk/agreement-sig`) mirrors it against the same fixtures. Domain-free: parties, roles and terms are opaque strings/JSON (the escrow-terms schema lives in daml-escrow); CMS names parties by identity token, so nothing here assumes T2/T3. |
| `identityclient` | The minimal T1 HTTP client for `daml-escrow-identity` — `Upsert`/`GetByOktaSub`/`GetByEmail`/`GetByToken`/`ManagesIdentity`. Pure wire-level transport, no resolution/reconciliation logic. | Both repos' local `internal/identityclient` packages were near-identical hand-rolled copies of this same subset (daml-escrow-cms's a strict subset of daml-escrow's). daml-escrow-cms uses this package directly, unmodified, since T1 is all it ever needs. daml-escrow composes its own local client on top via struct embedding, adding T2/T3 write-back and party-set CRUD that stay out of this module (see its own package doc). |

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
(`go get github.com/vdatacloud/daml-escrow-commons@vX.Y.Z`) — see `RELEASING.md` for cutting one. This repo is
public, so no auth/module-fetch setup is needed to consume a tagged version once one exists.

## Status

Scaffolded 2026-08-07 with `schema`, `hmacsig` and `validate`; since grown to the ten packages above, each with
unit tests (`go test ./...`). Consumers: `daml-escrow-cms` (`schema`, `hmacsig`, `identityclient`),
`daml-escrow-identity` (`hmacsig`) and `daml-escrow` (`agreementsig`, `apierror`, `apierror/canton`, `cantonid`, `hmacsig`,
`identityclient`, `metering`, `pseudoref`). `agreementsig`, `apierror` and `cantonid` are also mirrored in TypeScript by
`@vdatacloud/cx-commons`, tested against this repo's `testdata` fixtures -- change them here first.

What an API error may disclose (caller-sent values, values from the caller's own credentials, and -- once verified --
facts about its own resources; never another party's data or existence) is spelled out in the `apierror` package
doc; every service emitting the envelope follows it.

Made public 2026-08-10 (it already fit the bar: dependency-light, no domain logic, no ledger client code, no T2/T3
identity) — this also resolved a CI cross-repo-checkout problem in `daml-escrow-cms`/`daml-escrow-identity` for
free, with no PAT or deploy key to manage.
