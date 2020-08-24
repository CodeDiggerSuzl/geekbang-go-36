package main

import (
	"errors"
	"fmt"
)

func main() {

}
func echo(request string) (resp string, err error) {
	if request == "" {
		err = errors.New("empty error")
		return
	}
	resp = fmt.Sprintf("echo: %s", request)
	return
}
