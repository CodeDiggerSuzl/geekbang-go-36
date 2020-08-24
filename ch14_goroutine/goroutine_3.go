package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

func main() {
	var count uint32
	// trigger 函数会不断地获取一个名叫 count 的变量的值，并判断该值是否与参数 i 的值相同。
	// 如果相同，那么就立即调用 fn 代表的函数，然后把 count 变量的值加 1，最后显式地退出当前的循环。
	// 否则，我们就先让当前的 goroutine“睡眠” 一个纳秒再进入下一个迭代。
	trigger := func(i uint32, fn func()) {
		for {
			if n := atomic.LoadUint32(&count); n == i {
				fn()
				// 操作变量count的时候使用的都是原子操作。
				// 这是由于trigger函数会被多个 goroutine 并发地调用，所以它用到的非本地变量count，就被多个用户级线程共用了。
				// 因此，对它的操作就产生了竞态条件(race condition)，破坏了程序的并发安全性。
				atomic.AddUint32(&count, 1)
				break
			}
			time.Sleep(time.Nanosecond)
		}
	}

	// 在 go 函数中先声明了一个匿名的函数，并把它赋给了变量 fn。
	// 这个匿名函数做的事情很 简单，只是调用 fmt.Println 函数以打印 go 函数的参数 i 的值。
	for i := uint32(0); i < 10; i++ {
		go func(i uint32) {
			fn := func() {
				fmt.Println(i)
			}
			// 调用了一个名叫 trigger 的函数，并把 go 函数的参数 i 和刚刚声明的变量 fn 作为参数传给了它。
			// 注意，for 语句声明的局部变量 i 和 go 函数的参数 i 的类型都变了，都由 int 变为了 uint32
			trigger(i, fn)
		}(i)
	}
	// 该函数接受两个参数，一个是 uint32 类型的参数 i, 另一个是 func() 类型的参数 fn。
	// 你应该记得，func() 代表的是既无参数声明也无结果声明的函数类型。
	trigger(10, func() {})
}
