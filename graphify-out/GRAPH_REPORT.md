# Graph Report - daml-escrow-commons  (2026-11-08)

## Corpus Check
- 127 files · ~59,190 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 9 file(s) not represented in the graph (top: .ttf 5, (none) 4)

## Summary
- 902 nodes · 2328 edges · 46 communities (39 shown, 7 thin omitted)
- Extraction: 90% EXTRACTED · 10% INFERRED · 0% AMBIGUOUS · INFERRED: 238 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `b9b2ae25`
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
- agreement.go
- go_pkg_fmt
- notify.go
- next-version.sh
- test-next-version.sh
- shortid.go
- field.go
- apierror.go
- canton.go
- go_pkg_encoding_json
- logging_test.go
- telemetry_test.go
- apierror_test.go
- smtptest.go
- go_pkg_strings
- go_pkg_bytes
- notify/http.go
- MockServer
- smtptest_test.go
- mock.go
- apierror/http.go
- go_pkg_net_http
- Issuer
- agreementsig_test.go
- .Verify
- notify_test.go
- identityclient_test.go
- http_test.go
- LoadDirectory
- Registry

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

## Communities (46 total, 7 thin omitted)

### Community 0 - "GitHub & CI Workflow Skill"
Cohesion: 0.11
Nodes (17): 10. Useful Reference Commands, 11. References, 1. Pre-Commit Local Verification (MANDATORY — do this before every commit), 2. Branching Rules, 3. Commit Standards, 4. Staging & Pushing Changes, 5. Pull Request Creation, 6. CI Pipeline Overview (+9 more)

### Community 1 - "testing.T"
Cohesion: 0.22
Nodes (14): TestLedgerCommandEvent_Timing(), TestLedgerCommandEvent_Traffic(), TestLedgerCommandEvent_Validate(), TestNetworkFeeEvent_SpotPrice(), TestNetworkFeeEvent_Validate(), TestSettlementEvent_NetworkFee(), TestSettlementEvent_Validate(), withFee() (+6 more)

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
Cohesion: 0.07
Nodes (25): Algorithm, ecdsaSig, parseCanonicalECDSA(), parseEd25519(), parseP256(), Verify(), Classify(), FingerprintOf() (+17 more)

### Community 6 - "Releasing"
Cohesion: 0.18
Nodes (10): 1. Prerequisites (one-time, per machine that will `go get` this module), 2. Decide the version bump, 3. Tag and push, 4. Publish the GitHub release, 5. Update consumers, Automated (default), Future automation, Major releases (+2 more)

### Community 7 - "metering.go"
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
Cohesion: 0.11
Nodes (21): AgreementCanonical(), AgreementMessage(), checkID(), invalid(), TestDraftVectors(), AgreementVersion, Amends, decodeJSON() (+13 more)

### Community 15 - "go_pkg_fmt"
Cohesion: 0.15
Nodes (6): TestSchemaIsACopy(), TestFromLedgerError_Golden(), Schema(), fixture, ErrUnknownType, ValidationError

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
Cohesion: 0.16
Nodes (9): CodeForStatus(), encode(), Error, sortedJSON(), TestCodeForStatus(), ValidCode(), Detail, Upstream (+1 more)

### Community 23 - "canton.go"
Cohesion: 0.24
Nodes (13): ExpectedHash(), ForParty(), FromError(), FromLedgerError(), GRPCCodeName(), ParseLedgerError(), Signing(), TestFromError() (+5 more)

### Community 24 - "go_pkg_encoding_json"
Cohesion: 0.26
Nodes (11): Check(), Decisions(), cases(), generate(), TestCheckRefusals(), TestDecisions(), TestDecisionsGolden(), Decision (+3 more)

### Community 26 - "logging_test.go"
Cohesion: 0.10
Nodes (36): Newf(), TestActorMiddleware(), TestRequireScope(), ActorMiddleware(), RequireScope(), FromContext(), HasScope(), WithPrincipal() (+28 more)

### Community 27 - "telemetry_test.go"
Cohesion: 0.08
Nodes (18): Config, failing, Middleware(), Transport(), isURL(), metricOptions(), Propagator(), Setup() (+10 more)

### Community 28 - "apierror_test.go"
Cohesion: 0.18
Nodes (17): New(), Parse(), golden(), init(), TestCanonical_Golden(), TestCanonical_Properties(), TestDetailsInto(), TestError_IsAGoError() (+9 more)

### Community 29 - "smtptest.go"
Cohesion: 0.17
Nodes (8): readPart(), TestMockServer_RecordsWhatNetSMTPSends(), TestMockServer_TooBig(), WithMaxSize(), WithOnMessage(), Content, Message, Option

### Community 30 - "go_pkg_strings"
Cohesion: 0.05
Nodes (53): Compute(), Kind, Refs, New(), Parse(), TestAllAndParse(), TestInstitution(), TestKeyed() (+45 more)

### Community 31 - "go_pkg_bytes"
Cohesion: 0.07
Nodes (24): PendingInvitation, Artifact, checkName(), Deliver(), artifacts(), TestDeliver_ManyIsDeterministicZip(), TestDeliver_One(), TestDeliver_RefusesBadNames() (+16 more)

### Community 32 - "notify/http.go"
Cohesion: 0.26
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
Nodes (5): draftVector, decodeInto(), newRequestID(), RequestID(), requestIDKey

### Community 37 - "go_pkg_net_http"
Cohesion: 0.05
Nodes (34): claims, ctxKey, NewOIDCVerifier(), uniqueScopes(), OIDCVerifier, Principal, TokenVerifier, main() (+26 more)

### Community 38 - "Issuer"
Cohesion: 0.09
Nodes (12): RequestIDFrom(), Write(), Writer, authCode, Issuer, codeStore, Issuer, httpHandlerFunc() (+4 more)

### Community 39 - "agreementsig_test.go"
Cohesion: 0.18
Nodes (22): AgreementHash(), agreementCases(), agreementJSON(), flipS(), readJSON(), sha256Sum(), signatureCases(), spki() (+14 more)

### Community 40 - ".Verify"
Cohesion: 0.20
Nodes (11): signer(), TestOIDCVerifierReadsClaims(), TestOIDCVerifierRefuses(), token(), fakeVerifier, NewOIDCVerifierWithKeys(), Sign(), TestSignAndVerify_RoundTrip() (+3 more)

### Community 41 - "notify_test.go"
Cohesion: 0.27
Nodes (11): NewMockSender(), IsDevelopment(), TestIsDevelopment(), TestMockSender_RecordsOncePerIdempotencyKey(), TestMockSender_SentIsACopy(), TestMockSender_ValidatesAndFails(), TestNewMockSender_RefusedOutsideDevelopment(), TestValidate_AcceptsAWellFormedMessage() (+3 more)

### Community 42 - "identityclient_test.go"
Cohesion: 0.36
Nodes (9): New(), TestClient_CreatePendingInvitation(), TestClient_CreatePendingInvitation_ServerError(), TestClient_GetByEmail_ServerError(), TestClient_GetByOktaSub_Found(), TestClient_GetByToken_NotFound(), TestClient_ManagesIdentity(), TestClient_ManagesIdentity_ServerError() (+1 more)

### Community 43 - "http_test.go"
Cohesion: 0.43
Nodes (7): NewHTTPSender(), TestHTTPSender_RefusesAnInvalidMessageWithoutCalling(), TestHTTPSender_SendsAndReadsTheReceipt(), TestHTTPSender_Status(), TestHTTPSender_StatusAnswersToTheRelay(), TestHTTPSender_TransportFailuresAreRetryable(), TestNewHTTPSender_Refusals()

### Community 44 - "LoadDirectory"
Cohesion: 0.48
Nodes (6): LoadDirectory(), TestLoadDirectory_CompilesSchemas(), TestLoadDirectory_MissingDirectory(), TestValidate_InvalidPayloadReportsFailures(), TestValidate_UnknownType(), TestValidate_ValidPayload()

## Knowledge Gaps
- **63 isolated node(s):** `signing`, `requestIDKey`, `Error`, `ctxKey`, `github.com/vdatacloud/daml-escrow-commons` (+58 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 167 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **7 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Issuer` connect `Issuer` to `cantonid.go`, `go_pkg_net_http`, `agreement.go`, `go_pkg_fmt`, `shortid.go`, `go_pkg_encoding_json`, `go_pkg_strings`?**
  _High betweenness centrality (0.036) - this node is a cross-community bridge._
- **What connects `signing`, `requestIDKey`, `Error` to the rest of the system?**
  _63 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `GitHub & CI Workflow Skill` be split into smaller, more focused modules?**
  _Cohesion score 0.1111111111111111 - nodes in this community are weakly interconnected._
- **Why does `MockServer` connect `MockServer` to `go_pkg_net_http`, `smtptest_test.go`, `mock.go`, `smtptest.go`?**
  _High betweenness centrality (0.021) - this node is a cross-community bridge._
- **Should `scrub.go` be split into smaller, more focused modules?**
  _Cohesion score 0.0853302162478083 - nodes in this community are weakly interconnected._
- **Why does `Refs` connect `go_pkg_strings` to `logging_test.go`, `scrub.go`?**
  _High betweenness centrality (0.021) - this node is a cross-community bridge._
- **Should `cantonid.go` be split into smaller, more focused modules?**
  _Cohesion score 0.06676342525399129 - nodes in this community are weakly interconnected._