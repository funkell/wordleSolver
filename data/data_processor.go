package main

import (
	"bufio"
	"encoding/csv"
	"log"
	"os"
	"sort"
	"wordleSolver/utils"
)

const (
	// Ref: https://www.kaggle.com/datasets/wheelercode/english-word-frequency-list
	// ngram_freq.csv is a dataset for commonly used words in English.
	inputFileName = "data/ngram_freq.csv"
	// 5-letter-words.csv is the output file that will contain the 5-letter words,
	// and their frequencies, that are allowed to be guessed in Wordle.
	outputFileName = "data/5-letter-words.csv"
	// allowed_guesses.txt is a file that contains the list of allowed 5-letter words,
	// that can be guessed in Wordle.
	allowedWordsFileName = "data/allowed_guesses.txt"
)

// main performs the data processing to filter the
// allowed 5-letter words from the input file.
func main() {
	records := utils.ReadCSV(inputFileName)
	allowedWords := readAllowedWords(allowedWordsFileName)
	fiveLetterRecords := filterRecords(records, allowedWords)
	writeCSV(outputFileName, fiveLetterRecords)
	log.Println("5-letter-words.csv has been created successfully")
}

// readAllowedWords reads allowed words from a file and returns them as a map.
func readAllowedWords(fileName string) map[string]bool {
	allowedWords := make(map[string]bool)
	file, err := os.Open(fileName)
	checkError(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		allowedWords[scanner.Text()] = true
	}
	checkError(scanner.Err())
	return allowedWords
}

// filterRecords filters the records to include only those that are in the allowed words map.
func filterRecords(records [][]string, allowedWords map[string]bool) [][]string {
	var filteredRecords [][]string
	for _, record := range records {
		if allowedWords[record[0]] {
			filteredRecords = append(filteredRecords, record)
		}
	}
	sort.Slice(filteredRecords, func(i, j int) bool {
		return filteredRecords[i][0] < filteredRecords[j][0]
	})
	return filteredRecords
}

// writeCSV writes the filtered records to a CSV file.
func writeCSV(fileName string, records [][]string) {
	file, err := os.Create(fileName)
	checkError(err)
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, record := range records {
		checkError(writer.Write(record))
	}
}

// checkError logs and exits if an error is encountered.
func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
