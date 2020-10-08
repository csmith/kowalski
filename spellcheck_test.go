package kowalski

import (
	"bytes"
	"io/ioutil"
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

func TestSaveSpellChecker(t *testing.T) {
	buffer := &bytes.Buffer{}
	if err := SaveSpellChecker(buffer, testChecker); err != nil {
		t.Errorf("SaveSpellChecker() failed to save checker: %v", err)
	}

	saved, _ := ioutil.ReadFile("testdata/test_words.wl")
	if !reflect.DeepEqual(saved, buffer.Bytes()) {
		t.Errorf("Saved spell checker differs from golden: got %v, wanted %v", buffer.Bytes(), saved)
	}
}

func TestLoadSpellChecker(t *testing.T) {
	f, _ := os.Open("testdata/test_words.wl")
	s, err := LoadSpellChecker(f)

	if err != nil {
		t.Errorf("LoadSpellChecker() failed to load checker: %v", err)
	}

	if !reflect.DeepEqual(testChecker, s) {
		t.Errorf("Restored spell checker differs from original: got %v, wanted %v", s, testChecker)
	}
}
