package diff

import "sort"

// ObjectComparator handles comparison of JSON objects (maps).
// It compares keys from both objects using set operations:
// - Keys only in old → REMOVED
// - Keys only in new → ADDED
// - Keys in both → recursively compared via Dispatcher
type ObjectComparator struct{}

func (c *ObjectComparator) Compare(old, new interface{}, path string, ctx *Context) []DiffItem {
	oldMap := old.(map[string]interface{})
	newMap := new.(map[string]interface{})

	var diffs []DiffItem

	// Collect all keys from both maps
	keys := make(map[string]bool)
	for k := range oldMap {
		keys[k] = true
	}
	for k := range newMap {
		keys[k] = true
	}

	// Sort keys for stable output
	sortedKeys := make([]string, 0, len(keys))
	for k := range keys {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)

	for _, key := range sortedKeys {
		childPath := path + "." + key
		oldVal, oldHas := oldMap[key]
		newVal, newHas := newMap[key]

		if oldHas && !newHas {
			diffs = append(diffs, DiffItem{
				Op:       OpRemoved,
				Path:     childPath,
				OldValue: oldVal,
			})
		} else if !oldHas && newHas {
			diffs = append(diffs, DiffItem{
				Op:       OpAdded,
				Path:     childPath,
				NewValue: newVal,
			})
		} else {
			// Both have the key, recurse via dispatcher
			diffs = append(diffs, ctx.Dispatcher(oldVal, newVal, childPath, ctx)...)
		}
	}

	return diffs
}
