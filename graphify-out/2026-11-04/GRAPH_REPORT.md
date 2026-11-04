# Graph Report - daml-escrow-commons  (2026-11-02)

## Corpus Check
- 104 files · ~46,550 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 9 file(s) not represented in the graph (top: .ttf 5, (none) 4)

## Summary
- 693 nodes · 1692 edges · 27 communities (21 shown, 6 thin omitted)
- Extraction: 89% EXTRACTED · 11% INFERRED · 0% AMBIGUOUS · INFERRED: 185 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `75f652ca`
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
- next-version.sh
- test-next-version.sh
- Error
- field.go
- logging_test.go
- approval_test.go
- request.go
- telemetry_test.go
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
- `TestAccessLog()` --calls--> `RequestID()`  [EXTRACTED]
  logging/logging_test.go → apierror/http.go
- `keyInfo` --references--> `Kind`  [EXTRACTED]
  logging/scrub.go → auditref/auditref.go
- `refOf()` --references--> `Kind`  [EXTRACTED]
  logging/scrub.go → auditref/auditref.go
- `refLogger()` --references--> `Refs`  [EXTRACTED]
  logging/person_test.go → auditref/auditref.go
- `DefaultRules()` --references--> `Refs`  [EXTRACTED]
  logging/scrub.go → auditref/auditref.go

## Import Cycles
- None detected.

## Communities (27 total, 6 thin omitted)

### Community 0 - "GitHub & CI Workflow Skill"
Cohesion: 0.11
Nodes (17): 10. Useful Reference Commands, 11. References, 1. Pre-Commit Local Verification (MANDATORY — do this before every commit), 2. Branching Rules, 3. Commit Standards, 4. Staging & Pushing Changes, 5. Pull Request Creation, 6. CI Pipeline Overview (+9 more)

### Community 1 - "testing.T"
Cohesion: 0.08
Nodes (39): TestCanonical_Properties(), TestLogFields_ShortFormsOnly(), Sign(), TestSignAndVerify_RoundTrip(), TestVerify_MalformedHexFails(), TestVerify_TamperedMessageFails(), TestVerify_WrongSecretFails(), Verify() (+31 more)

### Community 2 - "widget.json"
Cohesion: 0.17
Nodes (11): minLength, type, properties, name, quantity, minimum, type, required (+3 more)

### Community 3 - "Repository Guardrails"
Cohesion: 0.22
Nodes (8): Branch Protection Strategy, Branching Strategy, CI Requirements, Code Review Requirements, Commit Standard, Pre-Commit Verification, Pull Request Rules, Repository Guardrails

### Community 4 - "scrub.go"
Cohesion: 0.09
Nodes (44): keyInfo, Rule, RuleFunc, TestServiceRules(), DefaultRules(), isRef(), keyOf(), KeyRedact() (+36 more)

### Community 5 - "cantonid.go"
Cohesion: 0.08
Nodes (33): ExpectedHash(), ForParty(), FromError(), FromLedgerError(), GRPCCodeName(), ParseLedgerError(), Signing(), TestFromError() (+25 more)

### Community 6 - "Releasing"
Cohesion: 0.18
Nodes (10): 1. Prerequisites (one-time, per machine that will `go get` this module), 2. Decide the version bump, 3. Tag and push, 4. Publish the GitHub release, 5. Update consumers, Automated (default), Future automation, Major releases (+2 more)

### Community 7 - "metering.go"
Cohesion: 0.08
Nodes (23): PendingInvitation, ChargeBearer, CommandOutcome, LedgerCommandEvent, requireBaseUnits(), requireDecimal(), NetworkFeeEvent, NetworkFeePayer (+15 more)

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
Cohesion: 0.06
Nodes (46): AgreementCanonical(), AgreementHash(), AgreementMessage(), checkID(), invalid(), agreementCases(), agreementJSON(), flipS() (+38 more)

### Community 18 - "test-next-version.sh"
Cohesion: 0.70
Nodes (4): c(), expect(), test-next-version.sh script, tg()

### Community 19 - "Error"
Cohesion: 0.09
Nodes (17): lessUTF16(), writeCanonical(), writeNumber(), writeString(), numberMode, encode(), Error, sortedJSON() (+9 more)

### Community 20 - "field.go"
Cohesion: 0.24
Nodes (10): ErrUnknownField, child(), contains(), Registry, lookupRef(), mergeProperty(), resolveNode(), splitPointer() (+2 more)

### Community 21 - "logging_test.go"
Cohesion: 0.05
Nodes (38): CodeForStatus(), New(), Newf(), Parse(), golden(), init(), TestCanonical_Golden(), TestCodeForStatus() (+30 more)

### Community 24 - "approval_test.go"
Cohesion: 0.26
Nodes (11): Check(), Decisions(), cases(), generate(), TestCheckRefusals(), TestDecisions(), TestDecisionsGolden(), Decision (+3 more)

### Community 26 - "request.go"
Cohesion: 0.09
Nodes (33): RequestIDFrom(), Client, Identity, AccessLogOption, actorKey, actorSlot, loggerKey, fieldMap() (+25 more)

### Community 27 - "telemetry_test.go"
Cohesion: 0.05
Nodes (23): Write(), Writer, accessLog, recorder, Config, failing, Middleware(), Transport() (+15 more)

### Community 30 - "auditref.go"
Cohesion: 0.06
Nodes (36): Compute(), Kind, Refs, New(), Parse(), TestAllAndParse(), TestInstitution(), TestKeyed() (+28 more)

### Community 31 - "deliver.go"
Cohesion: 0.29
Nodes (9): Artifact, checkName(), Deliver(), artifacts(), TestDeliver_ManyIsDeterministicZip(), TestDeliver_One(), TestDeliver_RefusesBadNames(), Zip() (+1 more)

## Knowledge Gaps
- **62 isolated node(s):** `signing`, `requestIDKey`, `Error`, `github.com/vdatacloud/daml-escrow-commons`, `loggerKey` (+57 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 142 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **6 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Refs` connect `auditref.go` to `request.go`, `scrub.go`?**
  _High betweenness centrality (0.026) - this node is a cross-community bridge._
- **What connects `signing`, `requestIDKey`, `Error` to the rest of the system?**
  _62 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `GitHub & CI Workflow Skill` be split into smaller, more focused modules?**
  _Cohesion score 0.1111111111111111 - nodes in this community are weakly interconnected._
- **Why does `Error` connect `Error` to `cantonid.go`, `logging_test.go`?**
  _High betweenness centrality (0.021) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.07918367346938776 - nodes in this community are weakly interconnected._
- **Why does `TestCheckRefusals()` connect `approval_test.go` to `testing.T`?**
  _High betweenness centrality (0.016) - this node is a cross-community bridge._
- **Should `scrub.go` be split into smaller, more focused modules?**
  _Cohesion score 0.0853302162478083 - nodes in this community are weakly interconnected._