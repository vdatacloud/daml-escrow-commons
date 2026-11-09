// Package legalorders is the wire shapes of escrow-api's legal-order acts
// (daml-escrow PLAN.md Phase 104.15): record an order with its document,
// record a review or note with a file, hold an order's objects, lift it.
// escrow-api serves them under escrow:legal-orders; the platform operator
// API relays them for the operator console's legal-order officer. Reading
// and listing orders is identity's (identityclient). An order's reference,
// authority and documents are restricted: never log them; name an order by
// its case id only.
package legalorders

import (
	"errors"
	"regexp"
)

// MaxFile caps an order's document or an attachment (25 MiB).
const MaxFile = 25 << 20

// FileTypes are the content types (as sniffed) an order's document or an
// attachment may be.
var FileTypes = map[string]bool{"application/pdf": true, "image/png": true, "image/jpeg": true}

// Multipart part names.
const (
	PartOrder    = "order"    // RecordRequest as JSON
	PartDocument = "document" // the order's document
	PartEvent    = "event"    // EventRequest as JSON
	PartFile     = "file"     // an event's attachment
)

var (
	caseID = regexp.MustCompile(`^[0-9A-Za-z._-]{1,128}$`)
	day    = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
)

// ValidCase: a case id that can name an order's document in object
// storage (letters, digits, '.', '_' and '-', at most 128).
func ValidCase(c string) bool { return caseID.MatchString(c) }

func validNotifyAfter(v string) bool { return v == "" || v == "never" || day.MatchString(v) }

var errCase = errors.New("case: letters, digits, '.', '_' and '-', at most 128")

// RecordRequest is an order to record. Dates are YYYY-MM-DD.
type RecordRequest struct {
	Case         string `json:"case" example:"LO-2026-014" validate:"required"`
	Kind         string `json:"kind" enums:"preservation,freeze" validate:"required"`
	OrderType    string `json:"orderType" example:"subpoena" validate:"required"`
	Reference    string `json:"reference" validate:"required"`
	Authority    string `json:"authority" validate:"required"`
	Jurisdiction string `json:"jurisdiction"`
	IssuedOn     string `json:"issuedOn" validate:"required"`
	ReceivedOn   string `json:"receivedOn" validate:"required"`
	ReviewOn     string `json:"reviewOn,omitempty"`
	// NotifyAfter: YYYY-MM-DD, "never", or "" (not yet decided).
	NotifyAfter string `json:"notifyAfter,omitempty"`
	// Targets: tenant:<party set id>, institution:<name> or escrow:<escrow id>.
	Targets []string `json:"targets" validate:"required"`
}

// Validate checks the request's shape; identity and escrow-api check the
// rest (targets against the escrow index, dates against today).
func (r RecordRequest) Validate() error {
	switch {
	case !ValidCase(r.Case):
		return errCase
	case r.Kind != "preservation" && r.Kind != "freeze":
		return errors.New("kind: preservation or freeze")
	case r.OrderType == "" || r.Reference == "" || r.Authority == "":
		return errors.New("orderType, reference and authority are required")
	case len(r.Reference) > 512 || len(r.Authority) > 512 || len(r.OrderType) > 64 || len(r.Jurisdiction) > 64:
		return errors.New("reference and authority: at most 512; orderType and jurisdiction: at most 64")
	case !day.MatchString(r.IssuedOn) || !day.MatchString(r.ReceivedOn):
		return errors.New("issuedOn and receivedOn are YYYY-MM-DD")
	case r.ReviewOn != "" && !day.MatchString(r.ReviewOn):
		return errors.New("reviewOn is YYYY-MM-DD")
	case !validNotifyAfter(r.NotifyAfter):
		return errors.New(`notifyAfter is YYYY-MM-DD, "never" or empty`)
	case len(r.Targets) == 0 || len(r.Targets) > 100:
		return errors.New("targets: 1 to 100")
	}
	return nil
}

// EventRequest is a review or note to record (with an optional file part).
type EventRequest struct {
	Case         string `json:"case" validate:"required"`
	Type         string `json:"type" enums:"review,note" validate:"required"`
	Outcome      string `json:"outcome,omitempty" enums:"continue,lift-requested,amended"`
	NextReviewOn string `json:"nextReviewOn,omitempty"`
	Note         string `json:"note,omitempty"`
	NotifyAfter  string `json:"notifyAfter,omitempty"`
}

// Validate checks the request's shape; identity checks the rest.
func (r EventRequest) Validate() error {
	switch {
	case !ValidCase(r.Case):
		return errCase
	case r.Type != "review" && r.Type != "note":
		return errors.New("type: review or note")
	case r.NextReviewOn != "" && !day.MatchString(r.NextReviewOn):
		return errors.New("nextReviewOn is YYYY-MM-DD")
	case !validNotifyAfter(r.NotifyAfter):
		return errors.New(`notifyAfter is YYYY-MM-DD, "never" or empty`)
	case len(r.Note) > 4000:
		return errors.New("note: at most 4000 characters")
	}
	return nil
}

// CaseRequest names an order.
type CaseRequest struct {
	Case string `json:"case" example:"LO-2026-014" validate:"required"`
}

// Validate checks the case id.
func (r CaseRequest) Validate() error {
	if r.Case == "" || len(r.Case) > 256 {
		return errors.New("case: required, at most 256")
	}
	return nil
}

// LiftRequest lifts an order with the authority's release reference.
type LiftRequest struct {
	Case      string `json:"case" validate:"required"`
	Reference string `json:"reference" validate:"required"`
}

// Validate checks the case and reference.
func (r LiftRequest) Validate() error {
	if err := (CaseRequest{Case: r.Case}).Validate(); err != nil {
		return err
	}
	if r.Reference == "" || len(r.Reference) > 512 {
		return errors.New("reference: the authority's release reference, required, at most 512")
	}
	return nil
}

// Sweep is what holding an order's objects did.
type Sweep struct {
	Escrows int `json:"escrows" validate:"required"`
	// Held counts the objects held, by kind: document, ledger-trail,
	// log-archive, order-document.
	Held map[string]int `json:"held" validate:"required"`
	// Unlockable counts objects in a bucket without Object Lock: only the
	// no-deletion checks protect them.
	Unlockable int `json:"unlockable" validate:"required"`
}

// Recorded is what recording an order did. HoldsError, when set, says the
// order is recorded but holding its objects failed: hold them again.
type Recorded struct {
	Case           string `json:"case" validate:"required"`
	DocumentURI    string `json:"documentUri" validate:"required"`
	DocumentSHA256 string `json:"documentSha256" validate:"required"`
	Targets        int    `json:"targets" validate:"required"`
	Holds          *Sweep `json:"holds,omitempty"`
	HoldsError     string `json:"holdsError,omitempty"`
}

// RecordedEvent is what recording a review or note did. HoldsError, when
// set, says its file is stored and the event recorded, but the file isn't
// held yet.
type RecordedEvent struct {
	DocumentURI    string `json:"documentUri,omitempty"`
	DocumentSHA256 string `json:"documentSha256,omitempty"`
	HoldsError     string `json:"holdsError,omitempty"`
}

// Lifted is what lifting did: identity's holds released, and the objects
// whose legal holds were lifted or are still held for another open order.
// ObjectsError, when set, says the order is lifted but its objects aren't
// released yet: lift again, which only releases them (AlreadyLifted).
type Lifted struct {
	Case            string `json:"case" validate:"required"`
	AlreadyLifted   bool   `json:"alreadyLifted,omitempty"`
	HoldsReleased   int    `json:"holdsReleased" validate:"required"`
	ObjectsReleased int    `json:"objectsReleased" validate:"required"`
	ObjectsKept     int    `json:"objectsKept" validate:"required"`
	ObjectsError    string `json:"objectsError,omitempty"`
}
