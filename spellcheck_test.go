package kowalski

import (
	"bytes"
	"os"
	"reflect"
	"testing"
)

var testChecker *SpellChecker

func init() {
	f, _ := os.Open("testdata/test_words.txt")
	defer f.Close()
	testChecker, _ = CreateSpellChecker(f, 10)
}

func TestCreateSpellChecker(t *testing.T) {
	tests := []struct {
		word   string
		valid  bool
		prefix bool
	}{
		{"foo", true, true},
		{"bar", true, true},
		{"baz", true, true},
		{"ba", false, true},
		{"b", false, true},
		{"ab", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.word, func(t *testing.T) {
			if got := testChecker.Valid(tt.word); got != tt.valid {
				t.Errorf("Valid() = %v, want %v", got, tt.valid)
			}
			if got := testChecker.Prefix(tt.word); got != tt.prefix {
				t.Errorf("Prefix() = %v, want %v", got, tt.prefix)
			}
		})
	}
}

func TestSaveLoadSpellChecker(t *testing.T) {
	buffer := &bytes.Buffer{}
	if err := SaveSpellChecker(buffer, testChecker); err != nil {
		t.Errorf("SaveSpellChecker() failed to save checker: %v", err)
	}

	saved, err := LoadSpellChecker(buffer)
	if err != nil {
		t.Errorf("Failed to load saved checker: $%v", err)
	}

	if !reflect.DeepEqual(saved, testChecker) {
		t.Errorf("Saved spell checker differs when loaded: got %v, wanted %v", saved, testChecker)
	}
}
