package main

import (
	"encoding/json"
	"os"
)

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
