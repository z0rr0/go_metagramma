package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"unicode/utf8"
)

var (
	loggerError = log.New(os.Stderr, fmt.Sprintf("ERROR "), log.Ldate|log.Ltime|log.Lshortfile)
)

type Word struct {
	L int
	W string
}

type Words []Word

func (a Words) Len() int           { return len(a) }
func (a Words) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Words) Less(i, j int) bool { return (a[i].L < a[j].L) && (a[i].W < a[j].W) }

type Leaf struct {
	Root      string
	Relations []int
}

func min(values ...int) int {
	m := values[0]
	for _, v := range values {
		if v < m {
			m = v
		}
	}
	return m
}

func levenshtein_distance(a, b string) int {
	n, m := len(a), len(b)
	if n > m {
		a, b = b, a
		n, m = m, n
	}
	currentRow := make([]int, n+1)
	previousRow := make([]int, n+1)
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

func createTree(lines []Word) []Leaf {
	var (
		result []Leaf
		candidates []int
	)

	for i, line := range lines {
		// candidates
		candidates = []int{}
		b := sort.Search(len(lines[:i]), func(j int) bool { return lines[j].L >= line.L })
		for t, v := range lines[b:i] {
			//fmt.Println(v.W, line.W)
			if (v.W != line.W) && (levenshtein_distance(v.W, line.W) == 1) {
				candidates = append(candidates, b + t)
			}
		}
		//fmt.Println(line, b, candidates	)
		for j := range candidates {
			result[j].Relations = append(result[j].Relations, i)
		}
		result = append(result, Leaf{Root: line.W, Relations: candidates})
	}
	return result
}

func main() {
	init := flag.String("i", "", "configuration file")
	flag.Parse()

	if *init == "" {
		loggerError.Println("empty init file name")
		return
	}
	data, err := readFile(*init)
	if err != nil {
		loggerError.Fatal(err)
	}
	fmt.Println(data)
	fmt.Println(createTree(data))
	//fmt.Println(levenshtein_distance("amc", "amcerf"))
	//fmt.Println(levenshtein_distance("", ""))
	//fmt.Println(levenshtein_distance("abc", "abe"))
	//fmt.Println(levenshtein_distance("abcdef", "abzzef"))
}
