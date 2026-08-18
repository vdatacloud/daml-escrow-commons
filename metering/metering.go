// Package metering defines the shared usage-event shapes daml-escrow and
// daml-escrow-cms both emit for daml-escrow-platform's rating/billing layer
// (see daml-escrow's PLAN.md Phase 34). These types deliberately carry zero
// pricing logic — they are raw, immutable usage facts; turning them into a
// priced line item is exclusively daml-escrow-platform's rating layer's
// job, kept separate so a rate-schedule change never touches instrumentation
// or reprices already-recorded usage retroactively.
//
// TenantID here is a daml-escrow-identity PartySet id (Phase 34's "Tenant =
// PartySet" decision) — this package has no opinion on, or dependency on,
// how that id is resolved; callers pass it in already resolved.
package metering

import (
	"time"

	"github.com/vdatacloud/daml-escrow-commons/validate"
)

// ChargeBearer designates who absorbs a settlement's platform fee, mirroring
// SWIFT MT103 field 71A's OUR/SHA/BEN convention for cross-border transfer
// charges. Distinct from any correspondent-bank/fiat-rail charge a
// FiatProvider already handles — this is Tripart's own platform fee only.
type ChargeBearer string

const (
	ChargeBearerOur    ChargeBearer = "OUR" // sender (depositor) pays all
	ChargeBearerShared ChargeBearer = "SHA" // shared between depositor and beneficiary
	ChargeBearerBen    ChargeBearer = "BEN" // beneficiary pays all
)

// Rail identifies which settlement rail a SettlementEvent moved funds over.
type Rail string

const (
	RailStablecoin Rail = "stablecoin"
	RailFiat       Rail = "fiat"
)

// LedgerCommandEvent records one Canton synchronizer command submission —
// the infra-cost side of metering. One event per ledger command, emitted at
// the point of submission regardless of outcome (a rejected command still
// consumed synchronizer traffic).
type LedgerCommandEvent struct {
	TenantID        string    `json:"tenantId"`
	EscrowID        string    `json:"escrowId"`
	CommandType     string    `json:"commandType"`
	ParticipantNode string    `json:"participantNode"`
	OccurredAt      time.Time `json:"occurredAt"`
}

// Validate checks LedgerCommandEvent's required fields. OccurredAt is not
// checked against wall-clock time (callers may legitimately backfill).
func (e LedgerCommandEvent) Validate() error {
	var errs validate.Errors
	errs.Add(validate.RequireNonEmpty("tenantId", e.TenantID))
	errs.Add(validate.RequireNonEmpty("escrowId", e.EscrowID))
	errs.Add(validate.RequireNonEmpty("commandType", e.CommandType))
	errs.Add(validate.RequireNonEmpty("participantNode", e.ParticipantNode))
	if e.OccurredAt.IsZero() {
		errs.Add(validate.RequireNonEmpty("occurredAt", ""))
	}
	return errs.ErrIfAny()
}

// SettlementEvent records one actual disbursement — the transaction-fee
// side of metering. Amount is the real settled amount, letting a future
// transaction-based platform fee be rated against it (a cut of GMV), not
// just command volume.
type SettlementEvent struct {
	TenantID     string       `json:"tenantId"`
	EscrowID     string       `json:"escrowId"`
	Amount       float64      `json:"amount"`
	Currency     string       `json:"currency"`
	Rail         Rail         `json:"rail"`
	ChargeBearer ChargeBearer `json:"chargeBearer"`
	OccurredAt   time.Time    `json:"occurredAt"`
}

// Validate checks SettlementEvent's required fields.
func (e SettlementEvent) Validate() error {
	var errs validate.Errors
	errs.Add(validate.RequireNonEmpty("tenantId", e.TenantID))
	errs.Add(validate.RequireNonEmpty("escrowId", e.EscrowID))
	errs.Add(validate.RequirePositive("amount", e.Amount))
	errs.Add(validate.RequireNonEmpty("currency", e.Currency))
	errs.Add(validate.RequireOneOf("rail", string(e.Rail), string(RailStablecoin), string(RailFiat)))
	errs.Add(validate.RequireOneOf("chargeBearer", string(e.ChargeBearer), string(ChargeBearerOur), string(ChargeBearerShared), string(ChargeBearerBen)))
	if e.OccurredAt.IsZero() {
		errs.Add(validate.RequireNonEmpty("occurredAt", ""))
	}
	return errs.ErrIfAny()
}
