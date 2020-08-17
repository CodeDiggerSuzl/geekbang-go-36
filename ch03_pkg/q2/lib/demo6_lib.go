package lib

import (
	in "geekbang-go-36/ch03/q2/lib/internal"
	"os"
)

func Hello(name string) {
	in.Hello(os.Stdout, name)
}
