package wordle

import "log"

// Constants for position results
const (
	NotPresent PositionResult = iota
	Present
	Correct
)

// PositionResult represents the result of a single character in a guess.
type PositionResult int8

// GuessResult represents the result of a guess.
type GuessResult [5]PositionResult

// ComputeResult calculates the result of a guess against the secret word.
func ComputeResult(guess, secret string) GuessResult {
	result := GuessResult{}
	secretCharCount := make(map[byte]int)

	// First pass: identify correct positions
	for i := range guess {
		if guess[i] == secret[i] {
			result[i] = Correct
		} else {
			secretCharCount[secret[i]]++
		}
	}

	// Second pass: identify present and not present characters
	for i := range guess {
		// Skip already identified correct positions
		if result[i] == Correct {
			continue
		}
		// Check if the character is present in the secret word
		// and not already accounted for
		if secretCharCount[guess[i]] > 0 {
			result[i] = Present
			secretCharCount[guess[i]]--
		} else {
			result[i] = NotPresent
		}
	}

	return result
}

var CorrectResultStr = "CCCCC"

// ParseResult parses the result string into a GuessResult.
func ParseResult(result string) GuessResult {
	if len(result) != 5 {
		log.Panicf("invalid result: %s", result)
	}
	parsedResult := GuessResult{}
	for i, r := range result {
		switch r {
		case 'C':
			parsedResult[i] = Correct
		case 'P':
			parsedResult[i] = Present
		case 'N':
			parsedResult[i] = NotPresent
		default:
			log.Panicf("invalid character in result: %c", r)
		}
	}
	return parsedResult
}
