package cantonid

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/json"
	"flag"
	"os"
	"strings"
	"testing"
)

var update = flag.Bool("update", false, "rewrite golden fixtures")

const ns = "1220ebb7b95ec6e5cfa8dbf6aa11981e4e1645b1926e4a5afb2c46d7ef49f7ca4288"

// parseCase is one row of testdata/parse.json -- the contract the
// TypeScript mirror (@vdatacloud/cx-commons sdk/canton-id) is tested
// against: how each input classifies, splits and shortens.
type parseCase struct {
	Input     string `json:"input"`
	Kind      Kind   `json:"kind"`
	Hint      string `json:"hint,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Short     string `json:"short"`
}

var parseInputs = []string{
	"relaytest::" + ns,        // external wallet party
	"Depositor::" + ns,        // participant-hosted party (shared namespace)
	"my party:1 - x_y::" + ns, // every hint character Canton allows
	ns,                        // a fingerprint (or hash)
	"Depositor",               // a bare hint
	"Depositor::1220ab",       // short namespace: not a party
	"Custodian::1220" + strings.Repeat("ff", 32) + "00", // namespace too long
	"bad hint!::" + ns,       // hint with a disallowed character
	"a::b::" + ns,            // hint containing the delimiter
	"::" + ns,                // empty hint
	strings.ToUpper(ns),      // uppercase hex: not canonical
	"1e20" + ns[4:],          // a non-sha256 multihash
	strings.Repeat("3f", 64), // a hex signature
	"VXDrn/op6YbtZHuSH+xdWUR8kW7+xMrP9w47KgxPZSE=", // a base64 hash
	"someone-else::" + strings.Repeat("0", 70),     // malformed echoed input
	"w-ebb7b95ec6e5cfa8dbf6aa11",                   // a w- ledger user
	"",
}

func TestParse_Golden(t *testing.T) {
	var got []parseCase
	for _, in := range parseInputs {
		c := parseCase{Input: in, Kind: Classify(in), Short: Short(in)}
		if p, err := ParsePartyID(in); err == nil {
			c.Hint, c.Namespace = p.Hint, string(p.Namespace)
			if p.String() != in || p.Short() != c.Short {
				t.Errorf("%q: round trip %q / short %q", in, p.String(), p.Short())
			}
		}
		got = append(got, c)
	}
	b, _ := json.MarshalIndent(got, "", "  ")
	if *update {
		if err := os.WriteFile("testdata/parse.json", append(b, '\n'), 0o644); err != nil {
			t.Fatal(err)
		}
		return
	}
	want, err := os.ReadFile("testdata/parse.json")
	if err != nil {
		t.Fatalf("%v (run go test -update)", err)
	}
	if strings.TrimRight(string(want), "\n") != string(b) {
		t.Errorf("testdata/parse.json differs:\n%s", b)
	}
}

func TestPartyID(t *testing.T) {
	p, err := ParsePartyID("relaytest::" + ns)
	if err != nil || p.Hint != "relaytest" || p.Namespace != Fingerprint(ns) {
		t.Fatalf("%+v %v", p, err)
	}
	if p.Short() != "relaytest::1220ebb7…4288" || p.String() != "relaytest::"+ns {
		t.Errorf("forms: %s / %s", p.Short(), p)
	}
	for _, bad := range []string{"Depositor", "Depositor::1220ab", "::" + ns, "a::b::" + ns, "bad hint!::" + ns} {
		if _, err := ParsePartyID(bad); err == nil {
			t.Errorf("%q must not parse", bad)
		}
	}
	var zero PartyID
	if !zero.IsZero() || zero.String() != "" || zero.Short() != "" {
		t.Error("zero party")
	}
	defer func() {
		if recover() == nil {
			t.Error("MustParsePartyID must panic on a bad id")
		}
	}()
	MustParsePartyID("nope")
}

func TestFingerprintOf_ControlsExternalParty(t *testing.T) {
	pub, _, _ := ed25519.GenerateKey(rand.Reader)
	fp := FingerprintOf(pub)
	if _, err := ParseFingerprint(string(fp)); err != nil {
		t.Fatalf("FingerprintOf must produce a valid fingerprint: %v", err)
	}
	party := PartyID{Hint: "wallet", Namespace: fp}
	other, _, _ := ed25519.GenerateKey(rand.Reader)
	if !party.ControlledBy(fp) || party.ControlledBy(FingerprintOf(other)) || party.ControlledBy("") {
		t.Error("ControlledBy")
	}
	// The formula is the one Canton 3.5 uses: live on LocalNet, the party
	// generate-topology proposes for a key has exactly this namespace (see
	// daml-escrow's TestAPILifecycle_WalletOnboardsAndSignsIn).
}

func TestJSON(t *testing.T) {
	type doc struct {
		Party PartyID     `json:"party,omitzero"`
		Key   Fingerprint `json:"key,omitzero"`
		Hash  Hash        `json:"hash,omitzero"`
	}
	in := doc{Party: MustParsePartyID("relaytest::" + ns), Key: Fingerprint(ns), Hash: Hash(ns)}
	b, err := json.Marshal(in)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != `{"party":"relaytest::`+ns+`","key":"`+ns+`","hash":"`+ns+`"}` {
		t.Fatalf("JSON carries full values: %s", b)
	}
	var out doc
	if err := json.Unmarshal(b, &out); err != nil || out != in {
		t.Fatalf("round trip: %v %+v", err, out)
	}
	if b, _ := json.Marshal(doc{}); string(b) != `{}` {
		t.Errorf("zero values omitted: %s", b)
	}
	for _, bad := range []string{`{"party":"Depositor"}`, `{"key":"1220ab"}`, `{"hash":"xyz"}`, `{"party":1}`} {
		if err := json.Unmarshal([]byte(bad), &out); err == nil {
			t.Errorf("%s must not decode", bad)
		}
	}
	if err := json.Unmarshal([]byte(`{"party":"","key":"","hash":""}`), &out); err != nil || !out.Party.IsZero() || !out.Key.IsZero() || !out.Hash.IsZero() {
		t.Errorf("empty strings decode to zero: %v %+v", err, out)
	}
}

func TestShortForms(t *testing.T) {
	for in, want := range map[string]string{
		"relaytest::" + ns:       "relaytest::1220ebb7…4288",
		ns:                       "1220ebb7…4288",
		strings.Repeat("3f", 64): "3f3f3f3f…3f3f",
		"Depositor::1220ab":      "Depositor::1220ab",
		"Depositor":              "Depositor",
		"":                       "",
		"🔑🔑🔑🔑🔑🔑🔑🔑🔑🔑🔑🔑🔑🔑🔑🔑🔑🔑": "🔑🔑🔑🔑🔑🔑🔑🔑…🔑🔑🔑🔑",
	} {
		if got := Short(in); got != want {
			t.Errorf("Short(%q) = %q, want %q", in, got, want)
		}
	}
	if Hash(ns).Short() != "1220ebb7…4288" || Fingerprint(ns).Short() != "1220ebb7…4288" {
		t.Error("typed short forms")
	}
}
