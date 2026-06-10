package diff

// NullComparator handles comparison involving null (nil) values.
// Cases:
// - Both nil → no difference
// - Old is nil, new is not → CHANGED (value was null, now has a value)
// - Old is not nil, new is nil → CHANGED (value existed, now is null)
type NullComparator struct{}

func (c *NullComparator) Compare(old, new interface{}, path string, ctx *Context) []DiffItem {
	if old == nil && new == nil {
		return nil
	}
	if old == nil {
		return []DiffItem{{Op: OpChanged, Path: path, OldValue: nil, NewValue: new}}
	}
	// new == nil
	return []DiffItem{{Op: OpChanged, Path: path, OldValue: old, NewValue: nil}}
}
