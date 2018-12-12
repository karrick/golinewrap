package main

import (
	"fmt"
	"io"
	"os"

	"github.com/karrick/gobls"
	"github.com/karrick/golf"
	"github.com/karrick/golinewrap"
	"github.com/karrick/gorill"
	"github.com/karrick/gows"
)

func main() {
	optWidth := golf.IntP('w', "width", 0, "width of output histogram. 0 implies use tty width")
	golf.Parse()

	if *optWidth == 0 {
		// ignore error; if cannot get window size, use 80
		*optWidth, _, _ = gows.GetWinSize()
		if *optWidth == 0 {
			*optWidth = 80
		}
	}

	// NOTE: When line wrapping based on terminal width, remember to save one
	// column for the newline character, or your output will be stuttered.
	lw, err := golinewrap.New(os.Stdout, *optWidth-1, "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	var ior io.Reader
	if golf.NArg() == 0 {
		ior = os.Stdin
	} else {
		ior = &gorill.FilesReader{Pathnames: golf.Args()}
	}

	scanner := gobls.NewScanner(ior)
	for scanner.Scan() {
		lw.WriteParagraph(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
