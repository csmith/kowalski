package kowalski

import (
	"strings"
)

// FindWords attempts to find substrings of the input that are valid words according to the checker.
// Duplicates may be present in the output if they occur at multiple positions.
func FindWords(checker *SpellChecker, input string) []string {
	var res []string

	findWords(checker, input, func(start, end int) {
		res = append(res, input[start:end])
	})

	return res
}

// findWords finds all substrings of the given input, calling func with their start and end offsets.
func findWords(checker *SpellChecker, input string, fn func(start, end int)) {
	lower := strings.ToLower(input)
	for i := 0; i < len(input); i++ {
		for j := i + 1; j < len(input)+1 && checker.Prefix(lower[i:j]); j++ {
			if checker.Valid(lower[i:j]) {
				fn(i, j)
			}
		}
	}
}

// WordSearch returns all words found by FindWords in the input word search grid. Words may occur horizontally,
// vertically or diagonally, and may read in either direction. If a word is found multiple times in different
// places it will be returned multiple times.
func WordSearch(checker *SpellChecker, input []string) []string {
	var res []string
	lines := wordSearchLines(input)
	for i := range lines {
		words := FindWords(checker, lines[i])
		for j := range words {
			if len(words[j]) >= 4 {
				res = append(res, words[j])
			}
		}
	}
	return res
}

func wordSearchLines(input []string) []string {
	var res []string

	if len(input) == 0 {
		return res
	}

	width := len(input[0])
	longestSide := max(len(input), width)

	// Horizontal
	for i := range input {
		res = append(res, input[i])
		res = append(res, reverseString(input[i]))
	}

	// Vertical
	for i := 0; i < width; i++ {
		chars := make([]byte, len(input))
		for j := range input {
			chars[j] = input[j][i]
		}
		res = append(res, string(chars))
		res = append(res, string(reverse(chars, 0)))
	}

	// Diagonals
	for i := 0; i < len(input)+longestSide; i++ {
		pos := strings.Builder{}
		neg := strings.Builder{}
		for j := 0; j < width; j++ {
			if k := i + j - longestSide; k >= 0 && k < len(input) {
				neg.WriteByte(input[k][j])
			}
			if k := i - j; k >= 0 && k < len(input) {
				pos.WriteByte(input[k][j])
			}
		}

		if pos.Len() > 1 {
			res = append(res, pos.String())
			res = append(res, reverseString(pos.String()))
		}

		if neg.Len() > 1 {
			res = append(res, neg.String())
			res = append(res, reverseString(neg.String()))
		}
	}

	return res
}
