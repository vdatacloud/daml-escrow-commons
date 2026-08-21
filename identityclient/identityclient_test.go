package identityclient

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_Upsert_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/identities" || r.Method != http.MethodPost {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer test-token" {
			t.Fatalf("expected bearer token header, got %q", got)
		}
		var body map[string]string
		_ = json.NewDecoder(r.Body).Decode(&body)
		if body["oktaSub"] != "sub-1" {
			t.Fatalf("expected oktaSub sub-1, got %q", body["oktaSub"])
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(Identity{OktaSub: "sub-1", IdentityToken: "tok-1", Email: "a@b.com"})
	}))
	defer server.Close()

	c := New(server.URL, "test-token")
	identity, err := c.Upsert(context.Background(), "sub-1", "a@b.com", "Alice", "Depositor")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if identity.IdentityToken != "tok-1" {
		t.Fatalf("expected token tok-1, got %q", identity.IdentityToken)
	}
}

func TestClient_GetByOktaSub_Found(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/identities/sub-1" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(Identity{OktaSub: "sub-1", DamlPartyID: "p-1"})
	}))
	defer server.Close()

	c := New(server.URL, "")
	identity, err := c.GetByOktaSub(context.Background(), "sub-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if identity.DamlPartyID != "p-1" {
		t.Fatalf("expected DamlPartyID p-1, got %q", identity.DamlPartyID)
	}
}

func TestClient_GetByToken_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	c := New(server.URL, "")
	identity, err := c.GetByToken(context.Background(), "missing")
	if err != nil {
		t.Fatalf("expected no error on 404, got %v", err)
	}
	if identity != nil {
		t.Fatalf("expected nil identity on 404, got %+v", identity)
	}
}

func TestClient_GetByEmail_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	c := New(server.URL, "")
	_, err := c.GetByEmail(context.Background(), "a@b.com")
	if err == nil {
		t.Fatal("expected error on 500")
	}
}

func TestClient_ManagesIdentity(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/party-sets/manages" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("manager") != "mgr-1" || r.URL.Query().Get("target") != "tgt-1" {
			t.Fatalf("unexpected query: %s", r.URL.RawQuery)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]bool{"manages": true})
	}))
	defer server.Close()

	c := New(server.URL, "")
	manages, err := c.ManagesIdentity(context.Background(), "mgr-1", "tgt-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !manages {
		t.Fatal("expected manages=true")
	}
}

func TestClient_ManagesIdentity_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	c := New(server.URL, "")
	_, err := c.ManagesIdentity(context.Background(), "mgr-1", "tgt-1")
	if err == nil {
		t.Fatal("expected error on 500")
	}
}
