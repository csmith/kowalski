package fst

import (
	"strings"

	"github.com/blevesearch/vellum"
)

var morseLetters = map[byte]string{
	'a': ".-",
	'b': "-...",
	'c': "-.-.",
	'd': "-..",
	'e': ".",
	'f': "..-.",
	'g': "--.",
	'h': "....",
	'i': "..",
	'j': ".---",
	'k': "-.-",
	'l': ".-..",
	'm': "--",
	'n': "-.",
	'o': "---",
	'p': ".--.",
	'q': "--.-",
	'r': ".-.",
	's': "...",
	't': "-",
	'u': "..-",
	'v': "...-",
	'w': ".--",
	'x': "-..-",
	'y': "-.--",
	'z': "--..",
}

const errorSentinel = -1

type morseAutomaton struct {
	chars string
}

// NewMorseAutomaton creates a vellum Automaton that matches morse code input.
// Any character other than '.' and '-' is ignored.
func NewMorseAutomaton(term string) vellum.Automaton {
	// Remove anything except - and .
	pruned := strings.Builder{}
	for _, c := range term {
		if c == '.' || c == '-' {
			pruned.WriteRune(c)
		}
	}

	return &morseAutomaton{
		chars: pruned.String(),
	}
}

func (m *morseAutomaton) Start() int {
	return 0
}

func (m *morseAutomaton) IsMatch(i int) bool {
	return i == len(m.chars)
}

func (m *morseAutomaton) CanMatch(i int) bool {
	return i <= len(m.chars) && i != errorSentinel
}

func (m *morseAutomaton) WillAlwaysMatch(i int) bool {
	return false
}

func (m *morseAutomaton) Accept(i int, b byte) int {
	if b == ' ' {
		// Skip over spaces
		return i
	}

	if b >= 'A' && b <= 'Z' {
		b += 32
	}

	morse, ok := morseLetters[b]
	if !ok {
		return errorSentinel
	}

	remainder := m.chars[i:]
	if strings.HasPrefix(remainder, morse) {
		return i + len(morse)
	}

	return errorSentinel
}
