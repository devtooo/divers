package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Struct to represent the JSON output format
type Data struct {
	Key   int      `json:"key"`
	Items []string `json:"items"`
	Cc    []CcItem `json:"cc"`
}

type CcItem struct {
	Ff string `json:"ff"`
	Qq string `json:"qq"`
}

func main() {
	// Open the CSV file
	file, err := os.Open("input.csv")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Read the CSV file
	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV:", err)
		return
	}

	// Map to hold the transformed data by key
	dataMap := make(map[string]*Data)

	// Iterate over CSV records, skipping the header row
	for _, record := range records[1:] {
		key := record[0]
		item := record[1]
		c1 := strings.TrimSpace(record[2])
		c2 := strings.TrimSpace(record[3])
		//_c3 := strings.TrimSpace(record[4])

		// If the key does not exist in the map, initialize a new Data struct
		if _, exists := dataMap[key]; !exists {
			dataMap[key] = &Data{
				Key:   atoi(key),
				Items: []string{},
				Cc:    []CcItem{},
			}
		}

		// Add the item to the Items list
		dataMap[key].Items = append(dataMap[key].Items, item)

		// If columns c1 and c2 are present, add them to the cc list
		if c1 != "" && c2 != "" {
			dataMap[key].Cc = append(dataMap[key].Cc, CcItem{
				Ff: c1,
				Qq: c2,
			})
		}
	}

	// Convert the map to JSON and print the results
	for _, data := range dataMap {
		jsonData, err := json.MarshalIndent(data, "", "    ")
		if err != nil {
			fmt.Println("Error marshalling to JSON:", err)
			return
		}
		fmt.Println(string(jsonData))
	}
}

// Helper function to convert string to int
func atoi(s string) int {
	result, _ := strconv.Atoi(s)
	return result
}
