* Improve efficiency of example fill-paragraph program.

Possible solution is create a io.Reader that converts single newlines
to spaces, and double-newlines to single newlines, kind of similar to
gocrlf.

* Want to be able to simplify writing many different strings to a single paragraph.

Look at the simple example, where we call WriteParagraph with monster
strings. It would be better to be able to invoke a function to build a
paragraph, then invoke a function to mark its end.
