# Graph Report - daml-escrow-commons  (2026-11-18)

## Corpus Check
- 143 files · ~64,774 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 9 file(s) not represented in the graph (top: .ttf 5, (none) 4)

## Summary
- 988 nodes · 2517 edges · 58 communities (52 shown, 6 thin omitted)
- Extraction: 90% EXTRACTED · 10% INFERRED · 0% AMBIGUOUS · INFERRED: 250 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `430b3334`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- shortid.go
- orgsetup.go
- widget.json
- signer
- scrub.go
- logging_test.go
- Releasing
- go_pkg_fmt
- CLAUDE.md
- /commons-contribution
- daml-escrow-commons
- .Verify
- install-git-hooks.sh
- github.com/vdatacloud/daml-escrow-commons
- authority.go
- auditref.go
- apierror_test.go
- next-version.sh
- test-next-version.sh
- apierror.go
- field.go
- auditref_test.go
- canton.go
- approval_test.go
- person_test.go
- telemetry_test.go
- testing.T
- smtptest.go
- Git and CI
- cantonid.go
- time.Time
- pull_request_template.md
- smtptest_test.go
- orgsetup_test.go
- deliver.go
- canton_test.go
- Issuer
- agreementsig_test.go
- Repositories
- notify_test.go
- legalorders/legalorders.go
- http_test.go
- encoding/json.RawMessage
- go_pkg_strings
- go_pkg_testing
- field_test.go
- erasure.go
- audit.go
- fields.go
- go_pkg_net_http
- logging.go
- LoadDirectory
- go_pkg_encoding_json
- signin.go
- render.go
- go_pkg_time

## God Nodes (most connected - your core abstractions)
1. `Issuer` - 28 edges
2. `Refs` - 25 edges
3. `New()` - 16 edges
4. `MockServer` - 16 edges
5. `Principal` - 15 edges
6. `Kind` - 14 edges
7. `Fingerprint` - 13 edges
8. `New()` - 13 edges
9. `New()` - 12 edges
10. `PartyID` - 12 edges

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

## Communities (58 total, 6 thin omitted)

### Community 0 - "shortid.go"
Cohesion: 0.29
Nodes (7): TestShortID(), Field, decodeInto(), Error, ShortID(), ShortIDs(), shortValue()

### Community 1 - "orgsetup.go"
Cohesion: 0.60
Nodes (3): Check, Setup, Status

### Community 2 - "widget.json"
Cohesion: 0.17
Nodes (11): minLength, type, properties, name, quantity, minimum, type, required (+3 more)

### Community 3 - "signer"
Cohesion: 0.36
Nodes (5): signer(), TestOIDCVerifierReadsClaims(), TestOIDCVerifierRefuses(), token(), NewOIDCVerifierWithKeys()

### Community 4 - "scrub.go"
Cohesion: 0.09
Nodes (44): keyInfo, Rule, RuleFunc, TestServiceRules(), DefaultRules(), isRef(), keyOf(), KeyRedact() (+36 more)

### Community 5 - "logging_test.go"
Cohesion: 0.07
Nodes (48): Newf(), newRequestID(), RequestID(), RequestIDFrom(), requestIDKey, TestActorMiddleware(), TestRequireScope(), claims (+40 more)

### Community 6 - "Releasing"
Cohesion: 0.18
Nodes (10): 1. Prerequisites (one-time, per machine that will `go get` this module), 2. Decide the version bump, 3. Tag and push, 4. Publish the GitHub release, 5. Update consumers, Automated (default), Future automation, Major releases (+2 more)

### Community 7 - "go_pkg_fmt"
Cohesion: 0.11
Nodes (22): ChargeBearer, CommandOutcome, LedgerCommandEvent, requireBaseUnits(), requireDecimal(), NetworkFeeEvent, NetworkFeePayer, Rail (+14 more)

### Community 8 - "CLAUDE.md"
Cohesion: 0.22
Nodes (7): Architecture, Commands, Conventions carried over from `daml-escrow` / `daml-escrow-cms`, Cross-repo context, graphify, Project Overview, Relationship to the other repos

### Community 9 - "/commons-contribution"
Cohesion: 0.33
Nodes (5): /commons-contribution, How to add something, Removing something, When NOT to add something, When to promote code into commons

### Community 10 - "daml-escrow-commons"
Cohesion: 0.33
Nodes (5): daml-escrow-commons, Status, Using this module, What's here, What's NOT here, and won't be

### Community 11 - ".Verify"
Cohesion: 0.43
Nodes (6): fakeVerifier, Sign(), TestSignAndVerify_RoundTrip(), TestVerify_MalformedHexFails(), TestVerify_TamperedMessageFails(), TestVerify_WrongSecretFails()

### Community 14 - "authority.go"
Cohesion: 0.54
Nodes (7): Auth(), Handler(), Logging(), NewReport(), TestReport(), Report, StandIn

### Community 15 - "auditref.go"
Cohesion: 0.20
Nodes (5): Compute(), Kind, Refs, cacheKey, Config

### Community 16 - "apierror_test.go"
Cohesion: 0.18
Nodes (17): New(), Parse(), golden(), init(), TestCanonical_Golden(), TestCanonical_Properties(), TestDetailsInto(), TestError_IsAGoError() (+9 more)

### Community 18 - "test-next-version.sh"
Cohesion: 0.70
Nodes (4): c(), expect(), test-next-version.sh script, tg()

### Community 19 - "apierror.go"
Cohesion: 0.16
Nodes (9): CodeForStatus(), encode(), Error, sortedJSON(), TestCodeForStatus(), ValidCode(), Detail, Upstream (+1 more)

### Community 20 - "field.go"
Cohesion: 0.24
Nodes (10): ErrUnknownField, child(), contains(), Registry, lookupRef(), mergeProperty(), resolveNode(), splitPointer() (+2 more)

### Community 21 - "auditref_test.go"
Cohesion: 0.36
Nodes (9): New(), Parse(), TestAllAndParse(), TestInstitution(), TestKeyed(), TestKeyedFailsClosed(), TestPlainInDevelopment(), TestRefsCache() (+1 more)

### Community 23 - "canton.go"
Cohesion: 0.32
Nodes (9): ExpectedHash(), ForParty(), FromError(), FromLedgerError(), GRPCCodeName(), Signing(), TestSigningDetails(), LedgerError (+1 more)

### Community 24 - "approval_test.go"
Cohesion: 0.26
Nodes (11): Check(), Decisions(), cases(), generate(), TestCheckRefusals(), TestDecisions(), TestDecisionsGolden(), Decision (+3 more)

### Community 26 - "person_test.go"
Cohesion: 0.28
Nodes (6): personKind(), refLogger(), TestNewWithRefs(), TestPersonKind(), TestRefScrubCoreClear(), TestRefScrubCoreDevelopment()

### Community 27 - "telemetry_test.go"
Cohesion: 0.08
Nodes (18): Config, failing, Middleware(), Transport(), isURL(), metricOptions(), Propagator(), Setup() (+10 more)

### Community 28 - "testing.T"
Cohesion: 0.22
Nodes (13): TestErasureRequestWireNames(), TestLegalOrderWireNames(), TestSigninWireNames(), RetentionPeriod(), TestLedgerCommandEvent_Timing(), TestLedgerCommandEvent_Traffic(), TestLedgerCommandEvent_Validate(), TestNetworkFeeEvent_SpotPrice() (+5 more)

### Community 29 - "smtptest.go"
Cohesion: 0.05
Nodes (23): actorSlot, Channel, refusal(), WithHTTPClient(), HTTPOption, HTTPSender, Message, cloneMessage() (+15 more)

### Community 30 - "Git and CI"
Cohesion: 0.40
Nodes (4): Branch protection, Doing it, Git and CI, Policy

### Community 31 - "cantonid.go"
Cohesion: 0.10
Nodes (19): Classify(), FingerprintOf(), Fingerprint, Hash, PartyID, MustParsePartyID(), ParseFingerprint(), ParseHash() (+11 more)

### Community 32 - "time.Time"
Cohesion: 0.21
Nodes (8): LegalOrder, LegalOrderEvent, LegalOrderTarget, LiftLegalOrderResponse, ListLegalOrdersRequest, OrderContact, RecordedLegalOrder, Renderer

### Community 33 - "pull_request_template.md"
Cohesion: 0.40
Nodes (4): Architecture impact, Security, Summary, Tests

### Community 34 - "smtptest_test.go"
Cohesion: 0.31
Nodes (8): NewMockServer(), start(), TestMockServer_FailWith(), TestMockServer_Protocol(), TestMockServer_RecordsWhatNetSMTPSends(), TestMockServer_SinglePart(), TestMockServer_TooBig(), TestNewMockServer_DevOnly()

### Community 35 - "orgsetup_test.go"
Cohesion: 0.50
Nodes (3): PathFor(), TestFinish(), TestPathFor()

### Community 36 - "deliver.go"
Cohesion: 0.32
Nodes (9): Artifact, checkName(), Deliver(), artifacts(), TestDeliver_ManyIsDeterministicZip(), TestDeliver_One(), TestDeliver_RefusesBadNames(), Zip() (+1 more)

### Community 37 - "canton_test.go"
Cohesion: 0.11
Nodes (11): TestSchemaIsACopy(), ParseLedgerError(), TestFromError(), TestFromLedgerError_Golden(), TestFromLedgerError_Mapping(), TestLedgerError(), Schema(), fixture (+3 more)

### Community 38 - "Issuer"
Cohesion: 0.09
Nodes (11): Write(), Writer, authCode, Issuer, codeStore, Issuer, httpHandlerFunc(), writeJSON() (+3 more)

### Community 39 - "agreementsig_test.go"
Cohesion: 0.07
Nodes (41): AgreementCanonical(), AgreementHash(), AgreementMessage(), checkID(), invalid(), agreementCases(), agreementJSON(), flipS() (+33 more)

### Community 41 - "notify_test.go"
Cohesion: 0.27
Nodes (11): NewMockSender(), IsDevelopment(), TestIsDevelopment(), TestMockSender_RecordsOncePerIdempotencyKey(), TestMockSender_SentIsACopy(), TestMockSender_ValidatesAndFails(), TestNewMockSender_RefusedOutsideDevelopment(), TestValidate_AcceptsAWellFormedMessage() (+3 more)

### Community 42 - "legalorders/legalorders.go"
Cohesion: 0.12
Nodes (14): CaseRequest, EventRequest, TestEventAndLiftValidate(), TestRecordRequestValidate(), TestWireNames(), valid(), ValidCase(), validNotifyAfter() (+6 more)

### Community 43 - "http_test.go"
Cohesion: 0.36
Nodes (8): NewHTTPSender(), staticToken(), TestHTTPSender_RefusesAnInvalidMessageWithoutCalling(), TestHTTPSender_SendsAndReadsTheReceipt(), TestHTTPSender_Status(), TestHTTPSender_StatusAnswersToTheRelay(), TestHTTPSender_TransportFailuresAreRetryable(), TestNewHTTPSender_Refusals()

### Community 44 - "encoding/json.RawMessage"
Cohesion: 0.43
Nodes (4): draftVector, Case, Detail, Event

### Community 45 - "go_pkg_strings"
Cohesion: 0.19
Nodes (11): TestClientPolicyFilter(), TestMiddleware(), TestParseClientPolicy(), ClientPolicy, WithClientPolicy(), WithDevBypass(), WithLogger(), Option (+3 more)

### Community 47 - "go_pkg_testing"
Cohesion: 0.22
Nodes (7): benchCore(), benchLine(), BenchmarkRawZap(), BenchmarkShapingCore(), benchRefs(), TestShapeLeavesCallerSlice(), TestShapingAllocs()

### Community 48 - "field_test.go"
Cohesion: 0.25
Nodes (7): TestSplitOrigins(), loadFields(), TestField_Errors(), TestField_KeepsOwnDescriptionAndRefKeywords(), TestField_RefFailures(), TestField_Resolves(), TestField_Root()

### Community 50 - "erasure.go"
Cohesion: 0.29
Nodes (6): ErasureRequest, ErasureRequestID, ListErasureRequests, RefuseErasureRequest, SubmitErasureRequest, VerifyErasureRequest

### Community 52 - "audit.go"
Cohesion: 0.27
Nodes (9): ActiveHolds, AuditIdentity, AuditRef, AuditRefsOf, LegalHold, Resolution, ResolveAuditRefs, ResolvedRef (+1 more)

### Community 53 - "fields.go"
Cohesion: 0.31
Nodes (7): accountNumber, Account(), Party(), Subject(), User(), personID, TestTypedFields()

### Community 54 - "go_pkg_net_http"
Cohesion: 0.07
Nodes (26): TestAuth(), main(), splitOrigins(), authorizeQuery(), browserIssuer(), challengeOf(), exchange(), signIn() (+18 more)

### Community 55 - "logging.go"
Cohesion: 0.43
Nodes (6): Config, build(), New(), TestNew(), TestNewIsNotSampled(), zapConfig()

### Community 56 - "LoadDirectory"
Cohesion: 0.48
Nodes (6): LoadDirectory(), TestLoadDirectory_CompilesSchemas(), TestLoadDirectory_MissingDirectory(), TestValidate_InvalidPayloadReportsFailures(), TestValidate_UnknownType(), TestValidate_ValidPayload()

### Community 57 - "go_pkg_encoding_json"
Cohesion: 0.12
Nodes (19): TestWireNames(), Client, Identity, New(), TestClient_CreatePendingInvitation(), TestClient_CreatePendingInvitation_ServerError(), TestClient_GetByEmail_ServerError(), TestClient_GetByOktaSub_Found() (+11 more)

### Community 58 - "signin.go"
Cohesion: 0.40
Nodes (5): SigninActionRequest, SigninActionResult, SigninFactor, SigninStatus, SigninStatusRequest

### Community 59 - "render.go"
Cohesion: 0.08
Nodes (21): TestDraftCanonical_EmptyTermsAreNull(), TestDraftVectors(), DraftCanonical(), DraftHash(), DraftMessage(), DraftVersion, Fake, CheckFonts() (+13 more)

### Community 63 - "go_pkg_time"
Cohesion: 0.40
Nodes (3): Response, Status, TestStage()

## Knowledge Gaps
- **67 isolated node(s):** `signing`, `requestIDKey`, `Error`, `ctxKey`, `github.com/vdatacloud/daml-escrow-commons` (+62 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 175 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **6 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Issuer` connect `Issuer` to `shortid.go`, `agreementsig_test.go`, `go_pkg_fmt`, `go_pkg_strings`, `go_pkg_net_http`, `go_pkg_encoding_json`, `render.go`, `go_pkg_time`?**
  _High betweenness centrality (0.035) - this node is a cross-community bridge._
- **What connects `signing`, `requestIDKey`, `Error` to the rest of the system?**
  _67 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `scrub.go` be split into smaller, more focused modules?**
  _Cohesion score 0.08771929824561403 - nodes in this community are weakly interconnected._
- **Why does `MockServer` connect `smtptest.go` to `smtptest_test.go`, `go_pkg_net_http`?**
  _High betweenness centrality (0.020) - this node is a cross-community bridge._
- **Should `logging_test.go` be split into smaller, more focused modules?**
  _Cohesion score 0.07322068612391193 - nodes in this community are weakly interconnected._
- **Why does `Refs` connect `auditref.go` to `scrub.go`, `logging_test.go`, `go_pkg_testing`, `auditref_test.go`, `logging.go`, `person_test.go`?**
  _High betweenness centrality (0.020) - this node is a cross-community bridge._
- **Should `go_pkg_fmt` be split into smaller, more focused modules?**
  _Cohesion score 0.10810810810810811 - nodes in this community are weakly interconnected._