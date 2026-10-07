// Package agreementsig is how a party signs one version of an agreement:
// the version's content is reduced to a canonical JSON document and hashed,
// and the signer signs a short fixed-English message naming the version,
// that hash and themselves. Anyone holding the version, the signer's public
// key and the signature can check the approval offline -- which is what lets
// a printable agreement record carry a verifiable signature page
// (daml-escrow PLAN.md Phase 81).
//
// Two schemas live here, each frozen once published (a change is a new
// schema, never an edit):
//
//   - DraftSchema (tripart.draft-version/1): one off-chain escrow draft
//     version, signed by a wallet (Phase 77). Kept so approvals already
//     recorded on the ledger stay verifiable.
//   - AgreementSchema (tripart.agreement-version/1): one version of an
//     agreement or an amendment -- id, version and parent hash, the rendered
//     document's and its canonical source's hashes (never where they are
//     stored), and the terms as an opaque JSON object named by its schema.
//     Numbers in terms are safe integers only; money is a decimal string.
//
// Signatures name their Algorithm (Ed25519, or ECDSA P-256 with SHA-256),
// so a wallet and a KMS-held key are verified alike.
//
// The package interprets no domain: terms, parties and roles are opaque
// strings and JSON to it. @vdatacloud/cx-commons (sdk/agreement-sig)
// mirrors it in TypeScript against the same testdata, so a browser
// recomputes the hash of what it shows before anyone signs.
package agreementsig
