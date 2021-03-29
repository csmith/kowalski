package kowalski

import (
	"context"
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
func MultiplexMatch(ctx context.Context, checkers []*SpellChecker, pattern string, opts ... MultiplexOption) ([][]string, error) {
	return multiplexWithErrors(checkers, func(checker *SpellChecker) ([]string, error) {
		return Match(ctx, checker, pattern)
	}, opts)
}

// MultiplexMultiMatch performs the MultiMatch operation over a number of different checkers.
func MultiplexMultiMatch(ctx context.Context, checkers []*SpellChecker, pattern string, opts ... MultiplexOption) ([][]string, error) {
	return multiplexWithErrors(checkers, func(checker *SpellChecker) ([]string, error) {
		return MultiMatch(ctx, checker, pattern)
	}, opts)
}

// MultiplexAnagram performs the Anagram operation over a number of different checkers.
func MultiplexAnagram(ctx context.Context, checkers []*SpellChecker, pattern string, opts ... MultiplexOption) ([][]string, error) {
	return multiplexWithErrors(checkers, func(checker *SpellChecker) ([]string, error) {
		return Anagram(ctx, checker, pattern)
	}, opts)
}

// MultiplexMultiAnagram performs the MultiAnagram operation over a number of different checkers.
func MultiplexMultiAnagram(ctx context.Context, checkers []*SpellChecker, pattern string, opts ... MultiplexOption) ([][]string, error) {
	return multiplexWithErrors(checkers, func(checker *SpellChecker) ([]string, error) {
		return MultiAnagram(ctx, checker, pattern)
	}, opts)
}

// MultiplexFindWords performs the FindWords operation over a number of different checkers.
func MultiplexFindWords(checkers []*SpellChecker, pattern string, opts ... MultiplexOption) [][]string {
	return multiplex(checkers, func(checker *SpellChecker) []string {
		return FindWords(checker, pattern)
	}, opts)
}

// MultiplexFromMorse performs the FromMorse operation over a number of different checkers.
func MultiplexFromMorse(checkers []*SpellChecker, pattern string, opts ... MultiplexOption) [][]string {
	return multiplex(checkers, func(checker *SpellChecker) []string {
		return FromMorse(checker, pattern)
	}, opts)
}

// MultiplexOffByOne performs the OffByOne operation over a number of different checkers.
func MultiplexOffByOne(ctx context.Context, checkers []*SpellChecker, pattern string, opts ... MultiplexOption) ([][]string, error) {
	return multiplexWithErrors(checkers, func(checker *SpellChecker) ([]string, error) {
		return OffByOne(ctx, checker, pattern)
	}, opts)
}

// MultiplexFromT9 performs the FromT9 operation over a number of different checkers.
func MultiplexFromT9(checkers []*SpellChecker, pattern string, opts ... MultiplexOption) [][]string {
	return multiplex(checkers, func(checker *SpellChecker) []string {
		return FromT9(checker, pattern)
	}, opts)
}

// MultiplexWordSearch performs the WordSearch operation over a number of different checkers.
func MultiplexWordSearch(checkers []*SpellChecker, pattern []string, opts ... MultiplexOption) [][]string {
	return multiplex(checkers, func(checker *SpellChecker) []string {
		return WordSearch(checker, pattern)
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

func multiplexWithErrors(checkers []*SpellChecker, f func(checker *SpellChecker)([]string, error), opts []MultiplexOption) ([][]string, error) {
	o := &multiplexOptions{}
	for i := range opts {
		opts[i](o)
	}

	res := make([][]string, len(checkers))
	errs := make([]error, len(checkers))
	wg := &sync.WaitGroup{}

	for i := range checkers {
		wg.Add(1)
		go func(i int) {
			res[i], errs[i] = f(checkers[i])
			wg.Done()
		}(i)
	}

	wg.Wait()

	for i := range errs {
		if errs[i] != nil {
			return nil, errs[i]
		}
	}

	if o.dedupe {
		return dedupe(res), nil
	}
	return res, nil
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
