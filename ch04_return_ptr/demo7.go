package main

import (
	"flag"
	"fmt"
)

func main() {
	var name string
	flag.StringVar(&name, "name", "all", "The greeting message")
	flag.Parse()
	fmt.Printf("hello %s\n", name)
}
