package legalorders

import (
	"encoding/json"
	"strings"
	"testing"
)

func valid() RecordRequest {
	return RecordRequest{Case: "LO-2026-014", Kind: "preservation", OrderType: "subpoena", Reference: "26-cv-1",
		Authority: "SDNY", IssuedOn: "2026-11-01", ReceivedOn: "2026-11-02", Targets: []string{"tenant:ps-1"}}
}

func TestRecordRequestValidate(t *testing.T) {
	if err := valid().Validate(); err != nil {
		t.Fatal(err)
	}
	for name, edit := range map[string]func(*RecordRequest){
		"case with a slash": func(r *RecordRequest) { r.Case = "LO/1" },
		"kind":              func(r *RecordRequest) { r.Kind = "seize" },
		"no reference":      func(r *RecordRequest) { r.Reference = "" },
		"long authority":    func(r *RecordRequest) { r.Authority = strings.Repeat("a", 513) },
		"issued":            func(r *RecordRequest) { r.IssuedOn = "1 Nov 2026" },
		"review":            func(r *RecordRequest) { r.ReviewOn = "soon" },
		"notify":            func(r *RecordRequest) { r.NotifyAfter = "later" },
		"no targets":        func(r *RecordRequest) { r.Targets = nil },
	} {
		r := valid()
		edit(&r)
		if r.Validate() == nil {
			t.Errorf("%s: accepted", name)
		}
	}
}

func TestEventAndLiftValidate(t *testing.T) {
	if err := (EventRequest{Case: "LO-1", Type: "note", Note: "n", NotifyAfter: "never"}).Validate(); err != nil {
		t.Fatal(err)
	}
	for _, e := range []EventRequest{{Case: "LO-1", Type: "memo"}, {Case: "", Type: "note"},
		{Case: "LO-1", Type: "review", NextReviewOn: "next week"}, {Case: "LO-1", Type: "note", Note: strings.Repeat("n", 4001)}} {
		if e.Validate() == nil {
			t.Errorf("accepted %+v", e)
		}
	}
	if (LiftRequest{Case: "LO-1"}).Validate() == nil || (LiftRequest{Case: "LO-1", Reference: "r"}).Validate() != nil {
		t.Error("lift needs the release reference")
	}
}

// The JSON names are escrow-api's and the console's: pinned.
func TestWireNames(t *testing.T) {
	b, _ := json.Marshal(Recorded{Case: "c", Holds: &Sweep{Held: map[string]int{"document": 1}}, HoldsError: "e"})
	want := `{"case":"c","documentUri":"","documentSha256":"","targets":0,"holds":{"escrows":0,"held":{"document":1},"unlockable":0},"holdsError":"e"}`
	if string(b) != want {
		t.Errorf("recorded:\n got %s\nwant %s", b, want)
	}
	b, _ = json.Marshal(Lifted{Case: "c", HoldsReleased: 1})
	if string(b) != `{"case":"c","holdsReleased":1,"objectsReleased":0,"objectsKept":0}` {
		t.Errorf("lifted %s", b)
	}
}
