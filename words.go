package kowalski

import (
	"bufio"
	"fmt"
	"github.com/willf/bloom"
	"os"
	"sort"
	"strings"
)

type Node struct {
	Master      *bloom.BloomFilter
	CrossChecks [2]*bloom.BloomFilter
	Roots       *bloom.BloomFilter
	Counter     int
}

func (n *Node) append(word string) {
	n.Master.AddString(word)
	n.CrossChecks[n.Counter].AddString(word)
	n.Counter = (n.Counter + 1) % len(n.CrossChecks)

	for i := range word {
		n.Roots.AddString(word[0:i+1])
	}
}

func (n *Node) valid(word string) bool {
	if n.Master.TestString(word) == false {
		return false
	}

	for i := range n.CrossChecks {
		if n.CrossChecks[i].TestString(word) {
			return true
		}
	}

	return false
}

// Match returns all valid words that match the given input, expanding '?' as a single character wildcard
func (n *Node) Match(word string) []string {
	res, _ := n.findMatch(strings.ToLower(word))
	return res
}

func (n *Node) findMatch(word string) ([]string, int) {
	maxLength := 0
	stems := []string{""}
	for offset := 0; offset < len(word) && len(stems) > 0; offset++ {
		newStems := make([]string, 0, len(stems))

		var chars []uint8
		if word[offset] == '?' {
			chars = []uint8{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25}
		} else {
			chars = []uint8{word[offset] - 'a'}
		}

		for _, nextChar := range chars {
			for s := range stems {
				word := fmt.Sprintf("%s%c", stems[s], 'a'+nextChar)
				if n.Roots.TestString(word) {
					newStems = append(newStems, word)
				}
			}
		}

		maxLength = offset
		stems = newStems
	}

	var res []string
	for s := range stems {
		if n.valid(stems[s]) {
			res = append(res, stems[s])
		}
	}
	sort.Strings(res)
	return res, maxLength
}

// Anagrams finds all anagrams of the given word, expanding '?' as a single wildcard character
func (n *Node) Anagrams(word string) []string {
	var (
		res        []string
		swapBefore = len(word)
	)

	sortedWord := func(w string) string {
		s := strings.Split(strings.ToLower(w), "")
		sort.Strings(s)
		return strings.Join(s, "")
	}(word)

	for w := []byte(sortedWord); w != nil; w = permute(w, swapBefore+1) {
		matches, count := n.findMatch(string(w))
		if len(matches) > 0 {
			res = append(res, matches...)
			swapBefore = len(word)
		} else {
			swapBefore = count
		}
	}

	sort.Strings(res)
	return unique(res)
}

// LoadWords reads all words from the given file and constructs a Trie for use in future operations
func LoadWords(file string) (*Node, error) {
	root := &Node{
		Master: bloom.NewWithEstimates(4000000, 0.001),
		CrossChecks: [2]*bloom.BloomFilter{
			bloom.NewWithEstimates(2000000, 0.001),
			bloom.NewWithEstimates(2000000, 0.001),
		},
		Roots: bloom.NewWithEstimates(4000000, 0.01),
	}

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.ToLower(scanner.Text())
		if isValidWord(line) {
			root.append(line)
		}
	}

	return root, scanner.Err()
}

func isValidWord(word string) bool {
	if len(word) == 0 {
		return false
	}

	for _, r := range word {
		if r < 'a' || r > 'z' {
			return false
		}
	}
	return true
}

// permute returns the next permutation of the given input, in lexicographical order.
// swapBefore can be used to force a swap within a certain number characters.
func permute(input []byte, swapBefore int) []byte {
	if swapBefore < len(input)-1 {
		input = append(input[0:swapBefore], func(w []byte) []byte {
			s := strings.Split(string(w), "")
			sort.Strings(s)
			return reverse([]byte(strings.Join(s, "")), 0)
		}(input[swapBefore:])...)
	}

	k, l := -1, -1
	for i := range input {
		if i+1 < len(input) && input[i] < input[i+1] {
			k = i
			l = -1
		} else if k >= 0 && input[k] < input[i] {
			l = i
		}
	}

	if k == -1 {
		return nil
	}

	input[k], input[l] = input[l], input[k]
	return reverse(input, k+1)
}

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
