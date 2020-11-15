package kowalski

import (
	"strings"
)

// CaesarShift performs a caesar shift of the given amount on all A-Z characters.
func CaesarShift(input string, count uint8) string {
	builder := strings.Builder{}

	shift := func(c, min byte) byte {
		res := c - min
		for res < min {
			res += 26
		}
		return min + (res + count) % 26
	}

	for i := range input {
		c := input[i]
		if c >= 'a' && c <= 'z' {
			builder.WriteByte(shift(c, 'a'))
		} else if c >= 'A' && c <= 'Z' {
			builder.WriteByte(shift(c, 'A'))
		} else {
			builder.WriteByte(c)
		}
	}
	return builder.String()
}

// CaesarShifts performs all 25 possible caesar shifts on the input.
func CaesarShifts(input string) [25]string {
	var res[25] string
	for i := 1; i <= 25; i++ {
		res[i-1] = CaesarShift(input, uint8(i))
	}
	return res
}
