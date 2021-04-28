package main

import (
	"flag"
	"fmt"
	"sort"
	"strings"

	"github.com/blevesearch/vellum"
	"github.com/blevesearch/vellum/regexp"
)

var fstModel = flag.String("fst-model", "", "Path to FST for fast word operations")

var fst *vellum.FST

func init() {
	flag.Parse()

	if *fstModel != "" {
		var err error

		fst, err = vellum.Open(*fstModel)
		if err != nil {
			panic(err)
		}
	}
}

func FstAnagram(input string, r Replier) {
	autonoma, err := NewAnagramAutomoma(input)
	if err != nil {
		r.reply("Error: %s", err.Error())
		return
	}

	iterator, err := fst.Search(autonoma, nil, nil)
	if err != nil {
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
	message.WriteString("Anagrams for '")
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

func init() {
	if *fstModel != "" {
		addTextCommand(FstAnagram, "Attempts to find anagrams from wikipedia, expanding '\\*' wildcards", "fstanagram", "fstagram")
	}
}

func FstRegex(input string, r Replier) {
	autonoma, err := regexp.New(input)
	if err != nil {
		r.reply("Error: %s", err.Error())
		return
	}

	iterator, err := fst.Search(autonoma, nil, nil)
	if err != nil {
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

func init() {
	if *fstModel != "" {
		addTextCommand(FstRegex, "Attempts to find word matches from wikipedia using regexp", "fstregex", "fstre")
	}
}

const errorMask = 1 << 62

type anagramAutonoma struct {
	Chars string
}

func NewAnagramAutomoma(term string) (vellum.Automaton, error) {
	// Ensure all wildcards are at the end of the terms
	term = strings.ToLower(term)
	wildcards := strings.Count(term, "*")
	rearranged := strings.ReplaceAll(term, "*", "") + strings.Repeat("*", wildcards)
	if len(term) >= 62 {
		return nil, fmt.Errorf("term is too long: %d (max length is 61)", len(term))
	}

	return &anagramAutonoma{
		Chars: rearranged,
	}, nil
}

func (a anagramAutonoma) Start() int {
	r := 0
	for i := range a.Chars {
		r |= 1 << i
	}
	return r
}

func (a anagramAutonoma) IsMatch(i int) bool {
	return i == 0
}

func (a anagramAutonoma) CanMatch(i int) bool {
	return (i & errorMask) == 0
}

func (a anagramAutonoma) WillAlwaysMatch(i int) bool {
	return false
}

func (a anagramAutonoma) Accept(i int, b byte) int {
	if b == ' ' {
		// Skip over spaces
		return i
	}

	for j := range a.Chars {
		if a.Chars[j] == b || a.Chars[j] == '*' {
			mask := 1 << j
			if (mask & i) == mask {
				// Not yet used
				return i ^ mask
			}
		}
	}

	return i | errorMask
}
