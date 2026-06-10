package render

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/fatih/color"
	"github.com/zhanjie/jsondiff/internal/diff"
)

var (
	addedColor   = color.New(color.FgGreen)
	removedColor = color.New(color.FgRed)
	changedColor = color.New(color.FgYellow)
	movedColor   = color.New(color.FgCyan)
	pathColor    = color.New(color.FgHiBlack)
)

// Render outputs the diff items to the given writer with colored terminal output.
func Render(diffs []diff.DiffItem, w io.Writer) {
	if len(diffs) == 0 {
		fmt.Fprintln(w, "No differences found.")
		return
	}

	for _, d := range diffs {
		renderItem(d, w)
	}
}

func renderItem(d diff.DiffItem, w io.Writer) {
	switch d.Op {
	case diff.OpAdded:
		tag := addedColor.Sprint("[ADDED]  ")
		path := pathColor.Sprint(d.Path)
		val := addedColor.Sprint(formatValue(d.NewValue))
		fmt.Fprintf(w, "%s  %s  : %s\n", tag, path, val)

	case diff.OpRemoved:
		tag := removedColor.Sprint("[REMOVED]")
		path := pathColor.Sprint(d.Path)
		val := removedColor.Sprint(formatValue(d.OldValue))
		fmt.Fprintf(w, "%s  %s  : %s\n", tag, path, val)

	case diff.OpChanged:
		tag := changedColor.Sprint("[CHANGED]")
		path := pathColor.Sprint(d.Path)
		oldVal := removedColor.Sprint(formatValue(d.OldValue))
		newVal := addedColor.Sprint(formatValue(d.NewValue))
		fmt.Fprintf(w, "%s  %s  : %s  →  %s\n", tag, path, oldVal, newVal)

	case diff.OpMoved:
		tag := movedColor.Sprint("[MOVED]  ")
		path := pathColor.Sprint(d.Path)
		moveInfo := movedColor.Sprintf("[%d] → [%d]", d.OldIndex, d.NewIndex)
		val := formatValue(d.NewValue)
		fmt.Fprintf(w, "%s  %s  %s  : %s\n", tag, path, moveInfo, val)

	default:
		// OpUnchanged or unknown op — skip silently
	}
}

// formatValue formats a JSON value for display.
func formatValue(v interface{}) string {
	if v == nil {
		return "null"
	}
	switch val := v.(type) {
	case json.Number:
		return val.String()
	case string:
		b, _ := json.Marshal(val)
		return string(b)
	case bool:
		if val {
			return "true"
		}
		return "false"
	case map[string]interface{}, []interface{}:
		b, _ := json.Marshal(val)
		return string(b)
	default:
		return fmt.Sprintf("%v", val)
	}
}
