# Graph Report - daml-escrow-commons  (2026-08-18)

## Corpus Check
- 18 files · ~6,850 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 129 nodes · 151 edges · 16 communities (12 shown, 4 thin omitted)
- Extraction: 85% EXTRACTED · 15% INFERRED · 0% AMBIGUOUS · INFERRED: 23 edges (avg confidence: 0.8)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `28a0427a`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- GitHub & CI Workflow Skill
- RequireNonEmpty
- widget.json
- Repository Guardrails
- Verify
- LoadDirectory
- Releasing
- SettlementEvent
- CLAUDE.md
- /commons-contribution
- daml-escrow-commons
- .claude/CLAUDE.md
- install-git-hooks.sh
- github.com/vdatacloud/daml-escrow-commons
- metering_test.go
- Errors

## God Nodes (most connected - your core abstractions)
1. `GitHub & CI Workflow Skill` - 12 edges
2. `Repository Guardrails` - 8 edges
3. `LoadDirectory()` - 7 edges
4. `Releasing` - 7 edges
5. `RequireNonEmpty()` - 6 edges
6. `Verify()` - 5 edges
7. `SettlementEvent` - 5 edges
8. `Registry` - 5 edges
9. `/commons-contribution` - 5 edges
10. `daml-escrow-commons` - 5 edges

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

## Communities (16 total, 4 thin omitted)

### Community 0 - "GitHub & CI Workflow Skill"
Cohesion: 0.11
Nodes (17): 10. Useful Reference Commands, 11. References, 1. Pre-Commit Local Verification (MANDATORY — do this before every commit), 2. Branching Rules, 3. Commit Standards, 4. Staging & Pushing Changes, 5. Pull Request Creation, 6. CI Pipeline Overview (+9 more)

### Community 1 - "RequireNonEmpty"
Cohesion: 0.29
Nodes (11): RequireNonEmpty(), RequireOneOf(), RequirePositive(), RequireValidEmail(), T, TestErrors_AggregatesAndReports(), TestErrors_ErrIfAny_NilWhenEmpty(), TestRequireNonEmpty() (+3 more)

### Community 2 - "widget.json"
Cohesion: 0.14
Nodes (13): name, quantity, minLength, type, properties, name, quantity, minimum (+5 more)

### Community 3 - "Repository Guardrails"
Cohesion: 0.22
Nodes (8): Branch Protection Strategy, Branching Strategy, CI Requirements, Code Review Requirements, Commit Standard, Pre-Commit Verification, Pull Request Rules, Repository Guardrails

### Community 4 - "Verify"
Cohesion: 0.47
Nodes (7): Sign(), T, TestSignAndVerify_RoundTrip(), TestVerify_MalformedHexFails(), TestVerify_TamperedMessageFails(), TestVerify_WrongSecretFails(), Verify()

### Community 5 - "LoadDirectory"
Cohesion: 0.18
Nodes (11): Schema, ErrUnknownType, Registry, LoadDirectory(), T, TestLoadDirectory_CompilesSchemas(), TestLoadDirectory_MissingDirectory(), TestValidate_InvalidPayloadReportsFailures() (+3 more)

### Community 6 - "Releasing"
Cohesion: 0.25
Nodes (7): 1. Prerequisites (one-time, per machine that will `go get` this module), 2. Decide the version bump, 3. Tag and push, 4. Publish the GitHub release, 5. Update consumers, Future automation, Releasing

### Community 7 - "SettlementEvent"
Cohesion: 0.43
Nodes (5): ChargeBearer, LedgerCommandEvent, Rail, SettlementEvent, Time

### Community 8 - "CLAUDE.md"
Cohesion: 0.29
Nodes (5): Architecture, Commands, Conventions carried over from `daml-escrow` / `daml-escrow-cms`, Project Overview, Relationship to the other repos

### Community 9 - "/commons-contribution"
Cohesion: 0.33
Nodes (5): /commons-contribution, How to add something, Removing something, When NOT to add something, When to promote code into commons

### Community 10 - "daml-escrow-commons"
Cohesion: 0.33
Nodes (5): daml-escrow-commons, Status, Using this module, What's here, What's NOT here, and won't be

### Community 14 - "metering_test.go"
Cohesion: 0.67
Nodes (3): T, TestLedgerCommandEvent_Validate(), TestSettlementEvent_Validate()

## Knowledge Gaps
- **53 isolated node(s):** `github.com/vdatacloud/daml-escrow-commons`, `$schema`, `title`, `type`, `type` (+48 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **4 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `RequireNonEmpty()` connect `RequireNonEmpty` to `SettlementEvent`?**
  _High betweenness centrality (0.011) - this node is a cross-community bridge._
- **Why does `SettlementEvent` connect `SettlementEvent` to `RequireNonEmpty`?**
  _High betweenness centrality (0.008) - this node is a cross-community bridge._
- **Are the 5 inferred relationships involving `LoadDirectory()` (e.g. with `TestLoadDirectory_CompilesSchemas()` and `TestLoadDirectory_MissingDirectory()`) actually correct?**
  _`LoadDirectory()` has 5 INFERRED edges - model-reasoned connections that need verification._
- **Are the 5 inferred relationships involving `RequireNonEmpty()` (e.g. with `.Validate()` and `.Validate()`) actually correct?**
  _`RequireNonEmpty()` has 5 INFERRED edges - model-reasoned connections that need verification._
- **What connects `github.com/vdatacloud/daml-escrow-commons`, `$schema`, `title` to the rest of the system?**
  _53 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `GitHub & CI Workflow Skill` be split into smaller, more focused modules?**
  _Cohesion score 0.1111111111111111 - nodes in this community are weakly interconnected._
- **Should `widget.json` be split into smaller, more focused modules?**
  _Cohesion score 0.14285714285714285 - nodes in this community are weakly interconnected._