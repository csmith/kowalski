package kowalski

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChunk(t *testing.T) {
	tests := []struct {
		name  string
		input string
		parts []int
		want  []string
	}{
		{"repeats a single size", "abcdef", []int{2}, []string{"ab", "cd", "ef"}},
		{"cycles through sizes", "abcdefgh", []int{1, 2}, []string{"a", "bc", "d", "ef", "g", "h"}},
		{"keeps a short final chunk", "abcde", []int{3}, []string{"abc", "de"}},
		{"size larger than the input", "abc", []int{10}, []string{"abc"}},
		{"empty input", "", []int{2}, nil},
		{"ignores a zero size", "abcdef", []int{3, 0}, []string{"abc", "def"}},
		{"ignores a negative size", "abcdef", []int{-1, 2}, []string{"ab", "cd", "ef"}},
		{"only unusable sizes", "abc", []int{0, -1}, []string{"abc"}},
		{"no sizes", "abc", nil, []string{"abc"}},
		{"no usable sizes or input", "", []int{0}, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, Chunk(tt.input, tt.parts...))
		})
	}
}

func TestChunk_doesNotModifyParts(t *testing.T) {
	parts := []int{0, 2}
	Chunk("abcd", parts...)
	assert.Equal(t, []int{0, 2}, parts)
}
