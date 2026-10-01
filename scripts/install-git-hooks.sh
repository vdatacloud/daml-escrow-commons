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

PRE_COMMIT_HOOK="$HOOKS_DIR/pre-commit"

echo "Installing git pre-commit hook..."

# Regenerates graphify-out/ locally for this checkout's own queries, but
# never stages it: .github/workflows/update-graph.yml regenerates and
# commits it on main after merge. Staging it here put graphify-out/ churn in
# every PR, which then conflicted with that bot's commits (same convention
# as daml-escrow's hook). CI rejects PRs that touch graphify-out/.
cat << 'EOF' > "$PRE_COMMIT_HOOK"
#!/bin/sh
echo "[pre-commit] Updating local Graphify memory map (not staged)..."
graphify update . || { echo "Graphify update failed. Aborting commit."; exit 1; }
exit 0
EOF

chmod +x "$PRE_COMMIT_HOOK"

echo "Pre-commit hook successfully installed at: .git/hooks/pre-commit"
