package jsonpath_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/zhanjie/jsondiff/internal/jsonpath"
)

func parseTestJSON(s string) interface{} {
	var v interface{}
	dec := json.NewDecoder(strings.NewReader(s))
	dec.UseNumber()
	if err := dec.Decode(&v); err != nil {
		panic(err)
	}
	return v
}

func TestExtract_SimpleKey(t *testing.T) {
	data := parseTestJSON(`{"name":"alice","age":30}`)
	result, err := jsonpath.Extract(data, "name")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "alice" {
		t.Errorf("expected 'alice', got %v", result)
	}
}

func TestExtract_NestedKey(t *testing.T) {
	data := parseTestJSON(`{"config":{"debug":true,"timeout":30}}`)
	result, err := jsonpath.Extract(data, "config.debug")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != true {
		t.Errorf("expected true, got %v", result)
	}
}

func TestExtract_ArrayIndex(t *testing.T) {
	data := parseTestJSON(`{"users":["alice","bob","charlie"]}`)
	result, err := jsonpath.Extract(data, "users.1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "bob" {
		t.Errorf("expected 'bob', got %v", result)
	}
}

func TestExtract_ObjectInArray(t *testing.T) {
	data := parseTestJSON(`{"users":[{"name":"alice"},{"name":"bob"}]}`)
	result, err := jsonpath.Extract(data, "users.0.name")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "alice" {
		t.Errorf("expected 'alice', got %v", result)
	}
}

func TestExtract_ReturnsObject(t *testing.T) {
	data := parseTestJSON(`{"config":{"a":1,"b":2}}`)
	result, err := jsonpath.Extract(data, "config")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("expected map, got %T", result)
	}
	if len(m) != 2 {
		t.Errorf("expected 2 keys, got %d", len(m))
	}
}

func TestExtract_ReturnsArray(t *testing.T) {
	data := parseTestJSON(`{"items":[1,2,3]}`)
	result, err := jsonpath.Extract(data, "items")
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

func TestExtract_PathNotFound(t *testing.T) {
	data := parseTestJSON(`{"a":1}`)
	_, err := jsonpath.Extract(data, "nonexistent")
	if err == nil {
		t.Error("expected error for non-existent path")
	}
}

func TestExtract_DeepPathNotFound(t *testing.T) {
	data := parseTestJSON(`{"a":{"b":1}}`)
	_, err := jsonpath.Extract(data, "a.c.d")
	if err == nil {
		t.Error("expected error for non-existent deep path")
	}
}

func TestExtract_NumberPreserved(t *testing.T) {
	data := parseTestJSON(`{"val":42}`)
	result, err := jsonpath.Extract(data, "val")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	num, ok := result.(json.Number)
	if !ok {
		t.Fatalf("expected json.Number, got %T", result)
	}
	if num.String() != "42" {
		t.Errorf("expected '42', got %q", num.String())
	}
}

func TestExtract_RootPath(t *testing.T) {
	data := parseTestJSON(`{"a":1}`)
	_, err := jsonpath.Extract(data, "")
	if err != nil {
		t.Logf("empty path returned error (expected): %v", err)
	}
}
