package kowalski

import (
	"reflect"
	"testing"
)

func TestLoadWords(t *testing.T) {
	words, err := LoadWords("testdata/test_words.txt")
	if err != nil {
		t.Errorf("LoadWords() failed to load test file: %v", err)
		return
	}

	if !words.Valid("foo") {
		t.Errorf("'foo' element should be valid")
		return
	}

	if words.Valid("bazinga") {
		t.Errorf("'bazinga' element should not be valid")
		return
	}
}

func TestNode_Match(t *testing.T) {
	words, _ := LoadWords("testdata/test_words.txt")
	tests := []struct {
		name  string
		query string
		want  []string
	}{
		{"exact match", "foo", []string{"foo"}},
		{"exact match with one wildcard", "fo?", []string{"foo"}},
		{"exact match with two wildcards", "f??", []string{"foo"}},
		{"exact match with all wildcards", "????", []string{"quux"}},
		{"multiple matches with one wildcard", "ba?", []string{"bar", "baz"}},
		{"multiple matches with all wildcards", "???", []string{"bar", "baz", "foo"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := words.Match(tt.query); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNode_Anagrams(t *testing.T) {
	words, _ := LoadWords("testdata/test_words.txt")
	tests := []struct {
		name  string
		query string
		want  []string
	}{
		{"exact match", "foo", []string{"foo"}},
		{"anagram match", "oof", []string{"foo"}},
		{"exact match with one wildcard", "fo?", []string{"foo"}},
		{"anagram match with one wildcard", "oo?", []string{"foo"}},
		{"exact match with two wildcards", "f??", []string{"foo"}},
		{"anagram match with two wildcards", "?f?", []string{"foo"}},
		{"exact match with all wildcards", "????", []string{"quux"}},
		{"multiple exact matches with one wildcard", "ba?", []string{"bar", "baz"}},
		{"multiple anagram matches with one wildcard", "b?a", []string{"bar", "baz"}},
		{"multiple exact matches with all wildcards", "???", []string{"bar", "baz", "foo"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := words.Anagrams(tt.query); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Anagrams() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_permute(t *testing.T) {
	tests := []struct {
		name       string
		input      []byte
		swapBefore int
		want       []byte
	}{
		{"Simple permutation 1", []byte("abc"), 3, []byte("acb")},
		{"Simple permutation 2", []byte("acb"), 3, []byte("bac")},
		{"Simple permutation 3", []byte("bac"), 3, []byte("bca")},
		{"Simple permutation 4", []byte("bca"), 3, []byte("cab")},
		{"Simple permutation 5", []byte("cab"), 3, []byte("cba")},
		{"Simple permutation 6", []byte("cba"), 3, nil},
		{"SwapBefore 1", []byte("abc"), 1, []byte("bac")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := permute(tt.input, tt.swapBefore); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("permute() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}
