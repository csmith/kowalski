# Kowalski

![Kowalski, analysis](kowalski.jpg)

Kowalski is a Go library for performing various operations to help solve puzzles,
and an accompanying Discord bot.

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

### Image processing

Various utilities to analyse images, find hidden parts, etc.

## Library usage

### SpellChecker

Most functions that involve English text use the `SpellChecker` struct, which can indicate
whether a word is a valid dictionary word or not. To create a new SpellChecker you must
provide it with an `io.reader` where it can read words line-by-line, and a rough estimate
of the number of words it will find:

```go
package example

import (
  "os"

  "github.com/csmith/kowalski/v5"
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
  "os"

  "github.com/csmith/kowalski/v5"
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

### Vellum

The `fst` package contains automata for use with the [Vellum](https://github.com/blevesearch/vellum/)
finite state transducer. This allows for quick matching of anagrams, regular expressions, etc,
against a pre-prepared transducer.

To use these you need to open a transducer using Vellum, create an iterator, and then iterate
over it:

```go
package example

import (
  "github.com/blevesearch/vellum"
  "github.com/csmith/kowalski/v5/fst"
)

func create() {
  transducer, err := vellum.Open("some_file.fst")
  if err != nil {
    panic(err)
  }

  i, err := transducer.Search(fst.NewMorseAutomaton("... --- .-.. .. -.-."), nil, nil)
  if err != nil {
    panic(err)
  }

  for err == nil {
    key, val := i.Current()
    println("Do something here with ", key, val)
    err = i.Next()
  }
}
```

## Discord bot

This repository also contains a Discord bot that allows users to perform analysis.

It currently supports these commands:

```
!anagram Attempts to find single-word anagrams, expanding '*' and '?' wildcards
!analysis Analyses text and provides a summary of potentially interesting findings [Aliases: !analyze, !analyse]
!chunk Splits the text into chunks of a given size
!colours Counts the colours within the image [Aliases: !colors]
!hidden Finds hidden pixels in images [Aliases: !hiddenpixels]
!letters Shows a frequency histogram of the number of letters in the input
!match Attempts to expand '?' wildcards to find a single-word match
!morse Attempts to split a morse code input to spell a single word
!multigram Attempts to find multi-word anagrams, expanding '?' wildcards [Aliases: !multianagram]
!multimatch Attempts to expand '?' wildcards to find multi-word matches
!obo Finds all words that are one character different from the input [Aliases: !offbyone, !ob1]
!rgb Splits an image into its red, green and blue channels
!shift Shows the result of the 25 possible caesar shifts [Aliases: !caesar]
!t9 Attempts to treat a series of numbers as T9 input to spell a single word
!transpose Transposes columns to rows and rows to columns
!wordsearch Searches for words in the given text grid
!help Shows this help text
!fstanagram Attempts to find anagrams from wikipedia, expanding '*' wildcards [Aliases: !fstagram]
!fstregex Attempts to find word matches from wikipedia using regexp [Aliases: !fstre]
!fstmorse Attempts to find word matches from wikipedia using morse
```
