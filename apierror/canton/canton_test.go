package canton

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/vdatacloud/daml-escrow-commons/apierror"
	"github.com/xeipuuv/gojsonschema"
)

var update = flag.Bool("update", false, "rewrite golden fixtures")

// fixture is one testdata/classify input: a stage and a JSON Ledger API
// response (body: Canton's error object, or a raw string).
type fixture struct {
	Stage  string          `json:"stage"`
	Status int             `json:"status"`
	Body   json.RawMessage `json:"body"`
}

func (f fixture) rawBody() []byte {
	var s string
	if json.Unmarshal(f.Body, &s) == nil {
		return []byte(s)
	}
	return f.Body
}

// TestFromLedgerError_Golden classifies every testdata/classify/*.input.json
// and compares the canonical envelope with *.want.json. The TypeScript
// mirror (@vdatacloud/cx-commons fromCantonError) runs the same fixtures,
// so a UI talking to Canton directly reports refusals identically.
func TestFromLedgerError_Golden(t *testing.T) {
	inputs, _ := filepath.Glob("testdata/classify/*.input.json")
	if len(inputs) == 0 {
		t.Fatal("no fixtures")
	}
	for _, in := range inputs {
		name := strings.TrimSuffix(filepath.Base(in), ".input.json")
		t.Run(name, func(t *testing.T) {
			raw, _ := os.ReadFile(in)
			var f fixture
			if err := json.Unmarshal(raw, &f); err != nil {
				t.Fatal(err)
			}
			got, err := FromLedgerError(f.Stage, ParseLedgerError(f.Status, f.rawBody())).Canonical(apierror.Full)
			if err != nil {
				t.Fatal(err)
			}
			res, err := gojsonschema.Validate(gojsonschema.NewBytesLoader(apierror.Schema()), gojsonschema.NewBytesLoader(got))
			if err != nil || !res.Valid() {
				t.Errorf("schema: %v %v", err, res.Errors())
			}
			path := strings.TrimSuffix(in, ".input.json") + ".want.json"
			if *update {
				_ = os.WriteFile(path, append(got, '\n'), 0o644)
				return
			}
			want, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("%v (run go test -update)", err)
			}
			if strings.TrimRight(string(want), "\n") != string(got) {
				t.Errorf("\n got  %s\n want %s", got, want)
			}
		})
	}
}

func TestFromLedgerError_Mapping(t *testing.T) {
	for _, c := range []struct {
		grpc, status int
		cause        string
		want         string
		http         int
	}{
		{7, 403, "", CodeLedgerPermissionDenied, 403},
		{16, 401, "", CodeLedgerUnauthenticated, 502},
		{3, 400, "Received 0 valid signatures ... Transaction hash to be signed: 1220ab", CodeLedgerSignature, 422},
		{3, 400, "bad field", CodeLedgerInvalidArgument, 400},
		{5, 404, "", CodeLedgerNotFound, 409},
		{6, 409, "", CodeLedgerConflict, 409},
		{10, 409, "", CodeLedgerConflict, 409},
		{9, 400, "", CodeLedgerRejected, 422},
		{14, 503, "", CodeLedgerUnavailable, 503},
		{4, 504, "", CodeLedgerUnavailable, 503},
		{13, 500, "", CodeLedgerRejected, 502},
	} {
		body := fmt.Sprintf(`{"code":"X","cause":%q,"grpcCodeValue":%d,"traceId":"t"}`, c.cause, c.grpc)
		e := FromLedgerError("s", ParseLedgerError(c.status, []byte(body)))
		if e.Code != c.want || e.Status != c.http || e.Upstream.Service != Service || e.Upstream.TraceID != "t" || e.Hint == "" {
			t.Errorf("grpc %d: got %s/%d %+v", c.grpc, e.Code, e.Status, e.Upstream)
		}
	}
}

func TestFromError(t *testing.T) {
	le := ParseLedgerError(403, []byte(`{"grpcCodeValue":7}`))
	if e := FromError("execute", fmt.Errorf("execute wallet transaction: %w", le)); e.Code != CodeLedgerPermissionDenied {
		t.Errorf("wrapped LedgerError: %s", e.Code)
	}
	e := FromError("execute", errors.New("dial tcp: connection refused"))
	if e.Code != CodeLedgerUnreachable || e.Status != http.StatusBadGateway || e.Upstream.Service != Service {
		t.Errorf("unreachable: %+v", e)
	}
}

func TestLedgerError(t *testing.T) {
	le := ParseLedgerError(400, []byte(`{"code":"C","cause":"why","traceId":"t1","grpcCodeValue":3,"errorCategory":8,"context":{"participant":"p1"}}`))
	if le.Code != "C" || le.GRPCCodeName() != "INVALID_ARGUMENT" || le.Participant != "p1" || le.Category != 8 ||
		!strings.HasPrefix(le.Error(), "JSON API error (400): {") {
		t.Fatalf("%+v", le)
	}
	if plain := ParseLedgerError(404, []byte("404 page not found")); plain.Code != "" || plain.Error() != "JSON API error (404): 404 page not found" {
		t.Fatalf("%+v", plain)
	}
	if GRPCCodeName(99) != "CODE_99" || GRPCCodeName(-1) != "CODE_-1" {
		t.Error("unknown codes")
	}
}

func TestSigningDetails(t *testing.T) {
	party := "w::1220" + strings.Repeat("ab", 32)
	if d := ForParty(party); d.PartyFingerprint != "1220"+strings.Repeat("ab", 32) || Namespace("Depositor") != "" {
		t.Fatalf("%+v", d)
	}
	if ExpectedHash("Transaction hash to be signed: 1220ff. Ensure") != "1220ff" || ExpectedHash("nothing") != "" {
		t.Error("ExpectedHash")
	}

	// Signing attaches to a classified refusal, keeping what's there.
	e := FromLedgerError("execute", ParseLedgerError(400, []byte(`{"cause":"Transaction hash to be signed: 1220ff.","grpcCodeValue":3}`)))
	d := Signing(e)
	d.SignatureReceived = "abcd"
	if got := Signing(e); got.ExpectedHash != "1220ff" || got.SignatureReceived != "abcd" {
		t.Errorf("%+v", got)
	}
	// ...and to one with no details, or details parsed from JSON.
	plain := apierror.New(401, CodeInvalidSignature, "x")
	Signing(plain).Party = party
	if plain.Details.(*SigningDetails).Party != party {
		t.Error("attach to none")
	}
	body, _ := plain.Canonical(apierror.Full)
	if parsed := apierror.Parse(401, body); Signing(parsed).Party != party {
		t.Error("from parsed JSON")
	}
}
