// Package auditref names people in logs and traces without naming them
// (daml-escrow PLAN.md Phase 96.6): a keyed, versioned reference that every
// service computes the same way, so one person's activity joins across
// services, and that only the identity service resolves back.
//
// A reference is "ref:<version>:<kind>:<16 hex>", the HMAC-SHA256 of a
// namespaced id under that version's key. The kinds are the identifiers a
// service holds for a person: the OIDC subject (s), the ledger user (u) and
// the party (p). The ledger user id and a party's hint derive from the email
// address, so neither may appear in a production log as it is.
//
// An institution (i) is an organization, not a person: logs carry its id in
// the clear, so they group by institution, unless Config.InstitutionRefs
// asks for references too (a deployment whose logs leave the platform).
//
// In development (Config.Environment dev, local or test) a Refs shows the
// identifiers in the clear, so a developer reads who did what. Anywhere else
// it computes references, and New fails without a key: the switch fails
// closed, never to plain.
package auditref

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"strings"
)

// Kind is the identifier a reference is computed from.
type Kind string

const (
	Subject Kind = "s" // OIDC subject
	User    Kind = "u" // ledger user id
	Party   Kind = "p" // Canton party id
	// Institution is an institution's id (its directory name or tenant).
	Institution Kind = "i"
)

// Prefix starts every reference.
const Prefix = "ref:"

// hexLen is a reference's digest length: 64 bits, collision-free for any
// realistic directory.
const hexLen = 16

// Config configures a Refs; the caller loads it.
type Config struct {
	// Environment is the stack's ENVIRONMENT. dev, local and test log
	// identifiers in the clear; anything else logs references.
	Environment string
	// Keys are the HMAC keys by version ("r1", "r2", ...). Current names the
	// one new references use; the others stay for resolution.
	Keys    map[string][]byte
	Current string
	// InstitutionRefs logs institutions by reference too; off, they are
	// logged in the clear in every environment.
	InstitutionRefs bool
}

// Refs computes the references a service logs.
type Refs struct {
	plain           bool
	institutionRefs bool
	keys            map[string][]byte
	current         string
}

// Plain reports whether environment logs identifiers in the clear.
func Plain(environment string) bool {
	switch strings.ToLower(strings.TrimSpace(environment)) {
	case "dev", "development", "local", "test":
		return true
	}
	return false
}

// New is the Refs for cfg. Outside development it needs the current key.
func New(cfg Config) (*Refs, error) {
	r := &Refs{plain: Plain(cfg.Environment), institutionRefs: cfg.InstitutionRefs, keys: map[string][]byte{}, current: cfg.Current}
	for v, k := range cfg.Keys {
		if v == "" || strings.Contains(v, ":") {
			return nil, fmt.Errorf("auditref: key version %q: want a name without ':'", v)
		}
		if len(k) > 0 {
			r.keys[v] = k
		}
	}
	if len(r.keys[r.current]) == 0 {
		if !r.plain || r.institutionRefs {
			return nil, fmt.Errorf("auditref: environment %q logs references and needs the current key (version %q)", cfg.Environment, cfg.Current)
		}
		r.current = ""
	}
	return r, nil
}

// Plain reports whether r logs identifiers in the clear.
func (r *Refs) Plain() bool { return r.plain }

// Of is what a log carries for id of kind: in development id itself,
// anywhere else its reference under the current key. An institution is
// carried in the clear unless InstitutionRefs, in every environment. ""
// stays "".
func (r *Refs) Of(kind Kind, id string) string {
	if id == "" {
		return ""
	}
	if kind == Institution {
		if !r.institutionRefs {
			return id
		}
	} else if r.plain {
		return id
	}
	return Compute(r.current, r.keys[r.current], kind, id)
}

// Subject, User and Party are Of for each kind.
func (r *Refs) Subject(sub string) string   { return r.Of(Subject, sub) }
func (r *Refs) User(userID string) string   { return r.Of(User, userID) }
func (r *Refs) Party(partyID string) string { return r.Of(Party, partyID) }

// Institution is Of(Institution, id).
func (r *Refs) Institution(id string) string { return r.Of(Institution, id) }

// All is id's reference under every key version, newest-named first -- what
// the identity service stores so a reference from any version resolves.
// It works in development too, given keys.
func (r *Refs) All(kind Kind, id string) []string {
	if id == "" {
		return nil
	}
	versions := make([]string, 0, len(r.keys))
	for v := range r.keys {
		versions = append(versions, v)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(versions)))
	out := make([]string, 0, len(versions))
	for _, v := range versions {
		out = append(out, Compute(v, r.keys[v], kind, id))
	}
	return out
}

// Compute is id's reference of kind under key version: the HMAC-SHA256 of
// "<kind>:<id>", so a party and a user with the same text never collide.
func Compute(version string, key []byte, kind Kind, id string) string {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(string(kind) + ":" + id))
	return Prefix + version + ":" + string(kind) + ":" + hex.EncodeToString(mac.Sum(nil))[:hexLen]
}

// ErrNotRef: the value isn't a reference.
var ErrNotRef = errors.New("auditref: not a reference")

// Parse splits a reference into its version and kind.
func Parse(ref string) (version string, kind Kind, err error) {
	rest, ok := strings.CutPrefix(ref, Prefix)
	if !ok {
		return "", "", ErrNotRef
	}
	parts := strings.Split(rest, ":")
	if len(parts) != 3 || parts[0] == "" || len(parts[2]) != hexLen {
		return "", "", ErrNotRef
	}
	switch k := Kind(parts[1]); k {
	case Subject, User, Party, Institution:
		return parts[0], k, nil
	}
	return "", "", ErrNotRef
}
