package logging

import (
	"github.com/vdatacloud/daml-escrow-commons/auditref"
	"github.com/vdatacloud/daml-escrow-commons/redact"
	"go.uber.org/zap"
)

// Typed fields (daml-escrow PLAN.md TD-4) say what a value is, so the
// shaping core references or masks it whatever its key: the key rules
// (personKind, redact.ForKey) stay as the safety net for plain fields.

// personID is a person's identifier as a field value. Its String is the
// masked identifier, so a core that doesn't shape still never writes an
// email or a full id.
type personID struct {
	kind auditref.Kind
	id   string
}

func (p personID) String() string { return redact.Text(p.id) }

// accountNumber is an account, routing, IBAN or card number as a field
// value; its String is the last four.
type accountNumber string

func (a accountNumber) String() string { return redact.Account(string(a)) }

// User is a field naming a person by their ledger user id.
func User(key, userID string) zap.Field { return zap.Stringer(key, personID{auditref.User, userID}) }

// Party is a field naming a person by their party id.
func Party(key, partyID string) zap.Field {
	return zap.Stringer(key, personID{auditref.Party, partyID})
}

// Subject is a field naming a person by their OIDC subject.
func Subject(key, sub string) zap.Field { return zap.Stringer(key, personID{auditref.Subject, sub}) }

// Account is a field holding an account, routing, IBAN or card number; it
// is logged as its last four.
func Account(key, number string) zap.Field { return zap.Stringer(key, accountNumber(number)) }
