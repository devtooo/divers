package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func main() {
	// Open the CSV file
	f, err := os.Open("data.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Create a new CSV reader
	reader := csv.NewReader(f)

	// Read the header
	header, err := reader.Read()
	if err != nil {
		panic(err)
	}

	// Read all records
	var records [][]string
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		records = append(records, record)
	}

	// Process each record and build JSON objects
	var output []map[string]interface{}
	for _, record := range records {
		obj := make(map[string]interface{})
		for i, value := range record {
			pathStr := header[i]
			path := strings.Split(pathStr, "/")
			// Remove empty strings from path parts
			var pathParts []string
			for _, part := range path {
				if part != "" {
					pathParts = append(pathParts, part)
				}
			}
			setValue(obj, pathParts, value)
		}
		output = append(output, obj)
	}

	// Convert the output to JSON
	jsonBytes, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(jsonBytes))
}

// setValue sets a value in a nested map based on the provided path
func setValue(obj map[string]interface{}, path []string, value interface{}) {
	if len(path) == 0 {
		return
	}
	key := path[0]
	if len(path) == 1 {
		obj[key] = value
		return
	}
	nextKey := path[1]
	if idx, err := strconv.Atoi(nextKey); err == nil {
		// Handle array index
		var arr []interface{}
		if existing, ok := obj[key]; ok {
			arr = existing.([]interface{})
		} else {
			arr = make([]interface{}, idx+1)
			obj[key] = arr
		}
		// Ensure array is large enough
		if len(arr) <= idx {
			newArr := make([]interface{}, idx+1)
			copy(newArr, arr)
			arr = newArr
			obj[key] = arr
		}
		if len(path) == 2 {
			arr[idx] = value
			return
		}
		var nextObj map[string]interface{}
		if arr[idx] != nil {
			nextObj = arr[idx].(map[string]interface{})
		} else {
			nextObj = make(map[string]interface{})
			arr[idx] = nextObj
		}
		setValue(nextObj, path[2:], value)
	} else {
		// Handle nested map
		var nextObj map[string]interface{}
		if existing, ok := obj[key]; ok {
			nextObj = existing.(map[string]interface{})
		} else {
			nextObj = make(map[string]interface{})
			obj[key] = nextObj
		}
		setValue(nextObj, path[1:], value)
	}
}
