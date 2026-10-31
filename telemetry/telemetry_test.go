package telemetry

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
)

// recording installs an in-memory tracer provider and the propagator for
// one test, restoring the previous globals after.
func recording(t *testing.T) *tracetest.SpanRecorder {
	t.Helper()
	prevTP, prevProp := otel.GetTracerProvider(), otel.GetTextMapPropagator()
	rec := tracetest.NewSpanRecorder()
	otel.SetTracerProvider(sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(rec)))
	otel.SetTextMapPropagator(Propagator())
	t.Cleanup(func() { otel.SetTracerProvider(prevTP); otel.SetTextMapPropagator(prevProp) })
	return rec
}

func attrs(s sdktrace.ReadOnlySpan) map[attribute.Key]attribute.Value {
	m := map[attribute.Key]attribute.Value{}
	for _, kv := range s.Attributes() {
		m[kv.Key] = kv.Value
	}
	return m
}

const (
	traceID     = "4bf92f3577b34da6a3ce929d0e0e4736"
	traceparent = "00-" + traceID + "-00f067aa0ba902b7-01"
	contractID  = "00b34d9a70b29cc6404ce837603ecdd1e084ce51b65c5b8e6f9c28ba47754a2f7d"
)

func TestMiddleware(t *testing.T) {
	rec := recording(t)
	h := Middleware("test", func(*http.Request) string { return "/escrows/{id}/fund" })(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !trace.SpanContextFromContext(r.Context()).IsValid() {
			t.Error("the handler's context carries the span")
		}
		w.WriteHeader(http.StatusCreated)
	}))
	req := httptest.NewRequest("POST", "/escrows/"+contractID+"/fund?email=joe@example.com", nil)
	req.Header.Set("traceparent", traceparent)
	h.ServeHTTP(httptest.NewRecorder(), req)

	spans := rec.Ended()
	if len(spans) != 1 {
		t.Fatalf("one span, got %d", len(spans))
	}
	s := spans[0]
	if s.Name() != "POST /escrows/{id}/fund" || s.SpanKind() != trace.SpanKindServer {
		t.Errorf("span %q kind %v", s.Name(), s.SpanKind())
	}
	if s.SpanContext().TraceID().String() != traceID {
		t.Errorf("the caller's trace continues: %s", s.SpanContext().TraceID())
	}
	a := attrs(s)
	if a[attrPath].AsString() != "/escrows/00b34d9a…2f7d/fund" || a[attrRoute].AsString() != "/escrows/{id}/fund" || a[attrStatus].AsInt64() != 201 {
		t.Errorf("attributes: %v", a)
	}
	for _, kv := range s.Attributes() {
		if strings.Contains(kv.Value.String(), "joe@example.com") || strings.Contains(kv.Value.String(), contractID) {
			t.Errorf("attribute %s leaks %q", kv.Key, kv.Value.String())
		}
	}
	if s.Status().Code == codes.Error {
		t.Error("a 201 isn't an error")
	}
}

func TestMiddlewareServerError(t *testing.T) {
	rec := recording(t)
	Middleware("test", nil)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
	})).ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/x", nil))
	if s := rec.Ended()[0]; s.Status().Code != codes.Error {
		t.Errorf("a 502 marks the span an error: %v", s.Status())
	}
}

// TestTransportPropagates: a client span, and the server it calls receives
// the trace -- so a server-side Middleware continues it.
func TestTransportPropagates(t *testing.T) {
	rec := recording(t)
	var got string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = r.Header.Get("traceparent")
		w.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()

	ctx, parent := otel.Tracer("test").Start(context.Background(), "parent")
	req, _ := http.NewRequestWithContext(ctx, "POST", server.URL+"/usage/"+contractID+"?q=joe@example.com", nil)
	resp, err := (&http.Client{Transport: Transport("test", nil)}).Do(req)
	if err != nil {
		t.Fatal(err)
	}
	_ = resp.Body.Close()
	parent.End()
	if req.Header.Get("traceparent") != "" {
		t.Error("the caller's request is not modified")
	}

	var client sdktrace.ReadOnlySpan
	for _, s := range rec.Ended() {
		if s.SpanKind() == trace.SpanKindClient {
			client = s
		}
	}
	if client == nil {
		t.Fatal("no client span")
	}
	if client.Parent().SpanID() != parent.SpanContext().SpanID() {
		t.Error("the client span is the caller's child")
	}
	if !strings.Contains(got, client.SpanContext().TraceID().String()) || !strings.Contains(got, client.SpanContext().SpanID().String()) {
		t.Errorf("traceparent %q doesn't carry the client span", got)
	}
	a := attrs(client)
	if a[attrPath].AsString() != "/usage/00b34d9a…2f7d" || a[attrStatus].AsInt64() != 202 || a[attrAddress].AsString() != "127.0.0.1" {
		t.Errorf("attributes: %v", a)
	}
}

type failing struct{}

func (failing) RoundTrip(*http.Request) (*http.Response, error) { return nil, errors.New("refused") }

func TestTransportError(t *testing.T) {
	rec := recording(t)
	req, _ := http.NewRequest("GET", "http://identity.internal/x", nil)
	if _, err := Transport("test", failing{}).RoundTrip(req); err == nil {
		t.Fatal("the error passes through")
	}
	if s := rec.Ended()[0]; s.Status().Code != codes.Error || len(s.Events()) == 0 {
		t.Errorf("a transport error is recorded: %v %v", s.Status(), s.Events())
	}
}

func TestSetup(t *testing.T) {
	prevTP, prevMP, prevProp := otel.GetTracerProvider(), otel.GetMeterProvider(), otel.GetTextMapPropagator()
	t.Cleanup(func() {
		otel.SetTracerProvider(prevTP)
		otel.SetMeterProvider(prevMP)
		otel.SetTextMapPropagator(prevProp)
	})

	shutdown, err := Setup(context.Background(), Config{Enabled: false})
	if err != nil || shutdown(context.Background()) != nil {
		t.Fatalf("disabled: %v", err)
	}
	if !strings.Contains(strings.Join(otel.GetTextMapPropagator().Fields(), ","), "traceparent") {
		t.Error("disabled still propagates trace context")
	}
	if _, ok := otel.GetTracerProvider().(*sdktrace.TracerProvider); ok {
		t.Error("disabled installs no SDK tracer provider")
	}

	for _, endpoint := range []string{"127.0.0.1:1", "http://127.0.0.1:1/", ""} {
		shutdown, err = Setup(context.Background(), Config{Enabled: true, Endpoint: endpoint, Insecure: true, Service: "test"})
		if err != nil {
			t.Fatalf("enabled (%q): %v", endpoint, err)
		}
		if _, ok := otel.GetTracerProvider().(*sdktrace.TracerProvider); !ok {
			t.Errorf("enabled (%q) installs the SDK tracer provider", endpoint)
		}
		_, span := otel.Tracer("test").Start(context.Background(), "probe")
		if ro, ok := span.(sdktrace.ReadOnlySpan); !ok || ro.Resource().String() == "" ||
			!strings.Contains(ro.Resource().String(), "service.name=test") {
			t.Errorf("enabled (%q): the resource names the service", endpoint)
		}
		span.End()
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		_ = shutdown(ctx) // nothing listens: a failed final export is fine here
		cancel()
	}
}

// TestSetupSampleRatio: a ratio keeps about that share of new traces, and a
// sampled parent's decision is kept whatever the ratio.
func TestSetupSampleRatio(t *testing.T) {
	prevTP, prevMP, prevProp := otel.GetTracerProvider(), otel.GetMeterProvider(), otel.GetTextMapPropagator()
	t.Cleanup(func() {
		otel.SetTracerProvider(prevTP)
		otel.SetMeterProvider(prevMP)
		otel.SetTextMapPropagator(prevProp)
	})
	shutdown, err := Setup(context.Background(), Config{Enabled: true, Endpoint: "127.0.0.1:1", Insecure: true, SampleRatio: 0.1, Service: "test"})
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = shutdown(ctx)
	}()
	sampled := 0
	for i := 0; i < 2000; i++ {
		_, span := otel.Tracer("test").Start(context.Background(), "root")
		if span.SpanContext().IsSampled() {
			sampled++
		}
		span.End()
	}
	if sampled < 100 || sampled > 320 {
		t.Errorf("ratio 0.1 kept %d of 2000 new traces", sampled)
	}
	tid, _ := trace.TraceIDFromHex(traceID)
	sid, _ := trace.SpanIDFromHex("00f067aa0ba902b7")
	parent := trace.ContextWithRemoteSpanContext(context.Background(),
		trace.NewSpanContext(trace.SpanContextConfig{TraceID: tid, SpanID: sid, TraceFlags: trace.FlagsSampled, Remote: true}))
	for i := 0; i < 50; i++ {
		_, span := otel.Tracer("test").Start(parent, "child")
		if !span.SpanContext().IsSampled() {
			t.Fatal("a sampled parent's trace is kept")
		}
		span.End()
	}
}
