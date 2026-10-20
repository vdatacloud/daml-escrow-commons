//go:build typst

package printing

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"testing"
	"time"
)

// Run with `go test -tags typst ./printing/` (needs typst 0.15.1 on
// PATH, e.g. `brew install typst`).
const pinned = "0.15.1"

func TestTypst_ByteIdenticalAndEscaped(t *testing.T) {
	ctx := context.Background()
	r, err := NewTypst(ctx, "typst", pinned)
	if err != nil {
		t.Fatal(err)
	}
	src := []byte(`#set text(font: "Source Serif 4")
#show heading: set text(font: "Plus Jakarta Sans")
= ` + Str("1. Parties *and* #roles") + `
#par[` + Str(`A "quoted" \ $5 #x [y] _z_`) + `]
`)
	at := time.Unix(1791900000, 0)
	a, err := r.Render(ctx, src, at)
	if err != nil {
		t.Fatal(err)
	}
	b, _ := r.Render(ctx, src, at)
	if !bytes.Equal(a, b) {
		t.Fatal("two renders of one source differ")
	}
	later, _ := r.Render(ctx, src, at.Add(time.Hour))
	if bytes.Equal(a, later) {
		t.Fatal("the creation time should be part of the output")
	}
	if r.Name() != "typst "+pinned {
		t.Fatal(r.Name())
	}
	if _, err := r.Render(ctx, []byte(`#nonexistent()`), at); err == nil {
		t.Fatal("a broken source rendered")
	}
}

func TestTypst_RefusesOtherVersion(t *testing.T) {
	if _, err := NewTypst(context.Background(), "typst", "0.0.1"); err == nil {
		t.Fatal("a different typst version was accepted")
	}
}

// goldenSource renders to goldenPDF's hash on every platform the CMS renders
// on (macOS arm64 locally, Linux x86-64 in CI and deployment): a signed
// document must re-render to the same bytes wherever it is re-rendered.
const goldenSource = `#set document(title: "Golden", author: ("daml-escrow-cms",))
#set page(paper: "us-letter", margin: 1in)
#set text(font: "Source Serif 4", size: 11pt, lang: "en")
#show heading: set text(font: "Plus Jakarta Sans")
= ` + `#"1. Golden clause"` + `
#par[#"The Cash Leg of "#strong[#"1,000.00 USD"]#" is released on delivery."]
#list([#"one"], [#"two"])
`

const goldenPDF = "f53898565e4b120a7eec287856f005ed3cb28b77230634cae8586e2e70a82ed2"

func TestTypst_GoldenHashAcrossPlatforms(t *testing.T) {
	r, err := NewTypst(context.Background(), "typst", pinned)
	if err != nil {
		t.Fatal(err)
	}
	pdf, err := r.Render(context.Background(), []byte(goldenSource), time.Unix(1791900000, 0))
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(pdf)
	if got := hex.EncodeToString(sum[:]); got != goldenPDF {
		t.Fatalf("golden PDF hash is %s on this platform, want %s: typst output differs across platforms", got, goldenPDF)
	}
}
