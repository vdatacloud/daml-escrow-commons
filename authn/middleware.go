package authn

import (
	"net/http"
	"strings"

	"go.uber.org/zap"

	"github.com/vdatacloud/daml-escrow-commons/apierror"
	"github.com/vdatacloud/daml-escrow-commons/auditref"
	commonslog "github.com/vdatacloud/daml-escrow-commons/logging"
)

// Option configures Middleware.
type Option func(*options)

type options struct {
	policy    *ClientPolicy
	logger    *zap.Logger
	devBypass bool
	devScopes []string
}

// WithClientPolicy withholds guarded scopes from clients not allowed them.
func WithClientPolicy(p *ClientPolicy) Option { return func(o *options) { o.policy = p } }

// WithLogger logs withheld scopes (client id and scope only).
func WithLogger(l *zap.Logger) Option { return func(o *options) { o.logger = l } }

// WithDevBypass trusts X-Dev-User, X-Dev-OktaSub and X-Dev-Client headers
// instead of a token and grants scopes. Development only: the caller decides
// (from its environment) whether to pass it, and must never in production.
// The client policy applies to a bypass request only when it names a client
// (X-Dev-Client), so dev flows keep working and the guard can still be tried.
// A request that carries a JWT bearer token and no X-Dev-User is not a bypass
// request: when the middleware has a verifier, the token is verified and its
// own scopes apply, so a real sign-in is never silently widened to the dev
// scopes (daml-escrow PLAN.md Phase 104).
func WithDevBypass(scopes ...string) Option {
	return func(o *options) { o.devBypass, o.devScopes = true, scopes }
}

// Middleware authenticates each request and stores its Principal in the
// context; a missing or invalid bearer token is a 401 apierror envelope.
// Mount public routes (health, discovery) outside it.
func Middleware(verifier TokenVerifier, opts ...Option) func(http.Handler) http.Handler {
	o := &options{logger: zap.NewNop()}
	for _, opt := range opts {
		opt(o)
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var p Principal
			raw, hasBearer := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
			hasBearer = hasBearer && raw != ""
			if o.devBypass && (r.Header.Get("X-Dev-User") != "" || !hasBearer || verifier == nil || !looksLikeJWT(raw)) {
				p = devPrincipal(r, o.devScopes)
			} else {
				if !hasBearer {
					apierror.Write(w, r, apierror.New(http.StatusUnauthorized, apierror.CodeUnauthenticated, "sign-in required"))
					return
				}
				var err error
				if p, err = verifier.Verify(r.Context(), raw); err != nil {
					apierror.Write(w, r, apierror.New(http.StatusUnauthorized, apierror.CodeUnauthenticated, "sign-in required"))
					return
				}
			}
			if !p.DevBypass || p.ClientID != "" {
				var withheld []string
				p.Scopes, withheld = o.policy.Filter(p.ClientID, p.Scopes)
				for _, s := range withheld {
					o.logger.Warn("authn: scope withheld: client not allowed it",
						zap.String("client", p.ClientID), zap.String("scope", s), zap.String("path", r.URL.Path))
				}
			}
			next.ServeHTTP(w, r.WithContext(WithPrincipal(r.Context(), p)))
		})
	}
}

// looksLikeJWT: a signed JWT's three dot-separated parts. Dev callers' static
// tokens (dev-identity-token, bypass-token-<email>) aren't, and stay dev.
func looksLikeJWT(raw string) bool { return strings.Count(raw, ".") == 2 }

func devPrincipal(r *http.Request, scopes []string) Principal {
	email := r.Header.Get("X-Dev-User")
	if email == "" {
		email = "dev@daml-escrow.local"
	}
	sub := r.Header.Get("X-Dev-OktaSub")
	if sub == "" {
		sub = "dev-" + strings.ToLower(strings.NewReplacer("|", "-", "@", "-", ".", "-", "_", "-").Replace(email))
	}
	return Principal{Subject: sub, Email: email, Scopes: append([]string(nil), scopes...),
		ClientID: r.Header.Get("X-Dev-Client"), DevBypass: true}
}

// ActorMiddleware names the caller in the request's audit access log by the
// subject's audit reference (never a name or email). Mount after Middleware.
func ActorMiddleware(refs *auditref.Refs) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if p, ok := FromContext(r.Context()); ok && p.Subject != "" {
				commonslog.SetActor(r.Context(), refs.Subject(p.Subject))
			}
			next.ServeHTTP(w, r)
		})
	}
}

// RequireScope is route middleware that refuses (403 apierror) a request
// whose principal may not use scope.
func RequireScope(scope string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !HasScope(r.Context(), scope) {
				apierror.Write(w, r, apierror.Newf(http.StatusForbidden, apierror.CodeForbidden, "requires scope %s", scope))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
