package schema

import (
	"errors"
	"reflect"
	"testing"
)

func loadFields(t *testing.T) *Registry {
	t.Helper()
	r, err := LoadDirectory("./testdata/fields")
	if err != nil {
		t.Fatalf("LoadDirectory failed: %v", err)
	}
	return r
}

func TestField_Resolves(t *testing.T) {
	r := loadFields(t)
	cases := []struct {
		pointer  string
		types    []string
		format   string
		enum     []any
		required bool
	}{
		{"/expires", []string{"string"}, "date-time", nil, false},
		{"/amount", []string{"string"}, "", nil, true},                        // through $ref
		{"/parties", []string{"array"}, "", nil, true},                        // through $ref
		{"/parties/*", []string{"string"}, "", nil, false},                    // array items through a nested $ref
		{"/mode", nil, "", []any{"Progressive", "AllOrNone"}, false},          // enum only
		{"/fee/tag", []string{"string"}, "", []any{"NoFee", "FlatFee"}, true}, // sibling narrows the $ref
		{"/milestones/0/label", []string{"string"}, "", nil, true},
		{"/milestones/*/a~1b", []string{"integer"}, "", nil, false}, // escaped segment
		{"/combined/x", []string{"integer"}, "", nil, true},         // allOf merged, required unioned
	}
	for _, c := range cases {
		f, err := r.Field("agreement", c.pointer)
		if err != nil {
			t.Fatalf("%s: %v", c.pointer, err)
		}
		if !reflect.DeepEqual(f.Types, c.types) || f.Format != c.format || !reflect.DeepEqual(f.Enum, c.enum) || f.Required != c.required {
			t.Errorf("%s: got types=%v format=%q enum=%v required=%v", c.pointer, f.Types, f.Format, f.Enum, f.Required)
		}
	}
}

func TestField_KeepsOwnDescriptionAndRefKeywords(t *testing.T) {
	f, err := loadFields(t).Field("agreement", "/amount")
	if err != nil {
		t.Fatal(err)
	}
	if f.Description != "Money as a decimal string." || f.Keywords["pattern"] == nil {
		t.Fatalf("got description %q keywords %v", f.Description, f.Keywords)
	}
	if !f.HasType("string") || f.HasType("number") {
		t.Fatalf("HasType wrong for %v", f.Types)
	}
}

func TestField_Root(t *testing.T) {
	f, err := loadFields(t).Field("agreement", "")
	if err != nil || !f.HasType("object") {
		t.Fatalf("root: %v %v", f, err)
	}
}

func TestField_Errors(t *testing.T) {
	r := loadFields(t)
	var unknownType *ErrUnknownType
	if _, err := r.Field("gadget", "/x"); !errors.As(err, &unknownType) {
		t.Fatalf("want ErrUnknownType, got %v", err)
	}
	var unknownField *ErrUnknownField
	for _, p := range []string{"/nope", "/expires/deeper", "/milestones/x/label"} {
		if _, err := r.Field("agreement", p); !errors.As(err, &unknownField) {
			t.Errorf("%s: want ErrUnknownField, got %v", p, err)
		}
	}
	if _, err := r.Field("agreement", "expires"); err == nil {
		t.Error("a pointer without a leading / must fail")
	}
}

func TestField_RefFailures(t *testing.T) {
	r := &Registry{raw: map[string]map[string]any{"bad": {
		"type": "object",
		"properties": map[string]any{
			"loop":     map[string]any{"$ref": "#/$defs/loop"},
			"external": map[string]any{"$ref": "other.json#/foo"},
			"dangling": map[string]any{"$ref": "#/$defs/missing"},
		},
		"$defs": map[string]any{"loop": map[string]any{"$ref": "#/$defs/loop"}},
	}}}
	for _, p := range []string{"/loop", "/external", "/dangling"} {
		var unknownField *ErrUnknownField
		if _, err := r.Field("bad", p); err == nil || errors.As(err, &unknownField) {
			t.Errorf("%s: want a $ref error, got %v", p, err)
		}
	}
}
