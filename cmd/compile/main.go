package main

import (
	"bytes"
	"flag"
	"github.com/csmith/kowalski/v3"
	"io"
	"io/ioutil"
	"log"
	"os"
)

var (
	inFile = flag.String("in", "-", "File to read words from, or '-' for stdin")
	outFile = flag.String("out", "words.wl", "File to write compiled spell checker to")
)

func main() {
	flag.Parse()

	var input io.Reader
	if *inFile == "-" {
		input = os.Stdin
	} else {
		f, err := os.Open(*inFile)
		if err != nil {
			log.Fatalf("Unable to open input: %v", err)
		}
		defer f.Close()
		input = f
	}

	b, err := ioutil.ReadAll(input)
	if err != nil {
		log.Fatalf("Unable to read input: %v", err)
	}

	count := bytes.Count(b, []byte{'\n'})
	reader := bytes.NewReader(b)

	checker, err := kowalski.CreateSpellChecker(reader, count)
	if err != nil {
		log.Fatalf("Unable to create checker: %v", err)
	}

	out, err := os.Create(*outFile)
	if err != nil {
		log.Fatalf("Unable to open output: %v", err)
	}
	defer out.Close()

	err = kowalski.SaveSpellChecker(out, checker)
	if err != nil {
		log.Fatalf("Unable to save checker: %v", err)
	}

	log.Printf("Spell checker with ~%d words successfully saved to %s", count, *outFile)
}
