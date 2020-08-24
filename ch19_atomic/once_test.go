package ch19

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func Test_Once1(t *testing.T) {
	var cnt uint32
	var once sync.Once
	fmt.Printf("The cnt: %d\n", cnt)
	once.Do(func() {
		atomic.AddUint32(&cnt, 1)
	})
	fmt.Printf("The cnt: %d\n", cnt)

	once.Do(func() {
		atomic.AddUint32(&cnt, 2)
	})
	fmt.Printf("The cnt: %d\n", cnt)

	fmt.Println()
}
func Test_Once2(t *testing.T) {
	once := sync.Once{}
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		once.Do(func() {
			for i := 0; i < 3; i++ {
				fmt.Printf("Do task. [1-%d]\n", i)
				time.Sleep(time.Second)
			}
		})
	}()
	go func() {
		defer wg.Done()
		time.Sleep(time.Millisecond * 500)
		once.Do(func() {
			fmt.Println("Do task. [2]")
		})
		fmt.Println("Done. [2]")
	}()
	go func() {
		defer wg.Done()
		time.Sleep(time.Millisecond * 500)
		once.Do(func() {
			fmt.Println("Do task. [3]")
		})
		fmt.Println("Done. [3]")
	}()
	wg.Wait()
	fmt.Println()
}

func Test_Once3(t *testing.T) {
	once := sync.Once{}
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		defer func() {
			if p := recover(); p != nil {
				fmt.Printf("fatal error:%v\n", p)
			}
		}()
		once.Do(func() {
			fmt.Println("Do task. [4]")
			panic(errors.New("something wrong"))
		})
	}()
	go func() {
		defer wg.Done()
		time.Sleep(time.Millisecond * 500)
		once.Do(func() {
			fmt.Println("Do task. [5]") //  will not run here cause the Once
		})
		fmt.Println("Done. [5]")
	}()
	wg.Wait()
}
