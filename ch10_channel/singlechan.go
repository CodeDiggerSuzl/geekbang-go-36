package main

import (
	"fmt"
	"math/rand"
)

func main() {
	// only send
	var sendChan = make(chan<- int, 1)
	// only recv
	var recvChan = make(<-chan int, 1)
	fmt.Printf("The useless channels: %v, %v\n", sendChan, recvChan)

	intCh2 := make(chan int, 3)
	sendInt(intCh2)

	_ = GetIntChan(getIntChan)
}

// a sendInt function
func sendInt(ch chan<- int) {
	ch <- rand.Intn(1000)
}

// Notifier 的接口类型
type Notifier interface {
	sendInt(ch chan<- int)
}

func getIntChan() <-chan int {
	num := 5
	ch := make(chan int, num)
	for i := 0; i < num; i++ {
		ch <- i
	}
	close(ch)
	return ch
}

// GetIntChan get
type GetIntChan func() <-chan int
