package main

import (
	"errors"
	"flag"
	"fmt"
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

func fstQuery(automaton vellum.Automaton, input string, r Replier) {
	iterator, err := transducer.Search(automaton, nil, nil)
	if err != nil {
		if errors.Is(err, vellum.ErrIteratorDone) {
			r.reply("No results found")
			return
		}

		r.reply("Error: %s", err.Error())
		return
	}

	type Match struct {
		Term  string
		Score uint64
	}

	var matches []Match
	for err == nil {
		key, val := iterator.Current()
		matches = append(matches, Match{Term: string(key), Score: val})
		err = iterator.Next()
	}

	sort.Slice(matches, func(i, j int) bool {
		return matches[i].Score > matches[j].Score
	})

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
