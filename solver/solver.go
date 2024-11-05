package solver

import (
	"log"
	"math"
	"math/big"
	"os"
	"strconv"
	"wordleSolver/utils"
	"wordleSolver/wordle"
)

const (
	wordsFileName = "../data/5-letter-words.csv"
	cacheFileName = "first_guess_cache.txt"
)

// WordleSolver defines the interface for a Wordle solver.
type WordleSolver interface {
	// NextGuess returns the next guess based on the previous guesses and feedback.
	NextGuess() string

	// AddResult adds the result of a guess to the solver.
	AddResult(guess string, result wordle.GuessResult)
}

// ResettableWordleSolver extends WordleSolver with a Reset method.
type ResettableWordleSolver interface {
	WordleSolver
	// Reset resets the solver to its initial state.
	Reset()
}

// word represents a word and its weight.
type word struct {
	word   string
	weight int64
}

// heuristicSolver implements the WordleSolver interface using a
// heuristic approach, that minimizes the expected number of remaining
// words after each guess.
type heuristicSolver struct {
	allWords      []word
	possibleWords []word
}

// calculateScoreForGuess calculates the score for a guess.
// The score is proportional to the expected number of
// remaining words after the guess.
func (h *heuristicSolver) calculateScoreForGuess(guess string) *big.Int {
	partitionWeights := make(map[wordle.GuessResult]*big.Int)
	partitionCounts := make(map[wordle.GuessResult]int64)

	// Determine the partition for each remaining word and
	// update the weights and counts for each partition.
	for _, possibleWord := range h.possibleWords {
		if possibleWord.word == guess {
			continue
		}
		wordResult := wordle.ComputeResult(guess, possibleWord.word)
		weight, ok := partitionWeights[wordResult]
		if !ok {
			weight = big.NewInt(0)
			partitionWeights[wordResult] = weight
		}
		weight.Add(weight, big.NewInt(possibleWord.weight))
		partitionCounts[wordResult]++
	}

	totalExpectedPartitionSize := big.NewInt(0)
	for result, partitionWeight := range partitionWeights {
		expectedPartitionSize := big.NewInt(partitionCounts[result])
		expectedPartitionSize.Mul(expectedPartitionSize, partitionWeight)
		totalExpectedPartitionSize.Add(totalExpectedPartitionSize, expectedPartitionSize)
	}

	return totalExpectedPartitionSize
}

// readCachedGuessIfPresent reads the cached guess from the file if it exists.
func readCachedGuessIfPresent() (string, bool) {
	_, err := os.Stat(cacheFileName)
	if os.IsNotExist(err) {
		return "", false
	}

	file, err := os.Open(cacheFileName)
	if err != nil {
		log.Panicf("error opening cache file: %v", err)
	}
	defer file.Close()

	content, err := os.ReadFile(cacheFileName)
	if err != nil {
		log.Printf("error reading from cache file: %v", err)
		return "", false
	}
	return string(content), true
}

// writeCachedGuess writes the guess to the cache file.
func writeCachedGuess(guess string) {
	file, err := os.Create(cacheFileName)
	if err != nil {
		log.Fatalf("error creating cache file: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(guess)
	if err != nil {
		log.Fatalf("error writing to cache file: %v", err)
	}
}

// NextGuess returns the next guess based on the previous guesses and feedback.
func (h *heuristicSolver) NextGuess() string {
	if len(h.possibleWords) == 0 {
		log.Fatalf("no possible words left")
	}
	if len(h.possibleWords) == 1 {
		return h.possibleWords[0].word
	}

	// If  no guesses have been made yet, check if there is a
	// cached first guess. First guess is cached since computing
	// it is expensive, and it is the same for all secret words.
	if len(h.possibleWords) == len(h.allWords) {
		cachedGuess, ok := readCachedGuessIfPresent()
		if ok {
			return cachedGuess
		}
	}

	// Find the guess that minimizes the expected number
	// of remaining words.
	var guess string
	minScore := big.NewInt(math.MaxInt64)
	minScore.Mul(minScore, minScore)
	for _, word := range h.allWords {
		score := h.calculateScoreForGuess(word.word)
		if score.Cmp(minScore) < 0 {
			minScore = score
			guess = word.word
		}
	}

	// If this is the first guess, cache it.
	if len(h.possibleWords) == len(h.allWords) {
		writeCachedGuess(guess)
	}
	return guess
}

// AddResult adds the result of a guess to the solver.
func (h *heuristicSolver) AddResult(guess string, result wordle.GuessResult) {
	// Filter out the words that are not consistent with the result.
	var newPossibleWords []word
	for _, word := range h.possibleWords {
		if wordle.ComputeResult(guess, word.word) == result {
			newPossibleWords = append(newPossibleWords, word)
		}
	}
	h.possibleWords = newPossibleWords
}

// Reset resets the solver to its initial state.
func (h *heuristicSolver) Reset() {
	h.possibleWords = h.allWords
}

// NewHeuristicSolver creates a new HeuristicSolver.
func NewHeuristicSolver() ResettableWordleSolver {
	words := getInitialWordsList()
	return &heuristicSolver{
		allWords:      words,
		possibleWords: words,
	}
}

// getInitialWordsList reads the initial words and their frequencies from a CSV file.
func getInitialWordsList() []word {
	records := utils.ReadCSV(wordsFileName)
	words := make([]word, 0)
	for _, record := range records {
		frequency, err := strconv.ParseInt(record[1], 10, 64)
		if err != nil {
			log.Fatalf("error parsing frequency for record: %s, err: %v", record, err)
		}
		words = append(words, word{record[0], frequency})
	}
	log.Printf("Read %d words from file\n", len(words))
	return words
}
