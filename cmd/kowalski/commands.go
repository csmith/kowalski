package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/csmith/kowalski/v3"
)

type Replier func(format string, a ...interface{})
type TextCommand func(input string, reply Replier)

type HelpInfo struct {
	Triggers []string
	Message  string
}

var textCommands = map[string]TextCommand{}
var help []HelpInfo

func addTextCommand(c TextCommand, helpText string, names ...string) {
	for i := range names {
		textCommands[names[i]] = c
	}

	help = append(help, HelpInfo{
		Triggers: names,
		Message:  helpText,
	})
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

func init() {
	addTextCommand(Anagram, "Attempts to find single-word anagrams, expanding '\\*' and '?' wildcards", "anagram")
}

func Analysis(input string, reply Replier) {
	input = strings.ToLower(input)
	res := kowalski.Analyse(checkers[0], input)
	if len(res) == 0 {
		reply("Analysis: nothing interesting found")
	}
	reply("Analysis:\n\t%s", strings.Join(res, "\n\t"))
}

func init() {
	addTextCommand(Analysis, "Analyses text and provides a summary of potentially interesting findings", "analysis", "analyze", "analyse")
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

func init() {
	addTextCommand(Chunk, "Splits the text into chunks of a given size", "chunk")
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

func init() {
	addTextCommand(Letters, "Shows a frequency histogram of the number of letters in the input", "letters")
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

func init() {
	addTextCommand(Match, "Attempts to expand '\\*' and '?' wildcards to find a single-word match", "match")
}

func Morse(input string, reply Replier) {
	res := merge(kowalski.MultiplexFromMorse(checkers, input, kowalski.Dedupe))
	reply("Matches for %s: %v", input, res)
}

func init() {
	addTextCommand(Morse, "Attempts to split a morse code input to spell a single word", "morse")
}

func OffByOne(input string, reply Replier) {
	input = strings.ToLower(input)
	if isValidWord(input) {
		res := merge(kowalski.MultiplexOffByOne(checkers, input, kowalski.Dedupe))
		reply("Off-by-ones for %s: %v", input, res)
	} else {
		reply("Invalid word: %s", input)
	}
}

func init() {
	addTextCommand(OffByOne, "Finds all words that are one character different from the input", "obo", "offbyone", "ob1")
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

func init() {
	addTextCommand(Shift, "Shows the result of the 25 possible caesar shifts", "shift", "caesar")
}

func T9(input string, reply Replier) {
	if isValidT9(input) {
		res := merge(kowalski.MultiplexFromT9(checkers, input, kowalski.Dedupe))
		reply("Matches for %s: %v", input, res)
	} else {
		reply("Invalid word: %s", input)
	}
}

func init() {
	addTextCommand(T9, "Attempts to treat a series of numbers as T9 input to spell a single word", "t9")
}

func Transpose(input string, reply Replier) {
	reply("Transposed:\n\n%s", strings.Join(kowalski.Transpose(strings.Split(input, "\n")), "\n"))
}

func init() {
	addTextCommand(Transpose, "Transposes columns to rows and rows to columns", "transpose")
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

func init() {
	addTextCommand(WordSearch, "Searches for words in the given text grid", "wordsearch")
}

func Help(_ string, reply Replier) {
	helpText := strings.Builder{}
	for i := range help {
		helpText.WriteString(fmt.Sprintf("\n\t**%s%s**", *prefix, help[i].Triggers[0]))
		helpText.WriteString(fmt.Sprintf(" _%s_", help[i].Message))
		if len(help[i].Triggers) > 1 {
			helpText.WriteString(" [Aliases: ")
			for j := range help[i].Triggers[1:] {
				if j > 0 {
					helpText.WriteString(", ")
				}
				helpText.WriteString(fmt.Sprintf("%s%s", *prefix, help[i].Triggers[1+j]))
			}
			helpText.WriteString("]")
		}
	}

	reply("Help:%s", helpText.String())
}

func init() {
	addTextCommand(Help, "Shows this help text", "help")
}
