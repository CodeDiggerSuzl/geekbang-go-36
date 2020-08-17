package main

import (
	"fmt"
)

func main() {
	ch1 := make(chan int, 2)
	ch1 <- 3
	ch1 <- 1
	// ch1 <- 5
	// elem1 := <-ch1
	fmt.Printf("The first element recv from channel :%v\n", <-ch1)
}
