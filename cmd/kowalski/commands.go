package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/csmith/kowalski/v3"
)

type Replier func(format string, a ...interface{})
type Command func(input string, reply Replier)

var commands = map[string]Command{}
var helpText = strings.Builder{}

func addCommand(c Command, help string, names ...string) {
	for i := range names {
		commands[names[i]] = c
	}

	helpText.WriteString(fmt.Sprintf("\n\t**%s%s**", *prefix, names[0]))
	helpText.WriteString(fmt.Sprintf(" _%s_", help))
	if len(names) > 1 {
		helpText.WriteString(" [Aliases: ")
		for i := range names[1:] {
			if i > 0 {
				helpText.WriteString(", ")
			}
			helpText.WriteString(fmt.Sprintf("%s%s", *prefix, names[1+i]))
		}
		helpText.WriteString("]")
	}
}

func init() {
	addCommand(Anagram, "Attempts to find single-word anagrams, expanding '\\*' and '?' wildcards", "anagram")
	addCommand(Analysis, "Analyses text and provides a summary of potentially interesting findings", "analysis", "analyze", "analyse")
	addCommand(Chunk, "Splits the text into chunks of a given size", "chunk")
	addCommand(Letters, "Shows a frequency histogram of the number of letters in the input", "letters")
	addCommand(Match, "Attempts to expand '\\*' and '?' wildcards to find a single-word match", "match")
	addCommand(Morse, "Attempts to split a morse code input to spell a single word", "morse")
	addCommand(Shift, "Shows the result of the 25 possible caesar shifts", "shift", "caesar")
	addCommand(T9, "Attempts to treat a series of numbers as T9 input to spell a single word", "t9")
	addCommand(Transpose, "Transposes columns to rows and rows to columns", "transpose")
	addCommand(WordSearch, "Searches for words in the given text grid", "wordsearch")
	addCommand(Help, "Shows this help text", "help")
}

func Anagram(input string, reply Replier) {
	input = strings.ToLower(input)
	if isValidWord(input) {
		res := merge(kowalski.MultiplexAnagram(checkers, input, kowalski.Dedupe))
		reply("Anagrams for %s: %v", input, res)
	} else {
		reply("Invalid word: %s", input)
	}
}

func Analysis(input string, reply Replier) {
	input = strings.ToLower(input)
	res := kowalski.Analyse(checkers[0], input)
	if len(res) == 0 {
		reply("Analysis: nothing interesting found")
	}
	reply("Analysis:\n\t%s", strings.Join(res, "\n\t"))
}

func Chunk(input string, reply Replier) {
	var parts []int
	words := strings.Split(input, " ")
	for i := range words {
		if v, err := strconv.Atoi(words[i]); err == nil {
			parts = append(parts, v)
		} else {
			break
		}
	}

	if len(parts) == 0 {
		reply("Usage: chunk <size> [size [size [...]]] <text>")
		return
	}

	text := strings.Join(words[len(parts):], "")
	reply("Chunked: %s", strings.Join(kowalski.Chunk(text, parts...), " "))
}

func Help(_ string, reply Replier) {
	reply("Help:%s", helpText.String())
}

func Letters(input string, reply Replier) {
	res := kowalski.LetterDistribution(input)
	max := 0
	for i := range res {
		if res[i] > max {
			max = res[i]
		}
	}
	const targetWidth = 20
	message := strings.Builder{}
	message.WriteString("Letter distribution:\n```")
	for i := range res {
		message.WriteByte(byte(i + 'A'))
		message.WriteString(": ")
		if res[i] > 0 {
			message.WriteRune('▕')
		}
		for j := 0; j < int(targetWidth*(float64(res[i])/float64(max))); j++ {
			message.WriteRune('█')
		}
		message.WriteString(fmt.Sprintf(" %d\n", res[i]))
	}
	message.WriteString("```")
	reply(message.String())
}

func Match(input string, reply Replier) {
	input = strings.ToLower(input)
	if isValidWord(input) {
		res := merge(kowalski.MultiplexMatch(checkers, input, kowalski.Dedupe))
		reply("Matches for %s: %v", input, res)
	} else {
		reply("Invalid word: %s", input)
	}
}

func Morse(input string, reply Replier) {
	res := merge(kowalski.MultiplexFromMorse(checkers, input, kowalski.Dedupe))
	reply("Matches for %s: %v", input, res)
}

func Shift(input string, reply Replier) {
	res := kowalski.CaesarShifts(input)
	out := strings.Builder{}
	out.WriteString("Caesar shifts:\n")
	for i, s := range res {
		score := kowalski.Score(checkers[0], s)
		if score > 0.5 {
			s = fmt.Sprintf("**%s**", s)
		}
		out.WriteString(fmt.Sprintf("\t%2d: %s\n", i+1, s))
	}
	reply(out.String())
}

func T9(input string, reply Replier) {
	if isValidT9(input) {
		res := merge(kowalski.MultiplexFromT9(checkers, input, kowalski.Dedupe))
		reply("Matches for %s: %v", input, res)
	} else {
		reply("Invalid word: %s", input)
	}
}

func Transpose(input string, reply Replier) {
	reply("Transposed:\n\n%s", strings.Join(kowalski.Transpose(strings.Split(input, "\n")), "\n"))
}

func WordSearch(input string, reply Replier) {
	input = strings.ToLower(input)
	res := kowalski.MultiplexWordSearch(checkers, strings.Split(input, "\n"))
	reply(
		"Words found:\n\nNormal: %s\n\nUD: %s",
		strings.Join(countReps(res[0]), ", "),
		strings.Join(countReps(subtract(res[1], res[0])), ", "),
	)
}
