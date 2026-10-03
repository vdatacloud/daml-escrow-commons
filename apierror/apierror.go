// Package apierror is the platform's canonical API error: one JSON envelope
// (schema/error.schema.json) every service answers a refusal with, and every
// client -- Go or the cx-commons TypeScript mirror -- parses, formats and
// logs the same way.
//
//	{"error":"human message","code":"LEDGER_SIGNATURE_REJECTED","status":422,
//	 "stage":"execute","hint":"what to do next","requestId":"…",
//	 "upstream":{"service":"canton","status":400,"code":"…","grpcCode":"…","cause":"…","traceId":"…","node":"…"},
//	 "details":{…full values…}}
//
// Long cryptographic identifiers appear shortened (ShortID) in message,
// hint and logs; their full values live in details. An error marshals at
// one of two detail levels: Full (everything) or Summary (no details, no
// upstream cause) -- for audiences that shouldn't see the full values.
//
// Canonical form: compact JSON, envelope fields in the order above, details
// object keys sorted, no HTML escaping, absent fields omitted. Go
// (Canonical) and TypeScript (canonicalApiError) produce identical bytes;
// testdata/canonical holds the fixtures both test against.
//
// # What an error may disclose
//
// Error bodies leave the service -- often from unauthenticated endpoints --
// and get pasted into tickets, chats and logs. Before adding a field, a
// detail or a word to a message, check it against these rules:
//
//   - May include: what the caller sent (echoed back); values derived only
//     from the caller's own credentials (e.g. its key's fingerprint); and,
//     once the caller is verified (an authenticated session, or a check it
//     just passed), facts about its OWN resources and the upstream's account
//     of its OWN request.
//   - Must never include: data about another party or user the caller
//     didn't send (no lookups on its behalf, no derived ids for unverified
//     claims); anything that reveals whether something exists for someone
//     else (an existence oracle -- checks before verification must not touch
//     storage, and must answer the same for real and made-up values);
//     identity data (names, emails, organizations); secrets (tokens, keys,
//     passwords, nonces or signatures the caller didn't send); raw internal
//     error text (database, network, hostnames) -- log it, return a generic
//     message with the request id instead.
//   - Details and upstream causes are Full-level only. Services must default
//     to Summary and enable Full only for developer environments; a Writer's
//     zero value is Full, so set Detail explicitly.
//   - Logs get LogFields (short forms only), never the full details.
//
// Domain packages should encode the verified/unverified distinction in
// types (see daml-escrow's verifiedParty) so a claim can't reach Details by
// accident.
package apierror

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
)

// Error is the canonical API error envelope.
type Error struct {
	// Message is the human-readable message (JSON "error", so clients that
	// only read that field keep working).
	Message string `json:"error"`
	// Code is a stable, machine-readable reason: UPPER_SNAKE_CASE.
	Code string `json:"code"`
	// Status is the HTTP status the error is (or would be) sent with.
	Status int `json:"status,omitempty"`
	// Stage names the step of a multi-step flow that failed.
	Stage string `json:"stage,omitempty"`
	// Hint is what to check or do next.
	Hint      string `json:"hint,omitempty"`
	RequestID string `json:"requestId,omitempty"`
	// Upstream is a backing service's own account of the failure.
	Upstream *Upstream `json:"upstream,omitempty"`
	// Details holds full-length values behind shortened ids, per code
	// (e.g. canton.SigningDetails). Omitted at Summary level.
	Details any `json:"details,omitempty"`
}

// Upstream is the backing service's report of a failure it caused.
type Upstream struct {
	// Service names it: "canton", "identity", "circle", ...
	Service string `json:"service"`
	Status  int    `json:"status,omitempty"`
	// Code is the service's own error id.
	Code string `json:"code,omitempty"`
	// GRPCCode is a gRPC status name ("PERMISSION_DENIED"), when the
	// service speaks gRPC.
	GRPCCode string `json:"grpcCode,omitempty"`
	// Cause is the service's explanation, bounded by MaxCause. Omitted at
	// Summary level.
	Cause   string `json:"cause,omitempty"`
	TraceID string `json:"traceId,omitempty"`
	// Node is the instance that answered (e.g. a Canton participant).
	Node string `json:"node,omitempty"`
}

// MaxCause bounds an upstream cause echoed to clients.
const MaxCause = 600

// Generic codes, one per HTTP status class a handler commonly returns.
// Domain packages define their own codes alongside these.
const (
	CodeInvalidRequest  = "INVALID_REQUEST"
	CodeUnauthenticated = "UNAUTHENTICATED"
	CodeForbidden       = "FORBIDDEN"
	CodeNotFound        = "NOT_FOUND"
	CodeConflict        = "CONFLICT"
	CodeUnprocessable   = "UNPROCESSABLE"
	CodeRateLimited     = "RATE_LIMITED"
	CodeInternal        = "INTERNAL"
	CodeNotImplemented  = "NOT_IMPLEMENTED"
	CodeUpstreamError   = "UPSTREAM_ERROR"
	CodeUnavailable     = "UNAVAILABLE"
	CodeTimeout         = "TIMEOUT"
)

// CodeForStatus is the generic code for an HTTP status.
func CodeForStatus(status int) string {
	switch status {
	case http.StatusBadRequest:
		return CodeInvalidRequest
	case http.StatusUnauthorized:
		return CodeUnauthenticated
	case http.StatusForbidden:
		return CodeForbidden
	case http.StatusNotFound:
		return CodeNotFound
	case http.StatusConflict:
		return CodeConflict
	case http.StatusUnprocessableEntity:
		return CodeUnprocessable
	case http.StatusTooManyRequests:
		return CodeRateLimited
	case http.StatusNotImplemented:
		return CodeNotImplemented
	case http.StatusBadGateway:
		return CodeUpstreamError
	case http.StatusServiceUnavailable:
		return CodeUnavailable
	case http.StatusGatewayTimeout:
		return CodeTimeout
	}
	if status >= 400 && status < 500 {
		return CodeInvalidRequest
	}
	return CodeInternal
}

var codeRe = regexp.MustCompile(`^[A-Z][A-Z0-9_]*$`)

// ValidCode reports whether code is UPPER_SNAKE_CASE.
func ValidCode(code string) bool { return codeRe.MatchString(code) }

// New is an error with status, code and message.
func New(status int, code, message string) *Error {
	return &Error{Status: status, Code: code, Message: message}
}

// Newf is New with a formatted message.
func Newf(status int, code, format string, args ...any) *Error {
	return New(status, code, fmt.Sprintf(format, args...))
}

// Error makes *Error a Go error: "CODE: message".
func (e *Error) Error() string { return e.Code + ": " + e.Message }

// WithStage, WithHint, WithDetails and WithUpstream set a field and return
// e, for chaining.
func (e *Error) WithStage(stage string) *Error  { e.Stage = stage; return e }
func (e *Error) WithHint(hint string) *Error    { e.Hint = hint; return e }
func (e *Error) WithDetails(details any) *Error { e.Details = details; return e }
func (e *Error) WithUpstream(u *Upstream) *Error {
	if u != nil && len(u.Cause) > MaxCause {
		c := *u
		c.Cause = u.Cause[:MaxCause] + "..."
		u = &c
	}
	e.Upstream = u
	return e
}

// Detail is how much of an error to marshal.
type Detail int

const (
	// Full includes details and the upstream cause.
	Full Detail = iota
	// Summary drops details and the upstream cause (which can repeat the
	// full values): code, message, hint and trace ids remain.
	Summary
)

// At returns a copy of e at detail level d.
func (e Error) At(d Detail) Error {
	if d == Summary {
		e.Details = nil
		if e.Upstream != nil {
			u := *e.Upstream
			u.Cause = ""
			e.Upstream = &u
		}
	}
	return e
}

// Canonical is e's canonical JSON at detail level d (see the package doc).
func (e Error) Canonical(d Detail) ([]byte, error) {
	return e.At(d).MarshalJSON()
}

// wire is Error without its methods, so MarshalJSON can encode it.
type wire Error

// MarshalJSON always produces the canonical form (Full level -- use At or
// Canonical for Summary): details keys sorted, no HTML escaping.
func (e Error) MarshalJSON() ([]byte, error) {
	w := wire(e)
	if w.Details != nil {
		sorted, err := sortedJSON(w.Details)
		if err != nil {
			return nil, fmt.Errorf("apierror: details: %w", err)
		}
		w.Details = sorted
	}
	return encode(w)
}

// UnmarshalJSON decodes a canonical error; details stay raw until
// DetailsInto.
func (e *Error) UnmarshalJSON(b []byte) error {
	var w struct {
		wire
		Details json.RawMessage `json:"details,omitempty"`
	}
	if err := json.Unmarshal(b, &w); err != nil {
		return err
	}
	*e = Error(w.wire)
	e.Details = nil
	if len(w.Details) > 0 && string(w.Details) != "null" {
		e.Details = w.Details
	}
	return nil
}

// DetailsInto decodes e.Details into v (e.g. *canton.SigningDetails),
// whatever form it holds -- a typed value or raw JSON from Parse.
func (e *Error) DetailsInto(v any) error {
	if e.Details == nil {
		return fmt.Errorf("apierror: %s has no details", e.Code)
	}
	raw, ok := e.Details.(json.RawMessage)
	if !ok {
		var err error
		if raw, err = json.Marshal(e.Details); err != nil {
			return fmt.Errorf("apierror: details: %w", err)
		}
	}
	if err := json.Unmarshal(raw, v); err != nil {
		return fmt.Errorf("apierror: details: %w", err)
	}
	return nil
}

// sortedJSON re-encodes v through a generic value, which encoding/json
// writes with object keys sorted -- the canonical details form whatever
// Go type v is.
func sortedJSON(v any) (json.RawMessage, error) {
	raw, ok := v.(json.RawMessage)
	if !ok {
		var err error
		if raw, err = json.Marshal(v); err != nil {
			return nil, err
		}
	}
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	var generic any
	if err := dec.Decode(&generic); err != nil {
		return nil, err
	}
	return encode(generic)
}

func encode(v any) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		return nil, err
	}
	return bytes.TrimRight(buf.Bytes(), "\n"), nil
}

// Parse reads any error response body into an Error: the canonical
// envelope, a legacy {"error":"..."} body, or plain text -- the last two
// get the generic code for status. Never returns nil.
func Parse(status int, body []byte) *Error {
	trimmed := bytes.TrimSpace(body)
	var e Error
	if json.Unmarshal(trimmed, &e) == nil && e.Message != "" {
		if e.Code == "" {
			e.Code = CodeForStatus(status)
		}
		if e.Status == 0 {
			e.Status = status
		}
		return &e
	}
	msg := string(trimmed)
	if msg == "" {
		msg = http.StatusText(status)
	}
	return &Error{Message: msg, Code: CodeForStatus(status), Status: status}
}
