#!/usr/bin/env bash
# next-version.sh -- prints the next semver tag for this module from the
# Conventional Commits since the latest v* tag, or prints nothing if no
# release is warranted. Used by .github/workflows/release.yml; run it locally
# to preview what a merge to main would release. Rules (RELEASING.md):
#
#   feat!: / fix!: / any "BREAKING CHANGE:" footer -> major (pre-1.0: minor)
#   feat:                                         -> minor
#   fix: / perf:                                  -> patch
#   only chore/docs/test/ci/refactor/style/build  -> no release
#
# Usage: scripts/next-version.sh [range-end]   (default HEAD)
set -euo pipefail

END="${1:-HEAD}"
LAST_TAG="$(git describe --tags --abbrev=0 --match 'v[0-9]*.[0-9]*.[0-9]*' "$END" 2>/dev/null || true)"

if [ -z "$LAST_TAG" ]; then
  RANGE="$END"
  MAJOR=0 MINOR=0 PATCH=0
else
  RANGE="${LAST_TAG}..${END}"
  IFS=. read -r MAJOR MINOR PATCH <<<"${LAST_TAG#v}"
fi

# Full messages (subject + body), NUL-separated, so footers are seen too.
bump=none
while IFS= read -r -d '' msg; do
  subject="${msg%%$'\n'*}"
  if [[ "$subject" =~ ^[a-z]+(\([^\)]*\))?!: ]] || grep -q '^BREAKING[ -]CHANGE:' <<<"$msg"; then
    bump=major; break
  elif [[ "$subject" =~ ^feat(\([^\)]*\))?: ]]; then
    bump=minor
  elif [[ "$subject" =~ ^(fix|perf)(\([^\)]*\))?: ]] && [ "$bump" = none ]; then
    bump=patch
  fi
done < <(git log --format='%B%x00' "$RANGE")

# Pre-1.0: a breaking change is a minor bump (RELEASING.md "Pre-1.0 note").
if [ "$bump" = major ] && [ "$MAJOR" -eq 0 ]; then
  bump=minor
fi

case "$bump" in
  major) echo "v$((MAJOR + 1)).0.0" ;;
  minor) echo "v${MAJOR}.$((MINOR + 1)).0" ;;
  patch) echo "v${MAJOR}.${MINOR}.$((PATCH + 1))" ;;
  none)  : ;;
esac
