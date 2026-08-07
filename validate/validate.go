// Package validate provides small, composable field-level checks for the
// Request-DTO ".Validate() error" convention used across daml-escrow and
// daml-escrow-cms. It deliberately stops at generic primitives (non-empty,
// positive amount, one-of, email shape) — anything domain-specific (currency
// codes, role exclusivity, milestone invariants) stays in the owning
// service's own Validate(), composed out of these.
package validate

import (
	"fmt"
	"strings"
)

// Errors aggregates multiple field failures into a single error, so a
// DTO's Validate() can report every problem in one response instead of
// stopping at the first.
type Errors []string

func (e Errors) Error() string {
	return strings.Join(e, "; ")
}

// Add appends a formatted failure if it isn't empty, and is a no-op
// otherwise — lets callers write `errs.Add(RequireNonEmpty("name", v))`
// without a nil check.
func (e *Errors) Add(err error) {
	if err != nil {
		*e = append(*e, err.Error())
	}
}

// ErrIfAny returns e as an error if it has any entries, or nil otherwise —
// the usual last line of a Validate() method.
func (e Errors) ErrIfAny() error {
	if len(e) == 0 {
		return nil
	}
	return e
}

// RequireNonEmpty fails if value is empty after trimming whitespace.
func RequireNonEmpty(field, value string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("%s is required", field)
	}
	return nil
}

// RequirePositive fails if value is <= 0.
func RequirePositive(field string, value float64) error {
	if value <= 0 {
		return fmt.Errorf("%s must be greater than zero", field)
	}
	return nil
}

// RequireOneOf fails if value is not present in allowed.
func RequireOneOf(field, value string, allowed ...string) error {
	for _, a := range allowed {
		if value == a {
			return nil
		}
	}
	return fmt.Errorf("%s must be one of %s, got %q", field, strings.Join(allowed, ", "), value)
}

// RequireValidEmail fails if value doesn't have the minimal shape of an
// email address (one "@", something on both sides, a "." after the "@").
// This is intentionally not a full RFC 5322 validator — it exists to catch
// obvious input mistakes, not to be the authority on deliverability.
func RequireValidEmail(field, value string) error {
	at := strings.IndexByte(value, '@')
	if at <= 0 || at == len(value)-1 {
		return fmt.Errorf("%s must be a valid email address", field)
	}
	domain := value[at+1:]
	if !strings.Contains(domain, ".") {
		return fmt.Errorf("%s must be a valid email address", field)
	}
	return nil
}
