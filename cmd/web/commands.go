package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/csmith/cryptography"
	"github.com/csmith/kowalski/v6"
)

func processAnagram(input string) (interface{}, error) {
	input = strings.ToLower(input)
	if !isValidWord(input) {
		return nil, fmt.Errorf("invalid word: %s", input)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	words, err := kowalski.MultiplexAnagram(ctx, checkers, input, kowalski.Dedupe)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"input":  input,
		"result": merge(words),
	}, nil
}

func processAnalysis(input string) (interface{}, error) {
	input = strings.ToLower(input)
	res := kowalski.Analyse(checkers[0], input)

	return map[string]interface{}{
		"input":  input,
		"result": res,
	}, nil
}

func processChunk(input string) (interface{}, error) {
	var parts []int
	words := strings.Split(input, " ")
	for i := range words {
		if v, err := strconv.Atoi(words[i]); err == nil {
			parts = append(parts, v)
		} else {
			break
		}
	}

	if len(parts) == 0 {
		return nil, fmt.Errorf("usage: chunk <size> [size [size [...]]] <text>")
	}

	text := strings.Join(words[len(parts):], "")
	return map[string]interface{}{
		"input":  input,
		"result": kowalski.Chunk(text, parts...),
	}, nil
}

func processLetters(input string) (interface{}, error) {
	res := cryptography.LetterDistribution([]byte(input))

	distribution := make(map[string]int)
	for i := range res {
		distribution[string(byte(i+'A'))] = res[i]
	}

	return map[string]interface{}{
		"input":        input,
		"distribution": distribution,
	}, nil
}

func processMatch(input string) (interface{}, error) {
	input = strings.ToLower(input)
	if !isValidWord(input) {
		return nil, fmt.Errorf("invalid word: %s", input)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	words, err := kowalski.MultiplexMatch(ctx, checkers, input, kowalski.Dedupe)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"input":  input,
		"result": merge(words),
	}, nil
}

func processMorse(input string) (interface{}, error) {
	res := merge(kowalski.MultiplexFromMorse(checkers, input, kowalski.Dedupe))

	return map[string]interface{}{
		"input":  input,
		"result": res,
	}, nil
}

func processMultiAnagram(input string) (interface{}, error) {
	input = strings.ToLower(input)
	if !isValidWord(input) {
		return nil, fmt.Errorf("invalid word: %s", input)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	words, err := kowalski.MultiplexMultiAnagram(ctx, checkers, input, kowalski.Dedupe)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"input":  input,
		"result": merge(words),
	}, nil
}

func processMultiMatch(input string) (interface{}, error) {
	input = strings.ToLower(input)
	if !isValidWord(input) {
		return nil, fmt.Errorf("invalid word: %s", input)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	words, err := kowalski.MultiplexMultiMatch(ctx, checkers, input, kowalski.Dedupe)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"input":  input,
		"result": merge(words),
	}, nil
}

func processOffByOne(input string) (interface{}, error) {
	input = strings.ToLower(input)
	if !isValidWord(input) {
		return nil, fmt.Errorf("invalid word: %s", input)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	words, err := kowalski.MultiplexOffByOne(ctx, checkers, input, kowalski.Dedupe)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"input":  input,
		"result": merge(words),
	}, nil
}

func processShift(input string) (interface{}, error) {
	res := cryptography.CaesarShifts([]byte(input))

	shifts := make([]map[string]interface{}, 0, len(res))
	for i, s := range res {
		score := kowalski.Score(checkers[0], string(s))
		shifts = append(shifts, map[string]interface{}{
			"shift": i,
			"text":  string(s),
			"score": score,
		})
	}

	return map[string]interface{}{
		"input":  input,
		"shifts": shifts,
	}, nil
}

func processT9(input string) (interface{}, error) {
	if !isValidT9(input) {
		return nil, fmt.Errorf("invalid T9 input: %s", input)
	}

	res := merge(kowalski.MultiplexFromT9(checkers, input, kowalski.Dedupe))

	return map[string]interface{}{
		"input":  input,
		"result": res,
	}, nil
}

func processTranspose(input string) (interface{}, error) {
	result := kowalski.Transpose(strings.Split(input, "\n"))

	return map[string]interface{}{
		"input":  input,
		"result": strings.Join(result, "\n"),
	}, nil
}

func processWordSearch(input string) (interface{}, error) {
	input = strings.ToLower(input)
	res := kowalski.MultiplexWordSearch(checkers, strings.Split(input, "\n"))

	return map[string]interface{}{
		"input":  input,
		"normal": countReps(res[0]),
		"updown": countReps(subtract(res[1], res[0])),
	}, nil
}

func processColours(file io.Reader) (interface{}, error) {
	colours, err := kowalski.ExtractColours(file)
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, 0, len(colours))
	for i := range colours {
		if i >= 25 {
			break
		}

		r, g, b, a := colours[i].Colour.RGBA()
		colour := map[string]interface{}{
			"hex":   fmt.Sprintf("#%02x%02x%02x", r/257, g/257, b/257),
			"r":     r / 257,
			"g":     g / 257,
			"b":     b / 257,
			"a":     a / 257,
			"count": colours[i].Count,
		}
		result = append(result, colour)
	}

	return map[string]interface{}{
		"totalColours": len(colours),
		"colours":      result,
		"truncated":    len(colours) > 25,
	}, nil
}

func processHiddenPixels(file io.Reader) (interface{}, error) {
	output, err := kowalski.HiddenPixels(file)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, output); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"image": base64.StdEncoding.EncodeToString(buf.Bytes()),
	}, nil
}

func processRGB(file io.Reader) (interface{}, error) {
	red, green, blue, err := kowalski.SplitRGB(file)
	if err != nil {
		return nil, err
	}

	readImage := func(r io.Reader) (string, error) {
		buf := new(bytes.Buffer)
		if _, err := io.Copy(buf, r); err != nil {
			return "", err
		}
		return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
	}

	redData, err := readImage(red)
	if err != nil {
		return nil, err
	}

	greenData, err := readImage(green)
	if err != nil {
		return nil, err
	}

	blueData, err := readImage(blue)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"red":   redData,
		"green": greenData,
		"blue":  blueData,
	}, nil
}

func processFirstLetters(input string) (interface{}, error) {
	result := kowalski.FirstLetters(input)

	return map[string]interface{}{
		"input":  input,
		"result": result,
	}, nil
}

func processReverse(input string) (interface{}, error) {
	result := kowalski.Reverse(input)

	return map[string]interface{}{
		"input":  input,
		"result": result,
	}, nil
}

func processCheckWords(input string) (interface{}, error) {
	// Get results from all checkers
	allResults := kowalski.MultiplexCheckWords(checkers, input)

	// Format the results for the frontend
	// Combine results from all checkers to show which checker(s) validated each word
	lineCount := len(allResults[0])
	var formattedLines [][]map[string]interface{}

	for lineIdx := 0; lineIdx < lineCount; lineIdx++ {
		var lineWords []map[string]interface{}

		// Get word count from first checker (all should have same structure)
		wordCount := len(allResults[0][lineIdx])

		for wordIdx := 0; wordIdx < wordCount; wordIdx++ {
			word := allResults[0][lineIdx][wordIdx].Word

			// Check which checker(s) validated this word
			validInCheckers := []int{}
			for checkerIdx, checkerResults := range allResults {
				if checkerResults[lineIdx][wordIdx].Valid {
					validInCheckers = append(validInCheckers, checkerIdx)
				}
			}

			lineWords = append(lineWords, map[string]interface{}{
				"word":     word,
				"valid":    len(validInCheckers) > 0,
				"checkers": validInCheckers,
			})
		}
		formattedLines = append(formattedLines, lineWords)
	}

	return map[string]interface{}{
		"input":  input,
		"result": formattedLines,
	}, nil
}

func subtract(input, exclusions []string) []string {
	var res []string
	for i := range input {
		excluded := false
		for j := range exclusions {
			if exclusions[j] == input[i] {
				excluded = true
				break
			}
		}
		if !excluded {
			res = append(res, input[i])
		}
	}
	return res
}

func countReps(input []string) []string {
	counts := make(map[string]int)
	for _, word := range input {
		counts[word]++
	}

	var res []string
	for word, count := range counts {
		if count > 1 {
			res = append(res, fmt.Sprintf("%s × %d", word, count))
		} else {
			res = append(res, word)
		}
	}
	return res
}
