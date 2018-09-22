package golinewrap_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/karrick/golinewrap"
)

func TestNewReturnsErrorWhenPrefixTooLong(t *testing.T) {
	prefix := "1234567890"
	width := len(prefix)
	_, err := golinewrap.New(os.Stderr, width, prefix)
	if want := "columns"; err == nil || !strings.Contains(err.Error(), want) {
		t.Errorf("GOT: %v; WANT: %v", err, want)
	}
}

func TestWriteRune(t *testing.T) {
	emit := func(t *testing.T, width int, prefix string, rr []rune) string {
		bb := new(bytes.Buffer)

		lw, err := golinewrap.New(bb, width, prefix)
		if err != nil {
			t.Fatal(err)
		}

		for _, r := range rr {
			_, err = lw.WriteRune(r)
			if err != nil {
				t.Fatal(err)
			}
		}

		return string(bb.Bytes())
	}

	t.Run("less than one line", func(t *testing.T) {
		got := emit(t, 5, ">", []rune{'1', '2'})
		if want := ">12"; got != want {
			t.Errorf("GOT: %q; WANT: %q", got, want)
		}
	})

	t.Run("exactly one line", func(t *testing.T) {
		got := emit(t, 5, ">", []rune{'1', '2', '3'})
		if want := ">123"; got != want {
			t.Errorf("GOT: %q; WANT: %q", got, want)
		}
	})

	t.Run("between first and second line", func(t *testing.T) {
		got := emit(t, 5, ">", []rune{'1', '2', '3', '4'})
		if want := ">123\n>4"; got != want {
			t.Errorf("GOT: %q; WANT: %q", got, want)
		}
	})

	t.Run("exactly two lines", func(t *testing.T) {
		got := emit(t, 5, ">", []rune{'1', '2', '3', '4', '5', '6'})
		if want := ">123\n>456"; got != want {
			t.Errorf("GOT: %q; WANT: %q", got, want)
		}
	})
}

func TestWriteWord(t *testing.T) {
	emit := func(t *testing.T, width int, prefix string, ww []string) string {
		bb := new(bytes.Buffer)

		lw, err := golinewrap.New(bb, width, prefix)
		if err != nil {
			t.Fatal(err)
		}

		for _, w := range ww {
			_, err = lw.WriteWord(w)
			if err != nil {
				t.Fatal(err)
			}
		}

		return string(bb.Bytes())
	}

	t.Run("less than one line", func(t *testing.T) {
		got := emit(t, 10, ">", []string{"one"})
		if want := ">one"; got != want {
			t.Errorf("GOT: %q; WANT: %q", got, want)
		}
	})

	t.Run("exactly one line", func(t *testing.T) {
		t.Run("one word", func(t *testing.T) {
			got := emit(t, 9, ">", []string{"exactly"})
			if want := ">exactly"; got != want {
				t.Errorf("GOT: %q; WANT: %q", got, want)
			}
		})

		t.Run("two words", func(t *testing.T) {
			got := emit(t, 9, ">", []string{"one", "two"})
			if want := ">one two"; got != want {
				t.Errorf("GOT: %q; WANT: %q", got, want)
			}
		})
	})

	t.Run("between one and two lines", func(t *testing.T) {
		got := emit(t, 9, ">", []string{"another", "test"})
		if want := ">another\n>test"; got != want {
			t.Errorf("GOT: %q; WANT: %q", got, want)
		}
	})

	t.Run("exactly two lines", func(t *testing.T) {
		got := emit(t, 9, ">", []string{"another", "another"})
		if want := ">another\n>another"; got != want {
			t.Errorf("GOT: %q; WANT: %q", got, want)
		}
	})

	t.Run("between two and three lines", func(t *testing.T) {
		got := emit(t, 9, ">", []string{"another", "another", "another"})
		if want := ">another\n>another\n>another"; got != want {
			t.Errorf("GOT: %q; WANT: %q", got, want)
		}
	})
}

func TestWriteParagraph(t *testing.T) {
	emit := func(t *testing.T, width int, prefix string, p string) string {
		bb := new(bytes.Buffer)

		lw, err := golinewrap.New(bb, width, prefix)
		if err != nil {
			t.Fatal(err)
		}

		_, err = lw.WriteParagraph(p)
		if err != nil {
			t.Fatal(err)
		}

		return string(bb.Bytes())
	}

	t.Run("without trailing newline", func(t *testing.T) {
		t.Run("less than one line", func(t *testing.T) {
			got := emit(t, 15, ">", "one two three")
			if want := ">one two three\n"; got != want {
				t.Errorf("GOT: %q; WANT: %q", got, want)
			}
		})

		t.Run("exactly one line", func(t *testing.T) {
			got := emit(t, 14, ">", "one two three")
			if want := ">one two three\n"; got != want {
				t.Errorf("GOT: %q; WANT: %q", got, want)
			}
		})

		t.Run("between one and two lines", func(t *testing.T) {
			got := emit(t, 13, ">", "one two three")
			if want := ">one two three\n"; got != want {
				t.Errorf("GOT: %q; WANT: %q", got, want)
			}
		})

		t.Run("exactly two lines", func(t *testing.T) {
			got := emit(t, 4, ">", "one two")
			if want := ">one\n>two\n"; got != want {
				t.Errorf("GOT: %q; WANT: %q", got, want)
			}
		})
	})

	t.Run("with trailing newline", func(t *testing.T) {
		t.Run("less than one line", func(t *testing.T) {
			got := emit(t, 15, ">", "one two three\n")
			if want := ">one two three\n"; got != want {
				t.Errorf("GOT: %q; WANT: %q", got, want)
			}
		})

		t.Run("exactly one line", func(t *testing.T) {
			got := emit(t, 14, ">", "one two three\n")
			if want := ">one two three\n"; got != want {
				t.Errorf("GOT: %q; WANT: %q", got, want)
			}
		})

		t.Run("between one and two lines", func(t *testing.T) {
			got := emit(t, 13, ">", "one two three\n")
			if want := ">one two three\n"; got != want {
				t.Errorf("GOT: %q; WANT: %q", got, want)
			}
		})

		t.Run("exactly two lines", func(t *testing.T) {
			got := emit(t, 4, ">", "one two\n")
			if want := ">one\n>two\n"; got != want {
				t.Errorf("GOT: %q; WANT: %q", got, want)
			}
		})
	})
}

func TestWrite(t *testing.T) {
	emit := func(t *testing.T, width int, prefix string, s string) string {
		bb := new(bytes.Buffer)

		lw, err := golinewrap.New(bb, width, prefix)
		if err != nil {
			t.Fatal(err)
		}

		_, err = lw.Write([]byte(s))
		if err != nil {
			t.Fatal(err)
		}

		return string(bb.Bytes())
	}

	t.Run("without trailing newline", func(t *testing.T) {
		got := emit(t, 13, ">", "One two three four five six seven eight nine ten.\nOne two three four five six seven eight nine ten.\nOne two three four five six seven eight nine ten.")
		want := ">One two three\n>four five six\n>seven eight\n>nine ten.\n>\n>One two three\n>four five six\n>seven eight\n>nine ten.\n>\n>One two three\n>four five six\n>seven eight\n>nine ten.\n"
		if got != want {
			t.Errorf("\nGOT:\n    %q\nWANT:\n    %q", got, want)
		}
	})

	t.Run("with trailing newline", func(t *testing.T) {
		got := emit(t, 13, ">", "One two three four five six seven eight nine ten.\nOne two three four five six seven eight nine ten.\nOne two three four five six seven eight nine ten.\n")
		want := ">One two three\n>four five six\n>seven eight\n>nine ten.\n>\n>One two three\n>four five six\n>seven eight\n>nine ten.\n>\n>One two three\n>four five six\n>seven eight\n>nine ten.\n>\n>\n"
		if got != want {
			t.Errorf("\nGOT:\n    %q\nWANT:\n    %q", got, want)
		}
	})
}
