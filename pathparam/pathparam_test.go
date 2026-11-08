package pathparam

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/go-chi/chi/v5"
)

// A party or subject sent as the browser sends it (encodeURIComponent)
// reaches the handler as itself, not still escaped (Phase 100 edge case).
func TestGetUnescapes(t *testing.T) {
	var got string
	r := chi.NewRouter()
	r.Post("/signing/wallet-links/{walletParty}/revoke", func(w http.ResponseWriter, r *http.Request) {
		got = Get(r, "walletParty")
	})
	for _, want := range []string{
		"devwallet::1220981ccf767aaa7e4bd2c182c89ceb78727b44e6618dfa36d49ad46a046cc0220a",
		"idp-ui-check-1@wallet.devlocal",
		"First Bank",
		"plain-id-123",
		"a+b,c;d=e$f&g",
	} {
		for _, path := range []string{
			"/signing/wallet-links/" + url.PathEscape(want) + "/revoke",       // what Go clients send
			"/signing/wallet-links/" + jsEncodeURIComponent(want) + "/revoke", // what the portal sends
		} {
			got = ""
			r.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodPost, path, nil))
			if got != want {
				t.Errorf("%s: got %q, want %q", path, got, want)
			}
		}
	}
}

// jsEncodeURIComponent is the browser's encodeURIComponent for these tests:
// everything but A-Z a-z 0-9 - _ . ! ~ * ' ( ) escaped.
func jsEncodeURIComponent(s string) string {
	const keep = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_.!~*'()"
	out := make([]byte, 0, len(s)*3)
	for i := 0; i < len(s); i++ {
		c := s[i]
		if containsByte(keep, c) {
			out = append(out, c)
			continue
		}
		out = append(out, '%', "0123456789ABCDEF"[c>>4], "0123456789ABCDEF"[c&15])
	}
	return string(out)
}

func containsByte(s string, c byte) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			return true
		}
	}
	return false
}

// Not valid escaping: returned as sent.
func TestGetKeepsInvalidEscaping(t *testing.T) {
	var got string
	r := chi.NewRouter()
	r.Get("/x/{v}", func(w http.ResponseWriter, r *http.Request) { got = Get(r, "v") })
	r.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/x/100%25", nil))
	if got != "100%" {
		t.Fatalf("got %q", got)
	}
	req := httptest.NewRequest(http.MethodGet, "/x/a", nil)
	req.URL.RawPath, req.URL.Path = "/x/bad%zz", "/x/bad%zz"
	r.ServeHTTP(httptest.NewRecorder(), req)
	if got != "bad%zz" {
		t.Fatalf("got %q", got)
	}
}
