package diff

import (
	"encoding/json"
	"fmt"
)

// matchByKey matches unmatched array elements using a specified key field value.
func matchByKey(unmatchedOld, unmatchedNew []int, old, new []interface{},
	keyField, path string, matchedOld, matchedNew map[int]bool,
	diffs *[]DiffItem, ctx *Context) {

	newByKey := make(map[string][]int)
	for _, nj := range unmatchedNew {
		njMap, ok := new[nj].(map[string]interface{})
		if !ok {
			continue
		}
		kv := keyStringValue(njMap[keyField])
		if kv != "" {
			newByKey[kv] = append(newByKey[kv], nj)
		}
	}

	for _, oi := range unmatchedOld {
		oiMap, ok := old[oi].(map[string]interface{})
		if !ok {
			continue
		}
		kv := keyStringValue(oiMap[keyField])
		if kv == "" {
			continue
		}
		for _, nj := range newByKey[kv] {
			if matchedNew[nj] {
				continue
			}
			childPath := fmt.Sprintf("%s[%d]", path, nj)
			nested := ctx.Dispatcher(old[oi], new[nj], childPath, ctx)
			*diffs = append(*diffs, nested...)
			matchedOld[oi] = true
			matchedNew[nj] = true
			break
		}
	}
}

// matchBySimilarity matches unmatched object elements based on key structure similarity.
func matchBySimilarity(unmatchedOld, unmatchedNew []int, old, new []interface{},
	path string, matchedOld, matchedNew map[int]bool,
	diffs *[]DiffItem, ctx *Context) {

	for _, oi := range unmatchedOld {
		if matchedOld[oi] {
			continue
		}
		oiMap, ok := old[oi].(map[string]interface{})
		if !ok {
			continue
		}

		bestJ := -1
		bestSim := 0.0

		for _, nj := range unmatchedNew {
			if matchedNew[nj] {
				continue
			}
			njMap, ok := new[nj].(map[string]interface{})
			if !ok {
				continue
			}
			sim := objectSimilarity(oiMap, njMap)
			if sim > bestSim && sim >= similarityThreshold {
				bestSim = sim
				bestJ = nj
			}
		}

		if bestJ >= 0 {
			childPath := fmt.Sprintf("%s[%d]", path, bestJ)
			nested := ctx.Dispatcher(old[oi], new[bestJ], childPath, ctx)
			*diffs = append(*diffs, nested...)
			matchedOld[oi] = true
			matchedNew[bestJ] = true
		}
	}
}

// objectSimilarity computes the Jaccard similarity of two objects' key sets.
func objectSimilarity(a, b map[string]interface{}) float64 {
	if len(a) == 0 && len(b) == 0 {
		return 1.0
	}

	allKeys := make(map[string]bool)
	for k := range a {
		allKeys[k] = true
	}
	for k := range b {
		allKeys[k] = true
	}

	shared := 0
	for k := range allKeys {
		if _, inA := a[k]; inA {
			if _, inB := b[k]; inB {
				shared++
			}
		}
	}

	return float64(shared) / float64(len(allKeys))
}

// keyStringValue converts a JSON value to a comparable string for key matching.
func keyStringValue(v interface{}) string {
	if v == nil {
		return ""
	}
	b, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(b)
}

// filterUnmatched returns indices that are not yet in the matched set.
func filterUnmatched(indices []int, matched map[int]bool) []int {
	var result []int
	for _, i := range indices {
		if !matched[i] {
			result = append(result, i)
		}
	}
	return result
}
