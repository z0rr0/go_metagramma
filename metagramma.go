package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

var (
	logger = log.New(os.Stderr, fmt.Sprintf("METAGRAMMA "), log.Ldate|log.Ltime|log.Lshortfile)
)

func isInit(dictFile, outFile string) error {
	data, err := readFile(dictFile)
	if err != nil {
		return err
	}
	leafs := Prepare(data)

	fmt.Printf("prepared %v items\n", len(leafs))
	err = SaveJSON(leafs, outFile)
	if err != nil {
		return err
	}
	return nil
}

func isSearch(dbFile, start, end string) ([]string, error) {
	leafs, err := ReadJSON(dbFile)
	if err != nil {
		return nil, err
	}
	logger.Printf("%v items are read from %v\n", len(leafs), dbFile)
	wordsChain, err := Search(leafs, start, end)
	if err != nil {
		return nil, err
	}
	return wordsChain, nil
}

func main() {
	var (
		err        error
		wordsChain []string
	)
	start := time.Now()
	defer func() {
		fmt.Printf("duration %v\n", time.Now().Sub(start))
	}()

	init := flag.String("i", "", "init 	dict file")
	output := flag.String("o", "", "output file")
	db := flag.String("d", "", "prepared JSON file")

	fromWord := flag.String("f", "", "from word (start)")
	toWord := flag.String("t", "", "to word (end)")
	flag.Parse()

	if *init != "" {
		if *output == "" {
			logger.Fatalln("prapred db error")
		}
		err = isInit(*init, *output)
	} else {
		switch {
		case *db == "":
			logger.Fatalln("no prepared JSON file set")
		case (*fromWord == "") || (*toWord == ""):
			logger.Fatalln("start or end word is not set")
		case len(*fromWord) != len(*toWord):
			logger.Fatalln("lengths of words are not equal")
		case *fromWord == *toWord:
			logger.Fatalln("words are equal")
		}
		wordsChain, err = isSearch(*db, *fromWord, *toWord)
		if err == nil {
			if len(wordsChain) < 2 {
				fmt.Println("not found way")
			} else {
				for i, w := range wordsChain {
					fmt.Printf("%v: %v\n", i, w)
				}
			}
		}
	}
	if err != nil {
		logger.Fatalln(err)
	}
}
