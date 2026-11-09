package authority

import (
	"encoding/json"
	"net/http/httptest"
	"testing"
)

func TestAuth(t *testing.T) {
	scopes := []string{"platform:admin", "identity:audit"}
	if got := Auth(false, "platform:admin=c1;identity:audit=c1", scopes); len(got) != 0 {
		t.Fatalf("guarded: %v", got)
	}
	if got := Auth(true, "platform:admin=c1", scopes); len(got) != 2 || got[0].Setting != "AUTH_BYPASS" || got[1].Setting != "AUTH_CLIENT_POLICY" {
		t.Fatalf("bypass and unguarded: %v", got)
	}
	if got := Auth(false, "garbage", scopes); len(got) != 1 {
		t.Fatalf("unparseable policy: %v", got)
	}
}

func TestReport(t *testing.T) {
	r := NewReport("identity-api", "dev", "v1", Auth(true, "", nil), Logging("dev"))
	if r.Authoritative || len(r.StandIns) != 2 {
		t.Fatalf("dev: %+v", r)
	}
	r = NewReport("identity-api", "production", "v1", Logging("production"))
	if !r.Authoritative || r.StandIns == nil {
		t.Fatalf("production: %+v", r)
	}
	rec := httptest.NewRecorder()
	Handler(r)(rec, httptest.NewRequest("GET", Path, nil))
	var back Report
	if err := json.Unmarshal(rec.Body.Bytes(), &back); err != nil || !back.Authoritative || back.Service != "identity-api" {
		t.Fatalf("wire: %s", rec.Body.String())
	}
}
