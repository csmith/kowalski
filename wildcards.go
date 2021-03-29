package kowalski

import (
	"context"
	"fmt"
	"sort"
	"strings"
)

// Match returns all valid words that match the given pattern, expanding '?' as a single character wildcard
func Match(ctx context.Context, checker *SpellChecker, pattern string) ([]string, error) {
	res, _, err := findMatch(ctx, checker, strings.ToLower(pattern), false, 0)
	return res, err
}

// MultiMatch returns valid sequences of words that match the given pattern, expanding '?' as a single character
// wildcard. To reduce the search space, multi-match will first try to look for matches consisting only of longer
// words, then gradually reduce that threshold until at least one match is found.
func MultiMatch(ctx context.Context, checker *SpellChecker, pattern string) ([]string, error) {
	i := len(pattern) / 2
	if i > 5 {
		i = 5
	}

	for i > 0 {
		res, _, err := findMatch(ctx, checker, strings.ToLower(pattern), true, i)
		if err != nil {
			return nil, err
		}

		if len(res) > 0 {
			return res, nil
		}
		i--
	}
	return nil, nil
}

// OffByOne returns all words that can be made by performing one character change on the input. The input is
// assumed to be a single, lowercase word containing a-z chars only.
func OffByOne(ctx context.Context, checker *SpellChecker, input string) ([]string, error) {
	words := map[string]bool{}
	for i := range input {
		res, err := Match(ctx, checker, fmt.Sprintf("%s?%s", input[0:i], input[i+1:]))
		if err != nil {
			return nil, err
		}

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
	return res, nil
}

// findMatch returns all valid words that match the given pattern, expanding '?' as a single character wildcard.
// It will aggressively skip sequences that don't form valid prefixes; the maximum valid prefix length is returned as
// the second parameter (for cases where matches are returned, this will equal len(word)).
func findMatch(ctx context.Context, checker *SpellChecker, word string, multiWord bool, minLength int) ([]string, int, error) {
	maxLength := 0
	stems := [][]string{{""}}
	for offset := 0; offset < len(word) && len(stems) > 0; offset++ {
		if ctx.Err() != nil {
			return nil, 0, ctx.Err()
		}

		newStems := make([][]string, 0, len(stems))

		var chars []uint8
		if word[offset] == '?' {
			chars = []uint8("abcdefghijklmnopqrstuvwxyz")
		} else {
			chars = []uint8{word[offset]}
		}

		for _, nextChar := range chars {
			if ctx.Err() != nil {
				return nil, 0, ctx.Err()
			}

			for s := range stems {
				stem := stems[s]
				if newWord := fmt.Sprintf("%s%c", stem[len(stem)-1], nextChar); checker.Prefix(newWord) {
					var newStem []string
					if len(stem) > 1 {
						newStem = append([]string{}, stem[0:len(stem)-1]...)
					}
					newStem = append(newStem, newWord)

					newStems = append(newStems, newStem)
					if multiWord && checker.Valid(newWord) && len(newWord) >= minLength {
						newStems = append(newStems, append(append([]string{}, newStem...), ""))
					}
				}
			}
		}

		maxLength = offset
		stems = newStems
	}

	var res []string
	for s := range stems {
		stem := stems[s]
		if len(stem[len(stem)-1]) >= minLength && checker.Valid(stem[len(stem)-1]) {
			res = append(res, strings.Join(stems[s], " "))
		}
	}
	sort.Strings(res)
	return res, maxLength, nil
}
