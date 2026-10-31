// Package logging builds the platform's audit-grade zap logger (daml-escrow
// PLAN.md Phase 96, logging rules L1-L6): JSON, UTC ISO-8601 time, never
// sampled -- a dropped line is a hole in the audit record -- with service,
// environment and version on every line, and a core that scrubs personal
// and payment data by field key (package redact) before anything is
// written. Request-scoped loggers add the request id and trace ids; the
// access log records each request's outcome and actor.
package logging

import (
	"fmt"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config is a service's logging configuration; the caller loads it.
type Config struct {
	// Level is debug, info, warn or error; "" is info.
	Level string
	// Format is json (the default, and the only one for production) or
	// console (local development).
	Format      string
	Service     string
	Environment string
	Version     string
}

// New builds the logger cfg describes: unsampled, scrubbed, and tagged with
// service, environment and version; it writes to stdout.
func New(cfg Config) (*zap.Logger, error) { return build(cfg, "stdout") }

func build(cfg Config, output string) (*zap.Logger, error) {
	zc, err := zapConfig(cfg, output)
	if err != nil {
		return nil, err
	}
	logger, err := zc.Build(zap.WrapCore(func(c zapcore.Core) zapcore.Core { return NewScrubCore(c) }))
	if err != nil {
		return nil, fmt.Errorf("logging: building the logger: %w", err)
	}
	var fields []zap.Field
	for _, f := range [][2]string{{"service", cfg.Service}, {"environment", cfg.Environment}, {"version", cfg.Version}} {
		if f[1] != "" {
			fields = append(fields, zap.String(f[0], f[1]))
		}
	}
	return logger.With(fields...), nil
}

func zapConfig(cfg Config, output string) (zap.Config, error) {
	level := zap.NewAtomicLevelAt(zapcore.InfoLevel)
	if cfg.Level != "" {
		if err := level.UnmarshalText([]byte(strings.ToLower(cfg.Level))); err != nil {
			return zap.Config{}, fmt.Errorf("logging: level %q: %w", cfg.Level, err)
		}
	}
	encoding := "json"
	switch strings.ToLower(cfg.Format) {
	case "", "json":
	case "console":
		encoding = "console"
	default:
		return zap.Config{}, fmt.Errorf("logging: format %q: want json or console", cfg.Format)
	}
	enc := zap.NewProductionEncoderConfig()
	enc.TimeKey = "ts"
	enc.EncodeTime = func(t time.Time, e zapcore.PrimitiveArrayEncoder) { e.AppendString(t.UTC().Format(time.RFC3339Nano)) }
	return zap.Config{
		Level:            level,
		Encoding:         encoding,
		EncoderConfig:    enc,
		Sampling:         nil, // L1: never sampled
		OutputPaths:      []string{output},
		ErrorOutputPaths: []string{"stderr"},
	}, nil
}
