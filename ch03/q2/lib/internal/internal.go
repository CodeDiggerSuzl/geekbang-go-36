package internal

import (
	"fmt"
	"io"
)

func Hello(w io.Writer, name string) {
	fmt.Fprintf(w, "hello%s!\n", name)
}
