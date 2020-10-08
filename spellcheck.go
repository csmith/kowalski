package kowalski

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"github.com/willf/bloom"
	"io"
	"strings"
)

// SpellChecker provides a way to tell whether a word exists in a dictionary.
type SpellChecker struct {
	primary     *bloom.BloomFilter
	secondaries [2]*bloom.BloomFilter
	roots       *bloom.BloomFilter
}

// LoadSpellChecker attempts to load a SpellChecker that was previously saved with SaveSpellChecker.
func LoadSpellChecker(reader io.Reader) (*SpellChecker, error) {
	var filters []*bloom.BloomFilter

	if err := gob.NewDecoder(reader).Decode(&filters); err != nil {
		return nil, err
	}

	if len(filters) != 4 {
		return nil, fmt.Errorf("saved spell checker contains %d filters, expected 4", len(filters))
	}

	return &SpellChecker{
		primary: filters[0],
		secondaries: [2]*bloom.BloomFilter{
			filters[1],
			filters[2],
		},
		roots: filters[3],
	}, nil
}

// SaveSpellChecker serialises the given checker and writes it to the writer.
// It can later be restored with LoadSpellChecker.
func SaveSpellChecker(writer io.Writer, checker *SpellChecker) error {
	filters := []*bloom.BloomFilter{
		checker.primary,
		checker.secondaries[0],
		checker.secondaries[1],
		checker.roots,
	}

	return gob.NewEncoder(writer).Encode(filters)
}

// CreateSpellChecker creates a new SpellChecker by reading words line-by-line from the given reader.
// The wordCount parameter should be an approximation of the number of words available.
//
// This is likely to be a relatively expensive operation; for routine use prefer saving the spell
// checker via SaveSpellChecker and restoring it with LoadSpellChecker.
func CreateSpellChecker(reader io.Reader, wordCount int) (*SpellChecker, error) {
	c := &SpellChecker{
		primary: bloom.NewWithEstimates(uint(wordCount), 0.001),
		secondaries: [2]*bloom.BloomFilter{
			bloom.NewWithEstimates(uint(wordCount/2), 0.001),
			bloom.NewWithEstimates(uint(wordCount/2), 0.001),
		},
		roots: bloom.NewWithEstimates(uint(wordCount*10), 0.1),
	}

	scanner := bufio.NewScanner(reader)
	counter := 0
	for scanner.Scan() {
		line := strings.ToLower(scanner.Text())
		c.addWord(line, counter)
		counter = 1 - counter
	}

	return c, scanner.Err()
}

// addWord adds a new word to the spell checker - that is, it adds it to the primary bloom filter and one of the
// two secondaries, and also adds all prefixes of the word to the roots filter.
func (c *SpellChecker) addWord(word string, secondary int) {
	c.primary.AddString(word)
	c.secondaries[secondary].AddString(word)

	for i := range word {
		c.roots.AddString(word[0 : i+1])
	}
}

// Valid determines - probabilistically - whether the given word was in the word list used to create this SpellChecker.
// There is a small chance of false positives, i.e. a word that wasn't in the word list might be incorrectly identified
// as valid; there is no chance of false negatives.
func (c *SpellChecker) Valid(word string) bool {
	if c.primary.TestString(word) == false {
		return false
	}

	for i := range c.secondaries {
		if c.secondaries[i].TestString(word) {
			return true
		}
	}

	return false
}

// Prefix determines - probabilistically - whether the given string is a prefix of any known word.
// There is a small chance of false positives, i.e. an input that is not a prefix to any word in the word list might be
// incorrectly identified as valid; there is no chance of false negatives.
func (c *SpellChecker) Prefix(prefix string) bool {
	return c.roots.TestString(prefix)
}
