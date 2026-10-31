package logging

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/vdatacloud/daml-escrow-commons/apierror"
	"github.com/vdatacloud/daml-escrow-commons/redact"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

// observed is a scrubbed logger and what it wrote.
func observed() (*zap.Logger, *observer.ObservedLogs) {
	core, logs := observer.New(zapcore.DebugLevel)
	return zap.New(NewScrubCore(core)), logs
}

func fieldMap(e observer.LoggedEntry) map[string]any { return e.ContextMap() }

func TestNew(t *testing.T) {
	l, err := New(Config{Level: "debug", Service: "escrow-api", Environment: "dev", Version: "v1"})
	if err != nil || l == nil {
		t.Fatalf("New: %v", err)
	}
	if !l.Core().Enabled(zapcore.DebugLevel) {
		t.Error("level debug should enable debug")
	}
	if _, err := New(Config{}); err != nil {
		t.Errorf("the zero config is info and json: %v", err)
	}
	if _, err := New(Config{Level: "loud"}); err == nil {
		t.Error("an unknown level is refused")
	}
	if _, err := New(Config{Format: "xml"}); err == nil {
		t.Error("an unknown format is refused")
	}
}

// TestNewIsNotSampled: rule L1 -- a burst of identical lines is all
// written (zap's production default would keep 100 a second, then 1 in 100),
// as JSON with the service tags and a UTC timestamp.
func TestNewIsNotSampled(t *testing.T) {
	path := filepath.Join(t.TempDir(), "log.json")
	l, err := build(Config{Service: "escrow-api", Environment: "dev"}, path)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 500; i++ {
		l.Info("same line", zap.String("email", "joe@example.com"))
	}
	_ = l.Sync()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	lines := strings.Split(strings.TrimSpace(string(raw)), "\n")
	if len(lines) != 500 {
		t.Fatalf("wrote %d of 500 identical lines", len(lines))
	}
	var first map[string]any
	if err := json.Unmarshal([]byte(lines[0]), &first); err != nil {
		t.Fatalf("not JSON: %v", err)
	}
	if first["service"] != "escrow-api" || first["environment"] != "dev" || first["email"] != redact.Redacted {
		t.Errorf("line: %v", first)
	}
	if ts, _ := first["ts"].(string); !strings.HasSuffix(ts, "Z") {
		t.Errorf("ts %q is not UTC ISO-8601", first["ts"])
	}
}

func TestScrubCore(t *testing.T) {
	l, logs := observed()
	l.Info("payment for card 4111111111111111 failed",
		zap.String("email", "joe@example.com"),
		zap.String("display_name", "Joe Bloggs"),
		zap.String("routing_number", "021000021"),
		zap.String("account_number", "000123456789"),
		zap.Int64("accountNumber", 987654321),
		zap.String("walletAddress", "0x52908400098527886E0F7030069857D2E4169EE7"),
		zap.String("Authorization", "Bearer eyJhbGciOi"),
		zap.String("note", "pan 4111-1111-1111-1111 on file"),
		zap.Error(errors.New("charge 5500000000000004 declined")),
		zap.Any("physical_address", map[string]string{"street": "1 Main St"}),
		zap.Stringer("cardNumber", stringer("4111111111111111")),
		zap.String("escrowId", "e-1"),
		zap.Int("status", 402),
		zap.String("userID", "joey@depositor.devlocal"),
	)
	e := logs.All()[0]
	if strings.Contains(e.Message, "4111111111111111") || !strings.Contains(e.Message, "…1111") {
		t.Errorf("message: %q", e.Message)
	}
	m := fieldMap(e)
	want := map[string]any{
		"email": redact.Redacted, "display_name": redact.Redacted, "routing_number": "…0021",
		"account_number": "…6789", "accountNumber": "…4321", "walletAddress": "0x529084…9EE7",
		"Authorization": redact.Redacted, "note": "pan …1111 on file", "error": "charge …0004 declined",
		"physical_address": redact.Redacted, "cardNumber": "…1111", "escrowId": "e-1", "status": int64(402),
		"userID": redact.Redacted,
	}
	for k, v := range want {
		if fmt.Sprint(m[k]) != fmt.Sprint(v) {
			t.Errorf("%s = %v, want %v", k, m[k], v)
		}
	}
}

// TestScrubCoreWith: fields bound with With are scrubbed too.
func TestScrubCoreWith(t *testing.T) {
	l, logs := observed()
	l.With(zap.String("email", "joe@example.com")).Info("bound")
	if m := fieldMap(logs.All()[0]); m["email"] != redact.Redacted {
		t.Errorf("email = %v", m["email"])
	}
}

type stringer string

func (s stringer) String() string { return string(s) }

func TestActorSlot(t *testing.T) {
	ctx := WithActorSlot(context.Background())
	SetActor(ctx, "u-joey")
	if Actor(ctx) != "u-joey" {
		t.Errorf("actor = %q", Actor(ctx))
	}
	SetActor(context.Background(), "nobody") // no slot: ignored, no panic
	if Actor(context.Background()) != "" {
		t.Error("no slot, no actor")
	}
}

func spanCtx(ctx context.Context) context.Context {
	tid, _ := trace.TraceIDFromHex("4bf92f3577b34da6a3ce929d0e0e4736")
	sid, _ := trace.SpanIDFromHex("00f067aa0ba902b7")
	return trace.ContextWithSpanContext(ctx, trace.NewSpanContext(trace.SpanContextConfig{TraceID: tid, SpanID: sid, TraceFlags: trace.FlagsSampled}))
}

func TestFromContext(t *testing.T) {
	l, logs := observed()
	FromContext(spanCtx(context.Background()), l).Info("x")
	m := fieldMap(logs.All()[0])
	if m["trace_id"] != "4bf92f3577b34da6a3ce929d0e0e4736" || m["span_id"] != "00f067aa0ba902b7" {
		t.Errorf("trace fields: %v", m)
	}
	FromContext(context.Background(), nil).Info("a no-op logger, no panic")
}

func TestAccessLog(t *testing.T) {
	l, logs := observed()
	cid := "00b34d9a70b29cc6404ce837603ecdd1e084ce51b65c5b8e6f9c28ba47754a2f7d"
	var inner []observer.LoggedEntry
	h := apierror.RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(spanCtx(r.Context()))
		AccessLog(l, WithRoute(func(*http.Request) string { return "/escrows/{id}/fund" }))(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			SetActor(r.Context(), "u-joey") // as the auth middleware would
			FromContext(r.Context(), nil).Info("handler line")
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte("ok"))
		})).ServeHTTP(w, r)
	}))
	req := httptest.NewRequest("POST", "/escrows/"+cid+"/fund?email=joe@example.com", nil)
	req.Header.Set(apierror.HeaderRequestID, "req-12345678")
	h.ServeHTTP(httptest.NewRecorder(), req)

	inner = logs.FilterMessage("handler line").All()
	if len(inner) != 1 || fieldMap(inner[0])["request_id"] != "req-12345678" || fieldMap(inner[0])["trace_id"] == nil {
		t.Fatalf("the handler's logger carries the request and trace ids: %v", inner)
	}
	access := logs.FilterMessage("http_request").All()
	if len(access) != 1 {
		t.Fatalf("one access line, got %d", len(access))
	}
	m := fieldMap(access[0])
	for k, v := range map[string]any{
		"method": "POST", "route": "/escrows/{id}/fund", "path": "/escrows/00b34d9a…2f7d/fund", "status": int64(201),
		"bytes": int64(2), "actor": "u-joey", "request_id": "req-12345678", "trace_id": "4bf92f3577b34da6a3ce929d0e0e4736",
	} {
		if m[k] != v {
			t.Errorf("%s = %v, want %v", k, m[k], v)
		}
	}
	if access[0].Level != zapcore.InfoLevel {
		t.Errorf("a 201 is logged at info, got %v", access[0].Level)
	}
	for _, e := range logs.All() {
		for _, v := range e.ContextMap() {
			if strings.Contains(fmt.Sprint(v), "joe@example.com") {
				t.Errorf("the query's email reached the log: %v", e.ContextMap())
			}
		}
	}
}

func TestAccessLogServerError(t *testing.T) {
	l, logs := observed()
	AccessLog(l)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusBadGateway)
	})).ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/x", nil))
	e := logs.FilterMessage("http_request").All()[0]
	if e.Level != zapcore.ErrorLevel || fieldMap(e)["status"] != int64(502) {
		t.Errorf("a 502 is logged at error: %v %v", e.Level, fieldMap(e))
	}
}
