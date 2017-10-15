package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"time"
)

var (
	logger = log.New(os.Stderr, fmt.Sprintf("METAGRAMMA "), log.Ldate|log.Ltime|log.Lshortfile)
)

func main() {
	var leafs []Leaf

	start := time.Now()
	defer func() {
		fmt.Printf("duration %v\n", time.Now().Sub(start))
	}()

	init := flag.String("i", "", "configuration file")
	flag.Parse()

	if *init == "" {
		logger.Println("empty init file name")
		return
	}
	data, err := readFile(*init)
	if err != nil {
		logger.Fatal(err)
	}
	ch := make(chan []Leaf)

	parts := 0
	prev, currentLen := 0, 0
	for i, w := range data {
		if w.L != currentLen {
			parts++
			//fmt.Printf("xaz, parts=%v, cur=%v, prev=%v, i=%v, data=%v\n", parts, currentLen, prev, i, data[prev:i])
			go Build(data[prev:i], prev, ch)
			prev, currentLen = i, w.L
		}
	}
	batchLeafs := make([][]Leaf, parts)
	for j := 0; j < parts; j++ {
		batchLeafs[j] = <-ch
	}
	close(ch)
	sort.Sort(Leafs(batchLeafs))

	for j := range batchLeafs {
		leafs = append(leafs, batchLeafs[j]...)
	}
	fmt.Println("done")
	//for i, leaf := range leafs {
	//	fmt.Printf("%v: %v - %v\n", i, leaf.Root, leaf.Relations)
	//}
}
