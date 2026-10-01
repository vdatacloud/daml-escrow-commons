# Graph Report - daml-escrow-commons  (2026-10-01)

## Corpus Check
- 22 files · ~10,341 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 3 file(s) not represented in the graph (top: (none) 3)

## Summary
- 178 nodes · 260 edges · 19 communities (14 shown, 5 thin omitted)
- Extraction: 90% EXTRACTED · 10% INFERRED · 0% AMBIGUOUS · INFERRED: 25 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `cdd653ee`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- GitHub & CI Workflow Skill
- RequireNonEmpty
- widget.json
- Repository Guardrails
- hmacsig.go
- LoadDirectory
- Releasing
- metering.go
- CLAUDE.md
- /commons-contribution
- daml-escrow-commons
- .claude/CLAUDE.md
- install-git-hooks.sh
- github.com/vdatacloud/daml-escrow-commons
- schema.go
- Client
- testing.T
- next-version.sh
- test-next-version.sh

## God Nodes (most connected - your core abstractions)
1. `GitHub & CI Workflow Skill` - 12 edges
2. `Client` - 9 edges
3. `Releasing` - 9 edges
4. `New()` - 8 edges
5. `SettlementEvent` - 8 edges
6. `RequireNonEmpty()` - 8 edges
7. `Repository Guardrails` - 8 edges
8. `LoadDirectory()` - 7 edges
9. `Identity` - 6 edges
10. `Verify()` - 5 edges

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

## Communities (19 total, 5 thin omitted)

### Community 0 - "GitHub & CI Workflow Skill"
Cohesion: 0.11
Nodes (17): 10. Useful Reference Commands, 11. References, 1. Pre-Commit Local Verification (MANDATORY — do this before every commit), 2. Branching Rules, 3. Commit Standards, 4. Staging & Pushing Changes, 5. Pull Request Creation, 6. CI Pipeline Overview (+9 more)

### Community 1 - "RequireNonEmpty"
Cohesion: 0.24
Nodes (9): requireBaseUnits(), RequireNonEmpty(), RequireOneOf(), RequirePositive(), TestErrors_AggregatesAndReports(), TestErrors_ErrIfAny_NilWhenEmpty(), TestRequireNonEmpty(), TestRequireOneOf() (+1 more)

### Community 2 - "widget.json"
Cohesion: 0.17
Nodes (11): minLength, type, properties, name, quantity, minimum, type, required (+3 more)

### Community 3 - "Repository Guardrails"
Cohesion: 0.22
Nodes (8): Branch Protection Strategy, Branching Strategy, CI Requirements, Code Review Requirements, Commit Standard, Pre-Commit Verification, Pull Request Rules, Repository Guardrails

### Community 4 - "hmacsig.go"
Cohesion: 0.27
Nodes (6): Sign(), TestSignAndVerify_RoundTrip(), TestVerify_MalformedHexFails(), TestVerify_TamperedMessageFails(), TestVerify_WrongSecretFails(), Verify()

### Community 5 - "LoadDirectory"
Cohesion: 0.25
Nodes (7): Registry, LoadDirectory(), TestLoadDirectory_CompilesSchemas(), TestLoadDirectory_MissingDirectory(), TestValidate_InvalidPayloadReportsFailures(), TestValidate_UnknownType(), TestValidate_ValidPayload()

### Community 6 - "Releasing"
Cohesion: 0.20
Nodes (9): 1. Prerequisites (one-time, per machine that will `go get` this module), 2. Decide the version bump, 3. Tag and push, 4. Publish the GitHub release, 5. Update consumers, Automated (default), Future automation, Manual It follows [Semantic Versioning](https://semver.org/) and (+1 more)

### Community 7 - "metering.go"
Cohesion: 0.21
Nodes (11): ChargeBearer, LedgerCommandEvent, TestLedgerCommandEvent_Validate(), TestNetworkFeeEvent_Validate(), TestSettlementEvent_NetworkFee(), TestSettlementEvent_Validate(), withFee(), NetworkFeeEvent (+3 more)

### Community 8 - "CLAUDE.md"
Cohesion: 0.29
Nodes (5): Architecture, Commands, Conventions carried over from `daml-escrow` / `daml-escrow-cms`, Project Overview, Relationship to the other repos

### Community 9 - "/commons-contribution"
Cohesion: 0.33
Nodes (5): /commons-contribution, How to add something, Removing something, When NOT to add something, When to promote code into commons

### Community 10 - "daml-escrow-commons"
Cohesion: 0.33
Nodes (5): daml-escrow-commons, Status, Using this module, What's here, What's NOT here, and won't be

### Community 14 - "schema.go"
Cohesion: 0.12
Nodes (5): ErrUnknownType, ValidationError, Errors, RequireValidEmail(), TestRequireValidEmail()

### Community 16 - "testing.T"
Cohesion: 0.23
Nodes (7): New(), TestClient_GetByEmail_ServerError(), TestClient_GetByOktaSub_Found(), TestClient_GetByToken_NotFound(), TestClient_ManagesIdentity(), TestClient_ManagesIdentity_ServerError(), TestClient_Upsert_Success()

### Community 18 - "test-next-version.sh"
Cohesion: 0.70
Nodes (4): c(), expect(), test-next-version.sh script, tg()

## Knowledge Gaps
- **55 isolated node(s):** `github.com/vdatacloud/daml-escrow-commons`, `$schema`, `title`, `type`, `type` (+50 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 82 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **5 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Client` connect `Client` to `testing.T`?**
  _High betweenness centrality (0.040) - this node is a cross-community bridge._
- **Why does `New()` connect `testing.T` to `Client`?**
  _High betweenness centrality (0.034) - this node is a cross-community bridge._
- **Why does `LoadDirectory()` connect `LoadDirectory` to `schema.go`?**
  _High betweenness centrality (0.030) - this node is a cross-community bridge._
- **Are the 6 inferred relationships involving `New()` (e.g. with `TestClient_GetByEmail_ServerError()` and `TestClient_GetByOktaSub_Found()`) actually correct?**
  _`New()` has 6 INFERRED edges - model-reasoned connections that need verification._
- **What connects `github.com/vdatacloud/daml-escrow-commons`, `$schema`, `title` to the rest of the system?**
  _55 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `GitHub & CI Workflow Skill` be split into smaller, more focused modules?**
  _Cohesion score 0.1111111111111111 - nodes in this community are weakly interconnected._
- **Should `schema.go` be split into smaller, more focused modules?**
  _Cohesion score 0.125 - nodes in this community are weakly interconnected._