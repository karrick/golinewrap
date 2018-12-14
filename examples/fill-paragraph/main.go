package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/karrick/golf"
	"github.com/karrick/golinewrap"
	"github.com/karrick/gorill"
)

func main() {
	optHelp := golf.BoolP('h', "help", false, "display help and exit")
	optWidth := golf.IntP('w', "width", 0, "width of output; 0 implies use tty width")
	golf.Parse()

	if *optHelp {
		// For help text, just line wrap to 79 characters.
		lw, err := golinewrap.New(os.Stderr, 79, "")
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		lw.Printf("%s", filepath.Base(os.Args[0]))
		lw.Printf("Reflow paragraphs to specified line width.")
		lw.Printf("Reads input from multiple files specified on the command line or from standard input when no files are specified.")
		golf.Usage()
		os.Exit(0)
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

	// WARNING: This is very inefficient implementation that slurps up entire
	// file.

	buf, err := ioutil.ReadAll(ior)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}

	for _, paragraph := range strings.Split(string(buf), "\n\n") {
		lw.WriteParagraph(strings.Replace(paragraph, "\n", " ", -1))
	}
}
