#!/bin/bash
set -e

# Resolve paths
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
ROOT_DIR="$( cd "$SCRIPT_DIR/.." && pwd )"
GIT_DIR="$ROOT_DIR/.git"

if [ ! -d "$GIT_DIR" ]; then
    echo "Error: .git directory not found. Are you in a git repository?"
    exit 1
fi

HOOKS_DIR="$GIT_DIR/hooks"
mkdir -p "$HOOKS_DIR"

PRE_PUSH_HOOK="$HOOKS_DIR/pre-push"

echo "Installing git pre-push hook..."

cat << 'EOF' > "$PRE_PUSH_HOOK"
#!/bin/bash

# Pre-push hook that runs make verify
echo "--------------------------------------------------------"
echo "Running local verification checks before pushing..."
echo "--------------------------------------------------------"

ROOT_DIR="$(git rev-parse --show-toplevel)"
cd "$ROOT_DIR"

make verify
status=$?

if [ $status -ne 0 ]; then
    echo "--------------------------------------------------------"
    echo "ERROR: Local verification checks failed."
    echo "Push aborted. Please fix the errors before pushing."
    echo "--------------------------------------------------------"
    exit 1
fi

echo "--------------------------------------------------------"
echo "SUCCESS: Verification checks passed. Proceeding with push."
echo "--------------------------------------------------------"
exit 0
EOF

chmod +x "$PRE_PUSH_HOOK"
echo "Pre-push hook successfully installed at: .git/hooks/pre-push"

# graphify-out/ is regenerated on main by .github/workflows/update-graph.yml
# after every merge, and that copy is the one everyone uses -- no local hook
# writes it (a local run just churns the tracked files, and CI refuses PRs
# that touch them). Remove the graphify pre-commit hook earlier versions of
# this script installed.
PRE_COMMIT_HOOK="$HOOKS_DIR/pre-commit"
if [ -f "$PRE_COMMIT_HOOK" ] && grep -q "graphify update" "$PRE_COMMIT_HOOK"; then
    rm "$PRE_COMMIT_HOOK"
    echo "Removed the old graphify pre-commit hook (graphify-out/ comes from CI on main)."
fi
