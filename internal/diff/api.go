package diff

// Options configures the diff behavior
type Options struct {
	// KeyField is the field name used to identify matching objects in arrays.
	// When set, array elements that are objects will be matched by this field's value.
	KeyField string
	// IgnoreArrayOrder compares arrays as multisets and suppresses move diffs.
	// The option applies recursively to nested arrays.
	IgnoreArrayOrder bool
}

// Compare recursively compares two JSON values with default options.
func Compare(old, new interface{}, path string) []DiffItem {
	return CompareWithOpts(old, new, path, Options{})
}

// CompareWithOpts recursively compares two JSON values with the given options.
func CompareWithOpts(old, new interface{}, path string, opts Options) []DiffItem {
	ctx := &Context{
		Opts:       opts,
		Dispatcher: Dispatch,
	}
	return Dispatch(old, new, path, ctx)
}
