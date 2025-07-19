package kowalski

import (
	"strings"
	"unicode"
)

// Chunk takes the input, and splits it up into chunks of the given length. If the input is longer than the list of
// part lengths, the lengths will be repeated.
func Chunk(input string, parts ...int) []string {
	var res []string
	remaining := input
	p := 0
	for len(remaining) > 0 {
		if len(remaining) >= parts[p] {
			res = append(res, remaining[0:parts[p]])
			remaining = remaining[parts[p]:]
		} else {
			res = append(res, remaining)
			remaining = ""
		}

		p = (p + 1) % len(parts)
	}
	return res
}

// Transpose rotates the input text so rows become columns, and columns become rows.
func Transpose(input []string) []string {
	var res []string
	i := 0
	found := true
	for found {
		found = false
		line := strings.Builder{}
		for j := range input {
			if len(input[j]) > i {
				line.WriteByte(input[j][i])
				found = true
			}
		}
		res = append(res, line.String())
		i++
	}
	return res
}

// FirstLetters extracts the first letter of each word in the input, preserving line breaks.
// Punctuation is ignored - only letters are extracted.
func FirstLetters(input string) string {
	lines := strings.Split(input, "\n")
	var result []string

	for _, line := range lines {
		if line == "" {
			result = append(result, "")
			continue
		}

		words := strings.Fields(line)
		var firstLetters []string
		for _, word := range words {
			// Find the first letter (ignoring punctuation)
			for _, r := range word {
				if unicode.IsLetter(r) {
					firstLetters = append(firstLetters, string(r))
					break
				}
			}
		}
		result = append(result, strings.Join(firstLetters, ""))
	}

	return strings.Join(result, "\n")
}

// Reverse reverses the input text, stripping punctuation but preserving spaces and line breaks.
func Reverse(input string) string {
	// Extract only letters, digits, spaces, and newlines
	var filtered []rune
	for _, r := range input {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == ' ' || r == '\n' {
			filtered = append(filtered, r)
		}
	}

	// Reverse the filtered runes
	for i, j := 0, len(filtered)-1; i < j; i, j = i+1, j-1 {
		filtered[i], filtered[j] = filtered[j], filtered[i]
	}

	return string(filtered)
}

// WordCheckResult represents a word and whether it's valid according to the spell checkers.
type WordCheckResult struct {
	Word  string
	Valid bool
}

// CheckWords splits the input into words and checks each one against the spell checker.
// Returns a slice of results for each line.
func CheckWords(checker *SpellChecker, input string) [][]WordCheckResult {
	lines := strings.Split(input, "\n")
	results := make([][]WordCheckResult, len(lines))

	for i, line := range lines {
		words := strings.Fields(line)
		lineResults := make([]WordCheckResult, 0, len(words))

		for _, word := range words {
			// Extract just the letters from the word
			cleanWord := strings.Builder{}
			for _, r := range word {
				if unicode.IsLetter(r) {
					cleanWord.WriteRune(unicode.ToLower(r))
				}
			}

			clean := cleanWord.String()
			if clean != "" {
				lineResults = append(lineResults, WordCheckResult{
					Word:  word,
					Valid: checker.Valid(clean),
				})
			}
		}
		results[i] = lineResults
	}

	return results
}
