// +build golinewrap_debug

package golinewrap

import (
	"fmt"
	"os"
)

// debug formats and prints arguments to stderr for development builds
func debug(f string, a ...interface{}) {
	os.Stderr.Write([]byte("golinewrap: " + fmt.Sprintf(f, a...)))
}
