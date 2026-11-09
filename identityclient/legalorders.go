package identityclient

import "time"

// Legal orders as daml-escrow-identity records them (daml-escrow PLAN.md
// Phases 97.1, 98.3a, 98.4): the wire shapes of its /api/v1/audit/orders*
// endpoints. Shared by escrow-api's `legal-order` CLI and the platform
// operator API (PLAN.md Phase 104.7), so both read and write the same record.
// An order's reference, authority, jurisdiction, document and contacts are
// restricted: never log them; name an order by its case id only.

// LegalOrderTarget is a target a legal order holds: tenant:<party set id>,
// institution:<name> or escrow:<escrow id>, with an escrow's parties' ledger
// ids (the people behind them can't be erased while the hold stands).
type LegalOrderTarget struct {
	Key     string   `json:"key"`
	Parties []string `json:"parties,omitempty"`
}

// LegalOrder is a legal order to record. Dates are YYYY-MM-DD.
type LegalOrder struct {
	Case           string             `json:"case"`
	Kind           string             `json:"kind"`
	OrderType      string             `json:"orderType"`
	Reference      string             `json:"reference"`
	Authority      string             `json:"authority"`
	Jurisdiction   string             `json:"jurisdiction"`
	IssuedOn       string             `json:"issuedOn"`
	ReceivedOn     string             `json:"receivedOn"`
	ReviewOn       string             `json:"reviewOn,omitempty"`
	Targets        []LegalOrderTarget `json:"targets"`
	DocumentURI    string             `json:"documentUri"`
	DocumentSHA256 string             `json:"documentSha256"`
	// NotifyAfter (PLAN.md Phase 98.4a): when the institution may tell the
	// parties: YYYY-MM-DD, "never", or "" (not yet decided).
	NotifyAfter string `json:"notifyAfter,omitempty"`
}

// RecordedLegalOrder is an order as identity returns it. A listing leaves
// out its reference, authority, document and events.
type RecordedLegalOrder struct {
	Case           string     `json:"case"`
	Kind           string     `json:"kind"`
	OrderType      string     `json:"orderType"`
	Reference      string     `json:"reference"`
	Authority      string     `json:"authority"`
	Jurisdiction   string     `json:"jurisdiction"`
	IssuedOn       time.Time  `json:"issuedOn"`
	ReceivedOn     time.Time  `json:"receivedOn"`
	ReviewOn       *time.Time `json:"reviewOn,omitempty"`
	Targets        []string   `json:"targets"`
	DocumentURI    string     `json:"documentUri"`
	DocumentSHA256 string     `json:"documentSha256"`
	RecordedBy     string     `json:"recordedBy"`
	RecordedAt     time.Time  `json:"recordedAt"`
	LiftedBy       string     `json:"liftedBy,omitempty"`
	LiftedAt       *time.Time `json:"liftedAt,omitempty"`
	LiftReference  string     `json:"liftReference,omitempty"`
	NotifyAfter    string     `json:"notifyAfter,omitempty"`
	// Events are the order's reviews and notes, oldest first.
	Events []LegalOrderEvent `json:"events,omitempty"`
}

// Order event types and review outcomes (PLAN.md Phase 98.3a).
const (
	OrderEventReview    = "review"
	OrderEventNote      = "note"
	ReviewContinue      = "continue"
	ReviewLiftRequested = "lift-requested"
	ReviewAmended       = "amended"
)

// LegalOrderEvent is a review of an order or a note on it. A review needs
// an outcome and may set the next review date (YYYY-MM-DD, after today); a
// note needs text and has neither. Either may carry an attachment stored in
// the legal-orders bucket (URI and hash together).
type LegalOrderEvent struct {
	Case           string `json:"case,omitempty"`
	Type           string `json:"type"`
	Outcome        string `json:"outcome,omitempty"`
	NextReviewOn   string `json:"nextReviewOn,omitempty"`
	Note           string `json:"note,omitempty"`
	DocumentURI    string `json:"documentUri,omitempty"`
	DocumentSHA256 string `json:"documentSha256,omitempty"`
	// NotifyAfter, when set, changes the order's notify-after.
	NotifyAfter string     `json:"notifyAfter,omitempty"`
	RecordedBy  string     `json:"recordedBy,omitempty"`
	RecordedAt  *time.Time `json:"recordedAt,omitempty"`
}

// OrderContact is a person behind a party of one of an order's escrow
// targets (PLAN.md Phase 98.4b), for the institution to notify when counsel
// allows. Personal data: never log it.
type OrderContact struct {
	Target      string `json:"target"`
	Party       string `json:"party"`
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
}

// ListLegalOrdersRequest lists orders: open (not lifted), pastReview (open
// and past their review date), for Case (why the operator is looking).
type ListLegalOrdersRequest struct {
	Open       bool   `json:"open"`
	PastReview bool   `json:"pastReview,omitempty"`
	Case       string `json:"case"`
}

// LiftLegalOrderResponse is what lifting an order released at identity.
type LiftLegalOrderResponse struct {
	Case          string `json:"case"`
	HoldsReleased int    `json:"holdsReleased"`
}
