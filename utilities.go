package kowalski

func unique(words []string) (res []string) {
	last := ""
	for _, w := range words {
		if w != last {
			res = append(res, w)
			last = w
		}
	}
	return
}

func reverse(input []byte, start int) []byte {
	for left, right := start, len(input)-1; left < right; left, right = left+1, right-1 {
		input[left], input[right] = input[right], input[left]
	}
	return input
}

func reverseString(input string) string {
	return string(reverse([]byte(input), 0))
}

func max(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}
