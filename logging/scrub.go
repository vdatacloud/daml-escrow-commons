package logging

import (
	"fmt"
	"strconv"

	"github.com/vdatacloud/daml-escrow-commons/redact"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewScrubCore wraps c so every entry is scrubbed before it is written
// (logging rule L6): a field whose key marks it (redact.ForKey) is
// shortened, masked or replaced, any card number in a string field or the
// message is masked, and a field of a marked key that isn't a string is
// replaced whole. It is the safety net under the call sites, which should
// already log identities and short forms.
//
// Objects and arrays (zap.Object, zap.Array, zap.Any of a struct) can't be
// looked into; under a marked key they are replaced, under any other key
// they pass as logged -- don't log a struct that holds personal data.
func NewScrubCore(c zapcore.Core) zapcore.Core { return &scrubCore{c} }

type scrubCore struct{ zapcore.Core }

func (c *scrubCore) With(fields []zapcore.Field) zapcore.Core {
	return &scrubCore{c.Core.With(scrubFields(fields))}
}

func (c *scrubCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(ent.Level) {
		return ce.AddCore(ent, c)
	}
	return ce
}

func (c *scrubCore) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	ent.Message = redact.MaskPANs(ent.Message)
	return c.Core.Write(ent, scrubFields(fields))
}

func scrubFields(fields []zapcore.Field) []zapcore.Field {
	out := make([]zapcore.Field, len(fields))
	for i, f := range fields {
		out[i] = scrubField(f)
	}
	return out
}

func scrubField(f zapcore.Field) zapcore.Field {
	kind := redact.ForKey(f.Key)
	switch f.Type {
	case zapcore.StringType:
		return zap.String(f.Key, redact.Value(f.Key, f.String))
	case zapcore.ByteStringType:
		if b, ok := f.Interface.([]byte); ok {
			return zap.String(f.Key, redact.Value(f.Key, string(b)))
		}
	case zapcore.StringerType:
		if s, ok := f.Interface.(fmt.Stringer); ok {
			return zap.String(f.Key, redact.Value(f.Key, s.String()))
		}
	case zapcore.ErrorType:
		if err, ok := f.Interface.(error); ok && err != nil {
			if msg := err.Error(); kind != redact.Plain || redact.MaskPANs(msg) != msg {
				return zap.String(f.Key, redact.Value(f.Key, msg))
			}
		}
		return f
	case zapcore.Int64Type, zapcore.Int32Type, zapcore.Int16Type, zapcore.Int8Type:
		return scrubNumber(f, kind, strconv.FormatInt(f.Integer, 10))
	case zapcore.Uint64Type, zapcore.Uint32Type, zapcore.Uint16Type, zapcore.Uint8Type, zapcore.UintptrType:
		return scrubNumber(f, kind, strconv.FormatUint(uint64(f.Integer), 10))
	case zapcore.SkipType, zapcore.NamespaceType:
		return f
	}
	if kind == redact.Plain {
		return f
	}
	return zap.String(f.Key, redact.Redacted)
}

func scrubNumber(f zapcore.Field, kind redact.Kind, digits string) zapcore.Field {
	switch kind {
	case redact.AccountNumber:
		return zap.String(f.Key, redact.Account(digits))
	case redact.Sensitive:
		return zap.String(f.Key, redact.Redacted)
	}
	return f
}
