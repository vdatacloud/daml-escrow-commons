package redact

import (
	"strings"
	"testing"
)

const party = "relaytest::1220ebb7f1b0b0e3e1f0f7a2c9d8e4b3a2f1e0d9c8b7a6f5e4d3c2b1a0f9e8d74288"

func TestID(t *testing.T) {
	cases := []struct{ in, want string }{
		{party, "relaytest::1220ebb7…4288"},
		{"0x52908400098527886E0F7030069857D2E4169EE7", "0x529084…9EE7"},
		{"short-id", "short-id"},
		{"", ""},
	}
	for _, c := range cases {
		if got := ID(c.in); got != c.want {
			t.Errorf("ID(%q) = %q, want %q", c.in, got, c.want)
		}
	}
	if got := IDs([]string{party, "x"}); got[0] != "relaytest::1220ebb7…4288" || got[1] != "x" {
		t.Errorf("IDs = %v", got)
	}
}

func TestAccount(t *testing.T) {
	cases := []struct{ in, want string }{
		{"021000021", "…0021"},                   // ABA routing number
		{"000123456789", "…6789"},                // account number
		{"GB82 WEST 1234 5698 7654 32", "…5432"}, // IBAN, grouped
		{"4111-1111-1111-1111", "…1111"},         // card number
		{"1234", "…"},                            // too short to show any of it
		{"", ""},
	}
	for _, c := range cases {
		if got := Account(c.in); got != c.want {
			t.Errorf("Account(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestSecret(t *testing.T) {
	if Secret("hunter2") != Redacted || Secret("") != "" {
		t.Error("Secret replaces a value and keeps an absent one absent")
	}
}

func TestMaskPANs(t *testing.T) {
	cases := []struct{ in, want string }{
		{"card 4111111111111111 declined", "card …1111 declined"},
		{"card 4111 1111 1111 1111 declined", "card …1111 declined"},
		{"card 5500-0000-0000-0004 ok", "card …0004 ok"},
		// Not Luhn-valid, or not a card's length: left alone.
		{"order 4111111111111112", "order 4111111111111112"},
		{"offset 123456789012", "offset 123456789012"},
		{"no digits here", "no digits here"},
	}
	for _, c := range cases {
		if got := MaskPANs(c.in); got != c.want {
			t.Errorf("MaskPANs(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestPath(t *testing.T) {
	cid := "00b34d9a70b29cc6404ce837603ecdd1e084ce51b65c5b8e6f9c28ba47754a2f7d"
	cases := []struct{ in, want string }{
		{"/api/v1/escrows/" + cid + "/fund", "/api/v1/escrows/00b34d9a…2f7d/fund"},
		{"/api/v1/identities/joe@example.com", "/api/v1/identities/" + Redacted},
		{"/api/v1/identities/joe%40example.com", "/api/v1/identities/" + Redacted},
		{"/api/v1/escrows?counterparty=joe@example.com", "/api/v1/escrows"},
		{"/health", "/health"},
	}
	for _, c := range cases {
		if got := Path(c.in); got != c.want {
			t.Errorf("Path(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestForKey(t *testing.T) {
	cases := map[string]Kind{
		"walletAddress": Identifier, "custody_payout_address": Identifier, "wallet": Identifier,
		"account_number": AccountNumber, "accountNumber": AccountNumber, "routing_number": AccountNumber,
		"iban": AccountNumber, "payout.iban": AccountNumber, "cardNumber": AccountNumber, "pan": AccountNumber,
		"account": AccountNumber,
		"email":   Sensitive, "userEmail": Sensitive, "display_name": Sensitive, "name": Sensitive,
		"physical_address": Sensitive, "phone": Sensitive, "Authorization": Sensitive, "access_token": Sensitive,
		"client-secret": Sensitive, "cvv": Sensitive, "sub": Sensitive, "okta_sub": Sensitive, "pin": Sensitive,
		"requestId": Plain, "escrowId": Plain, "institutionName": Plain, "span": Plain, "span_id": Plain,
		"status": Plain, "": Plain, "accountId": Plain, "company": Plain,
	}
	for key, want := range cases {
		if got := ForKey(key); got != want {
			t.Errorf("ForKey(%q) = %d, want %d", key, got, want)
		}
	}
}

func TestValue(t *testing.T) {
	cases := []struct{ key, in, want string }{
		{"walletAddress", "0x52908400098527886E0F7030069857D2E4169EE7", "0x529084…9EE7"},
		{"routing_number", "021000021", "…0021"},
		{"email", "joe@example.com", Redacted},
		{"error", "card 4111111111111111 declined", "card …1111 declined"},
		{"escrowId", "e-1", "e-1"},
	}
	for _, c := range cases {
		if got := Value(c.key, c.in); got != c.want {
			t.Errorf("Value(%q, %q) = %q, want %q", c.key, c.in, got, c.want)
		}
	}
	if strings.Contains(Value("note", "pay 4111111111111111"), "41111111") {
		t.Error("a card number in free text must not survive")
	}
}

func TestText(t *testing.T) {
	cases := []struct{ in, want string }{
		{"identity joe.bloggs+x@depositor.devlocal not found", "identity " + Redacted + " not found"},
		{"card 4111111111111111 for a@b.co", "card …1111 for " + Redacted},
		{"user@host without a domain dot", "user@host without a domain dot"},
		{"plain", "plain"},
	}
	for _, c := range cases {
		if got := Text(c.in); got != c.want {
			t.Errorf("Text(%q) = %q, want %q", c.in, got, c.want)
		}
	}
	if Value("userID", "joey@depositor.devlocal") != Redacted {
		t.Error("an email under any key is replaced")
	}
}
