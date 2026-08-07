package schema

import (
	"errors"
	"testing"
)

func TestLoadDirectory_MissingDirectory(t *testing.T) {
	_, err := LoadDirectory("./does-not-exist")
	if err == nil {
		t.Fatal("expected an error for a missing directory")
	}
}

func TestLoadDirectory_CompilesSchemas(t *testing.T) {
	r, err := LoadDirectory("./testdata")
	if err != nil {
		t.Fatalf("LoadDirectory failed: %v", err)
	}

	types := r.Types()
	if len(types) != 1 || types[0] != "widget" {
		t.Fatalf("expected [widget], got %v", types)
	}
}

func TestValidate_UnknownType(t *testing.T) {
	r, err := LoadDirectory("./testdata")
	if err != nil {
		t.Fatalf("LoadDirectory failed: %v", err)
	}

	err = r.Validate("gadget", []byte(`{}`))
	var unknown *ErrUnknownType
	if !errors.As(err, &unknown) {
		t.Fatalf("expected *ErrUnknownType, got %v (%T)", err, err)
	}
}

func TestValidate_ValidPayload(t *testing.T) {
	r, err := LoadDirectory("./testdata")
	if err != nil {
		t.Fatalf("LoadDirectory failed: %v", err)
	}

	err = r.Validate("widget", []byte(`{"name": "bolt", "quantity": 12}`))
	if err != nil {
		t.Fatalf("expected valid payload to pass, got %v", err)
	}
}

func TestValidate_InvalidPayloadReportsFailures(t *testing.T) {
	r, err := LoadDirectory("./testdata")
	if err != nil {
		t.Fatalf("LoadDirectory failed: %v", err)
	}

	err = r.Validate("widget", []byte(`{"quantity": -5}`))
	var valErr *ValidationError
	if !errors.As(err, &valErr) {
		t.Fatalf("expected *ValidationError, got %v (%T)", err, err)
	}
	if valErr.TypeName != "widget" {
		t.Errorf("expected TypeName widget, got %q", valErr.TypeName)
	}
	if len(valErr.Failures) == 0 {
		t.Error("expected at least one failure (missing name, negative quantity)")
	}
}
