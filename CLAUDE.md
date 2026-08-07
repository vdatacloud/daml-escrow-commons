# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**daml-escrow-commons** is a shared Go module for utilities used by more than one repo in this platform —
currently `../daml-escrow` (ledger/settlement) and `../daml-escrow-cms` (contract authoring/CLM). It holds tools,
utilities, and validators, deliberately excluding anything domain-specific: no ledger client code, no escrow
business logic, no T1/T2/T3 identity resolution. See `README.md` for the current package list and rationale, and
`.claude/skills/commons-contribution/SKILL.md` before adding anything new here.

## Commands

- `go build ./...` — build all packages.
- `go test ./...` — run all unit tests (no network/Docker dependencies; keep it that way — this module is a
  dependency of two other repos' test suites and must stay fast).
- `go vet ./...` — static checks.
- See `RELEASING.md` for cutting a tagged version — this module's path (`github.com/vdatacloud/daml-escrow-commons`)
  matches its repo location so tagged releases are `go get`-resolvable, not just usable via local `replace`.

## Architecture

One flat-ish package per concern (`schema`, `hmacsig`, `validate`, …), each independently importable — a consumer
that only needs `hmacsig` shouldn't have to pull in `schema`'s `gojsonschema` dependency. Every package must:

- Have zero dependency on either consumer repo's domain types (no `EscrowContract`, no `DraftEscrow`, no
  `damlPartyId`) — take primitives (`[]byte`, `string`, `float64`) and return primitives or small local types.
- Ship with full unit test coverage in the same PR that adds it — this module has no integration-test tier, so
  unit tests are the only verification either consumer repo gets before depending on a change.
- Wrap errors with context (`fmt.Errorf("...: %w", err)`), matching both consumer repos' convention.

## Conventions carried over from `daml-escrow` / `daml-escrow-cms`

- Git: conventional commit prefixes (`feat:`, `fix:`, `docs:`, `refactor:`, `test:`, `chore:`), commits signed
  off (`-s`)/GPG-signed, changes land via PR rather than direct pushes to `main`.
- No scattered config/env reads — this module takes configuration as constructor/function arguments from its
  caller, never reads its own env vars or files outside of what's explicitly passed in (e.g. `schema.LoadDirectory`
  takes a directory path argument, it does not default to a well-known path).

## Relationship to the other repos

- `../daml-escrow/CLAUDE.md` and `../daml-escrow-cms/CLAUDE.md` are each other's ledger-side and CLM-side
  authorities; this repo doesn't duplicate either. If you're tempted to add domain logic here, it almost certainly
  belongs in one of those two repos instead — see the promotion checklist in
  `.claude/skills/commons-contribution/SKILL.md`.
- `../daml-escrow/plans/CMS_SEPARATION_PLAN.md` and `../daml-escrow-cms/INTEGRATION.md` reference this module for
  shared schema validation; keep this repo's `schema` package behavior consistent with what those documents
  describe.
