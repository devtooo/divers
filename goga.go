package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// SetNestedValue sets a value in a nested map or list based on a JSON path with slashes
func SetNestedValue(jsonObject map[string]interface{}, path string, value string) {
	// Remove the leading "/" if present
	if strings.HasPrefix(path, "/") {
		path = strings.TrimPrefix(path, "/")
	}

	// Split the path by "/"
	keys := strings.Split(path, "/")

	// Traverse or create the map/list structure based on keys
	current := jsonObject
	for i, key := range keys {
		// Check if the key represents an index (for a list)
		index, err := strconv.Atoi(key)
		isIndex := err == nil

		if i == len(keys)-1 {
			// If we are at the last key, set the value
			if isIndex {
				// Handle list index if it exists
				currentList := ensureList(current, keys[i-1], index)
				currentList[index] = value
			} else {
				current[key] = value
			}
			return
		}

		if isIndex {
			// Handle nested list
			prevKey := keys[i-1]
			currentList := ensureList(current, prevKey, index)
			if currentList[index] == nil {
				currentList[index] = make(map[string]interface{})
			}
			current = currentList[index].(map[string]interface{})
		} else {
			// Handle nested map
			if _, exists := current[key]; !exists {
				current[key] = make(map[string]interface{})
			}
			current = current[key].(map[string]interface{})
		}
	}
}

// ensureList ensures that a key in the parent map is a list and has at least `size` elements.
func ensureList(parent map[string]interface{}, key string, size int) []interface{} {
	// Check if the key exists and is a list
	if _, exists := parent[key]; !exists {
		parent[key] = make([]interface{}, size+1)
	}
	if list, ok := parent[key].([]interface{}); ok {
		// Expand the list if necessary
		if len(list) <= size {
			newList := make([]interface{}, size+1)
			copy(newList, list)
			parent[key] = newList
		}
		return parent[key].([]interface{})
	}
	log.Fatalf("Expected a list at key %s", key)
	return nil
}

func main() {
	// Open the CSV file
	file, err := os.Open("data.csv")
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	// Read the CSV file
	reader := csv.NewReader(file)
	headers, err := reader.Read()
	if err != nil {
		log.Fatalf("Failed to read headers: %v", err)
	}

	var jsonArray []map[string]interface{}

	// Process each row
	for {
		record, err := reader.Read()
		if err != nil {
			break
		}

		// Create a new JSON object for each row
		jsonObject := make(map[string]interface{})

		for i, header := range headers {
			SetNestedValue(jsonObject, header, record[i])
		}

		// Append JSON object to array
		jsonArray = append(jsonArray, jsonObject)
	}

	// Convert array of JSON objects to a single JSON
	finalJSON, err := json.MarshalIndent(jsonArray, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal final JSON: %v", err)
	}

	// Output the final JSON
	fmt.Println(string(finalJSON))
}
