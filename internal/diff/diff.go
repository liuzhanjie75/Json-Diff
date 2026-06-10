package diff

// Comparator is the interface for all JSON type comparators.
// Each JSON type (Object, Array, Primitive, Null) has its own implementation.
type Comparator interface {
	// Compare compares two JSON values and returns a list of differences.
	// path is the current JSON path (e.g. "$", "$.users.0.name").
	Compare(old, new interface{}, path string, ctx *Context) []DiffItem
}

// Context holds shared configuration and the dispatcher for recursive comparisons.
type Context struct {
	Opts       Options
	Dispatcher func(old, new interface{}, path string, ctx *Context) []DiffItem
}

// --- Registry ---

var (
	objectComparator    = &ObjectComparator{}
	arrayComparator     = &ArrayComparator{}
	primitiveComparator = &PrimitiveComparator{}
	nullComparator      = &NullComparator{}
)

// Dispatch selects the appropriate comparator based on the JSON types
// and delegates the comparison to it.
func Dispatch(old, new interface{}, path string, ctx *Context) []DiffItem {
	// NullComparator handles nil cases
	if old == nil || new == nil {
		return nullComparator.Compare(old, new, path, ctx)
	}

	_, oldIsMap := old.(map[string]interface{})
	_, newIsMap := new.(map[string]interface{})
	if oldIsMap && newIsMap {
		return objectComparator.Compare(old, new, path, ctx)
	}

	_, oldIsArr := old.([]interface{})
	_, newIsArr := new.([]interface{})
	if oldIsArr && newIsArr {
		return arrayComparator.Compare(old, new, path, ctx)
	}

	// Everything else is a primitive comparison
	return primitiveComparator.Compare(old, new, path, ctx)
}
