# Cross-repo context

This repo is the shared-utilities leg of a three-repo system: `../daml-escrow` (ledger/settlement),
`../daml-escrow-cms` (contract authoring/CLM), and this one (`daml-escrow-commons`, tools/utilities/validators
both depend on).

- Before adding anything here, read `.claude/skills/commons-contribution/SKILL.md` — it has the promotion
  checklist (used by 2+ repos, no domain logic, ships with tests) that keeps this module from becoming a dumping
  ground.
- `../daml-escrow/CLAUDE.md` and `../daml-escrow-cms/CLAUDE.md` are the domain authorities for their respective
  repos. If a question is about escrow lifecycle, ledger state, or CLM/drafting flows, look there — not here.
- If either consumer repo's `graphify-out/` exists and you need to answer a question about that repo's code,
  prefer `graphify query "<question>"` there over grepping raw files.
