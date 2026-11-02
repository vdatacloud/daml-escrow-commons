package logging

import (
	"fmt"
	"reflect"
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/vdatacloud/daml-escrow-commons/auditref"
	"github.com/vdatacloud/daml-escrow-commons/redact"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Rule shapes one field on its way to the log (daml-escrow PLAN.md TD-4).
// Shape returns the field to write and true when the rule handled it --
// later rules don't see it -- or false to pass it to the next rule
// unchanged. A rule must be safe for concurrent use.
type Rule interface {
	Shape(f zapcore.Field) (zapcore.Field, bool)
}

// RuleFunc is a function as a Rule.
type RuleFunc func(f zapcore.Field) (zapcore.Field, bool)

// Shape calls fn.
func (fn RuleFunc) Shape(f zapcore.Field) (zapcore.Field, bool) { return fn(f) }

// NewShapingCore wraps c so every field of every entry runs through rules
// in order before it is written, and the message is masked (redact.Text).
// A field no rule handles is written as logged.
func NewShapingCore(c zapcore.Core, rules ...Rule) zapcore.Core {
	return &shapingCore{Core: c, rules: rules}
}

// DefaultRules are the platform's shaping steps (logging rules L2-L6), in
// order: typed fields (User, Party, Subject, Account) by their type; a value
// that already is an audit reference passes; a field whose key names a
// person becomes its reference under refs (none: skipped); a field whose
// key marks it (redact.ForKey) is shortened, masked or replaced; every
// other string has card numbers and emails masked and long ids shortened.
func DefaultRules(refs *auditref.Refs) []Rule {
	rules := []Rule{Typed(refs), PassRefs()}
	if refs != nil {
		rules = append(rules, PersonRefs(refs))
	}
	return append(rules, KeyRedact(), TextMask())
}

// NewScrubCore is NewShapingCore with DefaultRules and no references: the
// safety net under the call sites, which should already log identities and
// short forms.
//
// An array of strings (zap.Strings) is shaped element by element. Objects
// and other arrays (zap.Object, zap.Any of a struct) can't be looked into;
// under a marked key they are replaced, under any other key they pass as
// logged -- don't log a struct that holds personal data.
func NewScrubCore(c zapcore.Core) zapcore.Core { return NewShapingCore(c, DefaultRules(nil)...) }

// NewRefScrubCore is NewScrubCore that also logs each field naming a person
// (by its key, personKind, or its type) as its audit reference under refs
// -- in development the identifier itself, still masked.
func NewRefScrubCore(c zapcore.Core, refs *auditref.Refs) zapcore.Core {
	return NewShapingCore(c, DefaultRules(refs)...)
}

type shapingCore struct {
	zapcore.Core
	rules []Rule
}

func (c *shapingCore) With(fields []zapcore.Field) zapcore.Core {
	return &shapingCore{Core: c.Core.With(c.shape(fields)), rules: c.rules}
}

func (c *shapingCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(ent.Level) {
		return ce.AddCore(ent, c)
	}
	return ce
}

func (c *shapingCore) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	ent.Message = redact.Text(ent.Message)
	return c.Core.Write(ent, c.shape(fields))
}

// shape runs each field through the rules into a new slice: the caller's
// slice may be reused, so it is never written.
func (c *shapingCore) shape(fields []zapcore.Field) []zapcore.Field {
	out := make([]zapcore.Field, len(fields))
	for i, f := range fields {
		out[i] = f
		for _, r := range c.rules {
			if shaped, ok := r.Shape(f); ok {
				out[i] = shaped
				break
			}
		}
	}
	return out
}

// keyInfo is what the rules decide from a key alone, cached per key: keys
// are a small fixed set, so each is normalized once.
type keyInfo struct {
	person   auditref.Kind
	isPerson bool
	redact   redact.Kind
}

var (
	keyCache     sync.Map // string -> keyInfo
	keyCacheSize atomic.Int64
)

// maxCachedKeys bounds the cache in case a caller builds keys at runtime.
const maxCachedKeys = 4096

func keyOf(key string) keyInfo {
	if v, ok := keyCache.Load(key); ok {
		return v.(keyInfo)
	}
	var ki keyInfo
	ki.person, ki.isPerson = personKind(key)
	ki.redact = redact.ForKey(key)
	if keyCacheSize.Load() < maxCachedKeys {
		if _, loaded := keyCache.LoadOrStore(key, ki); !loaded {
			keyCacheSize.Add(1)
		}
	}
	return ki
}

// Typed shapes the fields built with User, Party, Subject and Account by
// their type, whatever their key: a person as their reference under refs
// (without refs, masked), an account number as its last four.
func Typed(refs *auditref.Refs) Rule {
	return RuleFunc(func(f zapcore.Field) (zapcore.Field, bool) {
		if f.Type != zapcore.StringerType {
			return f, false
		}
		switch v := f.Interface.(type) {
		case personID:
			return zap.String(f.Key, refOf(refs, v.kind, v.id)), true
		case accountNumber:
			return zap.String(f.Key, redact.Account(string(v))), true
		}
		return f, false
	})
}

// PassRefs passes a string field that already is an audit reference (an
// access log's actor, a caller's own refs.User): it names no one.
func PassRefs() Rule {
	return RuleFunc(func(f zapcore.Field) (zapcore.Field, bool) {
		if f.Type == zapcore.StringType && isRef(f.String) {
			return f, true
		}
		return f, false
	})
}

func isRef(s string) bool {
	if len(s) <= len(auditref.Prefix) || s[:len(auditref.Prefix)] != auditref.Prefix {
		return false
	}
	_, _, err := auditref.Parse(s)
	return err == nil
}

// PersonRefs logs a string (or string-array) field whose key names a person
// (personKind: oktaSub, userID, depositorParty, actor, ...) as its audit
// reference under refs. Other types under such a key pass to the next rule.
func PersonRefs(refs *auditref.Refs) Rule {
	return RuleFunc(func(f zapcore.Field) (zapcore.Field, bool) {
		ki := keyOf(f.Key)
		if !ki.isPerson {
			return f, false
		}
		switch f.Type {
		case zapcore.StringType:
			return zap.String(f.Key, refOf(refs, ki.person, f.String)), true
		case zapcore.StringerType:
			if s, ok := f.Interface.(fmt.Stringer); ok {
				return zap.String(f.Key, refOf(refs, ki.person, s.String())), true
			}
		case zapcore.ArrayMarshalerType:
			if ss, ok := stringArray(f); ok {
				for i := range ss {
					ss[i] = refOf(refs, ki.person, ss[i])
				}
				return zap.Strings(f.Key, ss), true
			}
		}
		return f, false
	})
}

// refOf is id's audit reference. A reference passes unchanged; an
// identifier logged in the clear (development, or not a person:
// Refs.SetClear) is still masked (emails, long ids).
func refOf(refs *auditref.Refs, kind auditref.Kind, id string) string {
	if isRef(id) {
		return id
	}
	if refs != nil {
		if out := refs.Of(kind, id); out != id {
			return out
		}
	}
	return redact.Text(id)
}

// KeyRedact shortens, masks or replaces a field whose key marks it
// (redact.ForKey: wallets, account numbers, secrets, names, emails); a
// value of a marked key that can't be shaped is replaced whole.
func KeyRedact() Rule {
	return RuleFunc(func(f zapcore.Field) (zapcore.Field, bool) {
		kind := keyOf(f.Key).redact
		if kind == redact.Plain {
			return f, false
		}
		switch f.Type {
		case zapcore.StringType:
			return zap.String(f.Key, redactAs(kind, f.String)), true
		case zapcore.ByteStringType:
			if b, ok := f.Interface.([]byte); ok {
				return zap.String(f.Key, redactAs(kind, string(b))), true
			}
		case zapcore.StringerType:
			if s, ok := f.Interface.(fmt.Stringer); ok {
				return zap.String(f.Key, redactAs(kind, s.String())), true
			}
		case zapcore.ErrorType:
			if err, ok := f.Interface.(error); ok && err != nil {
				return zap.String(f.Key, redactAs(kind, err.Error())), true
			}
			return f, true
		case zapcore.Int64Type, zapcore.Int32Type, zapcore.Int16Type, zapcore.Int8Type:
			return shapeNumber(f, kind, strconv.FormatInt(f.Integer, 10)), true
		case zapcore.Uint64Type, zapcore.Uint32Type, zapcore.Uint16Type, zapcore.Uint8Type, zapcore.UintptrType:
			return shapeNumber(f, kind, strconv.FormatUint(uint64(f.Integer), 10)), true
		case zapcore.ArrayMarshalerType:
			if ss, ok := stringArray(f); ok {
				for i := range ss {
					ss[i] = redactAs(kind, ss[i])
				}
				return zap.Strings(f.Key, ss), true
			}
		case zapcore.SkipType, zapcore.NamespaceType:
			return f, true
		}
		return zap.String(f.Key, redact.Redacted), true
	})
}

func redactAs(kind redact.Kind, s string) string {
	switch kind {
	case redact.Identifier:
		return redact.ID(s)
	case redact.AccountNumber:
		return redact.Account(s)
	case redact.Sensitive:
		return redact.Secret(s)
	}
	return redact.Text(s)
}

func shapeNumber(f zapcore.Field, kind redact.Kind, digits string) zapcore.Field {
	switch kind {
	case redact.AccountNumber:
		return zap.String(f.Key, redact.Account(digits))
	case redact.Sensitive:
		return zap.String(f.Key, redact.Redacted)
	}
	return f
}

// TextMask masks card numbers and emails and shortens long ids in every
// string, string array and error message (redact.Text); other types pass.
// It handles every field, so it ends a rule list.
func TextMask() Rule {
	return RuleFunc(func(f zapcore.Field) (zapcore.Field, bool) {
		switch f.Type {
		case zapcore.StringType:
			return zap.String(f.Key, redact.Text(f.String)), true
		case zapcore.ByteStringType:
			if b, ok := f.Interface.([]byte); ok {
				return zap.String(f.Key, redact.Text(string(b))), true
			}
		case zapcore.StringerType:
			if s, ok := f.Interface.(fmt.Stringer); ok {
				return zap.String(f.Key, redact.Text(s.String())), true
			}
		case zapcore.ErrorType:
			if err, ok := f.Interface.(error); ok && err != nil {
				if msg := err.Error(); redact.Text(msg) != msg {
					return zap.String(f.Key, redact.Text(msg)), true
				}
			}
		case zapcore.ArrayMarshalerType:
			if ss, ok := stringArray(f); ok {
				changed := false
				for i := range ss {
					if t := redact.Text(ss[i]); t != ss[i] {
						ss[i], changed = t, true
					}
				}
				if changed {
					return zap.Strings(f.Key, ss), true
				}
			}
		}
		return f, true
	})
}

// stringArray is a copy of an array field's elements when every one is a
// string: a string slice (zap.Strings) directly, anything else through an
// encoder.
func stringArray(f zapcore.Field) ([]string, bool) {
	if v := reflect.ValueOf(f.Interface); v.Kind() == reflect.Slice && v.Type().Elem().Kind() == reflect.String {
		out := make([]string, v.Len())
		for i := range out {
			out[i] = v.Index(i).String()
		}
		return out, true
	}
	enc := zapcore.NewMapObjectEncoder()
	f.AddTo(enc)
	elems, ok := enc.Fields[f.Key].([]interface{})
	if !ok {
		return nil, false
	}
	out := make([]string, len(elems))
	for i, e := range elems {
		if out[i], ok = e.(string); !ok {
			return nil, false
		}
	}
	return out, true
}
