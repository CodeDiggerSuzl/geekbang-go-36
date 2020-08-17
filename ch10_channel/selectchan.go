package main

import (
	"fmt"
	"time"
)

func main() {
	example()
}
func example() {
	intChan := make(chan int, 1)
	time.AfterFunc(time.Second, func() {
		close(intChan)
	})

	select {
	case _, ok := <-intChan:
		if !ok {
			fmt.Println("The candidate case is closed")
			break
		}
		fmt.Println("The candidate case is selected")
	}
}
