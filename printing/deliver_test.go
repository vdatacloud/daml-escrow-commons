package printing

import (
	"archive/zip"
	"bytes"
	"errors"
	"io"
	"testing"
	"time"
)

func artifacts() []Artifact {
	at := time.Date(2026, 10, 20, 12, 0, 0, 0, time.UTC)
	return []Artifact{
		{Name: "agreement.pdf", MediaType: PDFMediaType, Data: []byte("%PDF-1 a"), Modified: at},
		{Name: "amendment-2.pdf", MediaType: PDFMediaType, Data: []byte("%PDF-1 b"), Modified: at.Add(time.Hour)},
	}
}

// One artifact is delivered as itself.
func TestDeliver_One(t *testing.T) {
	a := artifacts()[0]
	d, err := Deliver("record", []Artifact{a})
	if err != nil || d.Name != a.Name || d.MediaType != PDFMediaType || !bytes.Equal(d.Data, a.Data) {
		t.Fatalf("%+v %v", d, err)
	}
	if _, err := Deliver("record", nil); !errors.Is(err, ErrNothingToPrint) {
		t.Fatal(err)
	}
}

// Two or more are a ZIP, in order, with their own times, the same bytes
// every time.
func TestDeliver_ManyIsDeterministicZip(t *testing.T) {
	d, err := Deliver("record", artifacts())
	if err != nil || d.Name != "record.zip" || d.MediaType != ZipMediaType {
		t.Fatalf("%+v %v", d, err)
	}
	again, _ := Deliver("record.zip", artifacts())
	if !bytes.Equal(d.Data, again.Data) || again.Name != "record.zip" {
		t.Fatal("bundle is not deterministic")
	}
	r, err := zip.NewReader(bytes.NewReader(d.Data), int64(len(d.Data)))
	if err != nil || len(r.File) != 2 {
		t.Fatalf("%v %d", err, len(r.File))
	}
	for i, f := range r.File {
		want := artifacts()[i]
		rc, _ := f.Open()
		got, _ := io.ReadAll(rc)
		_ = rc.Close()
		if f.Name != want.Name || f.Method != zip.Store || !f.Modified.Equal(want.Modified) || !bytes.Equal(got, want.Data) {
			t.Fatalf("entry %d: %s %d %v", i, f.Name, f.Method, f.Modified)
		}
	}
}

func TestDeliver_RefusesBadNames(t *testing.T) {
	for _, name := range []string{"", "../x.pdf", "a/b.pdf", `a\b.pdf`, ".."} {
		as := artifacts()
		as[0].Name = name
		if _, err := Deliver("record", as); err == nil {
			t.Errorf("%q accepted", name)
		}
		if _, err := Deliver("record", as[:1]); err == nil {
			t.Errorf("single %q accepted", name)
		}
	}
	dup := artifacts()
	dup[1].Name = dup[0].Name
	if _, err := Deliver("record", dup); err == nil {
		t.Fatal("duplicate names accepted")
	}
}
