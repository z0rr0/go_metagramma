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
	//for i, v := range leafs {
	//	fmt.Printf("%v: %v %v\n", i, v.Root, v.Relations)
	//}
	err = SaveJSON(leafs, outFile)
	if err != nil {
		return err
	}
	return nil
}

func isSearch(dbFile string) error {
	leafs, err := ReadJSON(dbFile)
	if err != nil {
		return err
	}
	fmt.Printf("read %v items\n", len(leafs))
	return nil
}

func main() {
	var err error
	start := time.Now()
	defer func() {
		fmt.Printf("duration %v\n", time.Now().Sub(start))
	}()

	init := flag.String("i", "", "configuration file")
	output := flag.String("o", "", "output file")
	db := flag.String("d", "", "prepared JSON file")
	flag.Parse()

	if *init != "" {
		if *output == "" {
			logger.Fatal("prapred db error")
		}
		err = isInit(*init, *output)
	} else {
		if *db == "" {
			logger.Fatal("no prepared JSON file set")
		}
		err = isSearch(*db)
	}
	if err != nil {
		logger.Fatalln(err)
	}
}
