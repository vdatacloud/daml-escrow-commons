package devoidc

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Issuer is a development OpenID Connect issuer: a fresh RSA key, the
// discovery document and JWKS, a token endpoint that mints a token for any
// identity, and the browser Authorization Code + PKCE flow.
type Issuer struct {
	url      string
	key      *rsa.PrivateKey
	kid      string
	audience string
	origins  []string // the browser origins sign-ins may return to
	codes    codeStore
}

// DefaultOrigin is the escrow portal's dev origin, allowed when New is given
// no origins.
const DefaultOrigin = "http://localhost:4321"

// New creates an issuer at url (its iss) with a fresh RSA signing key,
// minting tokens for audience by default. Browser sign-ins may return only to
// origins (DefaultOrigin when none are given).
func New(url, audience string, origins ...string) (*Issuer, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("devoidc: generate key: %w", err)
	}
	sum := sha256.Sum256(key.N.Bytes())
	if len(origins) == 0 {
		origins = []string{DefaultOrigin}
	}
	return &Issuer{url: url, key: key, kid: base64.RawURLEncoding.EncodeToString(sum[:8]), audience: audience,
		origins: append([]string(nil), origins...)}, nil
}

// URL is the issuer URL (iss).
func (is *Issuer) URL() string { return is.url }

// Handler serves discovery, JWKS, the token endpoint (/token, and /v1/token
// for browsers) and the sign-in page (/v1/authorize).
func (is *Issuer) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/.well-known/openid-configuration", is.discovery)
	mux.HandleFunc("/jwks", is.jwks)
	mux.HandleFunc("/token", is.token)
	mux.HandleFunc("/v1/authorize", is.authorize)
	mux.HandleFunc("/v1/token", func(w http.ResponseWriter, r *http.Request) {
		if is.allowOrigin(w, r) {
			is.token(w, r)
		}
	})
	return mux
}

func (is *Issuer) discovery(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, map[string]any{
		"issuer":                                is.url,
		"jwks_uri":                              is.url + "/jwks",
		"token_endpoint":                        is.url + "/token",
		"authorization_endpoint":                is.url + "/v1/authorize",
		"response_types_supported":              []string{"code"},
		"code_challenge_methods_supported":      []string{"S256"},
		"subject_types_supported":               []string{"public"},
		"id_token_signing_alg_values_supported": []string{"RS256"},
	})
}

func (is *Issuer) jwks(w http.ResponseWriter, r *http.Request) {
	pub := is.key.PublicKey
	writeJSON(w, map[string]any{"keys": []map[string]string{{
		"kty": "RSA", "use": "sig", "alg": "RS256", "kid": is.kid,
		"n": base64.RawURLEncoding.EncodeToString(pub.N.Bytes()),
		"e": base64.RawURLEncoding.EncodeToString(big.NewInt(int64(pub.E)).Bytes()),
	}}})
}

// token mints a JWT from form fields: sub and email (required), scope
// (space-separated -> scp), role, origin_domain, aud, ttl (seconds).
func (is *Issuer) token(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}
	// OAuth2 client credentials (PLAN.md Phase 70: escrow-api's service
	// ledger user): the client is the subject, no email; the "audience"
	// parameter (Auth0/Okta style) targets the participant. Any secret is
	// accepted -- dev only.
	if r.Form.Get("grant_type") == "client_credentials" {
		is.clientCredentials(w, r)
		return
	}
	if r.Form.Get("grant_type") == "authorization_code" {
		is.exchangeCode(w, r)
		return
	}
	sub, email := r.Form.Get("sub"), r.Form.Get("email")
	if sub == "" || email == "" {
		http.Error(w, "sub and email are required", http.StatusBadRequest)
		return
	}
	ttl := 3600
	if v, err := strconv.Atoi(r.Form.Get("ttl")); err == nil && v > 0 {
		ttl = v
	}
	aud := r.Form.Get("aud")
	if aud == "" {
		aud = is.audience
	}
	extra := map[string]string{}
	for _, k := range []string{"role", "origin_domain", "cid", "acr", "institution"} {
		extra[k] = r.Form.Get(k)
	}
	tok, err := is.userToken(sub, email, aud, strings.Fields(r.Form.Get("scope")), extra, strings.Fields(r.Form.Get("signer_roles")), ttl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]any{"access_token": tok, "token_type": "Bearer", "expires_in": ttl})
}

// userToken mints a person's token: scopes as scp; and the optional claims
// -- role, origin_domain, and a signing session's (PLAN.md Phase 79.5c, what
// Okta's signing client puts in its tokens): the client (cid), the sign-in's
// acr, the institution and the signer roles. Empty values are left out.
func (is *Issuer) userToken(sub, email, aud string, scopes []string, extra map[string]string, roles []string, ttl int) (string, error) {
	now := time.Now()
	claims := map[string]any{
		"iss": is.url, "sub": sub, "aud": aud, "email": email,
		"iat": now.Unix(), "exp": now.Add(time.Duration(ttl) * time.Second).Unix(),
		"scp": scopes, "auth_method": "oidc",
		// Every token is a fresh sign-in (PLAN.md Phase 79.4/79.5c).
		"auth_time": now.Unix(),
	}
	for k, v := range extra {
		if v != "" {
			claims[k] = v
		}
	}
	if len(roles) > 0 {
		claims["signer_roles"] = roles
	}
	return is.sign(claims)
}

func (is *Issuer) clientCredentials(w http.ResponseWriter, r *http.Request) {
	client := r.Form.Get("client_id")
	if id, _, ok := r.BasicAuth(); ok {
		client = id
	}
	if client == "" {
		http.Error(w, "client_id is required", http.StatusBadRequest)
		return
	}
	aud := r.Form.Get("audience")
	if aud == "" {
		aud = is.audience
	}
	ttl := 3600
	now := time.Now()
	claims := map[string]any{
		"iss": is.url, "sub": client, "aud": aud,
		"iat": now.Unix(), "exp": now.Add(time.Duration(ttl) * time.Second).Unix(),
		"scope": r.Form.Get("scope"),
		// scp, as Okta's access tokens carry a client's scopes (escrow-api
		// reads scp; PLAN.md Phase 103: the platform scheduler's escrow:jobs).
		"scp": strings.Fields(r.Form.Get("scope")),
	}
	tok, err := is.sign(claims)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, map[string]any{"access_token": tok, "token_type": "Bearer", "expires_in": ttl})
}

func (is *Issuer) sign(claims map[string]any) (string, error) {
	enc := base64.RawURLEncoding.EncodeToString
	h, _ := json.Marshal(map[string]string{"alg": "RS256", "typ": "JWT", "kid": is.kid})
	c, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	msg := enc(h) + "." + enc(c)
	digest := sha256.Sum256([]byte(msg))
	sig, err := rsa.SignPKCS1v15(rand.Reader, is.key, crypto.SHA256, digest[:])
	if err != nil {
		return "", fmt.Errorf("devoidc: sign: %w", err)
	}
	return msg + "." + enc(sig), nil
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}

// httpHandlerFunc serves the issuer get returns, once it exists (a test
// server's URL is known only after it starts).
func httpHandlerFunc(get func() *Issuer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) { get().Handler().ServeHTTP(w, r) }
}
