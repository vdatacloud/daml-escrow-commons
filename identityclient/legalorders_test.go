package identityclient

import (
	"encoding/json"
	"testing"
	"time"
)

// The legal-order shapes are identity's wire format: the JSON names are the
// contract with daml-escrow-identity and must not drift.
func TestLegalOrderWireNames(t *testing.T) {
	day := time.Date(2026, 11, 8, 0, 0, 0, 0, time.UTC)
	for _, c := range []struct {
		v    any
		want string
	}{
		{LegalOrder{Case: "c", Kind: "preservation", OrderType: "subpoena", Reference: "r", Authority: "a", Jurisdiction: "j",
			IssuedOn: "2026-11-01", ReceivedOn: "2026-11-02", Targets: []LegalOrderTarget{{Key: "escrow:e1", Parties: []string{"p"}}},
			DocumentURI: "storage://x", DocumentSHA256: "h"},
			`{"case":"c","kind":"preservation","orderType":"subpoena","reference":"r","authority":"a","jurisdiction":"j","issuedOn":"2026-11-01","receivedOn":"2026-11-02","targets":[{"key":"escrow:e1","parties":["p"]}],"documentUri":"storage://x","documentSha256":"h"}`},
		{LegalOrderEvent{Case: "c", Type: OrderEventReview, Outcome: ReviewContinue, NextReviewOn: "2026-12-01", RecordedAt: &day},
			`{"case":"c","type":"review","outcome":"continue","nextReviewOn":"2026-12-01","recordedAt":"2026-11-08T00:00:00Z"}`},
		{OrderContact{Target: "escrow:e1", Party: "p", DisplayName: "n", Email: "e"},
			`{"target":"escrow:e1","party":"p","displayName":"n","email":"e"}`},
		{ListLegalOrdersRequest{Open: true, Case: "why"}, `{"open":true,"case":"why"}`},
		{LiftLegalOrderResponse{Case: "c", HoldsReleased: 2}, `{"case":"c","holdsReleased":2}`},
	} {
		b, err := json.Marshal(c.v)
		if err != nil || string(b) != c.want {
			t.Errorf("%T:\n got %s\nwant %s", c.v, b, c.want)
		}
	}
	var o RecordedLegalOrder
	if err := json.Unmarshal([]byte(`{"case":"c","issuedOn":"2026-11-01T00:00:00Z","reviewOn":"2026-12-01T00:00:00Z",
		"targets":["escrow:e1"],"liftedAt":null,"events":[{"type":"note","note":"n"}]}`), &o); err != nil {
		t.Fatal(err)
	}
	if o.ReviewOn == nil || len(o.Events) != 1 || o.Events[0].Note != "n" || o.Targets[0] != "escrow:e1" {
		t.Fatalf("decoded %+v", o)
	}
}

// PLAN.md Phase 104.8: the erasure queue's JSON names are identity's.
func TestErasureRequestWireNames(t *testing.T) {
	b, _ := json.Marshal(SubmitErasureRequest{Email: "e", Channel: "email", ReceivedOn: "2026-11-09", Case: "c"})
	if string(b) != `{"email":"e","channel":"email","receivedOn":"2026-11-09","case":"c"}` {
		t.Errorf("submit %s", b)
	}
	b, _ = json.Marshal(VerifyErasureRequest{ID: 7, Method: "m"})
	if string(b) != `{"id":7,"method":"m"}` {
		t.Errorf("verify %s", b)
	}
}

func TestSigninWireNames(t *testing.T) {
	b, _ := json.Marshal(SigninActionRequest{Email: "a@b", Case: "c", Action: SigninUnlock, Method: "m"})
	if string(b) != `{"email":"a@b","case":"c","action":"unlock","method":"m"}` {
		t.Fatalf("%s", b)
	}
}
