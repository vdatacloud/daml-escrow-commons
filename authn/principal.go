// Package authn verifies a platform caller's bearer token and carries who
// they are through the request: the subject, email, scopes and the OAuth
// client (`cid`) the token was issued to.
//
// Every service on the platform trusts the same Okta authorization server and
// audience, so each had its own near-identical verifier and middleware
// (daml-escrow-identity, daml-escrow-platform, daml-escrow's escrow-api). This
// package is the one copy, plus the client policy (policy.go) that keeps
// operator scopes to the operator clients (daml-escrow PLAN.md Phase 104).
package authn

import "context"

// Principal is the authenticated caller.
type Principal struct {
	// Subject is the token's `sub` (the Okta user id for a person, the
	// client id for a service). Log it only through an audit reference.
	Subject string
	// Email is the `email` claim; empty for service tokens. Never log it.
	Email string
	// Scopes are the scopes this request may use: the token's scopes minus
	// any the client policy withheld.
	Scopes []string
	// ClientID is the OAuth client the token was issued to (`cid`, or `azp`).
	// Client ids name applications, not people, and log in the clear.
	ClientID string
	// DevBypass marks a principal built from X-Dev-* headers (dev only).
	DevBypass bool
}

// HasScope reports whether the principal may use scope.
func (p Principal) HasScope(scope string) bool {
	for _, s := range p.Scopes {
		if s == scope {
			return true
		}
	}
	return false
}

type ctxKey struct{}

// WithPrincipal returns ctx carrying p.
func WithPrincipal(ctx context.Context, p Principal) context.Context {
	return context.WithValue(ctx, ctxKey{}, p)
}

// FromContext returns the request's principal, if authenticated.
func FromContext(ctx context.Context) (Principal, bool) {
	p, ok := ctx.Value(ctxKey{}).(Principal)
	return p, ok
}

// HasScope reports whether the request's principal may use scope; false when
// there is no principal.
func HasScope(ctx context.Context, scope string) bool {
	p, ok := FromContext(ctx)
	return ok && p.HasScope(scope)
}
