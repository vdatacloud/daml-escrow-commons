# Graph Report - daml-escrow-commons  (2026-11-09)

## Corpus Check
- 129 files · ~59,876 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 9 file(s) not represented in the graph (top: .ttf 5, (none) 4)

## Summary
- 912 nodes · 2345 edges · 41 communities (35 shown, 6 thin omitted)
- Extraction: 90% EXTRACTED · 10% INFERRED · 0% AMBIGUOUS · INFERRED: 238 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `762c3767`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- GitHub & CI Workflow Skill
- widget.json
- Repository Guardrails
- scrub.go
- cantonid.go
- Releasing
- go_pkg_fmt
- CLAUDE.md
- /commons-contribution
- daml-escrow-commons
- .claude/CLAUDE.md
- install-git-hooks.sh
- github.com/vdatacloud/daml-escrow-commons
- agreement.go
- canton_test.go
- notify.go
- next-version.sh
- test-next-version.sh
- shortid.go
- field.go
- apierror.go
- canton.go
- approval_test.go
- context.Context
- go_pkg_net_http
- apierror_test.go
- smtptest.go
- go_pkg_strings
- go_pkg_context
- notify/http.go
- MockServer
- smtptest_test.go
- mock.go
- apierror/http.go
- go_pkg_testing
- Issuer
- testing.T
- notify_test.go
- http_test.go

## God Nodes (most connected - your core abstractions)
1. `Issuer` - 28 edges
2. `Refs` - 25 edges
3. `New()` - 16 edges
4. `MockServer` - 16 edges
5. `Kind` - 14 edges
6. `Principal` - 14 edges
7. `Fingerprint` - 13 edges
8. `New()` - 12 edges
9. `PartyID` - 12 edges
10. `fieldMap()` - 12 edges

## Surprising Connections (you probably didn't know these)
- `Middleware()` --calls--> `New()`  [EXTRACTED]
  authn/middleware.go → apierror/apierror.go
- `SigningDetails` --references--> `Fingerprint`  [EXTRACTED]
  apierror/canton/canton.go → cantonid/cantonid.go
- `SigningDetails` --references--> `Hash`  [EXTRACTED]
  apierror/canton/canton.go → cantonid/cantonid.go
- `SigningDetails` --references--> `PartyID`  [EXTRACTED]
  apierror/canton/canton.go → cantonid/cantonid.go
- `ForParty()` --references--> `PartyID`  [EXTRACTED]
  apierror/canton/canton.go → cantonid/cantonid.go

## Import Cycles
- None detected.

## Communities (41 total, 6 thin omitted)

### Community 0 - "GitHub & CI Workflow Skill"
Cohesion: 0.11
Nodes (17): 10. Useful Reference Commands, 11. References, 1. Pre-Commit Local Verification (MANDATORY — do this before every commit), 2. Branching Rules, 3. Commit Standards, 4. Staging & Pushing Changes, 5. Pull Request Creation, 6. CI Pipeline Overview (+9 more)

### Community 2 - "widget.json"
Cohesion: 0.17
Nodes (11): minLength, type, properties, name, quantity, minimum, type, required (+3 more)

### Community 3 - "Repository Guardrails"
Cohesion: 0.22
Nodes (8): Branch Protection Strategy, Branching Strategy, CI Requirements, Code Review Requirements, Commit Standard, Pre-Commit Verification, Pull Request Rules, Repository Guardrails

### Community 4 - "scrub.go"
Cohesion: 0.06
Nodes (51): Compute(), Kind, Refs, cacheKey, keyInfo, personKind(), TestPersonKind(), personID (+43 more)

### Community 5 - "cantonid.go"
Cohesion: 0.09
Nodes (20): parseEd25519(), Classify(), FingerprintOf(), Fingerprint, Hash, PartyID, MustParsePartyID(), ParseFingerprint() (+12 more)

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

### Community 14 - "agreement.go"
Cohesion: 0.17
Nodes (19): AgreementCanonical(), AgreementMessage(), checkID(), invalid(), TestDraftVectors(), AgreementVersion, Amends, decodeJSON() (+11 more)

### Community 15 - "canton_test.go"
Cohesion: 0.11
Nodes (6): TestFromLedgerError_Golden(), fixture, TestLegalOrderWireNames(), ErrUnknownType, Registry, ValidationError

### Community 16 - "notify.go"
Cohesion: 0.27
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
Cohesion: 0.15
Nodes (11): CodeForStatus(), encode(), Error, Parse(), sortedJSON(), TestCodeForStatus(), TestParse(), ValidCode() (+3 more)

### Community 23 - "canton.go"
Cohesion: 0.24
Nodes (13): ExpectedHash(), ForParty(), FromError(), FromLedgerError(), GRPCCodeName(), ParseLedgerError(), Signing(), TestFromError() (+5 more)

### Community 24 - "approval_test.go"
Cohesion: 0.29
Nodes (11): Check(), Decisions(), cases(), generate(), TestCheckRefusals(), TestDecisions(), TestDecisionsGolden(), Decision (+3 more)

### Community 26 - "context.Context"
Cohesion: 0.08
Nodes (36): Newf(), RequestIDFrom(), TestActorMiddleware(), TestRequireScope(), claims, ctxKey, fakeVerifier, ActorMiddleware() (+28 more)

### Community 27 - "go_pkg_net_http"
Cohesion: 0.07
Nodes (18): Config, failing, Middleware(), Transport(), isURL(), metricOptions(), Propagator(), Setup() (+10 more)

### Community 28 - "apierror_test.go"
Cohesion: 0.20
Nodes (15): New(), golden(), init(), TestCanonical_Golden(), TestCanonical_Properties(), TestDetailsInto(), TestError_IsAGoError(), TestLogFields_Golden() (+7 more)

### Community 29 - "smtptest.go"
Cohesion: 0.19
Nodes (6): readPart(), WithMaxSize(), WithOnMessage(), Content, Message, Option

### Community 30 - "go_pkg_strings"
Cohesion: 0.05
Nodes (55): New(), Parse(), TestAllAndParse(), TestInstitution(), TestKeyed(), TestKeyedFailsClosed(), TestPlainInDevelopment(), TestRefsCache() (+47 more)

### Community 31 - "go_pkg_context"
Cohesion: 0.06
Nodes (30): LegalOrder, LegalOrderEvent, LegalOrderTarget, LiftLegalOrderResponse, ListLegalOrdersRequest, OrderContact, RecordedLegalOrder, Artifact (+22 more)

### Community 32 - "notify/http.go"
Cohesion: 0.26
Nodes (6): refusal(), WithHTTPClient(), HTTPOption, HTTPSender, Receipt, TokenFunc

### Community 33 - "MockServer"
Cohesion: 0.19
Nodes (3): path(), readData(), MockServer

### Community 34 - "smtptest_test.go"
Cohesion: 0.31
Nodes (8): NewMockServer(), start(), TestMockServer_FailWith(), TestMockServer_Protocol(), TestMockServer_RecordsWhatNetSMTPSends(), TestMockServer_SinglePart(), TestMockServer_TooBig(), TestNewMockServer_DevOnly()

### Community 35 - "mock.go"
Cohesion: 0.36
Nodes (3): actorSlot, cloneMessage(), MockSender

### Community 36 - "apierror/http.go"
Cohesion: 0.22
Nodes (7): draftVector, TestFromResponse(), decodeInto(), FromResponse(), newRequestID(), RequestID(), requestIDKey

### Community 37 - "go_pkg_testing"
Cohesion: 0.07
Nodes (29): main(), splitOrigins(), TestSplitOrigins(), authorizeQuery(), browserIssuer(), challengeOf(), exchange(), signIn() (+21 more)

### Community 38 - "Issuer"
Cohesion: 0.09
Nodes (11): Write(), Writer, authCode, Issuer, codeStore, Issuer, httpHandlerFunc(), writeJSON() (+3 more)

### Community 39 - "testing.T"
Cohesion: 0.05
Nodes (66): AgreementHash(), agreementCases(), agreementJSON(), flipS(), readJSON(), sha256Sum(), signatureCases(), spki() (+58 more)

### Community 41 - "notify_test.go"
Cohesion: 0.27
Nodes (11): NewMockSender(), IsDevelopment(), TestIsDevelopment(), TestMockSender_RecordsOncePerIdempotencyKey(), TestMockSender_SentIsACopy(), TestMockSender_ValidatesAndFails(), TestNewMockSender_RefusedOutsideDevelopment(), TestValidate_AcceptsAWellFormedMessage() (+3 more)

### Community 43 - "http_test.go"
Cohesion: 0.31
Nodes (8): NewHTTPSender(), staticToken(), TestHTTPSender_RefusesAnInvalidMessageWithoutCalling(), TestHTTPSender_SendsAndReadsTheReceipt(), TestHTTPSender_Status(), TestHTTPSender_StatusAnswersToTheRelay(), TestHTTPSender_TransportFailuresAreRetryable(), TestNewHTTPSender_Refusals()

## Knowledge Gaps
- **66 isolated node(s):** `signing`, `requestIDKey`, `Error`, `ctxKey`, `github.com/vdatacloud/daml-escrow-commons` (+61 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 170 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **6 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Issuer` connect `Issuer` to `go_pkg_testing`, `testing.T`, `go_pkg_fmt`, `canton_test.go`, `shortid.go`, `go_pkg_net_http`, `go_pkg_strings`, `go_pkg_context`?**
  _High betweenness centrality (0.036) - this node is a cross-community bridge._
- **What connects `signing`, `requestIDKey`, `Error` to the rest of the system?**
  _66 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `GitHub & CI Workflow Skill` be split into smaller, more focused modules?**
  _Cohesion score 0.1111111111111111 - nodes in this community are weakly interconnected._
- **Why does `MockServer` connect `MockServer` to `go_pkg_net_http`, `smtptest_test.go`, `mock.go`, `smtptest.go`?**
  _High betweenness centrality (0.021) - this node is a cross-community bridge._
- **Should `scrub.go` be split into smaller, more focused modules?**
  _Cohesion score 0.06227106227106227 - nodes in this community are weakly interconnected._
- **Why does `Refs` connect `scrub.go` to `context.Context`, `go_pkg_strings`?**
  _High betweenness centrality (0.021) - this node is a cross-community bridge._
- **Should `cantonid.go` be split into smaller, more focused modules?**
  _Cohesion score 0.09024390243902439 - nodes in this community are weakly interconnected._