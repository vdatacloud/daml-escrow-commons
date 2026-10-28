package approval

import (
	"encoding/json"
	"errors"
	"flag"
	"os"
	"reflect"
	"testing"
)

// testdata/decisions.json is shared: @vdatacloud/cx-commons
// (sdk/pending-approvals) checks its decisionsFor against the same cases.
// After changing the rule, run `go test ./approval -update` and copy it to
// cx-commons.
var update = flag.Bool("update", false, "rewrite testdata/decisions.json from the current implementation")

type vector struct {
	Name      string            `json:"name"`
	Request   Request           `json:"request"`
	Viewer    string            `json:"viewer"`
	Decisions []Decision        `json:"decisions"`
	Refusals  map[string]string `json:"refusals"` // decision -> refusal code
}

func cases() []vector {
	req := Request{ProposedBy: "alice", Approvals: []string{"alice"}}
	stale := req
	stale.Stale = true
	aboutSam := req
	aboutSam.Subject = "sam"
	unlisted := Request{ProposedBy: "alice", Approvals: []string{"bob"}}
	return []vector{
		{Name: "the proposer only withdraws", Request: req, Viewer: "alice"},
		{Name: "another person approves or rejects", Request: req, Viewer: "bob"},
		{Name: "someone who approved decides nothing more", Request: Request{ProposedBy: "alice", Approvals: []string{"alice", "bob"}}, Viewer: "bob"},
		{Name: "the proposer counts as approved though unlisted", Request: unlisted, Viewer: "alice"},
		{Name: "a stale request is only rejected", Request: stale, Viewer: "bob"},
		{Name: "the proposer withdraws a stale request", Request: stale, Viewer: "alice"},
		{Name: "the subject may not approve", Request: aboutSam, Viewer: "sam"},
		{Name: "others approve a request about someone", Request: aboutSam, Viewer: "bob"},
		{Name: "no viewer sees every decision", Request: req, Viewer: ""},
		{Name: "no viewer sees no approve when stale", Request: stale, Viewer: ""},
	}
}

func generate() []vector {
	vs := cases()
	for i := range vs {
		v := &vs[i]
		v.Decisions = Decisions(v.Request, v.Viewer)
		v.Refusals = map[string]string{}
		for _, d := range All {
			var r *Refusal
			if errors.As(Check(v.Request, v.Viewer, d), &r) {
				v.Refusals[string(d)] = r.Code
			}
		}
	}
	return vs
}

func TestDecisionsGolden(t *testing.T) {
	const path = "testdata/decisions.json"
	got := generate()
	if *update {
		b, err := json.MarshalIndent(got, "", "  ")
		if err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, append(b, '\n'), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var want []vector
	if err := json.Unmarshal(b, &want); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("%s is out of date -- run with -update", path)
	}
}

func TestDecisions(t *testing.T) {
	for _, v := range generate() {
		t.Run(v.Name, func(t *testing.T) {
			want := map[string][]Decision{
				"the proposer only withdraws":                     {Withdraw},
				"another person approves or rejects":              {Approve, Reject},
				"someone who approved decides nothing more":       {},
				"the proposer counts as approved though unlisted": {Withdraw},
				"a stale request is only rejected":                {Reject},
				"the proposer withdraws a stale request":          {Withdraw},
				"the subject may not approve":                     {Reject},
				"others approve a request about someone":          {Approve, Reject},
				"no viewer sees every decision":                   {Approve, Reject, Withdraw},
				"no viewer sees no approve when stale":            {Reject, Withdraw},
			}[v.Name]
			if !reflect.DeepEqual(v.Decisions, want) {
				t.Fatalf("decisions %v, want %v", v.Decisions, want)
			}
		})
	}
}

func TestCheckRefusals(t *testing.T) {
	req := Request{ProposedBy: "alice", Approvals: []string{"alice"}, Stale: true, Subject: "sam"}
	for _, c := range []struct {
		viewer string
		d      Decision
		want   error
	}{
		{"bob", "abstain", ErrUnknownDecision},
		{"bob", Withdraw, ErrNotProposer},
		{"alice", Reject, ErrAlreadyApproved},
		{"bob", Approve, ErrStale},
		{"bob", Reject, nil},
		{"alice", Withdraw, nil},
	} {
		if err := Check(req, c.viewer, c.d); !errors.Is(err, c.want) {
			t.Errorf("Check(%s, %s) = %v, want %v", c.viewer, c.d, err, c.want)
		}
	}
	fresh := req
	fresh.Stale = false
	if err := Check(fresh, "sam", Approve); !errors.Is(err, ErrSubject) {
		t.Errorf("the subject approving: %v", err)
	}
	if err := Check(fresh, "", Approve); err != nil {
		t.Errorf("an empty viewer is no subject: %v", err)
	}
	if !Approve.Valid() || Decision("abstain").Valid() {
		t.Error("Valid knows the three decisions")
	}
	if ErrStale.Error() == "" || ErrStale.Code != "stale" {
		t.Error("a refusal has a message and a code")
	}
}
