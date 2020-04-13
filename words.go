package kowalski

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

type Node struct {
	Children [26]*Node
	Valid    bool
}

func (n *Node) traverse(nextChar uint8) *Node {
	if n.Children[nextChar-'a'] == nil {
		n.Children[nextChar-'a'] = &Node{}
	}
	return n.Children[nextChar-'a']
}

func (n *Node) append(word string) {
	if len(word) == 0 {
		n.Valid = true
	} else {
		n.traverse(word[0]).append(word[1:])
	}
}

// Match returns all valid words that match the given input, expanding '?' as a single character wildcard
func (n *Node) Match(word string) []string {
	type match struct {
		text string
		node *Node
	}

	stems := []*match{{"", n}}
	for offset := 0; offset < len(word) && len(stems) > 0; offset++ {
		newStems := make([]*match, 0, len(stems))

		var chars []uint8
		if word[offset] == '?' {
			chars = []uint8{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25}
		} else {
			chars = []uint8{word[offset] - 'a'}
		}

		for _, nextChar := range chars {
			for s := range stems {
				child := stems[s].node.Children[nextChar]
				if child != nil {
					newStems = append(newStems, &match{
						text: fmt.Sprintf("%s%c", stems[s].text, 'a'+nextChar),
						node: child,
					})
				}
			}
		}

		stems = newStems
	}

	var res []string
	for s := range stems {
		if stems[s].node.Valid {
			res = append(res, stems[s].text)
		}
	}
	sort.Strings(res)
	return res
}

// Anagrams finds all anagrams of the given word, expanding '?' as a single wildcard character
func (n *Node) Anagrams(word string) []string {
	chars := make([]int, len(word))
	for i, r := range word {
		chars[i] = int(r)
	}

	var res []string
	permutations(chars, func(rs []int) {
		builder := strings.Builder{}
		for _, r := range rs {
			builder.WriteRune(rune(r))
		}

		res = append(res, n.Match(builder.String())...)
	})

	sort.Strings(res)

	return unique(res)
}

// LoadWords reads all words from the given file and constructs a Trie for use in future operations
func LoadWords(file string) (*Node, error) {
	root := &Node{}

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

func permutations(arr []int, callback func([]int)) {
	var helper func([]int, int)

	helper = func(arr []int, n int) {
		if n == 1 {
			callback(arr)
		} else {
			for i := 0; i < n; i++ {
				helper(arr, n-1)
				if n%2 == 1 {
					tmp := arr[i]
					arr[i] = arr[n-1]
					arr[n-1] = tmp
				} else {
					tmp := arr[0]
					arr[0] = arr[n-1]
					arr[n-1] = tmp
				}
			}
		}
	}

	helper(arr, len(arr))
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
