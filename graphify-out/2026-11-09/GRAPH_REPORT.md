# Graph Report - daml-escrow-commons  (2026-11-09)

## Corpus Check
- 140 files · ~63,354 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 9 file(s) not represented in the graph (top: .ttf 5, (none) 4)

## Summary
- 980 nodes · 2473 edges · 58 communities (50 shown, 8 thin omitted)
- Extraction: 90% EXTRACTED · 10% INFERRED · 0% AMBIGUOUS · INFERRED: 243 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `dde9dda5`
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
- auditref.go
- Refs
- notify.go
- next-version.sh
- test-next-version.sh
- Error
- canonical.go
- auditref_test.go
- canton_test.go
- approval_test.go
- context.Context
- telemetry_test.go
- testing.T
- smtptest.go
- person_test.go
- go_pkg_reflect
- pathparam_test.go
- MockServer
- smtptest_test.go
- mock.go
- deliver.go
- go_pkg_strings
- Issuer
- agreementsig_test.go
- typst_test.go
- notify_test.go
- legalorders/legalorders.go
- http_test.go
- authn_test.go
- authority.go
- time.Time
- bench_test.go
- field_test.go
- schema/schema.go
- legalorders/legalorders_test.go
- Str
- fields.go
- authorize_test.go
- logging.go
- LoadDirectory
- identityclient_test.go
- metering_test.go

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

## Communities (58 total, 8 thin omitted)

### Community 0 - "GitHub & CI Workflow Skill"
Cohesion: 0.11
Nodes (17): 10. Useful Reference Commands, 11. References, 1. Pre-Commit Local Verification (MANDATORY — do this before every commit), 2. Branching Rules, 3. Commit Standards, 4. Staging & Pushing Changes, 5. Pull Request Creation, 6. CI Pipeline Overview (+9 more)

### Community 1 - "apierror_test.go"
Cohesion: 0.16
Nodes (20): CodeForStatus(), New(), Parse(), golden(), init(), TestCanonical_Golden(), TestCodeForStatus(), TestDetailsInto() (+12 more)

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
Cohesion: 0.10
Nodes (19): Classify(), FingerprintOf(), Fingerprint, Hash, PartyID, MustParsePartyID(), ParseFingerprint(), ParseHash() (+11 more)

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

### Community 14 - "auditref.go"
Cohesion: 0.21
Nodes (4): cacheKey, Config, Ref(), TestRef()

### Community 15 - "Refs"
Cohesion: 0.26
Nodes (3): Compute(), Kind, Refs

### Community 16 - "notify.go"
Cohesion: 0.16
Nodes (9): Channel, refusal(), HTTPSender, Message, Receipt, Recipient, Sender, Status (+1 more)

### Community 18 - "test-next-version.sh"
Cohesion: 0.70
Nodes (4): c(), expect(), test-next-version.sh script, tg()

### Community 19 - "Error"
Cohesion: 0.09
Nodes (17): draftVector, encode(), Error, sortedJSON(), TestShortID(), Detail, Field, decodeInto() (+9 more)

### Community 20 - "canonical.go"
Cohesion: 0.14
Nodes (15): lessUTF16(), writeCanonical(), writeNumber(), writeString(), numberMode, ErrUnknownField, child(), contains() (+7 more)

### Community 21 - "auditref_test.go"
Cohesion: 0.36
Nodes (9): New(), Parse(), TestAllAndParse(), TestInstitution(), TestKeyed(), TestKeyedFailsClosed(), TestPlainInDevelopment(), TestRefsCache() (+1 more)

### Community 23 - "canton_test.go"
Cohesion: 0.18
Nodes (15): ExpectedHash(), ForParty(), FromError(), FromLedgerError(), GRPCCodeName(), ParseLedgerError(), Signing(), TestFromError() (+7 more)

### Community 24 - "approval_test.go"
Cohesion: 0.29
Nodes (11): Check(), Decisions(), cases(), generate(), TestCheckRefusals(), TestDecisions(), TestDecisionsGolden(), Decision (+3 more)

### Community 26 - "context.Context"
Cohesion: 0.07
Nodes (47): Newf(), newRequestID(), RequestID(), RequestIDFrom(), TestActorMiddleware(), TestMiddleware(), TestRequireScope(), claims (+39 more)

### Community 27 - "telemetry_test.go"
Cohesion: 0.08
Nodes (18): Config, failing, Middleware(), Transport(), isURL(), metricOptions(), Propagator(), Setup() (+10 more)

### Community 28 - "testing.T"
Cohesion: 0.14
Nodes (17): TestCanonical_Properties(), TestLogFields_ShortFormsOnly(), New(), TestDevOIDC_ClientCredentials(), TestDevOIDC_SigningSessionClaims(), TestDevOIDC_TokenRequiresSubAndEmail(), TestDevOIDC_TokenVerifiesWithGoOIDC(), Start() (+9 more)

### Community 29 - "smtptest.go"
Cohesion: 0.17
Nodes (6): readPart(), WithMaxSize(), WithOnMessage(), Content, Message, Option

### Community 30 - "person_test.go"
Cohesion: 0.18
Nodes (11): fieldMap(), observed(), TestAccessLogServerError(), TestScrubCore(), TestScrubCoreWith(), refLogger(), TestRefScrubCore(), TestRefScrubCoreClear() (+3 more)

### Community 31 - "go_pkg_reflect"
Cohesion: 0.40
Nodes (3): main(), splitOrigins(), TestSplitOrigins()

### Community 32 - "pathparam_test.go"
Cohesion: 0.43
Nodes (5): Get(), containsByte(), jsEncodeURIComponent(), TestGetKeepsInvalidEscaping(), TestGetUnescapes()

### Community 33 - "MockServer"
Cohesion: 0.14
Nodes (5): Issuer, Token, path(), readData(), MockServer

### Community 34 - "smtptest_test.go"
Cohesion: 0.31
Nodes (8): NewMockServer(), start(), TestMockServer_FailWith(), TestMockServer_Protocol(), TestMockServer_RecordsWhatNetSMTPSends(), TestMockServer_SinglePart(), TestMockServer_TooBig(), TestNewMockServer_DevOnly()

### Community 35 - "mock.go"
Cohesion: 0.36
Nodes (3): actorSlot, cloneMessage(), MockSender

### Community 36 - "deliver.go"
Cohesion: 0.29
Nodes (9): Artifact, checkName(), Deliver(), artifacts(), TestDeliver_ManyIsDeterministicZip(), TestDeliver_One(), TestDeliver_RefusesBadNames(), Zip() (+1 more)

### Community 38 - "Issuer"
Cohesion: 0.10
Nodes (11): Write(), Writer, authCode, Issuer, codeStore, Issuer, httpHandlerFunc(), writeJSON() (+3 more)

### Community 39 - "agreementsig_test.go"
Cohesion: 0.07
Nodes (47): AgreementCanonical(), AgreementHash(), AgreementMessage(), checkID(), invalid(), agreementCases(), agreementJSON(), flipS() (+39 more)

### Community 40 - "typst_test.go"
Cohesion: 0.38
Nodes (6): CheckFonts(), NewTypst(), TestCheckFonts(), TestTypst_ByteIdenticalAndEscaped(), TestTypst_GoldenHashAcrossPlatforms(), TestTypst_RefusesOtherVersion()

### Community 41 - "notify_test.go"
Cohesion: 0.27
Nodes (11): NewMockSender(), IsDevelopment(), TestIsDevelopment(), TestMockSender_RecordsOncePerIdempotencyKey(), TestMockSender_SentIsACopy(), TestMockSender_ValidatesAndFails(), TestNewMockSender_RefusedOutsideDevelopment(), TestValidate_AcceptsAWellFormedMessage() (+3 more)

### Community 42 - "legalorders/legalorders.go"
Cohesion: 0.18
Nodes (10): CaseRequest, EventRequest, ValidCase(), validNotifyAfter(), Lifted, LiftRequest, Recorded, RecordedEvent (+2 more)

### Community 43 - "http_test.go"
Cohesion: 0.24
Nodes (10): NewHTTPSender(), staticToken(), TestHTTPSender_RefusesAnInvalidMessageWithoutCalling(), TestHTTPSender_SendsAndReadsTheReceipt(), TestHTTPSender_Status(), TestHTTPSender_StatusAnswersToTheRelay(), TestHTTPSender_TransportFailuresAreRetryable(), TestNewHTTPSender_Refusals() (+2 more)

### Community 44 - "authn_test.go"
Cohesion: 0.22
Nodes (6): signer(), TestOIDCVerifierReadsClaims(), TestOIDCVerifierRefuses(), token(), personKind(), TestPersonKind()

### Community 45 - "authority.go"
Cohesion: 0.16
Nodes (14): TestClientPolicyFilter(), TestParseClientPolicy(), ClientPolicy, options, NewClientPolicy(), ParseClientPolicy(), Auth(), Handler() (+6 more)

### Community 46 - "time.Time"
Cohesion: 0.05
Nodes (33): Response, Status, ActiveHolds, AuditIdentity, AuditRef, AuditRefsOf, ErasureRequest, ErasureRequestID (+25 more)

### Community 47 - "bench_test.go"
Cohesion: 0.54
Nodes (6): benchCore(), benchLine(), BenchmarkRawZap(), BenchmarkShapingCore(), benchRefs(), TestShapingAllocs()

### Community 48 - "field_test.go"
Cohesion: 0.39
Nodes (6): loadFields(), TestField_Errors(), TestField_KeepsOwnDescriptionAndRefKeywords(), TestField_RefFailures(), TestField_Resolves(), TestField_Root()

### Community 50 - "legalorders/legalorders_test.go"
Cohesion: 0.50
Nodes (4): TestEventAndLiftValidate(), TestRecordRequestValidate(), TestWireNames(), valid()

### Community 51 - "Str"
Cohesion: 0.70
Nodes (4): TestStr(), Content(), Lit(), Str()

### Community 53 - "fields.go"
Cohesion: 0.27
Nodes (7): accountNumber, Account(), Party(), Subject(), User(), personID, TestTypedFields()

### Community 54 - "authorize_test.go"
Cohesion: 0.38
Nodes (8): authorizeQuery(), browserIssuer(), challengeOf(), exchange(), signIn(), TestDevOIDC_AuthorizeRefusals(), TestDevOIDC_BrowserPortalSession(), TestDevOIDC_BrowserSigningSession()

### Community 55 - "logging.go"
Cohesion: 0.36
Nodes (7): Config, build(), New(), TestNew(), TestNewIsNotSampled(), zapConfig(), TestNewWithRefs()

### Community 56 - "LoadDirectory"
Cohesion: 0.25
Nodes (7): Registry, LoadDirectory(), TestLoadDirectory_CompilesSchemas(), TestLoadDirectory_MissingDirectory(), TestValidate_InvalidPayloadReportsFailures(), TestValidate_UnknownType(), TestValidate_ValidPayload()

### Community 57 - "identityclient_test.go"
Cohesion: 0.36
Nodes (9): New(), TestClient_CreatePendingInvitation(), TestClient_CreatePendingInvitation_ServerError(), TestClient_GetByEmail_ServerError(), TestClient_GetByOktaSub_Found(), TestClient_GetByToken_NotFound(), TestClient_ManagesIdentity(), TestClient_ManagesIdentity_ServerError() (+1 more)

### Community 59 - "metering_test.go"
Cohesion: 0.28
Nodes (8): TestLedgerCommandEvent_Timing(), TestLedgerCommandEvent_Traffic(), TestLedgerCommandEvent_Validate(), TestNetworkFeeEvent_SpotPrice(), TestNetworkFeeEvent_Validate(), TestSettlementEvent_NetworkFee(), TestSettlementEvent_Validate(), withFee()

## Knowledge Gaps
- **80 isolated node(s):** `signing`, `requestIDKey`, `Error`, `ctxKey`, `github.com/vdatacloud/daml-escrow-commons` (+75 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 187 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **8 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Issuer` connect `Issuer` to `go_pkg_strings`, `agreementsig_test.go`, `authn_test.go`, `auditref.go`, `canonical.go`, `testing.T`?**
  _High betweenness centrality (0.035) - this node is a cross-community bridge._
- **What connects `signing`, `requestIDKey`, `Error` to the rest of the system?**
  _80 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `GitHub & CI Workflow Skill` be split into smaller, more focused modules?**
  _Cohesion score 0.1111111111111111 - nodes in this community are weakly interconnected._
- **Why does `MockServer` connect `MockServer` to `smtptest_test.go`, `mock.go`, `smtptest.go`?**
  _High betweenness centrality (0.020) - this node is a cross-community bridge._
- **Should `scrub.go` be split into smaller, more focused modules?**
  _Cohesion score 0.08771929824561403 - nodes in this community are weakly interconnected._
- **Why does `Refs` connect `Refs` to `scrub.go`, `auditref.go`, `bench_test.go`, `auditref_test.go`, `logging.go`, `context.Context`, `person_test.go`?**
  _High betweenness centrality (0.019) - this node is a cross-community bridge._
- **Should `cantonid.go` be split into smaller, more focused modules?**
  _Cohesion score 0.10099573257467995 - nodes in this community are weakly interconnected._