package fst

import (
	"github.com/blevesearch/vellum"
	"github.com/blevesearch/vellum/regexp"
)

// regexProgramLimit bounds the compiled size of a regex. vellum builds the whole DFA up front, and each DFA state
// costs time proportional to the program size, so nested repeats under the default 10MB limit can take minutes to
// compile.
const regexProgramLimit = 1 << 20

// NewRegexAutomaton creates a vellum Automaton that matches the given regular expression. Expressions that compile
// to a large program are rejected.
func NewRegexAutomaton(pattern string) (vellum.Automaton, error) {
	return regexp.NewWithLimit(pattern, regexProgramLimit)
}
