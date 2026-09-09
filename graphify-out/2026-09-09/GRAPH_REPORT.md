# Graph Report - daml-escrow-commons  (2026-08-21)

## Corpus Check
- 20 files · ~8,224 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 149 nodes · 196 edges · 17 communities (14 shown, 3 thin omitted)
- Extraction: 84% EXTRACTED · 16% INFERRED · 0% AMBIGUOUS · INFERRED: 31 edges (avg confidence: 0.8)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `fef5b3ae`
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
- Client
- New

## God Nodes (most connected - your core abstractions)
1. `GitHub & CI Workflow Skill` - 12 edges
2. `New()` - 10 edges
3. `Client` - 8 edges
4. `Repository Guardrails` - 8 edges
5. `LoadDirectory()` - 7 edges
6. `Releasing` - 7 edges
7. `Verify()` - 6 edges
8. `Identity` - 6 edges
9. `RequireNonEmpty()` - 6 edges
10. `Sign()` - 5 edges

## Surprising Connections (you probably didn't know these)
- `Sign()` --calls--> `New()`  [INFERRED]
  hmacsig/hmacsig.go → identityclient/identityclient.go
- `Verify()` --calls--> `New()`  [INFERRED]
  hmacsig/hmacsig.go → identityclient/identityclient.go
- `TestSignAndVerify_RoundTrip()` --calls--> `Sign()`  [INFERRED]
  hmacsig/hmacsig_test.go → hmacsig/hmacsig.go
- `TestVerify_TamperedMessageFails()` --calls--> `Sign()`  [INFERRED]
  hmacsig/hmacsig_test.go → hmacsig/hmacsig.go
- `TestVerify_WrongSecretFails()` --calls--> `Sign()`  [INFERRED]
  hmacsig/hmacsig_test.go → hmacsig/hmacsig.go

## Import Cycles
- None detected.

## Communities (17 total, 3 thin omitted)

### Community 0 - "GitHub & CI Workflow Skill"
Cohesion: 0.11
Nodes (17): 10. Useful Reference Commands, 11. References, 1. Pre-Commit Local Verification (MANDATORY — do this before every commit), 2. Branching Rules, 3. Commit Standards, 4. Staging & Pushing Changes, 5. Pull Request Creation, 6. CI Pipeline Overview (+9 more)

### Community 1 - "RequireNonEmpty"
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

### Community 15 - "Client"
Cohesion: 0.42
Nodes (4): Context, Client, Identity, Request

### Community 16 - "New"
Cohesion: 0.50
Nodes (8): New(), T, TestClient_GetByEmail_ServerError(), TestClient_GetByOktaSub_Found(), TestClient_GetByToken_NotFound(), TestClient_ManagesIdentity(), TestClient_ManagesIdentity_ServerError(), TestClient_Upsert_Success()

## Knowledge Gaps
- **53 isolated node(s):** `github.com/vdatacloud/daml-escrow-commons`, `$schema`, `title`, `type`, `type` (+48 more)
  These have ≤1 connection - possible missing edges or undocumented components.
- **3 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `New()` connect `New` to `Verify`, `Client`?**
  _High betweenness centrality (0.024) - this node is a cross-community bridge._
- **Why does `Client` connect `Client` to `New`?**
  _High betweenness centrality (0.014) - this node is a cross-community bridge._
- **Are the 8 inferred relationships involving `New()` (e.g. with `Sign()` and `Verify()`) actually correct?**
  _`New()` has 8 INFERRED edges - model-reasoned connections that need verification._
- **Are the 5 inferred relationships involving `LoadDirectory()` (e.g. with `TestLoadDirectory_CompilesSchemas()` and `TestLoadDirectory_MissingDirectory()`) actually correct?**
  _`LoadDirectory()` has 5 INFERRED edges - model-reasoned connections that need verification._
- **What connects `github.com/vdatacloud/daml-escrow-commons`, `$schema`, `title` to the rest of the system?**
  _53 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `GitHub & CI Workflow Skill` be split into smaller, more focused modules?**
  _Cohesion score 0.1111111111111111 - nodes in this community are weakly interconnected._
- **Should `widget.json` be split into smaller, more focused modules?**
  _Cohesion score 0.14285714285714285 - nodes in this community are weakly interconnected._