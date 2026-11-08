# Graph Report - daml-escrow-commons  (2026-11-04)

## Corpus Check
- 111 files · ~51,986 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 9 file(s) not represented in the graph (top: .ttf 5, (none) 4)

## Summary
- 783 nodes · 1938 edges · 38 communities (32 shown, 6 thin omitted)
- Extraction: 89% EXTRACTED · 11% INFERRED · 0% AMBIGUOUS · INFERRED: 206 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `706bfc34`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- GitHub & CI Workflow Skill
- testing.T
- widget.json
- Repository Guardrails
- redact.go
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
- shortid.go
- field.go
- apierror.go
- canton.go
- approval_test.go
- request.go
- telemetry_test.go
- apierror_test.go
- smtptest.go
- scrub.go
- go_pkg_context
- notify/http.go
- MockServer
- smtptest_test.go
- mock.go
- apierror/http.go
- TestFromLedgerError_Golden

## God Nodes (most connected - your core abstractions)
1. `Refs` - 24 edges
2. `MockServer` - 16 edges
3. `New()` - 15 edges
4. `Kind` - 14 edges
5. `Fingerprint` - 13 edges
6. `PartyID` - 12 edges
7. `fieldMap()` - 12 edges
8. `refLogger()` - 12 edges
9. `AccessLog()` - 12 edges
10. `DefaultRules()` - 12 edges

## Surprising Connections (you probably didn't know these)
- `SigningDetails` --references--> `Fingerprint`  [EXTRACTED]
  apierror/canton/canton.go → cantonid/cantonid.go
- `SigningDetails` --references--> `PartyID`  [EXTRACTED]
  apierror/canton/canton.go → cantonid/cantonid.go
- `ForParty()` --references--> `PartyID`  [EXTRACTED]
  apierror/canton/canton.go → cantonid/cantonid.go
- `TestSigningDetails()` --calls--> `MustParsePartyID()`  [EXTRACTED]
  apierror/canton/canton_test.go → cantonid/cantonid.go
- `TestAccessLog()` --calls--> `RequestID()`  [EXTRACTED]
  logging/logging_test.go → apierror/http.go

## Import Cycles
- None detected.

## Communities (38 total, 6 thin omitted)

### Community 0 - "GitHub & CI Workflow Skill"
Cohesion: 0.11
Nodes (17): 10. Useful Reference Commands, 11. References, 1. Pre-Commit Local Verification (MANDATORY — do this before every commit), 2. Branching Rules, 3. Commit Standards, 4. Staging & Pushing Changes, 5. Pull Request Creation, 6. CI Pipeline Overview (+9 more)

### Community 1 - "testing.T"
Cohesion: 0.07
Nodes (51): TestCanonical_Properties(), TestLogFields_ShortFormsOnly(), New(), TestClient_CreatePendingInvitation(), TestClient_CreatePendingInvitation_ServerError(), TestClient_GetByEmail_ServerError(), TestClient_GetByOktaSub_Found(), TestClient_GetByToken_NotFound() (+43 more)

### Community 2 - "widget.json"
Cohesion: 0.17
Nodes (11): minLength, type, properties, name, quantity, minimum, type, required (+3 more)

### Community 3 - "Repository Guardrails"
Cohesion: 0.22
Nodes (8): Branch Protection Strategy, Branching Strategy, CI Requirements, Code Review Requirements, Commit Standard, Pre-Commit Verification, Pull Request Rules, Repository Guardrails

### Community 4 - "redact.go"
Cohesion: 0.10
Nodes (32): accountNumber, Account(), Party(), Subject(), User(), personID, TestTypedFields(), redactAs() (+24 more)

### Community 5 - "cantonid.go"
Cohesion: 0.13
Nodes (16): Classify(), FingerprintOf(), Fingerprint, PartyID, MustParsePartyID(), ParseFingerprint(), ParsePartyID(), Short() (+8 more)

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
Cohesion: 0.05
Nodes (55): AgreementCanonical(), AgreementHash(), AgreementMessage(), checkID(), invalid(), agreementCases(), agreementJSON(), flipS() (+47 more)

### Community 15 - "logging_test.go"
Cohesion: 0.17
Nodes (3): parseCase, ErrUnknownType, ValidationError

### Community 16 - "notify.go"
Cohesion: 0.24
Nodes (5): Channel, Message, Recipient, Sender, Status

### Community 18 - "test-next-version.sh"
Cohesion: 0.70
Nodes (4): c(), expect(), test-next-version.sh script, tg()

### Community 19 - "shortid.go"
Cohesion: 0.33
Nodes (6): TestShortID(), Field, Error, ShortID(), ShortIDs(), shortValue()

### Community 20 - "field.go"
Cohesion: 0.24
Nodes (10): ErrUnknownField, child(), contains(), Registry, lookupRef(), mergeProperty(), resolveNode(), splitPointer() (+2 more)

### Community 21 - "apierror.go"
Cohesion: 0.16
Nodes (9): CodeForStatus(), encode(), Error, sortedJSON(), TestCodeForStatus(), ValidCode(), Detail, Upstream (+1 more)

### Community 23 - "canton.go"
Cohesion: 0.18
Nodes (15): ExpectedHash(), ForParty(), FromError(), FromLedgerError(), GRPCCodeName(), ParseLedgerError(), Signing(), TestFromError() (+7 more)

### Community 24 - "approval_test.go"
Cohesion: 0.26
Nodes (11): Check(), Decisions(), cases(), generate(), TestCheckRefusals(), TestDecisions(), TestDecisionsGolden(), Decision (+3 more)

### Community 26 - "request.go"
Cohesion: 0.08
Nodes (31): RequestIDFrom(), Write(), Writer, Client, Identity, accessLog, AccessLogOption, actorKey (+23 more)

### Community 27 - "telemetry_test.go"
Cohesion: 0.07
Nodes (18): Config, failing, Middleware(), Transport(), isURL(), metricOptions(), Propagator(), Setup() (+10 more)

### Community 28 - "apierror_test.go"
Cohesion: 0.16
Nodes (18): New(), Newf(), Parse(), golden(), init(), TestCanonical_Golden(), TestDetailsInto(), TestError_IsAGoError() (+10 more)

### Community 29 - "smtptest.go"
Cohesion: 0.16
Nodes (8): readPart(), TestMockServer_RecordsWhatNetSMTPSends(), TestMockServer_TooBig(), WithMaxSize(), WithOnMessage(), Content, Message, Option

### Community 30 - "scrub.go"
Cohesion: 0.05
Nodes (54): Compute(), Kind, Refs, New(), Parse(), TestAllAndParse(), TestInstitution(), TestKeyed() (+46 more)

### Community 31 - "go_pkg_context"
Cohesion: 0.12
Nodes (20): Artifact, checkName(), Deliver(), artifacts(), TestDeliver_ManyIsDeterministicZip(), TestDeliver_One(), TestDeliver_RefusesBadNames(), Zip() (+12 more)

### Community 32 - "notify/http.go"
Cohesion: 0.21
Nodes (6): refusal(), WithHTTPClient(), HTTPOption, HTTPSender, Receipt, TokenFunc

### Community 33 - "MockServer"
Cohesion: 0.19
Nodes (3): path(), readData(), MockServer

### Community 34 - "smtptest_test.go"
Cohesion: 0.27
Nodes (6): NewMockServer(), start(), TestMockServer_FailWith(), TestMockServer_Protocol(), TestMockServer_SinglePart(), TestNewMockServer_DevOnly()

### Community 35 - "mock.go"
Cohesion: 0.36
Nodes (3): actorSlot, cloneMessage(), MockSender

### Community 36 - "apierror/http.go"
Cohesion: 0.29
Nodes (3): newRequestID(), RequestID(), requestIDKey

### Community 37 - "TestFromLedgerError_Golden"
Cohesion: 0.33
Nodes (4): draftVector, TestFromLedgerError_Golden(), decodeInto(), fixture

## Knowledge Gaps
- **62 isolated node(s):** `signing`, `requestIDKey`, `Error`, `github.com/vdatacloud/daml-escrow-commons`, `loggerKey` (+57 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 154 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **6 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `MockServer` connect `MockServer` to `telemetry_test.go`, `smtptest_test.go`, `mock.go`, `smtptest.go`?**
  _High betweenness centrality (0.024) - this node is a cross-community bridge._
- **What connects `signing`, `requestIDKey`, `Error` to the rest of the system?**
  _62 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `GitHub & CI Workflow Skill` be split into smaller, more focused modules?**
  _Cohesion score 0.1111111111111111 - nodes in this community are weakly interconnected._
- **Why does `Error` connect `apierror.go` to `apierror_test.go`, `canton.go`?**
  _High betweenness centrality (0.019) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.06597222222222222 - nodes in this community are weakly interconnected._
- **Should `redact.go` be split into smaller, more focused modules?**
  _Cohesion score 0.1039136302294197 - nodes in this community are weakly interconnected._
- **Should `cantonid.go` be split into smaller, more focused modules?**
  _Cohesion score 0.12807881773399016 - nodes in this community are weakly interconnected._