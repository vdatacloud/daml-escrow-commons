// Package printing is the platform's printed output: a Typst source compiled
// to a PDF, and printed artifacts delivered as one PDF or a ZIP of several
// (Deliver). Rendering is deterministic: the same source, renderer version
// and creation time give the same bytes (daml-escrow-cms PLAN.md 81.4 T1). The fonts are pinned files embedded in the binary
// (fonts/, checked against fonts/SHA256SUMS), system fonts are ignored, and
// the PDF creation time is fixed by the caller, so a stored source can be
// re-rendered to the hash its signers saw.
package printing

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// SourceMediaType is the media type recorded for a Typst source.
const SourceMediaType = "text/x-typst"

// PDFMediaType is the media type of a rendered document.
const PDFMediaType = "application/pdf"

// ErrUnavailable means no renderer is configured.
var ErrUnavailable = errors.New("printing: no renderer is configured")

// Renderer compiles a Typst source to a PDF.
type Renderer interface {
	// Name identifies the renderer and its version, recorded with every
	// rendition (e.g. "typst 0.15.1").
	Name() string
	Render(ctx context.Context, source []byte, created time.Time) ([]byte, error)
}

//go:embed fonts/*.ttf fonts/SHA256SUMS
var fonts embed.FS

// CheckFonts verifies the embedded fonts against fonts/SHA256SUMS, so a
// silently replaced font file can't change rendered output.
func CheckFonts() error {
	sums, err := fonts.ReadFile("fonts/SHA256SUMS")
	if err != nil {
		return err
	}
	want := map[string]string{}
	sc := bufio.NewScanner(bytes.NewReader(sums))
	for sc.Scan() {
		f := strings.Fields(sc.Text())
		if len(f) == 2 {
			want[f[1]] = f[0]
		}
	}
	entries, err := fonts.ReadDir("fonts")
	if err != nil {
		return err
	}
	seen := 0
	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".ttf") {
			continue
		}
		data, _ := fonts.ReadFile("fonts/" + e.Name())
		sum := sha256.Sum256(data)
		if hex.EncodeToString(sum[:]) != want[e.Name()] {
			return fmt.Errorf("printing: font %s does not match fonts/SHA256SUMS", e.Name())
		}
		seen++
	}
	if seen != len(want) || seen == 0 {
		return fmt.Errorf("printing: fonts/SHA256SUMS lists %d fonts, %d are embedded", len(want), seen)
	}
	return nil
}

// Typst runs the typst CLI. Version pins the expected `typst --version`, so
// a different binary can't silently produce different bytes.
type Typst struct {
	path    string
	version string

	once    sync.Once
	fontDir string
	setup   error
}

// NewTypst checks that the binary at path is the expected version and that
// the embedded fonts are intact.
func NewTypst(ctx context.Context, path, version string) (*Typst, error) {
	if err := CheckFonts(); err != nil {
		return nil, err
	}
	out, err := exec.CommandContext(ctx, path, "--version").Output()
	if err != nil {
		return nil, fmt.Errorf("printing: run %s --version: %w", path, err)
	}
	got := strings.Fields(string(out))
	if len(got) < 2 || got[0] != "typst" || got[1] != version {
		return nil, fmt.Errorf("printing: %s is %q, want typst %s", path, strings.TrimSpace(string(out)), version)
	}
	return &Typst{path: path, version: version}, nil
}

func (t *Typst) Name() string { return "typst " + t.version }

func (t *Typst) fonts() (string, error) {
	t.once.Do(func() {
		dir, err := os.MkdirTemp("", "typst-fonts-")
		if err != nil {
			t.setup = err
			return
		}
		entries, _ := fonts.ReadDir("fonts")
		for _, e := range entries {
			if !strings.HasSuffix(e.Name(), ".ttf") {
				continue
			}
			data, _ := fonts.ReadFile("fonts/" + e.Name())
			if err := os.WriteFile(filepath.Join(dir, e.Name()), data, 0o600); err != nil {
				t.setup = err
				return
			}
		}
		t.fontDir = dir
	})
	return t.fontDir, t.setup
}

func (t *Typst) Render(ctx context.Context, source []byte, created time.Time) ([]byte, error) {
	dir, err := t.fonts()
	if err != nil {
		return nil, fmt.Errorf("printing: prepare fonts: %w", err)
	}
	cmd := exec.CommandContext(ctx, t.path, "compile",
		"--ignore-system-fonts", "--font-path", dir,
		"--creation-timestamp", fmt.Sprint(created.UTC().Unix()),
		"-", "-")
	cmd.Stdin = bytes.NewReader(source)
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("printing: typst compile: %w: %s", err, strings.TrimSpace(stderr.String()))
	}
	if !bytes.HasPrefix(stdout.Bytes(), []byte("%PDF-")) {
		return nil, errors.New("printing: typst produced no PDF")
	}
	return stdout.Bytes(), nil
}

// Fake is NOT a renderer: it returns a tiny placeholder PDF derived from the
// source's hash, for unit tests and dev stacks without typst. Its output is
// deterministic but is not a readable document.
type Fake struct{}

func (Fake) Name() string { return "fake" }

func (Fake) Render(_ context.Context, source []byte, created time.Time) ([]byte, error) {
	sum := sha256.Sum256(source)
	return []byte(fmt.Sprintf("%%PDF-1.7\n%% fake rendition of source %x at %d\n%%%%EOF\n", sum, created.UTC().Unix())), nil
}
