// Package canton turns a Canton JSON Ledger API refusal into the canonical
// apierror envelope -- for any service that submits ledger commands
// (daml-escrow, daml-escrow-cms), and mirrored in @vdatacloud/cx-commons
// (fromCantonError) for a UI that talks to a participant or wallet gateway
// directly. Pure wire-format work: it parses Canton's error JSON and maps
// its gRPC code; it knows nothing of who any party is.
package canton

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/vdatacloud/daml-escrow-commons/apierror"
	"github.com/vdatacloud/daml-escrow-commons/cantonid"
)

// Service is Upstream.Service for a Canton participant.
const Service = "canton"

// Codes for ledger refusals and for the signing checks a service makes
// before submitting for an external party.
const (
	CodeLedgerPermissionDenied = "LEDGER_PERMISSION_DENIED"
	CodeLedgerUnauthenticated  = "LEDGER_UNAUTHENTICATED"
	CodeLedgerSignature        = "LEDGER_SIGNATURE_REJECTED"
	CodeLedgerInvalidArgument  = "LEDGER_INVALID_ARGUMENT"
	CodeLedgerNotFound         = "LEDGER_NOT_FOUND"
	CodeLedgerConflict         = "LEDGER_CONFLICT"
	CodeLedgerRejected         = "LEDGER_REJECTED"
	CodeLedgerUnavailable      = "LEDGER_UNAVAILABLE"
	CodeLedgerUnreachable      = "LEDGER_UNREACHABLE"

	CodeKeyNotPartyOwner = "KEY_DOES_NOT_CONTROL_PARTY"
	CodeInvalidSignature = "INVALID_SIGNATURE"
	CodeInvalidPublicKey = "INVALID_PUBLIC_KEY"
)

// LedgerError is a non-2xx JSON Ledger API response with Canton's
// structured fields parsed out. Error() is "JSON API error (status): body".
type LedgerError struct {
	Status int
	Body   string
	// Code is Canton's error id (FAILED_TO_EXECUTE_TRANSACTION, ...).
	Code  string
	Cause string
	// GRPCCode is the gRPC status code (7 = PERMISSION_DENIED); 0 if absent.
	GRPCCode    int
	Category    int
	TraceID     string
	Participant string
}

func (e *LedgerError) Error() string {
	return fmt.Sprintf("JSON API error (%d): %s", e.Status, e.Body)
}

// GRPCCodeName is the gRPC status name, e.g. "PERMISSION_DENIED".
func (e *LedgerError) GRPCCodeName() string { return GRPCCodeName(e.GRPCCode) }

// ParseLedgerError parses a JSON Ledger API error response.
func ParseLedgerError(status int, body []byte) *LedgerError {
	e := &LedgerError{Status: status, Body: string(body)}
	var b struct {
		Code          string            `json:"code"`
		Cause         string            `json:"cause"`
		TraceID       string            `json:"traceId"`
		GRPCCodeValue int               `json:"grpcCodeValue"`
		ErrorCategory int               `json:"errorCategory"`
		Context       map[string]string `json:"context"`
	}
	if json.Unmarshal(body, &b) == nil {
		e.Code, e.Cause, e.TraceID, e.GRPCCode, e.Category = b.Code, b.Cause, b.TraceID, b.GRPCCodeValue, b.ErrorCategory
		e.Participant = b.Context["participant"]
	}
	return e
}

var grpcCodeNames = []string{
	"OK", "CANCELLED", "UNKNOWN", "INVALID_ARGUMENT", "DEADLINE_EXCEEDED", "NOT_FOUND", "ALREADY_EXISTS",
	"PERMISSION_DENIED", "RESOURCE_EXHAUSTED", "FAILED_PRECONDITION", "ABORTED", "OUT_OF_RANGE",
	"UNIMPLEMENTED", "INTERNAL", "UNAVAILABLE", "DATA_LOSS", "UNAUTHENTICATED",
}

// GRPCCodeName is a gRPC status code's name ("CODE_n" if unknown).
func GRPCCodeName(code int) string {
	if code < 0 || code >= len(grpcCodeNames) {
		return fmt.Sprintf("CODE_%d", code)
	}
	return grpcCodeNames[code]
}

// SigningDetails are the full values behind a signing or authorization
// refusal. Disclosure (apierror package doc): set Party, PartyFingerprint
// and LedgerUser only for a party the caller is VERIFIED to act for -- for a
// party it merely claimed, leave them empty (the claim can be echoed in the
// message); PublicKeyFingerprint, SignatureReceived and SignedMessage are
// the caller's own and always fine. This is the full detail of a
// refusal (apierror details for the LEDGER_SIGNATURE_REJECTED,
// LEDGER_PERMISSION_DENIED, KEY_DOES_NOT_CONTROL_PARTY and
// INVALID_SIGNATURE codes). Ids the service derived or vouches for are
// typed (cantonid); values echoed from the request (ActAs) stay raw
// strings, since they may not parse.
type SigningDetails struct {
	Party cantonid.PartyID `json:"party,omitzero"`
	// PartyFingerprint is Party's namespace -- the only key that can sign
	// for an external party (compare PublicKeyFingerprint).
	PartyFingerprint cantonid.Fingerprint `json:"partyFingerprint,omitzero"`
	// PublicKeyFingerprint is the fingerprint of the key presented.
	PublicKeyFingerprint cantonid.Fingerprint `json:"publicKeyFingerprint,omitzero"`
	// SignatureReceived is the signature as received (hex).
	SignatureReceived string `json:"signatureReceived,omitempty"`
	// SignedBy is the key fingerprint the signature was submitted under.
	SignedBy cantonid.Fingerprint `json:"signedBy,omitzero"`
	// SignedMessage is what the signature had to cover (a nonce, or a
	// base64 multi-hash).
	SignedMessage string `json:"signedMessage,omitempty"`
	// ExpectedHash is the transaction hash the participant expected
	// signed (from its refusal).
	ExpectedHash cantonid.Hash `json:"expectedHash,omitzero"`
	// ActAs is a refused command's actAs, as the caller sent it.
	ActAs []string `json:"actAs,omitempty"`
	// LedgerUser is the Ledger API user the request was submitted as.
	LedgerUser string `json:"ledgerUser,omitempty"`
}

// ForParty is SigningDetails for party, with its namespace filled in.
func ForParty(party cantonid.PartyID) *SigningDetails {
	return &SigningDetails{Party: party, PartyFingerprint: party.Namespace}
}

var expectedHashRe = regexp.MustCompile(`hash to be signed: ([0-9a-fA-F]+)`)

// ExpectedHash is the transaction hash a signature refusal's cause names
// ("... Transaction hash to be signed: 1220ab... "), or "" if it names
// none (or not a well-formed one).
func ExpectedHash(cause string) cantonid.Hash {
	if m := expectedHashRe.FindStringSubmatch(cause); m != nil {
		if h, err := cantonid.ParseHash(strings.ToLower(m[1])); err == nil {
			return h
		}
	}
	return ""
}

// Signing returns e's details as *SigningDetails, attaching empty ones if
// e has none -- for a caller adding what it knows (party, signature) to a
// classified refusal.
func Signing(e *apierror.Error) *SigningDetails {
	if d, ok := e.Details.(*SigningDetails); ok && d != nil {
		return d
	}
	d := &SigningDetails{}
	_ = e.DetailsInto(d)
	e.Details = d
	return d
}

// FromError classifies err from a ledger call made at stage: a
// *LedgerError (anywhere in its chain) by FromLedgerError, anything else
// as the ledger being unreachable.
func FromError(stage string, err error) *apierror.Error {
	var le *LedgerError
	if errors.As(err, &le) {
		return FromLedgerError(stage, le)
	}
	return apierror.New(http.StatusBadGateway, CodeLedgerUnreachable, "the ledger could not be reached").
		WithStage(stage).
		WithHint("check the participant's JSON Ledger API is up and reachable").
		WithUpstream(&apierror.Upstream{Service: Service})
}

// FromLedgerError maps a ledger refusal by its gRPC code to an HTTP
// status, code, message and hint (identifiers shortened), with Canton's
// own report as Upstream. A signature refusal's expected hash goes in
// SigningDetails. Callers add context they hold (Signing(e)) and may
// replace the hint with a more specific one.
func FromLedgerError(stage string, le *LedgerError) *apierror.Error {
	up := &apierror.Upstream{
		Service: Service, Status: le.Status, Code: le.Code, Cause: le.Cause, TraceID: le.TraceID, Node: le.Participant,
	}
	if le.GRPCCode != 0 {
		up.GRPCCode = le.GRPCCodeName()
	}
	if le.Code == "" && le.GRPCCode == 0 {
		// Not a Canton error body: something in front of (or instead of)
		// the JSON Ledger API answered.
		up.Cause = strings.TrimSpace(le.Body)
		return apierror.New(http.StatusBadGateway, CodeLedgerRejected, "the ledger endpoint answered with a non-Canton error").
			WithStage(stage).
			WithHint("check the JSON Ledger API URL and version, and any proxy in front of it -- see upstream.cause").
			WithUpstream(up)
	}
	expected := ExpectedHash(le.Cause)

	status, code, msg, hint := http.StatusBadGateway, CodeLedgerRejected, "the ledger refused the request", "see upstream.cause"
	switch {
	case le.GRPCCode == 7 || le.Status == http.StatusForbidden:
		status, code, msg = http.StatusForbidden, CodeLedgerPermissionDenied, "the ledger denied the submitting ledger user"
		hint = "the ledger user needs the right this command requires (CanActAs, CanExecuteAs or CanReadAs) on the acting party"
	case le.GRPCCode == 16 || le.Status == http.StatusUnauthorized:
		status, code, msg = http.StatusBadGateway, CodeLedgerUnauthenticated, "the ledger rejected the submitting token"
		hint = "ledger auth is misconfigured (token issuer, audience or secret) -- an operator problem, not the user's"
	case expected != "" || strings.Contains(le.Cause, "valid signature"):
		status, code, msg = http.StatusUnprocessableEntity, CodeLedgerSignature, "the participant found no valid signature from the party's key"
		hint = "sign the prepared transaction's hash, recomputed from the unmodified prepared transaction, with the party's key"
		if expected != "" {
			hint += "; the participant expected hash " + expected.Short() + " (details.expectedHash)"
		}
	case le.GRPCCode == 3:
		status, code, msg = http.StatusBadRequest, CodeLedgerInvalidArgument, "the ledger rejected the request's arguments"
		hint = "see upstream.cause for the field; a command may reference a package or template this participant doesn't have"
	case le.GRPCCode == 5:
		status, code, msg = http.StatusConflict, CodeLedgerNotFound, "the ledger couldn't find something the request needs"
		hint = "a referenced contract may already be archived (stale contract id) or not visible to the acting party -- refresh and retry"
	case le.GRPCCode == 6 || le.GRPCCode == 10:
		status, code, msg = http.StatusConflict, CodeLedgerConflict, "the request conflicts with the ledger's current state"
		hint = "a duplicate or concurrent submission -- refresh and retry"
	case le.GRPCCode == 9:
		status, code, msg = http.StatusUnprocessableEntity, CodeLedgerRejected, "the ledger rejected the command (a contract precondition failed)"
		hint = "see upstream.cause: usually a Daml assertion in the choice"
	case le.GRPCCode == 14 || le.GRPCCode == 4 || le.Status == http.StatusServiceUnavailable:
		status, code, msg = http.StatusServiceUnavailable, CodeLedgerUnavailable, "the ledger is unavailable or timed out"
		hint = "retry shortly; if it persists, check the participant and synchronizer"
	}
	e := apierror.New(status, code, msg).WithStage(stage).WithHint(hint).WithUpstream(up)
	if expected != "" {
		e.Details = &SigningDetails{ExpectedHash: expected}
	}
	return e
}
