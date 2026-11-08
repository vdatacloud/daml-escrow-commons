package authn

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	jose "github.com/go-jose/go-jose/v4"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"

	"github.com/vdatacloud/daml-escrow-commons/auditref"
	commonslog "github.com/vdatacloud/daml-escrow-commons/logging"
)

const (
	issuer   = "https://issuer.example/oauth2/default"
	audience = "daml-escrow"
	console  = "0oaOperatorConsole"
	portal   = "0oaEscrowPortal"
)

func signer(t *testing.T) (jose.Signer, crypto.PublicKey) {
	t.Helper()
	k, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	s, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.RS256, Key: k}, nil)
	if err != nil {
		t.Fatal(err)
	}
	return s, k.Public()
}

func token(t *testing.T, s jose.Signer, claims map[string]any) string {
	t.Helper()
	base := map[string]any{"iss": issuer, "aud": audience, "sub": "00uPerson", "exp": time.Now().Add(time.Hour).Unix(), "iat": time.Now().Unix()}
	for k, v := range claims {
		base[k] = v
	}
	b, _ := json.Marshal(base)
	jws, err := s.Sign(b)
	if err != nil {
		t.Fatal(err)
	}
	raw, err := jws.CompactSerialize()
	if err != nil {
		t.Fatal(err)
	}
	return raw
}

func TestOIDCVerifierReadsClaims(t *testing.T) {
	s, pub := signer(t)
	v := NewOIDCVerifierWithKeys(issuer, audience, &oidc.StaticKeySet{PublicKeys: []crypto.PublicKey{pub}})
	cases := []struct {
		name   string
		claims map[string]any
		want   Principal
	}{
		{"okta scp and cid", map[string]any{"email": "a@b.c", "scp": []string{"openid", "platform:admin"}, "cid": console},
			Principal{Subject: "00uPerson", Email: "a@b.c", Scopes: []string{"openid", "platform:admin"}, ClientID: console}},
		{"scope string and azp, merged once", map[string]any{"scope": "a b", "scp": []string{"b", "c"}, "azp": portal},
			Principal{Subject: "00uPerson", Scopes: []string{"a", "b", "c"}, ClientID: portal}},
		{"cid wins over azp", map[string]any{"cid": console, "azp": portal},
			Principal{Subject: "00uPerson", ClientID: console}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := v.Verify(context.Background(), token(t, s, tc.claims))
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("got %+v, want %+v", got, tc.want)
			}
		})
	}
}

func TestOIDCVerifierRefuses(t *testing.T) {
	s, pub := signer(t)
	other, _ := signer(t)
	v := NewOIDCVerifierWithKeys(issuer, audience, &oidc.StaticKeySet{PublicKeys: []crypto.PublicKey{pub}})
	for name, raw := range map[string]string{
		"wrong audience": token(t, s, map[string]any{"aud": "someone-else"}),
		"wrong issuer":   token(t, s, map[string]any{"iss": "https://evil.example"}),
		"expired":        token(t, s, map[string]any{"exp": time.Now().Add(-time.Minute).Unix()}),
		"wrong key":      token(t, other, nil),
		"garbage":        "not.a.token",
	} {
		if _, err := v.Verify(context.Background(), raw); err == nil {
			t.Errorf("%s: verified", name)
		}
	}
}

func TestParseClientPolicy(t *testing.T) {
	p, err := ParseClientPolicy(" platform:admin = " + console + ", cli ; identity:audit=" + console + ";")
	if err != nil {
		t.Fatal(err)
	}
	if got := p.Guarded(); !reflect.DeepEqual(got, []string{"identity:audit", "platform:admin"}) {
		t.Fatalf("guarded %v", got)
	}
	if !p.Allows("cli", "platform:admin") || p.Allows("cli", "identity:audit") || !p.Allows(portal, "escrow:read") {
		t.Fatal("allows wrong")
	}
	if empty, err := ParseClientPolicy(""); err != nil || len(empty.Guarded()) != 0 {
		t.Fatalf("empty: %v %v", empty, err)
	}
	for _, bad := range []string{"platform:admin", "platform:admin=", "platform:admin= , ", "=x", "a=x;a=y"} {
		if _, err := ParseClientPolicy(bad); err == nil {
			t.Errorf("%q parsed", bad)
		}
	}
}

func TestClientPolicyFilter(t *testing.T) {
	p, _ := NewClientPolicy(map[string][]string{"platform:admin": {console}})
	kept, withheld := p.Filter(portal, []string{"openid", "platform:admin"})
	if !reflect.DeepEqual(kept, []string{"openid"}) || !reflect.DeepEqual(withheld, []string{"platform:admin"}) {
		t.Fatalf("kept %v withheld %v", kept, withheld)
	}
	if _, withheld := p.Filter("", []string{"platform:admin"}); len(withheld) != 1 {
		t.Fatal("a token without a client kept a guarded scope")
	}
	var none *ClientPolicy
	if kept, _ := none.Filter("", []string{"platform:admin"}); len(kept) != 1 {
		t.Fatal("nil policy withheld")
	}
}

type fakeVerifier struct {
	p   Principal
	err error
}

func (f fakeVerifier) Verify(context.Context, string) (Principal, error) { return f.p, f.err }

func TestMiddleware(t *testing.T) {
	policy, _ := NewClientPolicy(map[string][]string{"platform:admin": {console}})
	person := Principal{Subject: "00uPerson", Scopes: []string{"openid", "platform:admin"}}
	bearer := map[string]string{"Authorization": "Bearer x"}

	run := func(v TokenVerifier, hdr map[string]string, opts ...Option) (*httptest.ResponseRecorder, Principal) {
		var seen Principal
		h := Middleware(v, opts...)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { seen, _ = FromContext(r.Context()) }))
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/ops/jobs", nil)
		for k, v := range hdr {
			req.Header.Set(k, v)
		}
		h.ServeHTTP(rec, req)
		return rec, seen
	}

	t.Run("no token is a 401 envelope", func(t *testing.T) {
		rec, _ := run(fakeVerifier{p: person}, nil)
		if rec.Code != http.StatusUnauthorized || rec.Header().Get("Content-Type") != "application/json" {
			t.Fatalf("%d %s", rec.Code, rec.Header().Get("Content-Type"))
		}
	})
	t.Run("bad token is a 401", func(t *testing.T) {
		if rec, _ := run(fakeVerifier{err: errors.New("bad")}, bearer); rec.Code != http.StatusUnauthorized {
			t.Fatal(rec.Code)
		}
	})
	t.Run("allowed client keeps the operator scope", func(t *testing.T) {
		p := person
		p.ClientID = console
		_, seen := run(fakeVerifier{p: p}, bearer, WithClientPolicy(policy))
		if !seen.HasScope("platform:admin") {
			t.Fatalf("%+v", seen)
		}
	})
	t.Run("other client loses it, and it is logged without the person", func(t *testing.T) {
		core, logs := observer.New(zap.WarnLevel)
		p := person
		p.ClientID = portal
		rec, seen := run(fakeVerifier{p: p}, bearer, WithClientPolicy(policy), WithLogger(zap.New(core)))
		if rec.Code != http.StatusOK || seen.HasScope("platform:admin") || !seen.HasScope("openid") {
			t.Fatalf("%d %+v", rec.Code, seen)
		}
		if logs.Len() != 1 {
			t.Fatalf("%d log lines", logs.Len())
		}
		f := logs.All()[0].ContextMap()
		if f["client"] != portal || f["scope"] != "platform:admin" || len(f) != 3 {
			t.Fatalf("fields %v", f)
		}
	})
	t.Run("dev bypass without a client keeps its scopes", func(t *testing.T) {
		_, seen := run(nil, map[string]string{"X-Dev-User": "Op@Example.com"}, WithDevBypass("platform:admin"), WithClientPolicy(policy))
		if !seen.DevBypass || !seen.HasScope("platform:admin") || seen.Subject != "dev-op-example-com" {
			t.Fatalf("%+v", seen)
		}
	})
	t.Run("dev bypass naming a client is guarded", func(t *testing.T) {
		_, seen := run(nil, map[string]string{"X-Dev-Client": portal}, WithDevBypass("platform:admin"), WithClientPolicy(policy))
		if seen.HasScope("platform:admin") {
			t.Fatalf("%+v", seen)
		}
	})
}

func TestRequireScope(t *testing.T) {
	ok := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNoContent) })
	h := RequireScope("platform:admin")(ok)
	for _, tc := range []struct {
		p    *Principal
		want int
	}{
		{nil, http.StatusForbidden},
		{&Principal{Scopes: []string{"openid"}}, http.StatusForbidden},
		{&Principal{Scopes: []string{"platform:admin"}}, http.StatusNoContent},
	} {
		req := httptest.NewRequest(http.MethodPost, "/ops/jobs/x/run", nil)
		if tc.p != nil {
			req = req.WithContext(WithPrincipal(req.Context(), *tc.p))
		}
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		if rec.Code != tc.want {
			t.Errorf("%+v: %d, want %d", tc.p, rec.Code, tc.want)
		}
	}
}

func TestActorMiddleware(t *testing.T) {
	refs, err := auditref.New(auditref.Config{Environment: "dev"})
	if err != nil {
		t.Fatal(err)
	}
	var actor string
	h := ActorMiddleware(refs)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { actor = commonslog.Actor(r.Context()) }))
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx := commonslog.WithActorSlot(WithPrincipal(req.Context(), Principal{Subject: "00uPerson"}))
	h.ServeHTTP(httptest.NewRecorder(), req.WithContext(ctx))
	if actor != refs.Subject("00uPerson") {
		t.Fatalf("actor %q", actor)
	}
}
