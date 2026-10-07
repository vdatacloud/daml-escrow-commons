package agreementsig

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"errors"
	"fmt"
	"math/big"
)

// Algorithm names how a signature was made. Every recorded signature
// carries one, so a signer signs with the key type it holds: Ed25519 for
// wallets, ECDSA P-256 for keys in a cloud KMS or HSM (the custodian's) and,
// later, passkeys.
type Algorithm string

const (
	// Ed25519 signs the message's UTF-8 bytes (RFC 8032). Public key: the
	// raw 32 bytes or a DER SubjectPublicKeyInfo. Signature: 64 bytes.
	Ed25519 Algorithm = "ed25519"
	// ECDSAP256SHA256 signs the SHA-256 digest of the message's UTF-8 bytes
	// on curve P-256. Public key: a DER SubjectPublicKeyInfo. Signature:
	// ASN.1 DER (r, s), strictly encoded, with s in the lower half of the
	// curve order -- so each approval has exactly one valid encoding.
	ECDSAP256SHA256 Algorithm = "ecdsa-p256-sha256"
)

var (
	// ErrBadSignature: the signature isn't the key's over the message.
	ErrBadSignature = errors.New("agreementsig: the signature does not match this message and key")
	// ErrNonCanonicalSignature: a valid ECDSA signature in a form other than
	// the one this spec accepts (high s, or loose DER); NormalizeECDSA fixes
	// it before recording.
	ErrNonCanonicalSignature = errors.New("agreementsig: ECDSA signature is not in canonical low-s DER form")
	// ErrUnsupportedAlgorithm: not a known Algorithm.
	ErrUnsupportedAlgorithm = errors.New("agreementsig: unsupported signature algorithm")
)

// Verify checks that signature is publicKey's signature over message under
// alg.
func Verify(alg Algorithm, publicKey []byte, message string, signature []byte) error {
	switch alg {
	case Ed25519:
		pub, err := parseEd25519(publicKey)
		if err != nil {
			return err
		}
		if len(signature) != ed25519.SignatureSize || !ed25519.Verify(pub, []byte(message), signature) {
			return ErrBadSignature
		}
		return nil
	case ECDSAP256SHA256:
		pub, err := parseP256(publicKey)
		if err != nil {
			return err
		}
		r, s, err := parseCanonicalECDSA(signature)
		if err != nil {
			return err
		}
		digest := sha256.Sum256([]byte(message))
		if !ecdsa.Verify(pub, digest[:], r, s) {
			return ErrBadSignature
		}
		return nil
	default:
		return fmt.Errorf("%w: %q", ErrUnsupportedAlgorithm, alg)
	}
}

// NormalizeECDSA returns a P-256 DER signature in the canonical form Verify
// accepts: strict DER with low s. A KMS may return either half; signers
// pass its output through this before recording it.
func NormalizeECDSA(signature []byte) ([]byte, error) {
	var sig ecdsaSig
	rest, err := asn1.Unmarshal(signature, &sig)
	if err != nil || len(rest) != 0 || sig.R == nil || sig.S == nil || sig.R.Sign() <= 0 || sig.S.Sign() <= 0 {
		return nil, fmt.Errorf("%w: not a DER (r, s) sequence", ErrBadSignature)
	}
	if sig.S.Cmp(p256HalfOrder) > 0 {
		sig.S = new(big.Int).Sub(elliptic.P256().Params().N, sig.S)
	}
	return asn1.Marshal(sig)
}

type ecdsaSig struct{ R, S *big.Int }

var p256HalfOrder = new(big.Int).Rsh(elliptic.P256().Params().N, 1)

func parseCanonicalECDSA(signature []byte) (*big.Int, *big.Int, error) {
	var sig ecdsaSig
	rest, err := asn1.Unmarshal(signature, &sig)
	if err != nil || len(rest) != 0 || sig.R == nil || sig.S == nil || sig.R.Sign() <= 0 || sig.S.Sign() <= 0 {
		return nil, nil, ErrBadSignature
	}
	// Re-encoding must reproduce the input exactly (no loose DER), and s
	// must be the low one of the two valid values.
	der, err := asn1.Marshal(sig)
	if err != nil || !bytes.Equal(der, signature) || sig.S.Cmp(p256HalfOrder) > 0 {
		return nil, nil, ErrNonCanonicalSignature
	}
	return sig.R, sig.S, nil
}

func parseEd25519(b []byte) (ed25519.PublicKey, error) {
	if len(b) == ed25519.PublicKeySize {
		return ed25519.PublicKey(b), nil
	}
	parsed, err := x509.ParsePKIXPublicKey(b)
	if err != nil {
		return nil, fmt.Errorf("agreementsig: not a raw or DER Ed25519 public key: %w", err)
	}
	pub, ok := parsed.(ed25519.PublicKey)
	if !ok {
		return nil, fmt.Errorf("agreementsig: public key is %T, not Ed25519", parsed)
	}
	return pub, nil
}

func parseP256(b []byte) (*ecdsa.PublicKey, error) {
	parsed, err := x509.ParsePKIXPublicKey(b)
	if err != nil {
		return nil, fmt.Errorf("agreementsig: not a DER P-256 public key: %w", err)
	}
	pub, ok := parsed.(*ecdsa.PublicKey)
	if !ok || pub.Curve != elliptic.P256() {
		return nil, fmt.Errorf("agreementsig: public key is not on P-256")
	}
	return pub, nil
}
