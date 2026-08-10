# Repository Guardrails

This is the canonical source-control policy for this repo — mirrored from
`../daml-escrow/.gemini/repo_rules.md` and kept consistent with the same
sections in `../daml-escrow-cms` and `../daml-escrow-identity`. `.agents/skills/
github-ci/SKILL.md` covers the operational commands (git/gh invocations, CI
failure triage) that this file doesn't spell out — if the two ever disagree,
this file wins and the skill should be updated to match.

## Branching Strategy

There is no `develop` branch — `main` is the only long-lived branch.

main → production ready, PR-only, no direct pushes except the bugfix exception below

feature/\* → new work

fix/\* → bug fixes

docs/\* → documentation-only changes

chore/\* → tooling, dependencies, CI config

sec/\* → security vulnerability remediation

------------------------------------------------------------------------

## Commit Standard

Use conventional commits:

feat: fix: docs: refactor: test: chore: sec:

**MANDATORY:** Every commit MUST carry both:
- DCO sign-off (`-s` / `--signoff`)
- A GPG signature (via `git config commit.gpgsign true`, or `git commit -S`)

Only "Verified" (GPG-signed) commits are eligible for merge into `main`.

------------------------------------------------------------------------

## Pull Request Rules

**MANDATORY:** ALL changes to the codebase MUST be made as a formal Pull Request.

- Direct pushes to `origin/main` are NOT allowed.
- Exception: Bug fixes may be pushed directly to `main` ONLY after asking and receiving explicit user confirmation.

PR must include:

- description
- what's affected (which package(s) — see the promotion checklist in
  `.claude/skills/commons-contribution/SKILL.md` for whether a change even
  belongs in this repo)
- tests
- security considerations (this module is a dependency of every consumer
  repo's build — a bad change here breaks all of them at once)

------------------------------------------------------------------------

## Code Review Requirements

Minimum: **1 reviewer** for any change (this repo has no smart-contract-equivalent
2-reviewer tier — everything here is a shared library dependency, treated
uniformly).

## Pre-Commit Verification

**MANDATORY:** Before committing and pushing code, run `make verify` (build, vet,
unit tests) — this is also what the `pre-push` git hook (installed via
`make install-hooks`) runs automatically.

```bash
make verify
# equivalent to: go build ./... && go vet ./... && go test ./...
```

------------------------------------------------------------------------

## CI Requirements

Every PR to `main` runs the two jobs defined in `.github/workflows/ci.yml`:

- `lint` — `golangci-lint` across `./...`
- `unit-tests` — `go build ./...` then `go test -v ./...`

No integration-test tier exists here — this module has zero network/Docker
dependencies by design (see `CLAUDE.md`).

------------------------------------------------------------------------

## Branch Protection Strategy

**Current reality:** GitHub branch protection is not turned on for `main`, and
cannot be today — the `vdatacloud` org is on GitHub's Free plan, and both the
classic protection API and the rulesets API return `403: Upgrade to GitHub Pro
or make this repository public to enable this feature` for a private repo on
that plan (verified 2026-08-10 against `../daml-escrow`, same org). The rules
above are enforced **only** by this document, the `github-ci` skill, and
human/agent discipline — not by GitHub itself.

**If the org plan changes**, apply:

```bash
gh api repos/vdatacloud/daml-escrow-commons/branches/main/protection -X PUT --input - <<'EOF'
{
  "required_status_checks": {
    "strict": true,
    "contexts": ["Go Lint", "Go Unit Tests"]
  },
  "enforce_admins": true,
  "required_pull_request_reviews": {
    "required_approving_review_count": 1,
    "dismiss_stale_reviews": true
  },
  "restrictions": null,
  "allow_force_pushes": false,
  "allow_deletions": false,
  "required_signatures": true
}
EOF
```

Keep this section consistent with the equivalent section in `../daml-escrow`,
`../daml-escrow-cms`, and `../daml-escrow-identity` if any of them change.
