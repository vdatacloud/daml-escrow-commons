package apierror

import (
	"encoding/json"
	"errors"
	"flag"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xeipuuv/gojsonschema"
)

var update = flag.Bool("update", false, "rewrite golden fixtures")

// golden compares got with testdata/<name>, or rewrites it under -update.
// These fixtures are the cross-language contract: @vdatacloud/cx-commons
// tests its TypeScript mirror against copies of them.
func golden(t *testing.T, name string, got []byte) {
	t.Helper()
	path := filepath.Join("testdata", name)
	if *update {
		if err := os.WriteFile(path, append(got, '\n'), 0o644); err != nil {
			t.Fatal(err)
		}
		return
	}
	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("%s: %v (run go test -update)", path, err)
	}
	if strings.TrimRight(string(want), "\n") != string(got) {
		t.Errorf("%s:\n got  %s\n want %s", name, got, want)
	}
}

func validate(t *testing.T, doc []byte) {
	t.Helper()
	res, err := gojsonschema.Validate(gojsonschema.NewBytesLoader(Schema()), gojsonschema.NewBytesLoader(doc))
	if err != nil {
		t.Fatalf("schema: %v", err)
	}
	if !res.Valid() {
		t.Errorf("%s does not match the schema: %v", doc, res.Errors())
	}
}

// signing mirrors canton.SigningDetails' JSON in struct (unsorted) field
// order, proving canonical output sorts details keys whatever the Go type.
type signing struct {
	Party             string   `json:"party,omitempty"`
	SignatureReceived string   `json:"signatureReceived,omitempty"`
	ExpectedHash      string   `json:"expectedHash,omitempty"`
	ActAs             []string `json:"actAs,omitempty"`
}

var party = "relaytest::1220ebb7b95ec6e5cfa8dbf6aa11981e4e1645b1926e4a5afb2c46d7ef49f7ca4288"

var canonicalCases = map[string]*Error{
	"minimal": New(http.StatusNotFound, CodeNotFound, "escrow not found"),
	"full": New(http.StatusUnprocessableEntity, "LEDGER_SIGNATURE_REJECTED", "the participant found no valid signature from the party's key").
		WithStage("execute").
		WithHint("sign the prepared transaction's hash with key "+ShortID(party)+" <not the base64 text> & retry").
		WithUpstream(&Upstream{Service: "canton", Status: 400, Code: "FAILED_TO_EXECUTE_TRANSACTION", GRPCCode: "INVALID_ARGUMENT",
			Cause: "Received 0 valid signatures. Transaction hash to be signed: 1220aa", TraceID: "009fa53edf07ce7f", Node: "app-provider"}).
		WithDetails(signing{Party: party, SignatureReceived: strings.Repeat("3f", 64), ExpectedHash: "1220" + strings.Repeat("aa", 32), ActAs: []string{party}}),
	"upstream-no-cause": New(http.StatusServiceUnavailable, CodeUnavailable, "identity service unavailable").
		WithUpstream(&Upstream{Service: "identity", Status: 503}),
}

func init() {
	canonicalCases["full"].RequestID = "req-01HZX"
}

func TestCanonical_Golden(t *testing.T) {
	for name, e := range canonicalCases {
		for level, suffix := range map[Detail]string{Full: "full", Summary: "summary"} {
			got, err := e.Canonical(level)
			if err != nil {
				t.Fatal(err)
			}
			golden(t, "canonical/"+name+"."+suffix+".json", got)
			validate(t, got)

			// Round trip: parsing the canonical form and re-encoding it
			// gives the same bytes.
			again, err := Parse(e.Status, got).Canonical(Full)
			if err != nil || string(again) != string(got) {
				t.Errorf("%s.%s round trip:\n got  %s\n want %s", name, suffix, again, got)
			}
		}
	}
}

func TestCanonical_Properties(t *testing.T) {
	full, _ := canonicalCases["full"].Canonical(Full)
	s := string(full)
	if strings.Contains(s, `\u003c`) || !strings.Contains(s, "<not the base64 text> & retry") {
		t.Error("canonical JSON must not HTML-escape")
	}
	order := []string{`"error"`, `"code"`, `"status"`, `"stage"`, `"hint"`, `"requestId"`, `"upstream"`, `"details"`}
	last := -1
	for _, k := range order {
		i := strings.Index(s, k)
		if i < last {
			t.Errorf("envelope field %s out of order", k)
		}
		last = i
	}
	d := s[strings.Index(s, `"details"`):]
	keys := []string{`"actAs"`, `"expectedHash"`, `"party"`, `"signatureReceived"`}
	for i := 1; i < len(keys); i++ {
		if strings.Index(d, keys[i-1]) > strings.Index(d, keys[i]) {
			t.Errorf("details keys must be sorted: %s", d)
		}
	}

	sum, _ := canonicalCases["full"].Canonical(Summary)
	if strings.Contains(string(sum), "details") || strings.Contains(string(sum), "cause") || strings.Contains(string(sum), strings.Repeat("3f", 64)) {
		t.Errorf("summary must drop details and the upstream cause: %s", sum)
	}
	if !strings.Contains(string(sum), `"traceId":"009fa53edf07ce7f"`) {
		t.Error("summary keeps trace ids")
	}
	if canonicalCases["full"].Details == nil || canonicalCases["full"].Upstream.Cause == "" {
		t.Error("At(Summary) must not mutate the original")
	}
}

func TestParse(t *testing.T) {
	cases := []struct {
		name, body string
		status     int
		code, msg  string
	}{
		{"canonical", `{"error":"no","code":"LEDGER_CONFLICT","status":409}`, 409, "LEDGER_CONFLICT", "no"},
		{"canonical without status", `{"error":"no","code":"X_Y"}`, 418, "X_Y", "no"},
		{"legacy json", `{"error":"invalid request body"}`, 400, CodeInvalidRequest, "invalid request body"},
		{"plain text", "invalid or expired nonce\n", 401, CodeUnauthenticated, "invalid or expired nonce"},
		{"empty", "", 503, CodeUnavailable, "Service Unavailable"},
		{"json without error", `{"message":"x"}`, 500, CodeInternal, `{"message":"x"}`},
	}
	for _, c := range cases {
		e := Parse(c.status, []byte(c.body))
		if e.Code != c.code || e.Message != c.msg || e.Status != c.status {
			t.Errorf("%s: got %+v", c.name, e)
		}
	}
}

func TestDetailsInto(t *testing.T) {
	full, _ := canonicalCases["full"].Canonical(Full)
	e := Parse(422, full)
	var d signing
	if err := e.DetailsInto(&d); err != nil || d.Party != party || len(d.ActAs) != 1 {
		t.Fatalf("parsed details: %v %+v", err, d)
	}
	var d2 signing
	if err := canonicalCases["full"].DetailsInto(&d2); err != nil || d2.Party != party {
		t.Fatalf("typed details: %v %+v", err, d2)
	}
	if err := New(400, CodeInvalidRequest, "x").DetailsInto(&d2); err == nil {
		t.Error("no details must be an error")
	}
}

func TestError_IsAGoError(t *testing.T) {
	var err error = New(409, CodeConflict, "stale")
	var e *Error
	if !errors.As(err, &e) || err.Error() != "CONFLICT: stale" {
		t.Fatalf("%v", err)
	}
}

func TestWithUpstream_BoundsCause(t *testing.T) {
	e := New(502, CodeUpstreamError, "x").WithUpstream(&Upstream{Service: "canton", Cause: strings.Repeat("c", 5000)})
	if len(e.Upstream.Cause) != MaxCause+3 {
		t.Fatalf("cause length %d", len(e.Upstream.Cause))
	}
}

func TestCodeForStatus(t *testing.T) {
	for status, code := range map[int]string{400: CodeInvalidRequest, 401: CodeUnauthenticated, 403: CodeForbidden, 404: CodeNotFound,
		409: CodeConflict, 422: CodeUnprocessable, 429: CodeRateLimited, 418: CodeInvalidRequest, 500: CodeInternal, 501: CodeNotImplemented,
		502: CodeUpstreamError, 503: CodeUnavailable, 504: CodeTimeout} {
		if got := CodeForStatus(status); got != code || !ValidCode(got) {
			t.Errorf("CodeForStatus(%d) = %s", status, got)
		}
	}
	if ValidCode("lower_case") || ValidCode("") || !ValidCode("A1_B") {
		t.Error("ValidCode")
	}
}

func TestShortID(t *testing.T) {
	ns := "1220ebb7b95ec6e5cfa8dbf6aa11981e4e1645b1926e4a5afb2c46d7ef49f7ca4288"
	for in, want := range map[string]string{
		"relaytest::" + ns:      "relaytest::1220ebb7…4288",
		ns:                      "1220ebb7…4288",
		strings.Repeat("3f", 64): "3f3f3f3f…3f3f",
		"VXDrn/op6YbtZHuSH+xdWUR8kW7+xMrP9w47KgxPZSE=": "VXDrn/op…ZSE=",
		"Depositor":         "Depositor",
		"Depositor::1220ab": "Depositor::1220ab",
		"":                  "",
	} {
		if got := ShortID(in); got != want {
			t.Errorf("ShortID(%q) = %q, want %q", in, got, want)
		}
	}
	if got := ShortIDs([]string{"a::" + ns, "b"}); got[0] != "a::1220ebb7…4288" || got[1] != "b" {
		t.Errorf("ShortIDs: %v", got)
	}
}

func TestLogFields_ShortFormsOnly(t *testing.T) {
	fields := canonicalCases["full"].LogFields()
	got := map[string]string{}
	var keys []string
	for _, f := range fields {
		got[f.Key] = f.Value
		keys = append(keys, f.Key)
	}
	want := map[string]string{
		"code": "LEDGER_SIGNATURE_REJECTED", "status": "422", "stage": "execute", "requestId": "req-01HZX",
		"upstream.service": "canton", "upstream.code": "FAILED_TO_EXECUTE_TRANSACTION", "upstream.grpcCode": "INVALID_ARGUMENT",
		"upstream.traceId": "009fa53edf07ce7f", "upstream.node": "app-provider",
		"details.party": "relaytest::1220ebb7…4288", "details.signatureReceived": "3f3f3f3f…3f3f",
		"details.expectedHash": "1220aaaa…aaaa", "details.actAs": "relaytest::1220ebb7…4288",
	}
	for k, v := range want {
		if got[k] != v {
			t.Errorf("%s = %q, want %q", k, got[k], v)
		}
	}
	if keys[0] != "code" {
		t.Errorf("code first: %v", keys)
	}
	for _, f := range fields {
		if strings.Contains(f.Value, ebb7Full) {
			t.Errorf("log field %s carries a full id", f.Key)
		}
	}
}

const ebb7Full = "1220ebb7b95ec6e5cfa8dbf6aa11981e4e1645b1926e4a5afb2c46d7ef49f7ca4288"

func TestWriter(t *testing.T) {
	var logged *Error
	h := RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Writer{Detail: Summary, Log: func(_ *http.Request, e *Error) { logged = e }}.
			Write(w, r, New(http.StatusForbidden, "LEDGER_PERMISSION_DENIED", "denied").WithDetails(signing{Party: party}))
	}))

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set(HeaderRequestID, "client-req-7")
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusForbidden || rr.Header().Get("Content-Type") != "application/json" || rr.Header().Get(HeaderRequestID) != "client-req-7" {
		t.Fatalf("response: %d %v", rr.Code, rr.Header())
	}
	var body map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["requestId"] != "client-req-7" || body["details"] != nil {
		t.Errorf("summary body: %v", body)
	}
	if logged == nil || logged.RequestID != "client-req-7" || logged.Details == nil {
		t.Error("the log hook sees the full error")
	}
	validate(t, rr.Body.Bytes())

	// A malformed client id is replaced, and a fresh one generated.
	rr = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set(HeaderRequestID, "bad id with spaces")
	h.ServeHTTP(rr, req)
	if id := rr.Header().Get(HeaderRequestID); id == "bad id with spaces" || len(id) != 24 {
		t.Errorf("request id %q", id)
	}

	// Without the middleware or a status, Write still sends a valid error.
	rr = httptest.NewRecorder()
	Write(rr, httptest.NewRequest(http.MethodGet, "/", nil), &Error{Message: "boom"})
	if rr.Code != 500 || !strings.Contains(rr.Body.String(), `"code":"INTERNAL"`) {
		t.Errorf("defaults: %d %s", rr.Code, rr.Body.String())
	}
}

func TestFromResponse(t *testing.T) {
	resp := &http.Response{StatusCode: 404, Header: http.Header{HeaderRequestID: []string{"up-1"}}}
	e := FromResponse(resp, []byte("not found"))
	if e.Code != CodeNotFound || e.RequestID != "up-1" {
		t.Fatalf("%+v", e)
	}
}

func TestSchemaIsACopy(t *testing.T) {
	s := Schema()
	s[0] = 'X'
	if Schema()[0] == 'X' {
		t.Error("Schema must return a copy")
	}
}
