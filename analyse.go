package kowalski

import (
	"fmt"
	"github.com/csmith/kowalski/v2/data"
	"math"
	"regexp"
	"strings"
)

var nonLetterRegex = regexp.MustCompile("[^a-z]+")

// Analyse performs various forms of text analysis on the input and returns findings.
func Analyse(checker *SpellChecker, input string) []string {
	var results []string

	entropy := shannonEntropy(input)
	if entropy <= 0.5 {
		results = append(results, fmt.Sprintf("Shannon entropy is %.2f - very little variation in input", entropy))
	} else if entropy >= 3.5 && entropy <= 5 {
		results = append(results, fmt.Sprintf("Shannon entropy is %.2f - typical of English text", entropy))
	} else if entropy >= 7.5 {
		results = append(results, fmt.Sprintf("Shannon entropy is %.2f - very high, likely encrypted/compressed", entropy))
	}

	cleaned := nonLetterRegex.ReplaceAllString(strings.ToLower(input), "")
	if len(cleaned) > 0 {
		for name := range data.Index {
			if consistsOf(cleaned, data.Index[name]) {
				results = append(results, fmt.Sprintf("Consists entirely of %s", name))
			}
		}
	}

	shifts := CaesarShifts(input)
	for i, s := range shifts {
		if Score(checker, s) > 0.5 {
			results = append(results, fmt.Sprintf("Caesar shift of %d might be English: %s", i, s))
		}
	}

	if len(input) % 8 == 0 {
		results = append(results, "Multiple of 8 characters - might be encoded binary?")
	} else if len(cleaned) % 8 == 0 {
		results = append(results, "Multiple of 8 A-Z characters - might be encoded binary?")
	}

	return results
}

// Score assigns a score to an input showing how likely it is to be English text. A score of 1.0 means almost
// certainly English, a score of 0.0 means almost certainly not. This is fairly arbitrary and may not be very good.
func Score(checker *SpellChecker, input string) float64 {
	const targetDensity = 2.0
	density := float64(len(FindWords(checker, input))) / float64(len(input))
	densityScore := math.Max(1 - math.Abs(density - targetDensity), 0.1)

	entropy := shannonEntropy(input)
	entropyScore := 1.0
	if entropy < 3.5 {
		entropyScore = math.Max(entropy / 3.5, 0.1)
	} else if entropy > 5 {
		entropyScore = math.Max(1 - (entropy - 5) / 3, 0.1)
	}

	return densityScore * entropyScore
}

// consistsOf determines if the given input consists _entirely_ of terms in the given slice. The input is expected
// to be lowercase, and with any irrelevant characters removed.
func consistsOf(input string, terms []string) bool {
	for i := range terms {
		if strings.HasPrefix(input, strings.ToLower(terms[i])) {
			if len(input) == len(terms[i]) || consistsOf(input[len(terms[i]):], terms) {
				return true
			}
		}
	}
	return false
}

// shannonEntropy calculates the Shannon Entropy of the input.
func shannonEntropy(input string) float64 {
	var occurrences [256]float64
	for i := range input {
		occurrences[input[i]]++
	}

	var size = float64(len(input))
	var entropy float64 = 0
	for i := range occurrences {
		if occurrences[i] > 0 {
			prob := occurrences[i] / size
			entropy -= prob * math.Log2(prob)
		}
	}
	return entropy
}
