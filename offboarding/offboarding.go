// Package offboarding holds the wire shapes of escrow-api's termination
// cases, read-only (daml-escrow PLAN.md Phase 104.13): a case's progress,
// its signed certificate and its events, for the operator console's
// offboarding view. Every step stays on the `escrow-api tenant` CLI.
package offboarding

import (
	"encoding/json"
	"time"
)

// Path is escrow-api's list; a case is at Path + "/{id}" (system:admin).
const Path = "/api/v1/system/termination-cases"

// Case is a termination case: one tenant (a party set or an institution)
// from notice to deletion. Actors are operators.
type Case struct {
	ID                string     `json:"id" example:"tc-7f3e" validate:"required"`
	Tenant            string     `json:"tenant" example:"tenant:ps-1" validate:"required"`
	Institution       string     `json:"institution,omitempty"`
	PartySetID        string     `json:"partySetId,omitempty"`
	Notice            string     `json:"notice" validate:"required"`
	OpenedBy          string     `json:"openedBy" validate:"required"`
	OpenedAt          time.Time  `json:"openedAt" validate:"required"`
	ManifestSHA256    string     `json:"manifestSha256,omitempty"`
	ExportedAt        *time.Time `json:"exportedAt,omitempty"`
	ReceiptRef        string     `json:"receiptRef,omitempty"`
	ReceiptRecordedBy string     `json:"receiptRecordedBy,omitempty"`
	ReceiptRecordedAt *time.Time `json:"receiptRecordedAt,omitempty"`
	DeletedAt         *time.Time `json:"deletedAt,omitempty"`
	IdentityErasedAt  *time.Time `json:"identityErasedAt,omitempty"`
	// Certificate is the signed deletion certificate, byte for byte.
	Certificate json.RawMessage `json:"certificate,omitempty" swaggertype:"object"`
}

// Event is one step on a case (append-only).
type Event struct {
	At     time.Time       `json:"at" validate:"required"`
	Action string          `json:"action" example:"exported" validate:"required"`
	Actor  string          `json:"actor" validate:"required"`
	Detail json.RawMessage `json:"detail,omitempty" swaggertype:"object"`
}

// Detail is a case with its events, oldest first.
type Detail struct {
	Case   Case    `json:"case" validate:"required"`
	Events []Event `json:"events" validate:"required"`
}

// Stage is how far a case has come: opened, exported, receipted, deleted
// or closed (deleted and, for a party set, its people erased at identity).
func (c Case) Stage() string {
	switch {
	case c.DeletedAt != nil && (c.PartySetID == "" || c.IdentityErasedAt != nil):
		return "closed"
	case c.DeletedAt != nil:
		return "deleted"
	case c.ReceiptRecordedAt != nil:
		return "receipted"
	case c.ExportedAt != nil:
		return "exported"
	}
	return "opened"
}
