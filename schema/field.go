package schema

import (
	"fmt"
	"strconv"
	"strings"
)

// FieldInfo describes the schema of one field, found by JSON pointer. It is
// what a caller needs to bind a value to that field (a CMS clause
// placeholder, a smart-form input): the JSON types it accepts, its format,
// its enumeration, and whether its parent object requires it.
type FieldInfo struct {
	Pointer     string
	Types       []string // the "type" keyword; empty when the schema gives none (an enum-only field)
	Format      string
	Enum        []any
	Required    bool // the parent object lists the field as required; always false for the root and array items
	Description string
	// Keywords is the field's schema with local $refs and allOf merged in, for
	// the keywords FieldInfo doesn't lift out (bounds, patterns, items).
	Keywords map[string]any
}

// HasType reports whether the field accepts JSON type t.
func (f *FieldInfo) HasType(t string) bool {
	for _, ft := range f.Types {
		if ft == t {
			return true
		}
	}
	return false
}

// ErrUnknownField is returned by Field when the pointer names nothing in the
// schema.
type ErrUnknownField struct {
	TypeName string
	Pointer  string
}

func (e *ErrUnknownField) Error() string {
	return fmt.Sprintf("schema: %q has no field %q", e.TypeName, e.Pointer)
}

// maxRefDepth bounds $ref chasing so a cyclic schema fails instead of looping.
const maxRefDepth = 32

// Field resolves a JSON pointer (RFC 6901, e.g. "/terms/expiryDate") to the
// schema of that field in the schema registered under typeName. Object
// members are walked through "properties"; an array's items are reached with
// the segment "*" or any index ("/milestones/*/description"). Local $refs
// ("#/$defs/...") and allOf are merged; a $ref to another document is not
// supported. Returns *ErrUnknownType and *ErrUnknownField like Validate.
func (r *Registry) Field(typeName, pointer string) (*FieldInfo, error) {
	doc, ok := r.raw[typeName]
	if !ok {
		return nil, &ErrUnknownType{TypeName: typeName}
	}

	segments, err := splitPointer(pointer)
	if err != nil {
		return nil, err
	}

	node, err := resolveNode(doc, doc, 0)
	if err != nil {
		return nil, fmt.Errorf("schema: %q: %w", typeName, err)
	}
	required := false
	for _, seg := range segments {
		next, req, ok := child(node, seg)
		if !ok {
			return nil, &ErrUnknownField{TypeName: typeName, Pointer: pointer}
		}
		if node, err = resolveNode(doc, next, 0); err != nil {
			return nil, fmt.Errorf("schema: %q %s: %w", typeName, pointer, err)
		}
		required = req
	}

	info := &FieldInfo{Pointer: pointer, Required: required, Keywords: node}
	switch t := node["type"].(type) {
	case string:
		info.Types = []string{t}
	case []any:
		for _, v := range t {
			if s, ok := v.(string); ok {
				info.Types = append(info.Types, s)
			}
		}
	}
	info.Format, _ = node["format"].(string)
	info.Description, _ = node["description"].(string)
	info.Enum, _ = node["enum"].([]any)
	return info, nil
}

func splitPointer(pointer string) ([]string, error) {
	if pointer == "" || pointer == "/" {
		return nil, nil
	}
	if !strings.HasPrefix(pointer, "/") {
		return nil, fmt.Errorf("schema: pointer %q must start with /", pointer)
	}
	parts := strings.Split(pointer[1:], "/")
	for i, p := range parts {
		parts[i] = strings.ReplaceAll(strings.ReplaceAll(p, "~1", "/"), "~0", "~")
	}
	return parts, nil
}

// child steps from an object or array schema to one member, reporting
// whether the parent requires it.
func child(node map[string]any, seg string) (map[string]any, bool, bool) {
	if props, ok := node["properties"].(map[string]any); ok {
		if c, ok := props[seg].(map[string]any); ok {
			return c, contains(node["required"], seg), true
		}
	}
	if items, ok := node["items"].(map[string]any); ok {
		if _, err := strconv.Atoi(seg); seg == "*" || err == nil {
			return items, false, true
		}
	}
	return nil, false, false
}

// resolveNode merges a node's local $ref and allOf members into one schema;
// the node's own keywords win over the referenced ones, except "properties"
// (merged member by member) and "required" (unioned).
func resolveNode(doc, node map[string]any, depth int) (map[string]any, error) {
	if depth > maxRefDepth {
		return nil, fmt.Errorf("$ref nesting deeper than %d (cyclic?)", maxRefDepth)
	}
	var parts []map[string]any
	if ref, ok := node["$ref"].(string); ok {
		target, err := lookupRef(doc, ref)
		if err != nil {
			return nil, err
		}
		resolved, err := resolveNode(doc, target, depth+1)
		if err != nil {
			return nil, err
		}
		parts = append(parts, resolved)
	}
	if all, ok := node["allOf"].([]any); ok {
		for _, a := range all {
			if m, ok := a.(map[string]any); ok {
				resolved, err := resolveNode(doc, m, depth+1)
				if err != nil {
					return nil, err
				}
				parts = append(parts, resolved)
			}
		}
	}
	if len(parts) == 0 {
		return node, nil
	}

	out := map[string]any{}
	for _, p := range append(parts, node) {
		for k, v := range p {
			switch k {
			case "$ref", "allOf":
			case "properties":
				merged, _ := out["properties"].(map[string]any)
				if merged == nil {
					merged = map[string]any{}
				}
				if pm, ok := v.(map[string]any); ok {
					for pk, pv := range pm {
						merged[pk] = mergeProperty(merged[pk], pv)
					}
				}
				out["properties"] = merged
			case "required":
				out["required"] = union(out["required"], v)
			default:
				out[k] = v
			}
		}
	}
	return out, nil
}

// mergeProperty overlays a property schema on an earlier one, keyword by
// keyword, so {"$ref": variant, "properties": {"tag": {"enum": [...]}}}
// narrows the variant's tag rather than replacing it.
func mergeProperty(base, over any) any {
	b, ok1 := base.(map[string]any)
	o, ok2 := over.(map[string]any)
	if !ok1 || !ok2 {
		return over
	}
	out := map[string]any{}
	for k, v := range b {
		out[k] = v
	}
	for k, v := range o {
		out[k] = v
	}
	return out
}

func lookupRef(doc map[string]any, ref string) (map[string]any, error) {
	if !strings.HasPrefix(ref, "#") {
		return nil, fmt.Errorf("$ref %q: only local references are supported", ref)
	}
	segments, err := splitPointer(strings.TrimPrefix(ref, "#"))
	if err != nil {
		return nil, err
	}
	var cur any = doc
	for _, seg := range segments {
		m, ok := cur.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("$ref %q does not resolve", ref)
		}
		if cur, ok = m[seg]; !ok {
			return nil, fmt.Errorf("$ref %q does not resolve", ref)
		}
	}
	m, ok := cur.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("$ref %q is not a schema object", ref)
	}
	return m, nil
}

func contains(list any, s string) bool {
	l, _ := list.([]any)
	for _, v := range l {
		if v == s {
			return true
		}
	}
	return false
}

func union(a, b any) []any {
	out, _ := a.([]any)
	out = append([]any(nil), out...)
	bl, _ := b.([]any)
	for _, v := range bl {
		if s, ok := v.(string); !ok || !contains(out, s) {
			out = append(out, v)
		}
	}
	return out
}
