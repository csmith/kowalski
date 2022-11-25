package main

import (
	"errors"
	"flag"
	"fmt"
	regexp2 "regexp"
	"sort"
	"strings"

	"github.com/blevesearch/vellum"
	"github.com/blevesearch/vellum/regexp"
	"github.com/csmith/kowalski/v5/fst"
)

var fstModel = flag.String("fst-model", "", "Path to FST for fast word operations")

var transducer *vellum.FST

func init() {
	flag.Parse()

	if *fstModel != "" {
		var err error

		transducer, err = vellum.Open(*fstModel)
		if err != nil {
			panic(err)
		}
	}
}

func FstAnagram(input string, r Replier) {
	automaton, err := fst.NewAnagramAutomaton(input)
	if err != nil {
		r.reply("Error: %s", err.Error())
		return
	}

	fstQuery(automaton, input, r)
}

func init() {
	if *fstModel != "" {
		addCommand(textCommands, FstAnagram, "Attempts to find anagrams from wikipedia, expanding '\\*' wildcards", "fstanagram", "fstagram")
	}
}

func FstRegex(input string, r Replier) {
	automaton, err := regexp.New(input)
	if err != nil {
		r.reply("Error: %s", err.Error())
		return
	}

	fstQuery(automaton, input, r)
}

func init() {
	if *fstModel != "" {
		addCommand(textCommands, FstRegex, "Attempts to find word matches from wikipedia using regexp", "fstregex", "fstre")
	}
}

func FstMorse(input string, r Replier) {
	fstQuery(fst.NewMorseAutomaton(input), input, r)
}

func init() {
	if *fstModel != "" {
		addCommand(textCommands, FstMorse, "Attempts to find word matches from wikipedia using morse", "fstmorse")
	}
}

func WordLink(input string, r Replier) {
	parts := strings.Split(input, " ")
	if len(parts) != 2 {
		r.reply("Invalid input, must specify two words to link")
		return
	}

	leftAut, err := regexp.New(fmt.Sprintf("%s[ -]?([a-zA-Z]{3,})", regexp2.QuoteMeta(parts[0])))
	if err != nil {
		r.reply("Error: %s", err.Error())
		return
	}
	leftSide, err := fstResults(leftAut)
	if err != nil {
		r.reply("Error: %s", err.Error())
		return
	} else if len(leftSide) == 0 {
		r.reply("No words found that can follow '%s'", parts[0])
		return
	}

	rightAut, err := regexp.New(fmt.Sprintf("([a-zA-Z]{3,})[ -]?%s", regexp2.QuoteMeta(parts[1])))
	if err != nil {
		r.reply("Error: %s", err.Error())
		return
	}

	rightSide, err := fstResults(rightAut)
	if err != nil {
		r.reply("Error: %s", err.Error())
		return
	} else if len(leftSide) == 0 {
		r.reply("No words found that can precede '%s'", parts[1])
		return
	}

	var leftScores = make(map[string]uint64)
	for i := range leftSide {
		leftTerm := strings.TrimSpace(strings.TrimPrefix(leftSide[i].Term, parts[0]))
		if normal, _, _ := transducer.Get([]byte(leftTerm)); normal < 5_000_000 {
			leftScores[leftTerm] += leftSide[i].Score
		}
	}

	var rightScores = make(map[string]uint64)
	for i := range rightSide {
		rightTerm := strings.TrimSpace(strings.TrimSuffix(rightSide[i].Term, parts[1]))
		if normal, _, _ := transducer.Get([]byte(rightTerm)); normal < 5_000_000 {
			rightScores[rightTerm] += rightSide[i].Score
		}
	}

	var combined []fstMatch
	for rightTerm := range rightScores {
		if leftScore, ok := leftScores[rightTerm]; ok {
			var score uint64
			if leftScore < rightScores[rightTerm] {
				score = leftScore
			} else {
				score = rightScores[rightTerm]
			}

			if score > 2 {
				combined = append(combined, fstMatch{
					Term:  rightTerm,
					Score: score,
				})
			}
		}
	}

	if len(combined) == 0 {
		r.reply("No linking words found")
		return
	}

	sort.Slice(combined, func(i, j int) bool {
		return combined[i].Score > combined[j].Score
	})

	message := strings.Builder{}
	message.WriteString("Linking words for '")
	message.WriteString(parts[0])
	message.WriteString("' <> '")
	message.WriteString(parts[1])
	message.WriteString("': ")

	for i := range combined {
		message.WriteString(fmt.Sprintf("`%s` (%d) ", combined[i].Term, combined[i].Score))
		if message.Len() > 1900 {
			message.WriteString("[...]")
			break
		}
	}

	r.reply("%s", message.String())
}

func init() {
	if *fstModel != "" {
		addCommand(textCommands, WordLink, "Attempts to find a word that links two others", "wordlink", "link")
	}
}

type fstMatch struct {
	Term  string
	Score uint64
}

func fstQuery(automaton vellum.Automaton, input string, r Replier) {
	matches, err := fstResults(automaton)
	if err != nil {
		r.reply("Error: %s", err.Error())
		return
	} else if len(matches) == 0 {
		r.reply("No results found")
		return
	}

	message := strings.Builder{}
	message.WriteString("Matches for '")
	message.WriteString(input)
	message.WriteString("': ")

	for i := range matches {
		message.WriteString(fmt.Sprintf("`%s` (%d) ", matches[i].Term, matches[i].Score))
		if message.Len() > 1900 {
			message.WriteString("[...]")
			break
		}
	}

	r.reply("%s", message.String())
}

func fstResults(automaton vellum.Automaton) ([]fstMatch, error) {
	iterator, err := transducer.Search(automaton, nil, nil)
	if err != nil {
		if errors.Is(err, vellum.ErrIteratorDone) {
			return nil, nil
		}

		return nil, err
	}

	var matches []fstMatch
	for err == nil {
		key, val := iterator.Current()
		matches = append(matches, fstMatch{Term: string(key), Score: val})
		err = iterator.Next()
	}

	sort.Slice(matches, func(i, j int) bool {
		return matches[i].Score > matches[j].Score
	})

	return matches, nil
}
