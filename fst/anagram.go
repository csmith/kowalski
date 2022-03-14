package fst

import (
	"fmt"
	"strings"

	"github.com/blevesearch/vellum"
)

const errorMask = 1 << 62

type anagramAutomaton struct {
	Chars string
}

// NewAnagramAutomaton creates a vellum Automaton that will match the input letters in any order.
// Single digit wildcards ("*") are supported. At most 61 letters or wildcards may be specified.
func NewAnagramAutomaton(term string) (vellum.Automaton, error) {
	// Ensure all wildcards are at the end of the terms
	term = strings.ToLower(term)
	wildcards := strings.Count(term, "*")
	rearranged := strings.ReplaceAll(term, "*", "") + strings.Repeat("*", wildcards)
	if len(term) >= 62 {
		return nil, fmt.Errorf("term is too long: %d (max length is 61)", len(term))
	}

	return &anagramAutomaton{
		Chars: rearranged,
	}, nil
}

func (a *anagramAutomaton) Start() int {
	r := 0
	for i := range a.Chars {
		r |= 1 << i
	}
	return r
}

func (a *anagramAutomaton) IsMatch(i int) bool {
	return i == 0
}

func (a *anagramAutomaton) CanMatch(i int) bool {
	return (i & errorMask) == 0
}

func (a *anagramAutomaton) WillAlwaysMatch(i int) bool {
	return false
}

func (a *anagramAutomaton) Accept(i int, b byte) int {
	if b == ' ' {
		// Skip over spaces
		return i
	}

	for j := range a.Chars {
		if a.Chars[j] == b || a.Chars[j] == '*' {
			mask := 1 << j
			if (mask & i) == mask {
				// Not yet used
				return i ^ mask
			}
		}
	}

	return i | errorMask
}
