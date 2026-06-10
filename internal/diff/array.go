package diff

import "fmt"

// similarityThreshold is the minimum Jaccard similarity ratio to consider
// two objects as the "same element with modifications".
const similarityThreshold = 0.5

// ArrayComparator handles comparison of JSON arrays.
// It orchestrates a three-phase matching strategy:
//  1. Hash-based LCS (exact match)
//  2. Key-based or similarity-based matching for unmatched object elements
//  3. Remaining unmatched elements → REMOVED / ADDED / MOVED
type ArrayComparator struct{}

func (c *ArrayComparator) Compare(old, new interface{}, path string, ctx *Context) []DiffItem {
	oldArr := old.([]interface{})
	newArr := new.([]interface{})
	return compareArrays(oldArr, newArr, path, ctx)
}

func compareArrays(old, new []interface{}, path string, ctx *Context) []DiffItem {
	m, n := len(old), len(new)

	if m == 0 && n == 0 {
		return nil
	}
	if m == 0 {
		return allAdded(new, path)
	}
	if n == 0 {
		return allRemoved(old, path)
	}

	// Build hash arrays for exact matching
	oldHashes := make([]string, m)
	newHashes := make([]string, n)
	for i := 0; i < m; i++ {
		oldHashes[i] = HashJSON(old[i])
	}
	for j := 0; j < n; j++ {
		newHashes[j] = HashJSON(new[j])
	}

	// Phase 1: LCS on exact hash match
	lcsPairs := lcs(oldHashes, newHashes)
	matchedOld := make(map[int]bool)
	matchedNew := make(map[int]bool)

	var diffs []DiffItem
	for _, p := range lcsPairs {
		matchedOld[p[0]] = true
		matchedNew[p[1]] = true
		childPath := fmt.Sprintf("%s[%d]", path, p[1])
		diffs = append(diffs, ctx.Dispatcher(old[p[0]], new[p[1]], childPath, ctx)...)
	}

	// Collect unmatched indices
	unmatchedOld := collectUnmatched(m, matchedOld)
	unmatchedNew := collectUnmatched(n, matchedNew)

	// Phase 2: Match by key or similarity
	matched2Old := make(map[int]bool)
	matched2New := make(map[int]bool)

	if ctx.Opts.KeyField != "" {
		matchByKey(unmatchedOld, unmatchedNew, old, new, ctx.Opts.KeyField,
			path, matched2Old, matched2New, &diffs, ctx)
	} else {
		matchBySimilarity(unmatchedOld, unmatchedNew, old, new,
			path, matched2Old, matched2New, &diffs, ctx)
	}

	// Phase 3: Move detection for remaining hash-matched elements
	finalOld := filterUnmatched(unmatchedOld, matched2Old)
	finalNew := filterUnmatched(unmatchedNew, matched2New)
	detectMoves(finalOld, finalNew, oldHashes, newHashes, old, new,
		path, matched2Old, matched2New, &diffs)

	// Remaining → Removed / Added
	for _, oi := range unmatchedOld {
		if !matched2Old[oi] {
			diffs = append(diffs, DiffItem{
				Op: OpRemoved, Path: fmt.Sprintf("%s[%d]", path, oi), OldValue: old[oi],
			})
		}
	}
	for _, nj := range unmatchedNew {
		if !matched2New[nj] {
			diffs = append(diffs, DiffItem{
				Op: OpAdded, Path: fmt.Sprintf("%s[%d]", path, nj), NewValue: new[nj],
			})
		}
	}

	return diffs
}

// detectMoves finds elements that moved position (same hash, different index).
func detectMoves(finalOld, finalNew []int, oldHashes, newHashes []string,
	old, new []interface{}, path string,
	matchedOld, matchedNew map[int]bool, diffs *[]DiffItem) {

	for _, oi := range finalOld {
		if matchedOld[oi] {
			continue
		}
		for _, nj := range finalNew {
			if matchedNew[nj] {
				continue
			}
			if oldHashes[oi] == newHashes[nj] {
				*diffs = append(*diffs, DiffItem{
					Op: OpMoved, Path: fmt.Sprintf("%s[%d]", path, nj),
					OldValue: old[oi], NewValue: new[nj],
					OldIndex: oi, NewIndex: nj,
				})
				matchedOld[oi] = true
				matchedNew[nj] = true
				break
			}
		}
	}
}

func collectUnmatched(size int, matched map[int]bool) []int {
	var result []int
	for i := 0; i < size; i++ {
		if !matched[i] {
			result = append(result, i)
		}
	}
	return result
}

func allAdded(arr []interface{}, path string) []DiffItem {
	diffs := make([]DiffItem, len(arr))
	for j := 0; j < len(arr); j++ {
		diffs[j] = DiffItem{Op: OpAdded, Path: fmt.Sprintf("%s[%d]", path, j), NewValue: arr[j]}
	}
	return diffs
}

func allRemoved(arr []interface{}, path string) []DiffItem {
	diffs := make([]DiffItem, len(arr))
	for i := 0; i < len(arr); i++ {
		diffs[i] = DiffItem{Op: OpRemoved, Path: fmt.Sprintf("%s[%d]", path, i), OldValue: arr[i]}
	}
	return diffs
}
