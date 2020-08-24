package main

import (
	"errors"
	"fmt"
)

func main() {
	defer func() {
		if p := recover(); p != nil {
			fmt.Printf("panic: %s\n", p)
		}
		fmt.Println("exit defer func")
	}()
	panic(errors.New("ops"))

}
