package kowalski

import "strings"

// Chunk takes the input, and splits it up into chunks of the given length. If the input is longer than the list of
// part lengths, the lengths will be repeated.
func Chunk(input string, parts ...int) []string {
	var res []string
	remaining := input
	p := 0
	for len(remaining) > 0 {
		if len(remaining) >= parts[p] {
			res = append(res, remaining[0:parts[p]])
			remaining = remaining[parts[p]:]
		} else {
			res = append(res, remaining)
			remaining = ""
		}

		p = (p + 1) % len(parts)
	}
	return res
}

// Transpose rotates the input text so rows become columns, and columns become rows.
func Transpose(input []string) []string {
	var res []string
	i := 0
	found := true
	for found {
		found = false
		line := strings.Builder{}
		for j := range input {
			if len(input[j]) > i {
				line.WriteByte(input[j][i])
				found = true
			}
		}
		res = append(res, line.String())
		i++
	}
	return res
}
