package main

import "fmt"

func main() {
	ch1 := make(chan int, 2)
	// sender
	go func() {
		for i := 0; i < 10; i++ {
			fmt.Printf("Sender: sending element %v ...\n", i)
			ch1 <- i
		}
		fmt.Println("Sender closing...")
		close(ch1)
	}()
	// receiver
	for {
		elem, ok := <-ch1
		if !ok {
			fmt.Println("Receiver: closed channel")
			break
		}
		fmt.Printf("Receiver: received an element: %v\n", elem)
	}
	fmt.Println("end ~~~")
}
