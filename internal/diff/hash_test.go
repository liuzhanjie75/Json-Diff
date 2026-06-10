package diff

import (
	"encoding/json"
	"testing"
)

func TestCanonicalJSON_Null(t *testing.T) {
	result := CanonicalJSON(nil)
	if result != "null" {
		t.Errorf("expected 'null', got %q", result)
	}
}

func TestCanonicalJSON_String(t *testing.T) {
	result := CanonicalJSON("hello")
	if result != `"hello"` {
		t.Errorf("expected '\"hello\"', got %q", result)
	}
}

func TestCanonicalJSON_StringWithEscape(t *testing.T) {
	result := CanonicalJSON("he\"llo")
	if result != `"he\"llo"` {
		t.Errorf("expected escaped string, got %q", result)
	}
}

func TestCanonicalJSON_Number(t *testing.T) {
	result := CanonicalJSON(json.Number("42"))
	if result != "42" {
		t.Errorf("expected '42', got %q", result)
	}
}

func TestCanonicalJSON_Bool(t *testing.T) {
	if CanonicalJSON(true) != "true" {
		t.Error("expected 'true'")
	}
	if CanonicalJSON(false) != "false" {
		t.Error("expected 'false'")
	}
}

func TestCanonicalJSON_EmptyObject(t *testing.T) {
	result := CanonicalJSON(map[string]interface{}{})
	if result != "{}" {
		t.Errorf("expected '{}', got %q", result)
	}
}

func TestCanonicalJSON_EmptyArray(t *testing.T) {
	result := CanonicalJSON([]interface{}{})
	if result != "[]" {
		t.Errorf("expected '[]', got %q", result)
	}
}

func TestCanonicalJSON_ObjectKeysSorted(t *testing.T) {
	obj := map[string]interface{}{
		"z": json.Number("1"),
		"a": json.Number("2"),
		"m": json.Number("3"),
	}
	result := CanonicalJSON(obj)
	expected := `{"a":2,"m":3,"z":1}`
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestCanonicalJSON_ObjectKeysSortedEquivalence(t *testing.T) {
	obj1 := map[string]interface{}{"a": json.Number("1"), "b": json.Number("2")}
	obj2 := map[string]interface{}{"b": json.Number("2"), "a": json.Number("1")}

	c1 := CanonicalJSON(obj1)
	c2 := CanonicalJSON(obj2)
	if c1 != c2 {
		t.Errorf("canonical forms should be equal:\n  %q\n  %q", c1, c2)
	}
}

func TestCanonicalJSON_NestedObject(t *testing.T) {
	obj := map[string]interface{}{
		"outer": map[string]interface{}{
			"z": json.Number("1"),
			"a": json.Number("2"),
		},
	}
	result := CanonicalJSON(obj)
	expected := `{"outer":{"a":2,"z":1}}`
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestCanonicalJSON_Array(t *testing.T) {
	arr := []interface{}{json.Number("1"), "two", true, nil}
	result := CanonicalJSON(arr)
	expected := `[1,"two",true,null]`
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestHashJSON_Deterministic(t *testing.T) {
	obj := map[string]interface{}{"a": json.Number("1"), "b": json.Number("2")}
	h1 := HashJSON(obj)
	h2 := HashJSON(obj)
	if h1 != h2 {
		t.Errorf("HashJSON should be deterministic: %q != %q", h1, h2)
	}
}

func TestHashJSON_DifferentValues(t *testing.T) {
	obj1 := map[string]interface{}{"a": json.Number("1")}
	obj2 := map[string]interface{}{"a": json.Number("2")}
	if HashJSON(obj1) == HashJSON(obj2) {
		t.Error("different objects should have different hashes")
	}
}

func TestHashJSON_SameKeysDifferentOrder(t *testing.T) {
	obj1 := map[string]interface{}{"a": json.Number("1"), "b": json.Number("2")}
	obj2 := map[string]interface{}{"b": json.Number("2"), "a": json.Number("1")}
	if HashJSON(obj1) != HashJSON(obj2) {
		t.Error("objects with same keys in different order should have same hash")
	}
}

func TestHashJSON_Nil(t *testing.T) {
	h := HashJSON(nil)
	if h == "" {
		t.Error("HashJSON(nil) should return a non-empty hash")
	}
}
