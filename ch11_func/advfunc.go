package main

import (
	"errors"
	"fmt"
)

type operate2 func(x, y int) int

func calc(x, y int, op operate2) (int, error) {
	if op == nil {
		return 0, errors.New("invalid operation")
	}
	return op(x, y), nil
}

func main() {
	x, y := 12, 23
	op := func(x, y int) int {
		return x * y
	}
	result, err := calc(x, y, op)
	fmt.Printf("The result: %d (error: %v)\n", result, err)
}
