package golinewrap

import (
	"bytes"
	"fmt"
	"io"
	"unicode/utf8"
)

// Writer is a structure that writes to the underlying io.Writer, but forces
// line wrapping at the specified width.
type Writer struct {
	io.Writer
	bb           *bytes.Buffer
	width        int
	prefixLength int
	prefix       []byte
}

// New returns a new Writer using the specified width and prefix string for each
// line.
func New(w io.Writer, width int, prefix string) (*Writer, error) {
	prefixBytes := []byte(prefix)
	prefixLength := len([]rune(prefix))

	if width <= 0 || width <= prefixLength {
		return nil, fmt.Errorf("cannot create Writer unless width (%d) is greater than zero and greater than length of prefix: %d.", width, prefixLength)
	}

	ww := &Writer{
		Writer:       w,
		bb:           bytes.NewBuffer(make([]byte, 0, width)),
		width:        width,
		prefixLength: prefixLength,
	}

	if prefixLength > 0 {
		ww.prefix = prefixBytes
	}

	return ww, nil
}

// Write writes the contents of buf to the underlying io.Writer, wrapping lines
// as necessary to prevent line lengths from exceeding the pre-configured
// width. Each pilcrow rune, ¶, in the byte sequence causes a new paragraph to
// be emitted.
func (ww *Writer) Write(buf []byte) (int, error) {
	// When there is nothing to write, we still want to invoke the underlying
	// io.Writer with the empty buffer.
	if len(buf) == 0 {
		return ww.Writer.Write(buf)
	}

	var err error
	var tw int // total written

	// Need to output a paragraph for each pilcrow. Recall that a pilcrow rune
	// requires more than a single byte, and bytes.Split wants a slice of bytes
	// as the delimiter.
	paragraphs := bytes.Split(buf, []byte("¶"))

	for _, paragraph := range paragraphs {
		// When there is a prefix, and line buffer is empty, append the prefix
		// to the line buffer.
		if ww.prefixLength > 0 && ww.bb.Len() == 0 {
			if _, err = ww.bb.Write(ww.prefix); err != nil {
				return tw, err
			}
		}

		words := bytes.Fields(paragraph)

		for _, word := range words {
			// When there is not enough room in the line buffer to append this
			// word followed by a space rune, then write out line buffer
			// followed by a newline character.
			if ww.bb.Len()+utf8.RuneCount(word)+1 >= ww.width {
				// Maybe the line buffer is empty, and this word is too long to
				// fit on a line. Only flush line buffer if it is not empty.
				if ww.bb.Len() > 0 {
					// Change final space to newline rune.
					b := ww.bb.Bytes()
					b[len(b)-1] = '\n'

					nw, err := ww.bb.WriteTo(ww.Writer)
					tw += int(nw)
					if err != nil {
						return tw, err
					}

					ww.bb.Reset()

					if ww.prefixLength > 0 {
						if _, err = ww.bb.Write(ww.prefix); err != nil {
							return tw, err
						}
					}
				}
			}

			// The line buffer is ready for appending this word followed by a
			// space rune.
			if _, err = ww.bb.Write(append(word, ' ')); err != nil {
				return tw, err
			}
		}
		// POST: All words for this paragraph written into bb.

		// Change final space to newline rune.
		if ww.bb.Len() > 0 {
			b := ww.bb.Bytes()
			b[len(b)-1] = '\n'

			nw, err := ww.bb.WriteTo(ww.Writer)
			tw += int(nw)
			if err != nil {
				return tw, err
			}
			ww.bb.Reset()
		}
	}

	return tw, nil
}
