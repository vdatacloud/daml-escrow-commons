package logging

import (
	"errors"
	"io"
	"testing"

	"github.com/vdatacloud/daml-escrow-commons/auditref"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// benchLine is TD-4's measured line: six fields, a person among them.
func benchLine(l *zap.Logger) {
	l.Info("escrow funded",
		zap.String("escrowId", "00a1b2c3d4e5f60718293a4b5c6d7e8f9"),
		zap.String("depositorParty", "u-joey-depositor-devlocal::1220abcdef0123456789abcdef0123456789abcdef0123456789abcdef01234567"),
		zap.String("status", "FUNDED"),
		zap.Int64("amount", 250000),
		zap.Strings("assets", []string{"USDCx", "CC"}),
		zap.Error(errors.New("ledger refused: CONTRACT_NOT_FOUND")))
}

func benchCore() zapcore.Core {
	enc := zap.NewProductionEncoderConfig()
	return zapcore.NewCore(zapcore.NewJSONEncoder(enc), zapcore.AddSync(io.Discard), zapcore.InfoLevel)
}

func benchRefs(tb testing.TB) *auditref.Refs {
	refs, err := auditref.New(auditref.Config{Environment: "production", Keys: map[string][]byte{"r1": []byte("bench-key")}, Current: "r1"})
	if err != nil {
		tb.Fatal(err)
	}
	return refs
}

func BenchmarkRawZap(b *testing.B) {
	l := zap.New(benchCore())
	b.ReportAllocs()
	for b.Loop() {
		benchLine(l)
	}
}

func BenchmarkShapingCore(b *testing.B) {
	l := zap.New(NewRefScrubCore(benchCore(), benchRefs(b)))
	b.ReportAllocs()
	for b.Loop() {
		benchLine(l)
	}
}
