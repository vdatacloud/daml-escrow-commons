// Package gastank holds the wire shapes of escrow-api's Gas Tank watch
// (daml-escrow PLAN.md Phase 104.14): the enterprise Gas Tank's native
// balance on each watched chain against its alert minimum.
package gastank

import "time"

// Path is escrow-api's endpoint (system:admin).
const Path = "/api/v1/system/gas-tank"

// Status is one chain's tank.
type Status struct {
	Coin string `json:"coin" example:"arbeth" validate:"required"`
	// Balance and Minimum are in the native coin ("0.0412"); Balance is
	// empty when unreadable.
	Balance string `json:"balance,omitempty" example:"0.0412"`
	Minimum string `json:"minimum" example:"0.05" validate:"required"`
	// Low: below the minimum, or unreadable (unknown counts as failing).
	Low       bool      `json:"low" validate:"required"`
	Error     string    `json:"error,omitempty"`
	CheckedAt time.Time `json:"checkedAt" validate:"required"`
}

// Response is every watched chain's tank.
type Response struct {
	Networks []Status `json:"networks" validate:"required"`
}
