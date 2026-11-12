package identityclient

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

// A token source replaces the fixed token and is asked again on every request, so a refreshed
// token reaches the next call.
func TestWithTokenSource_SendsCurrentToken(t *testing.T) {
	var seen []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen = append(seen, r.Header.Get("Authorization"))
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(Identity{OktaSub: "sub-1"})
	}))
	defer server.Close()

	var n atomic.Int32
	c := New(server.URL, "fixed-token", WithTokenSource(func(context.Context) (string, error) {
		if n.Add(1) == 1 {
			return "first", nil
		}
		return "refreshed", nil
	}))
	for range 2 {
		if _, err := c.GetByOktaSub(context.Background(), "sub-1"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}
	if len(seen) != 2 || seen[0] != "Bearer first" || seen[1] != "Bearer refreshed" {
		t.Fatalf("authorization headers = %q, want first then refreshed", seen)
	}
}

// A token source that fails fails the request before anything is sent.
func TestWithTokenSource_ErrorFailsRequest(t *testing.T) {
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls.Add(1)
	}))
	defer server.Close()

	boom := errors.New("token endpoint down")
	c := New(server.URL, "", WithTokenSource(func(context.Context) (string, error) { return "", boom }))

	if _, err := c.Upsert(context.Background(), "sub-1", "a@b.com", "A", "Depositor"); !errors.Is(err, boom) {
		t.Fatalf("Upsert error = %v, want the token source's", err)
	}
	if _, err := c.ManagesIdentity(context.Background(), "m", "t"); !errors.Is(err, boom) {
		t.Fatalf("ManagesIdentity error = %v, want the token source's", err)
	}
	if _, err := c.CreatePendingInvitation(context.Background(), "a@b.com", "m", "ctx"); !errors.Is(err, boom) {
		t.Fatalf("CreatePendingInvitation error = %v, want the token source's", err)
	}
	if calls.Load() != 0 {
		t.Fatalf("server saw %d requests, want none", calls.Load())
	}
}

// Without a token or a source, no Authorization header is sent.
func TestNoToken_NoHeader(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if h := r.Header.Get("Authorization"); h != "" {
			t.Errorf("unexpected Authorization %q", h)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(Identity{})
	}))
	defer server.Close()
	if _, err := New(server.URL, "").GetByEmail(context.Background(), "a@b.com"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
