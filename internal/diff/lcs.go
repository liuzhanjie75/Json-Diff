package diff

import "sort"

// lcs computes the Longest Common Subsequence of two string slices.
// Returns pairs of [oldIndex, newIndex] that form the LCS,
// sorted by old index in ascending order.
func lcs(a, b []string) [][2]int {
	m, n := len(a), len(b)

	dp := make([][]int, m+1)
	for i := 0; i <= m; i++ {
		dp[i] = make([]int, n+1)
	}

	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if a[i-1] == b[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
			} else if dp[i-1][j] > dp[i][j-1] {
				dp[i][j] = dp[i-1][j]
			} else {
				dp[i][j] = dp[i][j-1]
			}
		}
	}

	var pairs [][2]int
	i, j := m, n
	for i > 0 && j > 0 {
		if a[i-1] == b[j-1] {
			pairs = append(pairs, [2]int{i - 1, j - 1})
			i--
			j--
		} else if dp[i-1][j] > dp[i][j-1] {
			i--
		} else {
			j--
		}
	}

	sort.Slice(pairs, func(a, b int) bool {
		return pairs[a][0] < pairs[b][0]
	})

	return pairs
}
