package main

import (
	"fmt"
	"log"
	"wordleSolver/solver"
	"wordleSolver/wordle"
)

// main runs the Wordle solver.
// The program interacts with the user to guess the secret word
// based on feedback provided by the user.
func main() {
	wordleSolver := solver.NewHeuristicSolver()
	attempts := 0

	for {
		attempts++
		guess := wordleSolver.NextGuess()
		fmt.Printf("Guessing word: %s\n", guess)
		fmt.Print("Enter result for guess: ")

		var result string
		_, err := fmt.Scanln(&result)
		if err != nil {
			log.Fatalf("error reading result: %v", err)
		}

		// Check if the guess is correct
		if result == wordle.CorrectResultStr {
			fmt.Printf("Guessed word: %s in %d attempts\n", guess, attempts)
			break
		}

		// Add the result of the guess to the solver
		wordleSolver.AddResult(guess, wordle.ParseResult(result))
	}
}
