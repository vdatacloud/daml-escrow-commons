package devoidc

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	"github.com/vdatacloud/daml-escrow-commons/authn"
)

// A minted token passes authn's verifier (discovery + JWKS + audience) and
// carries the client id the operator-scope policy checks (daml-escrow PLAN.md
// Phase 104).
func TestStartAndMintVerifyWithAuthn(t *testing.T) {
	is := Start(t, "daml-escrow")
	v, err := authn.NewOIDCVerifier(context.Background(), is.URL(), "daml-escrow")
	if err != nil {
		t.Fatal(err)
	}
	tok, err := is.Mint(Token{Subject: "00uOperator", Email: "op@platform.devlocal", Scopes: []string{"platform:admin"}, ClientID: "operator-console-dev"})
	if err != nil {
		t.Fatal(err)
	}
	p, err := v.Verify(context.Background(), tok)
	if err != nil {
		t.Fatal(err)
	}
	want := authn.Principal{Subject: "00uOperator", Email: "op@platform.devlocal", Scopes: []string{"platform:admin"}, ClientID: "operator-console-dev"}
	if !reflect.DeepEqual(p, want) {
		t.Fatalf("got %+v", p)
	}
	other, _ := is.Mint(Token{Subject: "x", Email: "x@y.z", Audience: "someone-else"})
	if _, err := v.Verify(context.Background(), other); err == nil {
		t.Fatal("a token for another audience verified")
	}
}

// Several browser apps (the escrow portal, the operator console) sign in to
// one issuer; each gets its own CORS answer, and nobody else does.
func TestSeveralOrigins(t *testing.T) {
	const console = "http://localhost:4322"
	is := Start(t, "daml-escrow", testPortal, console)
	for _, origin := range []string{testPortal, console} {
		q := authorizeQuery("openid email platform:admin", origin+"/callback")
		if resp, _ := http.Get(is.URL() + "/v1/authorize?" + q.Encode()); resp.StatusCode != http.StatusOK {
			t.Fatalf("%s: sign-in page %s", origin, resp.Status)
		}
		req, _ := http.NewRequest(http.MethodOptions, is.URL()+"/v1/token", nil)
		req.Header.Set("Origin", origin)
		resp, err := http.DefaultClient.Do(req)
		if err != nil || resp.Header.Get("Access-Control-Allow-Origin") != origin {
			t.Fatalf("%s: CORS %v %q", origin, err, resp.Header.Get("Access-Control-Allow-Origin"))
		}
	}
	req, _ := http.NewRequest(http.MethodOptions, is.URL()+"/v1/token", nil)
	req.Header.Set("Origin", "https://evil.example")
	if resp, _ := http.DefaultClient.Do(req); resp.Header.Get("Access-Control-Allow-Origin") != "" {
		t.Fatal("an unknown origin may read token responses")
	}
	q := authorizeQuery("openid", "https://evil.example/callback")
	if resp, _ := http.Get(is.URL() + "/v1/authorize?" + q.Encode()); resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("another site: %s", resp.Status)
	}
}

func TestNewDefaultsToThePortal(t *testing.T) {
	is, err := New("http://issuer.test", "daml-escrow")
	if err != nil || !reflect.DeepEqual(is.origins, []string{DefaultOrigin}) {
		t.Fatalf("%v %v", is.origins, err)
	}
}
