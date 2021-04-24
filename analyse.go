package kowalski

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/csmith/kowalski/v4/data"
)

var nonLetterRegex = regexp.MustCompile("[^a-z]+")

type analyser func(checker *SpellChecker, input string) []string

func analyseEntropy(_ *SpellChecker, input string) []string {
	var results []string

	entropy := shannonEntropy(input)
	if entropy <= 0.5 {
		results = append(results, fmt.Sprintf("Shannon entropy is %.2f - very little variation in input", entropy))
	} else if entropy >= 3.5 && entropy <= 5 {
		results = append(results, fmt.Sprintf("Shannon entropy is %.2f - typical of English text", entropy))
	} else if entropy >= 7.5 {
		results = append(results, fmt.Sprintf("Shannon entropy is %.2f - very high, likely encrypted/compressed", entropy))
	}

	return results
}

func analyseDataReferences(_ *SpellChecker, input string) []string {
	var results []string

	cleaned := nonLetterRegex.ReplaceAllString(strings.ToLower(input), "")
	if len(cleaned) > 0 {
		for name := range data.Index {
			if terms, ok := splitTerms(cleaned, nil, data.Index[name]); ok {
				if sameLength(data.Index[name]) {
					results = append(results, fmt.Sprintf("Consists entirely of %s", name))
				} else {
					results = append(results, fmt.Sprintf("Consists entirely of %s: %s", name, strings.Join(terms, " ")))
				}
			}
		}
	}

	return results
}

func analyseCaesarShifts(checker *SpellChecker, input string) []string {
	var results []string

	shifts := CaesarShifts(input)
	bestScore, bestShift := 0.0, 0
	for i, s := range shifts {
		score := Score(checker, s)
		if score > bestScore {
			bestScore = score
			bestShift = i
		}
	}
	if bestScore > 0.5 {
		results = append(results, fmt.Sprintf("Caesar shift of %d might be English: %s", bestShift + 1, shifts[bestShift]))
	}

	return results
}

func analyseAlternateChars(checker *SpellChecker, input string) []string {
	var results []string

	odds := strings.Builder{}
	evens := strings.Builder{}
	for i := range input {
		if i%2 == 0 {
			evens.WriteByte(input[i])
		} else {
			odds.WriteByte(input[i])
		}
	}

	if Score(checker, odds.String()) > 0.5 {
		results = append(results, fmt.Sprintf("Alternating characters might be English: %s", odds.String()))
	}

	if Score(checker, evens.String()) > 0.5 {
		results = append(results, fmt.Sprintf("Alternating characters might be English: %s", evens.String()))
	}

	return results
}

func analyseLength(_ *SpellChecker, input string) []string {
	var results []string

	cleaned := nonLetterRegex.ReplaceAllString(strings.ToLower(input), "")
	if len(input)%8 == 0 {
		results = append(results, "Multiple of 8 characters - might be encoded binary?")
	} else if len(cleaned)%8 == 0 {
		results = append(results, "Multiple of 8 A-Z characters - might be encoded binary?")
	}

	results = append(results, fmt.Sprintf("%d characters long (total)", len(input)))
	results = append(results, fmt.Sprintf("%d characters long (a-zA-Z)", len(cleaned)))

	return results
}

func analyseDistribution(_ *SpellChecker, input string) []string {
	var results []string

	dists := LetterDistribution(input)
	present := 0
	for i := range dists {
		if dists[i] > 0 {
			present++
		}
	}

	if present > 20 {
		message := strings.Builder{}
		message.WriteString("Contains all english letters")
		if present < 26 {
			message.WriteString(" except for: ")
			for i := range dists {
				if dists[i] == 0 {
					message.WriteByte(byte('A' + i))
				}
			}
		}
		results = append(results, message.String())
	}

	if present > 0 && present < 10 && present < len(input) {
		chars := strings.Builder{}
		for i := range dists {
			if dists[i] > 0 {
				chars.WriteByte(byte('A' + i))
			}
		}

		results = append(results, fmt.Sprintf("Contains only some letters: %s", chars.String()))
		if chars.String() == "ADFGX" {
			results = append(results, "Might be an ADFGX cipher?")
		} else if chars.String() == "ADFGVX" {
			results = append(results, "Might be an ADFGVX cipher?")
		}
	}

	return results
}

var rleRegex = regexp.MustCompile(`^(\d+\D)+$`)

func analyseRunLengthEncoding(_ *SpellChecker, input string) []string {
	var results []string

	if rleRegex.MatchString(input) {
		message := strings.Builder{}
		message.WriteString("Might be run-length encoded: ")
		num := 0
		for i := range input {
			if d, err := strconv.Atoi(string(input[i])); err == nil {
				num = 10 * num + d
			} else {
				message.WriteString(strings.Repeat(string(input[i]), num))
				num = 0
			}
		}
		if message.Len() > 250 {
			results = append(results, fmt.Sprintf("%s...", message.String()[0:247]))
		} else {
			results = append(results, message.String())
		}
	}

	return results
}

var analysers = []analyser{
	analyseEntropy,
	analyseDataReferences,
	analyseCaesarShifts,
	analyseAlternateChars,
	analyseLength,
	analyseDistribution,
	analyseRunLengthEncoding,
}

// Analyse performs various forms of text analysis on the input and returns findings.
func Analyse(checker *SpellChecker, input string) []string {
	var results []string

	for i := range analysers {
		results = append(results, analysers[i](checker, input)...)
	}

	return results
}

// Score assigns a score to an input showing how likely it is to be English text. A score of 1.0 means almost
// certainly English, a score of 0.0 means almost certainly not. This is fairly arbitrary and is not very good.
func Score(checker *SpellChecker, input string) float64 {
	const targetDensity = 2.0
	density := float64(len(FindWords(checker, input))) / float64(len(input))
	densityScore := math.Max(1-math.Abs(density-targetDensity), 0.1)

	entropy := shannonEntropy(input)
	entropyScore := 1.0
	if entropy < 3.5 {
		entropyScore = math.Max(entropy/3.5, 0.1)
	} else if entropy > 5 {
		entropyScore = math.Max(1-(entropy-5)/3, 0.1)
	}

	return densityScore * entropyScore
}

// splitTerms splits the input up into a list of the given terms. The input is expected
// to be lowercase, and with any irrelevant characters removed.
func splitTerms(input string, prefix, terms []string) ([]string, bool) {
	for i := range terms {
		if strings.HasPrefix(input, strings.ToLower(terms[i])) {
			newPrefix := append(prefix, terms[i])
			if len(input) == len(terms[i]) {
				return newPrefix, true
			} else if res, ok := splitTerms(input[len(terms[i]):], newPrefix, terms); ok {
				return res, true
			}
		}
	}
	return nil, false
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

func sameLength(terms []string) bool {
	if len(terms) == 0 {
		return true
	}

	target := len(terms[0])
	for i := range terms {
		if len(terms[i]) != target {
			return false
		}
	}
	return true
}

// LetterDistribution counts the number of the occurrences of each English letter (ignoring case).
func LetterDistribution(input string) [26]int {
	var res [26]int
	for i := range input {
		if input[i] >= 'a' && input[i] <= 'z' {
			res[input[i]-'a']++
		}
		if input[i] >= 'A' && input[i] <= 'Z' {
			res[input[i]-'A']++
		}
	}
	return res
}
