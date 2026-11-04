// Package notify is how a service asks for an email or SMS to be sent to a
// person (daml-escrow PLAN.md Phase 101). Sending lives in one place, the
// platform service (daml-escrow-platform): it holds the provider credentials,
// the templates, delivery status and retries, and meters every message per
// tenant. Every other service (escrow-api, CMS, identity, billing) names what
// to send through this package's Sender and never talks to a mail or SMS
// provider itself.
//
// A Message carries no addresses and no free text:
//   - the recipient is a person named by identity reference (their OIDC
//     subject or audit reference); the platform service looks the contact up
//     in the identity service at send time and never stores it, so erasure
//     stays in the identity service;
//   - the content is a template id, a locale and parameters; the wording,
//     the tenant's branding and the translation live with the template.
//
// MockSender records messages for tests. It is a stand-in, so NewMockSender
// refuses outside development: outside dev a service sends for real or its
// notices are off.
package notify

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Channel is how a message reaches the person.
type Channel string

const (
	ChannelEmail Channel = "email"
	ChannelSMS   Channel = "sms"
)

// Limits on a message's parameters: enough for names, amounts, links and
// codes, never a document.
const (
	MaxParams          = 32
	MaxParamValueBytes = 512
	MaxKeyBytes        = 128
)

// Recipient names a person, never an address. Exactly one field is set.
// The platform service resolves it through the identity service, so an
// address put here resolves to no one and the message fails; that lookup is
// the rule's enforcement, not a pattern check here (a dev identity
// provider's subject is an email, and a dev audit reference is in the
// clear).
type Recipient struct {
	// Subject is the person's identity-provider subject.
	Subject string `json:"subject,omitempty"`
	// AuditRef is the person's audit reference (auditref, "ref:..."), for a
	// caller that holds only that.
	AuditRef string `json:"auditRef,omitempty"`
}

// Message is one notice to one person.
type Message struct {
	// IdempotencyKey is the caller's key for this notice: sending the same
	// key again is the same notice, never a second one (an outbox relay
	// retries with it).
	IdempotencyKey string `json:"idempotencyKey"`
	// TenantID is the tenant the message is metered to and branded for (a
	// daml-escrow-identity PartySet id, as in metering).
	TenantID  string    `json:"tenantId"`
	Channel   Channel   `json:"channel"`
	Recipient Recipient `json:"recipient"`
	// Template is the message's template id, dotted lowercase and owned by
	// the calling product, e.g. "identity.reenrolment-started".
	Template string `json:"template"`
	// Locale is the person's language (BCP 47, e.g. "en", "fr-CA"); the
	// platform service falls back to the template's default.
	Locale string `json:"locale,omitempty"`
	// Params fill the template's placeholders.
	Params map[string]string `json:"params,omitempty"`
}

// Status is where a message is in delivery.
type Status string

const (
	StatusAccepted  Status = "accepted"  // queued for delivery
	StatusSent      Status = "sent"      // handed to the provider
	StatusDelivered Status = "delivered" // the provider reports delivery
	StatusFailed    Status = "failed"    // given up after retries
)

// Receipt is the platform service's answer to a Send.
type Receipt struct {
	// ID is the platform service's id for the message (the same for a
	// repeated IdempotencyKey).
	ID     string `json:"id"`
	Status Status `json:"status"`
}

// Sender sends notices. Send returns once the message is accepted for
// delivery, not delivered: delivery has its own retries and status.
type Sender interface {
	Send(ctx context.Context, msg Message) (Receipt, error)
}

var (
	// ErrInvalidMessage: the message breaks a rule in Validate. Not
	// retryable.
	ErrInvalidMessage = errors.New("notify: invalid message")
	// ErrUnavailable: the sender couldn't take the message now. Retry with
	// the same IdempotencyKey.
	ErrUnavailable = errors.New("notify: sender unavailable")
	// ErrNotInDev: a stand-in sender asked for outside development.
	ErrNotInDev = errors.New("notify: stand-in sender outside development")
)

var (
	templatePattern = regexp.MustCompile(`^[a-z][a-z0-9]*(\.[a-z0-9][a-z0-9-]*)+$`)
	localePattern   = regexp.MustCompile(`^[A-Za-z]{2,3}(-[A-Za-z0-9]{2,8})*$`)
	keyPattern      = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9_]*$`)
	idemPattern     = regexp.MustCompile(`^[A-Za-z0-9._:-]+$`)
)

// Validate checks the rules every Sender enforces before taking a message.
// The error wraps ErrInvalidMessage and names each broken rule.
func (m Message) Validate() error {
	var problems []string
	add := func(format string, args ...any) { problems = append(problems, fmt.Sprintf(format, args...)) }

	if m.IdempotencyKey == "" {
		add("idempotencyKey is required")
	} else if len(m.IdempotencyKey) > MaxKeyBytes || !idemPattern.MatchString(m.IdempotencyKey) {
		add("idempotencyKey must be at most %d of A-Z a-z 0-9 . _ : -", MaxKeyBytes)
	}
	if strings.TrimSpace(m.TenantID) == "" {
		add("tenantId is required")
	}
	if m.Channel != ChannelEmail && m.Channel != ChannelSMS {
		add("channel must be %q or %q", ChannelEmail, ChannelSMS)
	}
	switch r := m.Recipient; {
	case r.Subject == "" && r.AuditRef == "":
		add("recipient needs a subject or an auditRef")
	case r.Subject != "" && r.AuditRef != "":
		add("recipient takes a subject or an auditRef, not both")
	case len(r.Subject) > MaxKeyBytes*2 || len(r.AuditRef) > MaxKeyBytes:
		add("recipient reference is too long")
	}
	if !templatePattern.MatchString(m.Template) || len(m.Template) > MaxKeyBytes {
		add("template must be a dotted lowercase id such as identity.reenrolment-started")
	}
	if m.Locale != "" && !localePattern.MatchString(m.Locale) {
		add("locale must be a BCP 47 tag such as en or fr-CA")
	}
	if len(m.Params) > MaxParams {
		add("at most %d params", MaxParams)
	}
	keys := make([]string, 0, len(m.Params))
	for k := range m.Params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		v := m.Params[k]
		if len(k) > MaxKeyBytes || !keyPattern.MatchString(k) {
			add("param name %q must start with a letter and use letters, digits and _", k)
		}
		if len(v) > MaxParamValueBytes {
			add("param %q is over %d bytes", k, MaxParamValueBytes)
		}
		if !utf8.ValidString(v) || strings.IndexFunc(v, unicode.IsControl) >= 0 {
			add("param %q has control characters or invalid UTF-8", k)
		}
	}
	if len(problems) > 0 {
		return fmt.Errorf("%w: %s", ErrInvalidMessage, strings.Join(problems, "; "))
	}
	return nil
}

// IsDevelopment reports whether environment is a development one, the same
// set auditref treats as development.
func IsDevelopment(environment string) bool {
	switch strings.ToLower(strings.TrimSpace(environment)) {
	case "dev", "development", "local", "test":
		return true
	}
	return false
}
