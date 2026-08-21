// Package identityclient is the core HTTP client for talking to
// daml-escrow-identity — the shared T1 identity directory service. This is
// deliberately the *minimal* subset every consumer needs (T1 lookup/upsert
// plus the manager-override check), not a full wrapper of that service's
// entire API surface. daml-escrow-cms uses this client directly, unmodified,
// since T1 is all it ever needs (INTEGRATION.md §4). daml-escrow needs more
// (T2/T3 write-back, party-set CRUD, pending invitations) and composes its
// own local client on top of this one via struct embedding rather than
// pushing that superset here — see daml-escrow/internal/identityclient's
// package doc for why.
//
// This package never imports daml-escrow-identity's own Go module — that
// service has its own database and its own deploy lifecycle; pulling its
// module in as a library would bypass the HTTP boundary, AuthMiddleware,
// and scope enforcement between two independently-deployed services. It is
// called over HTTP like any other external dependency.
package identityclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

func New(baseURL, token string) *Client {
	return &Client{baseURL: baseURL, token: token, httpClient: &http.Client{}}
}

// Identity mirrors daml-escrow-identity's response shape. DamlUserID/
// DamlPartyID are present because the identity service's response carries
// them (for daml-escrow's benefit) -- a T1-only consumer like daml-escrow-cms
// simply never reads those two fields, matching the boundary its own
// INTEGRATION.md §4 already documents.
type Identity struct {
	OktaSub       string `json:"oktaSub"`
	IdentityToken string `json:"identityToken"`
	Email         string `json:"email"`
	DisplayName   string `json:"displayName,omitempty"`
	DamlUserID    string `json:"damlUserId,omitempty"`
	DamlPartyID   string `json:"damlPartyId,omitempty"`
	Role          string `json:"role,omitempty"`
}

// Upsert registers or updates the T1 record for oktaSub.
func (c *Client) Upsert(ctx context.Context, oktaSub, email, displayName, role string) (*Identity, error) {
	body, err := json.Marshal(map[string]string{
		"oktaSub":     oktaSub,
		"email":       email,
		"displayName": displayName,
		"role":        role,
	})
	if err != nil {
		return nil, fmt.Errorf("identityclient: marshalling identity upsert request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/api/v1/identities", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("identityclient: building identity upsert request: %w", err)
	}
	return c.doIdentityRequest(req)
}

// GetByOktaSub looks up an identity by its okta_sub -- the primary lookup
// path for a JIT-provisioning flow.
func (c *Client) GetByOktaSub(ctx context.Context, oktaSub string) (*Identity, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/api/v1/identities/"+url.PathEscape(oktaSub), nil)
	if err != nil {
		return nil, fmt.Errorf("identityclient: building identity lookup request: %w", err)
	}
	return c.doIdentityRequest(req)
}

// GetByEmail looks up an identity by email -- used to resolve a
// beneficiary_email-style placeholder once the invited counterparty has an
// identity on file.
func (c *Client) GetByEmail(ctx context.Context, email string) (*Identity, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/api/v1/identities/by-email/"+url.PathEscape(email), nil)
	if err != nil {
		return nil, fmt.Errorf("identityclient: building identity lookup request: %w", err)
	}
	return c.doIdentityRequest(req)
}

// GetByToken looks up an identity by its opaque identity_token -- the
// reference a T1-only consumer stores locally instead of a raw okta_sub.
func (c *Client) GetByToken(ctx context.Context, identityToken string) (*Identity, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/api/v1/identities/by-token/"+url.PathEscape(identityToken), nil)
	if err != nil {
		return nil, fmt.Errorf("identityclient: building identity lookup request: %w", err)
	}
	return c.doIdentityRequest(req)
}

// ManagesIdentity asks daml-escrow-identity whether managerOktaSub manages
// a party set containing targetOktaSub -- the manager-override check both
// consumer repos need for their own authorization flows.
func (c *Client) ManagesIdentity(ctx context.Context, managerOktaSub, targetOktaSub string) (bool, error) {
	q := url.Values{"manager": {managerOktaSub}, "target": {targetOktaSub}}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/api/v1/party-sets/manages?"+q.Encode(), nil)
	if err != nil {
		return false, fmt.Errorf("identityclient: building manages-identity request: %w", err)
	}
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("identityclient: calling daml-escrow-identity: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("identityclient: daml-escrow-identity returned status %d", resp.StatusCode)
	}
	var result struct {
		Manages bool `json:"manages"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, fmt.Errorf("identityclient: decoding manages-identity response: %w", err)
	}
	return result.Manages, nil
}

// doIdentityRequest is unexported -- a composing client (e.g. daml-escrow's
// local wrapper) reaches identity responses through the public methods
// above, not by calling this directly.
func (c *Client) doIdentityRequest(req *http.Request) (*Identity, error) {
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("identityclient: calling daml-escrow-identity: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("identityclient: daml-escrow-identity returned status %d", resp.StatusCode)
	}

	var identity Identity
	if err := json.NewDecoder(resp.Body).Decode(&identity); err != nil {
		return nil, fmt.Errorf("identityclient: decoding identity response: %w", err)
	}
	return &identity, nil
}
