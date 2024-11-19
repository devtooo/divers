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

		// If this is the last key, set the value
		if i == len(keys)-1 {
			if isIndex {
				// Handle list at the last level
				parentKey := keys[i-1]
				currentList := ensureList(current.(map[string]interface{}), parentKey, index)
				currentList[index] = value
			} else {
				current.(map[string]interface{})[key] = value
			}
			return
		}

		if isIndex {
			// Handle intermediate list
			parentKey := keys[i-1]
			currentList := ensureList(current.(map[string]interface{}), parentKey, index)
			if currentList[index] == nil {
				currentList[index] = make(map[string]interface{})
			}
			current = currentList[index].(map[string]interface{})
		} else {
			// Handle intermediate map
			if _, exists := current.(map[string]interface{})[key]; !exists {
				current.(map[string]interface{})[key] = make(map[string]interface{})
			}
			current = current.(map[string]interface{})[key]
		}
	}
}

// ensureList ensures that a key in the parent map is a list and has at least `size` elements.
func ensureList(parent map[string]interface{}, key string, size int) []interface{} {
	// Check if the key exists and is already a list
	if _, exists := parent[key]; !exists {
		parent[key] = make([]interface{}, size+1)
	}

	// Convert to list if valid
	list, ok := parent[key].([]interface{})
	if !ok {
		log.Fatalf("Expected a list at key '%s' but found a different type", key)
	}

	// Expand the list if necessary
	if len(list) <= size {
		newList := make([]interface{}, size+1)
		copy(newList, list)
		parent[key] = newList
	}

	return parent[key].([]interface{})
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
