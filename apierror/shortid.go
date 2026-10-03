package apierror

import (
	"sort"
	"strconv"
	"strings"
)

// ShortID is the one short form for long identifiers in messages, hints
// and logs: an id with a "hint::namespace" shape (a Canton party) keeps its
// hint and shortens its namespace ("relaytest::1220ebb7…4288"); any other
// long value -- hash, fingerprint, key, signature -- keeps its first 8 and
// last 4 characters ("3f9a0c1d…77e2"). Values short enough to read are
// returned unchanged. Full values belong in Error.Details.
func ShortID(s string) string {
	if hint, ns, ok := strings.Cut(s, "::"); ok {
		return hint + "::" + shorten(ns)
	}
	return shorten(s)
}

const shortHead, shortTail = 8, 4

func shorten(s string) string {
	r := []rune(s)
	if len(r) <= shortHead+shortTail+4 {
		return s
	}
	return string(r[:shortHead]) + "…" + string(r[len(r)-shortTail:])
}

// ShortIDs applies ShortID to each id.
func ShortIDs(ids []string) []string {
	out := make([]string, len(ids))
	for i, id := range ids {
		out[i] = ShortID(id)
	}
	return out
}

// Field is one key/value for a structured logger (zap.String, slog.String,
// ...), kept logger-agnostic.
type Field struct{ Key, Value string }

// LogFields are e's fields for a log line, in a fixed order, with every
// details value shortened by ShortID -- logs never carry the full values
// (they are in the response's details). Keys: code, status, stage,
// requestId, upstream.*, details.<key>.
func (e *Error) LogFields() []Field {
	f := []Field{{"code", e.Code}}
	add := func(k, v string) {
		if v != "" {
			f = append(f, Field{k, v})
		}
	}
	if e.Status != 0 {
		add("status", strconv.Itoa(e.Status))
	}
	add("stage", e.Stage)
	add("requestId", e.RequestID)
	if u := e.Upstream; u != nil {
		add("upstream.service", u.Service)
		add("upstream.code", u.Code)
		add("upstream.grpcCode", u.GRPCCode)
		add("upstream.traceId", u.TraceID)
		add("upstream.node", u.Node)
	}
	if e.Details != nil {
		raw, err := sortedJSON(e.Details)
		if err == nil {
			var m map[string]any
			if decodeInto(raw, &m) == nil {
				keys := make([]string, 0, len(m))
				for k := range m {
					keys = append(keys, k)
				}
				sort.Strings(keys)
				for _, k := range keys {
					add("details."+k, shortValue(m[k]))
				}
			}
		}
	}
	return f
}

func shortValue(v any) string {
	switch t := v.(type) {
	case string:
		return ShortID(t)
	case []any:
		parts := make([]string, 0, len(t))
		for _, x := range t {
			parts = append(parts, shortValue(x))
		}
		return strings.Join(parts, ",")
	case nil:
		return ""
	default:
		b, _ := encode(t)
		return ShortID(string(b))
	}
}
