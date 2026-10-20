package printing

import (
	"bytes"
	"context"
	"testing"
	"time"
)

func TestCheckFonts(t *testing.T) {
	if err := CheckFonts(); err != nil {
		t.Fatal(err)
	}
}

func TestStr(t *testing.T) {
	cases := map[string]string{
		`plain`:            `#"plain"`,
		`a "q" \ b`:        `#"a \"q\" \\ b"`,
		"line\r\nnext\tx":  `#"line\nnext\tx"`,
		`= #x *y* $z$ [w]`: `#"= #x *y* $z$ [w]"`,
	}
	for in, want := range cases {
		if got := Str(in); got != want {
			t.Errorf("Str(%q) = %s, want %s", in, got, want)
		}
	}
	if Content("x") != `[#"x"]` {
		t.Fatal(Content("x"))
	}
}

func TestFakeIsDeterministic(t *testing.T) {
	at := time.Unix(1791900000, 0)
	a, _ := Fake{}.Render(context.Background(), []byte("source"), at)
	b, _ := Fake{}.Render(context.Background(), []byte("source"), at)
	c, _ := Fake{}.Render(context.Background(), []byte("other"), at)
	if !bytes.Equal(a, b) || bytes.Equal(a, c) || !bytes.HasPrefix(a, []byte("%PDF-")) {
		t.Fatal("fake renderer is not deterministic per source")
	}
}
