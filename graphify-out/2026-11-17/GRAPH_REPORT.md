# Graph Report - daml-escrow-commons  (2026-11-16)

## Corpus Check
- 141 files · ~64,591 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 9 file(s) not represented in the graph (top: .ttf 5, (none) 4)

## Summary
- 994 nodes · 2516 edges · 58 communities (49 shown, 9 thin omitted)
- Extraction: 90% EXTRACTED · 10% INFERRED · 0% AMBIGUOUS · INFERRED: 249 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `225f9072`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- GitHub & CI Workflow Skill
- canonical.go
- widget.json
- Repository Guardrails
- redact.go
- context.Context
- Releasing
- time.Time
- CLAUDE.md
- /commons-contribution
- daml-escrow-commons
- .claude/CLAUDE.md
- install-git-hooks.sh
- github.com/vdatacloud/daml-escrow-commons
- scrub.go
- Refs
- Parse
- next-version.sh
- test-next-version.sh
- Error
- field.go
- auditref_test.go
- Client
- approval_test.go
- fieldMap
- telemetry_test.go
- testing.T
- smtptest.go
- Principal
- Fingerprint
- net/http.Handler
- MockServer
- smtptest_test.go
- mock.go
- deliver.go
- go_pkg_strings
- Issuer
- agreementsig_test.go
- Token
- notify_test.go
- legalorders/legalorders.go
- http_test.go
- .Verify
- authority.go
- bench_test.go
- field_test.go
- HTTPSender
- go.uber.org/zap/zapcore.Field
- audit.go
- fields.go
- authorize_test.go
- build
- LoadDirectory
- New
- NewTypst
- signature.go

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
- `TestAccessLog()` --calls--> `RequestID()`  [EXTRACTED]
  logging/logging_test.go → apierror/http.go
- `AccessLog()` --calls--> `RequestIDFrom()`  [EXTRACTED]
  logging/request.go → apierror/http.go
- `keyInfo` --references--> `Kind`  [EXTRACTED]
  logging/scrub.go → auditref/auditref.go
- `personID` --references--> `Kind`  [EXTRACTED]
  logging/fields.go → auditref/auditref.go

## Import Cycles
- None detected.

## Communities (58 total, 9 thin omitted)

### Community 0 - "GitHub & CI Workflow Skill"
Cohesion: 0.11
Nodes (17): 10. Useful Reference Commands, 11. References, 1. Pre-Commit Local Verification (MANDATORY — do this before every commit), 2. Branching Rules, 3. Commit Standards, 4. Staging & Pushing Changes, 5. Pull Request Creation, 6. CI Pipeline Overview (+9 more)

### Community 1 - "canonical.go"
Cohesion: 0.14
Nodes (19): AgreementCanonical(), AgreementMessage(), checkID(), invalid(), TestDraftVectors(), AgreementVersion, Amends, decodeJSON() (+11 more)

### Community 2 - "widget.json"
Cohesion: 0.17
Nodes (11): minLength, type, properties, name, quantity, minimum, type, required (+3 more)

### Community 3 - "Repository Guardrails"
Cohesion: 0.22
Nodes (8): Branch Protection Strategy, Branching Strategy, CI Requirements, Code Review Requirements, Commit Standard, Pre-Commit Verification, Pull Request Rules, Repository Guardrails

### Community 4 - "redact.go"
Cohesion: 0.14
Nodes (26): redactAs(), shapeNumber(), ForKey(), Kind, normalize(), Value(), Account(), hasDigitRun() (+18 more)

### Community 5 - "context.Context"
Cohesion: 0.31
Nodes (12): TestActorMiddleware(), ActorMiddleware(), FromContext(), HasScope(), WithPrincipal(), TestActorSlot(), Actor(), Attribution() (+4 more)

### Community 6 - "Releasing"
Cohesion: 0.18
Nodes (10): 1. Prerequisites (one-time, per machine that will `go get` this module), 2. Decide the version bump, 3. Tag and push, 4. Publish the GitHub release, 5. Update consumers, Automated (default), Future automation, Major releases (+2 more)

### Community 7 - "time.Time"
Cohesion: 0.05
Nodes (44): Response, Status, ErasureRequest, ErasureRequestID, LegalOrder, LegalOrderEvent, LegalOrderTarget, LiftLegalOrderResponse (+36 more)

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
Cohesion: 0.39
Nodes (14): keyInfo, Rule, RuleFunc, TestServiceRules(), DefaultRules(), isRef(), keyOf(), KeyRedact() (+6 more)

### Community 15 - "Refs"
Cohesion: 0.22
Nodes (6): Compute(), Kind, Refs, cacheKey, personKind(), TestPersonKind()

### Community 16 - "Parse"
Cohesion: 0.17
Nodes (12): Parse(), golden(), TestCanonical_Golden(), TestDetailsInto(), TestFromResponse(), TestLogFields_Golden(), TestParse(), TestSchemaIsACopy() (+4 more)

### Community 18 - "test-next-version.sh"
Cohesion: 0.70
Nodes (4): c(), expect(), test-next-version.sh script, tg()

### Community 19 - "Error"
Cohesion: 0.06
Nodes (31): draftVector, encode(), Error, New(), sortedJSON(), init(), TestError_IsAGoError(), TestShortID() (+23 more)

### Community 20 - "field.go"
Cohesion: 0.24
Nodes (10): ErrUnknownField, child(), contains(), Registry, lookupRef(), mergeProperty(), resolveNode(), splitPointer() (+2 more)

### Community 21 - "auditref_test.go"
Cohesion: 0.27
Nodes (10): New(), Parse(), TestAllAndParse(), TestInstitution(), TestKeyed(), TestKeyedFailsClosed(), TestPlainInDevelopment(), TestRefsCache() (+2 more)

### Community 24 - "approval_test.go"
Cohesion: 0.29
Nodes (11): Check(), Decisions(), cases(), generate(), TestCheckRefusals(), TestDecisions(), TestDecisionsGolden(), Decision (+3 more)

### Community 26 - "fieldMap"
Cohesion: 0.14
Nodes (20): AccessLogOption, fieldMap(), observed(), spanCtx(), TestAccessLog(), TestAccessLogServerError(), TestFromContext(), TestScrubCore() (+12 more)

### Community 27 - "telemetry_test.go"
Cohesion: 0.08
Nodes (18): Config, failing, Middleware(), Transport(), isURL(), metricOptions(), Propagator(), Setup() (+10 more)

### Community 28 - "testing.T"
Cohesion: 0.11
Nodes (25): TestDraftCanonical_EmptyTermsAreNull(), TestCanonical_Properties(), TestLogFields_ShortFormsOnly(), New(), TestDevOIDC_ClientCredentials(), TestDevOIDC_SigningSessionClaims(), TestDevOIDC_TokenRequiresSubAndEmail(), TestDevOIDC_TokenVerifiesWithGoOIDC() (+17 more)

### Community 29 - "smtptest.go"
Cohesion: 0.16
Nodes (6): readPart(), WithMaxSize(), WithOnMessage(), Content, Message, Option

### Community 30 - "Principal"
Cohesion: 0.29
Nodes (7): ctxKey, fakeVerifier, devPrincipal(), looksLikeJWT(), Middleware(), Principal, TokenVerifier

### Community 31 - "Fingerprint"
Cohesion: 0.09
Nodes (23): ExpectedHash(), ForParty(), Signing(), TestSigningDetails(), SigningDetails, Classify(), FingerprintOf(), Fingerprint (+15 more)

### Community 32 - "net/http.Handler"
Cohesion: 0.33
Nodes (5): Newf(), newRequestID(), RequestID(), TestRequireScope(), RequireScope()

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
Cohesion: 0.05
Nodes (23): CodeForStatus(), TestCodeForStatus(), ValidCode(), requestIDKey, signing, parseCase, main(), splitOrigins() (+15 more)

### Community 38 - "Issuer"
Cohesion: 0.09
Nodes (12): RequestIDFrom(), Write(), Writer, authCode, Issuer, codeStore, Issuer, httpHandlerFunc() (+4 more)

### Community 39 - "agreementsig_test.go"
Cohesion: 0.19
Nodes (21): AgreementHash(), agreementCases(), agreementJSON(), flipS(), readJSON(), sha256Sum(), signatureCases(), spki() (+13 more)

### Community 41 - "notify_test.go"
Cohesion: 0.27
Nodes (11): NewMockSender(), IsDevelopment(), TestIsDevelopment(), TestMockSender_RecordsOncePerIdempotencyKey(), TestMockSender_SentIsACopy(), TestMockSender_ValidatesAndFails(), TestNewMockSender_RefusedOutsideDevelopment(), TestValidate_AcceptsAWellFormedMessage() (+3 more)

### Community 42 - "legalorders/legalorders.go"
Cohesion: 0.13
Nodes (14): CaseRequest, EventRequest, TestEventAndLiftValidate(), TestRecordRequestValidate(), TestWireNames(), valid(), ValidCase(), validNotifyAfter() (+6 more)

### Community 43 - "http_test.go"
Cohesion: 0.36
Nodes (8): NewHTTPSender(), staticToken(), TestHTTPSender_RefusesAnInvalidMessageWithoutCalling(), TestHTTPSender_SendsAndReadsTheReceipt(), TestHTTPSender_Status(), TestHTTPSender_StatusAnswersToTheRelay(), TestHTTPSender_TransportFailuresAreRetryable(), TestNewHTTPSender_Refusals()

### Community 44 - ".Verify"
Cohesion: 0.22
Nodes (10): signer(), TestOIDCVerifierReadsClaims(), TestOIDCVerifierRefuses(), token(), NewOIDCVerifierWithKeys(), Sign(), TestSignAndVerify_RoundTrip(), TestVerify_MalformedHexFails() (+2 more)

### Community 45 - "authority.go"
Cohesion: 0.13
Nodes (19): TestClientPolicyFilter(), TestMiddleware(), TestParseClientPolicy(), ClientPolicy, WithClientPolicy(), WithDevBypass(), WithLogger(), Option (+11 more)

### Community 47 - "bench_test.go"
Cohesion: 0.35
Nodes (9): benchCore(), benchLine(), BenchmarkRawZap(), BenchmarkShapingCore(), benchRefs(), TestShapingAllocs(), NewRefScrubCore(), NewScrubCore() (+1 more)

### Community 48 - "field_test.go"
Cohesion: 0.39
Nodes (6): loadFields(), TestField_Errors(), TestField_KeepsOwnDescriptionAndRefKeywords(), TestField_RefFailures(), TestField_Resolves(), TestField_Root()

### Community 49 - "HTTPSender"
Cohesion: 0.17
Nodes (8): refusal(), WithHTTPClient(), HTTPOption, HTTPSender, Receipt, Sender, Status, TokenFunc

### Community 52 - "audit.go"
Cohesion: 0.27
Nodes (9): ActiveHolds, AuditIdentity, AuditRef, AuditRefsOf, LegalHold, Resolution, ResolveAuditRefs, ResolvedRef (+1 more)

### Community 53 - "fields.go"
Cohesion: 0.31
Nodes (7): accountNumber, Account(), Party(), Subject(), User(), personID, TestTypedFields()

### Community 54 - "authorize_test.go"
Cohesion: 0.14
Nodes (15): claims, NewOIDCVerifier(), uniqueScopes(), OIDCVerifier, authorizeQuery(), browserIssuer(), challengeOf(), exchange() (+7 more)

### Community 55 - "build"
Cohesion: 0.32
Nodes (7): Config, build(), New(), TestNew(), TestNewIsNotSampled(), zapConfig(), TestNewWithRefs()

### Community 56 - "LoadDirectory"
Cohesion: 0.25
Nodes (7): Registry, LoadDirectory(), TestLoadDirectory_CompilesSchemas(), TestLoadDirectory_MissingDirectory(), TestValidate_InvalidPayloadReportsFailures(), TestValidate_UnknownType(), TestValidate_ValidPayload()

### Community 57 - "New"
Cohesion: 0.19
Nodes (15): New(), TestClient_CreatePendingInvitation(), TestClient_CreatePendingInvitation_ServerError(), TestClient_GetByEmail_ServerError(), TestClient_GetByOktaSub_Found(), TestClient_GetByToken_NotFound(), TestClient_ManagesIdentity(), TestClient_ManagesIdentity_ServerError() (+7 more)

### Community 59 - "NewTypst"
Cohesion: 0.13
Nodes (12): Fake, CheckFonts(), NewTypst(), TestCheckFonts(), TestStr(), Content(), Lit(), Str() (+4 more)

### Community 66 - "signature.go"
Cohesion: 0.18
Nodes (6): Algorithm, ecdsaSig, parseCanonicalECDSA(), parseEd25519(), parseP256(), Verify()

## Knowledge Gaps
- **80 isolated node(s):** `signing`, `requestIDKey`, `Error`, `ctxKey`, `github.com/vdatacloud/daml-escrow-commons` (+75 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 187 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **9 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Issuer` connect `Issuer` to `canonical.go`, `signature.go`, `go_pkg_strings`, `testing.T`, `smtptest.go`?**
  _High betweenness centrality (0.034) - this node is a cross-community bridge._
- **What connects `signing`, `requestIDKey`, `Error` to the rest of the system?**
  _80 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `GitHub & CI Workflow Skill` be split into smaller, more focused modules?**
  _Cohesion score 0.1111111111111111 - nodes in this community are weakly interconnected._
- **Why does `MockServer` connect `MockServer` to `Token`, `smtptest_test.go`, `mock.go`, `smtptest.go`?**
  _High betweenness centrality (0.019) - this node is a cross-community bridge._
- **Should `canonical.go` be split into smaller, more focused modules?**
  _Cohesion score 0.14461538461538462 - nodes in this community are weakly interconnected._
- **Why does `Refs` connect `Refs` to `context.Context`, `go_pkg_strings`, `scrub.go`, `bench_test.go`, `auditref_test.go`, `build`, `fieldMap`?**
  _High betweenness centrality (0.019) - this node is a cross-community bridge._
- **Should `redact.go` be split into smaller, more focused modules?**
  _Cohesion score 0.14039408866995073 - nodes in this community are weakly interconnected._