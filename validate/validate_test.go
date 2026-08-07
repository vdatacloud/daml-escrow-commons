package validate

import "testing"

func TestRequireNonEmpty(t *testing.T) {
	if err := RequireNonEmpty("name", "bob"); err != nil {
		t.Errorf("expected no error for non-empty value, got %v", err)
	}
	if err := RequireNonEmpty("name", "   "); err == nil {
		t.Error("expected an error for whitespace-only value")
	}
}

func TestRequirePositive(t *testing.T) {
	if err := RequirePositive("amount", 10.5); err != nil {
		t.Errorf("expected no error for positive value, got %v", err)
	}
	if err := RequirePositive("amount", 0); err == nil {
		t.Error("expected an error for zero")
	}
	if err := RequirePositive("amount", -1); err == nil {
		t.Error("expected an error for negative value")
	}
}

func TestRequireOneOf(t *testing.T) {
	if err := RequireOneOf("currency", "USD", "USD", "EUR", "GBP"); err != nil {
		t.Errorf("expected no error for allowed value, got %v", err)
	}
	if err := RequireOneOf("currency", "JPY", "USD", "EUR", "GBP"); err == nil {
		t.Error("expected an error for disallowed value")
	}
}

func TestRequireValidEmail(t *testing.T) {
	valid := []string{"a@b.com", "user.name+tag@sub.example.co"}
	for _, v := range valid {
		if err := RequireValidEmail("email", v); err != nil {
			t.Errorf("expected %q to be valid, got %v", v, err)
		}
	}

	invalid := []string{"", "no-at-sign", "@nodomain", "user@", "user@nodot"}
	for _, v := range invalid {
		if err := RequireValidEmail("email", v); err == nil {
			t.Errorf("expected %q to be invalid", v)
		}
	}
}

func TestErrors_AggregatesAndReports(t *testing.T) {
	var errs Errors
	errs.Add(RequireNonEmpty("name", ""))
	errs.Add(RequirePositive("amount", -5))
	errs.Add(RequireNonEmpty("currency", "USD")) // passes, adds nothing

	err := errs.ErrIfAny()
	if err == nil {
		t.Fatal("expected an aggregated error")
	}
	if len(errs) != 2 {
		t.Fatalf("expected 2 failures, got %d: %v", len(errs), errs)
	}
}

func TestErrors_ErrIfAny_NilWhenEmpty(t *testing.T) {
	var errs Errors
	errs.Add(RequireNonEmpty("name", "bob"))
	if err := errs.ErrIfAny(); err != nil {
		t.Errorf("expected nil error when no failures, got %v", err)
	}
}
