package redact

import "strings"

// Kind is how a field's key says its value must be treated.
type Kind int

const (
	// Plain: nothing in the key marks it; free text is still checked for
	// card numbers (MaskPANs).
	Plain Kind = iota
	// Identifier: shortened to ID's form (wallet and payout addresses).
	Identifier
	// AccountNumber: masked to its last four (Account).
	AccountNumber
	// Sensitive: replaced outright (Secret) -- secrets and personal data.
	Sensitive
)

// normalize lowercases key and drops everything but letters and digits, so
// "account_number", "accountNumber" and "payout.account-number" compare
// alike.
func normalize(key string) string {
	return strings.Map(func(r rune) rune {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			return r
		case r >= 'A' && r <= 'Z':
			return r + ('a' - 'A')
		}
		return -1
	}, key)
}

// Key suffixes by kind, normalized. Identifier is checked before Sensitive,
// so "walletAddress" is shortened rather than redacted as an "address".
var (
	identifierSuffixes = []string{"walletaddress", "payoutaddress", "wallet", "walletid"}
	accountSuffixes    = []string{"accountnumber", "acctnumber", "accountno", "routingnumber", "abanumber", "iban",
		"sortcode", "cardnumber", "primaryaccountnumber", "bankaccount"}
	sensitiveSuffixes = []string{
		// secrets and credentials
		"password", "passwd", "passphrase", "secret", "token", "apikey", "privatekey", "authorization", "cookie",
		"credential", "credentials", "cvv", "cvv2", "cvc", "securitycode", "cardexpiry", "expirydate", "trackdata",
		// personal data
		"email", "emailaddress", "phone", "phonenumber", "mobile", "firstname", "lastname", "fullname", "givenname",
		"familyname", "displayname", "personname", "address", "postaladdress", "streetaddress", "dateofbirth",
		"birthdate", "dob", "ssn", "taxid", "nationalid", "passport", "passportnumber",
	}
	// Whole keys (not suffixes): too short or too common to match as one.
	sensitiveKeys = map[string]bool{"name": true, "pin": true, "sub": true, "subject": true, "oktasub": true}
	accountKeys   = map[string]bool{"account": true, "routing": true, "aba": true, "pan": true}
)

// ForKey is the Kind a field's key gives its value.
func ForKey(key string) Kind {
	k := normalize(key)
	if k == "" {
		return Plain
	}
	for _, s := range identifierSuffixes {
		if strings.HasSuffix(k, s) {
			return Identifier
		}
	}
	if accountKeys[k] {
		return AccountNumber
	}
	for _, s := range accountSuffixes {
		if strings.HasSuffix(k, s) {
			return AccountNumber
		}
	}
	if sensitiveKeys[k] {
		return Sensitive
	}
	for _, s := range sensitiveSuffixes {
		if strings.HasSuffix(k, s) {
			return Sensitive
		}
	}
	return Plain
}

// Value is value as a log may carry it under key: shortened, masked or
// replaced as ForKey says, and free text checked for card numbers.
func Value(key, value string) string {
	switch ForKey(key) {
	case Identifier:
		return ID(value)
	case AccountNumber:
		return Account(value)
	case Sensitive:
		return Secret(value)
	}
	return MaskPANs(value)
}
