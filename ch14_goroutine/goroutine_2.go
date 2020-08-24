package main

import (
	"fmt"
)

//"time"
// 用什么手段可以对 goroutine 的启用数量加以限制?
// 1. 使用 time 让主go程等待其他go 程的休眠
// func main() {
// 	num := 10
// 	for i := 0; i < num; i++ {
// 		go func() {
// 			fmt.Println(i)
// 		}()
// 	}
// 	time.Sleep(time.Millisecond * 2000)
// }

// 我在声明通道sign的时候是以chan struct{}作为其类型 的。
// 其中的类型字面量struct{}有些类似于空接口类型interface{}，它代表了既不包 含任何字段也不拥有任何方法的空结构体类型。
func main() {
	num := 10
	sign := make(chan struct{}, num)
	for i := 0; i < num; i++ {
		go func(i int) {
			fmt.Println(i)
			sign <- struct{}{}
		}(i)
	}
	for j := 0; j < num; j++ {
		<-sign
	}
}
