package golinewrap

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"
)

func debug(format string, a ...interface{}) (int, error) {
	if true {
		return 0, nil
	}
	return fmt.Fprintf(os.Stderr, format, a...)
}

// Writer is a structure that writes to the underlying io.Writer, but forces
// line wrapping at the specified width.
type Writer struct {
	io.Writer
	lb            *bytes.Buffer
	max           int // max number of columns to fill for each line
	remaining     int // remaining columns in the line buffer
	prefixColumns int // number of columns used by prefix
	prefix        string
}

// New returns a new Writer using the specified width and prefix string for each
// line.
func New(w io.Writer, width int, prefix string) (*Writer, error) {
	prefixColumns := utf8.RuneCountInString(prefix)

	if width <= 0 || width <= prefixColumns {
		return nil, fmt.Errorf("cannot create Writer unless width (%d) is greater than zero and greater than number of columns used by prefix: %d.", width, prefixColumns)
	}

	// NOTE: The line buffer is sized for single byte runes, and will be
	// extended as required when runes that require more than a single byte are
	// emitted.

	ww := &Writer{
		Writer:        w,
		lb:            bytes.NewBuffer(make([]byte, 0, width+1)),
		max:           width,
		prefixColumns: prefixColumns,
	}

	if prefixColumns > 0 {
		_, err := ww.lb.WriteString(prefix)
		if err != nil {
			return nil, err
		}
		ww.prefix = prefix
		ww.remaining = ww.max - ww.prefixColumns
	}

	return ww, nil
}

// flush flushes the contents of line buffer to underlying Writer. This method
// is called at the conclusion of every public method, not necessarily for each
// line.
func (ww *Writer) flush() (int, error) {
	debug("flush: %q\n", string(ww.lb.Bytes()))
	if ww.lb.Len() == 0 {
		return 0, nil
	}
	nw, err := ww.lb.WriteTo(ww.Writer)
	return int(nw), err
}

// newline appends newline to line buffer then flushes to underlying writer
// because this library is line based.
func (ww *Writer) newline() (int, error) {
	debug("newline\n")

	if _, err := ww.lb.WriteRune('\n'); err != nil {
		return 0, err
	}

	// After newline written, the entire line length is available.
	ww.remaining = ww.max

	// Because this library is meant to be line based, go ahead and flush the
	// contents of the line buffer after each newline.
	nw, err := ww.flush()
	if err != nil {
		return nw, err
	}

	return nw, ww.writePrefix()
}

func (ww *Writer) writePrefix() error {
	debug("write prefix\n")

	if ww.prefixColumns == 0 {
		return nil
	}

	ww.remaining -= ww.prefixColumns
	_, err := ww.lb.WriteString(ww.prefix)

	return err
}

// Write writes buf to the underlying io.Writer. It converts the input to a
// string, splits on newline, and emits each line as a paragraph.
func (ww *Writer) Write(buf []byte) (int, error) {
	var tw int
	var err error

	pp := strings.Split(string(buf), "\n")
	ifp := len(pp) - 1

	for i, p := range pp {
		nw, err := ww.WriteParagraph(p)
		tw += nw
		if err != nil {
			return tw, err
		}

		if i < ifp {
			nw, err = ww.newline()
			tw += nw
		}
	}

	return tw, err
}

// WriteParagraph writes p to the underlying io.Writer, wrapping lines as
// necessary to prevent line lengths from exceeding the pre-configured width.
func (ww *Writer) WriteParagraph(p string) (int, error) {
	debug("WriteParagraph(%q): %q; %d\n", p, string(ww.lb.Bytes()), ww.remaining)

	var tw int // total written

	for _, word := range strings.Fields(p) {
		nw, err := ww.writeWord(word)
		tw += nw
		if err != nil {
			return tw, err
		}
	}

	// All words for this paragraph have been written. Write newline to buffer
	// and flush.
	debug("# paragraph complete; line buffer: %q\n", string(ww.lb.Bytes()))

	nw, err := ww.newline()
	tw += nw

	// Do not need to flush again after newline, because we do not want the next
	// prefix to be flushed yet.
	return nw, err
}

// WriteRune writes r to the underlying io.Writer, wrapping lines as necessary
// to prevent line lengths from exceeding the pre-configured width.
func (ww *Writer) WriteRune(r rune) (int, error) {
	var err error
	var tw int

	debug("WriteRune(%q): %q; %d\n", r, string(ww.lb.Bytes()), ww.remaining)

	if ww.remaining < 2 {
		// Not enough room for r and a newline.
		if tw, err = ww.newline(); err != nil {
			return tw, err
		}
	}

	if _, err := ww.lb.WriteRune(r); err != nil {
		return tw, err
	}
	ww.remaining--

	nw, err := ww.flush()
	tw += nw
	return tw, err
}

// WriteWord writes w to the underlying io.Writer, wrapping lines as necessary
// to prevent line lengths from exceeding the pre-configured width.
func (ww *Writer) WriteWord(w string) (int, error) {
	tw, err := ww.writeWord(w)
	if err != nil {
		return tw, err
	}

	nw, err2 := ww.flush()
	tw += nw

	if err == nil {
		err = err2
	}

	return tw, err
}

func (ww *Writer) writeWord(w string) (int, error) {
	var err error
	var tw int // total written

	// Number of columns this word occupies, plus a column for the final space
	// or newline.
	rc := utf8.RuneCountInString(w)

	debug("WriteWord(%q); %q; %d\n", w, string(ww.lb.Bytes()), ww.remaining)

	if ww.remaining < rc {
		// Not enough room for w and a newline.
		if tw, err = ww.newline(); err != nil {
			return tw, err
		}
	}

	if ww.max-ww.remaining != ww.prefixColumns {
		debug("not the first word on line\n")
		// Write a space before this word if it's not the first thing after the
		// prefix.
		if _, err = ww.lb.WriteRune(' '); err != nil {
			return tw, err
		}
		ww.remaining--
	}

	if _, err = ww.lb.WriteString(w); err != nil {
		return tw, err
	}
	ww.remaining -= rc

	return tw, err
}
