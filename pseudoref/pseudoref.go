// Package pseudoref derives short, stable, non-reversible references to
// sensitive identifiers (tenant ids, escrow/transaction ids) for operational
// telemetry -- traces, timing records -- where enough of a reference is
// needed to cross-reference a record, but the identifier itself must not
// appear (daml-escrow PLAN.md Phase 68).
//
// A reference is HMAC-SHA256(key, id), hex-encoded and truncated to 12
// characters (48 bits): stable for the same key and id, so records
// correlate; not reversible, and not recomputable without the key -- unlike
// a raw prefix/suffix of the id, which leaks part of it and collides on
// shared prefixes (e.g. a common participant fingerprint). At this
// platform's volumes 48 bits is effectively collision-free.
//
// Only the holder of the key (the platform) can go from a real id to its
// reference to look records up; telemetry consumers see only references.
package pseudoref

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// Length is the number of hex characters in a reference.
const Length = 12

// Ref returns the pseudonymous reference for id under key, or "" when
// either is empty -- callers must never fall back to the raw id.
func Ref(key []byte, id string) string {
	if len(key) == 0 || id == "" {
		return ""
	}
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(id))
	return hex.EncodeToString(mac.Sum(nil))[:Length]
}
