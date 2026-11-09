package identityclient

import "time"

// Erasure requests (GDPR Art. 17) as daml-escrow-identity queues them
// (daml-escrow PLAN.md Phase 104.8): the wire shapes of its
// /api/v1/audit/erasure-requests* endpoints, for the platform operator API.
// A request is received (the person found by their email), verified by a
// second operator (four eyes: never the one who received it), then erased
// or refused. The person's email and name are personal data: shown to the
// operator handling the request, never logged; gone once erased.

// Erasure request statuses.
const (
	ErasureReceived = "received"
	ErasureVerified = "verified"
	ErasureErased   = "erased"
	ErasureRefused  = "refused"
)

// ErasureChannels are how a request may arrive.
var ErasureChannels = []string{"email", "web", "letter", "phone", "other"}

// SubmitErasureRequest records a request for the person with Email.
// ReceivedOn is YYYY-MM-DD; Case is the request's own case reference.
type SubmitErasureRequest struct {
	Email      string `json:"email" validate:"required"`
	Channel    string `json:"channel" enums:"email,web,letter,phone,other" validate:"required"`
	ReceivedOn string `json:"receivedOn" example:"2026-11-09" validate:"required"`
	Case       string `json:"case" example:"DSR-2026-031" validate:"required"`
}

// ListErasureRequests lists requests: Open (received or verified) or all.
type ListErasureRequests struct {
	Open bool `json:"open"`
}

// VerifyErasureRequest records how the requester was verified to be the
// person (a different operator from the one who received it).
type VerifyErasureRequest struct {
	ID     int64  `json:"id" validate:"required"`
	Method string `json:"method" example:"replied from the account's email; confirmed the last escrow id" validate:"required"`
}

// ErasureRequestID names a request to erase.
type ErasureRequestID struct {
	ID int64 `json:"id" validate:"required"`
}

// RefuseErasureRequest closes a request without erasing, with the reason
// (unverified, a legal obligation to keep the record).
type RefuseErasureRequest struct {
	ID     int64  `json:"id" validate:"required"`
	Reason string `json:"reason" validate:"required"`
}

// ErasureRequest is a queued request. Email and DisplayName are the
// person's while they exist (empty once erased).
type ErasureRequest struct {
	ID             int64      `json:"id" validate:"required"`
	Status         string     `json:"status" enums:"received,verified,erased,refused" validate:"required"`
	Case           string     `json:"case" validate:"required"`
	Channel        string     `json:"channel" validate:"required"`
	ReceivedOn     time.Time  `json:"receivedOn" validate:"required"`
	Email          string     `json:"email,omitempty"`
	DisplayName    string     `json:"displayName,omitempty"`
	ReceivedBy     string     `json:"receivedBy" validate:"required"`
	ReceivedAt     time.Time  `json:"receivedAt" validate:"required"`
	VerifiedBy     string     `json:"verifiedBy,omitempty"`
	VerifiedMethod string     `json:"verifiedMethod,omitempty"`
	VerifiedAt     *time.Time `json:"verifiedAt,omitempty"`
	ClosedBy       string     `json:"closedBy,omitempty"`
	ClosedAt       *time.Time `json:"closedAt,omitempty"`
	// Outcome is the refusal's reason, or what the erasure did.
	Outcome string `json:"outcome,omitempty"`
}
