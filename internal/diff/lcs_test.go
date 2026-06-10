package diff

import (
	"testing"
)

func TestLCS_EmptyInputs(t *testing.T) {
	pairs := lcs(nil, nil)
	if len(pairs) != 0 {
		t.Errorf("expected empty LCS for nil inputs, got %d", len(pairs))
	}
}

func TestLCS_OneEmpty(t *testing.T) {
	pairs := lcs([]string{"a", "b"}, nil)
	if len(pairs) != 0 {
		t.Errorf("expected empty LCS when one input is nil, got %d", len(pairs))
	}
}

func TestLCS_Identical(t *testing.T) {
	pairs := lcs([]string{"a", "b", "c"}, []string{"a", "b", "c"})
	if len(pairs) != 3 {
		t.Fatalf("expected 3 pairs, got %d", len(pairs))
	}
	for i, p := range pairs {
		if p[0] != i || p[1] != i {
			t.Errorf("pair[%d] = [%d,%d], want [%d,%d]", i, p[0], p[1], i, i)
		}
	}
}

func TestLCS_CompletelyDifferent(t *testing.T) {
	pairs := lcs([]string{"a", "b"}, []string{"c", "d"})
	if len(pairs) != 0 {
		t.Errorf("expected empty LCS for completely different inputs, got %d", len(pairs))
	}
}

func TestLCS_Subsequence(t *testing.T) {
	// LCS of [a,b,c,d] and [a,c,d] = [a,c,d] 鈫?3 pairs
	pairs := lcs([]string{"a", "b", "c", "d"}, []string{"a", "c", "d"})
	if len(pairs) != 3 {
		t.Fatalf("expected 3 pairs, got %d", len(pairs))
	}
	expected := [][2]int{{0, 0}, {2, 1}, {3, 2}}
	for i, p := range pairs {
		if p != expected[i] {
			t.Errorf("pair[%d] = %v, want %v", i, p, expected[i])
		}
	}
}

func TestLCS_Reversed(t *testing.T) {
	pairs := lcs([]string{"a", "b", "c"}, []string{"c", "b", "a"})
	if len(pairs) != 1 {
		t.Errorf("expected 1 pair for reversed input, got %d", len(pairs))
	}
}

func TestLCS_PartialOverlap(t *testing.T) {
	pairs := lcs([]string{"a", "b", "c", "d"}, []string{"b", "d", "e"})
	if len(pairs) != 2 {
		t.Fatalf("expected 2 pairs, got %d", len(pairs))
	}
}

func TestLCS_SortedByOldIndex(t *testing.T) {
	pairs := lcs([]string{"x", "a", "b"}, []string{"b", "x", "a"})
	for i := 1; i < len(pairs); i++ {
		if pairs[i][0] < pairs[i-1][0] {
			t.Errorf("pairs not sorted by old index: %v", pairs)
			break
		}
	}
}

func TestLCS_SingleElement(t *testing.T) {
	pairs := lcs([]string{"a"}, []string{"a"})
	if len(pairs) != 1 || pairs[0] != [2]int{0, 0} {
		t.Errorf("expected [(0,0)], got %v", pairs)
	}
}

func TestLCS_DuplicateElements(t *testing.T) {
	pairs := lcs([]string{"a", "a", "a"}, []string{"a", "a"})
	if len(pairs) != 2 {
		t.Errorf("expected 2 pairs for [a,a,a] vs [a,a], got %d", len(pairs))
	}
}
