package kowalski

import (
	"reflect"
	"testing"
)

func TestFromT9(t *testing.T) {
	tests := []struct {
		name  string
		query string
		want  []string
	}{
		{"foo", "366", []string{"foo"}},
		{"bar", "227", []string{"bar"}},
		{"baz", "229", []string{"baz"}},
		{"quux", "7889", []string{"quux"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FromT9(testChecker, tt.query); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromT9() = %v, want %v", got, tt.want)
			}
		})
	}
}
