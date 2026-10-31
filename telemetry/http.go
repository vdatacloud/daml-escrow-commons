package telemetry

import (
	"net/http"

	"github.com/vdatacloud/daml-escrow-commons/redact"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// Span attributes, OpenTelemetry semantic-convention names.
const (
	attrMethod  = "http.request.method"
	attrRoute   = "http.route"
	attrPath    = "url.path" // always redact.Path'd: ids short, emails and the query gone
	attrStatus  = "http.response.status_code"
	attrAddress = "server.address"
)

// Middleware traces each request as a server span: it continues the
// caller's trace (W3C traceparent), names the span "METHOD route" once
// routing is done (route reports the pattern, e.g. chi's
// RouteContext(r.Context()).RoutePattern(); nil uses r.Pattern), and
// records the method, route, redacted path and status. A 5xx marks the span
// an error. Mount it after apierror.RequestID and before logging.AccessLog,
// so the access log carries the span's ids.
func Middleware(tracerName string, route func(*http.Request) string) func(http.Handler) http.Handler {
	if route == nil {
		route = func(r *http.Request) string { return r.Pattern }
	}
	tracer := otel.Tracer(tracerName)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))
			ctx, span := tracer.Start(ctx, "HTTP "+r.Method, trace.WithSpanKind(trace.SpanKindServer),
				trace.WithAttributes(attribute.String(attrMethod, r.Method), attribute.String(attrPath, redact.Path(r.URL.Path))))
			defer span.End()
			rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
			r = r.WithContext(ctx)
			next.ServeHTTP(rec, r)
			if p := route(r); p != "" {
				span.SetName(r.Method + " " + p)
				span.SetAttributes(attribute.String(attrRoute, p))
			}
			span.SetAttributes(attribute.Int(attrStatus, rec.status))
			if rec.status >= 500 {
				span.SetStatus(codes.Error, http.StatusText(rec.status))
			}
		})
	}
}

// Transport traces each outbound request as a client span and propagates
// the trace to the server it calls (W3C traceparent). It records the
// method, host, redacted path and status; a transport error or a 5xx marks
// the span an error. base nil is http.DefaultTransport.
func Transport(tracerName string, base http.RoundTripper) http.RoundTripper {
	if base == nil {
		base = http.DefaultTransport
	}
	return &transport{base: base, tracer: otel.Tracer(tracerName)}
}

type transport struct {
	base   http.RoundTripper
	tracer trace.Tracer
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx, span := t.tracer.Start(req.Context(), "HTTP "+req.Method, trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(attribute.String(attrMethod, req.Method), attribute.String(attrAddress, req.URL.Hostname()),
			attribute.String(attrPath, redact.Path(req.URL.Path))))
	defer span.End()
	out := req.Clone(ctx)
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(out.Header))
	resp, err := t.base.RoundTrip(out)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "transport error")
		return resp, err
	}
	span.SetAttributes(attribute.Int(attrStatus, resp.StatusCode))
	if resp.StatusCode >= 500 {
		span.SetStatus(codes.Error, http.StatusText(resp.StatusCode))
	}
	return resp, nil
}

// statusRecorder captures a response's status.
type statusRecorder struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func (r *statusRecorder) WriteHeader(code int) {
	if !r.wroteHeader {
		r.status, r.wroteHeader = code, true
	}
	r.ResponseWriter.WriteHeader(code)
}

func (r *statusRecorder) Write(b []byte) (int, error) {
	r.wroteHeader = true
	return r.ResponseWriter.Write(b)
}

func (r *statusRecorder) Unwrap() http.ResponseWriter { return r.ResponseWriter }

func (r *statusRecorder) Flush() {
	if f, ok := r.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}
