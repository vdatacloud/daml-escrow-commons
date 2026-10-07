package agreementsig

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"unicode"
)

// AgreementSchema names the canonical form of one agreement version. A
// change to what is hashed, or how, is a new schema, never an edit of this
// one: signatures under it must stay verifiable for as long as the
// agreement governs anything.
const AgreementSchema = "tripart.agreement-version/1"

// Kind is what an agreement document is.
type Kind string

const (
	// KindAgreement is an original agreement.
	KindAgreement Kind = "agreement"
	// KindAmendment amends a ratified agreement and every amendment ratified
	// before it.
	KindAmendment Kind = "amendment"
)

// Document identifies the rendered document a version's signers read, by
// content only. Where it is stored is deliberately not signed: moving or
// re-mirroring storage must never invalidate a signature.
type Document struct {
	// SHA256 is the lowercase hex SHA-256 of the rendered file (e.g. the PDF).
	SHA256    string `json:"sha256"`
	MediaType string `json:"mediaType"`
	// SourceSHA256 is the lowercase hex SHA-256 of the canonical source the
	// file was rendered from, which stays stable when a renderer's output
	// bytes (fonts, metadata, timestamps) do not.
	SourceSHA256    string `json:"sourceSha256"`
	SourceMediaType string `json:"sourceMediaType"`
}

// Amends links an amendment to what it amends.
type Amends struct {
	// AgreementID is the original agreement's id.
	AgreementID string `json:"agreementId"`
	// RatifiedHashes are the version hashes of the ratified original and of
	// every amendment ratified since, oldest first.
	RatifiedHashes []string `json:"ratifiedHashes"`
}

// AgreementVersion is the content of one version of an agreement or an
// amendment that its parties sign.
type AgreementVersion struct {
	AgreementID string `json:"agreementId"`
	// Version counts from 1; each later version names its parent's hash.
	Version    int    `json:"version"`
	ParentHash string `json:"parentHash,omitempty"`
	Kind       Kind   `json:"kind"`
	// Amends is set for an amendment, nil for an agreement.
	Amends   *Amends  `json:"amends,omitempty"`
	Document Document `json:"document"`
	// TermsSchema names the schema Terms conform to (e.g. a versioned
	// escrow-terms schema id); Terms is a JSON object hashed as given. This
	// package does not interpret either.
	TermsSchema string          `json:"termsSchema"`
	Terms       json.RawMessage `json:"terms"`
}

var hexSHA256 = regexp.MustCompile(`^[0-9a-f]{64}$`)

// maxIDLength bounds identifiers, which appear in signed messages.
const maxIDLength = 256

// ErrInvalid wraps every validation failure.
var ErrInvalid = errors.New("agreementsig: invalid agreement version")

func invalid(format string, args ...interface{}) error {
	return fmt.Errorf("%w: %s", ErrInvalid, fmt.Sprintf(format, args...))
}

// Validate checks v's structure: identifiers, the version chain, hashes and
// that Terms is a JSON object with integer-only numbers.
func (v AgreementVersion) Validate() error {
	if err := checkID("agreement id", v.AgreementID); err != nil {
		return err
	}
	if v.Version < 1 || v.Version > maxSafeInteger {
		return invalid("version must be 1 or more")
	}
	switch {
	case v.Version == 1 && v.ParentHash != "":
		return invalid("version 1 has no parent")
	case v.Version > 1 && !hexSHA256.MatchString(v.ParentHash):
		return invalid("version %d needs its parent's hash (64 lowercase hex)", v.Version)
	}
	switch v.Kind {
	case KindAgreement:
		if v.Amends != nil {
			return invalid("an agreement amends nothing")
		}
	case KindAmendment:
		if v.Amends == nil {
			return invalid("an amendment names what it amends")
		}
		if err := checkID("amended agreement id", v.Amends.AgreementID); err != nil {
			return err
		}
		if v.Amends.AgreementID == v.AgreementID {
			return invalid("an amendment has its own id, not the agreement's")
		}
		if len(v.Amends.RatifiedHashes) == 0 {
			return invalid("an amendment names at least the ratified agreement's hash")
		}
		for _, h := range v.Amends.RatifiedHashes {
			if !hexSHA256.MatchString(h) {
				return invalid("ratified hash %q is not 64 lowercase hex", h)
			}
		}
	default:
		return invalid("kind must be %q or %q", KindAgreement, KindAmendment)
	}
	if !hexSHA256.MatchString(v.Document.SHA256) || !hexSHA256.MatchString(v.Document.SourceSHA256) {
		return invalid("document hashes must be 64 lowercase hex")
	}
	if err := checkID("document media type", v.Document.MediaType); err != nil {
		return err
	}
	if err := checkID("source media type", v.Document.SourceMediaType); err != nil {
		return err
	}
	if err := checkID("terms schema", v.TermsSchema); err != nil {
		return err
	}
	terms, err := decodeJSON(v.Terms)
	if err != nil {
		return invalid("terms are not JSON: %v", err)
	}
	if _, ok := terms.(map[string]interface{}); !ok {
		return invalid("terms must be a JSON object")
	}
	return nil
}

// checkID refuses empty, overlong or control-character identifiers: they
// appear verbatim in signed messages, where a newline could forge a line.
func checkID(what, s string) error {
	if s == "" {
		return invalid("%s is empty", what)
	}
	if len(s) > maxIDLength {
		return invalid("%s is longer than %d bytes", what, maxIDLength)
	}
	for _, r := range s {
		if unicode.IsControl(r) || r == ' ' || r == ' ' {
			return invalid("%s contains a control character", what)
		}
	}
	return nil
}

// AgreementCanonical validates v and returns its canonical document. Numbers
// in Terms must be safe integers; decimals (money) are strings.
func AgreementCanonical(v AgreementVersion) ([]byte, error) {
	if err := v.Validate(); err != nil {
		return nil, err
	}
	terms, err := decodeJSON(v.Terms)
	if err != nil {
		return nil, invalid("terms are not JSON: %v", err)
	}
	var parent, amends interface{}
	if v.ParentHash != "" {
		parent = v.ParentHash
	}
	if v.Amends != nil {
		hashes := make([]interface{}, len(v.Amends.RatifiedHashes))
		for i, h := range v.Amends.RatifiedHashes {
			hashes[i] = h
		}
		amends = map[string]interface{}{"agreementId": v.Amends.AgreementID, "ratifiedHashes": hashes}
	}
	doc := map[string]interface{}{
		"schema":      AgreementSchema,
		"agreementId": v.AgreementID,
		"version":     float64(v.Version),
		"parentHash":  parent,
		"kind":        string(v.Kind),
		"amends":      amends,
		"document": map[string]interface{}{
			"sha256":          v.Document.SHA256,
			"mediaType":       v.Document.MediaType,
			"sourceSha256":    v.Document.SourceSHA256,
			"sourceMediaType": v.Document.SourceMediaType,
		},
		"termsSchema": v.TermsSchema,
		"terms":       terms,
	}
	var buf bytes.Buffer
	if err := writeCanonical(&buf, doc, integersOnly); err != nil {
		return nil, invalid("%v", err)
	}
	return buf.Bytes(), nil
}

// AgreementHash is the lowercase hex SHA-256 of v's canonical document.
func AgreementHash(v AgreementVersion) (string, error) {
	doc, err := AgreementCanonical(v)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(doc)
	return hex.EncodeToString(sum[:]), nil
}

// AgreementMessage is the text a signer signs to approve version v, whose
// hash is versionHash. signer is the signer's own ledger identity as their
// wallet or key holder knows it; whoever records the signature checks that
// the signing key controls it. Fixed English, fully determined by its
// arguments.
func AgreementMessage(v AgreementVersion, versionHash, signer string) (string, error) {
	if err := checkID("signer", signer); err != nil {
		return "", err
	}
	if !hexSHA256.MatchString(versionHash) {
		return "", invalid("version hash must be 64 lowercase hex")
	}
	if err := v.Validate(); err != nil {
		return "", err
	}
	what := "agreement " + v.AgreementID
	if v.Kind == KindAmendment {
		what = "amendment " + v.AgreementID + " to agreement " + v.Amends.AgreementID
	}
	return fmt.Sprintf("Tripart\nI, %s, approve version %d of %s.\nVersion hash: sha256:%s", signer, v.Version, what, versionHash), nil
}
