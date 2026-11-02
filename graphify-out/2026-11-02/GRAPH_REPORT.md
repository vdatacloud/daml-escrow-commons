# Graph Report - daml-escrow-commons  (2026-11-01)

## Corpus Check
- 101 files · ~44,613 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 9 file(s) not represented in the graph (top: .ttf 5, (none) 4)

## Summary
- 652 nodes · 1540 edges · 34 communities (25 shown, 9 thin omitted)
- Extraction: 89% EXTRACTED · 11% INFERRED · 0% AMBIGUOUS · INFERRED: 165 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `e8dfc27d`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- GitHub & CI Workflow Skill
- testing.T
- widget.json
- Repository Guardrails
- redact.go
- cantonid.go
- Releasing
- go_pkg_fmt
- CLAUDE.md
- /commons-contribution
- daml-escrow-commons
- .claude/CLAUDE.md
- install-git-hooks.sh
- github.com/vdatacloud/daml-escrow-commons
- agreementsig_test.go
- go_pkg_testing
- telemetry_test.go
- next-version.sh
- test-next-version.sh
- canonical.go
- field.go
- apierror_test.go
- apierror/http.go
- approval_test.go
- logging_test.go
- telemetry.go
- net/http.ResponseWriter
- transport
- go_pkg_strings
- deliver.go
- go_pkg_encoding_json
- attrs

## God Nodes (most connected - your core abstractions)
1. `Refs` - 17 edges
2. `Fingerprint` - 13 edges
3. `New()` - 12 edges
4. `PartyID` - 12 edges
5. `AccessLog()` - 12 edges
6. `GitHub & CI Workflow Skill` - 12 edges
7. `AgreementVersion` - 11 edges
8. `New()` - 11 edges
9. `TestSigningDetails()` - 11 edges
10. `Kind` - 11 edges

## Surprising Connections (you probably didn't know these)
- `TestAccessLog()` --calls--> `RequestID()`  [EXTRACTED]
  logging/logging_test.go → apierror/http.go
- `Config` --references--> `Refs`  [EXTRACTED]
  logging/logging.go → auditref/auditref.go
- `refLogger()` --references--> `Refs`  [EXTRACTED]
  logging/person_test.go → auditref/auditref.go
- `NewRefScrubCore()` --references--> `Refs`  [EXTRACTED]
  logging/scrub.go → auditref/auditref.go
- `scrubCore` --references--> `Refs`  [EXTRACTED]
  logging/scrub.go → auditref/auditref.go

## Import Cycles
- None detected.

## Communities (34 total, 9 thin omitted)

### Community 0 - "GitHub & CI Workflow Skill"
Cohesion: 0.11
Nodes (17): 10. Useful Reference Commands, 11. References, 1. Pre-Commit Local Verification (MANDATORY — do this before every commit), 2. Branching Rules, 3. Commit Standards, 4. Staging & Pushing Changes, 5. Pull Request Creation, 6. CI Pipeline Overview (+9 more)

### Community 1 - "testing.T"
Cohesion: 0.10
Nodes (35): Sign(), TestSignAndVerify_RoundTrip(), TestVerify_MalformedHexFails(), TestVerify_TamperedMessageFails(), TestVerify_WrongSecretFails(), Verify(), New(), TestClient_CreatePendingInvitation() (+27 more)

### Community 2 - "widget.json"
Cohesion: 0.17
Nodes (11): minLength, type, properties, name, quantity, minimum, type, required (+3 more)

### Community 3 - "Repository Guardrails"
Cohesion: 0.22
Nodes (8): Branch Protection Strategy, Branching Strategy, CI Requirements, Code Review Requirements, Commit Standard, Pre-Commit Verification, Pull Request Rules, Repository Guardrails

### Community 4 - "redact.go"
Cohesion: 0.10
Nodes (28): NewRefScrubCore(), NewScrubCore(), scrubField(), scrubNumber(), stringArray(), scrubCore, ForKey(), Kind (+20 more)

### Community 5 - "cantonid.go"
Cohesion: 0.07
Nodes (34): ExpectedHash(), ForParty(), FromError(), FromLedgerError(), GRPCCodeName(), ParseLedgerError(), Signing(), TestFromError() (+26 more)

### Community 6 - "Releasing"
Cohesion: 0.18
Nodes (10): 1. Prerequisites (one-time, per machine that will `go get` this module), 2. Decide the version bump, 3. Tag and push, 4. Publish the GitHub release, 5. Update consumers, Automated (default), Future automation, Major releases (+2 more)

### Community 7 - "go_pkg_fmt"
Cohesion: 0.10
Nodes (21): ChargeBearer, CommandOutcome, LedgerCommandEvent, requireBaseUnits(), requireDecimal(), NetworkFeeEvent, NetworkFeePayer, Rail (+13 more)

### Community 8 - "CLAUDE.md"
Cohesion: 0.29
Nodes (5): Architecture, Commands, Conventions carried over from `daml-escrow` / `daml-escrow-cms`, Project Overview, Relationship to the other repos

### Community 9 - "/commons-contribution"
Cohesion: 0.33
Nodes (5): /commons-contribution, How to add something, Removing something, When NOT to add something, When to promote code into commons

### Community 10 - "daml-escrow-commons"
Cohesion: 0.33
Nodes (5): daml-escrow-commons, Status, Using this module, What's here, What's NOT here, and won't be

### Community 14 - "agreementsig_test.go"
Cohesion: 0.07
Nodes (44): AgreementCanonical(), AgreementHash(), AgreementMessage(), checkID(), invalid(), agreementCases(), agreementJSON(), flipS() (+36 more)

### Community 15 - "go_pkg_testing"
Cohesion: 0.09
Nodes (15): Fake, CheckFonts(), NewTypst(), TestCheckFonts(), TestFakeIsDeterministic(), TestStr(), Content(), Lit() (+7 more)

### Community 16 - "telemetry_test.go"
Cohesion: 0.20
Nodes (8): failing, Middleware(), Transport(), recording(), TestMiddleware(), TestMiddlewareServerError(), TestTransportError(), TestTransportPropagates()

### Community 18 - "test-next-version.sh"
Cohesion: 0.70
Nodes (4): c(), expect(), test-next-version.sh script, tg()

### Community 19 - "canonical.go"
Cohesion: 0.18
Nodes (11): lessUTF16(), writeCanonical(), writeNumber(), writeString(), numberMode, TestShortID(), Field, Error (+3 more)

### Community 20 - "field.go"
Cohesion: 0.24
Nodes (10): ErrUnknownField, child(), contains(), Registry, lookupRef(), mergeProperty(), resolveNode(), splitPointer() (+2 more)

### Community 21 - "apierror_test.go"
Cohesion: 0.06
Nodes (30): CodeForStatus(), encode(), Error, New(), Newf(), Parse(), sortedJSON(), golden() (+22 more)

### Community 23 - "apierror/http.go"
Cohesion: 0.22
Nodes (8): TestFromResponse(), FromResponse(), newRequestID(), RequestID(), Write(), requestIDKey, Writer, accessLog

### Community 24 - "approval_test.go"
Cohesion: 0.26
Nodes (11): Check(), Decisions(), cases(), generate(), TestCheckRefusals(), TestDecisions(), TestDecisionsGolden(), Decision (+3 more)

### Community 26 - "logging_test.go"
Cohesion: 0.09
Nodes (37): RequestIDFrom(), AccessLogOption, actorKey, actorSlot, Config, loggerKey, build(), New() (+29 more)

### Community 27 - "telemetry.go"
Cohesion: 0.15
Nodes (8): Config, isURL(), metricOptions(), Propagator(), Setup(), TestSetup(), TestSetupSampleRatio(), traceOptions()

### Community 30 - "go_pkg_strings"
Cohesion: 0.12
Nodes (14): Compute(), Kind, Refs, New(), Parse(), TestAllAndParse(), TestInstitution(), TestKeyed() (+6 more)

### Community 31 - "deliver.go"
Cohesion: 0.29
Nodes (9): Artifact, checkName(), Deliver(), artifacts(), TestDeliver_ManyIsDeterministicZip(), TestDeliver_One(), TestDeliver_RefusesBadNames(), Zip() (+1 more)

### Community 32 - "go_pkg_encoding_json"
Cohesion: 0.25
Nodes (3): Client, Identity, PendingInvitation

## Knowledge Gaps
- **62 isolated node(s):** `signing`, `requestIDKey`, `Error`, `github.com/vdatacloud/daml-escrow-commons`, `loggerKey` (+57 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 142 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **9 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `Refs` connect `go_pkg_strings` to `logging_test.go`, `redact.go`?**
  _High betweenness centrality (0.024) - this node is a cross-community bridge._
- **What connects `signing`, `requestIDKey`, `Error` to the rest of the system?**
  _62 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `GitHub & CI Workflow Skill` be split into smaller, more focused modules?**
  _Cohesion score 0.1111111111111111 - nodes in this community are weakly interconnected._
- **Why does `Error` connect `apierror_test.go` to `cantonid.go`?**
  _High betweenness centrality (0.022) - this node is a cross-community bridge._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.09966777408637874 - nodes in this community are weakly interconnected._
- **Why does `refLogger()` connect `logging_test.go` to `testing.T`, `redact.go`, `go_pkg_strings`?**
  _High betweenness centrality (0.018) - this node is a cross-community bridge._
- **Should `redact.go` be split into smaller, more focused modules?**
  _Cohesion score 0.10452961672473868 - nodes in this community are weakly interconnected._