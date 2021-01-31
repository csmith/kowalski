# Kowalski

![Kowalski, analysis](kowalski.jpg)

Kowalski is a Go library for performing various operations to help solve puzzles.

## Supported functions

### Wildcard matching

AKA Crossword Solving. Given an input with one or more missing letters (represented by
`?` characters), returns a list of dictionary words that match.

### Anagram solving

Given a set of letters, possibly including `?` wildcards, checks all possible anagrams
and returns a list of dictionary words that match.

### Morse decoding

Given a Morse-encoded word (represented with `-` and `.` characters) without spaces,
finds all valid dictionary words that could match.

## Usage

The current functions all revolve around the `SpellChecker` struct, which can indicate
whether a word is a valid dictionary word or not. To create a new SpellChecker you must
provide it with an `io.reader` where it can read words line-by-line, and a rough estimate
of the number of words it will find:

```go
package example

import (
    "github.com/csmith/kowalski/v3"
    "os"
)

func create() {
    f, _ := os.Open("file.txt")
    defer f.Close()
    
    checker, err := kowalski.CreateSpellChecker(f, 100)
}
```

As creating a spellchecker can be expensive and cumbersome, Kowalski supports
serialising data to disk and loading it back. Some examples of these serialised
models are available in the [models](models) directory. To load a model:

```go
package example

import (
    "github.com/csmith/kowalski/v3"
    "os"
)

func create() {
    f, _ := os.Open("file.wl")
    defer f.Close()
    
    checker, err := kowalski.LoadSpellChecker(f)
}
```

This repository also contains a command-line tool to generate a new SpellChecker and export
the serialised model:

```
go run cmd/compile -in wordlist.txt -out model.wl
```

## Discord bot

This repository also contains a Discord bot that allows users to perform analysis.

It currently supports these commands:

- `analysis <term>` performs some analysis on the input and returns
  hints about what it could be.

- `match <term>` returns all known words that match the given term,
  where '?' is a single-character wildcard.
  e.g. `match melism?` will return `melisma`.

- `anagram <term>` returns all known anagrams that match the given term,
  where '?' is a single-character wildcard.
  e.g. `anagram lismem?` will return `melisma`.

- `letters <term>` shows a chart of the distribution of the English letters
  (A-Z, ignoring case) in the given term.

- `morse <term>` returns all possible words that match the given morse code
  (specified using `-` and `.`), ignoring all spaces/pauses.
  
- `shift <term>` performs caesar shifts of 1-25 on the term and displays them.

- `t9 <term>` returns all possible words that match the given T9 input
  (specified using numbers).

- `wordsearch <grid>` returns all found words in the word search.
