package main

import (
	"flag"
	"geekbang-go-36/ch03/q2/lib"
)

var name string

func init() {
	flag.StringVar(&name, "name", "everyone", "The greeting message.")
}
func main() {
	flag.Parse()
	lib.Hello(name)
}
