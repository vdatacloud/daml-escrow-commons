package logging

import (
	"os"
	"strings"
	"testing"

	"github.com/vdatacloud/daml-escrow-commons/auditref"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func refLogger(t *testing.T, env string) (*zap.Logger, *observer.ObservedLogs, *auditref.Refs) {
	t.Helper()
	cfg := auditref.Config{Environment: env}
	if !auditref.Plain(env) {
		cfg.Keys, cfg.Current = map[string][]byte{"r1": []byte("ref-key")}, "r1"
	}
	refs, err := auditref.New(cfg)
	if err != nil {
		t.Fatal(err)
	}
	core, logs := observer.New(zapcore.DebugLevel)
	return zap.New(NewRefScrubCore(core, refs)), logs, refs
}

func TestPersonKind(t *testing.T) {
	for key, want := range map[string]auditref.Kind{
		"oktaSub": auditref.Subject, "sub": auditref.Subject, "callerOktaSub": auditref.Subject,
		"userID": auditref.User, "user": auditref.User, "actor": auditref.User, "custodianUserId": auditref.User,
		"agent_ledger_user": auditref.User, "custodianUser": auditref.User, "mediatorUser": auditref.User,
		"party": auditref.Party, "resolvedParty": auditref.Party, "depositor_party_id": auditref.Party,
		"parties": auditref.Party, "observerParties": auditref.Party,
	} {
		if got, ok := personKind(key); !ok || got != want {
			t.Errorf("personKind(%q) = %q, %v; want %q", key, got, ok, want)
		}
	}
	for _, key := range []string{"username", "email", "id", "escrowId", "walletAddress", "institution", "userAgent"} {
		if got, ok := personKind(key); ok {
			t.Errorf("personKind(%q) = %q, want none", key, got)
		}
	}
}

// TestRefScrubCore (daml-escrow PLAN.md Phase 96 field sweep): outside
// development a field naming a person is its audit reference; a reference
// passes unchanged; everything else is scrubbed as before.
func TestRefScrubCore(t *testing.T) {
	l, logs, refs := refLogger(t, "production")
	l.With(zap.String("oktaSub", "00u1joey")).Info("acted",
		zap.String("userID", "u-joey-depositor-devlocal"),
		zap.Stringer("depositorParty", stringer("u-joey::1220abcdef")),
		zap.Strings("observerParties", []string{"a::1220aa", "b::1220bb"}),
		zap.String("actor", refs.User("u-joey-depositor-devlocal")),
		zap.String("email", "joey@depositor.devlocal"),
		zap.Int("userID2", 7))
	m := fieldMap(logs.All()[0])
	for key, want := range map[string]any{
		"oktaSub":        refs.Subject("00u1joey"),
		"userID":         refs.User("u-joey-depositor-devlocal"),
		"depositorParty": refs.Party("u-joey::1220abcdef"),
		"actor":          refs.User("u-joey-depositor-devlocal"),
		"email":          "[redacted]",
	} {
		if m[key] != want {
			t.Errorf("%s = %v, want %v", key, m[key], want)
		}
	}
	if !strings.HasPrefix(m["userID"].(string), "ref:r1:u:") {
		t.Errorf("userID = %v, want a user reference", m["userID"])
	}
	if ps, _ := m["observerParties"].([]any); len(ps) != 2 || ps[0] != refs.Party("a::1220aa") {
		t.Errorf("observerParties = %v", m["observerParties"])
	}
}

// TestRefScrubCoreDevelopment: development shows the identifier, still
// scrubbed -- an email subject stays masked.
func TestRefScrubCoreDevelopment(t *testing.T) {
	l, logs, _ := refLogger(t, "dev")
	l.Info("acted", zap.String("userID", "u-joey-depositor-devlocal"), zap.String("oktaSub", "joey@depositor.devlocal"))
	m := fieldMap(logs.All()[0])
	if m["userID"] != "u-joey-depositor-devlocal" || m["oktaSub"] != "[redacted]" {
		t.Errorf("fields = %v", m)
	}
}

// TestNewWithRefs: New wires Config.Refs into the core.
func TestNewWithRefs(t *testing.T) {
	refs, err := auditref.New(auditref.Config{Environment: "production", Keys: map[string][]byte{"r1": []byte("k")}, Current: "r1"})
	if err != nil {
		t.Fatal(err)
	}
	path := t.TempDir() + "/log.json"
	l, err := build(Config{Service: "escrow-api", Environment: "production", Refs: refs}, path)
	if err != nil {
		t.Fatal(err)
	}
	l.Info("acted", zap.String("userID", "u-joey"))
	_ = l.Sync()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if b := string(raw); strings.Contains(b, `"u-joey"`) || !strings.Contains(b, refs.User("u-joey")) {
		t.Errorf("log = %s", b)
	}
}
