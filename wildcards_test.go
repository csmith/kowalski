package kowalski

import (
	"context"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMatch(t *testing.T) {
	tests := []struct {
		name  string
		query string
		want  []string
	}{
		{"no match", "fr?", nil},
		{"exact match", "foo", []string{"foo"}},
		{"exact match with one wildcard", "fo?", []string{"foo"}},
		{"exact match with two wildcards", "f??", []string{"foo"}},
		{"exact match with all wildcards", "????", []string{"quux"}},
		{"multiple matches with one wildcard", "ba?", []string{"bar", "baz"}},
		{"multiple matches with all wildcards", "???", []string{"bar", "baz", "foo"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := Match(context.Background(), testChecker, tt.query); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMultiMatch(t *testing.T) {
	tests := []struct {
		name  string
		query string
		want  []string
	}{
		// All the normal things should still work...
		{"no match", "fr?", nil},
		{"exact match", "foo", []string{"foo"}},
		{"exact match with one wildcard", "fo?", []string{"foo"}},
		{"exact match with two wildcards", "f??", []string{"foo"}},
		{"exact match with all wildcards", "????", []string{"quux"}},
		{"multiple matches with one wildcard", "ba?", []string{"bar", "baz"}},
		{"multiple matches with all wildcards", "???", []string{"bar", "baz", "foo"}},
		// And our special new ones...
		{"two words", "foobar", []string{"foo bar"}},
		{"two words with wildcards", "foob??", []string{"foo bar", "foo baz"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := MultiMatch(context.Background(), testChecker, tt.query); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMultiMatch_tooManyCandidates(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := MultiMatch(ctx, allTwoLetterWordsChecker(t), "??????")
	assert.ErrorIs(t, err, ErrTooManyCandidates)
}

func allTwoLetterWordsChecker(t *testing.T) *SpellChecker {
	var words strings.Builder
	for a := 'a'; a <= 'z'; a++ {
		for b := 'a'; b <= 'z'; b++ {
			words.WriteString(string([]rune{a, b, '\n'}))
		}
	}

	checker, err := CreateSpellChecker(strings.NewReader(words.String()), 26*26)
	require.NoError(t, err)
	return checker
}
