# Graph Report - daml-escrow-commons  (2026-10-05)

## Corpus Check
- 64 files · ~21,055 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 3 file(s) not represented in the graph (top: (none) 3)

## Summary
- 319 nodes · 630 edges · 20 communities (15 shown, 5 thin omitted)
- Extraction: 89% EXTRACTED · 11% INFERRED · 0% AMBIGUOUS · INFERRED: 72 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `1cb9a655`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- GitHub & CI Workflow Skill
- testing.T
- widget.json
- Repository Guardrails
- hmacsig.go
- cantonid.go
- Releasing
- metering.go
- CLAUDE.md
- /commons-contribution
- daml-escrow-commons
- .claude/CLAUDE.md
- install-git-hooks.sh
- github.com/vdatacloud/daml-escrow-commons
- http.go
- Client
- apierror_test.go
- next-version.sh
- test-next-version.sh
- apierror.go

## God Nodes (most connected - your core abstractions)
1. `Fingerprint` - 13 edges
2. `PartyID` - 12 edges
3. `GitHub & CI Workflow Skill` - 12 edges
4. `New()` - 11 edges
5. `TestSigningDetails()` - 11 edges
6. `Hash` - 11 edges
7. `RequireNonEmpty()` - 11 edges
8. `FromLedgerError()` - 10 edges
9. `Releasing` - 10 edges
10. `ParsePartyID()` - 9 edges

## Surprising Connections (you probably didn't know these)
- `init()` --calls--> `New()`  [INFERRED]
  apierror/apierror_test.go → apierror/apierror.go
- `TestWriter()` --calls--> `New()`  [INFERRED]
  apierror/apierror_test.go → apierror/apierror.go
- `TestLogFields_Golden()` --calls--> `encode()`  [INFERRED]
  apierror/apierror_test.go → apierror/apierror.go
- `TestCanonical_Golden()` --calls--> `Parse()`  [INFERRED]
  apierror/apierror_test.go → apierror/apierror.go
- `FromResponse()` --calls--> `Parse()`  [INFERRED]
  apierror/http.go → apierror/apierror.go

## Import Cycles
- None detected.

## Communities (20 total, 5 thin omitted)

### Community 0 - "GitHub & CI Workflow Skill"
Cohesion: 0.11
Nodes (17): 10. Useful Reference Commands, 11. References, 1. Pre-Commit Local Verification (MANDATORY — do this before every commit), 2. Branching Rules, 3. Commit Standards, 4. Staging & Pushing Changes, 5. Pull Request Creation, 6. CI Pipeline Overview (+9 more)

### Community 1 - "testing.T"
Cohesion: 0.11
Nodes (24): TestCanonical_Properties(), TestLogFields_ShortFormsOnly(), New(), TestClient_GetByEmail_ServerError(), TestClient_GetByOktaSub_Found(), TestClient_GetByToken_NotFound(), TestClient_ManagesIdentity(), TestClient_ManagesIdentity_ServerError() (+16 more)

### Community 2 - "widget.json"
Cohesion: 0.17
Nodes (11): minLength, type, properties, name, quantity, minimum, type, required (+3 more)

### Community 3 - "Repository Guardrails"
Cohesion: 0.22
Nodes (8): Branch Protection Strategy, Branching Strategy, CI Requirements, Code Review Requirements, Commit Standard, Pre-Commit Verification, Pull Request Rules, Repository Guardrails

### Community 4 - "hmacsig.go"
Cohesion: 0.23
Nodes (8): Sign(), TestSignAndVerify_RoundTrip(), TestVerify_MalformedHexFails(), TestVerify_TamperedMessageFails(), TestVerify_WrongSecretFails(), Verify(), Ref(), TestRef()

### Community 5 - "cantonid.go"
Cohesion: 0.09
Nodes (24): ExpectedHash(), ForParty(), Signing(), TestSigningDetails(), SigningDetails, Classify(), FingerprintOf(), Fingerprint (+16 more)

### Community 6 - "Releasing"
Cohesion: 0.18
Nodes (10): 1. Prerequisites (one-time, per machine that will `go get` this module), 2. Decide the version bump, 3. Tag and push, 4. Publish the GitHub release, 5. Update consumers, Automated (default), Future automation, Major releases (+2 more)

### Community 7 - "metering.go"
Cohesion: 0.11
Nodes (20): ChargeBearer, CommandOutcome, LedgerCommandEvent, requireBaseUnits(), requireDecimal(), NetworkFeeEvent, NetworkFeePayer, Rail (+12 more)

### Community 8 - "CLAUDE.md"
Cohesion: 0.29
Nodes (5): Architecture, Commands, Conventions carried over from `daml-escrow` / `daml-escrow-cms`, Project Overview, Relationship to the other repos

### Community 9 - "/commons-contribution"
Cohesion: 0.33
Nodes (5): /commons-contribution, How to add something, Removing something, When NOT to add something, When to promote code into commons

### Community 10 - "daml-escrow-commons"
Cohesion: 0.33
Nodes (5): daml-escrow-commons, Status, Using this module, What's here, What's NOT here, and won't be

### Community 14 - "http.go"
Cohesion: 0.17
Nodes (9): TestFromResponse(), TestWriter(), FromResponse(), newRequestID(), RequestID(), RequestIDFrom(), Write(), requestIDKey (+1 more)

### Community 16 - "apierror_test.go"
Cohesion: 0.09
Nodes (21): golden(), TestCanonical_Golden(), TestLogFields_Golden(), TestSchemaIsACopy(), TestShortID(), validate(), FromError(), FromLedgerError() (+13 more)

### Community 18 - "test-next-version.sh"
Cohesion: 0.70
Nodes (4): c(), expect(), test-next-version.sh script, tg()

### Community 21 - "apierror.go"
Cohesion: 0.09
Nodes (22): CodeForStatus(), encode(), Error, New(), Newf(), Parse(), sortedJSON(), init() (+14 more)

## Knowledge Gaps
- **59 isolated node(s):** `signing`, `requestIDKey`, `Error`, `github.com/vdatacloud/daml-escrow-commons`, `$schema` (+54 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 96 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **5 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Error` connect `apierror.go` to `apierror_test.go`, `cantonid.go`?**
  _High betweenness centrality (0.040) - this node is a cross-community bridge._
- **Are the 3 inferred relationships involving `Fingerprint` (e.g. with `TestJSON()` and `TestPartyID()`) actually correct?**
  _`Fingerprint` has 3 INFERRED edges - model-reasoned connections that need verification._
- **What connects `signing`, `requestIDKey`, `Error` to the rest of the system?**
  _59 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `GitHub & CI Workflow Skill` be split into smaller, more focused modules?**
  _Cohesion score 0.1111111111111111 - nodes in this community are weakly interconnected._
- **Why does `Client` connect `Client` to `testing.T`?**
  _High betweenness centrality (0.022) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.11491935483870967 - nodes in this community are weakly interconnected._
- **Should `cantonid.go` be split into smaller, more focused modules?**
  _Cohesion score 0.09059233449477352 - nodes in this community are weakly interconnected._