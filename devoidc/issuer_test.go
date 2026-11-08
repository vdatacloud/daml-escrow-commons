package devoidc

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/coreos/go-oidc/v3/oidc"
)

// A dev-oidc token passes the same go-oidc discovery + JWKS + audience
// verification escrow-api uses for Okta (api.NewRealVerifier), and carries
// the claims its middleware reads.
func TestDevOIDC_TokenVerifiesWithGoOIDC(t *testing.T) {
	var is *Issuer
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { is.Handler().ServeHTTP(w, r) }))
	defer server.Close()
	var err error
	if is, err = New(server.URL, "daml-escrow"); err != nil {
		t.Fatal(err)
	}

	resp, err := http.PostForm(server.URL+"/token", url.Values{
		"sub": {"joey@depositor.devlocal"}, "email": {"joey@depositor.devlocal"},
		"scope": {"escrow:read escrow:write"}, "role": {"Depositor"},
	})
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("token: %v %v", err, resp.Status)
	}
	var tok struct {
		AccessToken string `json:"access_token"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&tok)

	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, server.URL)
	if err != nil {
		t.Fatalf("discovery: %v", err)
	}
	idt, err := provider.Verifier(&oidc.Config{ClientID: "daml-escrow"}).Verify(ctx, tok.AccessToken)
	if err != nil {
		t.Fatalf("verify: %v", err)
	}
	var claims struct {
		Email  string   `json:"email"`
		Scopes []string `json:"scp"`
		Role   string   `json:"role"`
	}
	_ = idt.Claims(&claims)
	if idt.Subject != "joey@depositor.devlocal" || claims.Email != "joey@depositor.devlocal" || claims.Role != "Depositor" || strings.Join(claims.Scopes, " ") != "escrow:read escrow:write" {
		t.Fatalf("claims: sub=%s %+v", idt.Subject, claims)
	}

	// Wrong audience is rejected, as escrow-api's verifier would.
	if _, err := provider.Verifier(&oidc.Config{ClientID: "someone-else"}).Verify(ctx, tok.AccessToken); err == nil {
		t.Fatal("a token for another audience must not verify")
	}
}

func TestDevOIDC_TokenRequiresSubAndEmail(t *testing.T) {
	is, _ := New("http://issuer.test", "daml-escrow")
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/token", strings.NewReader("sub=x"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	is.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("got %d, want 400 without email", rec.Code)
	}
}

// PLAN.md Phase 79.6: a person's token can be a signing session's -- the
// signing client, institution and signer roles -- and is a fresh sign-in.
func TestDevOIDC_SigningSessionClaims(t *testing.T) {
	is, _ := New("http://issuer.test", "daml-escrow")
	form := url.Values{"sub": {"s"}, "email": {"s@bank.test"}, "scope": {"institution:sign"}, "cid": {"signing"},
		"institution": {"bank"}, "signer_roles": {"release-signer agreement-signer"}}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/token", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	is.Handler().ServeHTTP(rec, req)
	var out struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&out); err != nil {
		t.Fatal(err)
	}
	payload, err := base64.RawURLEncoding.DecodeString(strings.Split(out.AccessToken, ".")[1])
	if err != nil {
		t.Fatal(err)
	}
	var c struct {
		CID         string   `json:"cid"`
		Institution string   `json:"institution"`
		SignerRoles []string `json:"signer_roles"`
		AuthTime    int64    `json:"auth_time"`
	}
	if err := json.Unmarshal(payload, &c); err != nil {
		t.Fatal(err)
	}
	if c.CID != "signing" || c.Institution != "bank" || len(c.SignerRoles) != 2 || c.AuthTime == 0 {
		t.Fatalf("claims: %+v", c)
	}
}

// PLAN.md Phase 70: a client-credentials token (escrow-api's service ledger
// user) has the client as its subject and the requested audience, and
// verifies against the same JWKS a participant would use.
func TestDevOIDC_ClientCredentials(t *testing.T) {
	var is *Issuer
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { is.Handler().ServeHTTP(w, r) }))
	defer server.Close()
	var err error
	if is, err = New(server.URL, "daml-escrow"); err != nil {
		t.Fatal(err)
	}
	req, _ := http.NewRequest(http.MethodPost, server.URL+"/token", strings.NewReader(url.Values{
		"grant_type": {"client_credentials"}, "audience": {"https://daml.com/jwt/aud/participant/p1"}, "scope": {"daml_ledger_api"},
	}.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth("escrow-svc", "any-secret")
	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("token: %v %v", err, resp.Status)
	}
	var tok struct {
		AccessToken string `json:"access_token"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&tok)

	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, server.URL)
	if err != nil {
		t.Fatal(err)
	}
	idt, err := provider.Verifier(&oidc.Config{ClientID: "https://daml.com/jwt/aud/participant/p1"}).Verify(ctx, tok.AccessToken)
	if err != nil || idt.Subject != "escrow-svc" {
		t.Fatalf("verify: %v sub=%v", err, idt)
	}

	missing, _ := http.PostForm(server.URL+"/token", url.Values{"grant_type": {"client_credentials"}})
	if missing.StatusCode != http.StatusBadRequest {
		t.Fatalf("no client: %d", missing.StatusCode)
	}
}
