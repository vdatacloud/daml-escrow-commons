// Command dev-oidc runs the devoidc development OpenID Connect issuer on
// localhost for every repo's dev stack (daml-escrow PLAN.md Phases 69, 91.1
// and 104).
//
// DEV ONLY: anyone who can reach POST /token can mint a token for any
// identity. Keep it on localhost.
//
//	go run github.com/vdatacloud/daml-escrow-commons/cmd/dev-oidc -addr 127.0.0.1:8095 \
//	    -origins http://localhost:4321,http://localhost:4322
//	curl -s -d 'sub=joey&email=joey@depositor.devlocal&scope=escrow:read escrow:write' localhost:8095/token
package main

import (
	"flag"
	"log"
	"net/http"
	"slices"
	"strings"

	"github.com/vdatacloud/daml-escrow-commons/devoidc"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:8095", "listen address (keep it on localhost)")
	audience := flag.String("audience", "daml-escrow", "default aud claim (the services' AUTH_AUDIENCE)")
	origins := flag.String("origins", devoidc.DefaultOrigin, "comma-separated browser origins sign-ins may return to (escrow portal, operator console)")
	portal := flag.String("portal", "", "a single browser origin (older flag; added to -origins)")
	// The issuer URL (iss) can name the listener by another host: LocalNet
	// uses http://localhost:8095, the host the browser already trusts.
	issuerURL := flag.String("url", "", "issuer URL (default http://<addr>)")
	flag.Parse()

	if *issuerURL == "" {
		*issuerURL = "http://" + *addr
	}
	list := splitOrigins(*origins, *portal)
	is, err := devoidc.New(*issuerURL, *audience, list...)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("dev-oidc issuer %s for %s (DEV ONLY -- mints tokens for any identity)", is.URL(), strings.Join(list, ", "))
	log.Fatal(http.ListenAndServe(*addr, is.Handler()))
}

// splitOrigins is -origins plus -portal, each once, blanks dropped.
func splitOrigins(origins, portal string) []string {
	var out []string
	for _, o := range append(strings.Split(origins, ","), portal) {
		o = strings.TrimRight(strings.TrimSpace(o), "/")
		if o != "" && !slices.Contains(out, o) {
			out = append(out, o)
		}
	}
	return out
}
