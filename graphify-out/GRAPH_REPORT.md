# Graph Report - daml-escrow-commons  (2026-11-09)

## Corpus Check
- 131 files · ~61,125 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 9 file(s) not represented in the graph (top: .ttf 5, (none) 4)

## Summary
- 932 nodes · 2378 edges · 60 communities (54 shown, 6 thin omitted)
- Extraction: 90% EXTRACTED · 10% INFERRED · 0% AMBIGUOUS · INFERRED: 238 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `85353fd4`
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
- go_pkg_fmt
- CLAUDE.md
- /commons-contribution
- daml-escrow-commons
- .claude/CLAUDE.md
- install-git-hooks.sh
- github.com/vdatacloud/daml-escrow-commons
- agreement.go
- telemetry.go
- notify.go
- next-version.sh
- test-next-version.sh
- shortid.go
- field.go
- apierror.go
- canton_test.go
- approval_test.go
- context.Context
- go_pkg_net_http
- apierror_test.go
- smtptest.go
- authn_test.go
- go_pkg_strings
- notify/http.go
- MockServer
- smtptest_test.go
- mock.go
- go_pkg_encoding_json
- main.go
- Issuer
- agreementsig_test.go
- middleware.go
- notify_test.go
- legalorders/legalorders.go
- http_test.go
- .Verify
- auditref_test.go
- authorize_test.go
- logging_test.go
- ClientPolicy
- identityclient.go
- TestAgreementVectors
- oidc.go
- pathparam_test.go
- fields.go
- go_pkg_net_http_httptest
- logging.go
- bench_test.go
- field_test.go
- LoadDirectory
- go_pkg_testing

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

## Communities (60 total, 6 thin omitted)

### Community 0 - "GitHub & CI Workflow Skill"
Cohesion: 0.11
Nodes (17): 10. Useful Reference Commands, 11. References, 1. Pre-Commit Local Verification (MANDATORY — do this before every commit), 2. Branching Rules, 3. Commit Standards, 4. Staging & Pushing Changes, 5. Pull Request Creation, 6. CI Pipeline Overview (+9 more)

### Community 1 - "testing.T"
Cohesion: 0.22
Nodes (17): New(), TestClient_CreatePendingInvitation(), TestClient_CreatePendingInvitation_ServerError(), TestClient_GetByEmail_ServerError(), TestClient_GetByOktaSub_Found(), TestClient_GetByToken_NotFound(), TestClient_ManagesIdentity(), TestClient_ManagesIdentity_ServerError() (+9 more)

### Community 2 - "widget.json"
Cohesion: 0.17
Nodes (11): minLength, type, properties, name, quantity, minimum, type, required (+3 more)

### Community 3 - "Repository Guardrails"
Cohesion: 0.22
Nodes (8): Branch Protection Strategy, Branching Strategy, CI Requirements, Code Review Requirements, Commit Standard, Pre-Commit Verification, Pull Request Rules, Repository Guardrails

### Community 4 - "scrub.go"
Cohesion: 0.06
Nodes (50): Compute(), Kind, Refs, TestRefsCache(), cacheKey, keyInfo, personID, Rule (+42 more)

### Community 5 - "cantonid.go"
Cohesion: 0.10
Nodes (19): Classify(), FingerprintOf(), Fingerprint, Hash, PartyID, MustParsePartyID(), ParseFingerprint(), ParseHash() (+11 more)

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
Cohesion: 0.15
Nodes (20): AgreementCanonical(), AgreementMessage(), checkID(), invalid(), TestDraftVectors(), AgreementVersion, Amends, decodeJSON() (+12 more)

### Community 15 - "telemetry.go"
Cohesion: 0.15
Nodes (8): Config, isURL(), metricOptions(), Propagator(), Setup(), TestSetup(), TestSetupSampleRatio(), traceOptions()

### Community 16 - "notify.go"
Cohesion: 0.27
Nodes (5): Channel, Message, Recipient, Sender, Status

### Community 18 - "test-next-version.sh"
Cohesion: 0.70
Nodes (4): c(), expect(), test-next-version.sh script, tg()

### Community 19 - "shortid.go"
Cohesion: 0.29
Nodes (7): TestShortID(), Field, decodeInto(), Error, ShortID(), ShortIDs(), shortValue()

### Community 20 - "field.go"
Cohesion: 0.24
Nodes (10): ErrUnknownField, child(), contains(), Registry, lookupRef(), mergeProperty(), resolveNode(), splitPointer() (+2 more)

### Community 21 - "apierror.go"
Cohesion: 0.16
Nodes (9): CodeForStatus(), encode(), Error, sortedJSON(), TestCodeForStatus(), ValidCode(), Detail, Upstream (+1 more)

### Community 23 - "canton_test.go"
Cohesion: 0.17
Nodes (16): ExpectedHash(), ForParty(), FromError(), FromLedgerError(), GRPCCodeName(), ParseLedgerError(), Signing(), TestFromError() (+8 more)

### Community 24 - "approval_test.go"
Cohesion: 0.29
Nodes (11): Check(), Decisions(), cases(), generate(), TestCheckRefusals(), TestDecisions(), TestDecisionsGolden(), Decision (+3 more)

### Community 26 - "context.Context"
Cohesion: 0.20
Nodes (20): RequestIDFrom(), TestActorMiddleware(), ActorMiddleware(), FromContext(), HasScope(), AccessLogOption, actorKey, loggerKey (+12 more)

### Community 27 - "go_pkg_net_http"
Cohesion: 0.12
Nodes (10): failing, Middleware(), Transport(), attrs(), recording(), TestMiddleware(), TestMiddlewareServerError(), TestTransportError() (+2 more)

### Community 28 - "apierror_test.go"
Cohesion: 0.09
Nodes (21): New(), Parse(), golden(), init(), TestCanonical_Golden(), TestCanonical_Properties(), TestDetailsInto(), TestError_IsAGoError() (+13 more)

### Community 29 - "smtptest.go"
Cohesion: 0.19
Nodes (6): readPart(), WithMaxSize(), WithOnMessage(), Content, Message, Option

### Community 30 - "authn_test.go"
Cohesion: 0.22
Nodes (5): refLogger(), TestRefScrubCore(), TestRefScrubCoreClear(), TestRefScrubCoreDevelopment(), TestShapeLeavesCallerSlice()

### Community 31 - "go_pkg_strings"
Cohesion: 0.05
Nodes (34): LegalOrder, LegalOrderEvent, LegalOrderTarget, LiftLegalOrderResponse, ListLegalOrdersRequest, OrderContact, RecordedLegalOrder, personKind() (+26 more)

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

### Community 36 - "go_pkg_encoding_json"
Cohesion: 0.33
Nodes (4): newRequestID(), RequestID(), requestIDKey, TestLegalOrderWireNames()

### Community 37 - "main.go"
Cohesion: 0.18
Nodes (9): main(), splitOrigins(), TestSplitOrigins(), New(), TestDevOIDC_ClientCredentials(), TestDevOIDC_SigningSessionClaims(), TestDevOIDC_TokenRequiresSubAndEmail(), TestDevOIDC_TokenVerifiesWithGoOIDC() (+1 more)

### Community 38 - "Issuer"
Cohesion: 0.08
Nodes (11): Write(), Writer, authCode, Issuer, codeStore, Issuer, httpHandlerFunc(), writeJSON() (+3 more)

### Community 39 - "agreementsig_test.go"
Cohesion: 0.13
Nodes (18): flipS(), sha256Sum(), signatureCases(), spki(), TestDraftCanonical_EmptyTermsAreNull(), testKeys(), TestNormalizeECDSA(), TestVerify_KeyRefusals() (+10 more)

### Community 40 - "middleware.go"
Cohesion: 0.18
Nodes (15): Newf(), TestMiddleware(), TestRequireScope(), ctxKey, fakeVerifier, devPrincipal(), Middleware(), RequireScope() (+7 more)

### Community 41 - "notify_test.go"
Cohesion: 0.27
Nodes (11): NewMockSender(), IsDevelopment(), TestIsDevelopment(), TestMockSender_RecordsOncePerIdempotencyKey(), TestMockSender_SentIsACopy(), TestMockSender_ValidatesAndFails(), TestNewMockSender_RefusedOutsideDevelopment(), TestValidate_AcceptsAWellFormedMessage() (+3 more)

### Community 42 - "legalorders/legalorders.go"
Cohesion: 0.17
Nodes (10): CaseRequest, EventRequest, ValidCase(), validNotifyAfter(), Lifted, LiftRequest, Recorded, RecordedEvent (+2 more)

### Community 43 - "http_test.go"
Cohesion: 0.31
Nodes (8): NewHTTPSender(), staticToken(), TestHTTPSender_RefusesAnInvalidMessageWithoutCalling(), TestHTTPSender_SendsAndReadsTheReceipt(), TestHTTPSender_Status(), TestHTTPSender_StatusAnswersToTheRelay(), TestHTTPSender_TransportFailuresAreRetryable(), TestNewHTTPSender_Refusals()

### Community 44 - ".Verify"
Cohesion: 0.22
Nodes (10): signer(), TestOIDCVerifierReadsClaims(), TestOIDCVerifierRefuses(), token(), NewOIDCVerifierWithKeys(), Sign(), TestSignAndVerify_RoundTrip(), TestVerify_MalformedHexFails() (+2 more)

### Community 45 - "auditref_test.go"
Cohesion: 0.26
Nodes (10): New(), Parse(), TestAllAndParse(), TestInstitution(), TestKeyed(), TestKeyedFailsClosed(), TestPlainInDevelopment(), TestSetClear() (+2 more)

### Community 46 - "authorize_test.go"
Cohesion: 0.38
Nodes (8): authorizeQuery(), browserIssuer(), challengeOf(), exchange(), signIn(), TestDevOIDC_AuthorizeRefusals(), TestDevOIDC_BrowserPortalSession(), TestDevOIDC_BrowserSigningSession()

### Community 47 - "logging_test.go"
Cohesion: 0.36
Nodes (9): fieldMap(), observed(), spanCtx(), TestAccessLog(), TestAccessLogServerError(), TestFromContext(), TestScrubCore(), TestScrubCoreWith() (+1 more)

### Community 48 - "ClientPolicy"
Cohesion: 0.25
Nodes (6): TestClientPolicyFilter(), TestParseClientPolicy(), ClientPolicy, options, NewClientPolicy(), ParseClientPolicy()

### Community 49 - "identityclient.go"
Cohesion: 0.36
Nodes (3): Client, Identity, PendingInvitation

### Community 50 - "TestAgreementVectors"
Cohesion: 0.29
Nodes (10): AgreementHash(), agreementCases(), agreementJSON(), readJSON(), TestAgreementHash_ChangesWithEverySignedField(), TestAgreementMessage_Refusals(), TestAgreementVectors(), TestAgreementVersion_JSONRoundTrip() (+2 more)

### Community 51 - "oidc.go"
Cohesion: 0.27
Nodes (5): claims, NewOIDCVerifier(), uniqueScopes(), OIDCVerifier, TestStartAndMintVerifyWithAuthn()

### Community 52 - "pathparam_test.go"
Cohesion: 0.36
Nodes (5): Get(), containsByte(), jsEncodeURIComponent(), TestGetKeepsInvalidEscaping(), TestGetUnescapes()

### Community 53 - "fields.go"
Cohesion: 0.42
Nodes (6): accountNumber, Account(), Party(), Subject(), User(), TestTypedFields()

### Community 54 - "go_pkg_net_http_httptest"
Cohesion: 0.29
Nodes (4): Issuer, Start(), TestSeveralOrigins(), Token

### Community 55 - "logging.go"
Cohesion: 0.43
Nodes (6): Config, build(), New(), TestNew(), TestNewIsNotSampled(), zapConfig()

### Community 56 - "bench_test.go"
Cohesion: 0.54
Nodes (6): benchCore(), benchLine(), BenchmarkRawZap(), BenchmarkShapingCore(), benchRefs(), TestShapingAllocs()

### Community 57 - "field_test.go"
Cohesion: 0.39
Nodes (6): loadFields(), TestField_Errors(), TestField_KeepsOwnDescriptionAndRefKeywords(), TestField_RefFailures(), TestField_Resolves(), TestField_Root()

### Community 58 - "LoadDirectory"
Cohesion: 0.48
Nodes (6): LoadDirectory(), TestLoadDirectory_CompilesSchemas(), TestLoadDirectory_MissingDirectory(), TestValidate_InvalidPayloadReportsFailures(), TestValidate_UnknownType(), TestValidate_ValidPayload()

### Community 59 - "go_pkg_testing"
Cohesion: 0.40
Nodes (4): TestEventAndLiftValidate(), TestRecordRequestValidate(), TestWireNames(), valid()

## Knowledge Gaps
- **68 isolated node(s):** `signing`, `requestIDKey`, `Error`, `ctxKey`, `github.com/vdatacloud/daml-escrow-commons` (+63 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 174 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **6 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Issuer` connect `Issuer` to `go_pkg_encoding_json`, `main.go`, `go_pkg_fmt`, `agreementsig_test.go`, `shortid.go`, `go_pkg_net_http`, `authn_test.go`, `go_pkg_strings`?**
  _High betweenness centrality (0.035) - this node is a cross-community bridge._
- **What connects `signing`, `requestIDKey`, `Error` to the rest of the system?**
  _68 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `GitHub & CI Workflow Skill` be split into smaller, more focused modules?**
  _Cohesion score 0.1111111111111111 - nodes in this community are weakly interconnected._
- **Why does `MockServer` connect `MockServer` to `smtptest_test.go`, `mock.go`, `smtptest.go`, `telemetry.go`?**
  _High betweenness centrality (0.021) - this node is a cross-community bridge._
- **Should `scrub.go` be split into smaller, more focused modules?**
  _Cohesion score 0.06456140350877193 - nodes in this community are weakly interconnected._
- **Why does `Refs` connect `scrub.go` to `auditref_test.go`, `logging.go`, `bench_test.go`, `context.Context`, `authn_test.go`, `go_pkg_strings`?**
  _High betweenness centrality (0.020) - this node is a cross-community bridge._
- **Should `cantonid.go` be split into smaller, more focused modules?**
  _Cohesion score 0.09716599190283401 - nodes in this community are weakly interconnected._