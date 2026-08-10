# Graph Report - daml-escrow-commons  (2026-08-10)

## Corpus Check
- 16 files · ~5,936 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 117 nodes · 133 edges · 14 communities (11 shown, 3 thin omitted)
- Extraction: 86% EXTRACTED · 14% INFERRED · 0% AMBIGUOUS · INFERRED: 19 edges (avg confidence: 0.8)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `f5317de6`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- GitHub & CI Workflow Skill
- validate_test.go
- widget.json
- Repository Guardrails
- Verify
- Registry
- Releasing
- LoadDirectory
- CLAUDE.md
- /commons-contribution
- daml-escrow-commons
- .claude/CLAUDE.md
- install-git-hooks.sh
- github.com/vdatacloud/daml-escrow-commons

## God Nodes (most connected - your core abstractions)
1. `GitHub & CI Workflow Skill` - 12 edges
2. `Repository Guardrails` - 8 edges
3. `LoadDirectory()` - 7 edges
4. `Releasing` - 7 edges
5. `Verify()` - 5 edges
6. `Registry` - 5 edges
7. `/commons-contribution` - 5 edges
8. `daml-escrow-commons` - 5 edges
9. `Sign()` - 4 edges
10. `TestSignAndVerify_RoundTrip()` - 4 edges

## Surprising Connections (you probably didn't know these)
- `TestSignAndVerify_RoundTrip()` --calls--> `Sign()`  [INFERRED]
  hmacsig/hmacsig_test.go → hmacsig/hmacsig.go
- `TestVerify_TamperedMessageFails()` --calls--> `Sign()`  [INFERRED]
  hmacsig/hmacsig_test.go → hmacsig/hmacsig.go
- `TestVerify_WrongSecretFails()` --calls--> `Sign()`  [INFERRED]
  hmacsig/hmacsig_test.go → hmacsig/hmacsig.go
- `TestSignAndVerify_RoundTrip()` --calls--> `Verify()`  [INFERRED]
  hmacsig/hmacsig_test.go → hmacsig/hmacsig.go
- `TestVerify_MalformedHexFails()` --calls--> `Verify()`  [INFERRED]
  hmacsig/hmacsig_test.go → hmacsig/hmacsig.go

## Import Cycles
- None detected.

## Communities (14 total, 3 thin omitted)

### Community 0 - "GitHub & CI Workflow Skill"
Cohesion: 0.11
Nodes (17): 10. Useful Reference Commands, 11. References, 1. Pre-Commit Local Verification (MANDATORY — do this before every commit), 2. Branching Rules, 3. Commit Standards, 4. Staging & Pushing Changes, 5. Pull Request Creation, 6. CI Pipeline Overview (+9 more)

### Community 1 - "validate_test.go"
Cohesion: 0.20
Nodes (12): Errors, RequireNonEmpty(), RequireOneOf(), RequirePositive(), RequireValidEmail(), T, TestErrors_AggregatesAndReports(), TestErrors_ErrIfAny_NilWhenEmpty() (+4 more)

### Community 2 - "widget.json"
Cohesion: 0.14
Nodes (13): name, quantity, minLength, type, properties, name, quantity, minimum (+5 more)

### Community 3 - "Repository Guardrails"
Cohesion: 0.22
Nodes (8): Branch Protection Strategy, Branching Strategy, CI Requirements, Code Review Requirements, Commit Standard, Pre-Commit Verification, Pull Request Rules, Repository Guardrails

### Community 4 - "Verify"
Cohesion: 0.47
Nodes (7): Sign(), T, TestSignAndVerify_RoundTrip(), TestVerify_MalformedHexFails(), TestVerify_TamperedMessageFails(), TestVerify_WrongSecretFails(), Verify()

### Community 5 - "Registry"
Cohesion: 0.22
Nodes (4): Schema, ErrUnknownType, Registry, ValidationError

### Community 6 - "Releasing"
Cohesion: 0.25
Nodes (7): 1. Prerequisites (one-time, per machine that will `go get` this module), 2. Decide the version bump, 3. Tag and push, 4. Publish the GitHub release, 5. Update consumers, Future automation, Releasing

### Community 7 - "LoadDirectory"
Cohesion: 0.54
Nodes (7): LoadDirectory(), T, TestLoadDirectory_CompilesSchemas(), TestLoadDirectory_MissingDirectory(), TestValidate_InvalidPayloadReportsFailures(), TestValidate_UnknownType(), TestValidate_ValidPayload()

### Community 8 - "CLAUDE.md"
Cohesion: 0.29
Nodes (5): Architecture, Commands, Conventions carried over from `daml-escrow` / `daml-escrow-cms`, Project Overview, Relationship to the other repos

### Community 9 - "/commons-contribution"
Cohesion: 0.33
Nodes (5): /commons-contribution, How to add something, Removing something, When NOT to add something, When to promote code into commons

### Community 10 - "daml-escrow-commons"
Cohesion: 0.33
Nodes (5): daml-escrow-commons, Status, Using this module, What's here, What's NOT here, and won't be

## Knowledge Gaps
- **53 isolated node(s):** `github.com/vdatacloud/daml-escrow-commons`, `$schema`, `title`, `type`, `type` (+48 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **3 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `LoadDirectory()` connect `LoadDirectory` to `Registry`?**
  _High betweenness centrality (0.010) - this node is a cross-community bridge._
- **Are the 5 inferred relationships involving `LoadDirectory()` (e.g. with `TestLoadDirectory_CompilesSchemas()` and `TestLoadDirectory_MissingDirectory()`) actually correct?**
  _`LoadDirectory()` has 5 INFERRED edges - model-reasoned connections that need verification._
- **What connects `github.com/vdatacloud/daml-escrow-commons`, `$schema`, `title` to the rest of the system?**
  _53 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `GitHub & CI Workflow Skill` be split into smaller, more focused modules?**
  _Cohesion score 0.1111111111111111 - nodes in this community are weakly interconnected._
- **Should `widget.json` be split into smaller, more focused modules?**
  _Cohesion score 0.14285714285714285 - nodes in this community are weakly interconnected._