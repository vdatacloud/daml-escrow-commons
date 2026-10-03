// Package cantonid holds typed Canton identifiers -- key fingerprints,
// party ids and transaction/topology hashes -- with one parser and one
// short display form each, so every service (and the @vdatacloud/cx-commons
// TypeScript mirror, sdk/canton-id, tested against testdata/parse.json)
// reads and renders them the same way.
//
//	Fingerprint  1220<64 hex>               a key: "1220" (multihash sha-256, 32 bytes) + sha256(purpose || key)
//	PartyID      <hint>::<Fingerprint>      a party: its hint plus its namespace, which IS a fingerprint --
//	                                        the controlling key's for an external party, the hosting
//	                                        participant's (shared by all its parties) otherwise
//	Hash         1220<64 hex>               a transaction or topology hash (same multihash format)
//
// Short forms are for humans only (messages, hints, logs, UI labels with
// the full value on hover): lossy, never parsed back, never an identifier.
// JSON always carries the full value.
package cantonid

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const (
	shortHead = 8 // includes the "1220" multihash prefix
	shortTail = 4
	// signingKeyPurpose is Canton's hash purpose for public-key
	// fingerprints (HashPurpose.PublicKeyFingerprint).
	signingKeyPurpose = 12
)

var multihashRe = regexp.MustCompile(`^1220[0-9a-f]{64}$`)

// hintRe is Canton's party-hint alphabet: letters, digits, '-', '_', ':'
// and ' ' -- but never the "::" delimiter -- up to 185 characters.
var hintRe = regexp.MustCompile(`^[A-Za-z0-9_\- :]{1,185}$`)

func shortHex(s string) string {
	if len(s) <= shortHead+shortTail {
		return s
	}
	return s[:shortHead] + "…" + s[len(s)-shortTail:]
}

// ---- Fingerprint ----------------------------------------------------------

// Fingerprint identifies a public key: "1220" + 64 lowercase hex.
type Fingerprint string

// ParseFingerprint validates s as a fingerprint.
func ParseFingerprint(s string) (Fingerprint, error) {
	if !multihashRe.MatchString(s) {
		return "", fmt.Errorf("cantonid: %q is not a fingerprint (1220 + 64 lowercase hex)", s)
	}
	return Fingerprint(s), nil
}

// FingerprintOf is pub's Canton fingerprint: "1220" + hex(sha256(uint32be(12) || pub)).
func FingerprintOf(pub ed25519.PublicKey) Fingerprint {
	var purpose [4]byte
	binary.BigEndian.PutUint32(purpose[:], signingKeyPurpose)
	sum := sha256.Sum256(append(purpose[:], pub...))
	return Fingerprint("1220" + hex.EncodeToString(sum[:]))
}

func (f Fingerprint) String() string { return string(f) }

// Short is "1220ebb7…4288".
func (f Fingerprint) Short() string { return shortHex(string(f)) }

// IsZero reports an empty fingerprint (for omitzero).
func (f Fingerprint) IsZero() bool { return f == "" }

func (f *Fingerprint) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	if s == "" {
		*f = ""
		return nil
	}
	v, err := ParseFingerprint(s)
	if err != nil {
		return err
	}
	*f = v
	return nil
}

// ---- Hash -----------------------------------------------------------------

// Hash is a transaction or topology hash in Canton's multihash hex form.
type Hash string

// ParseHash validates s as a multihash hex hash.
func ParseHash(s string) (Hash, error) {
	if !multihashRe.MatchString(s) {
		return "", fmt.Errorf("cantonid: %q is not a hash (1220 + 64 lowercase hex)", s)
	}
	return Hash(s), nil
}

func (h Hash) String() string { return string(h) }

// Short is "12207029…89f1".
func (h Hash) Short() string { return shortHex(string(h)) }

func (h Hash) IsZero() bool { return h == "" }

func (h *Hash) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	if s == "" {
		*h = ""
		return nil
	}
	v, err := ParseHash(s)
	if err != nil {
		return err
	}
	*h = v
	return nil
}

// ---- PartyID --------------------------------------------------------------

// PartyID is a Canton party: Hint::Namespace.
type PartyID struct {
	Hint      string
	Namespace Fingerprint
}

// ErrNotPartyID: the value isn't hint::fingerprint.
var ErrNotPartyID = errors.New("cantonid: not a party id (hint::1220<64 hex>)")

// ParsePartyID validates s as hint::fingerprint.
func ParsePartyID(s string) (PartyID, error) {
	hint, ns, ok := strings.Cut(s, "::")
	if !ok || !hintRe.MatchString(hint) || strings.Contains(hint, "::") {
		return PartyID{}, fmt.Errorf("%w: %q", ErrNotPartyID, s)
	}
	fp, err := ParseFingerprint(ns)
	if err != nil {
		return PartyID{}, fmt.Errorf("%w: %q", ErrNotPartyID, s)
	}
	return PartyID{Hint: hint, Namespace: fp}, nil
}

// MustParsePartyID is ParsePartyID for values known valid (tests,
// constants); it panics otherwise.
func MustParsePartyID(s string) PartyID {
	p, err := ParsePartyID(s)
	if err != nil {
		panic(err)
	}
	return p
}

// String is the full id, "hint::1220…" (all 68 hex).
func (p PartyID) String() string {
	if p.IsZero() {
		return ""
	}
	return p.Hint + "::" + string(p.Namespace)
}

// Short is "relaytest::1220ebb7…4288".
func (p PartyID) Short() string {
	if p.IsZero() {
		return ""
	}
	return p.Hint + "::" + p.Namespace.Short()
}

func (p PartyID) IsZero() bool { return p.Hint == "" && p.Namespace == "" }

// ControlledBy reports whether key fp alone controls p -- true exactly for
// an external party whose namespace is fp. (A participant-hosted party's
// namespace is its participant's key, shared by every party it hosts.)
func (p PartyID) ControlledBy(fp Fingerprint) bool {
	return fp != "" && p.Namespace == fp
}

func (p PartyID) MarshalJSON() ([]byte, error) { return json.Marshal(p.String()) }

func (p *PartyID) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	if s == "" {
		*p = PartyID{}
		return nil
	}
	v, err := ParsePartyID(s)
	if err != nil {
		return err
	}
	*p = v
	return nil
}

// ---- untyped values -------------------------------------------------------

// Kind names what Classify recognized.
type Kind string

const (
	KindParty       Kind = "party"
	KindFingerprint Kind = "fingerprint" // also a hash: same format, context decides
	KindOther       Kind = "other"
)

// Classify recognizes s as a party id or a fingerprint/hash.
func Classify(s string) Kind {
	if _, err := ParsePartyID(s); err == nil {
		return KindParty
	}
	if multihashRe.MatchString(s) {
		return KindFingerprint
	}
	return KindOther
}

// Short is the short form of any string: a party or fingerprint/hash by
// its type's rule; anything else long (a signature, a base64 hash, an
// unparseable id echoed from input) keeps its first 8 and last 4
// characters. Prefer the typed Short methods when the type is known.
func Short(s string) string {
	if p, err := ParsePartyID(s); err == nil {
		return p.Short()
	}
	if multihashRe.MatchString(s) {
		return shortHex(s)
	}
	// Unparseable "x::y" (e.g. a malformed id from input): keep the hint,
	// shorten the rest -- still recognizable, never mistaken for valid.
	if hint, rest, ok := strings.Cut(s, "::"); ok {
		return hint + "::" + shortRunes(rest)
	}
	return shortRunes(s)
}

func shortRunes(s string) string {
	r := []rune(s)
	if len(r) <= shortHead+shortTail+4 {
		return s
	}
	return string(r[:shortHead]) + "…" + string(r[len(r)-shortTail:])
}
