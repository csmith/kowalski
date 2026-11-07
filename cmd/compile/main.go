package main

import (
	"bytes"
	"flag"
	"io"
	"log/slog"
	"os"

	"github.com/csmith/envflag/v2"
	"github.com/csmith/kowalski/v6"
	"github.com/csmith/slogflags"
)

var (
	inFile  = flag.String("in", "-", "File to read words from, or '-' for stdin")
	outFile = flag.String("out", "words.wl", "File to write compiled spell checker to")
)

func main() {
	envflag.Parse()
	_ = slogflags.Logger(slogflags.WithSetDefault(true))

	var input io.Reader
	if *inFile == "-" {
		input = os.Stdin
	} else {
		f, err := os.Open(*inFile)
		if err != nil {
			slog.Error("Unable to open input", "error", err)
			os.Exit(1)
		}
		defer f.Close()
		input = f
	}

	b, err := io.ReadAll(input)
	if err != nil {
		slog.Error("Unable to read input", "error", err)
		os.Exit(1)
	}

	count := bytes.Count(b, []byte{'\n'})
	reader := bytes.NewReader(b)

	checker, err := kowalski.CreateSpellChecker(reader, count)
	if err != nil {
		slog.Error("Unable to create checker", "error", err)
		os.Exit(1)
	}

	out, err := os.Create(*outFile)
	if err != nil {
		slog.Error("Unable to open output", "error", err)
		os.Exit(1)
	}
	defer out.Close()

	err = kowalski.SaveSpellChecker(out, checker)
	if err != nil {
		slog.Error("Unable to save checker", "error", err)
		os.Exit(1)
	}

	slog.Info("Spell checker successfully saved", "words", count, "path", *outFile)
}
