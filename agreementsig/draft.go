package agreementsig

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

// DraftSchema names the canonical form of one off-chain escrow draft
// version (daml-escrow PLAN.md Phase 77). Approvals signed under it are
// recorded on the ledger, so it stays verifiable unchanged; new agreements
// use AgreementSchema.
const DraftSchema = "tripart.draft-version/1"

// DraftVersion is the content of one draft version that its parties agree to.
type DraftVersion struct {
	RootID           string
	Version          int
	ContractType     string
	Amount           float64
	Currency         string
	DepositorID      string
	BeneficiaryEmail string
	MediatorID       string
	Terms            json.RawMessage
	Metadata         json.RawMessage
}

// DraftCanonical returns v's canonical document. Empty terms or metadata
// are null.
func DraftCanonical(v DraftVersion) ([]byte, error) {
	doc := map[string]interface{}{
		"schema":           DraftSchema,
		"rootId":           v.RootID,
		"version":          float64(v.Version),
		"contractType":     v.ContractType,
		"amount":           v.Amount,
		"currency":         v.Currency,
		"depositorId":      v.DepositorID,
		"beneficiaryEmail": v.BeneficiaryEmail,
		"mediatorId":       v.MediatorID,
		"terms":            nil,
		"metadata":         nil,
	}
	for key, raw := range map[string]json.RawMessage{"terms": v.Terms, "metadata": v.Metadata} {
		if len(bytes.TrimSpace(raw)) == 0 {
			continue
		}
		parsed, err := decodeJSON(raw)
		if err != nil {
			return nil, fmt.Errorf("draft %s is not JSON: %w", key, err)
		}
		doc[key] = parsed
	}
	var buf bytes.Buffer
	if err := writeCanonical(&buf, doc, anyNumber); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// DraftHash is the lowercase hex SHA-256 of v's canonical document.
func DraftHash(v DraftVersion) (string, error) {
	doc, err := DraftCanonical(v)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(doc)
	return hex.EncodeToString(sum[:]), nil
}

// DraftMessage is the text a party's wallet signs to agree to a draft
// version. It is fixed English (a wallet shows it verbatim; the app
// explains it in the reader's language) and fully determined by its
// arguments.
func DraftMessage(rootID string, version int, versionHash string) string {
	return fmt.Sprintf("Tripart\nI approve version %d of escrow draft %s.\nVersion hash: sha256:%s", version, rootID, versionHash)
}
