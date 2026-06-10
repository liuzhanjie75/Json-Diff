package render_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/fatih/color"
	"github.com/zhanjie/jsondiff/internal/diff"
	"github.com/zhanjie/jsondiff/internal/render"
)

func init() {
	// Disable colors for predictable test output
	color.NoColor = true
}

func TestRender_NoDiffs(t *testing.T) {
	var buf bytes.Buffer
	render.Render(nil, &buf)
	if !strings.Contains(buf.String(), "No differences") {
		t.Errorf("expected 'No differences' message, got %q", buf.String())
	}
}

func TestRender_EmptyDiffs(t *testing.T) {
	var buf bytes.Buffer
	render.Render([]diff.DiffItem{}, &buf)
	if !strings.Contains(buf.String(), "No differences") {
		t.Errorf("expected 'No differences' message, got %q", buf.String())
	}
}

func TestRender_Added(t *testing.T) {
	var buf bytes.Buffer
	diffs := []diff.DiffItem{
		{Op: diff.OpAdded, Path: "$.name", NewValue: "alice"},
	}
	render.Render(diffs, &buf)
	output := buf.String()
	if !strings.Contains(output, "[ADDED]") {
		t.Errorf("expected [ADDED] tag, got %q", output)
	}
	if !strings.Contains(output, "$.name") {
		t.Errorf("expected path $.name, got %q", output)
	}
	if !strings.Contains(output, `"alice"`) {
		t.Errorf("expected value alice, got %q", output)
	}
}

func TestRender_Removed(t *testing.T) {
	var buf bytes.Buffer
	diffs := []diff.DiffItem{
		{Op: diff.OpRemoved, Path: "$.age", OldValue: 30},
	}
	render.Render(diffs, &buf)
	output := buf.String()
	if !strings.Contains(output, "[REMOVED]") {
		t.Errorf("expected [REMOVED] tag, got %q", output)
	}
	if !strings.Contains(output, "$.age") {
		t.Errorf("expected path $.age, got %q", output)
	}
}

func TestRender_Changed(t *testing.T) {
	var buf bytes.Buffer
	diffs := []diff.DiffItem{
		{Op: diff.OpChanged, Path: "$.val", OldValue: "old", NewValue: "new"},
	}
	render.Render(diffs, &buf)
	output := buf.String()
	if !strings.Contains(output, "[CHANGED]") {
		t.Errorf("expected [CHANGED] tag, got %q", output)
	}
	if !strings.Contains(output, "→") {
		t.Errorf("expected arrow (→), got %q", output)
	}
}

func TestRender_Moved(t *testing.T) {
	var buf bytes.Buffer
	diffs := []diff.DiffItem{
		{Op: diff.OpMoved, Path: "$[0]", OldValue: "x", NewValue: "x", OldIndex: 2, NewIndex: 0},
	}
	render.Render(diffs, &buf)
	output := buf.String()
	if !strings.Contains(output, "[MOVED]") {
		t.Errorf("expected [MOVED] tag, got %q", output)
	}
	if !strings.Contains(output, "[2] → [0]") {
		t.Errorf("expected move info, got %q", output)
	}
}

func TestRender_MultipleDiffs(t *testing.T) {
	var buf bytes.Buffer
	diffs := []diff.DiffItem{
		{Op: diff.OpAdded, Path: "$.a", NewValue: 1},
		{Op: diff.OpRemoved, Path: "$.b", OldValue: 2},
		{Op: diff.OpChanged, Path: "$.c", OldValue: 3, NewValue: 4},
	}
	render.Render(diffs, &buf)
	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) != 3 {
		t.Errorf("expected 3 lines, got %d", len(lines))
	}
}

func TestRender_NullValue(t *testing.T) {
	var buf bytes.Buffer
	diffs := []diff.DiffItem{
		{Op: diff.OpAdded, Path: "$.x", NewValue: nil},
	}
	render.Render(diffs, &buf)
	if !strings.Contains(buf.String(), "null") {
		t.Errorf("expected 'null' in output, got %q", buf.String())
	}
}
