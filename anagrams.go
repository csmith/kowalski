package kowalski

import (
	"context"
	"sort"
	"strings"
)

// Anagram finds all single-word anagrams of the given word, expanding '?' as a single wildcard character
func Anagram(ctx context.Context, checker *SpellChecker, word string) ([]string, error) {
	return anagram(ctx, checker, word, false, 0)
}

// MultiAnagram finds all single- and multi-word anagrams of the given word, expanding '?' as a single wildcard
// character. To avoid duplicates, words are sorted lexicographically (i.e., "a ball" will be returned and "ball a"
// will not).
func MultiAnagram(ctx context.Context, checker *SpellChecker, word string) ([]string, error) {
	// TODO: Allow configuring of the min length
	return anagram(ctx, checker, word, true, 2)
}

func anagram(ctx context.Context, checker *SpellChecker, word string, multiWord bool, minLength int) ([]string, error) {
	var (
		res        []string
		swapBefore = len(word)
	)

	sortedWord := func(w string) string {
		s := strings.Split(strings.ToLower(w), "")
		sort.Strings(s)
		return strings.Join(s, "")
	}(word)

	for w := []byte(sortedWord); w != nil; w = permute(w, swapBefore+1) {
		matches, count, err := findMatch(ctx, checker, string(w), multiWord, minLength)
		if err != nil {
			return nil, err
		}

		if len(matches) > 0 {
			res = append(res, matches...)
			swapBefore = len(word)
		} else {
			swapBefore = count
		}
	}

	if multiWord {
		res = onlyAscendingWords(res)
	}

	sort.Strings(res)
	return unique(res), nil
}

// onlyAscendingWords returns a slice that contains all the entries from input that meet the following criteria:
// * the entry is a single word
// * the entry is multiple-space separated words in increasing lexicographical order
func onlyAscendingWords(input []string) []string {
	var filtered []string
	for i := range input {
		parts := strings.Split(input[i], " ")
		if len(parts) > 1 {
			last := parts[0]
			valid := true
			for _, p := range parts[1:] {
				if strings.Compare(last, p) == 1 {
					valid = false
					break
				}
			}
			if valid {
				filtered = append(filtered, input[i])
			}
		} else {
			filtered = append(filtered, input[i])
		}
	}
	return filtered
}

// permute returns the next permutation of the given input, in lexicographical order.
// swapBefore can be used to force a swap within a certain number characters (i.e., skip all permutations that
// affect characters after the one with index swapBefore).
func permute(input []byte, swapBefore int) []byte {
	if swapBefore < len(input)-1 {
		input = append(input[0:swapBefore], func(w []byte) []byte {
			s := strings.Split(string(w), "")
			sort.Strings(s)
			return reverse([]byte(strings.Join(s, "")), 0)
		}(input[swapBefore:])...)
	}

	k, l := -1, -1
	for i := range input {
		if i+1 < len(input) && input[i] < input[i+1] {
			k = i
			l = -1
		} else if k >= 0 && input[k] < input[i] {
			l = i
		}
	}

	if k == -1 {
		return nil
	}

	input[k], input[l] = input[l], input[k]
	return reverse(input, k+1)
}
