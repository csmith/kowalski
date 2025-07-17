package main

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"sort"
	"strings"

	"github.com/blevesearch/vellum"
	vellumRegexp "github.com/blevesearch/vellum/regexp"
	"github.com/csmith/kowalski/v6/fst"
)

var fstTransducer *vellum.FST

func initFST(path string) {
	var err error
	fstTransducer, err = vellum.Open(path)
	if err != nil {
		log.Printf("Failed to open FST model: %v", err)
	}
}

type fstMatch struct {
	Term  string `json:"term"`
	Score uint64 `json:"score"`
}

func processFstAnagram(input string) (interface{}, error) {
	automaton, err := fst.NewAnagramAutomaton(input)
	if err != nil {
		return nil, err
	}

	matches, err := fstQuery(automaton)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"input":   input,
		"matches": matches,
	}, nil
}

func processFstRegex(input string) (interface{}, error) {
	automaton, err := vellumRegexp.New(input)
	if err != nil {
		return nil, err
	}

	matches, err := fstQuery(automaton)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"input":   input,
		"matches": matches,
	}, nil
}

func processFstMorse(input string) (interface{}, error) {
	automaton := fst.NewMorseAutomaton(input)
	matches, err := fstQuery(automaton)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"input":   input,
		"matches": matches,
	}, nil
}

func processWordLink(input string) (interface{}, error) {
	input = strings.ToLower(input)
	parts := strings.Split(input, " ")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid input, must specify two words to link")
	}

	leftAut, err := vellumRegexp.New(fmt.Sprintf("%s[ -]?([a-zA-Z]{3,})", regexp.QuoteMeta(parts[0])))
	if err != nil {
		return nil, err
	}
	leftSide, err := fstResults(leftAut)
	if err != nil {
		return nil, err
	} else if len(leftSide) == 0 {
		return nil, fmt.Errorf("no words found that can follow '%s'", parts[0])
	}

	rightAut, err := vellumRegexp.New(fmt.Sprintf("([a-zA-Z]{3,})[ -]?%s", regexp.QuoteMeta(parts[1])))
	if err != nil {
		return nil, err
	}

	rightSide, err := fstResults(rightAut)
	if err != nil {
		return nil, err
	} else if len(rightSide) == 0 {
		return nil, fmt.Errorf("no words found that can precede '%s'", parts[1])
	}

	var leftScores = make(map[string]uint64)
	for i := range leftSide {
		leftTerm := strings.TrimSpace(strings.TrimPrefix(leftSide[i].Term, parts[0]))
		if normal, _, _ := fstTransducer.Get([]byte(leftTerm)); normal < 5_000_000 {
			leftScores[leftTerm] += leftSide[i].Score
		}
	}

	var rightScores = make(map[string]uint64)
	for i := range rightSide {
		rightTerm := strings.TrimSpace(strings.TrimSuffix(rightSide[i].Term, parts[1]))
		if normal, _, _ := fstTransducer.Get([]byte(rightTerm)); normal < 5_000_000 {
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
		return nil, fmt.Errorf("no linking words found")
	}

	sort.Slice(combined, func(i, j int) bool {
		return combined[i].Score > combined[j].Score
	})

	return map[string]interface{}{
		"input": input,
		"words": parts,
		"links": combined,
	}, nil
}

func fstQuery(automaton vellum.Automaton) ([]fstMatch, error) {
	return fstResults(automaton)
}

func fstResults(automaton vellum.Automaton) ([]fstMatch, error) {
	iterator, err := fstTransducer.Search(automaton, nil, nil)
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
