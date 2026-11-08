package devoidc

import (
	"net/http/httptest"
	"testing"
	"time"
)

// Token is what Mint puts in a person's token. Empty fields are left out.
type Token struct {
	Subject, Email string
	// Audience defaults to the issuer's.
	Audience string
	Scopes   []string
	// ClientID is the OAuth client (cid): the operator console, the escrow
	// portal, the signing client.
	ClientID                string
	Role, OriginDomain, ACR string
	Institution             string
	SignerRoles             []string
	TTL                     time.Duration // default one hour
}

// Mint signs t as a person's access token, as the token endpoint would.
func (is *Issuer) Mint(t Token) (string, error) {
	aud := t.Audience
	if aud == "" {
		aud = is.audience
	}
	ttl := int(t.TTL / time.Second)
	if ttl <= 0 {
		ttl = 3600
	}
	extra := map[string]string{"role": t.Role, "origin_domain": t.OriginDomain, "cid": t.ClientID, "acr": t.ACR, "institution": t.Institution}
	return is.userToken(t.Subject, t.Email, aud, t.Scopes, extra, t.SignerRoles, ttl)
}

// Start runs an issuer for audience on a local test server (no network
// beyond loopback) and closes it when the test ends: point a verifier at
// URL(), and Mint tokens for each persona.
func Start(tb testing.TB, audience string, origins ...string) *Issuer {
	tb.Helper()
	var is *Issuer
	srv := httptest.NewServer(httpHandlerFunc(func() *Issuer { return is }))
	tb.Cleanup(srv.Close)
	var err error
	if is, err = New(srv.URL, audience, origins...); err != nil {
		tb.Fatal(err)
	}
	return is
}
