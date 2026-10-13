# Graph Report - daml-escrow-commons  (2026-10-13)

## Corpus Check
- 77 files · ~29,469 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 3 file(s) not represented in the graph (top: (none) 3)

## Summary
- 412 nodes · 878 edges · 23 communities (18 shown, 5 thin omitted)
- Extraction: 89% EXTRACTED · 11% INFERRED · 0% AMBIGUOUS · INFERRED: 96 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `fdc072cb`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- GitHub & CI Workflow Skill
- testing.T
- widget.json
- Repository Guardrails
- draft.go
- Fingerprint
- Releasing
- metering.go
- CLAUDE.md
- /commons-contribution
- daml-escrow-commons
- .claude/CLAUDE.md
- install-git-hooks.sh
- github.com/vdatacloud/daml-escrow-commons
- agreementsig_test.go
- Client
- apierror_test.go
- next-version.sh
- test-next-version.sh
- agreement.go
- field.go
- Error

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

## Communities (23 total, 5 thin omitted)

### Community 0 - "GitHub & CI Workflow Skill"
Cohesion: 0.11
Nodes (17): 10. Useful Reference Commands, 11. References, 1. Pre-Commit Local Verification (MANDATORY — do this before every commit), 2. Branching Rules, 3. Commit Standards, 4. Staging & Pushing Changes, 5. Pull Request Creation, 6. CI Pipeline Overview (+9 more)

### Community 1 - "testing.T"
Cohesion: 0.07
Nodes (38): Parse(), golden(), TestCanonical_Golden(), TestCanonical_Properties(), TestDetailsInto(), TestFromResponse(), TestLogFields_Golden(), TestLogFields_ShortFormsOnly() (+30 more)

### Community 2 - "widget.json"
Cohesion: 0.17
Nodes (11): minLength, type, properties, name, quantity, minimum, type, required (+3 more)

### Community 3 - "Repository Guardrails"
Cohesion: 0.22
Nodes (8): Branch Protection Strategy, Branching Strategy, CI Requirements, Code Review Requirements, Commit Standard, Pre-Commit Verification, Pull Request Rules, Repository Guardrails

### Community 4 - "draft.go"
Cohesion: 0.17
Nodes (14): TestDraftCanonical_EmptyTermsAreNull(), TestDraftVectors(), DraftCanonical(), DraftHash(), DraftMessage(), DraftVersion, Sign(), TestSignAndVerify_RoundTrip() (+6 more)

### Community 5 - "Fingerprint"
Cohesion: 0.09
Nodes (23): ExpectedHash(), ForParty(), Signing(), TestSigningDetails(), SigningDetails, Classify(), FingerprintOf(), Fingerprint (+15 more)

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

### Community 14 - "agreementsig_test.go"
Cohesion: 0.11
Nodes (27): AgreementHash(), agreementCases(), agreementJSON(), flipS(), readJSON(), sha256Sum(), signatureCases(), spki() (+19 more)

### Community 15 - "Client"
Cohesion: 0.27
Nodes (4): Write(), Writer, Client, Identity

### Community 16 - "apierror_test.go"
Cohesion: 0.09
Nodes (18): CodeForStatus(), TestCodeForStatus(), TestSchemaIsACopy(), TestShortID(), TestWriter(), validate(), ValidCode(), newRequestID() (+10 more)

### Community 18 - "test-next-version.sh"
Cohesion: 0.70
Nodes (4): c(), expect(), test-next-version.sh script, tg()

### Community 19 - "agreement.go"
Cohesion: 0.17
Nodes (14): AgreementCanonical(), AgreementMessage(), checkID(), invalid(), AgreementVersion, Amends, decodeJSON(), lessUTF16() (+6 more)

### Community 20 - "field.go"
Cohesion: 0.24
Nodes (10): ErrUnknownField, child(), contains(), Registry, lookupRef(), mergeProperty(), resolveNode(), splitPointer() (+2 more)

### Community 21 - "Error"
Cohesion: 0.07
Nodes (26): draftVector, encode(), Error, New(), Newf(), sortedJSON(), init(), TestError_IsAGoError() (+18 more)

## Knowledge Gaps
- **60 isolated node(s):** `signing`, `requestIDKey`, `Error`, `github.com/vdatacloud/daml-escrow-commons`, `Registry` (+55 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 107 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **5 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Error` connect `Error` to `apierror_test.go`, `testing.T`, `Fingerprint`?**
  _High betweenness centrality (0.033) - this node is a cross-community bridge._
- **Are the 3 inferred relationships involving `Fingerprint` (e.g. with `TestJSON()` and `TestPartyID()`) actually correct?**
  _`Fingerprint` has 3 INFERRED edges - model-reasoned connections that need verification._
- **What connects `signing`, `requestIDKey`, `Error` to the rest of the system?**
  _60 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `GitHub & CI Workflow Skill` be split into smaller, more focused modules?**
  _Cohesion score 0.1111111111111111 - nodes in this community are weakly interconnected._
- **Why does `Client` connect `Client` to `apierror_test.go`, `testing.T`?**
  _High betweenness centrality (0.019) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.07346938775510205 - nodes in this community are weakly interconnected._
- **Should `Fingerprint` be split into smaller, more focused modules?**
  _Cohesion score 0.0915915915915916 - nodes in this community are weakly interconnected._