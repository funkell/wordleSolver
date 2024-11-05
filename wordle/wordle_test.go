package wordle

import (
	"testing"
)

// TestParseResult tests the ParseResult function with various inputs.
func TestParseResult(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expected  GuessResult
		wantPanic bool
	}{
		{"ValidInput", "CCPNN", GuessResult{Correct, Correct, Present, NotPresent, NotPresent}, false},
		{"InvalidLength", "CCP", GuessResult{}, true},
		{"InvalidCharacter", "CCPXN", GuessResult{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil && !tt.wantPanic {
					t.Errorf("unexpected panic for input %s", tt.input)
				} else if r == nil && tt.wantPanic {
					t.Errorf("expected panic for input %s", tt.input)
				}
			}()
			result := ParseResult(tt.input)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

// TestComputeResult tests the ComputeResult function with various inputs.
func TestComputeResult(t *testing.T) {
	tests := []struct {
		name     string
		guess    string
		secret   string
		expected GuessResult
	}{
		{"AllCorrect", "apple", "apple", GuessResult{Correct, Correct, Correct, Correct, Correct}},
		{"Mixed", "apple", "ample", GuessResult{Correct, NotPresent, Correct, Correct, Correct}},
		{"AllNotPresent", "apple", "zzzzz", GuessResult{NotPresent, NotPresent, NotPresent, NotPresent, NotPresent}},
		{"SomePresent", "apple", "peach", GuessResult{Present, Present, NotPresent, NotPresent, Present}},
		{"SomeCorrectSomePresent", "apple", "apric", GuessResult{Correct, Correct, NotPresent, NotPresent, NotPresent}},
		{"RepeatedLetters", "apple", "allee", GuessResult{Correct, NotPresent, NotPresent, Present, Correct}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ComputeResult(tt.guess, tt.secret)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
