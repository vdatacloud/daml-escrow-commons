package metering

import (
	"testing"
	"time"
)

func TestLedgerCommandEvent_Validate(t *testing.T) {
	valid := LedgerCommandEvent{
		TenantID:        "ps-123",
		EscrowID:        "escrow-456",
		CommandType:     "ExerciseCommand:Fund",
		ParticipantNode: "bank",
		OccurredAt:      time.Now(),
	}
	if err := valid.Validate(); err != nil {
		t.Fatalf("expected valid event to pass, got %v", err)
	}

	cases := []struct {
		name  string
		event LedgerCommandEvent
	}{
		{"missing tenantId", LedgerCommandEvent{EscrowID: "e", CommandType: "c", ParticipantNode: "p", OccurredAt: time.Now()}},
		{"missing escrowId", LedgerCommandEvent{TenantID: "t", CommandType: "c", ParticipantNode: "p", OccurredAt: time.Now()}},
		{"missing commandType", LedgerCommandEvent{TenantID: "t", EscrowID: "e", ParticipantNode: "p", OccurredAt: time.Now()}},
		{"missing participantNode", LedgerCommandEvent{TenantID: "t", EscrowID: "e", CommandType: "c", OccurredAt: time.Now()}},
		{"zero occurredAt", LedgerCommandEvent{TenantID: "t", EscrowID: "e", CommandType: "c", ParticipantNode: "p"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.event.Validate(); err == nil {
				t.Fatalf("expected validation error for %s", tc.name)
			}
		})
	}
}

func TestSettlementEvent_Validate(t *testing.T) {
	valid := SettlementEvent{
		TenantID:     "ps-123",
		EscrowID:     "escrow-456",
		Amount:       1000.50,
		Currency:     "USD",
		Rail:         RailStablecoin,
		ChargeBearer: ChargeBearerShared,
		OccurredAt:   time.Now(),
	}
	if err := valid.Validate(); err != nil {
		t.Fatalf("expected valid event to pass, got %v", err)
	}

	cases := []struct {
		name  string
		event SettlementEvent
	}{
		{"missing tenantId", SettlementEvent{EscrowID: "e", Amount: 1, Currency: "USD", Rail: RailFiat, ChargeBearer: ChargeBearerOur, OccurredAt: time.Now()}},
		{"missing escrowId", SettlementEvent{TenantID: "t", Amount: 1, Currency: "USD", Rail: RailFiat, ChargeBearer: ChargeBearerOur, OccurredAt: time.Now()}},
		{"zero amount", SettlementEvent{TenantID: "t", EscrowID: "e", Amount: 0, Currency: "USD", Rail: RailFiat, ChargeBearer: ChargeBearerOur, OccurredAt: time.Now()}},
		{"negative amount", SettlementEvent{TenantID: "t", EscrowID: "e", Amount: -5, Currency: "USD", Rail: RailFiat, ChargeBearer: ChargeBearerOur, OccurredAt: time.Now()}},
		{"missing currency", SettlementEvent{TenantID: "t", EscrowID: "e", Amount: 1, Rail: RailFiat, ChargeBearer: ChargeBearerOur, OccurredAt: time.Now()}},
		{"invalid rail", SettlementEvent{TenantID: "t", EscrowID: "e", Amount: 1, Currency: "USD", Rail: "crypto", ChargeBearer: ChargeBearerOur, OccurredAt: time.Now()}},
		{"invalid chargeBearer", SettlementEvent{TenantID: "t", EscrowID: "e", Amount: 1, Currency: "USD", Rail: RailFiat, ChargeBearer: "XXX", OccurredAt: time.Now()}},
		{"zero occurredAt", SettlementEvent{TenantID: "t", EscrowID: "e", Amount: 1, Currency: "USD", Rail: RailFiat, ChargeBearer: ChargeBearerOur}},
		{"network fee without asset", withFee(valid, "21000000000000", "", NetworkFeePaidByGasTank)},
		{"network fee without payer", withFee(valid, "21000000000000", "sepeth", "")},
		{"network fee bad payer", withFee(valid, "21000000000000", "sepeth", "BANK")},
		{"network fee not an integer", withFee(valid, "0.5", "sepeth", NetworkFeePaidBySenderWallet)},
		{"network fee negative", withFee(valid, "-1", "sepeth", NetworkFeePaidBySenderWallet)},
		{"payer without amount", withFee(valid, "", "sepeth", NetworkFeePaidByGasTank)},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.event.Validate(); err == nil {
				t.Fatalf("expected validation error for %s", tc.name)
			}
		})
	}
}

func withFee(e SettlementEvent, baseUnits, asset string, payer NetworkFeePayer) SettlementEvent {
	e.NetworkFeeBaseUnits, e.NetworkFeeAsset, e.NetworkFeePaidFrom = baseUnits, asset, payer
	return e
}

// The network-fee fields are optional (all-or-nothing) and carry a raw
// base-unit integer string that can exceed float64/int64 range.
func TestSettlementEvent_NetworkFee(t *testing.T) {
	base := SettlementEvent{TenantID: "t", EscrowID: "e", Amount: 1, Currency: "USD", Rail: RailStablecoin, ChargeBearer: ChargeBearerShared, OccurredAt: time.Now()}
	for name, e := range map[string]SettlementEvent{
		"no fee":                   base,
		"gas tank paid":            withFee(base, "21000000000000", "sepeth", NetworkFeePaidByGasTank),
		"sender wallet paid":       withFee(base, "0", "tarbeth", NetworkFeePaidBySenderWallet),
		"beyond int64 (raw units)": withFee(base, "123456789012345678901234567890", "eth", NetworkFeePaidByGasTank),
	} {
		if err := e.Validate(); err != nil {
			t.Errorf("%s: expected valid, got %v", name, err)
		}
	}
}

func TestNetworkFeeEvent_Validate(t *testing.T) {
	valid := NetworkFeeEvent{
		TenantID: "ps-1", EscrowID: "escrow-1", Provider: "bitgo", TransferID: "w1:t1",
		FeeBaseUnits: "123456789012345678901234567890", FeeAsset: "sepeth",
		PaidFrom: NetworkFeePaidByGasTank, OccurredAt: time.Now(),
	}
	if err := valid.Validate(); err != nil {
		t.Fatalf("expected valid event to pass, got %v", err)
	}
	mutate := func(f func(*NetworkFeeEvent)) NetworkFeeEvent { e := valid; f(&e); return e }
	cases := map[string]NetworkFeeEvent{
		"missing tenantId":   mutate(func(e *NetworkFeeEvent) { e.TenantID = "" }),
		"missing escrowId":   mutate(func(e *NetworkFeeEvent) { e.EscrowID = "" }),
		"missing provider":   mutate(func(e *NetworkFeeEvent) { e.Provider = "" }),
		"missing transferId": mutate(func(e *NetworkFeeEvent) { e.TransferID = "" }),
		"missing fee":        mutate(func(e *NetworkFeeEvent) { e.FeeBaseUnits = "" }),
		"decimal fee":        mutate(func(e *NetworkFeeEvent) { e.FeeBaseUnits = "0.5" }),
		"negative fee":       mutate(func(e *NetworkFeeEvent) { e.FeeBaseUnits = "-1" }),
		"missing asset":      mutate(func(e *NetworkFeeEvent) { e.FeeAsset = "" }),
		"bad payer":          mutate(func(e *NetworkFeeEvent) { e.PaidFrom = "BANK" }),
		"zero occurredAt":    mutate(func(e *NetworkFeeEvent) { e.OccurredAt = time.Time{} }),
	}
	for name, e := range cases {
		if err := e.Validate(); err == nil {
			t.Errorf("%s: expected validation error", name)
		}
	}
	if err := mutate(func(e *NetworkFeeEvent) { e.FeeBaseUnits = "0" }).Validate(); err != nil {
		t.Errorf("zero fee should be valid (a sponsored or free transfer), got %v", err)
	}
}

func TestLedgerCommandEvent_Timing(t *testing.T) {
	now := time.Now().UTC()
	base := LedgerCommandEvent{TenantID: "t", EscrowID: "e", CommandType: "Fund", ParticipantNode: "bank", OccurredAt: now}
	timed := base
	timed.SubmittedAt, timed.CompletedAt, timed.Outcome = now, now.Add(1500*time.Millisecond), CommandOutcomeOK
	for name, e := range map[string]LedgerCommandEvent{"untimed (v0.4.0 shape)": base, "timed": timed} {
		if err := e.Validate(); err != nil {
			t.Errorf("%s: expected valid, got %v", name, err)
		}
	}
	bad := map[string]func(*LedgerCommandEvent){
		"outcome only":        func(e *LedgerCommandEvent) { e.SubmittedAt, e.CompletedAt = time.Time{}, time.Time{} },
		"missing completedAt": func(e *LedgerCommandEvent) { e.CompletedAt = time.Time{} },
		"completed before":    func(e *LedgerCommandEvent) { e.CompletedAt = e.SubmittedAt.Add(-time.Second) },
		"unknown outcome":     func(e *LedgerCommandEvent) { e.Outcome = "MAYBE" },
		"missing outcome":     func(e *LedgerCommandEvent) { e.Outcome = "" },
	}
	for name, f := range bad {
		e := timed
		f(&e)
		if err := e.Validate(); err == nil {
			t.Errorf("%s: expected validation error", name)
		}
	}
}

func TestNetworkFeeEvent_SpotPrice(t *testing.T) {
	base := NetworkFeeEvent{TenantID: "t", EscrowID: "e", Provider: "bitgo", TransferID: "w1:t1",
		FeeBaseUnits: "69216383616241", FeeAsset: "sepeth", PaidFrom: NetworkFeePaidBySenderWallet, OccurredAt: time.Now()}
	priced := base
	priced.SpotPriceUSD, priced.PriceSource, priced.PricedAt = "2648.4051136778", "bitgo:usdRate", time.Now()
	for name, e := range map[string]NetworkFeeEvent{"unpriced (v0.4.0 shape)": base, "priced": priced} {
		if err := e.Validate(); err != nil {
			t.Errorf("%s: expected valid, got %v", name, err)
		}
	}
	for _, bad := range []string{"-1", "1e3", "2,648.40", ".5", "5.", "1.2.3", "abc"} {
		e := priced
		e.SpotPriceUSD = bad
		if err := e.Validate(); err == nil {
			t.Errorf("spot price %q: expected validation error", bad)
		}
	}
	partial := base
	partial.SpotPriceUSD = "2648.41"
	if err := partial.Validate(); err == nil {
		t.Error("price without source/pricedAt: expected validation error")
	}
}
