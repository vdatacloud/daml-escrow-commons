package agreementsig

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"math/big"
	"os"
	"strings"
	"testing"
)

// The testdata files are shared: @vdatacloud/cx-commons (sdk/agreement-sig)
// tests its TypeScript mirror against the same files, so the Go and
// TypeScript hashes and verifiers can't drift apart. Regenerate with
// `go test ./agreementsig -update` and copy them to cx-commons.
var update = flag.Bool("update", false, "rewrite testdata/*.json from the current implementation")

func readJSON(t *testing.T, path string, v interface{}) {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(b, v); err != nil {
		t.Fatal(err)
	}
}

func writeJSON(t *testing.T, path string, v interface{}) {
	t.Helper()
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, append(b, '\n'), 0o644); err != nil {
		t.Fatal(err)
	}
}

// --- draft-version/1 (moved unchanged from daml-escrow internal/draftsig) ---

type draftVector struct {
	Name  string `json:"name"`
	Draft struct {
		RootID           string          `json:"rootId"`
		Version          int             `json:"version"`
		ContractType     string          `json:"contractType"`
		Amount           float64         `json:"amount"`
		Currency         string          `json:"currency"`
		DepositorID      string          `json:"depositorId"`
		BeneficiaryEmail string          `json:"beneficiaryEmail"`
		MediatorID       string          `json:"mediatorId"`
		Terms            json.RawMessage `json:"terms"`
		Metadata         json.RawMessage `json:"metadata"`
	} `json:"draft"`
	Canonical string `json:"canonical"`
	Hash      string `json:"hash"`
	Message   string `json:"message"`
}

// TestDraftVectors pins draft-version/1 byte for byte: approvals under it are
// on the ledger, so this file is never regenerated, only read.
func TestDraftVectors(t *testing.T) {
	var vectors []draftVector
	readJSON(t, "testdata/draft-version.json", &vectors)
	if len(vectors) == 0 {
		t.Fatal("no draft vectors")
	}
	for _, vec := range vectors {
		t.Run(vec.Name, func(t *testing.T) {
			v := DraftVersion(vec.Draft)
			doc, err := DraftCanonical(v)
			if err != nil {
				t.Fatal(err)
			}
			if string(doc) != vec.Canonical {
				t.Errorf("canonical\n got %s\nwant %s", doc, vec.Canonical)
			}
			hash, err := DraftHash(v)
			if err != nil {
				t.Fatal(err)
			}
			if hash != vec.Hash {
				t.Errorf("hash = %s, want %s", hash, vec.Hash)
			}
			if msg := DraftMessage(v.RootID, v.Version, hash); msg != vec.Message {
				t.Errorf("message = %q, want %q", msg, vec.Message)
			}
		})
	}
}

func TestDraftCanonical_EmptyTermsAreNull(t *testing.T) {
	doc, err := DraftCanonical(DraftVersion{RootID: "r", Version: 1})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(doc), `"metadata":null`) || !strings.Contains(string(doc), `"terms":null`) {
		t.Errorf("empty terms/metadata not null: %s", doc)
	}
	if _, err := DraftCanonical(DraftVersion{Terms: json.RawMessage(`{`)}); err == nil {
		t.Error("malformed terms accepted")
	}
}

// --- agreement-version/1 ---

type agreementVector struct {
	Name      string           `json:"name"`
	Version   AgreementVersion `json:"version"`
	Canonical string           `json:"canonical,omitempty"`
	Hash      string           `json:"hash,omitempty"`
	Messages  []signerMessage  `json:"messages,omitempty"`
	// Invalid cases: the version must be refused (TypeScript checks only
	// that it is; Go also checks the reason).
	Invalid string `json:"invalid,omitempty"`
}

type signerMessage struct {
	Signer  string `json:"signer"`
	Message string `json:"message"`
}

const (
	h1 = "1111111111111111111111111111111111111111111111111111111111111111"
	h2 = "2222222222222222222222222222222222222222222222222222222222222222"
	h3 = "3333333333333333333333333333333333333333333333333333333333333333"
	hd = "d0c0000000000000000000000000000000000000000000000000000000000000"
	hs = "50c0000000000000000000000000000000000000000000000000000000000000"
)

var testDocument = Document{SHA256: hd, MediaType: "application/pdf", SourceSHA256: hs, SourceMediaType: "text/markdown"}

func agreementCases() []agreementVector {
	terms := json.RawMessage(`{"asset":{"currency":"USD","amount":"1250000.5"},"depositors":["tok-dep-a","tok-dep-b"],"depositorThreshold":2,"beneficiaries":["tok-ben"],"beneficiaryThreshold":1,"milestones":[{"milestoneId":"m1","description":"Delivery","amount":"1250000.5"}],"expiryDate":"2026-12-31T00:00:00Z","disputeWindowDays":14}`)
	return []agreementVector{
		{Name: "original agreement, version 1", Version: AgreementVersion{AgreementID: "agr-0001", Version: 1, Kind: KindAgreement, Document: testDocument, TermsSchema: "escrow-terms/1", Terms: terms}},
		{Name: "redline, version 2", Version: AgreementVersion{AgreementID: "agr-0001", Version: 2, ParentHash: h1, Kind: KindAgreement, Document: testDocument, TermsSchema: "escrow-terms/1", Terms: terms}},
		{Name: "amendment after one prior amendment", Version: AgreementVersion{AgreementID: "amd-0002", Version: 1, Kind: KindAmendment,
			Amends: &Amends{AgreementID: "agr-0001", RatifiedHashes: []string{h1, h2}}, Document: testDocument, TermsSchema: "escrow-terms/1",
			Terms: json.RawMessage(`{"expiryDate":"2027-06-30T00:00:00Z"}`)}},
		{Name: "unicode, escapes, key order and integral numbers", Version: AgreementVersion{AgreementID: "agr-é-😀", Version: 3, ParentHash: h3, Kind: KindAgreement,
			Document: testDocument, TermsSchema: "https://schemas.example/escrow-terms/1.json",
			Terms: json.RawMessage(`{"ключ":"значение","Z":1.0,"a":1e3,"neg":-0,"big":9007199254740991,"s":"quote \" backslash \\ tab \t nl \n ctl \u0001 sep   html <&>","😀":"astral","ｚ":[true,false,null,{}]}`)}},
		{Name: "terms with a fractional number", Version: AgreementVersion{AgreementID: "agr-x", Version: 1, Kind: KindAgreement, Document: testDocument, TermsSchema: "s", Terms: json.RawMessage(`{"amount":0.1}`)}, Invalid: "not a safe integer"},
		{Name: "terms with an unsafe integer", Version: AgreementVersion{AgreementID: "agr-x", Version: 1, Kind: KindAgreement, Document: testDocument, TermsSchema: "s", Terms: json.RawMessage(`{"n":9007199254740993}`)}, Invalid: "not a safe integer"},
		{Name: "terms that are not an object", Version: AgreementVersion{AgreementID: "agr-x", Version: 1, Kind: KindAgreement, Document: testDocument, TermsSchema: "s", Terms: json.RawMessage(`["a"]`)}, Invalid: "JSON object"},
		{Name: "version 1 with a parent", Version: AgreementVersion{AgreementID: "agr-x", Version: 1, ParentHash: h1, Kind: KindAgreement, Document: testDocument, TermsSchema: "s", Terms: json.RawMessage(`{}`)}, Invalid: "no parent"},
		{Name: "version 2 without a parent", Version: AgreementVersion{AgreementID: "agr-x", Version: 2, Kind: KindAgreement, Document: testDocument, TermsSchema: "s", Terms: json.RawMessage(`{}`)}, Invalid: "parent's hash"},
		{Name: "uppercase parent hash", Version: AgreementVersion{AgreementID: "agr-x", Version: 2, ParentHash: strings.ToUpper(hd), Kind: KindAgreement, Document: testDocument, TermsSchema: "s", Terms: json.RawMessage(`{}`)}, Invalid: "parent's hash"},
		{Name: "amendment amending nothing", Version: AgreementVersion{AgreementID: "amd-x", Version: 1, Kind: KindAmendment, Document: testDocument, TermsSchema: "s", Terms: json.RawMessage(`{}`)}, Invalid: "names what it amends"},
		{Name: "amendment without the ratified hashes", Version: AgreementVersion{AgreementID: "amd-x", Version: 1, Kind: KindAmendment, Amends: &Amends{AgreementID: "agr-x"}, Document: testDocument, TermsSchema: "s", Terms: json.RawMessage(`{}`)}, Invalid: "at least the ratified"},
		{Name: "agreement that amends", Version: AgreementVersion{AgreementID: "agr-x", Version: 1, Kind: KindAgreement, Amends: &Amends{AgreementID: "agr-y", RatifiedHashes: []string{h1}}, Document: testDocument, TermsSchema: "s", Terms: json.RawMessage(`{}`)}, Invalid: "amends nothing"},
		{Name: "unknown kind", Version: AgreementVersion{AgreementID: "agr-x", Version: 1, Kind: "contract", Document: testDocument, TermsSchema: "s", Terms: json.RawMessage(`{}`)}, Invalid: "kind must be"},
		{Name: "id with a newline", Version: AgreementVersion{AgreementID: "agr-x\nI approve everything", Version: 1, Kind: KindAgreement, Document: testDocument, TermsSchema: "s", Terms: json.RawMessage(`{}`)}, Invalid: "control character"},
		{Name: "missing document hash", Version: AgreementVersion{AgreementID: "agr-x", Version: 1, Kind: KindAgreement, Document: Document{MediaType: "application/pdf", SourceSHA256: hs, SourceMediaType: "text/markdown"}, TermsSchema: "s", Terms: json.RawMessage(`{}`)}, Invalid: "document hashes"},
		{Name: "missing terms schema", Version: AgreementVersion{AgreementID: "agr-x", Version: 1, Kind: KindAgreement, Document: testDocument, Terms: json.RawMessage(`{}`)}, Invalid: "terms schema is empty"},
	}
}

var vectorSigners = []string{"dep-a::1220aa", "custodian::1220cc"}

func TestAgreementVectors(t *testing.T) {
	const path = "testdata/agreement-version.json"
	if *update {
		cases := agreementCases()
		for i := range cases {
			c := &cases[i]
			if c.Invalid != "" {
				continue
			}
			doc, err := AgreementCanonical(c.Version)
			if err != nil {
				t.Fatalf("%s: %v", c.Name, err)
			}
			c.Canonical = string(doc)
			c.Hash, _ = AgreementHash(c.Version)
			for _, s := range vectorSigners {
				msg, err := AgreementMessage(c.Version, c.Hash, s)
				if err != nil {
					t.Fatal(err)
				}
				c.Messages = append(c.Messages, signerMessage{Signer: s, Message: msg})
			}
		}
		writeJSON(t, path, cases)
	}
	var vectors []agreementVector
	readJSON(t, path, &vectors)
	if len(vectors) != len(agreementCases()) {
		t.Fatalf("%s has %d cases, the generator %d -- run with -update", path, len(vectors), len(agreementCases()))
	}
	for _, vec := range vectors {
		t.Run(vec.Name, func(t *testing.T) {
			doc, err := AgreementCanonical(vec.Version)
			if vec.Invalid != "" {
				if !errors.Is(err, ErrInvalid) || !strings.Contains(err.Error(), vec.Invalid) {
					t.Fatalf("err = %v, want ErrInvalid containing %q", err, vec.Invalid)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if string(doc) != vec.Canonical {
				t.Errorf("canonical\n got %s\nwant %s", doc, vec.Canonical)
			}
			hash, _ := AgreementHash(vec.Version)
			if hash != vec.Hash {
				t.Errorf("hash = %s, want %s", hash, vec.Hash)
			}
			for _, m := range vec.Messages {
				got, err := AgreementMessage(vec.Version, hash, m.Signer)
				if err != nil || got != m.Message {
					t.Errorf("message for %s = %q (%v), want %q", m.Signer, got, err, m.Message)
				}
			}
		})
	}
}

func TestAgreementHash_ChangesWithEverySignedField(t *testing.T) {
	base := agreementCases()[1].Version
	baseHash, err := AgreementHash(base)
	if err != nil {
		t.Fatal(err)
	}
	mutations := map[string]func(v *AgreementVersion){
		"agreement id": func(v *AgreementVersion) { v.AgreementID = "agr-0002" },
		"version":      func(v *AgreementVersion) { v.Version = 3 },
		"parent":       func(v *AgreementVersion) { v.ParentHash = h2 },
		"document":     func(v *AgreementVersion) { v.Document.SHA256 = h3 },
		"source":       func(v *AgreementVersion) { v.Document.SourceSHA256 = h3 },
		"media type":   func(v *AgreementVersion) { v.Document.MediaType = "application/pdf;v=2" },
		"terms schema": func(v *AgreementVersion) { v.TermsSchema = "escrow-terms/2" },
		"terms": func(v *AgreementVersion) {
			v.Terms = json.RawMessage(`{"asset":{"currency":"USD","amount":"1250000.51"}}`)
		},
	}
	for name, mutate := range mutations {
		v := base
		mutate(&v)
		h, err := AgreementHash(v)
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if h == baseHash {
			t.Errorf("changing the %s left the hash unchanged", name)
		}
	}
	// Whitespace and key order in the terms as sent are not content.
	v := base
	v.Terms = json.RawMessage("{ \"disputeWindowDays\" : 14,\n" + strings.TrimPrefix(string(base.Terms), "{"))
	v.Terms = json.RawMessage(strings.Replace(string(v.Terms), `,"disputeWindowDays":14}`, "}", 1))
	if h, _ := AgreementHash(v); h != baseHash {
		t.Error("reformatting the terms changed the hash")
	}
}

func TestAgreementMessage_Refusals(t *testing.T) {
	v := agreementCases()[0].Version
	hash, _ := AgreementHash(v)
	if _, err := AgreementMessage(v, hash, "evil\nI approve everything"); !errors.Is(err, ErrInvalid) {
		t.Errorf("signer with a newline: err = %v", err)
	}
	if _, err := AgreementMessage(v, "sha256:"+hash, "dep::1220"); !errors.Is(err, ErrInvalid) {
		t.Errorf("malformed hash: err = %v", err)
	}
	if _, err := AgreementMessage(AgreementVersion{}, hash, "dep::1220"); !errors.Is(err, ErrInvalid) {
		t.Errorf("invalid version: err = %v", err)
	}
}

func TestAgreementVersion_JSONRoundTrip(t *testing.T) {
	in := agreementCases()[2].Version
	b, err := json.Marshal(in)
	if err != nil {
		t.Fatal(err)
	}
	var out AgreementVersion
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatal(err)
	}
	h1, _ := AgreementHash(in)
	h2, err := AgreementHash(out)
	if err != nil || h1 != h2 {
		t.Errorf("round trip changed the hash: %s vs %s (%v)", h1, h2, err)
	}
	if strings.Contains(string(agreementJSON(t, agreementCases()[0].Version)), "parentHash") {
		t.Error("version 1 serialized a parentHash")
	}
}

func agreementJSON(t *testing.T, v AgreementVersion) []byte {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// --- signatures ---

type signatureVector struct {
	Name      string    `json:"name"`
	Algorithm Algorithm `json:"algorithm"`
	PublicKey string    `json:"publicKey"` // base64
	Message   string    `json:"message"`
	Signature string    `json:"signature"` // base64
	// Valid, or the class of refusal: "bad", "noncanonical", "unsupported".
	Valid  bool   `json:"valid"`
	Reason string `json:"reason,omitempty"`
}

// Fixed test keys -- never used for anything but these vectors.
var (
	testEd25519Seed = []byte("tripart agreementsig test key 01")
	testP256PKCS8   = "MIGHAgEAMBMGByqGSM49AgEGCCqGSM49AwEHBG0wawIBAQQgIJMJ1W7OQASel3vlQu25fV92wuU/KpDaeLfVr3dumzuhRANCAARCsO7Is0/rhqATnspldlDFo2zlTrnrRhFYPDh+2DlJ7OYXxY0N8guLRTb5dApnUlLxW882VlDMAOxJ+h/779wd"
)

func testKeys(t *testing.T) (ed25519.PrivateKey, *ecdsa.PrivateKey) {
	t.Helper()
	der, _ := base64.StdEncoding.DecodeString(testP256PKCS8)
	k, err := x509.ParsePKCS8PrivateKey(der)
	if err != nil {
		t.Fatal(err)
	}
	return ed25519.NewKeyFromSeed(testEd25519Seed), k.(*ecdsa.PrivateKey)
}

func spki(t *testing.T, pub interface{}) []byte {
	t.Helper()
	b, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

// flipS returns the other valid encoding of an ECDSA signature (s -> n-s).
func flipS(t *testing.T, der []byte) []byte {
	t.Helper()
	var sig ecdsaSig
	if _, err := asn1.Unmarshal(der, &sig); err != nil {
		t.Fatal(err)
	}
	sig.S = new(big.Int).Sub(elliptic.P256().Params().N, sig.S)
	out, err := asn1.Marshal(sig)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func signatureCases(t *testing.T) []signatureVector {
	edKey, p256Key := testKeys(t)
	var agr []agreementVector
	readJSON(t, "testdata/agreement-version.json", &agr)
	message := agr[0].Messages[1].Message // the custodian's approval of version 1
	b64 := base64.StdEncoding.EncodeToString

	edPub := edKey.Public().(ed25519.PublicKey)
	edSig := ed25519.Sign(edKey, []byte(message))
	raw, err := ecdsa.SignASN1(rand.Reader, p256Key, sha256Sum(message))
	if err != nil {
		t.Fatal(err)
	}
	low, err := NormalizeECDSA(raw)
	if err != nil {
		t.Fatal(err)
	}
	p256Pub := spki(t, &p256Key.PublicKey)
	tampered := append([]byte(nil), edSig...)
	tampered[0] ^= 1

	return []signatureVector{
		{Name: "ed25519, raw public key", Algorithm: Ed25519, PublicKey: b64(edPub), Message: message, Signature: b64(edSig), Valid: true},
		{Name: "ed25519, DER public key", Algorithm: Ed25519, PublicKey: b64(spki(t, edPub)), Message: message, Signature: b64(edSig), Valid: true},
		{Name: "ed25519, other message", Algorithm: Ed25519, PublicKey: b64(edPub), Message: message + ".", Signature: b64(edSig), Reason: "bad"},
		{Name: "ed25519, tampered signature", Algorithm: Ed25519, PublicKey: b64(edPub), Message: message, Signature: b64(tampered), Reason: "bad"},
		{Name: "p-256, low s", Algorithm: ECDSAP256SHA256, PublicKey: b64(p256Pub), Message: message, Signature: b64(low), Valid: true},
		{Name: "p-256, high s", Algorithm: ECDSAP256SHA256, PublicKey: b64(p256Pub), Message: message, Signature: b64(flipS(t, low)), Reason: "noncanonical"},
		{Name: "p-256, other message", Algorithm: ECDSAP256SHA256, PublicKey: b64(p256Pub), Message: message + ".", Signature: b64(low), Reason: "bad"},
		{Name: "p-256 signature checked as ed25519", Algorithm: Ed25519, PublicKey: b64(edPub), Message: message, Signature: b64(low), Reason: "bad"},
		{Name: "unknown algorithm", Algorithm: "rsa-pss-sha256", PublicKey: b64(p256Pub), Message: message, Signature: b64(low), Reason: "unsupported"},
	}
}

func sha256Sum(message string) []byte {
	sum := sha256.Sum256([]byte(message))
	return sum[:]
}

func TestSignatureVectors(t *testing.T) {
	const path = "testdata/signatures.json"
	if *update {
		writeJSON(t, path, signatureCases(t))
	}
	var vectors []signatureVector
	readJSON(t, path, &vectors)
	if len(vectors) != len(signatureCases(t)) {
		t.Fatalf("%s is out of date -- run with -update", path)
	}
	for _, vec := range vectors {
		t.Run(vec.Name, func(t *testing.T) {
			pub, _ := base64.StdEncoding.DecodeString(vec.PublicKey)
			sig, _ := base64.StdEncoding.DecodeString(vec.Signature)
			err := Verify(vec.Algorithm, pub, vec.Message, sig)
			switch {
			case vec.Valid:
				if err != nil {
					t.Fatalf("valid signature refused: %v", err)
				}
			case vec.Reason == "bad":
				if !errors.Is(err, ErrBadSignature) {
					t.Fatalf("err = %v, want ErrBadSignature", err)
				}
			case vec.Reason == "noncanonical":
				if !errors.Is(err, ErrNonCanonicalSignature) {
					t.Fatalf("err = %v, want ErrNonCanonicalSignature", err)
				}
			case vec.Reason == "unsupported":
				if !errors.Is(err, ErrUnsupportedAlgorithm) {
					t.Fatalf("err = %v, want ErrUnsupportedAlgorithm", err)
				}
			default:
				t.Fatalf("vector has no outcome")
			}
		})
	}
}

func TestNormalizeECDSA(t *testing.T) {
	_, key := testKeys(t)
	pub := spki(t, &key.PublicKey)
	for i := 0; i < 16; i++ { // both halves of s turn up within a few signatures
		raw, err := ecdsa.SignASN1(rand.Reader, key, sha256Sum("m"))
		if err != nil {
			t.Fatal(err)
		}
		for _, sig := range [][]byte{raw, flipS(t, raw)} {
			norm, err := NormalizeECDSA(sig)
			if err != nil {
				t.Fatal(err)
			}
			if err := Verify(ECDSAP256SHA256, pub, "m", norm); err != nil {
				t.Fatalf("normalized signature refused: %v", err)
			}
		}
	}
	if _, err := NormalizeECDSA([]byte{0x30, 0x00}); err == nil {
		t.Error("empty sequence accepted")
	}
}

func TestVerify_KeyRefusals(t *testing.T) {
	edKey, p256Key := testKeys(t)
	if err := Verify(Ed25519, []byte{1, 2, 3}, "m", make([]byte, 64)); err == nil {
		t.Error("short ed25519 key accepted")
	}
	if err := Verify(Ed25519, spki(t, &p256Key.PublicKey), "m", make([]byte, 64)); err == nil {
		t.Error("P-256 key accepted as ed25519")
	}
	if err := Verify(ECDSAP256SHA256, spki(t, edKey.Public()), "m", []byte{0x30, 0x00}); err == nil {
		t.Error("ed25519 key accepted as P-256")
	}
	p384, _ := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err := Verify(ECDSAP256SHA256, spki(t, &p384.PublicKey), "m", []byte{0x30, 0x00}); err == nil {
		t.Error("P-384 key accepted as P-256")
	}
	// Loose DER (a long-form length for a short sequence) is refused even
	// when the numbers inside verify.
	raw, _ := ecdsa.SignASN1(rand.Reader, p256Key, sha256Sum("m"))
	low, _ := NormalizeECDSA(raw)
	loose := append([]byte{0x30, 0x81, low[1]}, low[2:]...)
	if err := Verify(ECDSAP256SHA256, spki(t, &p256Key.PublicKey), "m", loose); err == nil {
		t.Error("loose DER accepted")
	}
}
