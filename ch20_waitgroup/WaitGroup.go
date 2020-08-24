package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	coordinateWithWaitGroup()
}

func coordinateWithWaitGroup() {
	var wg sync.WaitGroup
	wg.Add(2)
	num := int32(0)
	max := int32(10)
	go addNum(&num, 3, max, wg.Done)
	go addNum(&num, 4, max, wg.Done)
	wg.Wait()
}
func addNum(numP *int32, id, max int32, deferFunc func()) {
	defer func() {
		deferFunc()
	}()
	for i := 0; ; i++ {
		currNum := atomic.LoadInt32(numP)
		if currNum > max {
			break
		}
		newNum := currNum + 2
		time.Sleep(time.Millisecond * 200)
		if atomic.CompareAndSwapInt32(numP, currNum, newNum) {
			fmt.Printf("The number: %d [%d-%d]\n,currNum:%d,newNum: %d", newNum, id, i, currNum, newNum)
		} else {
			fmt.Printf("The CAS operation failed. [%d-%d]\n", id, i)
		}
	}
}
