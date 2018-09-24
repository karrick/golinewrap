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
			got := emit(t, 16, ">", "one two three")
			if want := ">one two three\n>\n"; got != want {
				t.Errorf("GOT: %q; WANT: %q", got, want)
			}
		})

		t.Run("exactly one line", func(t *testing.T) {
			got := emit(t, 15, ">", "one two three")
			if want := ">one two three\n>\n"; got != want {
				t.Errorf("GOT: %q; WANT: %q", got, want)
			}
		})

		t.Run("between one and two lines", func(t *testing.T) {
			got := emit(t, 14, ">", "one two three")
			if want := ">one two\n>three\n>\n"; got != want {
				t.Errorf("GOT: %q; WANT: %q", got, want)
			}
		})

		t.Run("exactly two lines", func(t *testing.T) {
			got := emit(t, 5, ">", "one two")
			if want := ">one\n>two\n>\n"; got != want {
				t.Errorf("GOT: %q; WANT: %q", got, want)
			}
		})
	})

	t.Run("with trailing newline", func(t *testing.T) {
		t.Run("less than one line", func(t *testing.T) {
			got := emit(t, 16, ">", "one two three\n")
			if want := ">one two three\n>\n"; got != want {
				t.Errorf("GOT: %q; WANT: %q", got, want)
			}
		})

		t.Run("exactly one line", func(t *testing.T) {
			got := emit(t, 15, ">", "one two three\n")
			if want := ">one two three\n>\n"; got != want {
				t.Errorf("GOT: %q; WANT: %q", got, want)
			}
		})

		t.Run("between one and two lines", func(t *testing.T) {
			got := emit(t, 14, ">", "one two three\n")
			if want := ">one two\n>three\n>\n"; got != want {
				t.Errorf("GOT: %q; WANT: %q", got, want)
			}
		})

		t.Run("exactly two lines", func(t *testing.T) {
			got := emit(t, 5, ">", "one two\n")
			if want := ">one\n>two\n>\n"; got != want {
				t.Errorf("GOT: %q; WANT: %q", got, want)
			}
		})
	})
}

func TestWriteParagraphMultiple(t *testing.T) {
	bb := new(bytes.Buffer)

	lw, err := golinewrap.New(bb, 80, "")
	if err != nil {
		t.Fatal(err)
	}

	_, err = lw.WriteParagraph("Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec euismod velit nec sollicitudin euismod. Lorem ipsum dolor sit amet, consectetur adipiscing elit. In molestie quam ut faucibus lobortis. Mauris sit amet felis dapibus, condimentum metus quis, volutpat nulla. Morbi magna felis, pellentesque vel pellentesque vitae, suscipit quis felis. Donec porta tincidunt nisl id tempus. Cras eros mi, dapibus in laoreet quis, hendrerit et nisl. Quisque dapibus lectus sem, a laoreet turpis accumsan at.")
	if err != nil {
		t.Fatal(err)
	}

	_, err = lw.WriteParagraph("Vestibulum vitae mi vitae nulla suscipit mollis id in orci. Phasellus sit amet ante tellus. Proin bibendum erat eget turpis euismod volutpat. Duis a eros tincidunt, rhoncus sapien eget, placerat est. Sed ultrices neque in congue hendrerit. Orci varius natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. Nullam eu pretium sem, nec viverra lorem. Nulla vehicula nibh vitae vestibulum consequat. Nulla rhoncus ex turpis, nec eleifend mauris finibus quis.")
	if err != nil {
		t.Fatal(err)
	}

	_, err = lw.WriteParagraph("Morbi maximus sed magna sed vestibulum. Praesent varius massa quis ex laoreet placerat a at lorem. Proin sollicitudin fringilla tincidunt. Ut vel libero vitae lorem tincidunt iaculis. Pellentesque nec hendrerit justo, quis ornare lectus. Curabitur ullamcorper sem eget purus molestie venenatis. Integer in rhoncus est, a tincidunt tellus.")
	if err != nil {
		t.Fatal(err)
	}

	_, err = lw.WriteParagraph("Nam id mattis erat, mattis porta mi. Sed posuere erat et leo hendrerit convallis. Praesent id ex tincidunt, ullamcorper nunc eu, venenatis odio. Nunc interdum vitae enim nec mollis. Nullam nulla metus, ornare eget fermentum sed, iaculis ut felis. Suspendisse rhoncus accumsan tellus, quis laoreet diam pretium sit amet. Fusce velit ante, imperdiet quis felis ac, aliquet rhoncus ante. Vivamus imperdiet leo arcu, non venenatis enim semper ut. Suspendisse non est ligula. Vivamus tristique aliquet dolor, sed pellentesque velit. Integer urna lorem, aliquet et porttitor nec, ultricies et nunc. Nullam lacinia arcu eu porttitor auctor.")
	if err != nil {
		t.Fatal(err)
	}

	_, err = lw.WriteParagraph("Mauris porttitor pretium felis, ut venenatis ex ultrices non. Morbi eget urna venenatis, ultricies lorem vitae, suscipit lorem. Nam accumsan urna sit amet augue faucibus dapibus at quis ex. Aliquam sit amet fringilla mauris, tempor egestas tellus. Proin vel facilisis diam. Sed vulputate tortor tortor, ac pulvinar mi blandit vel. Ut tincidunt eros a congue venenatis. Aenean efficitur facilisis libero, eu laoreet augue. Maecenas libero ante, tincidunt id massa sit amet, pulvinar fermentum felis. Nulla quis tristique dolor. Nunc ut pulvinar nunc, nec lacinia sapien. Ut dolor augue, porta vel blandit nec, posuere in tellus. Aliquam nibh arcu, feugiat eu dapibus in, mattis mollis erat.")
	if err != nil {
		t.Fatal(err)
	}

	got := string(bb.Bytes())

	want := `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec euismod velit
nec sollicitudin euismod. Lorem ipsum dolor sit amet, consectetur adipiscing
elit. In molestie quam ut faucibus lobortis. Mauris sit amet felis dapibus,
condimentum metus quis, volutpat nulla. Morbi magna felis, pellentesque vel
pellentesque vitae, suscipit quis felis. Donec porta tincidunt nisl id tempus.
Cras eros mi, dapibus in laoreet quis, hendrerit et nisl. Quisque dapibus
lectus sem, a laoreet turpis accumsan at.

Vestibulum vitae mi vitae nulla suscipit mollis id in orci. Phasellus sit amet
ante tellus. Proin bibendum erat eget turpis euismod volutpat. Duis a eros
tincidunt, rhoncus sapien eget, placerat est. Sed ultrices neque in congue
hendrerit. Orci varius natoque penatibus et magnis dis parturient montes,
nascetur ridiculus mus. Nullam eu pretium sem, nec viverra lorem. Nulla
vehicula nibh vitae vestibulum consequat. Nulla rhoncus ex turpis, nec eleifend
mauris finibus quis.

Morbi maximus sed magna sed vestibulum. Praesent varius massa quis ex laoreet
placerat a at lorem. Proin sollicitudin fringilla tincidunt. Ut vel libero
vitae lorem tincidunt iaculis. Pellentesque nec hendrerit justo, quis ornare
lectus. Curabitur ullamcorper sem eget purus molestie venenatis. Integer in
rhoncus est, a tincidunt tellus.

Nam id mattis erat, mattis porta mi. Sed posuere erat et leo hendrerit
convallis. Praesent id ex tincidunt, ullamcorper nunc eu, venenatis odio. Nunc
interdum vitae enim nec mollis. Nullam nulla metus, ornare eget fermentum sed,
iaculis ut felis. Suspendisse rhoncus accumsan tellus, quis laoreet diam
pretium sit amet. Fusce velit ante, imperdiet quis felis ac, aliquet rhoncus
ante. Vivamus imperdiet leo arcu, non venenatis enim semper ut. Suspendisse non
est ligula. Vivamus tristique aliquet dolor, sed pellentesque velit. Integer
urna lorem, aliquet et porttitor nec, ultricies et nunc. Nullam lacinia arcu eu
porttitor auctor.

Mauris porttitor pretium felis, ut venenatis ex ultrices non. Morbi eget urna
venenatis, ultricies lorem vitae, suscipit lorem. Nam accumsan urna sit amet
augue faucibus dapibus at quis ex. Aliquam sit amet fringilla mauris, tempor
egestas tellus. Proin vel facilisis diam. Sed vulputate tortor tortor, ac
pulvinar mi blandit vel. Ut tincidunt eros a congue venenatis. Aenean efficitur
facilisis libero, eu laoreet augue. Maecenas libero ante, tincidunt id massa
sit amet, pulvinar fermentum felis. Nulla quis tristique dolor. Nunc ut
pulvinar nunc, nec lacinia sapien. Ut dolor augue, porta vel blandit nec,
posuere in tellus. Aliquam nibh arcu, feugiat eu dapibus in, mattis mollis
erat.

`

	gotLines := strings.Split(got, "\n")
	wantLines := strings.Split(want, "\n")

	for i := 0; i < len(gotLines) || i < len(wantLines); i++ {
		if i >= len(gotLines) {
			t.Errorf("INDEX: %d; WANT: %q", i, wantLines[i])
		} else if i >= len(wantLines) {
			t.Errorf("INDEX: %d; GOT: %q", i, gotLines[i])
		} else if gotLines[i] != wantLines[i] {
			t.Errorf("\nINDEX: %d; GOT:\n    %s\nWANT:\n    %s", i, gotLines[i], wantLines[i])
		}
	}

	if got != want {
		t.Logf("\nGOT:\n    %q\nWANT:\n    %q", got, want)
	}
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
		want := ">One two\n>three four\n>five six\n>seven eight\n>nine ten.\n>\n>One two\n>three four\n>five six\n>seven eight\n>nine ten.\n>\n>One two\n>three four\n>five six\n>seven eight\n>nine ten.\n>\n"
		if got != want {
			t.Errorf("\nGOT:\n    %q\nWANT:\n    %q", got, want)
		}
	})

	t.Run("with trailing newline", func(t *testing.T) {
		got := emit(t, 13, ">", "One two three four five six seven eight nine ten.\nOne two three four five six seven eight nine ten.\nOne two three four five six seven eight nine ten.\n")
		want := ">One two\n>three four\n>five six\n>seven eight\n>nine ten.\n>\n>One two\n>three four\n>five six\n>seven eight\n>nine ten.\n>\n>One two\n>three four\n>five six\n>seven eight\n>nine ten.\n>\n>\n>\n"
		if got != want {
			t.Errorf("\nGOT:\n    %q\nWANT:\n    %q", got, want)
		}
	})
}
