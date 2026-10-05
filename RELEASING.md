# Releasing

## Automated (default)

Every push to `main` runs `.github/workflows/release.yml`: it computes the next version from the
Conventional Commits since the last `v*` tag (`scripts/next-version.sh`, rules in section 2 below), re-runs
`go build`/`go vet`/`go test`, then pushes an annotated tag and publishes a GitHub release with generated
notes. **Merging a PR is releasing it** — so the PR's commit prefixes decide the bump:

| Commits since the last tag | Bump |
|---|---|
| `type!:` or a `BREAKING CHANGE:` footer (breaks a consumer: another repo, the website) | minor |
| `feat:` (a substantial feature) or a `Deprecated:` footer | minor |
| `fix:` / `perf:` / `refactor:` / `revert:` / `build:` / `chore(deps):` (small changes, bug fixes) | patch |
| only `chore`/`docs`/`test`/`ci`/`style` | **no release** |

Majors are never automatic -- see "Major releases" below.

Preview what merging would release: `scripts/next-version.sh <branch>` (empty output = no release). Tags are
made by `github-actions[bot]` (annotated, not GPG-signed). After a release, update consumers (section 5).

## Major releases

A major version is a **product release** (the website or the daml-escrow platform), decided deliberately, not
triggered by a commit prefix. Until the first one, every repo stays at `v0.x`, where a minor bump is semver's
signal for a breaking change -- mark breaks with `type!:` or a `BREAKING CHANGE:` footer anyway, so the release
notes say so. To cut a major, tag it by hand (manual steps below) with the next major (`v1.0.0`), after
`make verify`; automation then continues from that tag.

The manual process below remains the fallback (e.g. a re-release, or if the workflow is unavailable).

## Manual It follows [Semantic Versioning](https://semver.org/) and
reads version bumps off [Conventional Commits](https://www.conventionalcommits.org/) prefixes, matching the commit
convention already used across this repo and `daml-escrow`/`daml-escrow-cms`.

Go's own module system is the reason this matters here more than in a typical app repo: a consumer pins to a git
tag (`go get github.com/vdatacloud/daml-escrow-commons@v0.2.0`), so a tag is a real, load-bearing public API
commitment the moment anything depends on it by version instead of by local `replace`.

## 1. Prerequisites (one-time, per machine that will `go get` this module)

This is a **private** repo, so plain `go get` fails until a consumer's environment is configured:

```bash
# Tell Go not to look this module up on the public proxy/checksum DB
go env -w GOPRIVATE=github.com/vdatacloud/*

# Tell git to use your SSH identity for github.com fetches (adjust the alias/key
# to whatever this machine actually uses — see daml-escrow's ~/.ssh/config
# github-dhushon-personal entry for the pattern)
git config --global url."git@github.com:".insteadOf "https://github.com/"
```

Without `GOPRIVATE`, `go get` tries to verify the module against the public checksum database and fails with a
404 on this private repo. Without the SSH rewrite, it tries an anonymous HTTPS clone and fails the same way.

## 2. Decide the version bump

From the repo root:

```bash
git fetch --tags
git log $(git describe --tags --abbrev=0 2>/dev/null || echo "$(git rev-list --max-parents=0 HEAD)")..HEAD --oneline
```

(The `2>/dev/null || ...` fallback handles the very first release, when no tag exists yet — it walks from the
root commit instead.)

Read the commit prefixes in that range:

Apply the same rules as the automated table above (majors only as a deliberate product release, see "Major
releases"). Default to `v0.1.0` for the first release.

## 3. Tag and push

```bash
NEW_VERSION=v0.1.0   # substitute the version decided in step 2

git checkout main
git pull
go build ./... && go vet ./... && go test ./...   # or: make verify

git tag -s "$NEW_VERSION" -m "$NEW_VERSION"
git push origin "$NEW_VERSION"
```

Use `-s` (GPG-signed) to match this repo's commit-signing convention; fall back to `-a` (annotated, unsigned)
only if signing isn't set up on the machine doing the release.

## 4. Publish the GitHub release

```bash
gh release create "$NEW_VERSION" --title "$NEW_VERSION" --generate-notes
```

`--generate-notes` builds the changelog from merged PRs since the previous tag — review and edit it for anything
that needs plain-language explanation (a schema-validation behavior change, a breaking rename) before publishing.

## 5. Update consumers

In each consumer repo (`daml-escrow`, and `daml-escrow-cms` once it has a `go.mod`):

```bash
# Remove or comment out the replace directive for local dev, then:
go get github.com/vdatacloud/daml-escrow-commons@$NEW_VERSION
go mod tidy
```

Keep the `replace` directive present but commented-out (not deleted) between releases — see `go.mod` in
`daml-escrow` for the pattern — so switching back to local-path development for the next round of changes is a
one-line uncomment, not a re-add.

## Future automation

Once there's a second real consumer actively pinning versions (not just the inert local `replace` staged in
`daml-escrow/go.mod` today), it's worth wiring up `release-please` or an equivalent tag-on-merge GitHub Action so
this stops being manual. Not done yet because there's no real release history to bootstrap it from — see
`.claude/skills/commons-contribution/SKILL.md`'s promotion philosophy (don't build ahead of a real need).
