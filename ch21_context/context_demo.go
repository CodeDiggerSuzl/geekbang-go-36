package main

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"
)

func main() {
	coordinateWithContext()
}
func coordinateWithContext() {
	total := 12
	var num int32
	fmt.Printf("the num : %d {with context.Context}\n", num)
	cxt, cancelFunc := context.WithCancel(context.Background())
	for i := 0; i < total; i++ {
		go addNum2(&num, i, func() {
			if atomic.LoadInt32(&num) == int32(total) {
				cancelFunc()
			}
		})
	}
	<-cxt.Done()
	fmt.Print("end")
}
func addNum2(numP *int32, id int, deferFunc func()) {
	defer func() {
		deferFunc()
	}()
	for i := 0; ; i++ {
		currNum := atomic.LoadInt32(numP)
		newNum := currNum + 1
		time.Sleep(time.Millisecond * 200)
		if atomic.CompareAndSwapInt32(numP, currNum, newNum) {
			fmt.Printf("The number: %d [%d-%d]\n", newNum, id, i)
			break
		} else {
			//fmt.Printf("The CAS operation failed. [%d-%d]\n", id, i)
		}
	}
}
