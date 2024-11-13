package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/yalp/jsonpath"
)

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
			if err := jsonpath.Set(jsonObject, header, record[i]); err != nil {
				log.Printf("Failed to set JSON path %s: %v", header, err)
			}
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
