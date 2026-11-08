package devoidc

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/coreos/go-oidc/v3/oidc"
)

const testPortal = "http://localhost:4321"

func browserIssuer(t *testing.T) (*Issuer, *httptest.Server) {
	t.Helper()
	var is *Issuer
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { is.Handler().ServeHTTP(w, r) }))
	t.Cleanup(server.Close)
	var err error
	if is, err = New(server.URL, "daml-escrow", testPortal); err != nil {
		t.Fatal(err)
	}
	return is, server
}

func challengeOf(verifier string) string {
	sum := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(sum[:])
}

func authorizeQuery(scope, redirect string) url.Values {
	return url.Values{"response_type": {"code"}, "client_id": {"escrow-signing-dev"}, "redirect_uri": {redirect},
		"state": {"st-1"}, "scope": {scope}, "code_challenge": {challengeOf("verifier-1")}, "code_challenge_method": {"S256"}}
}

var noRedirect = &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}

// signIn submits the sign-in page and returns the code it redirects back with.
func signIn(t *testing.T, server *httptest.Server, q url.Values, form url.Values) string {
	t.Helper()
	for k, v := range q {
		form[k] = v
	}
	resp, err := noRedirect.PostForm(server.URL+"/v1/authorize", form)
	if err != nil || resp.StatusCode != http.StatusFound {
		t.Fatalf("sign-in: %v %v", err, resp.Status)
	}
	back, _ := url.Parse(resp.Header.Get("Location"))
	if !strings.HasPrefix(back.String(), testPortal+"/callback") || back.Query().Get("state") != "st-1" {
		t.Fatalf("redirect %s", back)
	}
	return back.Query().Get("code")
}

func exchange(server *httptest.Server, code, verifier string) (*http.Response, error) {
	req, _ := http.NewRequest(http.MethodPost, server.URL+"/v1/token", strings.NewReader(url.Values{
		"grant_type": {"authorization_code"}, "code": {code}, "redirect_uri": {testPortal + "/callback"},
		"client_id": {"escrow-signing-dev"}, "code_verifier": {verifier}}.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Origin", testPortal)
	return http.DefaultClient.Do(req)
}

// PLAN.md Phase 91.1: a browser signing session through the portal's
// Authorization Code + PKCE flow yields a token escrow-api's verifier
// accepts, carrying the signing session's claims.
func TestDevOIDC_BrowserSigningSession(t *testing.T) {
	_, server := browserIssuer(t)
	q := authorizeQuery("openid email escrow:read institution:sign", testPortal+"/callback")

	page, err := http.Get(server.URL + "/v1/authorize?" + q.Encode() + "&login_hint=admin%40bank.devlocal")
	if err != nil || page.StatusCode != http.StatusOK {
		t.Fatalf("page: %v %v", err, page.Status)
	}
	body, _ := io.ReadAll(page.Body)
	if !strings.Contains(string(body), `value="admin@bank.devlocal"`) || !strings.Contains(string(body), `name="signer_roles"`) {
		t.Fatalf("the signing page lacks the email or the signer fields:\n%s", body)
	}

	code := signIn(t, server, q, url.Values{"email": {"admin@bank.devlocal"}, "institution": {"bank"}, "signer_roles": {"Approver Releaser"}})
	if resp, _ := exchange(server, code, "wrong-verifier"); resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("a wrong PKCE verifier was accepted: %s", resp.Status)
	}
	// The failed attempt used the code up.
	if resp, _ := exchange(server, code, "verifier-1"); resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("a used code was accepted: %s", resp.Status)
	}

	code = signIn(t, server, q, url.Values{"email": {"admin@bank.devlocal"}, "institution": {"bank"}, "signer_roles": {"Approver Releaser"}})
	resp, err := exchange(server, code, "verifier-1")
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("exchange: %v %v", err, resp.Status)
	}
	if resp.Header.Get("Access-Control-Allow-Origin") != testPortal {
		t.Fatalf("the portal can't read the token response: %q", resp.Header.Get("Access-Control-Allow-Origin"))
	}
	var tok struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&tok)
	if tok.ExpiresIn != 600 {
		t.Fatalf("a signing session's token lasts %d s, want 600", tok.ExpiresIn)
	}

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
		Email       string   `json:"email"`
		Scopes      []string `json:"scp"`
		ClientID    string   `json:"cid"`
		Institution string   `json:"institution"`
		Roles       []string `json:"signer_roles"`
		AuthTime    int64    `json:"auth_time"`
	}
	_ = idt.Claims(&claims)
	if claims.Email != "admin@bank.devlocal" || claims.ClientID != "escrow-signing-dev" || claims.Institution != "bank" ||
		strings.Join(claims.Roles, " ") != "Approver Releaser" || claims.AuthTime == 0 ||
		strings.Join(claims.Scopes, " ") != "openid email escrow:read institution:sign" {
		t.Fatalf("claims %+v", claims)
	}
}

// A portal session (no institution:sign) asks only for the email and
// carries no signing claims.
func TestDevOIDC_BrowserPortalSession(t *testing.T) {
	_, server := browserIssuer(t)
	q := authorizeQuery("openid email escrow:read", testPortal+"/callback")
	page, _ := http.Get(server.URL + "/v1/authorize?" + q.Encode())
	body, _ := io.ReadAll(page.Body)
	if strings.Contains(string(body), `name="signer_roles"`) {
		t.Fatal("a portal sign-in asked for signer roles")
	}
	code := signIn(t, server, q, url.Values{"email": {"joey@depositor.devlocal"}, "institution": {"bank"}})
	resp, _ := exchange(server, code, "verifier-1")
	var tok struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&tok)
	parts := strings.Split(tok.AccessToken, ".")
	if len(parts) != 3 || tok.ExpiresIn != 3600 {
		t.Fatalf("token %+v", tok)
	}
	raw, _ := base64.RawURLEncoding.DecodeString(parts[1])
	if strings.Contains(string(raw), `"institution"`) {
		t.Fatalf("a portal token carries an institution: %s", raw)
	}
}

// The page hands codes only back to the portal, and only for PKCE requests.
func TestDevOIDC_AuthorizeRefusals(t *testing.T) {
	_, server := browserIssuer(t)
	for name, q := range map[string]url.Values{
		"another site": authorizeQuery("openid", "https://evil.example/callback"),
		"a lookalike":  authorizeQuery("openid", testPortal+".evil.example/callback"),
		"no PKCE": func() url.Values {
			v := authorizeQuery("openid", testPortal+"/callback")
			v.Del("code_challenge")
			return v
		}(),
		"plain PKCE": func() url.Values {
			v := authorizeQuery("openid", testPortal+"/callback")
			v.Set("code_challenge_method", "plain")
			return v
		}(),
		"implicit response": func() url.Values {
			v := authorizeQuery("openid", testPortal+"/callback")
			v.Set("response_type", "token")
			return v
		}(),
	} {
		if resp, _ := http.Get(server.URL + "/v1/authorize?" + q.Encode()); resp.StatusCode != http.StatusBadRequest {
			t.Errorf("%s: %s, want 400", name, resp.Status)
		}
	}
}
