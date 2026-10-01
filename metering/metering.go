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
	"fmt"
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
// OccurredAt on both event types below must be set with time.Now().UTC() (or
// an equivalent explicit UTC conversion), not a bare time.Now() -- the
// receiving daml-escrow-platform database stores it in a TIMESTAMPTZ column
// precisely so multiple independent emitters (daml-escrow, and
// daml-escrow-cms once it gains its own ledger-submission point) stay
// chronologically comparable regardless of which timezone each process
// runs in. The column type alone doesn't enforce this -- every emitter must
// standardize on UTC at the call site.
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

	// NetworkFee* (optional, all-or-nothing) record the on-chain network
	// (gas) fee the settlement's transfer incurred, exactly as the custody
	// provider reported it -- a raw fact, never converted or priced here
	// (daml-escrow PLAN.md Phase 67: per-tenant gas attribution).
	//
	// Deprecated: emit a NetworkFeeEvent instead. A SettlementEvent is
	// emitted when the ledger disbursement happens, before the custody
	// provider's transfer that actually incurs the fee, so these fields can't
	// be populated at emission time without double-counting the settlement.
	// Kept (and still validated) so v0.3.0 emitters/consumers keep working.
	// NetworkFeeBaseUnits is a non-negative integer in the fee asset's
	// smallest unit (e.g. wei) as a decimal string, since such amounts
	// overflow float64 precision. NetworkFeeAsset is that asset's ticker
	// (e.g. "sepeth"). NetworkFeePaidFrom says whose funds actually paid
	// it -- the sending wallet's own native balance, or the custody
	// provider's enterprise Gas Tank (a pooled platform cost to recover).
	// Absent on fiat-rail events and on stablecoin events whose provider
	// reports no fee.
	NetworkFeeBaseUnits string          `json:"networkFeeBaseUnits,omitempty"`
	NetworkFeeAsset     string          `json:"networkFeeAsset,omitempty"`
	NetworkFeePaidFrom  NetworkFeePayer `json:"networkFeePaidFrom,omitempty"`
}

// NetworkFeeEvent records the on-chain network (gas) fee one custody-
// provider transfer incurred, emitted right after that transfer -- when the
// fee is actually known -- rather than folded into the SettlementEvent that
// precedes it (daml-escrow PLAN.md Phase 67: per-tenant gas attribution).
// Also covers fees with no settlement attached (refunds, sweeps,
// consolidation). A raw fact exactly as the provider reported it: no
// conversion or pricing here.
//
// FeeBaseUnits is a non-negative integer in FeeAsset's smallest unit (e.g.
// wei) as a decimal string -- such amounts overflow float64/int64.
// TransferID is the provider's own transfer id; Provider names the custody
// provider (e.g. "bitgo"); PaidFrom says whose funds paid it. OccurredAt in
// UTC, as for the other event types.
type NetworkFeeEvent struct {
	TenantID     string          `json:"tenantId"`
	EscrowID     string          `json:"escrowId"`
	Provider     string          `json:"provider"`
	TransferID   string          `json:"transferId"`
	FeeBaseUnits string          `json:"feeBaseUnits"`
	FeeAsset     string          `json:"feeAsset"`
	PaidFrom     NetworkFeePayer `json:"paidFrom"`
	OccurredAt   time.Time       `json:"occurredAt"`
}

// Validate checks NetworkFeeEvent's required fields.
func (e NetworkFeeEvent) Validate() error {
	var errs validate.Errors
	errs.Add(validate.RequireNonEmpty("tenantId", e.TenantID))
	errs.Add(validate.RequireNonEmpty("escrowId", e.EscrowID))
	errs.Add(validate.RequireNonEmpty("provider", e.Provider))
	errs.Add(validate.RequireNonEmpty("transferId", e.TransferID))
	errs.Add(validate.RequireNonEmpty("feeBaseUnits", e.FeeBaseUnits))
	errs.Add(requireBaseUnits("feeBaseUnits", e.FeeBaseUnits))
	errs.Add(validate.RequireNonEmpty("feeAsset", e.FeeAsset))
	errs.Add(validate.RequireOneOf("paidFrom", string(e.PaidFrom), string(NetworkFeePaidBySenderWallet), string(NetworkFeePaidByGasTank)))
	if e.OccurredAt.IsZero() {
		errs.Add(validate.RequireNonEmpty("occurredAt", ""))
	}
	return errs.ErrIfAny()
}

// requireBaseUnits rejects anything but a base-10 non-negative integer
// string (empty is left to RequireNonEmpty).
func requireBaseUnits(field, v string) error {
	for _, c := range v {
		if c < '0' || c > '9' {
			return fmt.Errorf("%s must be a non-negative integer string, got %q", field, v)
		}
	}
	return nil
}

// NetworkFeePayer identifies whose funds paid a network fee.
type NetworkFeePayer string

const (
	NetworkFeePaidBySenderWallet NetworkFeePayer = "SENDER_WALLET"
	NetworkFeePaidByGasTank      NetworkFeePayer = "GAS_TANK"
)

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
	errs.Add(e.validateNetworkFee())
	return errs.ErrIfAny()
}

// validateNetworkFee: the NetworkFee* fields are all set or all empty, and
// when set the amount is a non-negative base-10 integer string.
func (e SettlementEvent) validateNetworkFee() error {
	if e.NetworkFeeBaseUnits == "" && e.NetworkFeeAsset == "" && e.NetworkFeePaidFrom == "" {
		return nil
	}
	var errs validate.Errors
	errs.Add(validate.RequireNonEmpty("networkFeeBaseUnits", e.NetworkFeeBaseUnits))
	errs.Add(validate.RequireNonEmpty("networkFeeAsset", e.NetworkFeeAsset))
	errs.Add(validate.RequireOneOf("networkFeePaidFrom", string(e.NetworkFeePaidFrom), string(NetworkFeePaidBySenderWallet), string(NetworkFeePaidByGasTank)))
	errs.Add(requireBaseUnits("networkFeeBaseUnits", e.NetworkFeeBaseUnits))
	return errs.ErrIfAny()
}
