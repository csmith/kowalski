package kowalski

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMultiplex_ReturnsResultsInCheckerOrder(t *testing.T) {
	checkers := []*SpellChecker{{}, {}, {}}
	index := func(checker *SpellChecker) int {
		for i := range checkers {
			if checkers[i] == checker {
				return i
			}
		}
		return -1
	}

	res, err := multiplex(
		checkers,
		func(checker *SpellChecker) (int, error) {
			return index(checker), nil
		},
	)

	require.NoError(t, err)
	assert.Equal(t, []int{0, 1, 2}, res)
}

func TestMultiplex_ReturnsEarliestError(t *testing.T) {
	checkers := []*SpellChecker{{}, {}, {}}
	errs := map[*SpellChecker]error{
		checkers[1]: errors.New("second"),
		checkers[2]: errors.New("third"),
	}

	res, err := multiplex(
		checkers,
		func(checker *SpellChecker) (string, error) {
			return "ok", errs[checker]
		},
	)

	assert.Nil(t, res)
	assert.EqualError(t, err, "second")
}

func TestMultiplexWords_Dedupe(t *testing.T) {
	checkers := []*SpellChecker{{}, {}, {}}
	words := map[*SpellChecker][]string{
		checkers[0]: {"a", "b", "c"},
		checkers[1]: {"b", "d"},
		checkers[2]: {"a", "d", "e"},
	}
	f := func(checker *SpellChecker) ([]string, error) {
		return words[checker], nil
	}

	res, err := multiplexWords(checkers, f, nil)
	require.NoError(t, err)
	assert.Equal(t, [][]string{{"a", "b", "c"}, {"b", "d"}, {"a", "d", "e"}}, res)

	res, err = multiplexWords(checkers, f, []MultiplexOption{Dedupe})
	require.NoError(t, err)
	assert.Equal(t, [][]string{{"a", "b", "c"}, {"d"}, {"e"}}, res)
}
