package identityclient

import "time"

// daml-escrow PLAN.md Phase 104.9: sign-in support for people, through
// identity's Okta API app (identity:signin-support). Only people identity
// holds a record for; each act names the case (the ticket) and how the
// caller was verified, and is recorded.

// Sign-in support actions.
const (
	SigninUnlock        = "unlock"
	SigninResetPassword = "reset-password"
	SigninResetMFA      = "reset-mfa"
)

// SigninActions are the acts an operator can take.
var SigninActions = []string{SigninUnlock, SigninResetPassword, SigninResetMFA}

// SigninStatusRequest asks for a person's sign-in state.
type SigninStatusRequest struct {
	Email string `json:"email" example:"joey@depositor.example" validate:"required"`
	Case  string `json:"case" example:"SUP-2026-118" validate:"required"`
}

// SigninActionRequest unlocks a person or resets their password or MFA.
type SigninActionRequest struct {
	Email  string `json:"email" example:"joey@depositor.example" validate:"required"`
	Case   string `json:"case" example:"SUP-2026-118" validate:"required"`
	Action string `json:"action" example:"unlock" enums:"unlock,reset-password,reset-mfa" validate:"required"`
	// Method is how the caller was verified to be the person.
	Method string `json:"method" example:"call-back to the number on file" validate:"required"`
}

// SigninFactor is an enrolled MFA factor (its type only).
type SigninFactor struct {
	Type   string `json:"type" example:"signed_nonce" validate:"required"`
	Status string `json:"status" example:"ACTIVE" validate:"required"`
}

// SigninStatus is a person's sign-in state at the identity provider.
type SigninStatus struct {
	// Status is the provider's (ACTIVE, LOCKED_OUT, PASSWORD_EXPIRED,
	// RECOVERY, SUSPENDED, ...).
	Status      string         `json:"status" example:"LOCKED_OUT" validate:"required"`
	DisplayName string         `json:"displayName,omitempty"`
	LastLogin   *time.Time     `json:"lastLogin,omitempty"`
	Factors     []SigninFactor `json:"factors" validate:"required"`
}

// SigninActionResult is what an act did.
type SigninActionResult struct {
	Action string `json:"action" validate:"required"`
	// Outcome says what the person sees next ("an email with a reset link").
	Outcome string `json:"outcome" validate:"required"`
}
