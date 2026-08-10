# graphify
- **graphify** (`.claude/skills/graphify/SKILL.md` if vendored, otherwise the globally installed skill) - any input to knowledge graph. Trigger: `/graphify`
When the user types `/graphify`, use the installed graphify skill or instructions before doing anything else. `graphify-out/graph.json` exists in this repo — prefer `graphify query "<question>"` over raw grep/read for codebase questions, per `daml-escrow/.claude/CLAUDE.md`'s convention. Run `graphify update .` after changes (the pre-commit hook installed via `scripts/install-git-hooks.sh` does this automatically and stages `graphify-out/`).

# Cross-repo context

This repo is the shared-utilities leg of a four-repo system: `../daml-escrow` (ledger/settlement),
`../daml-escrow-cms` (contract authoring/CLM), `../daml-escrow-identity` (shared T1 identity directory), and this
one (`daml-escrow-commons`, tools/utilities/validators the others depend on).

- Before adding anything here, read `.claude/skills/commons-contribution/SKILL.md` — it has the promotion
  checklist (used by 2+ repos, no domain logic, ships with tests) that keeps this module from becoming a dumping
  ground.
- `../daml-escrow/CLAUDE.md`, `../daml-escrow-cms/CLAUDE.md`, and `../daml-escrow-identity/CLAUDE.md` are the
  domain authorities for their respective repos. If a question is about escrow lifecycle, ledger state, CLM/
  drafting flows, or identity resolution, look there — not here.
- Source-control policy (branching, commit signing, PR/review requirements, branch protection strategy) is
  mirrored across all four repos — see `.gemini/repo_rules.md` here and `.agents/skills/github-ci/SKILL.md` for
  the operational mechanics. Keep it consistent with the other three repos' equivalents if either changes.
- If another repo's `graphify-out/` exists and you need to answer a question about that repo's code, prefer
  `graphify query "<question>"` there over grepping raw files. For a question spanning repos, merge graphs on
  demand: `graphify merge-graphs ../daml-escrow/graphify-out/graph.json ../daml-escrow-cms/graphify-out/graph.json ../daml-escrow-identity/graphify-out/graph.json graphify-out/graph.json --out <scratch-path>/cross-repo-graph.json`
  (generate into a scratch location, not committed — it goes stale as soon as any side's own graph updates).
