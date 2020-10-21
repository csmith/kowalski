package kowalski

import (
	"sync"
)

type MultiplexOption func(*multiplexOptions)

type multiplexOptions struct {
	dedupe bool
}

// Dedupe removes duplicate entries from multiplexed results. That is, if the first checker provides words A, B and C,
// the second checker provides B and D, and the third A, D, and E, then the result will be: {A,B,C},{D},{E}.
func Dedupe(options *multiplexOptions) {
	options.dedupe = true
}

// MultiplexMatch performs the Match operation over a number of different checkers.
func MultiplexMatch(checkers []*SpellChecker, pattern string, opts ... MultiplexOption) [][]string {
	return multiplex(checkers, func(checker *SpellChecker) []string {
		return Match(checker, pattern)
	}, opts)
}

// MultiplexAnagram performs the Anagram operation over a number of different checkers.
func MultiplexAnagram(checkers []*SpellChecker, pattern string, opts ... MultiplexOption) [][]string {
	return multiplex(checkers, func(checker *SpellChecker) []string {
		return Anagram(checker, pattern)
	}, opts)
}

// MultiplexFromMorse performs the FromMorse operation over a number of different checkers.
func MultiplexFromMorse(checkers []*SpellChecker, pattern string, opts ... MultiplexOption) [][]string {
	return multiplex(checkers, func(checker *SpellChecker) []string {
		return FromMorse(checker, pattern)
	}, opts)
}

// MultiplexFromT9 performs the FromT9 operation over a number of different checkers.
func MultiplexFromT9(checkers []*SpellChecker, pattern string, opts ... MultiplexOption) [][]string {
	return multiplex(checkers, func(checker *SpellChecker) []string {
		return FromT9(checker, pattern)
	}, opts)
}

func multiplex(checkers []*SpellChecker, f func(checker *SpellChecker)[]string, opts []MultiplexOption) [][]string {
	o := &multiplexOptions{}
	for i := range opts {
		opts[i](o)
	}

	res := make([][]string, len(checkers))
	wg := &sync.WaitGroup{}

	for i := range checkers {
		wg.Add(1)
		go func(i int) {
			res[i] = f(checkers[i])
			wg.Done()
		}(i)
	}

	wg.Wait()

	if o.dedupe {
		return dedupe(res)
	}
	return res
}

func dedupe(data [][]string) [][]string {
	res := make([][]string, len(data))

	existing := make(map[string]bool)
	for i := range data {
		words := data[i]
		for j := range words {
			if !existing[words[j]] {
				res[i] = append(res[i], words[j])
				existing[words[j]] = true
			}
		}
	}

	return res
}
