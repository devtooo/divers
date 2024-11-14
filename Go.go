package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

// SetNestedValue sets a value in a nested map based on a JSON path
func SetNestedValue(jsonObject map[string]interface{}, path string, value string) {
	// Remove the '$.' prefix from the JSON path
	path = strings.TrimPrefix(path, "$.")
	keys := strings.Split(path, ".")

	// Traverse or create the map structure based on keys
	current := jsonObject
	for i, key := range keys {
		// If we're at the last key, set the value
		if i == len(keys)-1 {
			current[key] = value
			return
		}

		// If the key does not exist or is not a map, create a new map
		if _, exists := current[key]; !exists {
			current[key] = make(map[string]interface{})
		}

		// Move to the next level in the map
		if next, ok := current[key].(map[string]interface{}); ok {
			current = next
		} else {
			log.Fatalf("Failed to traverse JSON path: %s", path)
		}
	}
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
