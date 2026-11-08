package authn

import (
	"context"
	"fmt"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
)

// TokenVerifier turns a raw bearer token into a Principal.
type TokenVerifier interface {
	Verify(ctx context.Context, rawToken string) (Principal, error)
}

// OIDCVerifier verifies access tokens from the shared Okta authorization
// server: signature against its JWKS, issuer, audience and expiry.
type OIDCVerifier struct {
	verifier *oidc.IDTokenVerifier
}

// NewOIDCVerifier discovers issuer's keys and checks tokens for audience.
func NewOIDCVerifier(ctx context.Context, issuer, audience string) (*OIDCVerifier, error) {
	provider, err := oidc.NewProvider(ctx, issuer)
	if err != nil {
		return nil, fmt.Errorf("authn: discover issuer: %w", err)
	}
	return &OIDCVerifier{verifier: provider.Verifier(&oidc.Config{ClientID: audience})}, nil
}

// claims are the claims read from an access token. Okta carries scopes as an
// `scp` array, other providers (and dev-oidc) as a `scope` string: both are
// read. Okta names the client `cid`; OIDC providers generally `azp`.
type claims struct {
	Email string   `json:"email"`
	Scope string   `json:"scope"`
	Scp   []string `json:"scp"`
	CID   string   `json:"cid"`
	AZP   string   `json:"azp"`
}

func (c claims) principal(subject string) Principal {
	client := c.CID
	if client == "" {
		client = c.AZP
	}
	return Principal{Subject: subject, Email: c.Email, Scopes: uniqueScopes(strings.Fields(c.Scope), c.Scp), ClientID: client}
}

func uniqueScopes(lists ...[]string) []string {
	seen := map[string]bool{}
	var out []string
	for _, l := range lists {
		for _, s := range l {
			if s != "" && !seen[s] {
				seen[s] = true
				out = append(out, s)
			}
		}
	}
	return out
}

// Verify checks rawToken and returns its principal (before any client policy).
func (v *OIDCVerifier) Verify(ctx context.Context, rawToken string) (Principal, error) {
	tok, err := v.verifier.Verify(ctx, rawToken)
	if err != nil {
		return Principal{}, fmt.Errorf("authn: verify token: %w", err)
	}
	var c claims
	if err := tok.Claims(&c); err != nil {
		return Principal{}, fmt.Errorf("authn: read claims: %w", err)
	}
	return c.principal(tok.Subject), nil
}

// NewOIDCVerifierWithKeys checks tokens against a fixed key set instead of
// discovering issuer's: for tests, here and in consumers (oidc.StaticKeySet).
func NewOIDCVerifierWithKeys(issuer, audience string, keys oidc.KeySet) *OIDCVerifier {
	return &OIDCVerifier{verifier: oidc.NewVerifier(issuer, keys, &oidc.Config{ClientID: audience})}
}
