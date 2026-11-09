# Graph Report - daml-escrow-commons  (2026-11-09)

## Corpus Check
- 139 files · ~63,075 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 9 file(s) not represented in the graph (top: .ttf 5, (none) 4)

## Summary
- 973 nodes · 2463 edges · 62 communities (56 shown, 6 thin omitted)
- Extraction: 90% EXTRACTED · 10% INFERRED · 0% AMBIGUOUS · INFERRED: 243 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `8adf11d2`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- GitHub & CI Workflow Skill
- apierror_test.go
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
- go_pkg_errors
- notify.go
- next-version.sh
- test-next-version.sh
- shortid.go
- field.go
- Error
- canton.go
- approval_test.go
- request.go
- telemetry_test.go
- testing.T
- smtptest.go
- logging_test.go
- go_pkg_testing
- go_pkg_net_http
- MockServer
- smtptest_test.go
- mock.go
- deliver.go
- main.go
- Issuer
- agreementsig_test.go
- go_pkg_context
- notify_test.go
- legalorders/legalorders.go
- http_test.go
- authn_test.go
- authority.go
- time.Time
- Principal
- middleware.go
- context.Context
- go_pkg_strings
- Str
- signatureCases
- fields.go
- authorize_test.go
- logging.go
- LoadDirectory
- identityclient_test.go
- Hash
- metering_test.go
- .Verify
- apierror/http.go

## God Nodes (most connected - your core abstractions)
1. `Issuer` - 28 edges
2. `Refs` - 25 edges
3. `New()` - 16 edges
4. `MockServer` - 16 edges
5. `Principal` - 15 edges
6. `Kind` - 14 edges
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

## Communities (62 total, 6 thin omitted)

### Community 0 - "GitHub & CI Workflow Skill"
Cohesion: 0.11
Nodes (17): 10. Useful Reference Commands, 11. References, 1. Pre-Commit Local Verification (MANDATORY — do this before every commit), 2. Branching Rules, 3. Commit Standards, 4. Staging & Pushing Changes, 5. Pull Request Creation, 6. CI Pipeline Overview (+9 more)

### Community 1 - "apierror_test.go"
Cohesion: 0.14
Nodes (12): CodeForStatus(), Parse(), TestCodeForStatus(), TestDetailsInto(), TestFromResponse(), TestParse(), ValidCode(), FromResponse() (+4 more)

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
Cohesion: 0.15
Nodes (14): Classify(), FingerprintOf(), Fingerprint, PartyID, MustParsePartyID(), ParseFingerprint(), ParsePartyID(), Short() (+6 more)

### Community 6 - "Releasing"
Cohesion: 0.18
Nodes (10): 1. Prerequisites (one-time, per machine that will `go get` this module), 2. Decide the version bump, 3. Tag and push, 4. Publish the GitHub release, 5. Update consumers, Automated (default), Future automation, Major releases (+2 more)

### Community 7 - "metering.go"
Cohesion: 0.15
Nodes (19): ChargeBearer, CommandOutcome, LedgerCommandEvent, requireBaseUnits(), requireDecimal(), NetworkFeeEvent, NetworkFeePayer, Rail (+11 more)

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

### Community 15 - "go_pkg_errors"
Cohesion: 0.07
Nodes (26): Compute(), Kind, Refs, New(), Parse(), TestAllAndParse(), TestInstitution(), TestKeyed() (+18 more)

### Community 16 - "notify.go"
Cohesion: 0.16
Nodes (8): Channel, refusal(), HTTPSender, Message, Receipt, Recipient, Sender, Status

### Community 18 - "test-next-version.sh"
Cohesion: 0.70
Nodes (4): c(), expect(), test-next-version.sh script, tg()

### Community 19 - "shortid.go"
Cohesion: 0.21
Nodes (10): encode(), sortedJSON(), TestShortID(), Field, decodeInto(), Error, ShortID(), ShortIDs() (+2 more)

### Community 20 - "field.go"
Cohesion: 0.24
Nodes (10): ErrUnknownField, child(), contains(), Registry, lookupRef(), mergeProperty(), resolveNode(), splitPointer() (+2 more)

### Community 21 - "Error"
Cohesion: 0.21
Nodes (3): Error, Detail, Upstream

### Community 23 - "canton.go"
Cohesion: 0.18
Nodes (18): New(), init(), TestError_IsAGoError(), TestWithUpstream_BoundsCause(), ExpectedHash(), ForParty(), FromError(), FromLedgerError() (+10 more)

### Community 24 - "approval_test.go"
Cohesion: 0.29
Nodes (11): Check(), Decisions(), cases(), generate(), TestCheckRefusals(), TestDecisions(), TestDecisionsGolden(), Decision (+3 more)

### Community 26 - "request.go"
Cohesion: 0.19
Nodes (18): TestActorMiddleware(), ActorMiddleware(), WithPrincipal(), AccessLogOption, actorKey, loggerKey, TestActorSlot(), AccessLog() (+10 more)

### Community 27 - "telemetry_test.go"
Cohesion: 0.07
Nodes (18): Config, failing, Middleware(), Transport(), isURL(), metricOptions(), Propagator(), Setup() (+10 more)

### Community 28 - "testing.T"
Cohesion: 0.13
Nodes (22): golden(), TestCanonical_Golden(), TestCanonical_Properties(), TestLogFields_Golden(), TestLogFields_ShortFormsOnly(), TestSchemaIsACopy(), TestWriter(), validate() (+14 more)

### Community 29 - "smtptest.go"
Cohesion: 0.17
Nodes (8): readPart(), TestMockServer_RecordsWhatNetSMTPSends(), TestMockServer_TooBig(), WithMaxSize(), WithOnMessage(), Content, Message, Option

### Community 30 - "logging_test.go"
Cohesion: 0.19
Nodes (14): fieldMap(), observed(), spanCtx(), TestAccessLog(), TestAccessLogServerError(), TestFromContext(), TestScrubCore(), TestScrubCoreWith() (+6 more)

### Community 31 - "go_pkg_testing"
Cohesion: 0.14
Nodes (8): TestAuth(), TestSplitOrigins(), Issuer, Start(), TestNewDefaultsToThePortal(), TestSeveralOrigins(), TestStartAndMintVerifyWithAuthn(), Token

### Community 32 - "go_pkg_net_http"
Cohesion: 0.31
Nodes (5): Get(), containsByte(), jsEncodeURIComponent(), TestGetKeepsInvalidEscaping(), TestGetUnescapes()

### Community 33 - "MockServer"
Cohesion: 0.19
Nodes (3): path(), readData(), MockServer

### Community 34 - "smtptest_test.go"
Cohesion: 0.27
Nodes (6): NewMockServer(), start(), TestMockServer_FailWith(), TestMockServer_Protocol(), TestMockServer_SinglePart(), TestNewMockServer_DevOnly()

### Community 35 - "mock.go"
Cohesion: 0.36
Nodes (3): actorSlot, cloneMessage(), MockSender

### Community 36 - "deliver.go"
Cohesion: 0.29
Nodes (9): Artifact, checkName(), Deliver(), artifacts(), TestDeliver_ManyIsDeterministicZip(), TestDeliver_One(), TestDeliver_RefusesBadNames(), Zip() (+1 more)

### Community 37 - "main.go"
Cohesion: 0.17
Nodes (7): main(), splitOrigins(), New(), TestDevOIDC_ClientCredentials(), TestDevOIDC_SigningSessionClaims(), TestDevOIDC_TokenRequiresSubAndEmail(), TestDevOIDC_TokenVerifiesWithGoOIDC()

### Community 38 - "Issuer"
Cohesion: 0.09
Nodes (11): Write(), Writer, authCode, Issuer, codeStore, Issuer, httpHandlerFunc(), writeJSON() (+3 more)

### Community 39 - "agreementsig_test.go"
Cohesion: 0.12
Nodes (18): AgreementHash(), agreementCases(), agreementJSON(), TestAgreementHash_ChangesWithEverySignedField(), TestAgreementMessage_Refusals(), TestAgreementVectors(), TestAgreementVersion_JSONRoundTrip(), TestDraftCanonical_EmptyTermsAreNull() (+10 more)

### Community 40 - "go_pkg_context"
Cohesion: 0.23
Nodes (6): CheckFonts(), NewTypst(), TestCheckFonts(), TestTypst_ByteIdenticalAndEscaped(), TestTypst_GoldenHashAcrossPlatforms(), TestTypst_RefusesOtherVersion()

### Community 41 - "notify_test.go"
Cohesion: 0.27
Nodes (11): NewMockSender(), IsDevelopment(), TestIsDevelopment(), TestMockSender_RecordsOncePerIdempotencyKey(), TestMockSender_SentIsACopy(), TestMockSender_ValidatesAndFails(), TestNewMockSender_RefusedOutsideDevelopment(), TestValidate_AcceptsAWellFormedMessage() (+3 more)

### Community 42 - "legalorders/legalorders.go"
Cohesion: 0.17
Nodes (10): CaseRequest, EventRequest, ValidCase(), validNotifyAfter(), Lifted, LiftRequest, Recorded, RecordedEvent (+2 more)

### Community 43 - "http_test.go"
Cohesion: 0.22
Nodes (11): NewHTTPSender(), staticToken(), TestHTTPSender_RefusesAnInvalidMessageWithoutCalling(), TestHTTPSender_SendsAndReadsTheReceipt(), TestHTTPSender_Status(), TestHTTPSender_StatusAnswersToTheRelay(), TestHTTPSender_TransportFailuresAreRetryable(), TestNewHTTPSender_Refusals() (+3 more)

### Community 44 - "authn_test.go"
Cohesion: 0.33
Nodes (5): signer(), TestOIDCVerifierReadsClaims(), TestOIDCVerifierRefuses(), token(), NewOIDCVerifierWithKeys()

### Community 45 - "authority.go"
Cohesion: 0.19
Nodes (13): TestClientPolicyFilter(), TestParseClientPolicy(), ClientPolicy, options, NewClientPolicy(), ParseClientPolicy(), Auth(), Handler() (+5 more)

### Community 46 - "time.Time"
Cohesion: 0.05
Nodes (32): draftVector, fixture, Response, Status, ActiveHolds, AuditIdentity, AuditRef, AuditRefsOf (+24 more)

### Community 47 - "Principal"
Cohesion: 0.19
Nodes (10): claims, ctxKey, fakeVerifier, NewOIDCVerifier(), uniqueScopes(), OIDCVerifier, Principal, FromContext() (+2 more)

### Community 48 - "middleware.go"
Cohesion: 0.23
Nodes (10): Newf(), TestMiddleware(), TestRequireScope(), devPrincipal(), Middleware(), RequireScope(), WithClientPolicy(), WithDevBypass() (+2 more)

### Community 49 - "context.Context"
Cohesion: 0.42
Nodes (3): Client, Identity, PendingInvitation

### Community 50 - "go_pkg_strings"
Cohesion: 0.18
Nodes (3): TestRecordRequestValidate(), valid(), Errors

### Community 51 - "Str"
Cohesion: 0.70
Nodes (4): TestStr(), Content(), Lit(), Str()

### Community 52 - "signatureCases"
Cohesion: 0.29
Nodes (10): flipS(), readJSON(), sha256Sum(), signatureCases(), spki(), testKeys(), TestNormalizeECDSA(), TestSignatureVectors() (+2 more)

### Community 53 - "fields.go"
Cohesion: 0.36
Nodes (6): accountNumber, Account(), Party(), Subject(), User(), TestTypedFields()

### Community 54 - "authorize_test.go"
Cohesion: 0.38
Nodes (8): authorizeQuery(), browserIssuer(), challengeOf(), exchange(), signIn(), TestDevOIDC_AuthorizeRefusals(), TestDevOIDC_BrowserPortalSession(), TestDevOIDC_BrowserSigningSession()

### Community 55 - "logging.go"
Cohesion: 0.43
Nodes (6): Config, build(), New(), TestNew(), TestNewIsNotSampled(), zapConfig()

### Community 56 - "LoadDirectory"
Cohesion: 0.25
Nodes (7): Registry, LoadDirectory(), TestLoadDirectory_CompilesSchemas(), TestLoadDirectory_MissingDirectory(), TestValidate_InvalidPayloadReportsFailures(), TestValidate_UnknownType(), TestValidate_ValidPayload()

### Community 57 - "identityclient_test.go"
Cohesion: 0.36
Nodes (9): New(), TestClient_CreatePendingInvitation(), TestClient_CreatePendingInvitation_ServerError(), TestClient_GetByEmail_ServerError(), TestClient_GetByOktaSub_Found(), TestClient_GetByToken_NotFound(), TestClient_ManagesIdentity(), TestClient_ManagesIdentity_ServerError() (+1 more)

### Community 58 - "Hash"
Cohesion: 0.25
Nodes (4): Hash, ParseHash(), shortHex(), TestShortForms()

### Community 59 - "metering_test.go"
Cohesion: 0.28
Nodes (8): TestLedgerCommandEvent_Timing(), TestLedgerCommandEvent_Traffic(), TestLedgerCommandEvent_Validate(), TestNetworkFeeEvent_SpotPrice(), TestNetworkFeeEvent_Validate(), TestSettlementEvent_NetworkFee(), TestSettlementEvent_Validate(), withFee()

### Community 60 - ".Verify"
Cohesion: 0.52
Nodes (5): Sign(), TestSignAndVerify_RoundTrip(), TestVerify_MalformedHexFails(), TestVerify_TamperedMessageFails(), TestVerify_WrongSecretFails()

### Community 61 - "apierror/http.go"
Cohesion: 0.40
Nodes (4): newRequestID(), RequestID(), RequestIDFrom(), requestIDKey

## Knowledge Gaps
- **77 isolated node(s):** `signing`, `requestIDKey`, `Error`, `ctxKey`, `github.com/vdatacloud/daml-escrow-commons` (+72 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 184 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **6 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Issuer` connect `Issuer` to `go_pkg_net_http`, `apierror_test.go`, `main.go`, `agreementsig_test.go`, `go_pkg_errors`, `go_pkg_strings`, `shortid.go`, `apierror/http.go`, `go_pkg_testing`?**
  _High betweenness centrality (0.035) - this node is a cross-community bridge._
- **What connects `signing`, `requestIDKey`, `Error` to the rest of the system?**
  _77 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `GitHub & CI Workflow Skill` be split into smaller, more focused modules?**
  _Cohesion score 0.1111111111111111 - nodes in this community are weakly interconnected._
- **Why does `MockServer` connect `MockServer` to `telemetry_test.go`, `smtptest_test.go`, `mock.go`, `smtptest.go`?**
  _High betweenness centrality (0.020) - this node is a cross-community bridge._
- **Should `apierror_test.go` be split into smaller, more focused modules?**
  _Cohesion score 0.14 - nodes in this community are weakly interconnected._
- **Why does `Refs` connect `go_pkg_errors` to `request.go`, `scrub.go`, `logging_test.go`, `logging.go`?**
  _High betweenness centrality (0.020) - this node is a cross-community bridge._
- **Should `scrub.go` be split into smaller, more focused modules?**
  _Cohesion score 0.08771929824561403 - nodes in this community are weakly interconnected._