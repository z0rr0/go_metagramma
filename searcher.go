package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
)

// Vertex is struct for graph's vertex.
type Vertex struct {
	Num  int
	Size int
	Prev *Vertex
}

// ReadJSON reads prepared result from the file.
func ReadJSON(filename string) ([]Leaf, error) {
	var leafs []Leaf
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&leafs)
	if err != nil {
		return nil, err
	}
	return leafs, nil
}

func searchWord(leafs []Leaf, w string) int {
	i := sort.Search(len(leafs), func(i int) bool { return leafs[i].Greater(w) })
	if (i < len(leafs)) && (leafs[i].Root == w) {
		return i
	}
	return -1
}

// Search searches the shortest path between two words.
// It is based on Dijkstra's algorithm.
func Search(leafs []Leaf, start, end string) ([]string, error) {
	var minPath int
	s, e := searchWord(leafs, start), searchWord(leafs, end)
	if s < 0 {
		return nil, fmt.Errorf("word '%v' is not found in the data file", start)
	}
	if e < 0 {
		return nil, fmt.Errorf("word '%v' is not found in the data file", end)
	}
	current := &Vertex{Num: s, Size: 0, Prev: nil}
	black := map[int]*Vertex{}
	grey := map[int]*Vertex{current.Num: current}

	for {
		black[current.Num] = grey[current.Num]
		delete(grey, current.Num)

		vertexes := leafs[current.Num].Relations
		sort.Sort(sort.IntSlice(vertexes))
		for _, v := range vertexes {
			if _, ok := black[v]; ok {
				continue
			}
			if _, ok := grey[v]; ok {
				// skip grey vertex
				continue
			}
			grey[v] = &Vertex{Num: v, Size: current.Size + 1, Prev: current}
			if v == e {
				current = grey[v]
				grey = map[int]*Vertex{} // "for" exit condition
				break
			}
		}
		if len(grey) == 0 {
			break
		}
		minPath = -1
		for s, v := range grey {
			if (minPath < 0) || (v.Size < current.Size) {
				current, minPath = grey[s], v.Size
			}
		}
	}
	if leafs[current.Num].Root != end {
		// not found
		return nil, errors.New("not found way")
	}

	result := []string{}
	for current != nil {
		result = append(result, leafs[current.Num].Root)
		current = current.Prev
	}
	// reverse result
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}
	return result, nil
}
