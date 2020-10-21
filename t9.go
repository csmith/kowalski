package kowalski

import "fmt"

var t9mapping = map[uint8][]rune {
	'2': []rune("abc"),
	'3': []rune("def"),
	'4': []rune("ghi"),
	'5': []rune("jkl"),
	'6': []rune("mno"),
	'7': []rune("pqrs"),
	'8': []rune("tuv"),
	'9': []rune("wxyz"),
}

// FromT9 takes an input that represents a sequence of key presses on a T9 keyboard and returns possible
// words that match. The input should not contain spaces (the "0" digit) - words should be solved independently,
// to avoid an explosion of possible results.
func FromT9(checker *SpellChecker, input string) []string {
	return fromT9(checker, input, "")
}

func fromT9(checker *SpellChecker, input, prefix string) []string {
	var res []string
	if opts, ok := t9mapping[input[0]]; ok {
		for i := range opts {
			next := fmt.Sprintf("%s%c", prefix, opts[i])
			if len(input) > 1 && checker.Prefix(next) {
				res = append(res, fromT9(checker, input[1:], next)...)
			} else if len(input) == 1 && checker.Valid(next) {
				res = append(res, next)
			}
		}
	}
	return res
}
