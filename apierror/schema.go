package apierror

import _ "embed"

//go:embed schema/error.schema.json
var schemaJSON []byte

// Schema is the JSON Schema (draft-07) of the canonical error envelope --
// the contract shared with @vdatacloud/cx-commons (sdk/api-error) and any
// non-Go client.
func Schema() []byte {
	out := make([]byte, len(schemaJSON))
	copy(out, schemaJSON)
	return out
}
