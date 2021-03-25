package kowalski

import (
	"fmt"
	"sort"
	"strings"
)

// Match returns all valid words that match the given pattern, expanding '?' as a single character wildcard
func Match(checker *SpellChecker, pattern string) []string {
	res, _ := findMatch(checker, strings.ToLower(pattern))
	return res
}

// OffByOne returns all words that can be made by performing one character change on the input. The input is
// assumed to be a single, lowercase word containing a-z chars only.
func OffByOne(checker *SpellChecker, input string) []string {
	words := map[string]bool{}
	for i := range input {
		res := Match(checker, fmt.Sprintf("%s?%s", input[0:i], input[i+1:]))
		for j := range res {
			words[res[j]] = true
		}
	}

	var res []string
	for w := range words {
		if w != input {
			res = append(res, w)
		}
	}
	return res
}

// findMatch returns all valid words that match the given pattern, expanding '?' as a single character wildcard.
// It will aggressively skip sequences that don't form valid prefixes; the maximum valid prefix length is returned as
// the second parameter (for cases where matches are returned, this will equal len(word)).
func findMatch(checker *SpellChecker, word string) ([]string, int) {
	maxLength := 0
	stems := []string{""}
	for offset := 0; offset < len(word) && len(stems) > 0; offset++ {
		newStems := make([]string, 0, len(stems))

		var chars []uint8
		if word[offset] == '?' {
			chars = []uint8("abcdefghijklmnopqrstuvwxyz")
		} else {
			chars = []uint8{word[offset]}
		}

		for _, nextChar := range chars {
			for s := range stems {
				if word := fmt.Sprintf("%s%c", stems[s], nextChar); checker.Prefix(word) {
					newStems = append(newStems, word)
				}
			}
		}

		maxLength = offset
		stems = newStems
	}

	var res []string
	for s := range stems {
		if checker.Valid(stems[s]) {
			res = append(res, stems[s])
		}
	}
	sort.Strings(res)
	return res, maxLength
}
