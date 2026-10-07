package agreementsig

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strconv"
	"unicode/utf16"
)

// maxSafeInteger is the largest integer every JSON implementation in the
// platform (Go's float64, JavaScript's Number) represents exactly.
const maxSafeInteger = 1<<53 - 1

// numberMode is how a canonical document writes numbers.
type numberMode int

const (
	// anyNumber writes any finite number as ECMAScript's JSON.stringify does
	// (draft-version/1, whose drafts carry float amounts).
	anyNumber numberMode = iota
	// integersOnly accepts only integers in the safe range and refuses
	// everything else (agreement-version/1: money is a decimal string, so no
	// implementation has to agree on float formatting).
	integersOnly
)

// decodeJSON parses raw as one JSON value, numbers as json.Number.
func decodeJSON(raw []byte) (interface{}, error) {
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	var v interface{}
	if err := dec.Decode(&v); err != nil {
		return nil, err
	}
	if dec.More() {
		return nil, fmt.Errorf("trailing data after the JSON value")
	}
	return v, nil
}

// writeCanonical writes v as a canonical JSON document: object keys sorted
// by UTF-16 code units (as JavaScript sorts strings), no insignificant
// whitespace, strings and numbers serialized as JSON.stringify does (RFC
// 8785's rules for the values that occur here).
func writeCanonical(buf *bytes.Buffer, v interface{}, mode numberMode) error {
	switch t := v.(type) {
	case nil:
		buf.WriteString("null")
	case bool:
		buf.WriteString(strconv.FormatBool(t))
	case json.Number:
		f, err := strconv.ParseFloat(string(t), 64)
		if err != nil {
			return fmt.Errorf("agreementsig: number %s: %w", t, err)
		}
		return writeNumber(buf, f, mode)
	case float64:
		return writeNumber(buf, t, mode)
	case string:
		writeString(buf, t)
	case []interface{}:
		buf.WriteByte('[')
		for i, e := range t {
			if i > 0 {
				buf.WriteByte(',')
			}
			if err := writeCanonical(buf, e, mode); err != nil {
				return err
			}
		}
		buf.WriteByte(']')
	case map[string]interface{}:
		keys := make([]string, 0, len(t))
		for k := range t {
			keys = append(keys, k)
		}
		sort.Slice(keys, func(i, j int) bool { return lessUTF16(keys[i], keys[j]) })
		buf.WriteByte('{')
		for i, k := range keys {
			if i > 0 {
				buf.WriteByte(',')
			}
			writeString(buf, k)
			buf.WriteByte(':')
			if err := writeCanonical(buf, t[k], mode); err != nil {
				return err
			}
		}
		buf.WriteByte('}')
	default:
		return fmt.Errorf("agreementsig: unsupported value %T", v)
	}
	return nil
}

func writeNumber(buf *bytes.Buffer, f float64, mode numberMode) error {
	if math.IsInf(f, 0) || math.IsNaN(f) {
		return fmt.Errorf("agreementsig: number out of range")
	}
	if mode == integersOnly && (f != math.Trunc(f) || math.Abs(f) > maxSafeInteger) {
		return fmt.Errorf("agreementsig: %v is not a safe integer (write decimals as strings)", f)
	}
	// encoding/json formats float64 as ECMAScript does (shortest round-trip,
	// exponent form outside [1e-6, 1e21)) -- except -0, which
	// JSON.stringify writes as 0.
	if f == 0 {
		f = 0
	}
	b, err := json.Marshal(f)
	if err != nil {
		return err
	}
	buf.Write(b)
	return nil
}

// writeString escapes as JSON.stringify does: the two-character escapes,
// other control characters as \u00xx, everything else literal (no HTML or
// U+2028/2029 escaping, unlike encoding/json).
func writeString(buf *bytes.Buffer, s string) {
	buf.WriteByte('"')
	for _, r := range s {
		switch r {
		case '"':
			buf.WriteString(`\"`)
		case '\\':
			buf.WriteString(`\\`)
		case '\b':
			buf.WriteString(`\b`)
		case '\f':
			buf.WriteString(`\f`)
		case '\n':
			buf.WriteString(`\n`)
		case '\r':
			buf.WriteString(`\r`)
		case '\t':
			buf.WriteString(`\t`)
		default:
			if r < 0x20 {
				fmt.Fprintf(buf, `\u%04x`, r)
			} else {
				buf.WriteRune(r)
			}
		}
	}
	buf.WriteByte('"')
}

func lessUTF16(a, b string) bool {
	ua, ub := utf16.Encode([]rune(a)), utf16.Encode([]rune(b))
	for i := 0; i < len(ua) && i < len(ub); i++ {
		if ua[i] != ub[i] {
			return ua[i] < ub[i]
		}
	}
	return len(ua) < len(ub)
}
