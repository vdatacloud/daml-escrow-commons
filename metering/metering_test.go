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
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.event.Validate(); err == nil {
				t.Fatalf("expected validation error for %s", tc.name)
			}
		})
	}
}
