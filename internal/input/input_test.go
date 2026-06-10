package input_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/zhanjie/jsondiff/internal/input"
)

func TestResolve_InlineObject(t *testing.T) {
	result, err := input.Resolve(`{"a":1}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("expected map, got %T", result)
	}
	if len(m) != 1 {
		t.Errorf("expected 1 key, got %d", len(m))
	}
}

func TestResolve_InlineArray(t *testing.T) {
	result, err := input.Resolve(`[1,2,3]`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	arr, ok := result.([]interface{})
	if !ok {
		t.Fatalf("expected []interface{}, got %T", result)
	}
	if len(arr) != 3 {
		t.Errorf("expected 3 elements, got %d", len(arr))
	}
}

func TestResolve_InlineWithWhitespace(t *testing.T) {
	result, err := input.Resolve(`  {"a":1}  `)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Error("expected non-nil result")
	}
}

func TestResolve_FilePath(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.json")
	os.WriteFile(path, []byte(`{"key":"value"}`), 0644)

	result, err := input.Resolve(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("expected map, got %T", result)
	}
	if m["key"] != "value" {
		t.Errorf("expected key=value, got %v", m["key"])
	}
}

func TestResolve_FileNotFound(t *testing.T) {
	_, err := input.Resolve("nonexistent_file.json")
	if err == nil {
		t.Error("expected error for non-existent file")
	}
}

func TestResolve_EmptyInput(t *testing.T) {
	_, err := input.Resolve("")
	if err == nil {
		t.Error("expected error for empty input")
	}
}

func TestResolve_InvalidJSON(t *testing.T) {
	_, err := input.Resolve(`{invalid}`)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestResolve_RejectsTrailingJSONContent(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{name: "trailing text", input: `{"a":1} trailing`},
		{name: "multiple values", input: `{"a":1} {"b":2}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := input.Resolve(tt.input); err == nil {
				t.Fatalf("expected error for %q", tt.input)
			}
		})
	}
}

func TestResolve_PrefersExistingNumericFilename(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "2024.json")
	if err := os.WriteFile(path, []byte(`{"source":"file"}`), 0644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("get working directory: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("change working directory: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Errorf("restore working directory: %v", err)
		}
	})

	result, err := input.Resolve("2024.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("expected map, got %T", result)
	}
	if m["source"] != "file" {
		t.Fatalf("expected file contents, got %v", m)
	}
}

func TestResolve_UseNumber(t *testing.T) {
	result, err := input.Resolve(`{"big":99999999999999999}`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m := result.(map[string]interface{})
	_, isNum := m["big"].(json.Number)
	if !isNum {
		t.Errorf("expected json.Number, got %T", m["big"])
	}
}

func TestResolve_NullJSON(t *testing.T) {
	result, err := input.Resolve(`null`)
	if err != nil {
		t.Fatalf("unexpected error for 'null': %v", err)
	}
	if result != nil {
		t.Errorf("expected nil result for 'null', got %v", result)
	}
}

func TestResolve_TrueJSON(t *testing.T) {
	result, err := input.Resolve(`true`)
	if err != nil {
		t.Fatalf("unexpected error for 'true': %v", err)
	}
	if result != true {
		t.Errorf("expected true, got %v", result)
	}
}

func TestResolve_FalseJSON(t *testing.T) {
	result, err := input.Resolve(`false`)
	if err != nil {
		t.Fatalf("unexpected error for 'false': %v", err)
	}
	if result != false {
		t.Errorf("expected false, got %v", result)
	}
}

func TestResolve_NumberJSON(t *testing.T) {
	result, err := input.Resolve(`42`)
	if err != nil {
		t.Fatalf("unexpected error for '42': %v", err)
	}
	num, ok := result.(interface{ String() string })
	if !ok {
		t.Fatalf("expected json.Number, got %T", result)
	}
	if num.String() != "42" {
		t.Errorf("expected '42', got %q", num.String())
	}
}

func TestResolve_StringJSON(t *testing.T) {
	result, err := input.Resolve(`"hello"`)
	if err != nil {
		t.Fatalf("unexpected error for '\"hello\"': %v", err)
	}
	s, ok := result.(string)
	if !ok {
		t.Fatalf("expected string, got %T", result)
	}
	if s != "hello" {
		t.Errorf("expected 'hello', got %q", s)
	}
}
