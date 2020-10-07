package kowalski

import (
	"reflect"
	"testing"
)

func TestFromMorse(t *testing.T) {
	words, _ := LoadWords("testdata/test_words.txt")
	tests := []struct {
		name  string
		query string
		want  []string
	}{
		{"foo", "..-.------", []string{"foo"}},
		{"bar", "-....-.-.", []string{"bar"}},
		{"baz", "-... .- --..", []string{"baz"}},
		{"quux", "--.-..-..--..-", []string{"quux"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := words.FromMorse(tt.query); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromMorse() = %v, want %v", got, tt.want)
			}
		})
	}
}
