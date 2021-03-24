package kowalski

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
