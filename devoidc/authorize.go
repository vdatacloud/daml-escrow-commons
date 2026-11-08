package devoidc

// Browser sign-in (PLAN.md Phase 91.1): a browser app's Authorization Code +
// PKCE flow (the escrow portal's frontend/src/lib/auth.ts, the operator console), so a browser can sign in as any dev
// persona -- signing sessions included -- the way it signs in through Okta.
// DEV ONLY like the rest of this command: the page asks for no password.

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"html/template"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"sync"
	"time"
)

// signingScope marks a signing session's request (the portal's
// SIGNING_SCOPE); the page then asks for the institution and signer roles,
// which Okta's signing client reads from the person's profile.
const signingScope = "institution:sign"

const codeTTL = 2 * time.Minute

// authCode is one sign-in waiting for its token exchange.
type authCode struct {
	clientID, redirectURI, challenge string
	email, scope, acr                string
	institution                      string
	roles                            []string
	expires                          time.Time
}

type codeStore struct {
	mu    sync.Mutex
	codes map[string]authCode
}

func (s *codeStore) put(c authCode) (string, error) {
	b := make([]byte, 24)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	code := base64.RawURLEncoding.EncodeToString(b)
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.codes == nil {
		s.codes = map[string]authCode{}
	}
	s.codes[code] = c
	return code, nil
}

// take returns the code's sign-in once: a code is single-use.
func (s *codeStore) take(code string) (authCode, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	c, ok := s.codes[code]
	delete(s.codes, code)
	if !ok || time.Now().After(c.expires) {
		return authCode{}, false
	}
	return c, true
}

var authorizePage = template.Must(template.New("authorize").Parse(`<!doctype html>
<html><head><meta charset="utf-8"><title>dev-oidc sign-in</title>
<style>body{font-family:system-ui,sans-serif;max-width:28rem;margin:3rem auto;padding:0 1rem}
label{display:block;margin:.75rem 0 .25rem}input{width:100%;padding:.4rem;box-sizing:border-box}
button{margin-top:1rem;padding:.5rem 1rem}.note{color:#a33;font-size:.85rem}</style></head>
<body><h1>dev-oidc sign-in</h1>
<p class="note">Development only: signs in as whoever you name, with no password.</p>
<form method="post">
{{range $k, $v := .Hidden}}<input type="hidden" name="{{$k}}" value="{{$v}}">{{end}}
<label for="email">Email</label><input id="email" name="email" value="{{.Email}}" required autofocus>
{{if .Signing}}<p>Signing session for <code>{{.ClientID}}</code>.</p>
<label for="institution">Institution</label><input id="institution" name="institution" value="bank">
<label for="signer_roles">Signer roles (space-separated, may be empty)</label><input id="signer_roles" name="signer_roles">
{{end}}<button type="submit">Sign in</button></form></body></html>`))

// authorizeParams are the request's parameters the page carries through.
var authorizeParams = []string{"client_id", "redirect_uri", "state", "scope", "code_challenge", "code_challenge_method", "acr_values", "response_type"}

// authorize shows the sign-in page (GET) and, on submit (POST), redirects
// back to the portal with a code.
func (is *Issuer) authorize(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	f := r.Form
	if f.Get("response_type") != "code" || f.Get("client_id") == "" || f.Get("code_challenge") == "" ||
		f.Get("code_challenge_method") != "S256" {
		http.Error(w, "needs response_type=code, client_id and an S256 code_challenge", http.StatusBadRequest)
		return
	}
	// Only back to an allowed origin: this page would otherwise hand a code
	// to any site.
	if !is.allowedRedirect(f.Get("redirect_uri")) {
		http.Error(w, "redirect_uri must be on one of "+strings.Join(is.origins, ", "), http.StatusBadRequest)
		return
	}
	signing := strings.Contains(" "+f.Get("scope")+" ", " "+signingScope+" ")
	if r.Method == http.MethodGet {
		hidden := map[string]string{}
		for _, k := range authorizeParams {
			if v := f.Get(k); v != "" {
				hidden[k] = v
			}
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_ = authorizePage.Execute(w, map[string]any{
			"Hidden": hidden, "Email": f.Get("login_hint"), "Signing": signing, "ClientID": f.Get("client_id"),
		})
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "GET or POST", http.StatusMethodNotAllowed)
		return
	}
	email := strings.TrimSpace(f.Get("email"))
	if !strings.Contains(email, "@") {
		http.Error(w, "an email is required", http.StatusBadRequest)
		return
	}
	c := authCode{clientID: f.Get("client_id"), redirectURI: f.Get("redirect_uri"), challenge: f.Get("code_challenge"),
		email: email, scope: f.Get("scope"), acr: f.Get("acr_values"), expires: time.Now().Add(codeTTL)}
	if signing {
		c.institution, c.roles = strings.TrimSpace(f.Get("institution")), strings.Fields(f.Get("signer_roles"))
	}
	code, err := is.codes.put(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	back, _ := url.Parse(c.redirectURI)
	q := back.Query()
	q.Set("code", code)
	if s := f.Get("state"); s != "" {
		q.Set("state", s)
	}
	back.RawQuery = q.Encode()
	http.Redirect(w, r, back.String(), http.StatusFound)
}

// exchangeCode is the token endpoint's authorization-code grant: the code,
// its redirect_uri and client, and the PKCE verifier must all match.
func (is *Issuer) exchangeCode(w http.ResponseWriter, r *http.Request) {
	c, ok := is.codes.take(r.Form.Get("code"))
	sum := sha256.Sum256([]byte(r.Form.Get("code_verifier")))
	verified := base64.RawURLEncoding.EncodeToString(sum[:])
	if !ok || c.redirectURI != r.Form.Get("redirect_uri") || c.clientID != r.Form.Get("client_id") ||
		subtle.ConstantTimeCompare([]byte(verified), []byte(c.challenge)) != 1 {
		w.WriteHeader(http.StatusBadRequest)
		writeJSON(w, map[string]string{"error": "invalid_grant"})
		return
	}
	extra := map[string]string{"cid": c.clientID, "acr": c.acr, "institution": c.institution}
	ttl := 3600
	if c.institution != "" {
		ttl = 600 // a signing session's token, like Okta's signing client
	}
	tok, err := is.userToken(c.email, c.email, is.audience, strings.Fields(c.scope), extra, c.roles, ttl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]any{"access_token": tok, "token_type": "Bearer", "expires_in": ttl, "scope": c.scope})
}

// allowedRedirect: redirect is under one of the allowed origins.
func (is *Issuer) allowedRedirect(redirect string) bool {
	for _, o := range is.origins {
		if strings.HasPrefix(redirect, o+"/") {
			return true
		}
	}
	return false
}

// allowOrigin lets an allowed origin's page read the token response.
func (is *Issuer) allowOrigin(w http.ResponseWriter, r *http.Request) bool {
	if o := r.Header.Get("Origin"); o != "" && slices.Contains(is.origins, o) {
		w.Header().Set("Access-Control-Allow-Origin", o)
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Vary", "Origin")
	}
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return false
	}
	return true
}
