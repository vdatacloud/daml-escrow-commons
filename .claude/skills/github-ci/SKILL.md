---
name: github-ci
description: Git and GitHub rules for the daml-escrow repos. Use before you branch, commit, push, open or update a PR, read or fix CI, or change branch protection.
---

# Git and CI

Shared by every daml-escrow repo (`docs/REPOSITORIES.md`); the source is `daml-escrow`. Edit it there, then copy it
over unchanged.

## Policy

- **`main` only.** No `develop`. Every change lands through a PR. A direct push to `main` needs the user's explicit
  go-ahead for that push, and only for a critical fix.
- **Branches** (lowercase kebab-case, from `main`, deleted after merge):
  - Topic branches, one coherent change with its own PR: `feature/` (new capability), `fix/` (bugs, security fixes
    included, with `fix:` commits so they release), `docs/` (documentation only), `chore/` (tooling,
    dependencies, CI, repo hygiene). Add a phase number when it helps: `feature/111-ledger-auth`.
  - Batching branches, to save Actions runs and reviews: `integration/phase-<N>` for a phase's slices, merged once;
    `integration/small-changes` for unrelated small items, recreated from `main` after each merge. Never open a PR
    to `main` for one small change.
  - Part of a phase: its integration branch. Small and not urgent: `integration/small-changes`. Must merge on its
    own: a topic branch.
- **Every commit is signed off and signed:** `git commit -s` with `commit.gpgsign=true` (a key registered on GitHub
  whose identities include the commit email, so it shows Verified). The `repo-policy` workflow fails a PR with any
  commit that is unverified or lacks `Signed-off-by:`. Never pass `-c commit.gpgsign=false`; if signing fails, stop
  and say so.
- **Conventional prefixes:** `feat:` `fix:` `perf:` `refactor:` `revert:` `build:` `docs:` `test:` `ci:` `style:`
  `chore:`. Imperative subject of 72 characters or fewer; a body that says why for anything non-trivial. The prefix
  sets the release bump (`RELEASING.md`), whatever the branch prefix.
- **Review:** one reviewer for now (one developer). Changes to DAML templates and choices, custody or settlement
  code, identity mapping or auth tick the PR template's high-risk line.
- **Never** commit secrets, `.env` values, tokens or client data; rewrite history on a shared branch; skip, disable
  or quarantine a test to get green; or push an empty commit to re-run CI.
- **CI files:** Claude asks before editing `.github/workflows/` or the pre-push hook (`scripts/install-git-hooks.sh`).

## Doing it

1. Run `make verify` (also the pre-push hook from `make install-hooks`). CI runs only on PRs to `main`, so an
   integration branch is checked locally until its PR opens.
2. Stage named files (never `git add .`), check `git diff --cached --stat`, commit with `-s`.
3. Push with `git push -u origin <branch>`, sparingly to an open PR. Open the PR against `main`; the body follows
   `.github/pull_request_template.md`.
4. Watch: `gh pr checks <n>`. On a failure, `gh run view <run-id> --log-failed > <scratch>/ci.log`, read the failing
   lines, fix the root cause, re-run that check locally, push. One re-run at most for a suspected flake.
5. Done when every check is green, there is no conflict and the template is filled in. The user merges.

The jobs a repo runs are the files in its `.github/workflows/`; read them rather than relying on a list here.

## Branch protection

`main` is protected by a repository ruleset (PR required, no approval needed yet, signed commits, required checks,
no force push or deletion). Only the org's `vdatacloud-graph-bot` app may bypass it, so the graph workflow can
push its commit (a one-hour token from the `GRAPH_APP_ID` and `GRAPH_APP_PRIVATE_KEY` org secrets). Check it with
`gh api repos/{owner}/{repo}/rulesets`. Changing it is a user decision: propose it, don't apply it unasked.
