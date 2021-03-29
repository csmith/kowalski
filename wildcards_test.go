package kowalski

import (
	"context"
	"reflect"
	"testing"
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
