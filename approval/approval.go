// Package approval is the rule for who may decide a request that needs
// several people's approval (daml-escrow PLAN.md Phase 91, Technical Debt
// TD-1): its proposer may only withdraw it; nobody approves or rejects
// twice, the proposer's proposal counting as their approval; nobody
// approves a request that can no longer take effect as proposed, nor one
// about themselves (an appointment's appointee, say).
//
// The package interprets no domain: the act, the people and their ids are
// opaque strings to it. A server checks a decision with Check before
// sending it; a page offers only Decisions. @vdatacloud/cx-commons
// (sdk/pending-approvals) applies the same rule in the browser, checked
// against testdata/decisions.json.
package approval

import "slices"

// Decision is what someone does with a pending request.
type Decision string

const (
	Approve  Decision = "approve"
	Reject   Decision = "reject"
	Withdraw Decision = "withdraw"
)

// Valid says whether d is one of All.
func (d Decision) Valid() bool { return slices.Contains(All, d) }

// All is every decision, in the order a page offers them.
var All = []Decision{Approve, Reject, Withdraw}

// Request is what the rule needs to know about one pending request.
type Request struct {
	// ProposedBy proposed it; the proposal is their approval.
	ProposedBy string `json:"proposedBy"`
	// Approvals are who approved so far (the proposer may or may not be listed).
	Approvals []string `json:"approvals,omitempty"`
	// Stale: something it depends on changed since it was proposed, so
	// approving it can't take effect (it can still be rejected or withdrawn).
	Stale bool `json:"stale,omitempty"`
	// Subject is who the act is about, if it is about someone who may not
	// approve it -- e.g. the person an appointment appoints.
	Subject string `json:"subject,omitempty"`
}

// Refusal is why a decision isn't allowed. Code is stable (tests and pages
// match on it); the message is for people, so it carries no package prefix:
// a caller wraps it in its own context.
type Refusal struct {
	Code    string
	message string
}

func (r *Refusal) Error() string { return r.message }

// The refusals, for errors.Is.
var (
	ErrUnknownDecision = &Refusal{"unknown-decision", "a decision is approve, reject or withdraw"}
	ErrNotProposer     = &Refusal{"not-proposer", "only its proposer withdraws a request"}
	ErrAlreadyApproved = &Refusal{"already-approved", "you've already approved it"}
	ErrStale           = &Refusal{"stale", "it can no longer take effect as proposed; withdraw or reject it"}
	ErrSubject         = &Refusal{"subject", "nobody approves a request about themselves"}
)

// Check says whether viewer may make decision d on r: nil, or one of the
// refusals above.
func Check(r Request, viewer string, d Decision) error {
	if !d.Valid() {
		return ErrUnknownDecision
	}
	approved := viewer == r.ProposedBy || slices.Contains(r.Approvals, viewer)
	switch {
	case d == Withdraw && viewer != r.ProposedBy:
		return ErrNotProposer
	case d != Withdraw && approved:
		return ErrAlreadyApproved
	case d == Approve && r.Stale:
		return ErrStale
	case d == Approve && r.Subject != "" && viewer == r.Subject:
		return ErrSubject
	}
	return nil
}

// Decisions are the decisions viewer may make on r, in All's order. With
// no viewer (an operator's view, or one not yet known) they are those
// nobody is refused: every decision, less approve when r is stale -- the
// server's Check then has the last word.
func Decisions(r Request, viewer string) []Decision {
	out := []Decision{}
	for _, d := range All {
		if viewer == "" {
			if d != Approve || !r.Stale {
				out = append(out, d)
			}
		} else if Check(r, viewer, d) == nil {
			out = append(out, d)
		}
	}
	return out
}
