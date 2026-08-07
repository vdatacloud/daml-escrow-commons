---
name: commons-contribution
description: "Use before adding, moving, or removing code in daml-escrow-commons — the promotion checklist for what belongs in a shared module used by both daml-escrow and daml-escrow-cms, versus what should stay local to one repo."
---

# /commons-contribution

Checklist for deciding whether something belongs in `daml-escrow-commons`, and how to add it correctly if so.

## When to promote code into commons

All of these must be true, not just one:

1. **Two or more repos actually need it**, not "might need it someday." A single current consumer (even if a
   second is plausible) stays local — `daml-escrow`'s `internal/services/schema_service.go` and
   `internal/services/compliance.go` still exist as the production implementations; they have not yet been
   migrated to depend on this module's `schema`/`hmacsig` packages, because doing that migration wasn't the ask
   that created this repo. Promotion into commons and migration of the *consumer* are two separate steps — don't
   assume the second happened just because the first did.
2. **It's a pure utility, not domain logic.** No `EscrowContract`, `DraftEscrow`, `damlPartyId`, milestone/dispute
   rules, custody thresholds, or anything from `ESCROW-PROCESS.md`'s DIRECTIVEs. If you're writing an `import` for
   a domain type from either consumer repo, stop — it doesn't belong here.
3. **It doesn't touch T2/T3 identity.** `daml-escrow-cms` is explicitly scoped to T1-only identity
   (`../daml-escrow-cms/INTEGRATION.md` §4). A utility that assumes ledger User IDs or `damlPartyId`s would tempt
   CMS to reach across that boundary — keep it in `daml-escrow` instead.
4. **It can ship with full unit test coverage in the same change.** This module has no integration-test tier — two
   consumer repos' test suites and production behavior end up depending on whatever ships here untested.

## How to add something

1. One package per concern, named for what it does (`schema`, `hmacsig`, `validate`), not for which repo asked for
   it.
2. Take configuration as arguments, not env vars or well-known file paths — the caller (in either consumer repo)
   owns config loading via its own `internal/config` (or equivalent); this module stays stateless and portable.
3. Wrap errors with context (`fmt.Errorf("<package>: ...: %w", err)`).
4. Write the tests first or alongside — table-driven where the primitive has more than 2-3 cases, one test
   function per behavior otherwise (see `validate/validate_test.go` for the pattern this repo already uses).
5. Update `README.md`'s package table with what it does and *why it's shared, not local* — that second half is
   the part reviewers actually need, and the part most likely to be missing.
6. Run `go build ./... && go vet ./... && go test ./...` before committing — this module is a dependency, a
   build break here breaks two other repos' builds once they actually depend on it.

## When NOT to add something

- It solves a problem only one repo currently has. Leave it there; revisit when a second repo actually needs the
  same thing (per `../daml-escrow/plans/CMS_SEPARATION_PLAN.md`'s general split philosophy: don't split/share
  ahead of a real second consumer).
- It's a naming/style convention rather than code (e.g. "always wrap errors with `%w`"). Document conventions like
  that in each repo's own `CLAUDE.md`, not as code here.
- It would require either consumer repo to import the other's domain types to use it. That's a sign the utility
  is actually domain logic wearing a utility's clothes.

## Removing something

If a package's last real consumer stops using it (migrated to something else, or the capability moved), delete it
in the same change that removes the last call site — don't leave dead packages "just in case." Note the removal
and why in the commit message; there's no separate deprecation-log file to update.
