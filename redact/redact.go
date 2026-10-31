// Package redact shortens and masks values before they reach a log line or
// a trace attribute (daml-escrow PLAN.md Phase 96, logging rules L2-L6).
// Logs are part of the platform's audit record, so they carry identities,
// never names, and never more of an identifier or an account number than a
// reader needs to recognize it:
//
//   - ids -- party and contract ids, hashes, signatures, wallet addresses --
//     take the one short form, cantonid.Short ("relaytest::1220ebb7…4288",
//     "0x529084…9ee7");
//   - account, routing, IBAN and card numbers show their last four only
//     ("…6789"), as PCI DSS masking does;
//   - secrets and personal data (names, emails, phones, addresses) are
//     replaced outright.
//
// Short and masked forms are for humans only: lossy, never parsed back.
package redact

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/vdatacloud/daml-escrow-commons/cantonid"
)

// Redacted replaces a value that must not appear at all.
const Redacted = "[redacted]"

// Ellipsis joins the kept parts of a shortened value.
const Ellipsis = "…"

// ID is an identifier's short form (cantonid.Short): a party id keeps its
// hint, any other long value its first 8 and last 4 characters.
func ID(s string) string { return cantonid.Short(s) }

// IDs applies ID to each id.
func IDs(ids []string) []string {
	out := make([]string, len(ids))
	for i, id := range ids {
		out[i] = ID(id)
	}
	return out
}

// Account masks an account, routing, IBAN or card number to its last four
// characters ("…6789"), ignoring spaces and dashes. A value of four
// characters or fewer is masked entirely: its last four would be all of it.
func Account(s string) string {
	var kept []rune
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			kept = append(kept, r)
		}
	}
	if len(kept) == 0 {
		return ""
	}
	if len(kept) <= 4 {
		return Ellipsis
	}
	return Ellipsis + string(kept[len(kept)-4:])
}

// Secret replaces a secret or a piece of personal data; "" stays "" so an
// absent value still reads as absent.
func Secret(s string) string {
	if s == "" {
		return ""
	}
	return Redacted
}

// panRe finds 13-19 digit runs, optionally grouped by single spaces or
// dashes -- the shapes a card number takes in free text.
var panRe = regexp.MustCompile(`\b\d(?:[ -]?\d){12,18}\b`)

// MaskPANs masks every Luhn-valid 13-19 digit run in s to its last four, the
// safety net for a card number that reached free text (a message, an error).
func MaskPANs(s string) string {
	if !strings.ContainsAny(s, "0123456789") {
		return s
	}
	return panRe.ReplaceAllStringFunc(s, func(m string) string {
		digits := strings.Map(func(r rune) rune {
			if r >= '0' && r <= '9' {
				return r
			}
			return -1
		}, m)
		if len(digits) < 13 || len(digits) > 19 || !luhn(digits) {
			return m
		}
		return Account(digits)
	})
}

func luhn(digits string) bool {
	sum, double := 0, false
	for i := len(digits) - 1; i >= 0; i-- {
		d := int(digits[i] - '0')
		if double {
			d *= 2
			if d > 9 {
				d -= 9
			}
		}
		sum += d
		double = !double
	}
	return sum%10 == 0
}

// emailRe finds email addresses in free text.
var emailRe = regexp.MustCompile(`[A-Za-z0-9._%+\-]+@[A-Za-z0-9.\-]+\.[A-Za-z]{2,}`)

// MaskEmails replaces every email address in s: an email is personal data
// (rule L2), wherever it turns up -- a dev user's subject, an error's text.
func MaskEmails(s string) string {
	if !strings.Contains(s, "@") {
		return s
	}
	return emailRe.ReplaceAllString(s, Redacted)
}

// idRe finds the ids free text carries in full -- a party id ("hint::" and
// its fingerprint), a contract id or hash (44+ hex digits: longer than a
// trace id's 32, so log correlation ids stay whole) and an EVM address.
var idRe = regexp.MustCompile(`(?:[A-Za-z0-9_.\-]+::)?\b[0-9a-f]{44,}\b|\b0x[0-9A-Fa-f]{40}\b`)

// ShortenIDs gives every id in s ID's short form (rule L3) -- the safety
// net for an error message that quotes a ledger refusal's contract and
// party ids.
func ShortenIDs(s string) string {
	if len(s) < 42 {
		return s
	}
	return idRe.ReplaceAllStringFunc(s, ID)
}

// Text is free text fit for a log: card numbers masked to their last four,
// email addresses replaced and ids shortened.
func Text(s string) string { return ShortenIDs(MaskEmails(MaskPANs(s))) }

// Path is a request path fit for a log: a segment that looks like an email
// is redacted, and a long one (an id) takes ID's short form. Query strings
// never belong in a log; Path drops one if present.
func Path(p string) string {
	if i := strings.IndexByte(p, '?'); i >= 0 {
		p = p[:i]
	}
	segs := strings.Split(p, "/")
	for i, s := range segs {
		switch {
		case strings.Contains(s, "@") || strings.Contains(s, "%40"):
			segs[i] = Redacted
		case len(s) > 16:
			segs[i] = ID(s)
		}
	}
	return strings.Join(segs, "/")
}
