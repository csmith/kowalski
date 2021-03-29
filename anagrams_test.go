package kowalski

import (
	"context"
	"reflect"
	"testing"
)

func TestAnagrams(t *testing.T) {
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
			if got, _ := Anagram(context.Background(), testChecker, tt.query); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Anagram() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPermute(t *testing.T) {
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
