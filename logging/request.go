package logging

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/vdatacloud/daml-escrow-commons/apierror"
	"github.com/vdatacloud/daml-escrow-commons/redact"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type loggerKey struct{}

// WithLogger is ctx carrying l, for FromContext.
func WithLogger(ctx context.Context, l *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, l)
}

// FromContext is the request's logger -- the one AccessLog put in ctx, else
// fallback -- with ctx's trace and span ids added when ctx carries a span.
// A nil fallback with no logger in ctx is a no-op logger.
func FromContext(ctx context.Context, fallback *zap.Logger) *zap.Logger {
	l, _ := ctx.Value(loggerKey{}).(*zap.Logger)
	if l == nil {
		l = fallback
	}
	if l == nil {
		l = zap.NewNop()
	}
	return l.With(TraceFields(ctx)...)
}

// TraceFields are ctx's trace_id and span_id, or none without a span.
func TraceFields(ctx context.Context) []zap.Field {
	sc := trace.SpanContextFromContext(ctx)
	if !sc.IsValid() {
		return nil
	}
	return []zap.Field{zap.String("trace_id", sc.TraceID().String()), zap.String("span_id", sc.SpanID().String())}
}

// The actor slot: the access log runs before authentication (it must time
// and record the whole request, refusals included), so it puts an empty slot
// in the request's context and the auth middleware fills it once it knows
// who is acting. The actor is an identity -- a ledger user id, or a
// pseudoref reference -- never a name, an email or a raw OIDC subject (L2).
type actorSlot struct {
	mu sync.Mutex
	id string
}

type actorKey struct{}

// WithActorSlot is ctx with an empty actor slot (AccessLog adds one).
func WithActorSlot(ctx context.Context) context.Context {
	return context.WithValue(ctx, actorKey{}, &actorSlot{})
}

// SetActor records who is acting on the request ctx belongs to; a ctx
// without a slot ignores it.
func SetActor(ctx context.Context, id string) {
	if s, ok := ctx.Value(actorKey{}).(*actorSlot); ok {
		s.mu.Lock()
		s.id = id
		s.mu.Unlock()
	}
}

// Actor is the actor SetActor recorded, or "".
func Actor(ctx context.Context) string {
	if s, ok := ctx.Value(actorKey{}).(*actorSlot); ok {
		s.mu.Lock()
		defer s.mu.Unlock()
		return s.id
	}
	return ""
}

// AccessLogOption configures AccessLog.
type AccessLogOption func(*accessLog)

type accessLog struct {
	route func(*http.Request) string
}

// WithRoute names a request's route pattern after it is served (e.g. chi's
// RouteContext(r.Context()).RoutePattern()); the default is the standard
// library's r.Pattern.
func WithRoute(f func(*http.Request) string) AccessLogOption {
	return func(a *accessLog) { a.route = f }
}

// AccessLog writes one audit line per request, after it is served: method,
// route pattern, path (redact.Path: ids shortened, emails and the query
// dropped), status, bytes, duration, request id, actor, and trace and span
// ids. A 5xx is logged at error, anything else at info. Before serving, it
// puts the request's logger (with the request id) and an empty actor slot
// in the context. Mount it after apierror.RequestID and the tracing
// middleware, so both ids are known.
func AccessLog(logger *zap.Logger, opts ...AccessLogOption) func(http.Handler) http.Handler {
	a := &accessLog{route: func(r *http.Request) string { return r.Pattern }}
	for _, o := range opts {
		o(a)
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ctx := r.Context()
			reqLogger := logger
			if id := apierror.RequestIDFrom(ctx); id != "" {
				reqLogger = logger.With(zap.String("request_id", id))
			}
			ctx = WithActorSlot(WithLogger(ctx, reqLogger))
			rec := &recorder{ResponseWriter: w, status: http.StatusOK}
			r = r.WithContext(ctx)
			next.ServeHTTP(rec, r)

			fields := append([]zap.Field{
				zap.String("method", r.Method),
				zap.String("route", a.route(r)),
				zap.String("path", redact.Path(r.URL.Path)),
				zap.Int("status", rec.status),
				zap.Int64("bytes", rec.bytes),
				zap.Float64("duration_ms", float64(time.Since(start).Microseconds())/1000),
				zap.String("actor", Actor(ctx)),
			}, TraceFields(ctx)...)
			if rec.status >= 500 {
				reqLogger.Error("http_request", fields...)
				return
			}
			reqLogger.Info("http_request", fields...)
		})
	}
}

// recorder captures the status and size of a response.
type recorder struct {
	http.ResponseWriter
	status      int
	bytes       int64
	wroteHeader bool
}

func (r *recorder) WriteHeader(code int) {
	if !r.wroteHeader {
		r.status, r.wroteHeader = code, true
	}
	r.ResponseWriter.WriteHeader(code)
}

func (r *recorder) Write(b []byte) (int, error) {
	r.wroteHeader = true
	n, err := r.ResponseWriter.Write(b)
	r.bytes += int64(n)
	return n, err
}

// Unwrap lets http.ResponseController reach the underlying writer
// (Flush, Hijack, deadlines).
func (r *recorder) Unwrap() http.ResponseWriter { return r.ResponseWriter }

// Flush passes a flush through for streaming handlers.
func (r *recorder) Flush() {
	if f, ok := r.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}
