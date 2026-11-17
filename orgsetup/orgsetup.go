// Package orgsetup holds the wire shape of an organization's setup
// readiness (daml-escrow PLAN.md 110 I4-e): escrow-api checks, read-only,
// what an organization needs before it signs on the ledger -- its party,
// its own signer, its automation party, its administrators, its signing
// policy, its automation mandate and the identity-provider evidence that
// vouches for its staff -- and says what is missing. The operator console
// shows it; every step stays with the organization's administrators or the
// ceremonies (`escrow-api institution-key`, `institution-signers`).
package orgsetup

import (
	"net/url"
	"time"
)

// PathFor is escrow-api's readiness for the named organization (its
// ledger.institutions name): the operator (system:admin) or the
// organization's administrators.
func PathFor(name string) string {
	return "/api/v1/institutions/" + url.PathEscape(name) + "/setup"
}

// Status is one check's outcome.
type Status string

const (
	// Ready: nothing to do.
	Ready Status = "ready"
	// Attention: it works, but not as production should (a key escrow-api
	// holds, a mandate about to expire, a dev-only shortcut).
	Attention Status = "attention"
	// Missing: the organization can't sign, or can't be vouched for, until
	// it is done.
	Missing Status = "missing"
	// NotApplicable: the organization has chosen not to (no automation).
	NotApplicable Status = "n/a"
)

// Check IDs, in the order a setup goes.
const (
	CheckParty          = "party"
	CheckSigning        = "signing"
	CheckAdministrators = "administrators"
	CheckPolicy         = "policy"
	CheckAutomation     = "automation"
	CheckMandate        = "mandate"
	CheckIdentity       = "identity"
)

// Check is one step of the setup and where it stands.
type Check struct {
	ID      string `json:"id" example:"signing" validate:"required"`
	Status  Status `json:"status" example:"ready" validate:"required" enums:"ready,attention,missing,n/a"`
	Summary string `json:"summary" example:"Its own signer signs for its party and its automation party" validate:"required"`
	// Hint is what to do next, for anything but ready.
	Hint string `json:"hint,omitempty"`
}

// Setup is an organization's readiness: Ready when no check is Missing.
type Setup struct {
	Institution string `json:"institution" example:"demo-bank" validate:"required"`
	// Party is the organization's full party id, when it has one (display
	// its short form).
	Party     string    `json:"party,omitempty"`
	Ready     bool      `json:"ready"`
	Checks    []Check   `json:"checks" validate:"required"`
	CheckedAt time.Time `json:"checkedAt" validate:"required"`
}

// Finish sets Ready from the checks.
func (s *Setup) Finish() {
	s.Ready = true
	for _, c := range s.Checks {
		if c.Status == Missing {
			s.Ready = false
		}
	}
}
