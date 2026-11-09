// Package authority reports whether a service runs on authoritative
// mechanisms (daml-escrow PLAN.md Phase 104.11, docs/PLATFORM_OPS.md
// "Production authority check"). A stand-in is a dev convenience, a mock or a
// placeholder a configuration selects; outside ENVIRONMENT=dev a service
// refuses to start while it selects any, so a report outside dev lists none.
// The operator console shows each service's report per environment.
package authority

import (
	"encoding/json"
	"net/http"
	"slices"

	"github.com/vdatacloud/daml-escrow-commons/auditref"
	"github.com/vdatacloud/daml-escrow-commons/authn"
)

// Path is where each service serves its report, behind its operator scope.
const Path = "/api/v1/system/authority"

// StandIn is one non-authoritative mechanism and the setting selecting it.
type StandIn struct {
	Setting string `json:"setting" example:"AUTH_BYPASS" validate:"required"`
	Why     string `json:"why" example:"requests are accepted without a verified token" validate:"required"`
}

// Report is a service's production authority.
type Report struct {
	Service     string `json:"service" example:"escrow-api" validate:"required"`
	Environment string `json:"environment" example:"staging" validate:"required"`
	Version     string `json:"version,omitempty" example:"v0.84.0"`
	// Authoritative: the service selects no stand-in.
	Authoritative bool      `json:"authoritative" validate:"required"`
	StandIns      []StandIn `json:"standIns" validate:"required"`
}

// NewReport is the service's report over the stand-ins it selects.
func NewReport(service, environment, version string, standIns ...[]StandIn) Report {
	all := []StandIn{}
	for _, s := range standIns {
		all = append(all, s...)
	}
	return Report{Service: service, Environment: environment, Version: version, Authoritative: len(all) == 0, StandIns: all}
}

// Auth lists the authentication stand-ins: the header bypass, and operator
// scopes the client policy doesn't keep to their clients.
func Auth(bypass bool, clientPolicy string, operatorScopes []string) []StandIn {
	var out []StandIn
	if bypass {
		out = append(out, StandIn{"AUTH_BYPASS", "requests are accepted without a verified token (X-Dev-* headers)"})
	}
	policy, err := authn.ParseClientPolicy(clientPolicy)
	guarded := policy.Guarded()
	for _, s := range operatorScopes {
		if err != nil || !slices.Contains(guarded, s) {
			out = append(out, StandIn{"AUTH_CLIENT_POLICY", "operator scopes are accepted from any client, not only the operator clients"})
			break
		}
	}
	return out
}

// Logging lists the logging stand-in: an environment whose logs name people
// in the clear instead of by audit reference.
func Logging(environment string) []StandIn {
	if auditref.Plain(environment) {
		return []StandIn{{"ENVIRONMENT", "logs name people in the clear, not by audit reference"}}
	}
	return nil
}

// Handler serves the report. The caller's route checks the operator scope.
func Handler(r Report) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(r)
	}
}
