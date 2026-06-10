package diff

import (
	"encoding/json"
	"testing"
)

func TestObjectSimilarity_Identical(t *testing.T) {
	a := map[string]interface{}{"x": 1, "y": 2, "z": 3}
	b := map[string]interface{}{"x": 1, "y": 2, "z": 3}
	sim := objectSimilarity(a, b)
	if sim != 1.0 {
		t.Errorf("expected 1.0 for identical key sets, got %f", sim)
	}
}

func TestObjectSimilarity_SameKeysDifferentValues(t *testing.T) {
	a := map[string]interface{}{"x": 1, "y": 2}
	b := map[string]interface{}{"x": 100, "y": 200}
	sim := objectSimilarity(a, b)
	if sim != 1.0 {
		t.Errorf("expected 1.0 (same keys), got %f", sim)
	}
}

func TestObjectSimilarity_PartialOverlap(t *testing.T) {
	a := map[string]interface{}{"x": 1, "y": 2}
	b := map[string]interface{}{"x": 1, "z": 3}
	sim := objectSimilarity(a, b)
	expected := 1.0 / 3.0
	if sim < expected-0.01 || sim > expected+0.01 {
		t.Errorf("expected ~0.333, got %f", sim)
	}
}

func TestObjectSimilarity_NoOverlap(t *testing.T) {
	a := map[string]interface{}{"x": 1}
	b := map[string]interface{}{"y": 2}
	sim := objectSimilarity(a, b)
	if sim != 0.0 {
		t.Errorf("expected 0.0 for no overlap, got %f", sim)
	}
}

func TestObjectSimilarity_BothEmpty(t *testing.T) {
	sim := objectSimilarity(map[string]interface{}{}, map[string]interface{}{})
	if sim != 1.0 {
		t.Errorf("expected 1.0 for both empty, got %f", sim)
	}
}

func TestObjectSimilarity_MostlySameKeys(t *testing.T) {
	a := map[string]interface{}{"a": 1, "b": 2, "c": 3, "d": 4, "e": 5}
	b := map[string]interface{}{"a": 1, "b": 2, "c": 3, "d": 4, "f": 6}
	sim := objectSimilarity(a, b)
	if sim < similarityThreshold {
		t.Errorf("expected similarity >= threshold, got %f", sim)
	}
}

func TestKeyStringValue_String(t *testing.T) {
	if keyStringValue("hello") != `"hello"` {
		t.Error(`expected '"hello"'`)
	}
}

func TestKeyStringValue_Number(t *testing.T) {
	if keyStringValue(json.Number("42")) != "42" {
		t.Error("expected '42'")
	}
}

func TestKeyStringValue_Bool(t *testing.T) {
	if keyStringValue(true) != "true" {
		t.Error("expected 'true'")
	}
	if keyStringValue(false) != "false" {
		t.Error("expected 'false'")
	}
}

func TestKeyStringValue_PreservesJSONTypes(t *testing.T) {
	tests := []struct {
		name string
		a    interface{}
		b    interface{}
	}{
		{name: "string and number", a: "1", b: json.Number("1")},
		{name: "string and boolean", a: "true", b: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if keyStringValue(tt.a) == keyStringValue(tt.b) {
				t.Fatalf("expected distinct key encodings for %T and %T", tt.a, tt.b)
			}
		})
	}
}

func TestKeyStringValue_Nil(t *testing.T) {
	if keyStringValue(nil) != "" {
		t.Error("expected empty string for nil")
	}
}

func TestFilterUnmatched_AllUnmatched(t *testing.T) {
	result := filterUnmatched([]int{0, 1, 2}, map[int]bool{})
	if len(result) != 3 {
		t.Errorf("expected 3, got %d", len(result))
	}
}

func TestFilterUnmatched_SomeMatched(t *testing.T) {
	result := filterUnmatched([]int{0, 1, 2}, map[int]bool{1: true})
	if len(result) != 2 {
		t.Errorf("expected 2, got %d", len(result))
	}
	if result[0] != 0 || result[1] != 2 {
		t.Errorf("expected [0,2], got %v", result)
	}
}

func TestFilterUnmatched_AllMatched(t *testing.T) {
	result := filterUnmatched([]int{0, 1}, map[int]bool{0: true, 1: true})
	if len(result) != 0 {
		t.Errorf("expected 0, got %d", len(result))
	}
}

func TestFilterUnmatched_EmptyInput(t *testing.T) {
	result := filterUnmatched(nil, map[int]bool{})
	if len(result) != 0 {
		t.Errorf("expected 0, got %d", len(result))
	}
}
