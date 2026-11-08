// Package pathparam reads a chi route parameter as the value the client
// meant, not as it was escaped on the wire.
package pathparam

import (
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
)

// Get is chi.URLParam, unescaped. When a path holds a character the browser
// escapes but Go's path encoding doesn't -- ":" in every Canton party
// (`ada::1220…`), "@" in an email or a dev subject -- net/url keeps RawPath
// and chi matches on it, so chi.URLParam hands back "ada%3A%3A1220…" (found
// in daml-escrow PLAN.md Phase 100: Unlink on My signing never matched its
// link). A value that isn't valid escaping is returned as is.
func Get(r *http.Request, key string) string {
	v := chi.URLParam(r, key)
	if u, err := url.PathUnescape(v); err == nil {
		return u
	}
	return v
}
