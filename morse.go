package kowalski

import (
	"fmt"
	"regexp"
	"strings"
)

var morseLetters = map[string]rune {
	".-": 'a',
	"-...": 'b',
	"-.-.": 'c',
	"-..": 'd',
	".": 'e',
	"..-.": 'f',
	"--.": 'g',
	"....": 'h',
	"..": 'i',
	".---": 'j',
	"-.-": 'k',
	".-..": 'l',
	"--": 'm',
	"-.": 'n',
	"---": 'o',
	".--.": 'p',
	"--.-": 'q',
	".-.": 'r',
	"...": 's',
	"-": 't',
	"..-": 'u',
	"...-": 'v',
	".--": 'w',
	"-..-": 'x',
	"-.--": 'y',
	"--..": 'z',
}

var nonMorseRegexp = regexp.MustCompile("[^.\\-]")

// FromMorse takes a sequence of morse signals (as ASCII dots and hyphens) and returns a set of possible words
// that could be constructed from them.
func (n *Node) FromMorse(input string) []string {
	return n.fromMorse(nonMorseRegexp.ReplaceAllString(input, ""), "")
}

func (n *Node) fromMorse(input string, prefix string) []string {
	var res []string

	for p := range morseLetters {
		if strings.HasPrefix(input, p) {
			next := fmt.Sprintf("%s%c", prefix, morseLetters[p])
			left := input[len(p):]
			if len(left) == 0 && n.Valid(next) {
				res = append(res, next)
			} else if len(left) > 0 && n.IsPrefix(next) {
				res = append(res, n.fromMorse(left, next)...)
			}
		}
	}

	return res
}
