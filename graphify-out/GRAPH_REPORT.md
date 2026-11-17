# Graph Report - daml-escrow-commons  (2026-11-17)

## Corpus Check
- 143 files · ~65,114 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 9 file(s) not represented in the graph (top: .ttf 5, (none) 4)

## Summary
- 1003 nodes · 2532 edges · 65 communities (56 shown, 9 thin omitted)
- Extraction: 90% EXTRACTED · 10% INFERRED · 0% AMBIGUOUS · INFERRED: 250 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `603428d1`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- GitHub & CI Workflow Skill
- pathparam_test.go
- widget.json
- Repository Guardrails
- redact.go
- context.Context
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
- TestCanonical_Golden
- next-version.sh
- test-next-version.sh
- Error
- canonical.go
- New
- New
- approval_test.go
- fieldMap
- telemetry_test.go
- testing.T
- smtptest.go
- validate.go
- PartyID
- time.Time
- MockServer
- smtptest_test.go
- mock.go
- deliver.go
- go_pkg_strings
- Issuer
- agreementsig_test.go
- TestSigningDetails
- notify_test.go
- legalorders/legalorders.go
- http_test.go
- encoding/json.RawMessage
- authority.go
- Parse
- NewRefScrubCore
- field_test.go
- HTTPSender
- erasure.go
- go.uber.org/zap/zapcore.Field
- audit.go
- fields.go
- authorize_test.go
- build
- LoadDirectory
- New
- signin.go
- NewTypst
- RequireNonEmpty
- legalorders/legalorders_test.go
- shortHex
- gastank.go
- Fingerprint

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
- `SigningDetails` --references--> `PartyID`  [EXTRACTED]
  apierror/canton/canton.go → cantonid/cantonid.go
- `ForParty()` --references--> `PartyID`  [EXTRACTED]
  apierror/canton/canton.go → cantonid/cantonid.go
- `TestSigningDetails()` --calls--> `MustParsePartyID()`  [EXTRACTED]
  apierror/canton/canton_test.go → cantonid/cantonid.go

## Import Cycles
- None detected.

## Communities (65 total, 9 thin omitted)

### Community 0 - "GitHub & CI Workflow Skill"
Cohesion: 0.11
Nodes (17): 10. Useful Reference Commands, 11. References, 1. Pre-Commit Local Verification (MANDATORY — do this before every commit), 2. Branching Rules, 3. Commit Standards, 4. Staging & Pushing Changes, 5. Pull Request Creation, 6. CI Pipeline Overview (+9 more)

### Community 1 - "pathparam_test.go"
Cohesion: 0.16
Nodes (11): Check, PathFor(), TestFinish(), TestPathFor(), Setup, Status, Get(), containsByte() (+3 more)

### Community 2 - "widget.json"
Cohesion: 0.17
Nodes (11): minLength, type, properties, name, quantity, minimum, type, required (+3 more)

### Community 3 - "Repository Guardrails"
Cohesion: 0.22
Nodes (8): Branch Protection Strategy, Branching Strategy, CI Requirements, Code Review Requirements, Commit Standard, Pre-Commit Verification, Pull Request Rules, Repository Guardrails

### Community 4 - "redact.go"
Cohesion: 0.13
Nodes (26): redactAs(), shapeNumber(), ForKey(), Kind, normalize(), Value(), Account(), hasDigitRun() (+18 more)

### Community 5 - "context.Context"
Cohesion: 0.07
Nodes (45): Newf(), newRequestID(), RequestID(), RequestIDFrom(), TestActorMiddleware(), TestMiddleware(), TestRequireScope(), ctxKey (+37 more)

### Community 6 - "Releasing"
Cohesion: 0.18
Nodes (10): 1. Prerequisites (one-time, per machine that will `go get` this module), 2. Decide the version bump, 3. Tag and push, 4. Publish the GitHub release, 5. Update consumers, Automated (default), Future automation, Major releases (+2 more)

### Community 7 - "metering.go"
Cohesion: 0.23
Nodes (10): ChargeBearer, requireBaseUnits(), requireDecimal(), NetworkFeeEvent, NetworkFeePayer, Rail, RetentionEvent, RetentionTier (+2 more)

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
Cohesion: 0.34
Nodes (16): keyInfo, Rule, RuleFunc, TestServiceRules(), DefaultRules(), isRef(), keyOf(), KeyRedact() (+8 more)

### Community 15 - "Refs"
Cohesion: 0.22
Nodes (6): Compute(), Kind, Refs, cacheKey, personKind(), TestPersonKind()

### Community 16 - "TestCanonical_Golden"
Cohesion: 0.29
Nodes (7): golden(), TestCanonical_Golden(), TestLogFields_Golden(), TestSchemaIsACopy(), TestWriter(), validate(), Schema()

### Community 18 - "test-next-version.sh"
Cohesion: 0.70
Nodes (4): c(), expect(), test-next-version.sh script, tg()

### Community 19 - "Error"
Cohesion: 0.17
Nodes (6): encode(), Error, sortedJSON(), Detail, Upstream, wire

### Community 20 - "canonical.go"
Cohesion: 0.14
Nodes (15): lessUTF16(), writeCanonical(), writeNumber(), writeString(), numberMode, ErrUnknownField, child(), contains() (+7 more)

### Community 21 - "New"
Cohesion: 0.20
Nodes (10): New(), Parse(), TestAllAndParse(), TestInstitution(), TestKeyed(), TestKeyedFailsClosed(), TestPlainInDevelopment(), TestRefsCache() (+2 more)

### Community 23 - "New"
Cohesion: 0.18
Nodes (13): New(), init(), TestError_IsAGoError(), TestWithUpstream_BoundsCause(), FromError(), FromLedgerError(), GRPCCodeName(), ParseLedgerError() (+5 more)

### Community 24 - "approval_test.go"
Cohesion: 0.29
Nodes (11): Check(), Decisions(), cases(), generate(), TestCheckRefusals(), TestDecisions(), TestDecisionsGolden(), Decision (+3 more)

### Community 26 - "fieldMap"
Cohesion: 0.20
Nodes (12): fieldMap(), observed(), TestAccessLogServerError(), TestFromContext(), TestScrubCore(), TestScrubCoreWith(), refLogger(), TestRefScrubCore() (+4 more)

### Community 27 - "telemetry_test.go"
Cohesion: 0.07
Nodes (18): Config, failing, Middleware(), Transport(), isURL(), metricOptions(), Propagator(), Setup() (+10 more)

### Community 28 - "testing.T"
Cohesion: 0.15
Nodes (18): TestCanonical_Properties(), TestLogFields_ShortFormsOnly(), TestWireNames(), TestErasureRequestWireNames(), TestLegalOrderWireNames(), TestSigninWireNames(), RetentionPeriod(), TestLedgerCommandEvent_Timing() (+10 more)

### Community 29 - "smtptest.go"
Cohesion: 0.17
Nodes (6): readPart(), WithMaxSize(), WithOnMessage(), Content, Message, Option

### Community 30 - "validate.go"
Cohesion: 0.16
Nodes (9): Errors, RequirePositive(), RequireValidEmail(), TestErrors_AggregatesAndReports(), TestErrors_ErrIfAny_NilWhenEmpty(), TestRequireNonEmpty(), TestRequireOneOf(), TestRequirePositive() (+1 more)

### Community 31 - "PartyID"
Cohesion: 0.22
Nodes (9): Classify(), PartyID, MustParsePartyID(), ParsePartyID(), Short(), shortRunes(), TestParse_Golden(), TestPartyID() (+1 more)

### Community 32 - "time.Time"
Cohesion: 0.21
Nodes (8): LegalOrder, LegalOrderEvent, LegalOrderTarget, LiftLegalOrderResponse, ListLegalOrdersRequest, OrderContact, RecordedLegalOrder, Renderer

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
Cohesion: 0.07
Nodes (17): TestShortID(), Field, decodeInto(), requestIDKey, Error, ShortID(), ShortIDs(), shortValue() (+9 more)

### Community 38 - "Issuer"
Cohesion: 0.06
Nodes (20): Write(), Writer, main(), splitOrigins(), TestSplitOrigins(), authCode, Issuer, codeStore (+12 more)

### Community 39 - "agreementsig_test.go"
Cohesion: 0.07
Nodes (48): AgreementCanonical(), AgreementHash(), AgreementMessage(), checkID(), invalid(), agreementCases(), agreementJSON(), flipS() (+40 more)

### Community 40 - "TestSigningDetails"
Cohesion: 0.31
Nodes (7): ExpectedHash(), ForParty(), Signing(), TestSigningDetails(), SigningDetails, Hash, ParseHash()

### Community 41 - "notify_test.go"
Cohesion: 0.27
Nodes (11): NewMockSender(), IsDevelopment(), TestIsDevelopment(), TestMockSender_RecordsOncePerIdempotencyKey(), TestMockSender_SentIsACopy(), TestMockSender_ValidatesAndFails(), TestNewMockSender_RefusedOutsideDevelopment(), TestValidate_AcceptsAWellFormedMessage() (+3 more)

### Community 42 - "legalorders/legalorders.go"
Cohesion: 0.18
Nodes (10): CaseRequest, EventRequest, ValidCase(), validNotifyAfter(), Lifted, LiftRequest, Recorded, RecordedEvent (+2 more)

### Community 43 - "http_test.go"
Cohesion: 0.22
Nodes (11): NewHTTPSender(), staticToken(), TestHTTPSender_RefusesAnInvalidMessageWithoutCalling(), TestHTTPSender_SendsAndReadsTheReceipt(), TestHTTPSender_Status(), TestHTTPSender_StatusAnswersToTheRelay(), TestHTTPSender_TransportFailuresAreRetryable(), TestNewHTTPSender_Refusals() (+3 more)

### Community 44 - "encoding/json.RawMessage"
Cohesion: 0.31
Nodes (5): draftVector, fixture, Case, Detail, Event

### Community 45 - "authority.go"
Cohesion: 0.16
Nodes (14): TestClientPolicyFilter(), TestParseClientPolicy(), ClientPolicy, options, NewClientPolicy(), ParseClientPolicy(), Auth(), Handler() (+6 more)

### Community 46 - "Parse"
Cohesion: 0.25
Nodes (8): CodeForStatus(), Parse(), TestCodeForStatus(), TestDetailsInto(), TestFromResponse(), TestParse(), ValidCode(), FromResponse()

### Community 47 - "NewRefScrubCore"
Cohesion: 0.43
Nodes (7): benchCore(), benchLine(), BenchmarkRawZap(), BenchmarkShapingCore(), benchRefs(), TestShapingAllocs(), NewRefScrubCore()

### Community 48 - "field_test.go"
Cohesion: 0.39
Nodes (6): loadFields(), TestField_Errors(), TestField_KeepsOwnDescriptionAndRefKeywords(), TestField_RefFailures(), TestField_Resolves(), TestField_Root()

### Community 49 - "HTTPSender"
Cohesion: 0.25
Nodes (5): refusal(), HTTPSender, Receipt, Sender, Status

### Community 50 - "erasure.go"
Cohesion: 0.29
Nodes (6): ErasureRequest, ErasureRequestID, ListErasureRequests, RefuseErasureRequest, SubmitErasureRequest, VerifyErasureRequest

### Community 52 - "audit.go"
Cohesion: 0.27
Nodes (9): ActiveHolds, AuditIdentity, AuditRef, AuditRefsOf, LegalHold, Resolution, ResolveAuditRefs, ResolvedRef (+1 more)

### Community 53 - "fields.go"
Cohesion: 0.31
Nodes (7): accountNumber, Account(), Party(), Subject(), User(), personID, TestTypedFields()

### Community 54 - "authorize_test.go"
Cohesion: 0.09
Nodes (22): signer(), TestOIDCVerifierReadsClaims(), TestOIDCVerifierRefuses(), token(), claims, NewOIDCVerifier(), NewOIDCVerifierWithKeys(), uniqueScopes() (+14 more)

### Community 55 - "build"
Cohesion: 0.32
Nodes (7): Config, build(), New(), TestNew(), TestNewIsNotSampled(), zapConfig(), TestNewWithRefs()

### Community 56 - "LoadDirectory"
Cohesion: 0.25
Nodes (7): Registry, LoadDirectory(), TestLoadDirectory_CompilesSchemas(), TestLoadDirectory_MissingDirectory(), TestValidate_InvalidPayloadReportsFailures(), TestValidate_UnknownType(), TestValidate_ValidPayload()

### Community 57 - "New"
Cohesion: 0.24
Nodes (12): New(), TestClient_CreatePendingInvitation(), TestClient_CreatePendingInvitation_ServerError(), TestClient_GetByEmail_ServerError(), TestClient_GetByOktaSub_Found(), TestClient_GetByToken_NotFound(), TestClient_ManagesIdentity(), TestClient_ManagesIdentity_ServerError() (+4 more)

### Community 58 - "signin.go"
Cohesion: 0.40
Nodes (5): SigninActionRequest, SigninActionResult, SigninFactor, SigninStatus, SigninStatusRequest

### Community 59 - "NewTypst"
Cohesion: 0.13
Nodes (12): Fake, CheckFonts(), NewTypst(), TestCheckFonts(), TestStr(), Content(), Lit(), Str() (+4 more)

### Community 60 - "RequireNonEmpty"
Cohesion: 0.60
Nodes (3): CommandOutcome, LedgerCommandEvent, RequireNonEmpty()

### Community 61 - "legalorders/legalorders_test.go"
Cohesion: 0.50
Nodes (4): TestEventAndLiftValidate(), TestRecordRequestValidate(), TestWireNames(), valid()

### Community 66 - "Fingerprint"
Cohesion: 0.22
Nodes (6): FingerprintOf(), Fingerprint, ParseFingerprint(), TestFingerprintOf_ControlsExternalParty(), TestJSON(), TestShortForms()

## Knowledge Gaps
- **80 isolated node(s):** `signing`, `requestIDKey`, `Error`, `ctxKey`, `github.com/vdatacloud/daml-escrow-commons` (+75 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 188 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **9 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Issuer` connect `Issuer` to `canonical.go`, `go_pkg_strings`, `agreementsig_test.go`?**
  _High betweenness centrality (0.036) - this node is a cross-community bridge._
- **What connects `signing`, `requestIDKey`, `Error` to the rest of the system?**
  _80 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `GitHub & CI Workflow Skill` be split into smaller, more focused modules?**
  _Cohesion score 0.1111111111111111 - nodes in this community are weakly interconnected._
- **Why does `Refs` connect `Refs` to `context.Context`, `go_pkg_strings`, `scrub.go`, `NewRefScrubCore`, `New`, `build`, `fieldMap`?**
  _High betweenness centrality (0.028) - this node is a cross-community bridge._
- **Should `redact.go` be split into smaller, more focused modules?**
  _Cohesion score 0.13333333333333333 - nodes in this community are weakly interconnected._
- **Why does `MockServer` connect `MockServer` to `telemetry_test.go`, `smtptest_test.go`, `mock.go`, `smtptest.go`?**
  _High betweenness centrality (0.019) - this node is a cross-community bridge._
- **Should `context.Context` be split into smaller, more focused modules?**
  _Cohesion score 0.07344632768361582 - nodes in this community are weakly interconnected._