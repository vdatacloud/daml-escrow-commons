# Graph Report - daml-escrow-commons  (2026-10-03)

## Corpus Check
- 64 files · ~20,527 words
- Verdict: corpus is large enough that graph structure adds value.
- Unclassified: 3 file(s) not represented in the graph (top: (none) 3)

## Summary
- 418 nodes · 730 edges · 32 communities (27 shown, 5 thin omitted)
- Extraction: 90% EXTRACTED · 10% INFERRED · 0% AMBIGUOUS · INFERRED: 72 edges (avg confidence: 0.85)
- Token cost: 0 input · 0 output

## Graph Freshness
- Built from commit: `e99cea0c`
- Run `git rev-parse HEAD` and compare to check if the graph is stale.
- Run `graphify update .` after code changes (no API cost).

## Community Hubs (Navigation)
- GitHub & CI Workflow Skill
- testing.T
- widget.json
- Repository Guardrails
- hmacsig.go
- Fingerprint
- Releasing
- metering.go
- CLAUDE.md
- /commons-contribution
- daml-escrow-commons
- .claude/CLAUDE.md
- install-git-hooks.sh
- github.com/vdatacloud/daml-escrow-commons
- LoadDirectory
- Client
- cantonid.go
- next-version.sh
- test-next-version.sh
- properties
- definitions
- Error
- encode
- properties
- properties
- error.schema.json
- status
- code
- service
- details
- hint
- stage

## God Nodes (most connected - your core abstractions)
1. `Fingerprint` - 13 edges
2. `PartyID` - 12 edges
3. `GitHub & CI Workflow Skill` - 12 edges
4. `New()` - 11 edges
5. `TestSigningDetails()` - 11 edges
6. `Hash` - 11 edges
7. `RequireNonEmpty()` - 11 edges
8. `FromLedgerError()` - 10 edges
9. `ParsePartyID()` - 9 edges
10. `Client` - 9 edges

## Surprising Connections (you probably didn't know these)
- `init()` --calls--> `New()`  [INFERRED]
  apierror/apierror_test.go → apierror/apierror.go
- `TestLogFields_Golden()` --calls--> `encode()`  [INFERRED]
  apierror/apierror_test.go → apierror/apierror.go
- `TestWriter()` --calls--> `RequestID()`  [INFERRED]
  apierror/apierror_test.go → apierror/http.go
- `TestSigningDetails()` --calls--> `ParseLedgerError()`  [INFERRED]
  apierror/canton/canton_test.go → apierror/canton/canton.go
- `SigningDetails` --references--> `Fingerprint`  [EXTRACTED]
  apierror/canton/canton.go → cantonid/cantonid.go

## Import Cycles
- None detected.

## Communities (32 total, 5 thin omitted)

### Community 0 - "GitHub & CI Workflow Skill"
Cohesion: 0.11
Nodes (17): 10. Useful Reference Commands, 11. References, 1. Pre-Commit Local Verification (MANDATORY — do this before every commit), 2. Branching Rules, 3. Commit Standards, 4. Staging & Pushing Changes, 5. Pull Request Creation, 6. CI Pipeline Overview (+9 more)

### Community 1 - "testing.T"
Cohesion: 0.08
Nodes (43): CodeForStatus(), New(), Parse(), golden(), init(), TestCanonical_Golden(), TestCanonical_Properties(), TestCodeForStatus() (+35 more)

### Community 2 - "widget.json"
Cohesion: 0.17
Nodes (11): minLength, type, properties, name, quantity, minimum, type, required (+3 more)

### Community 3 - "Repository Guardrails"
Cohesion: 0.22
Nodes (8): Branch Protection Strategy, Branching Strategy, CI Requirements, Code Review Requirements, Commit Standard, Pre-Commit Verification, Pull Request Rules, Repository Guardrails

### Community 4 - "hmacsig.go"
Cohesion: 0.23
Nodes (8): Sign(), TestSignAndVerify_RoundTrip(), TestVerify_MalformedHexFails(), TestVerify_TamperedMessageFails(), TestVerify_WrongSecretFails(), Verify(), Ref(), TestRef()

### Community 5 - "Fingerprint"
Cohesion: 0.08
Nodes (25): ExpectedHash(), ForParty(), FromLedgerError(), GRPCCodeName(), Signing(), TestSigningDetails(), LedgerError, SigningDetails (+17 more)

### Community 6 - "Releasing"
Cohesion: 0.20
Nodes (9): 1. Prerequisites (one-time, per machine that will `go get` this module), 2. Decide the version bump, 3. Tag and push, 4. Publish the GitHub release, 5. Update consumers, Automated (default), Future automation, Manual It follows [Semantic Versioning](https://semver.org/) and (+1 more)

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

### Community 14 - "LoadDirectory"
Cohesion: 0.25
Nodes (7): Registry, LoadDirectory(), TestLoadDirectory_CompilesSchemas(), TestLoadDirectory_MissingDirectory(), TestValidate_InvalidPayloadReportsFailures(), TestValidate_UnknownType(), TestValidate_ValidPayload()

### Community 16 - "cantonid.go"
Cohesion: 0.10
Nodes (8): newRequestID(), RequestID(), RequestIDFrom(), requestIDKey, Kind, parseCase, ErrUnknownType, ValidationError

### Community 18 - "test-next-version.sh"
Cohesion: 0.70
Nodes (4): c(), expect(), test-next-version.sh script, tg()

### Community 19 - "properties"
Cohesion: 0.08
Nodes (26): description, items, type, $ref, type, description, type, $ref (+18 more)

### Community 20 - "definitions"
Cohesion: 0.10
Nodes (21): definitions, fingerprint, hash, partyId, signingDetails, upstream, description, pattern (+13 more)

### Community 21 - "Error"
Cohesion: 0.15
Nodes (6): Error, Newf(), Detail, Write(), Upstream, Writer

### Community 22 - "encode"
Cohesion: 0.17
Nodes (11): encode(), sortedJSON(), TestShortID(), Field, decodeInto(), Error, ShortID(), ShortIDs() (+3 more)

### Community 23 - "properties"
Cohesion: 0.17
Nodes (12): maxLength, type, pattern, type, description, type, cause, grpcCode (+4 more)

### Community 24 - "properties"
Cohesion: 0.20
Nodes (10): description, minLength, type, properties, error, requestId, upstream, pattern (+2 more)

### Community 25 - "error.schema.json"
Cohesion: 0.22
Nodes (8): additionalProperties, allOf, description, $id, required, $schema, title, type

### Community 26 - "status"
Cohesion: 0.40
Nodes (5): status, description, maximum, minimum, type

### Community 27 - "code"
Cohesion: 0.50
Nodes (4): description, pattern, type, code

### Community 28 - "service"
Cohesion: 0.50
Nodes (4): service, description, minLength, type

### Community 29 - "details"
Cohesion: 0.67
Nodes (3): description, type, details

### Community 30 - "hint"
Cohesion: 0.67
Nodes (3): description, type, hint

### Community 31 - "stage"
Cohesion: 0.67
Nodes (3): stage, description, type

## Knowledge Gaps
- **125 isolated node(s):** `signing`, `requestIDKey`, `$schema`, `$id`, `title` (+120 more)
  These have ≤1 connection - possible missing edges or undocumented components. (Counts symbols only; 162 node(s) total have ≤1 connection when file, concept and rationale nodes are included.)
- **5 thin communities (<3 nodes) omitted from report** — run `graphify query` to explore isolated nodes.

## Suggested Questions
_Questions this graph is uniquely positioned to answer:_

- **Why does `definitions` connect `definitions` to `error.schema.json`?**
  _High betweenness centrality (0.034) - this node is a cross-community bridge._
- **Why does `properties` connect `properties` to `definitions`?**
  _High betweenness centrality (0.024) - this node is a cross-community bridge._
- **Why does `signingDetails` connect `definitions` to `properties`?**
  _High betweenness centrality (0.024) - this node is a cross-community bridge._
- **Are the 3 inferred relationships involving `Fingerprint` (e.g. with `TestJSON()` and `TestPartyID()`) actually correct?**
  _`Fingerprint` has 3 INFERRED edges - model-reasoned connections that need verification._
- **What connects `signing`, `requestIDKey`, `$schema` to the rest of the system?**
  _125 weakly-connected nodes found - possible documentation gaps or missing edges._
- **Should `GitHub & CI Workflow Skill` be split into smaller, more focused modules?**
  _Cohesion score 0.1111111111111111 - nodes in this community are weakly interconnected._
- **Should `testing.T` be split into smaller, more focused modules?**
  _Cohesion score 0.08081632653061224 - nodes in this community are weakly interconnected._