package main

import "testing"

func TestLevenshteinDistance(t *testing.T) {
	values := [][2]string{
		{"", ""},
		{"ab", "ba"},
		{"azcde", "abcde"},
		{"abcd", "abcd"},
		{"юяяя", "яяяф"},
		{"abcd", "яяяф"},
	}
	result := []int{0, 2, 1, 0, 2, 4}
	for i, strs := range values {
		if j := LevenshteinDistance(strs[0], strs[1]); j != result[i] {
			t.Errorf("error for string #%v", i)
		}
	}
}

func TestPrepare(t *testing.T) {
	// sorted values
	values := []Word{
		{L: 2, W: "ab"},
		{L: 2, W: "ba"},
		{L: 2, W: "be"},
		{L: 3, W: "abc"},
		{L: 3, W: "abe"},
		{L: 3, W: "afe"},
		{L: 4, W: "wxyy"},
		{L: 4, W: "wyxy"},
		{L: 4, W: "wyyy"},
	}
	expected := []Leaf{
		{Root: "ab", Relations: []int{}},
		{Root: "ba", Relations: []int{2}},
		{Root: "be", Relations: []int{1}},
		{Root: "abc", Relations: []int{4}},
		{Root: "abe", Relations: []int{3, 5}},
		{Root: "afe", Relations: []int{4}},
		{Root: "wxyy", Relations: []int{8}},
		{Root: "wyxy", Relations: []int{8}},
		{Root: "wyyy", Relations: []int{6, 7}},
	}
	prepared := Prepare(values)

	if i, j := len(prepared), len(expected); i != j {
		t.Fatalf("lengths are not equal: %v != %v", i, j)
	}
	for i := range expected {
		if !expected[i].Equal(&prepared[i]) {
			t.Errorf("it is not equal [%v]: %v", i, prepared[i])
		}
	}
}

func TestSearch(t *testing.T) {
	leafs := []Leaf{
		{Root: "ab", Relations: []int{}},
		{Root: "ba", Relations: []int{2}},
		{Root: "be", Relations: []int{1}},
		{Root: "abc", Relations: []int{4}},
		{Root: "abe", Relations: []int{3, 5}},
		{Root: "afe", Relations: []int{4}},
		{Root: "wxyy", Relations: []int{8}},
		{Root: "wyxy", Relations: []int{8}},
		{Root: "wyyy", Relations: []int{6, 7}},
	}
	failed := [][2]string{
		{"ab", "ba"},
		{"ab", "kw"},
		{"", "kw"},
		{"ba", "kw"},
		{"wk", "kw"},
	}
	for i, v := range failed {
		if _, err := Search(leafs, v[0], v[1]); err == nil {
			t.Errorf("unexpected [%v]: %v\n", i, err)
		}
	}
	success := [][2]string{
		{"ba", "be"},
		{"abc", "afe"},
		{"afe", "abc"},
		{"wxyy", "wyyy"},
		{"wxyy", "wyxy"},
		{"wyxy", "wxyy"},
	}
	expected := [][]string{
		{"ba", "be"},
		{"abc", "abe", "afe"},
		{"afe", "abe", "abc"},
		{"wxyy", "wyyy"},
		{"wxyy", "wyyy", "wyxy"},
		{"wyxy", "wyyy", "wxyy"},
	}
	for i, v := range success {
		s, err := Search(leafs, v[0], v[1])
		if err != nil {
			t.Errorf("failed [%v]: %v\n", i, err)
			continue
		}
		if a, b := len(s), len(expected[i]); a != b {
			t.Errorf("failed lengths: %v != %v\n", a, b)
			continue
		}
		for j, p := range expected[i] {
			if p != s[j] {
				t.Errorf("[%v] '%v' != '%v'\n", i, p, s[j])
			}
		}
	}
}
