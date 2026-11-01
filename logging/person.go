package logging

import (
	"strings"

	"github.com/vdatacloud/daml-escrow-commons/auditref"
)

// Keys that name a person, normalized (lowercase, letters and digits only),
// and the audit-reference kind each takes (daml-escrow docs/AUDIT_LOGGING.md
// §A.1). Whole keys first; then suffixes, so "depositorParty",
// "custodianUser" and "approverUserId" match but "username" (a person's
// name: redacted) and "userAgent" do not.
var (
	personKeys = map[string]auditref.Kind{
		"sub": auditref.Subject, "subject": auditref.Subject, "oktasub": auditref.Subject,
		"user": auditref.User, "userid": auditref.User, "actor": auditref.User,
		"party": auditref.Party, "partyid": auditref.Party, "parties": auditref.Party, "partyids": auditref.Party,
	}
	personSuffixes = []struct {
		suffix string
		kind   auditref.Kind
	}{
		{"oktasub", auditref.Subject},
		{"userid", auditref.User}, {"user", auditref.User},
		{"partyid", auditref.Party}, {"partyids", auditref.Party}, {"party", auditref.Party}, {"parties", auditref.Party},
	}
)

// personKind is the audit-reference kind key's value takes when it names a
// person: a subject, a ledger user or a party.
func personKind(key string) (auditref.Kind, bool) {
	k := strings.Map(func(r rune) rune {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			return r
		case r >= 'A' && r <= 'Z':
			return r + ('a' - 'A')
		}
		return -1
	}, key)
	if kind, ok := personKeys[k]; ok {
		return kind, true
	}
	for _, s := range personSuffixes {
		if strings.HasSuffix(k, s.suffix) {
			return s.kind, true
		}
	}
	return "", false
}
