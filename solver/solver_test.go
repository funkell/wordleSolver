package solver

import (
	"log"
	"math/rand"
	"sort"
	"testing"
	"wordleSolver/wordle"
)

// calculateMean calculates the mean number of attempts.
func calculateMean(attemptsCount map[int]int) float64 {
	totalAttempts := 0
	totalWords := 0
	for attempts, count := range attemptsCount {
		totalAttempts += attempts * count
		totalWords += count
	}
	return float64(totalAttempts) / float64(totalWords)
}

// calculateMedian calculates the median number of attempts.
func calculateMedian(attemptsCount map[int]int) float64 {
	totalWords := 0
	for _, count := range attemptsCount {
		totalWords += count
	}

	mid := totalWords / 2
	countSoFar := 0
	var median int

	// Extract keys and sort them
	var attempts []int
	for attempt := range attemptsCount {
		attempts = append(attempts, attempt)
	}
	sort.Ints(attempts)

	// Iterate over sorted keys
	for _, attempt := range attempts {
		countSoFar += attemptsCount[attempt]
		if countSoFar >= mid {
			median = attempt
			break
		}
	}

	return float64(median)
}

// TestWordleSolverGuesses tests the Wordle solver with some
// predefined secret words.
func TestWordleSolverGuesses(t *testing.T) {
	wordleSolver := NewHeuristicSolver()

	tests := []struct {
		name   string
		secret string
	}{
		{"SecretWordApple", "apple"},
		{"SecretWordPeach", "peach"},
		{"SecretWordGrape", "grape"},
		{"SecretWordMango", "mango"},
		{"SecretWordLemon", "lemon"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wordleSolver.Reset()
			attempts := 0
			for {
				attempts++
				guess := wordleSolver.NextGuess()
				// If the guess is correct, break the loop
				if guess == tt.secret {
					break
				}
				result := wordle.ComputeResult(guess, tt.secret)
				wordleSolver.AddResult(guess, result)
			}
			log.Printf("Secret word: %s, Guessed in %d attempts\n", tt.secret, attempts)
			if attempts > 6 {
				t.Errorf("expected to guess the word %s in less than 6 attempts, but took %d", tt.secret, attempts)
			}
		})
	}
}

// TestWordleSolverAllWords tests the Wordle solver by running it against a sample of words
// and calculates the mean and median number of attempts required to guess the words.
func TestWordleSolverAllWords(t *testing.T) {
	wordleSolver := NewHeuristicSolver()
	words := getInitialWordsList()

	const sampleSize = 100
	// Shuffle the words
	rand.Shuffle(len(words), func(i, j int) {
		words[i], words[j] = words[j], words[i]
	})
	// Select the first sampleSize words
	sampleWords := words[:sampleSize]

	attemptsCount := make(map[int]int)

	for _, word := range sampleWords {
		wordleSolver.Reset()
		attempts := 0
		for {
			attempts++
			guess := wordleSolver.NextGuess()
			if guess == word.word {
				break
			}
			result := wordle.ComputeResult(guess, word.word)
			wordleSolver.AddResult(guess, result)
		}
		attemptsCount[attempts]++
	}

	// Calculate mean and median
	mean := calculateMean(attemptsCount)
	median := calculateMedian(attemptsCount)

	log.Println("Statistics of attempts:")
	var sortedAttempts []int
	for attempts := range attemptsCount {
		sortedAttempts = append(sortedAttempts, attempts)
	}
	sort.Ints(sortedAttempts)

	for _, attempts := range sortedAttempts {
		log.Printf("Number of words guessed in %d attempts: %d\n", attempts, attemptsCount[attempts])
	}
	log.Printf("Mean number of attempts: %.2f\n", mean)
	log.Printf("Median number of attempts: %.2f\n", median)
}
