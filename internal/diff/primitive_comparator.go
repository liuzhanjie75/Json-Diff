package diff

import (
	"encoding/json"
	"reflect"
)

// PrimitiveComparator handles comparison of JSON primitive types:
// string, number (json.Number), and boolean.
// If the types differ or values are not equal, it reports a CHANGED diff.
type PrimitiveComparator struct{}

func (c *PrimitiveComparator) Compare(old, new interface{}, path string, ctx *Context) []DiffItem {
	// Type mismatch (e.g. string vs number, object vs primitive)
	if reflect.TypeOf(old) != reflect.TypeOf(new) {
		return []DiffItem{{Op: OpChanged, Path: path, OldValue: old, NewValue: new}}
	}

	// Same type, compare values
	if !equalScalars(old, new) {
		return []DiffItem{{Op: OpChanged, Path: path, OldValue: old, NewValue: new}}
	}

	return nil
}

// equalScalars compares two scalar JSON values for equality.
// json.Number values are compared by their string representation
// to avoid floating-point precision issues.
func equalScalars(a, b interface{}) bool {
	aNum, aIsNum := a.(json.Number)
	bNum, bIsNum := b.(json.Number)
	if aIsNum && bIsNum {
		return aNum.String() == bNum.String()
	}
	return reflect.DeepEqual(a, b)
}
