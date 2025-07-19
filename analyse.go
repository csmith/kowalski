package kowalski

import (
	"fmt"
	"math"
	"math/big"
	"regexp"
	"strconv"
	"strings"

	"github.com/csmith/cryptography"
	"github.com/csmith/kowalski/v6/data"
)

var nonLetterRegex = regexp.MustCompile("[^a-z]+")

type analyser func(checker *SpellChecker, input string) []string

func analyseEntropy(_ *SpellChecker, input string) []string {
	var results []string

	entropy := cryptography.ShannonEntropy([]byte(input))
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

	shifts := cryptography.CaesarShifts([]byte(input))
	bestScore, bestShift := 0.0, 0
	for i, s := range shifts {
		if i > 0 {
			score := Score(checker, string(s))
			if score > bestScore {
				bestScore = score
				bestShift = i
			}
		}
	}
	if bestScore > 0.5 {
		results = append(results, fmt.Sprintf("Caesar shift of %d might be English: %s (%.5f)", bestShift, shifts[bestShift], bestScore))
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

	if score := Score(checker, odds.String()); score > 0.5 {
		results = append(results, fmt.Sprintf("Alternating characters might be English: %s (%.5f)", odds.String(), score))
	}

	if score := Score(checker, evens.String()); score > 0.5 {
		results = append(results, fmt.Sprintf("Alternating characters might be English: %s (%.5f)", evens.String(), score))
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

	dists := cryptography.LetterDistribution([]byte(input))
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
				num = 10*num + d
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

func analyseWordCount(_ *SpellChecker, input string) []string {
	var results []string

	if strings.Contains(input, " ") {
		words := strings.Fields(input)
		results = append(results, fmt.Sprintf("%d words", len(words)))
	}

	return results
}

func analysePrimes(checker *SpellChecker, input string) []string {
	var results []string

	output := strings.Builder{}
	for i := range input {
		if big.NewInt(int64(i + 1)).ProbablyPrime(0) {
			output.WriteByte(input[i])
		}
	}

	if score := Score(checker, output.String()); score > 0.5 {
		results = append(results, fmt.Sprintf("Prime characters might be English: %s (%.5f)", output.String(), score))
	}

	return results
}

func analyseCommonLetters(_ *SpellChecker, input string) []string {
	words := strings.Fields(strings.ToLower(input))

	var matches [26]int
	for i := range words {
		for c := range words[i] {
			letter := words[i][c]
			if letter >= 'a' && letter <= 'z' {
				matches[letter-'a']++
			}
		}
	}

	common := ""
	for i := range matches {
		if matches[i] == len(words) {
			common += string(rune('A' + i))
		}
	}

	if len(common) > 0 {
		return []string{fmt.Sprintf("All words contain the letters: %s", common)}
	} else {
		return nil
	}
}

var analysers = []analyser{
	analyseEntropy,
	analyseDataReferences,
	analyseCaesarShifts,
	analyseAlternateChars,
	analysePrimes,
	analyseCommonLetters,
	analyseLength,
	analyseDistribution,
	analyseRunLengthEncoding,
	analyseWordCount,
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
	density := scoreWord(checker, input)
	entropy := scoreEntropy(input)
	bigram := scoreBigrams(input)
	ioc := scoreIoc(input)

	return density * entropy * bigram * ioc
}

// scoreWord returns a score for the text based on how many english words occur within it.
func scoreWord(checker *SpellChecker, input string) float64 {
	words := make([]int, len(input))
	findWords(checker, input, func(start, end int) {
		for i := start; i < end; i++ {
			words[i]++
		}
	})

	mean := float64(0)
	total := 0
	for i := range words {
		if input[i] == ' ' {
			continue
		}

		total++
		mean += float64(words[i])
	}
	mean /= float64(total)

	const targetDensity = 4.0
	return math.Min(math.Pow(math.Max(mean/targetDensity, 0.01), 2), 1.0)
}

// scoreEntropy returns a score for the text based on whether it has a Shannon entropy typical of English text.
func scoreEntropy(input string) float64 {
	entropy := cryptography.ShannonEntropy([]byte(input))
	score := 1.0
	if entropy < 3.5 {
		score = math.Max(entropy/3.5, 0.1)
	} else if entropy > 5 {
		score = math.Max(1-(entropy-5)/3, 0.1)
	}
	return score
}

// scoreBigrams returns a score for the text based on whether it has a bigram distribution typical of English text.
func scoreBigrams(input string) float64 {
	score := float64(0)
	cleaned := strings.ToUpper(nonLetterRegex.ReplaceAllString(strings.ToLower(input), ""))
	for i := range cleaned {
		if i+1 < len(cleaned) {
			b := data.Bigrams[cleaned[i:i+2]]
			score += math.Log10(math.Max(b, 0.0001))
		}
	}
	return math.Pow((10+score/float64(len(cleaned)))/float64(10), 2)
}

// scoreIoc returns a score for the text based on its Index of Coincidence compared to English.
func scoreIoc(input string) float64 {
	return 1 - math.Min(math.Abs(cryptography.IndexOfCoincidence([]byte(input))-cryptography.IndexOfCoincidenceEnglish), 1)
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
