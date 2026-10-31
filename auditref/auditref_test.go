package auditref

import (
	"strings"
	"testing"
)

var keys = map[string][]byte{"r1": []byte("key-one"), "r2": []byte("key-two")}

const party = "u-joey::12202dbfcf5f9a41157bf4c3672cff7c6e5d60123317f39708fdde1f658115535ab8"

func TestPlainInDevelopment(t *testing.T) {
	for _, env := range []string{"dev", "LOCAL", " test ", "development"} {
		r, err := New(Config{Environment: env})
		if err != nil || !r.Plain() {
			t.Fatalf("%q: plain, no key needed: %v", env, err)
		}
		if r.User("u-joey-depositor-devlocal") != "u-joey-depositor-devlocal" || r.Party(party) != party || r.Subject("") != "" {
			t.Errorf("%q shows identifiers in the clear", env)
		}
	}
}

// TestKeyedFailsClosed: outside development, no key is a startup error,
// never a fall back to plain.
func TestKeyedFailsClosed(t *testing.T) {
	for _, env := range []string{"prod", "staging", ""} {
		if _, err := New(Config{Environment: env}); err == nil {
			t.Errorf("%q without a key must fail", env)
		}
		if _, err := New(Config{Environment: env, Keys: keys, Current: "r3"}); err == nil {
			t.Errorf("%q with no key for the current version must fail", env)
		}
	}
	if _, err := New(Config{Environment: "prod", Keys: map[string][]byte{"r:1": []byte("k")}, Current: "r:1"}); err == nil {
		t.Error("a version holding ':' is refused")
	}
}

func TestKeyed(t *testing.T) {
	r, err := New(Config{Environment: "prod", Keys: keys, Current: "r2"})
	if err != nil || r.Plain() {
		t.Fatalf("keyed: %v", err)
	}
	u := r.User("u-joey-depositor-devlocal")
	if !strings.HasPrefix(u, "ref:r2:u:") || len(u) != len("ref:r2:u:")+hexLen || strings.Contains(u, "joey") {
		t.Errorf("user ref %q", u)
	}
	if u != r.User("u-joey-depositor-devlocal") {
		t.Error("stable")
	}
	if r.Of(Party, "x") == r.Of(User, "x") || r.Of(Party, "x")[len("ref:r2:p:"):] == r.Of(User, "x")[len("ref:r2:u:"):] {
		t.Error("kinds are namespaced: same text, different refs")
	}
	if r.Party(party) == r.Party(strings.Replace(party, "joey", "jimmy", 1)) {
		t.Error("parties on one node differ")
	}
	// Another service with the same key computes the same ref: the join.
	other, _ := New(Config{Environment: "prod", Keys: map[string][]byte{"r2": []byte("key-two")}, Current: "r2"})
	if other.User("u-joey-depositor-devlocal") != u {
		t.Error("services sharing a key agree")
	}
}

func TestAllAndParse(t *testing.T) {
	r, _ := New(Config{Environment: "prod", Keys: keys, Current: "r2"})
	all := r.All(Subject, "00u1abcd")
	if len(all) != 2 || !strings.HasPrefix(all[0], "ref:r2:s:") || !strings.HasPrefix(all[1], "ref:r1:s:") {
		t.Fatalf("All = %v", all)
	}
	if all[0] != r.Subject("00u1abcd") {
		t.Error("the current version's ref is among them")
	}
	old, _ := New(Config{Environment: "prod", Keys: map[string][]byte{"r1": []byte("key-one")}, Current: "r1"})
	if all[1] != old.Subject("00u1abcd") {
		t.Error("a ref logged under r1 still resolves after rotating to r2")
	}
	v, k, err := Parse(all[1])
	if err != nil || v != "r1" || k != Subject {
		t.Errorf("Parse = %q %q %v", v, k, err)
	}
	for _, bad := range []string{"u-joey", "ref:r1:x:0123456789abcdef", "ref:r1:s:short", "ref::s:0123456789abcdef"} {
		if _, _, err := Parse(bad); err == nil {
			t.Errorf("Parse(%q) must fail", bad)
		}
	}
	if r.All(User, "") != nil {
		t.Error("no id, no refs")
	}
}

// TestInstitution: institutions group logs, so they are in the clear by
// default in every environment; InstitutionRefs references them too.
func TestInstitution(t *testing.T) {
	prod, _ := New(Config{Environment: "prod", Keys: keys, Current: "r1"})
	dev, _ := New(Config{Environment: "dev"})
	if prod.Institution("trust-co") != "trust-co" || dev.Institution("trust-co") != "trust-co" {
		t.Error("institutions are in the clear by default")
	}
	refs, _ := New(Config{Environment: "prod", Keys: keys, Current: "r1", InstitutionRefs: true})
	i := refs.Institution("trust-co")
	if !strings.HasPrefix(i, "ref:r1:i:") || i != refs.Institution("trust-co") || i == refs.Institution("bank") {
		t.Errorf("institution ref %q", i)
	}
	if _, k, err := Parse(i); err != nil || k != Institution {
		t.Errorf("Parse = %q %v", k, err)
	}
	if _, err := New(Config{Environment: "dev", InstitutionRefs: true}); err == nil {
		t.Error("institution refs need the key, in development too")
	}
	devRefs, _ := New(Config{Environment: "dev", Keys: keys, Current: "r1", InstitutionRefs: true})
	if devRefs.User("u-joey") != "u-joey" || !strings.HasPrefix(devRefs.Institution("trust-co"), "ref:r1:i:") {
		t.Error("the two switches are independent")
	}
}
