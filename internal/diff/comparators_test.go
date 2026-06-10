package diff

import (
	"encoding/json"
	"testing"
)

// --- NullComparator tests ---

func TestNullComparator_BothNil(t *testing.T) {
	c := &NullComparator{}
	diffs := c.Compare(nil, nil, "$", newCtx())
	if len(diffs) != 0 {
		t.Errorf("expected no diffs for nil==nil, got %d", len(diffs))
	}
}

func TestNullComparator_OldNil(t *testing.T) {
	c := &NullComparator{}
	diffs := c.Compare(nil, parseJSON(`42`), "$", newCtx())
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(diffs))
	}
	if diffs[0].Op != OpChanged {
		t.Errorf("expected OpChanged, got %v", diffs[0].Op)
	}
	if diffs[0].Path != "$" {
		t.Errorf("expected path $, got %s", diffs[0].Path)
	}
}

func TestNullComparator_NewNil(t *testing.T) {
	c := &NullComparator{}
	diffs := c.Compare(parseJSON(`"hello"`), nil, "$", newCtx())
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(diffs))
	}
	if diffs[0].Op != OpChanged {
		t.Errorf("expected OpChanged, got %v", diffs[0].Op)
	}
}

func TestNullComparator_NilToObject(t *testing.T) {
	c := &NullComparator{}
	diffs := c.Compare(nil, parseJSON(`{"a":1}`), "$", newCtx())
	if len(diffs) != 1 || diffs[0].Op != OpChanged {
		t.Error("expected OpChanged for nil鈫抩bject")
	}
}

func TestNullComparator_ObjectToNil(t *testing.T) {
	c := &NullComparator{}
	diffs := c.Compare(parseJSON(`{"a":1}`), nil, "$", newCtx())
	if len(diffs) != 1 || diffs[0].Op != OpChanged {
		t.Error("expected OpChanged for object鈫抧il")
	}
}

func TestNullComparator_NilToArray(t *testing.T) {
	c := &NullComparator{}
	diffs := c.Compare(nil, parseJSON(`[1,2]`), "$", newCtx())
	if len(diffs) != 1 || diffs[0].Op != OpChanged {
		t.Error("expected OpChanged for nil鈫抋rray")
	}
}

func TestNullComparator_NilToNull(t *testing.T) {
	c := &NullComparator{}
	diffs := c.Compare(nil, nil, "$.a", newCtx())
	if len(diffs) != 0 {
		t.Errorf("expected no diffs for nil鈫抧ull, got %d", len(diffs))
	}
}

// --- PrimitiveComparator tests ---

func TestPrimitiveComparator_EqualStrings(t *testing.T) {
	c := &PrimitiveComparator{}
	diffs := c.Compare("hello", "hello", "$", newCtx())
	if len(diffs) != 0 {
		t.Errorf("expected no diffs for equal strings, got %d", len(diffs))
	}
}

func TestPrimitiveComparator_DifferentStrings(t *testing.T) {
	c := &PrimitiveComparator{}
	diffs := c.Compare("hello", "world", "$", newCtx())
	if len(diffs) != 1 || diffs[0].Op != OpChanged {
		t.Errorf("expected OpChanged for different strings")
	}
}

func TestPrimitiveComparator_EqualNumbers(t *testing.T) {
	c := &PrimitiveComparator{}
	diffs := c.Compare(json.Number("42"), json.Number("42"), "$", newCtx())
	if len(diffs) != 0 {
		t.Errorf("expected no diffs for equal numbers, got %d", len(diffs))
	}
}

func TestPrimitiveComparator_DifferentNumbers(t *testing.T) {
	c := &PrimitiveComparator{}
	diffs := c.Compare(json.Number("42"), json.Number("43"), "$", newCtx())
	if len(diffs) != 1 || diffs[0].Op != OpChanged {
		t.Errorf("expected OpChanged for different numbers")
	}
}

func TestPrimitiveComparator_EqualBools(t *testing.T) {
	c := &PrimitiveComparator{}
	diffs := c.Compare(true, true, "$", newCtx())
	if len(diffs) != 0 {
		t.Errorf("expected no diffs for equal bools, got %d", len(diffs))
	}
}

func TestPrimitiveComparator_DifferentBools(t *testing.T) {
	c := &PrimitiveComparator{}
	diffs := c.Compare(true, false, "$", newCtx())
	if len(diffs) != 1 || diffs[0].Op != OpChanged {
		t.Errorf("expected OpChanged for different bools")
	}
}

func TestPrimitiveComparator_TypeMismatch_StringVsNumber(t *testing.T) {
	c := &PrimitiveComparator{}
	diffs := c.Compare("42", json.Number("42"), "$", newCtx())
	if len(diffs) != 1 || diffs[0].Op != OpChanged {
		t.Errorf("expected OpChanged for string vs number type mismatch")
	}
}

func TestPrimitiveComparator_TypeMismatch_BoolVsString(t *testing.T) {
	c := &PrimitiveComparator{}
	diffs := c.Compare(true, "true", "$", newCtx())
	if len(diffs) != 1 || diffs[0].Op != OpChanged {
		t.Errorf("expected OpChanged for bool vs string type mismatch")
	}
}

func TestPrimitiveComparator_TypeMismatch_ObjectVsString(t *testing.T) {
	c := &PrimitiveComparator{}
	obj := map[string]interface{}{"a": 1}
	diffs := c.Compare(obj, "string", "$", newCtx())
	if len(diffs) != 1 || diffs[0].Op != OpChanged {
		t.Errorf("expected OpChanged for object vs string type mismatch")
	}
}

func TestPrimitiveComparator_NumberPrecision(t *testing.T) {
	c := &PrimitiveComparator{}
	diffs := c.Compare(json.Number("1.0"), json.Number("1.0"), "$", newCtx())
	if len(diffs) != 0 {
		t.Errorf("expected no diffs for same number string, got %d", len(diffs))
	}

	diffs = c.Compare(json.Number("1.0"), json.Number("1.00"), "$", newCtx())
	if len(diffs) != 1 {
		t.Errorf("expected OpChanged for 1.0 vs 1.00 (different string repr)")
	}
}

func TestEqualScalars(t *testing.T) {
	tests := []struct {
		name string
		a, b interface{}
		want bool
	}{
		{"same string", "hello", "hello", true},
		{"diff string", "hello", "world", false},
		{"same bool", true, true, true},
		{"diff bool", true, false, false},
		{"same number", json.Number("1"), json.Number("1"), true},
		{"diff number", json.Number("1"), json.Number("2"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := equalScalars(tt.a, tt.b); got != tt.want {
				t.Errorf("EqualScalars(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

// --- ObjectComparator tests ---

func TestObjectComparator_IdenticalObjects(t *testing.T) {
	c := &ObjectComparator{}
	old := parseJSON(`{"a":1,"b":"x"}`)
	newVal := parseJSON(`{"a":1,"b":"x"}`)

	diffs := c.Compare(old, newVal, "$", newCtx())
	if len(diffs) != 0 {
		t.Errorf("expected no diffs, got %d", len(diffs))
	}
}

func TestObjectComparator_AddedKeys(t *testing.T) {
	c := &ObjectComparator{}
	old := parseJSON(`{}`)
	newVal := parseJSON(`{"a":1,"b":2}`)

	diffs := c.Compare(old, newVal, "$", newCtx())
	if len(diffs) != 2 {
		t.Fatalf("expected 2 diffs, got %d", len(diffs))
	}
	for _, d := range diffs {
		if d.Op != OpAdded {
			t.Errorf("expected OpAdded, got %v for path %s", d.Op, d.Path)
		}
	}
}

func TestObjectComparator_RemovedKeys(t *testing.T) {
	c := &ObjectComparator{}
	old := parseJSON(`{"a":1,"b":2}`)
	newVal := parseJSON(`{}`)

	diffs := c.Compare(old, newVal, "$", newCtx())
	if len(diffs) != 2 {
		t.Fatalf("expected 2 diffs, got %d", len(diffs))
	}
	for _, d := range diffs {
		if d.Op != OpRemoved {
			t.Errorf("expected OpRemoved, got %v for path %s", d.Op, d.Path)
		}
	}
}

func TestObjectComparator_MixedChanges(t *testing.T) {
	c := &ObjectComparator{}
	old := parseJSON(`{"a":1,"b":2,"c":3}`)
	newVal := parseJSON(`{"a":1,"b":20,"d":4}`)

	diffs := c.Compare(old, newVal, "$", newCtx())

	opCounts := make(map[Op]int)
	for _, d := range diffs {
		opCounts[d.Op]++
	}
	if opCounts[OpChanged] != 1 {
		t.Errorf("expected 1 changed (b), got %d", opCounts[OpChanged])
	}
	if opCounts[OpRemoved] != 1 {
		t.Errorf("expected 1 removed (c), got %d", opCounts[OpRemoved])
	}
	if opCounts[OpAdded] != 1 {
		t.Errorf("expected 1 added (d), got %d", opCounts[OpAdded])
	}
}

func TestObjectComparator_NestedRecursion(t *testing.T) {
	c := &ObjectComparator{}
	old := parseJSON(`{"outer":{"inner":1}}`)
	newVal := parseJSON(`{"outer":{"inner":2}}`)

	diffs := c.Compare(old, newVal, "$", newCtx())
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(diffs))
	}
	if diffs[0].Path != "$.outer.inner" {
		t.Errorf("expected path $.outer.inner, got %s", diffs[0].Path)
	}
}

func TestObjectComparator_KeysSortedAlphabetically(t *testing.T) {
	c := &ObjectComparator{}
	old := parseJSON(`{"z":1,"a":2,"m":3}`)
	newVal := parseJSON(`{}`)

	diffs := c.Compare(old, newVal, "$", newCtx())
	if len(diffs) != 3 {
		t.Fatalf("expected 3 diffs, got %d", len(diffs))
	}
	expectedPaths := []string{"$.a", "$.m", "$.z"}
	for i, d := range diffs {
		if d.Path != expectedPaths[i] {
			t.Errorf("diff[%d].Path = %s, want %s", i, d.Path, expectedPaths[i])
		}
	}
}

func TestObjectComparator_EmptyObjects(t *testing.T) {
	c := &ObjectComparator{}
	diffs := c.Compare(map[string]interface{}{}, map[string]interface{}{}, "$", newCtx())
	if len(diffs) != 0 {
		t.Errorf("expected no diffs for empty objects, got %d", len(diffs))
	}
}

// --- ArrayComparator tests ---

func TestArrayComparator_BothEmpty(t *testing.T) {
	c := &ArrayComparator{}
	diffs := c.Compare(parseJSON(`[]`), parseJSON(`[]`), "$", newCtx())
	if len(diffs) != 0 {
		t.Errorf("expected no diffs for empty arrays, got %d", len(diffs))
	}
}

func TestArrayComparator_OldEmpty(t *testing.T) {
	c := &ArrayComparator{}
	diffs := c.Compare(parseJSON(`[]`), parseJSON(`[1,2,3]`), "$", newCtx())
	if len(diffs) != 3 {
		t.Fatalf("expected 3 diffs, got %d", len(diffs))
	}
	for _, d := range diffs {
		if d.Op != OpAdded {
			t.Errorf("expected OpAdded, got %v", d.Op)
		}
	}
}

func TestArrayComparator_NewEmpty(t *testing.T) {
	c := &ArrayComparator{}
	diffs := c.Compare(parseJSON(`[1,2,3]`), parseJSON(`[]`), "$", newCtx())
	if len(diffs) != 3 {
		t.Fatalf("expected 3 diffs, got %d", len(diffs))
	}
	for _, d := range diffs {
		if d.Op != OpRemoved {
			t.Errorf("expected OpRemoved, got %v", d.Op)
		}
	}
}

func TestArrayComparator_Identical(t *testing.T) {
	c := &ArrayComparator{}
	diffs := c.Compare(parseJSON(`[1,2,3]`), parseJSON(`[1,2,3]`), "$", newCtx())
	if len(diffs) != 0 {
		t.Errorf("expected no diffs for identical arrays, got %d", len(diffs))
	}
}

func TestArrayComparator_ElementAdded(t *testing.T) {
	c := &ArrayComparator{}
	diffs := c.Compare(parseJSON(`[1,2]`), parseJSON(`[1,2,3]`), "$", newCtx())
	added := 0
	for _, d := range diffs {
		if d.Op == OpAdded {
			added++
		}
	}
	if added != 1 {
		t.Errorf("expected 1 added, got %d", added)
	}
}

func TestArrayComparator_ElementRemoved(t *testing.T) {
	c := &ArrayComparator{}
	diffs := c.Compare(parseJSON(`[1,2,3]`), parseJSON(`[1,3]`), "$", newCtx())
	removed := 0
	for _, d := range diffs {
		if d.Op == OpRemoved {
			removed++
		}
	}
	if removed != 1 {
		t.Errorf("expected 1 removed, got %d", removed)
	}
}

func TestArrayComparator_MoveDetection(t *testing.T) {
	c := &ArrayComparator{}
	old := parseJSON(`["a","b","c"]`)
	newVal := parseJSON(`["c","a","b"]`)

	diffs := c.Compare(old, newVal, "$", newCtx())

	movedCount := 0
	for _, d := range diffs {
		if d.Op == OpMoved {
			movedCount++
		}
	}
	if movedCount == 0 {
		t.Error("expected at least one move for array rearrangement")
	}
}

func TestArrayComparator_MoveDetectionReversed(t *testing.T) {
	c := &ArrayComparator{}
	old := parseJSON(`[1,2,3,4,5]`)
	newVal := parseJSON(`[5,4,3,2,1]`)

	diffs := c.Compare(old, newVal, "$", newCtx())

	moveCount := 0
	for _, d := range diffs {
		if d.Op == OpMoved {
			moveCount++
		}
	}
	if moveCount == 0 {
		t.Error("expected moves for reversed array")
	}
}

func TestArrayComparator_NestedObjectArray(t *testing.T) {
	c := &ArrayComparator{}
	old := parseJSON(`[{"id":"a","val":1},{"id":"b","val":2}]`)
	newVal := parseJSON(`[{"id":"a","val":1},{"id":"b","val":20}]`)

	diffs := c.Compare(old, newVal, "$", newCtx())

	changed := 0
	for _, d := range diffs {
		if d.Op == OpChanged {
			changed++
		}
	}
	if changed != 1 {
		t.Errorf("expected 1 changed (b.val), got %d", changed)
	}
}

func TestArrayComparator_KeyFieldMatching(t *testing.T) {
	c := &ArrayComparator{}
	old := parseJSON(`[{"name":"alice","age":30},{"name":"bob","age":25}]`)
	newVal := parseJSON(`[{"name":"bob","age":25},{"name":"alice","age":31}]`)

	diffs := c.Compare(old, newVal, "$", newCtxWithKey("name"))

	hasChanged := false
	for _, d := range diffs {
		if d.Op == OpChanged {
			hasChanged = true
		}
	}
	if !hasChanged {
		t.Error("expected changed for alice.age")
	}
}

func TestArrayComparator_SingleElement(t *testing.T) {
	c := &ArrayComparator{}
	diffs := c.Compare(parseJSON(`[1]`), parseJSON(`[1]`), "$", newCtx())
	if len(diffs) != 0 {
		t.Errorf("expected no diffs, got %d", len(diffs))
	}
}

func TestArrayComparator_SingleElementChanged(t *testing.T) {
	c := &ArrayComparator{}
	diffs := c.Compare(parseJSON(`[1]`), parseJSON(`[2]`), "$", newCtx())
	if len(diffs) == 0 {
		t.Error("expected diffs for changed single element")
	}
}

func TestArrayComparator_DuplicateElements(t *testing.T) {
	c := &ArrayComparator{}
	diffs := c.Compare(parseJSON(`[1,1,1]`), parseJSON(`[1,1]`), "$", newCtx())
	removed := 0
	for _, d := range diffs {
		if d.Op == OpRemoved {
			removed++
		}
	}
	if removed != 1 {
		t.Errorf("expected 1 removed from [1,1,1]鈫抂1,1], got %d", removed)
	}
}
