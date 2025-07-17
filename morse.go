package kowalski

import (
	"fmt"
	"regexp"
	"strings"
)

var morseLetters = map[string]rune{
	".-":   'a',
	"-...": 'b',
	"-.-.": 'c',
	"-..":  'd',
	".":    'e',
	"..-.": 'f',
	"--.":  'g',
	"....": 'h',
	"..":   'i',
	".---": 'j',
	"-.-":  'k',
	".-..": 'l',
	"--":   'm',
	"-.":   'n',
	"---":  'o',
	".--.": 'p',
	"--.-": 'q',
	".-.":  'r',
	"...":  's',
	"-":    't',
	"..-":  'u',
	"...-": 'v',
	".--":  'w',
	"-..-": 'x',
	"-.--": 'y',
	"--..": 'z',
}

var nonMorseRegexp = regexp.MustCompile("[^.\\-]")

// FromMorse takes a sequence of morse signals (as ASCII dots and hyphens) and returns a set of possible words
// that could be constructed from them.
func FromMorse(checker *SpellChecker, input string) []string {
	return fromMorse(checker, nonMorseRegexp.ReplaceAllString(input, ""), "")
}

func fromMorse(checker *SpellChecker, input string, prefix string) []string {
	var res []string

	for p := range morseLetters {
		if strings.HasPrefix(input, p) {
			next := fmt.Sprintf("%s%c", prefix, morseLetters[p])
			left := input[len(p):]
			if len(left) == 0 && checker.Valid(next) {
				res = append(res, next)
			} else if len(left) > 0 && checker.Prefix(next) {
				res = append(res, fromMorse(checker, left, next)...)
			}
		}
	}

	return res
}
