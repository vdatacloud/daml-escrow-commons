# Graph Report - daml-escrow-commons  (2026-11-09)

## Corpus Check
- 140 files · ~63,594 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 9 file(s) not represented in the graph (top: .ttf 5, (none) 4)

## Summary
- 981 nodes · 2475 edges · 65 communities (56 shown, 9 thin omitted)
- Extraction: 90% EXTRACTED · 10% INFERRED · 0% AMBIGUOUS · INFERRED: 243 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `0fa1a78c`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- GitHub & CI Workflow Skill
- Parse
- widget.json
- Repository Guardrails
- redact.go
- Error
- Releasing
- metering.go
- CLAUDE.md
- /commons-contribution
- daml-escrow-commons
- .claude/CLAUDE.md
- install-git-hooks.sh
- github.com/vdatacloud/daml-escrow-commons
- scrub.go
- Refs
- notify.go
- next-version.sh
- test-next-version.sh
- encode
- canonical.go
- New
- encoding/json.RawMessage
- approval_test.go
- request.go
- telemetry_test.go
- testing.T
- smtptest.go
- fieldMap
- Principal
- context.Context
- MockServer
- smtptest_test.go
- mock.go
- deliver.go
- go_pkg_strings
- Issuer
- agreementsig_test.go
- NewTypst
- notify_test.go
- legalorders/legalorders.go
- http_test.go
- .Verify
- TestMiddleware
- time.Time
- NewRefScrubCore
- field_test.go
- HTTPSender
- legalorders/legalorders_test.go
- go.uber.org/zap/zapcore.Field
- audit.go
- fields.go
- authorize_test.go
- build
- LoadDirectory
- identityclient_test.go
- OIDCVerifier
- Typst
- TestAccessLog
- erasure.go
- signin.go
- keys.go
- gastank.go

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
- `keyInfo` --references--> `Kind`  [EXTRACTED]
  logging/scrub.go → auditref/auditref.go
- `personID` --references--> `Kind`  [EXTRACTED]
  logging/fields.go → auditref/auditref.go
- `refOf()` --references--> `Kind`  [EXTRACTED]
  logging/scrub.go → auditref/auditref.go
- `ActorMiddleware()` --references--> `Refs`  [EXTRACTED]
  authn/middleware.go → auditref/auditref.go

## Import Cycles
- None detected.

## Communities (65 total, 9 thin omitted)

### Community 0 - "GitHub & CI Workflow Skill"
Cohesion: 0.11
Nodes (17): 10. Useful Reference Commands, 11. References, 1. Pre-Commit Local Verification (MANDATORY — do this before every commit), 2. Branching Rules, 3. Commit Standards, 4. Staging & Pushing Changes, 5. Pull Request Creation, 6. CI Pipeline Overview (+9 more)

### Community 1 - "Parse"
Cohesion: 0.17
Nodes (12): Parse(), golden(), TestCanonical_Golden(), TestDetailsInto(), TestFromResponse(), TestLogFields_Golden(), TestParse(), TestSchemaIsACopy() (+4 more)

### Community 2 - "widget.json"
Cohesion: 0.17
Nodes (11): minLength, type, properties, name, quantity, minimum, type, required (+3 more)

### Community 3 - "Repository Guardrails"
Cohesion: 0.22
Nodes (8): Branch Protection Strategy, Branching Strategy, CI Requirements, Code Review Requirements, Commit Standard, Pre-Commit Verification, Pull Request Rules, Repository Guardrails

### Community 4 - "redact.go"
Cohesion: 0.17
Nodes (22): redactAs(), Value(), Account(), hasDigitRun(), ID(), IDs(), luhn(), MaskEmails() (+14 more)

### Community 5 - "Error"
Cohesion: 0.05
Nodes (39): Error, New(), init(), TestError_IsAGoError(), TestWithUpstream_BoundsCause(), ExpectedHash(), ForParty(), FromError() (+31 more)

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

### Community 14 - "scrub.go"
Cohesion: 0.35
Nodes (14): keyInfo, Rule, RuleFunc, TestServiceRules(), DefaultRules(), isRef(), keyOf(), KeyRedact() (+6 more)

### Community 15 - "Refs"
Cohesion: 0.22
Nodes (6): Compute(), Kind, Refs, cacheKey, personKind(), TestPersonKind()

### Community 16 - "notify.go"
Cohesion: 0.27
Nodes (5): Channel, Message, Recipient, Sender, Status

### Community 18 - "test-next-version.sh"
Cohesion: 0.70
Nodes (4): c(), expect(), test-next-version.sh script, tg()

### Community 19 - "encode"
Cohesion: 0.21
Nodes (10): encode(), sortedJSON(), TestShortID(), Field, decodeInto(), Error, ShortID(), ShortIDs() (+2 more)

### Community 20 - "canonical.go"
Cohesion: 0.14
Nodes (15): lessUTF16(), writeCanonical(), writeNumber(), writeString(), numberMode, ErrUnknownField, child(), contains() (+7 more)

### Community 21 - "New"
Cohesion: 0.20
Nodes (10): New(), Parse(), TestAllAndParse(), TestInstitution(), TestKeyed(), TestKeyedFailsClosed(), TestPlainInDevelopment(), TestRefsCache() (+2 more)

### Community 23 - "encoding/json.RawMessage"
Cohesion: 0.31
Nodes (5): draftVector, fixture, Case, Detail, Event

### Community 24 - "approval_test.go"
Cohesion: 0.29
Nodes (11): Check(), Decisions(), cases(), generate(), TestCheckRefusals(), TestDecisions(), TestDecisionsGolden(), Decision (+3 more)

### Community 26 - "request.go"
Cohesion: 0.22
Nodes (15): RequestIDFrom(), TestActorMiddleware(), ActorMiddleware(), AccessLogOption, actorKey, loggerKey, TestActorSlot(), AccessLog() (+7 more)

### Community 27 - "telemetry_test.go"
Cohesion: 0.07
Nodes (18): Config, failing, Middleware(), Transport(), isURL(), metricOptions(), Propagator(), Setup() (+10 more)

### Community 28 - "testing.T"
Cohesion: 0.13
Nodes (22): TestCanonical_Properties(), TestLogFields_ShortFormsOnly(), New(), TestDevOIDC_ClientCredentials(), TestDevOIDC_SigningSessionClaims(), TestDevOIDC_TokenRequiresSubAndEmail(), TestDevOIDC_TokenVerifiesWithGoOIDC(), TestNewDefaultsToThePortal() (+14 more)

### Community 29 - "smtptest.go"
Cohesion: 0.17
Nodes (6): readPart(), WithMaxSize(), WithOnMessage(), Content, Message, Option

### Community 30 - "fieldMap"
Cohesion: 0.21
Nodes (11): fieldMap(), observed(), TestAccessLogServerError(), TestScrubCore(), TestScrubCoreWith(), refLogger(), TestRefScrubCore(), TestRefScrubCoreClear() (+3 more)

### Community 31 - "Principal"
Cohesion: 0.19
Nodes (11): Newf(), TestRequireScope(), ctxKey, fakeVerifier, devPrincipal(), looksLikeJWT(), Middleware(), RequireScope() (+3 more)

### Community 32 - "context.Context"
Cohesion: 0.33
Nodes (5): FromContext(), HasScope(), Client, Identity, PendingInvitation

### Community 33 - "MockServer"
Cohesion: 0.19
Nodes (3): path(), readData(), MockServer

### Community 34 - "smtptest_test.go"
Cohesion: 0.31
Nodes (8): NewMockServer(), start(), TestMockServer_FailWith(), TestMockServer_Protocol(), TestMockServer_RecordsWhatNetSMTPSends(), TestMockServer_SinglePart(), TestMockServer_TooBig(), TestNewMockServer_DevOnly()

### Community 35 - "mock.go"
Cohesion: 0.36
Nodes (3): actorSlot, cloneMessage(), MockSender

### Community 36 - "deliver.go"
Cohesion: 0.29
Nodes (9): Artifact, checkName(), Deliver(), artifacts(), TestDeliver_ManyIsDeterministicZip(), TestDeliver_One(), TestDeliver_RefusesBadNames(), Zip() (+1 more)

### Community 37 - "go_pkg_strings"
Cohesion: 0.06
Nodes (26): CodeForStatus(), TestCodeForStatus(), ValidCode(), requestIDKey, signing, Auth(), Handler(), Logging() (+18 more)

### Community 38 - "Issuer"
Cohesion: 0.09
Nodes (11): Write(), Writer, authCode, Issuer, codeStore, Issuer, httpHandlerFunc(), writeJSON() (+3 more)

### Community 39 - "agreementsig_test.go"
Cohesion: 0.08
Nodes (42): AgreementCanonical(), AgreementHash(), AgreementMessage(), checkID(), invalid(), agreementCases(), agreementJSON(), flipS() (+34 more)

### Community 40 - "NewTypst"
Cohesion: 0.24
Nodes (10): CheckFonts(), NewTypst(), TestCheckFonts(), TestStr(), Content(), Lit(), Str(), TestTypst_ByteIdenticalAndEscaped() (+2 more)

### Community 41 - "notify_test.go"
Cohesion: 0.27
Nodes (11): NewMockSender(), IsDevelopment(), TestIsDevelopment(), TestMockSender_RecordsOncePerIdempotencyKey(), TestMockSender_SentIsACopy(), TestMockSender_ValidatesAndFails(), TestNewMockSender_RefusedOutsideDevelopment(), TestValidate_AcceptsAWellFormedMessage() (+3 more)

### Community 42 - "legalorders/legalorders.go"
Cohesion: 0.18
Nodes (10): CaseRequest, EventRequest, ValidCase(), validNotifyAfter(), Lifted, LiftRequest, Recorded, RecordedEvent (+2 more)

### Community 43 - "http_test.go"
Cohesion: 0.36
Nodes (8): NewHTTPSender(), staticToken(), TestHTTPSender_RefusesAnInvalidMessageWithoutCalling(), TestHTTPSender_SendsAndReadsTheReceipt(), TestHTTPSender_Status(), TestHTTPSender_StatusAnswersToTheRelay(), TestHTTPSender_TransportFailuresAreRetryable(), TestNewHTTPSender_Refusals()

### Community 44 - ".Verify"
Cohesion: 0.22
Nodes (10): signer(), TestOIDCVerifierReadsClaims(), TestOIDCVerifierRefuses(), token(), NewOIDCVerifierWithKeys(), Sign(), TestSignAndVerify_RoundTrip(), TestVerify_MalformedHexFails() (+2 more)

### Community 45 - "TestMiddleware"
Cohesion: 0.20
Nodes (11): TestClientPolicyFilter(), TestMiddleware(), TestParseClientPolicy(), ClientPolicy, WithClientPolicy(), WithDevBypass(), WithLogger(), Option (+3 more)

### Community 46 - "time.Time"
Cohesion: 0.21
Nodes (8): LegalOrder, LegalOrderEvent, LegalOrderTarget, LiftLegalOrderResponse, ListLegalOrdersRequest, OrderContact, RecordedLegalOrder, Renderer

### Community 47 - "NewRefScrubCore"
Cohesion: 0.33
Nodes (9): benchCore(), benchLine(), BenchmarkRawZap(), BenchmarkShapingCore(), benchRefs(), TestShapingAllocs(), NewRefScrubCore(), NewScrubCore() (+1 more)

### Community 48 - "field_test.go"
Cohesion: 0.39
Nodes (6): loadFields(), TestField_Errors(), TestField_KeepsOwnDescriptionAndRefKeywords(), TestField_RefFailures(), TestField_Resolves(), TestField_Root()

### Community 49 - "HTTPSender"
Cohesion: 0.23
Nodes (6): refusal(), WithHTTPClient(), HTTPOption, HTTPSender, Receipt, TokenFunc

### Community 50 - "legalorders/legalorders_test.go"
Cohesion: 0.50
Nodes (4): TestEventAndLiftValidate(), TestRecordRequestValidate(), TestWireNames(), valid()

### Community 52 - "audit.go"
Cohesion: 0.27
Nodes (9): ActiveHolds, AuditIdentity, AuditRef, AuditRefsOf, LegalHold, Resolution, ResolveAuditRefs, ResolvedRef (+1 more)

### Community 53 - "fields.go"
Cohesion: 0.27
Nodes (8): accountNumber, Account(), Party(), Subject(), User(), personID, TraceFields(), TestTypedFields()

### Community 54 - "authorize_test.go"
Cohesion: 0.20
Nodes (12): authorizeQuery(), browserIssuer(), challengeOf(), exchange(), signIn(), TestDevOIDC_AuthorizeRefusals(), TestDevOIDC_BrowserPortalSession(), TestDevOIDC_BrowserSigningSession() (+4 more)

### Community 55 - "build"
Cohesion: 0.32
Nodes (7): Config, build(), New(), TestNew(), TestNewIsNotSampled(), zapConfig(), TestNewWithRefs()

### Community 56 - "LoadDirectory"
Cohesion: 0.25
Nodes (7): Registry, LoadDirectory(), TestLoadDirectory_CompilesSchemas(), TestLoadDirectory_MissingDirectory(), TestValidate_InvalidPayloadReportsFailures(), TestValidate_UnknownType(), TestValidate_ValidPayload()

### Community 57 - "identityclient_test.go"
Cohesion: 0.36
Nodes (9): New(), TestClient_CreatePendingInvitation(), TestClient_CreatePendingInvitation_ServerError(), TestClient_GetByEmail_ServerError(), TestClient_GetByOktaSub_Found(), TestClient_GetByToken_NotFound(), TestClient_ManagesIdentity(), TestClient_ManagesIdentity_ServerError() (+1 more)

### Community 58 - "OIDCVerifier"
Cohesion: 0.25
Nodes (5): claims, NewOIDCVerifier(), uniqueScopes(), OIDCVerifier, TestStartAndMintVerifyWithAuthn()

### Community 60 - "TestAccessLog"
Cohesion: 0.33
Nodes (7): newRequestID(), RequestID(), spanCtx(), TestAccessLog(), TestFromContext(), FromContext(), SetInstitution()

### Community 61 - "erasure.go"
Cohesion: 0.29
Nodes (6): ErasureRequest, ErasureRequestID, ListErasureRequests, RefuseErasureRequest, SubmitErasureRequest, VerifyErasureRequest

### Community 62 - "signin.go"
Cohesion: 0.40
Nodes (5): SigninActionRequest, SigninActionResult, SigninFactor, SigninStatus, SigninStatusRequest

### Community 63 - "keys.go"
Cohesion: 0.60
Nodes (4): shapeNumber(), ForKey(), Kind, normalize()

## Knowledge Gaps
- **80 isolated node(s):** `signing`, `requestIDKey`, `Error`, `ctxKey`, `github.com/vdatacloud/daml-escrow-commons` (+75 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 187 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **9 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Issuer` connect `Issuer` to `canonical.go`, `testing.T`, `go_pkg_strings`, `agreementsig_test.go`?**
  _High betweenness centrality (0.035) - this node is a cross-community bridge._
- **What connects `signing`, `requestIDKey`, `Error` to the rest of the system?**
  _80 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `GitHub & CI Workflow Skill` be split into smaller, more focused modules?**
  _Cohesion score 0.1111111111111111 - nodes in this community are weakly interconnected._
- **Why does `MockServer` connect `MockServer` to `telemetry_test.go`, `smtptest_test.go`, `mock.go`, `smtptest.go`?**
  _High betweenness centrality (0.020) - this node is a cross-community bridge._
- **Should `Error` be split into smaller, more focused modules?**
  _Cohesion score 0.050595238095238096 - nodes in this community are weakly interconnected._
- **Why does `Refs` connect `Refs` to `go_pkg_strings`, `scrub.go`, `NewRefScrubCore`, `New`, `build`, `request.go`, `fieldMap`?**
  _High betweenness centrality (0.019) - this node is a cross-community bridge._
- **Should `metering.go` be split into smaller, more focused modules?**
  _Cohesion score 0.12121212121212122 - nodes in this community are weakly interconnected._