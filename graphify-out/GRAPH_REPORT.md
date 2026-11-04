# Graph Report - daml-escrow-commons  (2026-11-04)

## Corpus Check
- 109 files · ~50,141 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 9 file(s) not represented in the graph (top: .ttf 5, (none) 4)

## Summary
- 740 nodes · 1838 edges · 32 communities (25 shown, 7 thin omitted)
- Extraction: 89% EXTRACTED · 11% INFERRED · 0% AMBIGUOUS · INFERRED: 202 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `575c0a20`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- GitHub & CI Workflow Skill
- testing.T
- widget.json
- Repository Guardrails
- scrub.go
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
- logging_test.go
- notify.go
- next-version.sh
- test-next-version.sh
- canonical.go
- field.go
- Error
- FromLedgerError
- approval_test.go
- request.go
- telemetry_test.go
- TestSigningDetails
- PartyID
- auditref.go
- deliver.go

## God Nodes (most connected - your core abstractions)
1. `Refs` - 24 edges
2. `New()` - 15 edges
3. `Kind` - 14 edges
4. `Fingerprint` - 13 edges
5. `PartyID` - 12 edges
6. `fieldMap()` - 12 edges
7. `refLogger()` - 12 edges
8. `AccessLog()` - 12 edges
9. `DefaultRules()` - 12 edges
10. `GitHub & CI Workflow Skill` - 12 edges

## Surprising Connections (you probably didn't know these)
- `SigningDetails` --references--> `Fingerprint`  [EXTRACTED]
  apierror/canton/canton.go → cantonid/cantonid.go
- `SigningDetails` --references--> `PartyID`  [EXTRACTED]
  apierror/canton/canton.go → cantonid/cantonid.go
- `ForParty()` --references--> `PartyID`  [EXTRACTED]
  apierror/canton/canton.go → cantonid/cantonid.go
- `TestSigningDetails()` --calls--> `MustParsePartyID()`  [EXTRACTED]
  apierror/canton/canton_test.go → cantonid/cantonid.go
- `keyInfo` --references--> `Kind`  [EXTRACTED]
  logging/scrub.go → auditref/auditref.go

## Import Cycles
- None detected.

## Communities (32 total, 7 thin omitted)

### Community 0 - "GitHub & CI Workflow Skill"
Cohesion: 0.11
Nodes (17): 10. Useful Reference Commands, 11. References, 1. Pre-Commit Local Verification (MANDATORY — do this before every commit), 2. Branching Rules, 3. Commit Standards, 4. Staging & Pushing Changes, 5. Pull Request Creation, 6. CI Pipeline Overview (+9 more)

### Community 1 - "testing.T"
Cohesion: 0.07
Nodes (45): golden(), TestCanonical_Golden(), TestCanonical_Properties(), TestLogFields_Golden(), TestLogFields_ShortFormsOnly(), TestSchemaIsACopy(), TestWriter(), validate() (+37 more)

### Community 2 - "widget.json"
Cohesion: 0.17
Nodes (11): minLength, type, properties, name, quantity, minimum, type, required (+3 more)

### Community 3 - "Repository Guardrails"
Cohesion: 0.22
Nodes (8): Branch Protection Strategy, Branching Strategy, CI Requirements, Code Review Requirements, Commit Standard, Pre-Commit Verification, Pull Request Rules, Repository Guardrails

### Community 4 - "scrub.go"
Cohesion: 0.07
Nodes (51): accountNumber, Account(), Party(), Subject(), User(), keyInfo, personID, Rule (+43 more)

### Community 5 - "cantonid.go"
Cohesion: 0.16
Nodes (15): Classify(), FingerprintOf(), Fingerprint, MustParsePartyID(), ParseFingerprint(), ParsePartyID(), Short(), shortHex() (+7 more)

### Community 6 - "Releasing"
Cohesion: 0.18
Nodes (10): 1. Prerequisites (one-time, per machine that will `go get` this module), 2. Decide the version bump, 3. Tag and push, 4. Publish the GitHub release, 5. Update consumers, Automated (default), Future automation, Major releases (+2 more)

### Community 7 - "metering.go"
Cohesion: 0.08
Nodes (24): PendingInvitation, ChargeBearer, CommandOutcome, LedgerCommandEvent, requireBaseUnits(), requireDecimal(), NetworkFeeEvent, NetworkFeePayer (+16 more)

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
Cohesion: 0.07
Nodes (45): AgreementCanonical(), AgreementHash(), AgreementMessage(), checkID(), invalid(), agreementCases(), agreementJSON(), flipS() (+37 more)

### Community 15 - "logging_test.go"
Cohesion: 0.07
Nodes (31): CodeForStatus(), TestCodeForStatus(), ValidCode(), requestIDKey, signing, parseCase, NewHTTPSender(), staticToken() (+23 more)

### Community 16 - "notify.go"
Cohesion: 0.10
Nodes (21): actorSlot, Channel, refusal(), Message, cloneMessage(), NewMockSender(), MockSender, IsDevelopment() (+13 more)

### Community 18 - "test-next-version.sh"
Cohesion: 0.70
Nodes (4): c(), expect(), test-next-version.sh script, tg()

### Community 19 - "canonical.go"
Cohesion: 0.14
Nodes (15): lessUTF16(), writeCanonical(), writeNumber(), writeString(), numberMode, encode(), sortedJSON(), TestShortID() (+7 more)

### Community 20 - "field.go"
Cohesion: 0.24
Nodes (10): ErrUnknownField, child(), contains(), Registry, lookupRef(), mergeProperty(), resolveNode(), splitPointer() (+2 more)

### Community 21 - "Error"
Cohesion: 0.13
Nodes (11): Error, New(), Newf(), Parse(), init(), TestDetailsInto(), TestError_IsAGoError(), TestParse() (+3 more)

### Community 23 - "FromLedgerError"
Cohesion: 0.21
Nodes (10): FromError(), FromLedgerError(), GRPCCodeName(), ParseLedgerError(), TestFromError(), TestFromLedgerError_Golden(), TestFromLedgerError_Mapping(), TestLedgerError() (+2 more)

### Community 24 - "approval_test.go"
Cohesion: 0.26
Nodes (11): Check(), Decisions(), cases(), generate(), TestCheckRefusals(), TestDecisions(), TestDecisionsGolden(), Decision (+3 more)

### Community 26 - "request.go"
Cohesion: 0.09
Nodes (34): newRequestID(), RequestID(), RequestIDFrom(), Client, Identity, AccessLogOption, actorKey, loggerKey (+26 more)

### Community 27 - "telemetry_test.go"
Cohesion: 0.05
Nodes (25): TestFromResponse(), FromResponse(), Write(), Writer, accessLog, recorder, Config, failing (+17 more)

### Community 28 - "TestSigningDetails"
Cohesion: 0.31
Nodes (7): ExpectedHash(), ForParty(), Signing(), TestSigningDetails(), SigningDetails, Hash, ParseHash()

### Community 30 - "auditref.go"
Cohesion: 0.07
Nodes (29): Compute(), Kind, Refs, New(), Parse(), TestAllAndParse(), TestInstitution(), TestKeyed() (+21 more)

### Community 31 - "deliver.go"
Cohesion: 0.32
Nodes (9): Artifact, checkName(), Deliver(), artifacts(), TestDeliver_ManyIsDeterministicZip(), TestDeliver_One(), TestDeliver_RefusesBadNames(), Zip() (+1 more)

## Knowledge Gaps
- **62 isolated node(s):** `signing`, `requestIDKey`, `Error`, `github.com/vdatacloud/daml-escrow-commons`, `loggerKey` (+57 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 143 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **7 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Refs` connect `auditref.go` to `request.go`, `scrub.go`?**
  _High betweenness centrality (0.025) - this node is a cross-community bridge._
- **What connects `signing`, `requestIDKey`, `Error` to the rest of the system?**
  _62 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `GitHub & CI Workflow Skill` be split into smaller, more focused modules?**
  _Cohesion score 0.1111111111111111 - nodes in this community are weakly interconnected._
- **Why does `Error` connect `Error` to `canonical.go`, `FromLedgerError`, `TestSigningDetails`, `logging_test.go`?**
  _High betweenness centrality (0.020) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.07012987012987013 - nodes in this community are weakly interconnected._
- **Why does `TestCheckRefusals()` connect `approval_test.go` to `testing.T`?**
  _High betweenness centrality (0.016) - this node is a cross-community bridge._
- **Should `scrub.go` be split into smaller, more focused modules?**
  _Cohesion score 0.0711849957374254 - nodes in this community are weakly interconnected._