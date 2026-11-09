# Graph Report - daml-escrow-commons  (2026-11-09)

## Corpus Check
- 133 files · ~61,867 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 9 file(s) not represented in the graph (top: .ttf 5, (none) 4)

## Summary
- 950 nodes · 2402 edges · 55 communities (47 shown, 8 thin omitted)
- Extraction: 90% EXTRACTED · 10% INFERRED · 0% AMBIGUOUS · INFERRED: 238 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `78f297c5`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- GitHub & CI Workflow Skill
- canton_test.go
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
- go_pkg_encoding_json
- Refs
- notify.go
- next-version.sh
- test-next-version.sh
- shortid.go
- field.go
- apierror.go
- canton.go
- approval_test.go
- logging_test.go
- telemetry_test.go
- apierror_test.go
- smtptest.go
- go_pkg_strings
- notify/http.go
- MockServer
- smtptest_test.go
- mock.go
- deliver.go
- go_pkg_net_http
- Issuer
- testing.T
- go_pkg_time
- notify_test.go
- legalorders/legalorders.go
- http_test.go
- authn_test.go
- auditref_test.go
- time.Time
- audit.go
- Typst
- Client
- erasure.go
- Str
- fields.go
- logging.go
- bench_test.go

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
- `SigningDetails` --references--> `PartyID`  [EXTRACTED]
  apierror/canton/canton.go → cantonid/cantonid.go
- `ForParty()` --references--> `PartyID`  [EXTRACTED]
  apierror/canton/canton.go → cantonid/cantonid.go
- `TestSigningDetails()` --calls--> `MustParsePartyID()`  [EXTRACTED]
  apierror/canton/canton_test.go → cantonid/cantonid.go

## Import Cycles
- None detected.

## Communities (55 total, 8 thin omitted)

### Community 0 - "GitHub & CI Workflow Skill"
Cohesion: 0.11
Nodes (17): 10. Useful Reference Commands, 11. References, 1. Pre-Commit Local Verification (MANDATORY — do this before every commit), 2. Branching Rules, 3. Commit Standards, 4. Staging & Pushing Changes, 5. Pull Request Creation, 6. CI Pipeline Overview (+9 more)

### Community 1 - "canton_test.go"
Cohesion: 0.13
Nodes (9): ParseLedgerError(), TestFromError(), TestFromLedgerError_Golden(), TestFromLedgerError_Mapping(), TestLedgerError(), fixture, ErrUnknownType, Registry (+1 more)

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
Cohesion: 0.11
Nodes (17): Classify(), FingerprintOf(), Fingerprint, PartyID, MustParsePartyID(), ParseFingerprint(), ParsePartyID(), Short() (+9 more)

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

### Community 14 - "go_pkg_encoding_json"
Cohesion: 0.13
Nodes (22): AgreementCanonical(), AgreementMessage(), checkID(), invalid(), TestDraftVectors(), AgreementVersion, Amends, decodeJSON() (+14 more)

### Community 15 - "Refs"
Cohesion: 0.24
Nodes (4): Compute(), Kind, Refs, cacheKey

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

### Community 23 - "canton.go"
Cohesion: 0.23
Nodes (10): ExpectedHash(), ForParty(), FromLedgerError(), GRPCCodeName(), Signing(), TestSigningDetails(), LedgerError, SigningDetails (+2 more)

### Community 24 - "approval_test.go"
Cohesion: 0.29
Nodes (11): Check(), Decisions(), cases(), generate(), TestCheckRefusals(), TestDecisions(), TestDecisionsGolden(), Decision (+3 more)

### Community 26 - "logging_test.go"
Cohesion: 0.07
Nodes (49): Newf(), newRequestID(), RequestID(), RequestIDFrom(), requestIDKey, TestActorMiddleware(), TestClientPolicyFilter(), TestMiddleware() (+41 more)

### Community 27 - "telemetry_test.go"
Cohesion: 0.07
Nodes (19): Token, Config, failing, Middleware(), Transport(), isURL(), metricOptions(), Propagator() (+11 more)

### Community 28 - "apierror_test.go"
Cohesion: 0.15
Nodes (20): New(), Parse(), golden(), init(), TestCanonical_Golden(), TestCanonical_Properties(), TestDetailsInto(), TestError_IsAGoError() (+12 more)

### Community 29 - "smtptest.go"
Cohesion: 0.19
Nodes (6): readPart(), WithMaxSize(), WithOnMessage(), Content, Message, Option

### Community 30 - "go_pkg_strings"
Cohesion: 0.15
Nodes (11): TestEventAndLiftValidate(), TestWireNames(), personKind(), refLogger(), TestPersonKind(), TestRefScrubCore(), TestRefScrubCoreClear(), TestRefScrubCoreDevelopment() (+3 more)

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

### Community 36 - "deliver.go"
Cohesion: 0.32
Nodes (9): Artifact, checkName(), Deliver(), artifacts(), TestDeliver_ManyIsDeterministicZip(), TestDeliver_One(), TestDeliver_RefusesBadNames(), Zip() (+1 more)

### Community 37 - "go_pkg_net_http"
Cohesion: 0.06
Nodes (34): claims, ctxKey, NewOIDCVerifier(), NewOIDCVerifierWithKeys(), uniqueScopes(), OIDCVerifier, Principal, TokenVerifier (+26 more)

### Community 38 - "Issuer"
Cohesion: 0.10
Nodes (11): Write(), Writer, authCode, Issuer, codeStore, Issuer, httpHandlerFunc(), writeJSON() (+3 more)

### Community 39 - "testing.T"
Cohesion: 0.06
Nodes (63): AgreementHash(), agreementCases(), agreementJSON(), flipS(), readJSON(), sha256Sum(), signatureCases(), spki() (+55 more)

### Community 40 - "go_pkg_time"
Cohesion: 0.24
Nodes (7): CheckFonts(), NewTypst(), TestCheckFonts(), TestFakeIsDeterministic(), TestTypst_ByteIdenticalAndEscaped(), TestTypst_GoldenHashAcrossPlatforms(), TestTypst_RefusesOtherVersion()

### Community 41 - "notify_test.go"
Cohesion: 0.27
Nodes (11): NewMockSender(), IsDevelopment(), TestIsDevelopment(), TestMockSender_RecordsOncePerIdempotencyKey(), TestMockSender_SentIsACopy(), TestMockSender_ValidatesAndFails(), TestNewMockSender_RefusedOutsideDevelopment(), TestValidate_AcceptsAWellFormedMessage() (+3 more)

### Community 42 - "legalorders/legalorders.go"
Cohesion: 0.14
Nodes (12): CaseRequest, EventRequest, TestRecordRequestValidate(), valid(), ValidCase(), validNotifyAfter(), Lifted, LiftRequest (+4 more)

### Community 43 - "http_test.go"
Cohesion: 0.31
Nodes (8): NewHTTPSender(), staticToken(), TestHTTPSender_RefusesAnInvalidMessageWithoutCalling(), TestHTTPSender_SendsAndReadsTheReceipt(), TestHTTPSender_Status(), TestHTTPSender_StatusAnswersToTheRelay(), TestHTTPSender_TransportFailuresAreRetryable(), TestNewHTTPSender_Refusals()

### Community 44 - "authn_test.go"
Cohesion: 0.31
Nodes (4): signer(), TestOIDCVerifierReadsClaims(), TestOIDCVerifierRefuses(), token()

### Community 45 - "auditref_test.go"
Cohesion: 0.31
Nodes (10): New(), Parse(), TestAllAndParse(), TestInstitution(), TestKeyed(), TestKeyedFailsClosed(), TestPlainInDevelopment(), TestRefsCache() (+2 more)

### Community 46 - "time.Time"
Cohesion: 0.21
Nodes (8): LegalOrder, LegalOrderEvent, LegalOrderTarget, LiftLegalOrderResponse, ListLegalOrdersRequest, OrderContact, RecordedLegalOrder, Renderer

### Community 47 - "audit.go"
Cohesion: 0.27
Nodes (9): ActiveHolds, AuditIdentity, AuditRef, AuditRefsOf, LegalHold, Resolution, ResolveAuditRefs, ResolvedRef (+1 more)

### Community 49 - "Client"
Cohesion: 0.38
Nodes (3): Client, Identity, PendingInvitation

### Community 50 - "erasure.go"
Cohesion: 0.29
Nodes (6): ErasureRequest, ErasureRequestID, ListErasureRequests, RefuseErasureRequest, SubmitErasureRequest, VerifyErasureRequest

### Community 51 - "Str"
Cohesion: 0.70
Nodes (4): TestStr(), Content(), Lit(), Str()

### Community 53 - "fields.go"
Cohesion: 0.31
Nodes (7): accountNumber, Account(), Party(), Subject(), User(), personID, TestTypedFields()

### Community 55 - "logging.go"
Cohesion: 0.27
Nodes (7): Config, build(), New(), TestNew(), TestNewIsNotSampled(), zapConfig(), TestNewWithRefs()

### Community 56 - "bench_test.go"
Cohesion: 0.54
Nodes (6): benchCore(), benchLine(), BenchmarkRawZap(), BenchmarkShapingCore(), benchRefs(), TestShapingAllocs()

## Knowledge Gaps
- **77 isolated node(s):** `signing`, `requestIDKey`, `Error`, `ctxKey`, `github.com/vdatacloud/daml-escrow-commons` (+72 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 183 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **8 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Issuer` connect `Issuer` to `go_pkg_net_http`, `go_pkg_fmt`, `testing.T`, `go_pkg_time`, `authn_test.go`, `go_pkg_encoding_json`, `shortid.go`, `go_pkg_strings`, `go_pkg_errors`?**
  _High betweenness centrality (0.036) - this node is a cross-community bridge._
- **What connects `signing`, `requestIDKey`, `Error` to the rest of the system?**
  _77 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `GitHub & CI Workflow Skill` be split into smaller, more focused modules?**
  _Cohesion score 0.1111111111111111 - nodes in this community are weakly interconnected._
- **Why does `MockServer` connect `MockServer` to `telemetry_test.go`, `smtptest_test.go`, `mock.go`, `smtptest.go`?**
  _High betweenness centrality (0.020) - this node is a cross-community bridge._
- **Should `canton_test.go` be split into smaller, more focused modules?**
  _Cohesion score 0.13157894736842105 - nodes in this community are weakly interconnected._
- **Why does `Refs` connect `Refs` to `scrub.go`, `auditref_test.go`, `logging.go`, `bench_test.go`, `logging_test.go`, `go_pkg_strings`, `go_pkg_errors`?**
  _High betweenness centrality (0.020) - this node is a cross-community bridge._
- **Should `scrub.go` be split into smaller, more focused modules?**
  _Cohesion score 0.08771929824561403 - nodes in this community are weakly interconnected._