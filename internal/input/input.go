package input

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

// Resolve parses the argument as either an inline JSON string or a file path.
// Inline JSON is detected when the argument starts with '{', '[', '"',
// a digit, '-', or is one of the JSON literals (true, false, null).
// Otherwise, it's treated as a file path.
func Resolve(arg string) (interface{}, error) {
	trimmed := strings.TrimSpace(arg)
	if len(trimmed) == 0 {
		return nil, fmt.Errorf("empty input")
	}

	var data []byte
	if _, err := os.Stat(arg); err == nil {
		data, err = os.ReadFile(arg)
		if err != nil {
			return nil, fmt.Errorf("cannot read file %q: %w", arg, err)
		}
	} else if looksLikeJSON(trimmed) {
		// Inline JSON string
		data = []byte(trimmed)
	} else {
		// File path
		var err error
		data, err = os.ReadFile(arg)
		if err != nil {
			return nil, fmt.Errorf("cannot read file %q: %w", arg, err)
		}
	}

	return parseJSON(data)
}

// looksLikeJSON returns true if the trimmed string appears to be inline JSON
// rather than a file path.
func looksLikeJSON(s string) bool {
	c := s[0]
	if c == '{' || c == '[' || c == '"' {
		return true
	}
	// JSON numbers start with a digit or '-'
	if c == '-' || (c >= '0' && c <= '9') {
		return true
	}
	// JSON literals
	return s == "true" || s == "false" || s == "null"
}

func parseJSON(data []byte) (interface{}, error) {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()

	var result interface{}
	if err := decoder.Decode(&result); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		if err == nil {
			return nil, fmt.Errorf("invalid JSON: multiple JSON values")
		}
		return nil, fmt.Errorf("invalid JSON: trailing content: %w", err)
	}
	return result, nil
}
