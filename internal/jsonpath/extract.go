package jsonpath

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/tidwall/gjson"
)

// Extract extracts a sub-value from a JSON value using the given path.
// Path syntax follows gjson conventions: "users.0.name", "items.#(id==1)", etc.
func Extract(data interface{}, path string) (interface{}, error) {
	// Serialize the interface{} back to JSON bytes
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize JSON: %w", err)
	}

	result := gjson.GetBytes(bytes, path)
	if !result.Exists() {
		return nil, fmt.Errorf("path %q does not exist in JSON", path)
	}

	// Parse the result back to interface{} with UseNumber
	return parseResult(result)
}

func parseResult(result gjson.Result) (interface{}, error) {
	// gjson.Result.Value() returns the native Go value
	// but we need to use UseNumber for number precision
	raw := result.Raw
	if raw == "" {
		return result.Value(), nil
	}

	var v interface{}
	decoder := json.NewDecoder(strings.NewReader(raw))
	decoder.UseNumber()
	if err := decoder.Decode(&v); err != nil {
		// Fallback to gjson's native value
		return result.Value(), nil
	}
	return v, nil
}
