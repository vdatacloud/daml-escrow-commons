package orgsetup

import "testing"

func TestFinish(t *testing.T) {
	s := Setup{Checks: []Check{{ID: CheckParty, Status: Ready}, {ID: CheckMandate, Status: Attention}, {ID: CheckAutomation, Status: NotApplicable}}}
	if s.Finish(); !s.Ready {
		t.Fatal("attention and n/a don't block")
	}
	s.Checks = append(s.Checks, Check{ID: CheckSigning, Status: Missing})
	if s.Finish(); s.Ready {
		t.Fatal("a missing check blocks")
	}
}

func TestPathFor(t *testing.T) {
	if got := PathFor("demo-bank"); got != "/api/v1/institutions/demo-bank/setup" {
		t.Fatal(got)
	}
	if got := PathFor("a/b"); got != "/api/v1/institutions/a%2Fb/setup" {
		t.Fatal(got)
	}
}
