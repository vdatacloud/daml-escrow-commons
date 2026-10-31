# Graph Report - daml-escrow-commons  (2026-10-28)

## Corpus Check
- 87 files · ~35,149 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 9 file(s) not represented in the graph (top: .ttf 5, (none) 4)

## Summary
- 476 nodes · 1043 edges · 28 communities (21 shown, 7 thin omitted)
- Extraction: 89% EXTRACTED · 11% INFERRED · 0% AMBIGUOUS · INFERRED: 111 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `bf98ceb3`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- GitHub & CI Workflow Skill
- testing.T
- widget.json
- Repository Guardrails
- go_pkg_errors
- cantonid.go
- Releasing
- metering.go
- CLAUDE.md
- /commons-contribution
- daml-escrow-commons
- .claude/CLAUDE.md
- install-git-hooks.sh
- github.com/vdatacloud/daml-escrow-commons
- agreementsig_test.go
- render.go
- apierror_test.go
- next-version.sh
- test-next-version.sh
- agreement.go
- field.go
- apierror.go
- http.go
- approval_test.go
- context.Context
- Registry

## God Nodes (most connected - your core abstractions)
1. `Fingerprint` - 13 edges
2. `PartyID` - 12 edges
3. `GitHub & CI Workflow Skill` - 12 edges
4. `AgreementVersion` - 11 edges
5. `New()` - 11 edges
6. `TestSigningDetails()` - 11 edges
7. `Hash` - 11 edges
8. `RequireNonEmpty()` - 11 edges
9. `signatureCases()` - 10 edges
10. `FromLedgerError()` - 10 edges

## Surprising Connections (you probably didn't know these)
- `TestAgreementVectors()` --calls--> `AgreementCanonical()`  [INFERRED]
  agreementsig/agreementsig_test.go → agreementsig/agreement.go
- `TestAgreementMessage_Refusals()` --calls--> `AgreementMessage()`  [INFERRED]
  agreementsig/agreementsig_test.go → agreementsig/agreement.go
- `TestAgreementVectors()` --calls--> `AgreementMessage()`  [INFERRED]
  agreementsig/agreementsig_test.go → agreementsig/agreement.go
- `TestDraftVectors()` --calls--> `DraftMessage()`  [INFERRED]
  agreementsig/agreementsig_test.go → agreementsig/draft.go
- `DraftCanonical()` --calls--> `decodeJSON()`  [INFERRED]
  agreementsig/draft.go → agreementsig/canonical.go

## Import Cycles
- None detected.

## Communities (28 total, 7 thin omitted)

### Community 0 - "GitHub & CI Workflow Skill"
Cohesion: 0.11
Nodes (17): 10. Useful Reference Commands, 11. References, 1. Pre-Commit Local Verification (MANDATORY — do this before every commit), 2. Branching Rules, 3. Commit Standards, 4. Staging & Pushing Changes, 5. Pull Request Creation, 6. CI Pipeline Overview (+9 more)

### Community 1 - "testing.T"
Cohesion: 0.09
Nodes (36): golden(), TestCanonical_Golden(), TestCanonical_Properties(), TestLogFields_Golden(), TestLogFields_ShortFormsOnly(), TestWriter(), validate(), New() (+28 more)

### Community 2 - "widget.json"
Cohesion: 0.17
Nodes (11): minLength, type, properties, name, quantity, minimum, type, required (+3 more)

### Community 3 - "Repository Guardrails"
Cohesion: 0.22
Nodes (8): Branch Protection Strategy, Branching Strategy, CI Requirements, Code Review Requirements, Commit Standard, Pre-Commit Verification, Pull Request Rules, Repository Guardrails

### Community 4 - "go_pkg_errors"
Cohesion: 0.23
Nodes (8): Sign(), TestSignAndVerify_RoundTrip(), TestVerify_MalformedHexFails(), TestVerify_TamperedMessageFails(), TestVerify_WrongSecretFails(), Verify(), Ref(), TestRef()

### Community 5 - "cantonid.go"
Cohesion: 0.10
Nodes (23): ExpectedHash(), ForParty(), Signing(), TestSigningDetails(), SigningDetails, Classify(), FingerprintOf(), Fingerprint (+15 more)

### Community 6 - "Releasing"
Cohesion: 0.18
Nodes (10): 1. Prerequisites (one-time, per machine that will `go get` this module), 2. Decide the version bump, 3. Tag and push, 4. Publish the GitHub release, 5. Update consumers, Automated (default), Future automation, Major releases (+2 more)

### Community 7 - "metering.go"
Cohesion: 0.14
Nodes (20): PendingInvitation, ChargeBearer, CommandOutcome, LedgerCommandEvent, requireBaseUnits(), requireDecimal(), NetworkFeeEvent, NetworkFeePayer (+12 more)

### Community 8 - "CLAUDE.md"
Cohesion: 0.29
Nodes (5): Architecture, Commands, Conventions carried over from `daml-escrow` / `daml-escrow-cms`, Project Overview, Relationship to the other repos

### Community 9 - "/commons-contribution"
Cohesion: 0.33
Nodes (5): /commons-contribution, How to add something, Removing something, When NOT to add something, When to promote code into commons

### Community 10 - "daml-escrow-commons"
Cohesion: 0.33
Nodes (5): daml-escrow-commons, Status, Using this module, What's here, What's NOT here, and won't be

### Community 14 - "agreementsig_test.go"
Cohesion: 0.09
Nodes (33): AgreementHash(), agreementCases(), agreementJSON(), flipS(), readJSON(), sha256Sum(), signatureCases(), spki() (+25 more)

### Community 15 - "render.go"
Cohesion: 0.08
Nodes (22): Artifact, checkName(), Deliver(), artifacts(), TestDeliver_ManyIsDeterministicZip(), TestDeliver_One(), TestDeliver_RefusesBadNames(), Zip() (+14 more)

### Community 16 - "apierror_test.go"
Cohesion: 0.08
Nodes (18): TestSchemaIsACopy(), TestShortID(), GRPCCodeName(), ParseLedgerError(), TestFromError(), TestFromLedgerError_Golden(), TestFromLedgerError_Mapping(), TestLedgerError() (+10 more)

### Community 18 - "test-next-version.sh"
Cohesion: 0.70
Nodes (4): c(), expect(), test-next-version.sh script, tg()

### Community 19 - "agreement.go"
Cohesion: 0.17
Nodes (14): AgreementCanonical(), AgreementMessage(), checkID(), invalid(), AgreementVersion, Amends, decodeJSON(), lessUTF16() (+6 more)

### Community 20 - "field.go"
Cohesion: 0.24
Nodes (10): ErrUnknownField, child(), contains(), Registry, lookupRef(), mergeProperty(), resolveNode(), splitPointer() (+2 more)

### Community 21 - "apierror.go"
Cohesion: 0.09
Nodes (24): draftVector, CodeForStatus(), encode(), Error, New(), Newf(), Parse(), sortedJSON() (+16 more)

### Community 23 - "http.go"
Cohesion: 0.18
Nodes (8): TestFromResponse(), FromResponse(), newRequestID(), RequestID(), RequestIDFrom(), Write(), requestIDKey, Writer

### Community 24 - "approval_test.go"
Cohesion: 0.26
Nodes (11): Check(), Decisions(), cases(), generate(), TestCheckRefusals(), TestDecisions(), TestDecisionsGolden(), Decision (+3 more)

### Community 26 - "context.Context"
Cohesion: 0.29
Nodes (3): Client, Identity, Renderer

## Knowledge Gaps
- **60 isolated node(s):** `signing`, `requestIDKey`, `Error`, `github.com/vdatacloud/daml-escrow-commons`, `Registry` (+55 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 116 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **7 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Error` connect `apierror.go` to `cantonid.go`?**
  _High betweenness centrality (0.029) - this node is a cross-community bridge._
- **What connects `signing`, `requestIDKey`, `Error` to the rest of the system?**
  _60 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `GitHub & CI Workflow Skill` be split into smaller, more focused modules?**
  _Cohesion score 0.1111111111111111 - nodes in this community are weakly interconnected._
- **Why does `TestCheckRefusals()` connect `approval_test.go` to `testing.T`?**
  _High betweenness centrality (0.021) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.0919661733615222 - nodes in this community are weakly interconnected._
- **Why does `Client` connect `context.Context` to `apierror_test.go`, `testing.T`?**
  _High betweenness centrality (0.016) - this node is a cross-community bridge._
- **Should `cantonid.go` be split into smaller, more focused modules?**
  _Cohesion score 0.10121457489878542 - nodes in this community are weakly interconnected._