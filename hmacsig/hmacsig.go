// Package hmacsig provides the HMAC-SHA256 sign/verify primitive shared by
// every webhook-style integration across the platform — today daml-escrow's
// oracle webhook (ORACLE_WEBHOOK_SECRET, internal/services/compliance.go) and
// fiat-settlement webhook, and prospectively daml-escrow-cms's own import/OCR
// callback and third-party-CLM substitution points (INTEGRATION.md). Callers
// own the exact message format they sign (there is no canonical shape,
// deliberately, since payloads differ by feed) — this package only owns the
// keyed-hash primitive so every caller uses constant-time comparison and the
// same hex encoding instead of reimplementing it.
package hmacsig

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

// ErrInvalidSignature is returned by Verify when the signature does not match.
var ErrInvalidSignature = errors.New("hmacsig: invalid signature")

// Sign returns the lowercase-hex HMAC-SHA256 of message under secret.
func Sign(secret, message []byte) string {
	mac := hmac.New(sha256.New, secret)
	mac.Write(message)
	return hex.EncodeToString(mac.Sum(nil))
}

// Verify reports whether hexSignature is the correct HMAC-SHA256 of message
// under secret, using a constant-time comparison. A malformed (non-hex)
// hexSignature is treated as invalid rather than an error, matching how
// callers already handle it — a bad signature and a bad encoding both mean
// "don't trust this payload".
func Verify(secret, message []byte, hexSignature string) error {
	sigBytes, err := hex.DecodeString(hexSignature)
	if err != nil {
		return ErrInvalidSignature
	}

	mac := hmac.New(sha256.New, secret)
	mac.Write(message)
	expected := mac.Sum(nil)

	if !hmac.Equal(sigBytes, expected) {
		return ErrInvalidSignature
	}

	return nil
}
