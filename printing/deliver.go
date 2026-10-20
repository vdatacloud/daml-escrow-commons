package printing

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"path"
	"strings"
	"time"
)

// ZipMediaType is the media type of a bundle of artifacts.
const ZipMediaType = "application/zip"

// Artifact is one printed file: a rendered PDF, or any other document.
type Artifact struct {
	// Name is the file name inside a bundle and of a single download. It
	// must be a plain file name, unique among the artifacts delivered.
	Name      string
	MediaType string
	Data      []byte
	// Modified is recorded in a bundle; set it from the artifact's own
	// facts (e.g. when it was ratified), never the time of printing, so a
	// bundle of the same artifacts is the same bytes.
	Modified time.Time
}

// Delivery is what a print request downloads.
type Delivery struct {
	Name      string
	MediaType string
	Data      []byte
}

var ErrNothingToPrint = errors.New("printing: no artifacts selected")

// Deliver is the platform's print rule: one artifact is delivered as itself
// (a PDF stays a PDF); two or more are bundled into a ZIP named name (".zip"
// is added), in the order given. The bundle is deterministic -- entries are
// stored uncompressed (PDFs already are) with each artifact's own Modified
// time -- so the same artifacts always give the same bytes and a stored
// bundle can be checked by its hash.
func Deliver(name string, artifacts []Artifact) (Delivery, error) {
	switch len(artifacts) {
	case 0:
		return Delivery{}, ErrNothingToPrint
	case 1:
		a := artifacts[0]
		if err := checkName(a.Name); err != nil {
			return Delivery{}, err
		}
		return Delivery{Name: a.Name, MediaType: a.MediaType, Data: a.Data}, nil
	}
	data, err := Zip(artifacts)
	if err != nil {
		return Delivery{}, err
	}
	if !strings.HasSuffix(name, ".zip") {
		name += ".zip"
	}
	return Delivery{Name: name, MediaType: ZipMediaType, Data: data}, nil
}

// Zip bundles artifacts deterministically (see Deliver).
func Zip(artifacts []Artifact) ([]byte, error) {
	var buf bytes.Buffer
	w := zip.NewWriter(&buf)
	seen := map[string]bool{}
	for _, a := range artifacts {
		if err := checkName(a.Name); err != nil {
			return nil, err
		}
		if seen[a.Name] {
			return nil, fmt.Errorf("printing: two artifacts are named %q", a.Name)
		}
		seen[a.Name] = true
		h := &zip.FileHeader{Name: a.Name, Method: zip.Store}
		if !a.Modified.IsZero() {
			h.Modified = a.Modified.UTC()
		}
		f, err := w.CreateHeader(h)
		if err != nil {
			return nil, fmt.Errorf("printing: bundle %s: %w", a.Name, err)
		}
		if _, err := f.Write(a.Data); err != nil {
			return nil, fmt.Errorf("printing: bundle %s: %w", a.Name, err)
		}
	}
	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("printing: bundle: %w", err)
	}
	return buf.Bytes(), nil
}

func checkName(name string) error {
	if name == "" || name != path.Base(name) || strings.ContainsAny(name, `\/`) || name == "." || name == ".." {
		return fmt.Errorf("printing: %q is not a plain file name", name)
	}
	return nil
}
