// Package schema provides a directory-of-JSON-Schema-files validator shared
// across daml-escrow (validating EscrowMetadata payloads) and daml-escrow-cms
// (validating structurally-extracted contract terms before draft creation).
// Both repos validate against the same schema authority
// (daml-escrow/architecture/schemas/*.json); this package is the one
// implementation both depend on rather than each wrapping gojsonschema
// independently.
package schema

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

// Registry holds compiled JSON Schemas keyed by type name (the schema's
// filename without the .json extension, e.g. "escrow" for "escrow.json").
type Registry struct {
	schemas map[string]*gojsonschema.Schema
}

// LoadDirectory compiles every *.json file in dir into the Registry, keyed by
// filename without extension. Returns an error if the directory can't be read
// or any schema fails to compile.
func LoadDirectory(dir string) (*Registry, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("schema: failed to read directory %q: %w", dir, err)
	}

	r := &Registry{schemas: make(map[string]*gojsonschema.Schema)}

	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
			continue
		}

		typeName := strings.TrimSuffix(file.Name(), ".json")
		absPath, err := filepath.Abs(filepath.Join(dir, file.Name()))
		if err != nil {
			return nil, fmt.Errorf("schema: failed to resolve path for %q: %w", file.Name(), err)
		}

		loader := gojsonschema.NewReferenceLoader("file://" + absPath)
		compiled, err := gojsonschema.NewSchema(loader)
		if err != nil {
			return nil, fmt.Errorf("schema: failed to compile %q: %w", file.Name(), err)
		}

		r.schemas[typeName] = compiled
	}

	return r, nil
}

// Types returns the type names currently loaded, for diagnostics/health
// endpoints.
func (r *Registry) Types() []string {
	types := make([]string, 0, len(r.schemas))
	for t := range r.schemas {
		types = append(types, t)
	}
	return types
}

// ValidationError aggregates the individual field-level failures
// gojsonschema reports, so callers get one typed error rather than parsing
// strings.
type ValidationError struct {
	TypeName string
	Failures []string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("schema: %q validation failed: %s", e.TypeName, strings.Join(e.Failures, "; "))
}

// ErrUnknownType is returned by Validate when typeName has no loaded schema.
type ErrUnknownType struct {
	TypeName string
}

func (e *ErrUnknownType) Error() string {
	return fmt.Sprintf("schema: unknown type %q", e.TypeName)
}

// Validate checks data (raw JSON bytes) against the schema registered under
// typeName. Returns *ValidationError on schema-rule failures and
// *ErrUnknownType if typeName was never loaded — both are distinguishable via
// errors.As by callers that want to react differently (e.g. treat an unknown
// type as a 400 vs. a 404).
func (r *Registry) Validate(typeName string, data []byte) error {
	compiled, ok := r.schemas[typeName]
	if !ok {
		return &ErrUnknownType{TypeName: typeName}
	}

	result, err := compiled.Validate(gojsonschema.NewBytesLoader(data))
	if err != nil {
		return fmt.Errorf("schema: failed to run validation for %q: %w", typeName, err)
	}

	if !result.Valid() {
		failures := make([]string, 0, len(result.Errors()))
		for _, e := range result.Errors() {
			failures = append(failures, e.String())
		}
		return &ValidationError{TypeName: typeName, Failures: failures}
	}

	return nil
}
