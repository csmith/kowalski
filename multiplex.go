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
func MultiplexMatch(ctx context.Context, checkers []*SpellChecker, pattern string, opts ...MultiplexOption) ([][]string, error) {
	return multiplexWords(
		checkers,
		func(checker *SpellChecker) ([]string, error) {
			return Match(ctx, checker, pattern)
		},
		opts,
	)
}

// MultiplexMultiMatch performs the MultiMatch operation over a number of different checkers.
func MultiplexMultiMatch(ctx context.Context, checkers []*SpellChecker, pattern string, opts ...MultiplexOption) ([][]string, error) {
	return multiplexWords(
		checkers,
		func(checker *SpellChecker) ([]string, error) {
			return MultiMatch(ctx, checker, pattern)
		},
		opts,
	)
}

// MultiplexAnagram performs the Anagram operation over a number of different checkers.
func MultiplexAnagram(ctx context.Context, checkers []*SpellChecker, pattern string, opts ...MultiplexOption) ([][]string, error) {
	return multiplexWords(
		checkers,
		func(checker *SpellChecker) ([]string, error) {
			return Anagram(ctx, checker, pattern)
		},
		opts,
	)
}

// MultiplexMultiAnagram performs the MultiAnagram operation over a number of different checkers.
func MultiplexMultiAnagram(ctx context.Context, checkers []*SpellChecker, pattern string, opts ...MultiplexOption) ([][]string, error) {
	return multiplexWords(
		checkers,
		func(checker *SpellChecker) ([]string, error) {
			return MultiAnagram(ctx, checker, pattern)
		},
		opts,
	)
}

// MultiplexFindWords performs the FindWords operation over a number of different checkers.
func MultiplexFindWords(checkers []*SpellChecker, pattern string, opts ...MultiplexOption) [][]string {
	res, _ := multiplexWords(
		checkers,
		func(checker *SpellChecker) ([]string, error) {
			return FindWords(checker, pattern), nil
		},
		opts,
	)
	return res
}

// MultiplexFromMorse performs the FromMorse operation over a number of different checkers.
func MultiplexFromMorse(checkers []*SpellChecker, pattern string, opts ...MultiplexOption) [][]string {
	res, _ := multiplexWords(
		checkers,
		func(checker *SpellChecker) ([]string, error) {
			return FromMorse(checker, pattern), nil
		},
		opts,
	)
	return res
}

// MultiplexOffByOne performs the OffByOne operation over a number of different checkers.
func MultiplexOffByOne(ctx context.Context, checkers []*SpellChecker, pattern string, opts ...MultiplexOption) ([][]string, error) {
	return multiplexWords(
		checkers,
		func(checker *SpellChecker) ([]string, error) {
			return OffByOne(ctx, checker, pattern)
		},
		opts,
	)
}

// MultiplexFromT9 performs the FromT9 operation over a number of different checkers.
func MultiplexFromT9(checkers []*SpellChecker, pattern string, opts ...MultiplexOption) [][]string {
	res, _ := multiplexWords(
		checkers,
		func(checker *SpellChecker) ([]string, error) {
			return FromT9(checker, pattern), nil
		},
		opts,
	)
	return res
}

// MultiplexWordSearch performs the WordSearch operation over a number of different checkers.
func MultiplexWordSearch(checkers []*SpellChecker, pattern []string, opts ...MultiplexOption) [][]string {
	res, _ := multiplexWords(
		checkers,
		func(checker *SpellChecker) ([]string, error) {
			return WordSearch(checker, pattern), nil
		},
		opts,
	)
	return res
}

// MultiplexCheckWords performs the CheckWords operation over a number of different checkers.
// Returns results for each checker separately.
func MultiplexCheckWords(checkers []*SpellChecker, input string) [][][]WordCheckResult {
	res, _ := multiplex(
		checkers,
		func(checker *SpellChecker) ([][]WordCheckResult, error) {
			return CheckWords(checker, input), nil
		},
	)
	return res
}

func multiplexWords(
	checkers []*SpellChecker,
	f func(*SpellChecker) ([]string, error),
	opts []MultiplexOption,
) ([][]string, error) {
	o := &multiplexOptions{}
	for i := range opts {
		opts[i](o)
	}

	res, err := multiplex(checkers, f)
	if err != nil {
		return nil, err
	}

	if o.dedupe {
		return dedupe(res), nil
	}
	return res, nil
}

// multiplex runs f against each checker concurrently, returning results in checker order. If any call fails, the
// error from the earliest checker is returned and the results are discarded.
func multiplex[T any](checkers []*SpellChecker, f func(*SpellChecker) (T, error)) ([]T, error) {
	res := make([]T, len(checkers))
	errs := make([]error, len(checkers))
	wg := &sync.WaitGroup{}

	for i := range checkers {
		wg.Go(
			func() {
				res[i], errs[i] = f(checkers[i])
			},
		)
	}

	wg.Wait()

	for i := range errs {
		if errs[i] != nil {
			return nil, errs[i]
		}
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
