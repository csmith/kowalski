package kowalski

import (
	"reflect"
	"testing"
)

func TestFindWords(t *testing.T) {
	tests := []struct {
		name  string
		query string
		want  []string
	}{
		{"single word", "foo", []string{"foo"}},
		{"repeated words", "foofoofoo", []string{"foo", "foo", "foo"}},
		{"offset", "xyzfooxyz", []string{"foo"}},
		{"no words", "xyz", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FindWords(testChecker, tt.query); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindWords() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWordSearchLines(t *testing.T) {
	tests := []struct {
		name  string
		query []string
		want  []string
	}{
		{
			"square",
			[]string{"123", "456", "789"},
			[]string{
				// Horizontal
				"123", "456", "789",
				// Vertical
				"147", "258", "369",
				// Diagonals
				"42", "753", "86",
				"48", "159", "26",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expected := make(map[string]int)
			for i := range tt.want {
				expected[tt.want[i]]++
				expected[reverseString(tt.want[i])]++
			}

			actual := wordSearchLines(tt.query)
			for i := range actual {
				expected[actual[i]]--
			}

			for k, v := range expected {
				if v > 0 {
					t.Errorf("FindWords() = %v, want %v, found %d too few instances of %s", actual, tt.want, v, k)
				} else if v < 0 {
					t.Errorf("FindWords() = %v, want %v, found %d too many instances of %s", actual, tt.want, -v, k)
				}
			}
		})
	}
}
