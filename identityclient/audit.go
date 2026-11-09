package identityclient

import "time"

// Audit-reference resolution and person/tenant holds as daml-escrow-identity
// serves them (daml-escrow PLAN.md Phases 96.6, 97.2), for the platform
// operator API's console (104.8). Every call needs identity:audit and a case;
// identity records it. A resolution names people: shown to the operator, never
// logged.

// ResolveAuditRefs asks who is behind logged references, for a case.
type ResolveAuditRefs struct {
	Refs []string `json:"refs" validate:"required"`
	Case string   `json:"case" validate:"required"`
}

// AuditRefsOf asks for a person's references (to search the logs for them),
// for a case: by OktaSub or by Email, one of them.
type AuditRefsOf struct {
	OktaSub string `json:"oktaSub,omitempty"`
	Email   string `json:"email,omitempty"`
	Case    string `json:"case" validate:"required"`
}

// AuditRef is one reference: ref:<version>:<kind>:<16 hex>.
type AuditRef struct {
	Ref        string `json:"ref" validate:"required"`
	Kind       string `json:"kind" validate:"required"`
	KeyVersion string `json:"keyVersion" validate:"required"`
}

// AuditIdentity is who a reference names.
type AuditIdentity struct {
	OktaSub     string `json:"oktaSub" validate:"required"`
	Email       string `json:"email"`
	DisplayName string `json:"displayName,omitempty"`
	DamlUserID  string `json:"damlUserId,omitempty"`
	DamlPartyID string `json:"damlPartyId,omitempty"`
	Role        string `json:"role,omitempty"`
}

// ResolvedRef is a reference and the identity it names.
type ResolvedRef struct {
	AuditRef
	Identity *AuditIdentity `json:"identity"`
}

// Resolution is what resolving found; Unresolved name no one (never
// recorded, or unlinked by erasure).
type Resolution struct {
	Resolved   []ResolvedRef `json:"resolved" validate:"required"`
	Unresolved []string      `json:"unresolved" validate:"required"`
}

// LegalHold places or releases a hold for a case: on a person (OktaSub or
// Email) or on a target (Tenant: tenant:<party set id>, institution:<name>
// or escrow:<escrow id>, with an escrow's Parties). Exactly one of the three.
type LegalHold struct {
	OktaSub string   `json:"oktaSub,omitempty"`
	Email   string   `json:"email,omitempty"`
	Tenant  string   `json:"tenant,omitempty"`
	Parties []string `json:"parties,omitempty"`
	Case    string   `json:"case" validate:"required"`
}

// ActiveHolds lists the active target holds, for a case.
type ActiveHolds struct {
	Case string `json:"case" validate:"required"`
}

// TargetHold is an active hold on a tenant, an institution or an escrow.
type TargetHold struct {
	Key      string    `json:"key" validate:"required"`
	Case     string    `json:"case" validate:"required"`
	PlacedAt time.Time `json:"placedAt" validate:"required"`
}
