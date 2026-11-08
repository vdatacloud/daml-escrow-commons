package authn

import (
	"fmt"
	"sort"
	"strings"
)

// ClientPolicy keeps guarded scopes to the OAuth clients allowed to use them
// (daml-escrow PLAN.md Phase 104). Operator scopes are granted by Okta only to
// the operator clients, but every service trusts the same audience; the policy
// is the services' own check that a token carrying an operator scope came from
// an operator client. It is keyed by scope, not route, so every existing scope
// check is guarded without each route having to remember to be.
//
// A scope the policy doesn't name passes through unchanged. A guarded scope is
// kept only when the token's client is on its list, so a token without a
// client id never keeps one.
type ClientPolicy struct {
	allowed map[string]map[string]bool
}

// NewClientPolicy builds a policy from scope -> allowed client ids. Every
// scope needs at least one non-empty client id.
func NewClientPolicy(rules map[string][]string) (*ClientPolicy, error) {
	p := &ClientPolicy{allowed: map[string]map[string]bool{}}
	for scope, clients := range rules {
		scope = strings.TrimSpace(scope)
		if scope == "" {
			return nil, fmt.Errorf("authn: client policy: empty scope")
		}
		set := map[string]bool{}
		for _, c := range clients {
			if c = strings.TrimSpace(c); c != "" {
				set[c] = true
			}
		}
		if len(set) == 0 {
			return nil, fmt.Errorf("authn: client policy: scope %q has no clients", scope)
		}
		p.allowed[scope] = set
	}
	return p, nil
}

// ParseClientPolicy reads the configuration form
// "scope=client1,client2;scope2=client3" (whitespace ignored). An empty
// string is an empty policy.
func ParseClientPolicy(s string) (*ClientPolicy, error) {
	rules := map[string][]string{}
	for _, entry := range strings.Split(s, ";") {
		if strings.TrimSpace(entry) == "" {
			continue
		}
		scope, clients, ok := strings.Cut(entry, "=")
		if !ok {
			return nil, fmt.Errorf("authn: client policy: %q is not scope=clients", strings.TrimSpace(entry))
		}
		scope = strings.TrimSpace(scope)
		if _, dup := rules[scope]; dup {
			return nil, fmt.Errorf("authn: client policy: scope %q listed twice", scope)
		}
		rules[scope] = strings.Split(clients, ",")
	}
	return NewClientPolicy(rules)
}

// Guarded lists the scopes the policy guards, sorted.
func (p *ClientPolicy) Guarded() []string {
	if p == nil {
		return nil
	}
	out := make([]string, 0, len(p.allowed))
	for s := range p.allowed {
		out = append(out, s)
	}
	sort.Strings(out)
	return out
}

// Allows reports whether clientID may use scope. A nil policy allows all.
func (p *ClientPolicy) Allows(clientID, scope string) bool {
	if p == nil {
		return true
	}
	set, guarded := p.allowed[scope]
	return !guarded || set[clientID]
}

// Filter splits scopes into those clientID may use and those withheld.
func (p *ClientPolicy) Filter(clientID string, scopes []string) (kept, withheld []string) {
	for _, s := range scopes {
		if p.Allows(clientID, s) {
			kept = append(kept, s)
		} else {
			withheld = append(withheld, s)
		}
	}
	return kept, withheld
}
