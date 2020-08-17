package main

import (
	"flag"
	"fmt"
)

func main() {
	name := getTheFlag()
	flag.Parse()
	fmt.Printf("Helle %s", name)
}
func getTheFlag() *string {
	return flag.String("name", "everyone", "The greeting message")
}
