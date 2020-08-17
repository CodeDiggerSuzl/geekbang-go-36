package main

import "fmt"

type Printer func(content string) (n int, err error)

func printToStd(s string) (n int, err error) {
	return fmt.Println(s)
}

func main() {
	var p Printer
	p = printToStd
	p("hello")
}
