// Package main provides a command-line tool for solving Wordle puzzles.
// Game link: https://www.nytimes.com/games/wordle/index.html
// The program interacts with the user to guess the secret word based on feedback provided by the user.
//
// Usage:
// 1. Run the program.
// 2. The program will make an initial guess and prompt the user to enter the result.
// 3. The result should be a string of 5 characters where:
//   - 'C' indicates the character is correct and in the correct position.
//   - 'P' indicates the character is present in the word but in the wrong position.
//   - 'N' indicates the character is not present in the word.
//
// 4. The program will continue to make guesses and prompt for results until the word is guessed correctly.
//
// Example run:
// Guessing word: soare
// Enter result for guess: PNNPP
// Guessing word: reist
// Enter result for guess: PPNCN
// Guessing word: press
// Enter result for guess: NCCCN
// Guessing word: fresh
// Enter result for guess: CCCCC
// Guessed word: fresh in 4 attempts
package main
