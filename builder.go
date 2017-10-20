package main

import (
	"bufio"
	"encoding/json"
	"os"
	"sort"
	"strings"
	"time"
	"unicode/utf8"
)

// Word is a struct for one word data.
type Word struct {
	L int
	W string
}

// Words is a slice of Word structures.
type Words []Word

// Len returns a length or elements in the Words
func (a Words) Len() int { return len(a) }

// Swap swaps elements in the Words.
func (a Words) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

// Less compares two elements in the Words.
func (a Words) Less(i, j int) bool {
	if a[i].L == a[j].L {
		return a[i].W < a[j].W
	}
	return a[i].L < a[j].L
}

// Leaf is result items of prepared graph.
type Leaf struct {
	Root      string `json:"root"`
	Relations []int  `json:"relations"`
}

// Leafs is a slice of Leaf's slices.
type Leafs [][]Leaf

func (l *Leaf) Equal(x *Leaf) bool {
	if (l.Root != x.Root) || (len(l.Relations) != len(x.Relations)) {
		return false
	}
	for i := range l.Relations {
		if l.Relations[i] != x.Relations[i] {
			return false
		}
	}
	return true
}

func leafSize(l []Leaf) int {
	if len(l) == 0 {
		return 0
	}
	return len(l[0].Root)
}

// Len returns a length or elements in the Leafs slice.
func (l Leafs) Len() int { return len(l) }

// Swap swaps elements in the Leafs slice.
func (l Leafs) Swap(i, j int) { l[i], l[j] = l[j], l[i] }

// Less compares two elements in the Leafs slice.
func (l Leafs) Less(i, j int) bool {
	return leafSize(l[i]) < leafSize(l[j])

}

// min returns minimal int element from values.
func min(values ...int) int {
	m := values[0]
	for _, v := range values {
		if v < m {
			m = v
		}
	}
	return m
}

// LevenshteinDistance calculates Levenshtein distance for two strings.
func LevenshteinDistance(a, b string) int {
	n, m := len(a), len(b)
	if n > m {
		a, b = b, a
		n, m = m, n
	}
	currentRow := make([]int, n+1)
	previousRow := make([]int, n+1)
	for i := range currentRow {
		currentRow[i] = i
	}
	for i := 1; i <= m; i++ {
		for j := range currentRow {
			previousRow[j] = currentRow[j]
			if j == 0 {
				currentRow[j] = i
				continue
			} else {
				currentRow[j] = 0
			}
			add, del, change := previousRow[j]+1, currentRow[j-1]+1, previousRow[j-1]
			if a[j-1] != b[i-1] {
				change += 1
			}
			currentRow[j] = min(add, del, change)
		}
	}
	return currentRow[n]
}

// customLD is a custom Levenshtein distance method.
// It wraps LevenshteinDistance function for special cases to don't calculate LD matrix.
func (a *Word) customLD(b *Word) bool {
	if (a.L != b.L) || (a.W == b.W) {
		return false
	}
	return LevenshteinDistance(a.W, b.W) == 1
}

// readFile reads initial file and prepares sorted Word's slice.
func readFile(name string) ([]Word, error) {
	var (
		s      string
		result []Word
	)
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s = strings.Trim(scanner.Text(), " ")
		result = append(result, Word{L: utf8.RuneCountInString(s), W: s})
	}
	sort.Sort(Words(result))
	return result, nil
}

// buildAll is common (for single goroutine) builder of the graph.
// It isn't used now.
func buildAll(lines []Word) []Leaf {
	var (
		result     []Leaf
		candidates []int
	)

	for i, line := range lines {
		// candidates
		candidates = []int{}
		b := sort.Search(len(lines[:i]), func(j int) bool { return lines[j].L >= line.L })
		for t, v := range lines[b:i] {
			if line.customLD(&v) {
				candidates = append(candidates, b+t)
			}
		}
		for _, j := range candidates {
			result[j].Relations = append(result[j].Relations, i)
		}
		result = append(result, Leaf{Root: line.W, Relations: candidates})
	}
	return result
}

// Build returns a result graph for incoming lines.
func Build(lines []Word, offset int, ch chan []Leaf) {
	var n int
	start := time.Now()
	if len(lines) != 0 {
		n = lines[0].L
	}
	logger.Printf("start Build, len=%v, offset=%v\n", n, offset)

	var candidates []int
	result := make([]Leaf, len(lines))

	for i, line := range lines {
		candidates = []int{}
		for t, v := range lines[:i] {
			if line.customLD(&v) {
				candidates = append(candidates, t+offset)
			}
		}
		for _, j := range candidates {
			idx := j - offset
			result[idx].Relations = append(result[idx].Relations, i+offset)
		}
		result[i] = Leaf{Root: line.W, Relations: candidates}
	}
	logger.Printf(
		"end Build, len=%v, offset=%v, items=%v, duration=%v\n",
		n, offset, len(result), time.Now().Sub(start),
	)
	ch <- result
}

// Prepare builds result Word's array.
func Prepare(data []Word) []Leaf {
	var leafs []Leaf

	ch := make(chan []Leaf)
	parts := 0
	prev, currentLen := 0, 0

	for i, w := range data {
		if w.L != currentLen {
			parts++
			go Build(data[prev:i], prev, ch)
			prev, currentLen = i, w.L
		}
	}
	parts++
	go Build(data[prev:], prev, ch)

	batchLeafs := make([][]Leaf, parts)
	for j := 0; j < parts; j++ {
		batchLeafs[j] = <-ch
	}
	close(ch)
	sort.Sort(Leafs(batchLeafs))

	for j := range batchLeafs {
		leafs = append(leafs, batchLeafs[j]...)
	}
	return leafs
}

// SaveJSON saves prepared result to the file.
func SaveJSON(leafs []Leaf, filename string) error {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(&leafs)
}
