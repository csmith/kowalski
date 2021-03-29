package kowalski

import (
	"context"
	"sort"
	"strings"
)

// Anagram finds all anagrams of the given word, expanding '?' as a single wildcard character
func Anagram(ctx context.Context, checker *SpellChecker, word string) ([]string, error) {
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
		matches, count, err := findMatch(ctx, checker, string(w))
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

	sort.Strings(res)
	return unique(res), nil
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
