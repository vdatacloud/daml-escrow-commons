package apierror

import (
	"sort"
	"strconv"
	"strings"

	"github.com/vdatacloud/daml-escrow-commons/cantonid"
)

// ShortID is the one short form for long identifiers in messages, hints
// and logs -- cantonid.Short: a party id keeps its hint and shortens its
// namespace ("relaytest::1220ebb7…4288"); a fingerprint or hash keeps its
// first 8 and last 4 characters ("1220ebb7…4288"); any other long value
// (a signature, an id echoed from input) the same. Typed values should use
// their own Short methods (cantonid.PartyID.Short, ...). Full values
// belong in Error.Details.
func ShortID(s string) string { return cantonid.Short(s) }

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
