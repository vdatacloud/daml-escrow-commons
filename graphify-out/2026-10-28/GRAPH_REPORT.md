# Graph Report - daml-escrow-commons  (2026-10-20)

## Corpus Check
- 84 files · ~33,494 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 9 file(s) not represented in the graph (top: .ttf 5, (none) 4)

## Summary
- 460 nodes · 1003 edges · 26 communities (20 shown, 6 thin omitted)
- Extraction: 89% EXTRACTED · 11% INFERRED · 0% AMBIGUOUS · INFERRED: 108 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `e41ccde2`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- GitHub & CI Workflow Skill
- testing.T
- widget.json
- Repository Guardrails
- draft.go
- cantonid.go
- Releasing
- go_pkg_fmt
- CLAUDE.md
- /commons-contribution
- daml-escrow-commons
- .claude/CLAUDE.md
- install-git-hooks.sh
- github.com/vdatacloud/daml-escrow-commons
- agreementsig_test.go
- render.go
- canton_test.go
- next-version.sh
- test-next-version.sh
- canonical.go
- field.go
- canton.go
- http.go
- deliver.go

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
- `init()` --calls--> `New()`  [INFERRED]
  apierror/apierror_test.go → apierror/apierror.go
- `SigningDetails` --references--> `Fingerprint`  [EXTRACTED]
  apierror/canton/canton.go → cantonid/cantonid.go
- `SigningDetails` --references--> `PartyID`  [EXTRACTED]
  apierror/canton/canton.go → cantonid/cantonid.go
- `ForParty()` --references--> `PartyID`  [EXTRACTED]
  apierror/canton/canton.go → cantonid/cantonid.go
- `TestSigningDetails()` --calls--> `MustParsePartyID()`  [EXTRACTED]
  apierror/canton/canton_test.go → cantonid/cantonid.go

## Import Cycles
- None detected.

## Communities (26 total, 6 thin omitted)

### Community 0 - "GitHub & CI Workflow Skill"
Cohesion: 0.11
Nodes (17): 10. Useful Reference Commands, 11. References, 1. Pre-Commit Local Verification (MANDATORY — do this before every commit), 2. Branching Rules, 3. Commit Standards, 4. Staging & Pushing Changes, 5. Pull Request Creation, 6. CI Pipeline Overview (+9 more)

### Community 1 - "testing.T"
Cohesion: 0.07
Nodes (47): CodeForStatus(), Parse(), golden(), init(), TestCanonical_Golden(), TestCanonical_Properties(), TestCodeForStatus(), TestDetailsInto() (+39 more)

### Community 2 - "widget.json"
Cohesion: 0.17
Nodes (11): minLength, type, properties, name, quantity, minimum, type, required (+3 more)

### Community 3 - "Repository Guardrails"
Cohesion: 0.22
Nodes (8): Branch Protection Strategy, Branching Strategy, CI Requirements, Code Review Requirements, Commit Standard, Pre-Commit Verification, Pull Request Rules, Repository Guardrails

### Community 4 - "draft.go"
Cohesion: 0.19
Nodes (12): TestDraftCanonical_EmptyTermsAreNull(), TestDraftVectors(), DraftCanonical(), DraftHash(), DraftMessage(), DraftVersion, Sign(), TestSignAndVerify_RoundTrip() (+4 more)

### Community 5 - "cantonid.go"
Cohesion: 0.10
Nodes (18): parseEd25519(), Classify(), FingerprintOf(), Fingerprint, PartyID, MustParsePartyID(), ParseFingerprint(), ParsePartyID() (+10 more)

### Community 6 - "Releasing"
Cohesion: 0.18
Nodes (10): 1. Prerequisites (one-time, per machine that will `go get` this module), 2. Decide the version bump, 3. Tag and push, 4. Publish the GitHub release, 5. Update consumers, Automated (default), Future automation, Major releases (+2 more)

### Community 7 - "go_pkg_fmt"
Cohesion: 0.12
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
Cohesion: 0.08
Nodes (36): AgreementCanonical(), AgreementHash(), AgreementMessage(), checkID(), invalid(), agreementCases(), agreementJSON(), flipS() (+28 more)

### Community 15 - "render.go"
Cohesion: 0.08
Nodes (17): Client, Identity, PendingInvitation, Fake, CheckFonts(), NewTypst(), TestCheckFonts(), TestFakeIsDeterministic() (+9 more)

### Community 16 - "canton_test.go"
Cohesion: 0.10
Nodes (11): ParseLedgerError(), TestFromError(), TestFromLedgerError_Golden(), TestFromLedgerError_Mapping(), TestLedgerError(), fixture, Ref(), TestRef() (+3 more)

### Community 18 - "test-next-version.sh"
Cohesion: 0.70
Nodes (4): c(), expect(), test-next-version.sh script, tg()

### Community 19 - "canonical.go"
Cohesion: 0.16
Nodes (12): lessUTF16(), writeCanonical(), writeNumber(), writeString(), numberMode, TestShortID(), Field, decodeInto() (+4 more)

### Community 20 - "field.go"
Cohesion: 0.24
Nodes (10): ErrUnknownField, child(), contains(), Registry, lookupRef(), mergeProperty(), resolveNode(), splitPointer() (+2 more)

### Community 21 - "canton.go"
Cohesion: 0.11
Nodes (19): encode(), Error, New(), Newf(), sortedJSON(), ExpectedHash(), ForParty(), FromError() (+11 more)

### Community 23 - "http.go"
Cohesion: 0.17
Nodes (9): TestFromResponse(), TestWriter(), FromResponse(), newRequestID(), RequestID(), RequestIDFrom(), Write(), requestIDKey (+1 more)

### Community 24 - "deliver.go"
Cohesion: 0.29
Nodes (9): Artifact, checkName(), Deliver(), artifacts(), TestDeliver_ManyIsDeterministicZip(), TestDeliver_One(), TestDeliver_RefusesBadNames(), Zip() (+1 more)

## Knowledge Gaps
- **60 isolated node(s):** `signing`, `requestIDKey`, `Error`, `github.com/vdatacloud/daml-escrow-commons`, `Registry` (+55 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 115 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **6 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Error` connect `canton.go` to `testing.T`?**
  _High betweenness centrality (0.030) - this node is a cross-community bridge._
- **What connects `signing`, `requestIDKey`, `Error` to the rest of the system?**
  _60 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `GitHub & CI Workflow Skill` be split into smaller, more focused modules?**
  _Cohesion score 0.1111111111111111 - nodes in this community are weakly interconnected._
- **Why does `Fingerprint` connect `cantonid.go` to `canton.go`?**
  _High betweenness centrality (0.016) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.07456140350877193 - nodes in this community are weakly interconnected._
- **Why does `Client` connect `render.go` to `testing.T`?**
  _High betweenness centrality (0.016) - this node is a cross-community bridge._
- **Should `cantonid.go` be split into smaller, more focused modules?**
  _Cohesion score 0.1 - nodes in this community are weakly interconnected._