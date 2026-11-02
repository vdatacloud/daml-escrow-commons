package logging

import (
	"bytes"
	"strings"
	"testing"

	"github.com/vdatacloud/daml-escrow-commons/auditref"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

// TestTypedFields (TD-4): a typed field is referenced or masked by its
// type, whatever its key -- even a key no rule recognizes.
func TestTypedFields(t *testing.T) {
	l, logs, refs := refLogger(t, "production")
	l.Info("acted",
		User("approver", "joey@depositor.devlocal"),
		Party("counterparty", "u-joey::1220aa"),
		Subject("who", "00u1joey"),
		Account("payout", "021000021123456789"))
	m := fieldMap(logs.All()[0])
	if m["approver"] != refs.User("joey@depositor.devlocal") || m["counterparty"] != refs.Party("u-joey::1220aa") ||
		m["who"] != refs.Subject("00u1joey") || m["payout"] != "…6789" {
		t.Errorf("fields = %v", m)
	}

	// Without the shaping core, a typed field still writes no email or full
	// number.
	var buf bytes.Buffer
	raw := zap.New(zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()), zapcore.AddSync(&buf), zapcore.InfoLevel))
	raw.Info("acted", User("approver", "joey@depositor.devlocal"), Account("payout", "021000021123456789"))
	if out := buf.String(); strings.Contains(out, "joey@") || strings.Contains(out, "021000021123456789") {
		t.Errorf("unshaped output = %s", out)
	}
}

// TestServiceRules (TD-4): a service's own rules run before the defaults;
// a field it handles is final, one it passes still gets the defaults.
func TestServiceRules(t *testing.T) {
	refs, err := auditref.New(auditref.Config{Environment: "production", Keys: map[string][]byte{"r1": []byte("k")}, Current: "r1"})
	if err != nil {
		t.Fatal(err)
	}
	upper := RuleFunc(func(f zapcore.Field) (zapcore.Field, bool) {
		if f.Key == "code" && f.Type == zapcore.StringType {
			return zap.String(f.Key, strings.ToUpper(f.String)), true
		}
		return f, false
	})
	core, logs := observer.New(zapcore.InfoLevel)
	l := zap.New(NewShapingCore(core, append([]Rule{upper}, DefaultRules(refs)...)...))
	l.Info("x", zap.String("code", "abc"), zap.String("userID", "joey@depositor.devlocal"), zap.String("note", "mail joey@x.io"))
	m := fieldMap(logs.All()[0])
	if m["code"] != "ABC" || m["userID"] != refs.User("joey@depositor.devlocal") || m["note"] != "mail [redacted]" {
		t.Errorf("fields = %v", m)
	}
}

// TestShapeLeavesCallerSlice: shaping never writes the caller's fields.
func TestShapeLeavesCallerSlice(t *testing.T) {
	l, _, _ := refLogger(t, "production")
	fields := []zap.Field{zap.String("userID", "joey@depositor.devlocal")}
	l.Info("x", fields...)
	if fields[0].String != "joey@depositor.devlocal" {
		t.Errorf("caller's field = %q", fields[0].String)
	}
}

// TestShapingAllocs guards TD-4's cost target: the benchmark line shapes
// in under 10 allocations (raw zap: 4). BenchmarkShapingCore reports time.
func TestShapingAllocs(t *testing.T) {
	l := zap.New(NewRefScrubCore(benchCore(), benchRefs(t)))
	benchLine(l)
	if n := testing.AllocsPerRun(200, func() { benchLine(l) }); n >= 10 {
		t.Errorf("shaping the benchmark line takes %.0f allocations, want under 10", n)
	}
}
