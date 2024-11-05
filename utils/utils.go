package utils

import (
	"encoding/csv"
	"log"
	"os"
)

// ReadCSV reads a CSV file and returns the records.
func ReadCSV(fileName string) [][]string {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	return records
}
