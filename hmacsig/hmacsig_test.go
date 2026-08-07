package hmacsig

import (
	"errors"
	"testing"
)

func TestSignAndVerify_RoundTrip(t *testing.T) {
	secret := []byte("test-secret")
	message := []byte("escrow-123|0|MilestoneVerified|ais")

	sig := Sign(secret, message)
	if err := Verify(secret, message, sig); err != nil {
		t.Fatalf("expected valid signature to verify, got %v", err)
	}
}

func TestVerify_WrongSecretFails(t *testing.T) {
	message := []byte("payload")
	sig := Sign([]byte("secret-a"), message)

	err := Verify([]byte("secret-b"), message, sig)
	if !errors.Is(err, ErrInvalidSignature) {
		t.Fatalf("expected ErrInvalidSignature, got %v", err)
	}
}

func TestVerify_TamperedMessageFails(t *testing.T) {
	secret := []byte("test-secret")
	sig := Sign(secret, []byte("original"))

	err := Verify(secret, []byte("tampered"), sig)
	if !errors.Is(err, ErrInvalidSignature) {
		t.Fatalf("expected ErrInvalidSignature, got %v", err)
	}
}

func TestVerify_MalformedHexFails(t *testing.T) {
	err := Verify([]byte("secret"), []byte("payload"), "not-hex-!!")
	if !errors.Is(err, ErrInvalidSignature) {
		t.Fatalf("expected ErrInvalidSignature for malformed hex, got %v", err)
	}
}
