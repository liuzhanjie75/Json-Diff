package diff

import (
	"testing"
)

func TestCompare_Integration(t *testing.T) {
	tests := []struct {
		name      string
		old, new  string
		wantOps   map[Op]int // expected count of each op type
		wantCount int        // total expected diff count
	}{
		{
			name:      "identical objects",
			old:       `{"a":1,"b":"hello"}`,
			new:       `{"a":1,"b":"hello"}`,
			wantCount: 0,
		},
		{
			name:      "added key",
			old:       `{"a":1}`,
			new:       `{"a":1,"b":2}`,
			wantOps:   map[Op]int{OpAdded: 1},
			wantCount: 1,
		},
		{
			name:      "removed key",
			old:       `{"a":1,"b":2}`,
			new:       `{"a":1}`,
			wantOps:   map[Op]int{OpRemoved: 1},
			wantCount: 1,
		},
		{
			name:      "changed value",
			old:       `{"a":1}`,
			new:       `{"a":2}`,
			wantOps:   map[Op]int{OpChanged: 1},
			wantCount: 1,
		},
		{
			name:      "type change string to number",
			old:       `{"a":"hello"}`,
			new:       `{"a":123}`,
			wantOps:   map[Op]int{OpChanged: 1},
			wantCount: 1,
		},
		{
			name:      "type change number to bool",
			old:       `{"a":1}`,
			new:       `{"a":true}`,
			wantOps:   map[Op]int{OpChanged: 1},
			wantCount: 1,
		},
		{
			name:      "type change object to array",
			old:       `{"a":{"x":1}}`,
			new:       `{"a":[1,2]}`,
			wantOps:   map[Op]int{OpChanged: 1},
			wantCount: 1,
		},
		{
			name:      "nested object changes",
			old:       `{"config":{"debug":true,"timeout":30}}`,
			new:       `{"config":{"debug":false,"timeout":30,"verbose":true}}`,
			wantOps:   map[Op]int{OpChanged: 1, OpAdded: 1},
			wantCount: 2,
		},
		{
			name:      "null equals null",
			old:       `{"a":null}`,
			new:       `{"a":null}`,
			wantCount: 0,
		},
		{
			name:      "null to value",
			old:       `{"a":null}`,
			new:       `{"a":1}`,
			wantOps:   map[Op]int{OpChanged: 1},
			wantCount: 1,
		},
		{
			name:      "value to null",
			old:       `{"a":1}`,
			new:       `{"a":null}`,
			wantOps:   map[Op]int{OpChanged: 1},
			wantCount: 1,
		},
		{
			name:      "empty objects",
			old:       `{}`,
			new:       `{}`,
			wantCount: 0,
		},
		{
			name:      "empty arrays",
			old:       `[]`,
			new:       `[]`,
			wantCount: 0,
		},
		{
			name:      "array element added",
			old:       `[1,2,3]`,
			new:       `[1,2,3,4]`,
			wantOps:   map[Op]int{OpAdded: 1},
			wantCount: 1,
		},
		{
			name:      "array element removed",
			old:       `[1,2,3]`,
			new:       `[1,2]`,
			wantOps:   map[Op]int{OpRemoved: 1},
			wantCount: 1,
		},
		{
			name:      "array move detection",
			old:       `["a","b","c"]`,
			new:       `["c","a","b"]`,
			wantCount: -1, // at least 1 move, count varies by LCS
		},
		{
			name:      "scalar string equal",
			old:       `"hello"`,
			new:       `"hello"`,
			wantCount: 0,
		},
		{
			name:      "scalar string changed",
			old:       `"hello"`,
			new:       `"world"`,
			wantOps:   map[Op]int{OpChanged: 1},
			wantCount: 1,
		},
		{
			name:      "scalar number equal",
			old:       `42`,
			new:       `42`,
			wantCount: 0,
		},
		{
			name:      "scalar number changed",
			old:       `42`,
			new:       `43`,
			wantOps:   map[Op]int{OpChanged: 1},
			wantCount: 1,
		},
		{
			name:      "scalar bool equal",
			old:       `true`,
			new:       `true`,
			wantCount: 0,
		},
		{
			name:      "scalar bool changed",
			old:       `true`,
			new:       `false`,
			wantOps:   map[Op]int{OpChanged: 1},
			wantCount: 1,
		},
		{
			name:      "nil to object",
			old:       `null`,
			new:       `{"a":1}`,
			wantOps:   map[Op]int{OpChanged: 1},
			wantCount: 1,
		},
		{
			name:      "object to nil",
			old:       `{"a":1}`,
			new:       `null`,
			wantOps:   map[Op]int{OpChanged: 1},
			wantCount: 1,
		},
		{
			name:      "nil to nil",
			old:       `null`,
			new:       `null`,
			wantCount: 0,
		},
		{
			name:      "multiple changes in flat object",
			old:       `{"a":1,"b":2,"c":3}`,
			new:       `{"a":10,"b":2,"d":4}`,
			wantOps:   map[Op]int{OpChanged: 1, OpRemoved: 1, OpAdded: 1},
			wantCount: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldVal := parseJSON(tt.old)
			newVal := parseJSON(tt.new)
			diffs := Compare(oldVal, newVal, "$")

			if tt.wantCount >= 0 && len(diffs) != tt.wantCount {
				t.Errorf("expected %d diffs, got %d: %+v", tt.wantCount, len(diffs), diffs)
			}
			if tt.wantCount == -1 && len(diffs) == 0 {
				t.Error("expected at least one diff, got 0")
			}

			if tt.wantOps != nil {
				gotOps := make(map[Op]int)
				for _, d := range diffs {
					gotOps[d.Op]++
				}
				for op, want := range tt.wantOps {
					if gotOps[op] != want {
						t.Errorf("op %v: expected %d, got %d", op, want, gotOps[op])
					}
				}
			}
		})
	}
}

func TestCompare_WithKeyField(t *testing.T) {
	old := parseJSON(`[
		{"id":"a","val":1},
		{"id":"b","val":2}
	]`)
	newVal := parseJSON(`[
		{"id":"a","val":10},
		{"id":"b","val":2},
		{"id":"c","val":3}
	]`)

	diffs := CompareWithOpts(old, newVal, "$", Options{KeyField: "id"})

	// Expect: a.val changed (1鈫?0), c added
	addedCount := 0
	changedCount := 0
	for _, d := range diffs {
		switch d.Op {
		case OpAdded:
			addedCount++
		case OpChanged:
			changedCount++
		}
	}
	if changedCount != 1 {
		t.Errorf("expected 1 changed (a.val), got %d", changedCount)
	}
	if addedCount != 1 {
		t.Errorf("expected 1 added (c), got %d", addedCount)
	}
}

func TestCompare_WithKeyFieldDoesNotFallbackToSimilarity(t *testing.T) {
	old := parseJSON(`[{"id":"old","name":"item","value":1}]`)
	newVal := parseJSON(`[{"id":"new","name":"item","value":2}]`)

	diffs := CompareWithOpts(old, newVal, "$", Options{KeyField: "id"})

	if len(diffs) != 2 {
		t.Fatalf("expected removed and added diffs, got %d: %+v", len(diffs), diffs)
	}
	opCounts := make(map[Op]int)
	for _, d := range diffs {
		opCounts[d.Op]++
	}
	if opCounts[OpRemoved] != 1 || opCounts[OpAdded] != 1 {
		t.Fatalf("expected one removed and one added diff, got %+v", opCounts)
	}
	if opCounts[OpChanged] != 0 {
		t.Fatalf("expected no field-level changes for different key values, got %+v", diffs)
	}
}

func TestCompare_WithKeyFieldPreservesJSONTypes(t *testing.T) {
	old := parseJSON(`[{"id":"1","value":"old"}]`)
	newVal := parseJSON(`[{"id":1,"value":"new"}]`)

	diffs := CompareWithOpts(old, newVal, "$", Options{KeyField: "id"})

	if len(diffs) != 2 {
		t.Fatalf("expected removed and added diffs, got %d: %+v", len(diffs), diffs)
	}
	opCounts := make(map[Op]int)
	for _, d := range diffs {
		opCounts[d.Op]++
	}
	if opCounts[OpRemoved] != 1 || opCounts[OpAdded] != 1 {
		t.Fatalf("expected one removed and one added diff, got %+v", opCounts)
	}
}

func TestCompare_SimilarityMatching(t *testing.T) {
	// Two objects with same structure but one field different
	old := parseJSON(`[{"name":"x","age":1,"city":"a"}]`)
	newVal := parseJSON(`[{"name":"x","age":2,"city":"a"}]`)

	diffs := Compare(old, newVal, "$")

	// Should detect field-level change, not full remove+add
	hasChanged := false
	for _, d := range diffs {
		if d.Op == OpChanged {
			hasChanged = true
		}
	}
	if !hasChanged {
		t.Error("expected OpChanged for similar objects with one field different")
	}
}

func TestCompare_DeepNesting(t *testing.T) {
	old := parseJSON(`{"a":{"b":{"c":{"d":1}}}}`)
	newVal := parseJSON(`{"a":{"b":{"c":{"d":2}}}}`)

	diffs := Compare(old, newVal, "$")
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(diffs))
	}
	if diffs[0].Path != "$.a.b.c.d" {
		t.Errorf("expected path $.a.b.c.d, got %s", diffs[0].Path)
	}
	if diffs[0].Op != OpChanged {
		t.Errorf("expected OpChanged, got %v", diffs[0].Op)
	}
}

func TestCompare_LargeIntArray(t *testing.T) {
	old := parseJSON(`[1,2,3,4,5]`)
	newVal := parseJSON(`[5,4,3,2,1]`)

	diffs := Compare(old, newVal, "$")

	// 3 is in the same position (LCS match), others should be moves
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

func TestCompare_IgnoreArrayOrder(t *testing.T) {
	tests := []struct {
		name      string
		old       string
		new       string
		wantCount int
		wantOps   map[Op]int
	}{
		{
			name:      "reordered primitives are equal",
			old:       `[1,2]`,
			new:       `[2,1]`,
			wantCount: 0,
		},
		{
			name:      "duplicate counts are significant",
			old:       `[1,1,2]`,
			new:       `[1,2,2]`,
			wantCount: 2,
			wantOps:   map[Op]int{OpRemoved: 1, OpAdded: 1},
		},
		{
			name:      "nested arrays ignore order recursively",
			old:       `[[1,2],[3,4]]`,
			new:       `[[4,3],[2,1]]`,
			wantCount: 0,
		},
		{
			name:      "reordered objects are equal",
			old:       `[{"id":"a","value":1},{"id":"b","value":2}]`,
			new:       `[{"id":"b","value":2},{"id":"a","value":1}]`,
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			old := parseJSON(tt.old)
			newVal := parseJSON(tt.new)
			diffs := CompareWithOpts(old, newVal, "$", Options{IgnoreArrayOrder: true})

			if len(diffs) != tt.wantCount {
				t.Fatalf("expected %d diffs, got %d: %+v", tt.wantCount, len(diffs), diffs)
			}
			if tt.wantOps != nil {
				gotOps := make(map[Op]int)
				for _, d := range diffs {
					gotOps[d.Op]++
				}
				for op, want := range tt.wantOps {
					if gotOps[op] != want {
						t.Fatalf("op %v: expected %d, got %d", op, want, gotOps[op])
					}
				}
			}
			for _, d := range diffs {
				if d.Op == OpMoved {
					t.Fatalf("ignore-array-order must not report moves: %+v", diffs)
				}
			}
		})
	}
}

func TestCompare_IgnoreArrayOrderWithKeyReportsFieldChanges(t *testing.T) {
	old := parseJSON(`[{"id":"a","value":1},{"id":"b","value":2}]`)
	newVal := parseJSON(`[{"id":"b","value":3},{"id":"a","value":1}]`)

	diffs := CompareWithOpts(old, newVal, "$", Options{
		KeyField:         "id",
		IgnoreArrayOrder: true,
	})

	if len(diffs) != 1 {
		t.Fatalf("expected one field change, got %d: %+v", len(diffs), diffs)
	}
	if diffs[0].Op != OpChanged || diffs[0].Path != "$[0].value" {
		t.Fatalf("unexpected diff: %+v", diffs[0])
	}
}

func TestOp_String(t *testing.T) {
	tests := []struct {
		op   Op
		want string
	}{
		{OpUnchanged, "UNCHANGED"},
		{OpAdded, "ADDED"},
		{OpRemoved, "REMOVED"},
		{OpChanged, "CHANGED"},
		{OpMoved, "MOVED"},
	}
	for _, tt := range tests {
		if got := tt.op.String(); got != tt.want {
			t.Errorf("Op(%d).String() = %q, want %q", tt.op, got, tt.want)
		}
	}
}
