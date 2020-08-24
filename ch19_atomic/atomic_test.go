package ch19

import (
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

func TestAtomic(t *testing.T) {
	t.Helper()
	num := uint32(18)
	fmt.Printf("The number : %d\n", num)
	delta := int32(-3)
	// 减法改成 uint
	atomic.AddUint32(&num, uint32(delta))

	fmt.Printf("The number: %d\n", num)
	//
	atomic.AddUint32(&num, ^uint32(-(-3)-1))
	fmt.Printf("The number: %d\n", num)

	fmt.Println("-----")
	forAndCAS1()
	fmt.Println("-----")
	forAndCAS2()
}

// 自旋锁
func forAndCAS1() {
	sign := make(chan struct{}, 2)
	num := int32(0)
	fmt.Printf("The num is %d\n", num)
	// 定时增加 num 的值
	go func() {
		defer func() {
			sign <- struct{}{}
		}()

		for {
			time.Sleep(time.Millisecond * 500)
			newNum := atomic.AddInt32(&num, 2)
			fmt.Printf("The number %d\n", newNum)
			if newNum == 10 {
				break
			}
		}
	}()

	go func() {
		defer func() {
			sign <- struct{}{}
		}()
		// 使用 CAS 实现自旋锁 spinlock
		for {
			if atomic.CompareAndSwapInt32(&num, 10, 0) {
				fmt.Println("The number has go to 0")
				break
			}
			time.Sleep(time.Millisecond * 500)
		}
	}()
	<-sign
	<-sign
}

// 更加简单的互斥锁的操作

func forAndCAS2() {
	sign := make(chan struct{}, 2)
	num := int32(0)
	fmt.Printf("The number: %d\n", num)
	max := int32(20)
	go func(id int, max int32) { // 定时增加num的值。
		defer func() {
			sign <- struct{}{}
		}()
		for i := 0; ; i++ {
			currNum := atomic.LoadInt32(&num)
			if currNum >= max {
				break
			}
			newNum := currNum + 2
			time.Sleep(time.Millisecond * 200)
			if atomic.CompareAndSwapInt32(&num, currNum, newNum) {
				fmt.Printf("*The number: %d [%d-%d]\n", newNum, id, i)
			} else {
				fmt.Printf("*The CAS operation failed. [%d-%d]\n", id, i)
			}
		}
	}(1, max)
	go func(id int, max int32) { // 定时增加num的值。
		defer func() {
			sign <- struct{}{}
		}()
		for j := 0; ; j++ {
			currNum := atomic.LoadInt32(&num)
			if currNum >= max {
				break
			}
			newNum := currNum + 2
			time.Sleep(time.Millisecond * 200)
			if atomic.CompareAndSwapInt32(&num, currNum, newNum) {
				fmt.Printf("-The number: %d [%d-%d]\n", newNum, id, j)
			} else {
				fmt.Printf("-The CAS operation failed. [%d-%d]\n", id, j)
			}
		}
	}(2, max)
	<-sign
	<-sign
}
