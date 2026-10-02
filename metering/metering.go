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

	// SubmittedAt/CompletedAt/Outcome (optional, v0.5.0; all set or all
	// empty) time the command's submission to the ledger and its
	// completion, as the emitter measured them (UTC), and whether it
	// succeeded -- per-tenant timing facts (daml-escrow PLAN.md Phase 68).
	SubmittedAt time.Time      `json:"submittedAt,omitempty"`
	CompletedAt time.Time      `json:"completedAt,omitempty"`
	Outcome     CommandOutcome `json:"outcome,omitempty"`

	// TrafficCostBytes (optional, v0.6.0) is the synchronizer traffic the
	// participant paid to order this command's confirmation request, as
	// Canton reports it on the command's completion (Ledger API v2
	// Completion.paid_traffic_cost) -- daml-escrow PLAN.md Phase 68. 0 means
	// none reported: a synchronizer without traffic fees, or a command
	// rejected before ordering. It is traffic consumed; whether it was bought
	// or covered by the free base rate is not this event's concern.
	TrafficCostBytes int64 `json:"trafficCostBytes,omitempty"`
	// TrafficPriceUSDPerMB/TrafficPriceSource/TrafficPricedAt (optional,
	// v0.6.0; all set or all empty) record the synchronizer's published
	// traffic price at the command, as a raw fact: USD per MB (10^6 bytes)
	// as a non-negative decimal string (e.g. Splice's extraTrafficPrice
	// "16.67"), where it came from (e.g. "scan:extraTrafficPrice"), and when
	// it applied (UTC). Operational cost = TrafficCostBytes / 10^6 x price;
	// what a tenant is charged is the rating layer's job.
	TrafficPriceUSDPerMB string    `json:"trafficPriceUsdPerMb,omitempty"`
	TrafficPriceSource   string    `json:"trafficPriceSource,omitempty"`
	TrafficPricedAt      time.Time `json:"trafficPricedAt,omitempty"`
}

// CommandOutcome is how a timed ledger command ended.
type CommandOutcome string

const (
	CommandOutcomeOK    CommandOutcome = "OK"
	CommandOutcomeError CommandOutcome = "ERROR"
)

// Validate checks LedgerCommandEvent's required fields. OccurredAt is not
// checked against wall-clock time (callers may legitimately backfill).
func (e LedgerCommandEvent) Validate() error {
	var errs validate.Errors
	errs.Add(validate.RequireNonEmpty("tenantId", e.TenantID))
	errs.Add(validate.RequireNonEmpty("escrowId", e.EscrowID))
	errs.Add(validate.RequireNonEmpty("commandType", e.CommandType))
	errs.Add(validate.RequireNonEmpty("participantNode", e.ParticipantNode))
	errs.Add(e.validateTiming())
	errs.Add(e.validateTraffic())
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

	// SpotPriceUSD/PriceSource/PricedAt (optional, v0.5.0; all set or all
	// empty) record the fee asset's USD spot price at the transfer, as a
	// raw fact: SpotPriceUSD is a non-negative decimal string (e.g.
	// "2648.4051136778", USD per whole FeeAsset unit), PriceSource where it
	// came from (e.g. "bitgo:usdRate"), PricedAt when it applied (UTC).
	// Operational cost = FeeBaseUnits scaled to whole units x SpotPriceUSD;
	// what a tenant is charged is the rating layer's job, never this
	// event's -- the price here is history and is never adjusted.
	SpotPriceUSD string    `json:"spotPriceUsd,omitempty"`
	PriceSource  string    `json:"priceSource,omitempty"`
	PricedAt     time.Time `json:"pricedAt,omitempty"`
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
	errs.Add(e.validateSpotPrice())
	return errs.ErrIfAny()
}

// validateSpotPrice: SpotPriceUSD/PriceSource/PricedAt are all set or all
// empty; the price is a non-negative decimal string.
func (e NetworkFeeEvent) validateSpotPrice() error {
	if e.SpotPriceUSD == "" && e.PriceSource == "" && e.PricedAt.IsZero() {
		return nil
	}
	var errs validate.Errors
	errs.Add(validate.RequireNonEmpty("spotPriceUsd", e.SpotPriceUSD))
	errs.Add(requireDecimal("spotPriceUsd", e.SpotPriceUSD))
	errs.Add(validate.RequireNonEmpty("priceSource", e.PriceSource))
	if e.PricedAt.IsZero() {
		errs.Add(validate.RequireNonEmpty("pricedAt", ""))
	}
	return errs.ErrIfAny()
}

// requireDecimal rejects anything but a non-negative base-10 decimal string
// ("2648.41", "0", "1") -- no sign, exponent or thousands separators (empty
// is left to RequireNonEmpty).
func requireDecimal(field, v string) error {
	dot := false
	for i, c := range v {
		switch {
		case c >= '0' && c <= '9':
		case c == '.' && !dot && i > 0 && i < len(v)-1:
			dot = true
		default:
			return fmt.Errorf("%s must be a non-negative decimal string, got %q", field, v)
		}
	}
	return nil
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

// validateTiming: SubmittedAt/CompletedAt/Outcome are all set or all empty;
// when set, completion is not before submission.
func (e LedgerCommandEvent) validateTiming() error {
	if e.SubmittedAt.IsZero() && e.CompletedAt.IsZero() && e.Outcome == "" {
		return nil
	}
	var errs validate.Errors
	if e.SubmittedAt.IsZero() {
		errs.Add(validate.RequireNonEmpty("submittedAt", ""))
	}
	if e.CompletedAt.IsZero() {
		errs.Add(validate.RequireNonEmpty("completedAt", ""))
	}
	if !e.SubmittedAt.IsZero() && !e.CompletedAt.IsZero() && e.CompletedAt.Before(e.SubmittedAt) {
		errs.Add(fmt.Errorf("completedAt must not be before submittedAt"))
	}
	errs.Add(validate.RequireOneOf("outcome", string(e.Outcome), string(CommandOutcomeOK), string(CommandOutcomeError)))
	return errs.ErrIfAny()
}

// validateTraffic: TrafficCostBytes is non-negative; the traffic price
// fields are all set or all empty, the price a non-negative decimal.
func (e LedgerCommandEvent) validateTraffic() error {
	var errs validate.Errors
	if e.TrafficCostBytes < 0 {
		errs.Add(fmt.Errorf("trafficCostBytes must not be negative"))
	}
	if e.TrafficPriceUSDPerMB == "" && e.TrafficPriceSource == "" && e.TrafficPricedAt.IsZero() {
		return errs.ErrIfAny()
	}
	errs.Add(validate.RequireNonEmpty("trafficPriceUsdPerMb", e.TrafficPriceUSDPerMB))
	errs.Add(requireDecimal("trafficPriceUsdPerMb", e.TrafficPriceUSDPerMB))
	errs.Add(validate.RequireNonEmpty("trafficPriceSource", e.TrafficPriceSource))
	if e.TrafficPricedAt.IsZero() {
		errs.Add(validate.RequireNonEmpty("trafficPricedAt", ""))
	}
	return errs.ErrIfAny()
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
