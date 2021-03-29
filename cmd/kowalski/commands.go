package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/csmith/kowalski/v4"
)

type Replier interface {
	reply(format string, a ...interface{})
	replyWithFiles(files []*discordgo.File, format string, a ...interface{})
}

type TextCommand func(input string, r Replier)
type FileCommand func(input string, urls []string, r Replier)

type HelpInfo struct {
	Triggers []string
	Message  string
}

var textCommands = map[string]TextCommand{}
var fileCommands = map[string]FileCommand{}
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

func addFileCommand(c FileCommand, helpText string, names ...string) {
	for i := range names {
		fileCommands[names[i]] = c
	}

	help = append(help, HelpInfo{
		Triggers: names,
		Message:  helpText,
	})
}

func Anagram(input string, r Replier) {
	input = strings.ToLower(input)
	if isValidWord(input) {
		ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
		defer cancel()

		words, err := kowalski.MultiplexAnagram(ctx, checkers, input, kowalski.Dedupe)
		if err != nil {
			r.reply("Error: %v", err)
		} else {
			r.reply("Anagrams for %s: %v", input, merge(words))
		}
	} else {
		r.reply("Invalid word: %s", input)
	}
}

func init() {
	addTextCommand(Anagram, "Attempts to find single-word anagrams, expanding '\\*' and '?' wildcards", "anagram")
}

func Analysis(input string, r Replier) {
	input = strings.ToLower(input)
	res := kowalski.Analyse(checkers[0], input)
	if len(res) == 0 {
		r.reply("Analysis: nothing interesting found")
	}
	r.reply("Analysis:\n\t%s", strings.Join(res, "\n\t"))
}

func init() {
	addTextCommand(Analysis, "Analyses text and provides a summary of potentially interesting findings", "analysis", "analyze", "analyse")
}

func Chunk(input string, r Replier) {
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
		r.reply("Usage: chunk <size> [size [size [...]]] <text>")
		return
	}

	text := strings.Join(words[len(parts):], "")
	r.reply("Chunked: %s", strings.Join(kowalski.Chunk(text, parts...), " "))
}

func init() {
	addTextCommand(Chunk, "Splits the text into chunks of a given size", "chunk")
}

func Colours(_ string, urls []string, r Replier) {
	res, err := http.Get(urls[0])
	if err != nil {
		r.reply("Unable to download image: %v", err)
		return
	}

	defer res.Body.Close()
	colours, err := kowalski.ExtractColours(res.Body)
	if err != nil {
		r.reply("Unable to decode image: %v", err)
		return
	}

	text := strings.Builder{}
	text.WriteString(fmt.Sprintf("%d colours found:\n```\nHex         R   G   B   A Pixels\n", len(colours)))
	for i := range colours {
		if i >= 25 {
			text.WriteString("... truncated ...\n")
			break
		}

		r, g, b, a := colours[i].Colour.RGBA()
		if a == 65535 {
			text.WriteString(fmt.Sprintf("#%02x%02x%02x   %3[1]d %3[2]d %3[3]d   - %d\n", r/257, g/257, b/257, colours[i].Count))
		} else {
			text.WriteString(fmt.Sprintf("#%02x%02x%02x#%02x %3[1]d %3[2]d %3[3]d %3[4]d %d\n", r/257, g/257, b/257, a/257, colours[i].Count))
		}
	}
	text.WriteString("```")
	r.reply(text.String())
}

func init() {
	addFileCommand(Colours, "Counts the colours within the image", "colours", "colors")
}

func Letters(input string, r Replier) {
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
	r.reply(message.String())
}

func init() {
	addTextCommand(Letters, "Shows a frequency histogram of the number of letters in the input", "letters")
}

func Match(input string, r Replier) {
	input = strings.ToLower(input)
	if isValidWord(input) {
		ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
		defer cancel()

		words, err := kowalski.MultiplexMatch(ctx, checkers, input, kowalski.Dedupe)
		if err != nil {
			r.reply("Error: %v", err)
		} else {
			r.reply("Matches for %s: %v", input, merge(words))
		}
	} else {
		r.reply("Invalid word: %s", input)
	}
}

func init() {
	addTextCommand(Match, "Attempts to expand '?' wildcards to find a single-word match", "match")
}

func Morse(input string, r Replier) {
	res := merge(kowalski.MultiplexFromMorse(checkers, input, kowalski.Dedupe))
	r.reply("Matches for %s: %v", input, res)
}

func init() {
	addTextCommand(Morse, "Attempts to split a morse code input to spell a single word", "morse")
}

func MultiAnagram(input string, r Replier) {
	input = strings.ToLower(input)
	if isValidWord(input) {
		ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
		defer cancel()

		words, err := kowalski.MultiplexMultiAnagram(ctx, checkers, input, kowalski.Dedupe)
		if err != nil {
			r.reply("Error: %v", err)
		} else {
			r.reply("Multi anagrams for %s: %v", input, strings.Join(merge(words), ", "))
		}
	} else {
		r.reply("Invalid word: %s", input)
	}
}

func init() {
	addTextCommand(MultiAnagram, "Attempts to find multi-word anagrams, expanding '?' wildcards", "multigram", "multianagram")
}

func MultiMatch(input string, r Replier) {
	input = strings.ToLower(input)
	if isValidWord(input) {
		ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
		defer cancel()

		words, err := kowalski.MultiplexMultiMatch(ctx, checkers, input, kowalski.Dedupe)
		if err != nil {
			r.reply("Error: %v", err)
		} else {
			r.reply("Multi matches for %s: %s", input, strings.Join(merge(words), ", "))
		}
	} else {
		r.reply("Invalid word: %s", input)
	}
}

func init() {
	addTextCommand(MultiMatch, "Attempts to expand '?' wildcards to find multi-word matches", "multimatch")
}

func OffByOne(input string, r Replier) {
	input = strings.ToLower(input)
	if isValidWord(input) {
		ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
		defer cancel()

		words, err := kowalski.MultiplexOffByOne(ctx, checkers, input, kowalski.Dedupe)
		if err != nil {
			r.reply("Error: %v", err)
		} else {
			r.reply("Off-by-ones for %s: %v", input, merge(words))
		}
	} else {
		r.reply("Invalid word: %s", input)
	}
}

func init() {
	addTextCommand(OffByOne, "Finds all words that are one character different from the input", "obo", "offbyone", "ob1")
}

func RGB(_ string, urls []string, r Replier) {
	res, err := http.Get(urls[0])
	if err != nil {
		r.reply("Unable to download image: %v", err)
		return
	}

	defer res.Body.Close()
	red, green, blue, err := kowalski.SplitRGB(res.Body)
	if err != nil {
		r.reply("Unable to split image: %v", err)
		return
	}

	r.replyWithFiles([]*discordgo.File{
		{
			Name:        "red.png",
			ContentType: "image/png",
			Reader:      red,
		},
		{
			Name:        "green.png",
			ContentType: "image/png",
			Reader:      green,
		},
		{
			Name:        "blue.png",
			ContentType: "image/png",
			Reader:      blue,
		},
	}, "")
}

func init() {
	addFileCommand(RGB, "Splits an image into its red, green and blue channels", "rgb")
}

func Shift(input string, r Replier) {
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
	r.reply(out.String())
}

func init() {
	addTextCommand(Shift, "Shows the result of the 25 possible caesar shifts", "shift", "caesar")
}

func T9(input string, r Replier) {
	if isValidT9(input) {
		res := merge(kowalski.MultiplexFromT9(checkers, input, kowalski.Dedupe))
		r.reply("Matches for %s: %v", input, res)
	} else {
		r.reply("Invalid word: %s", input)
	}
}

func init() {
	addTextCommand(T9, "Attempts to treat a series of numbers as T9 input to spell a single word", "t9")
}

func Transpose(input string, r Replier) {
	r.reply("Transposed:\n\n%s", strings.Join(kowalski.Transpose(strings.Split(input, "\n")), "\n"))
}

func init() {
	addTextCommand(Transpose, "Transposes columns to rows and rows to columns", "transpose")
}

func WordSearch(input string, r Replier) {
	input = strings.ToLower(input)
	res := kowalski.MultiplexWordSearch(checkers, strings.Split(input, "\n"))
	r.reply(
		"Words found:\n\nNormal: %s\n\nUD: %s",
		strings.Join(countReps(res[0]), ", "),
		strings.Join(countReps(subtract(res[1], res[0])), ", "),
	)
}

func init() {
	addTextCommand(WordSearch, "Searches for words in the given text grid", "wordsearch")
}

func Help(_ string, r Replier) {
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

	r.reply("Help:%s", helpText.String())
}

func init() {
	addTextCommand(Help, "Shows this help text", "help")
}
