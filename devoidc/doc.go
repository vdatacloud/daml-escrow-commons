// Package devoidc is a minimal OpenID Connect issuer for development and
// tests: it lets each dev persona authenticate with its own signed JWT,
// through the same discovery + JWKS verification production uses with Okta,
// instead of the X-Dev-User header bypass.
//
// DEV ONLY: anyone who can reach its token endpoint can mint a token for any
// identity, scope or client. Run it on localhost (cmd/dev-oidc), or in a test
// process (Start). It generates a fresh RSA key at start, so tokens don't
// survive a restart.
//
// It moved here from daml-escrow's cmd/dev-oidc (daml-escrow PLAN.md Phase
// 104) so every service and browser app on the platform (the escrow portal,
// the operator console, identity, platform) signs in to the same mock.
package devoidc
