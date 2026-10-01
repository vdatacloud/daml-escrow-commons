# Graph Report - daml-escrow-commons  (2026-09-09)

## Corpus Check
- 20 files · ~8,333 words
- Verdict: corpus is large enough that graph structure adds value.

## Summary
- 144 nodes · 193 edges · 16 communities (13 shown, 3 thin omitted)
- Extraction: 87% EXTRACTED · 13% INFERRED · 0% AMBIGUOUS · INFERRED: 25 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `01eecc15`
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
- Client
- testing.T

## God Nodes (most connected - your core abstractions)
1. `GitHub & CI Workflow Skill` - 12 edges
2. `Client` - 9 edges
3. `New()` - 8 edges
4. `Repository Guardrails` - 8 edges
5. `LoadDirectory()` - 7 edges
6. `Releasing` - 7 edges
7. `Identity` - 6 edges
8. `RequireNonEmpty()` - 6 edges
9. `Verify()` - 5 edges
10. `SettlementEvent` - 5 edges

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

## Communities (16 total, 3 thin omitted)

### Community 0 - "GitHub & CI Workflow Skill"
Cohesion: 0.11
Nodes (17): 10. Useful Reference Commands, 11. References, 1. Pre-Commit Local Verification (MANDATORY — do this before every commit), 2. Branching Rules, 3. Commit Standards, 4. Staging & Pushing Changes, 5. Pull Request Creation, 6. CI Pipeline Overview (+9 more)

### Community 1 - "RequireNonEmpty"
Cohesion: 0.18
Nodes (11): Errors, RequireNonEmpty(), RequireOneOf(), RequirePositive(), RequireValidEmail(), TestErrors_AggregatesAndReports(), TestErrors_ErrIfAny_NilWhenEmpty(), TestRequireNonEmpty() (+3 more)

### Community 2 - "widget.json"
Cohesion: 0.17
Nodes (11): minLength, type, properties, name, quantity, minimum, type, required (+3 more)

### Community 3 - "Repository Guardrails"
Cohesion: 0.22
Nodes (8): Branch Protection Strategy, Branching Strategy, CI Requirements, Code Review Requirements, Commit Standard, Pre-Commit Verification, Pull Request Rules, Repository Guardrails

### Community 4 - "Verify"
Cohesion: 0.46
Nodes (6): Sign(), TestSignAndVerify_RoundTrip(), TestVerify_MalformedHexFails(), TestVerify_TamperedMessageFails(), TestVerify_WrongSecretFails(), Verify()

### Community 5 - "LoadDirectory"
Cohesion: 0.17
Nodes (10): github.com/xeipuuv/gojsonschema.Schema, ErrUnknownType, Registry, LoadDirectory(), TestLoadDirectory_CompilesSchemas(), TestLoadDirectory_MissingDirectory(), TestValidate_InvalidPayloadReportsFailures(), TestValidate_UnknownType() (+2 more)

### Community 6 - "Releasing"
Cohesion: 0.25
Nodes (7): 1. Prerequisites (one-time, per machine that will `go get` this module), 2. Decide the version bump, 3. Tag and push, 4. Publish the GitHub release, 5. Update consumers, Future automation, Releasing

### Community 7 - "SettlementEvent"
Cohesion: 0.43
Nodes (5): time.Time, ChargeBearer, LedgerCommandEvent, Rail, SettlementEvent

### Community 8 - "CLAUDE.md"
Cohesion: 0.29
Nodes (5): Architecture, Commands, Conventions carried over from `daml-escrow` / `daml-escrow-cms`, Project Overview, Relationship to the other repos

### Community 9 - "/commons-contribution"
Cohesion: 0.33
Nodes (5): /commons-contribution, How to add something, Removing something, When NOT to add something, When to promote code into commons

### Community 10 - "daml-escrow-commons"
Cohesion: 0.33
Nodes (5): daml-escrow-commons, Status, Using this module, What's here, What's NOT here, and won't be

### Community 15 - "Client"
Cohesion: 0.36
Nodes (5): context.Context, net/http.Client, net/http.Request, Client, Identity

### Community 16 - "testing.T"
Cohesion: 0.33
Nodes (10): testing.T, New(), TestClient_GetByEmail_ServerError(), TestClient_GetByOktaSub_Found(), TestClient_GetByToken_NotFound(), TestClient_ManagesIdentity(), TestClient_ManagesIdentity_ServerError(), TestClient_Upsert_Success() (+2 more)

## Knowledge Gaps
- **52 isolated node(s):** `github.com/vdatacloud/daml-escrow-commons`, `$schema`, `title`, `type`, `type` (+47 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 69 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **3 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `New()` connect `testing.T` to `Client`?**
  _High betweenness centrality (0.070) - this node is a cross-community bridge._
- **Why does `Client` connect `Client` to `testing.T`?**
  _High betweenness centrality (0.055) - this node is a cross-community bridge._
- **Are the 6 inferred relationships involving `New()` (e.g. with `TestClient_GetByEmail_ServerError()` and `TestClient_GetByOktaSub_Found()`) actually correct?**
  _`New()` has 6 INFERRED edges - model-reasoned connections that need verification._
- **What connects `github.com/vdatacloud/daml-escrow-commons`, `$schema`, `title` to the rest of the system?**
  _52 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `GitHub & CI Workflow Skill` be split into smaller, more focused modules?**
  _Cohesion score 0.1111111111111111 - nodes in this community are weakly interconnected._