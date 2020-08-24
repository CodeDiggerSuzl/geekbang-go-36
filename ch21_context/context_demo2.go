package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	// context.Background()  返回一个空的 Context，这个空的 Context 一般用于整个 Context 树的根节点。
	// 然后我们使用context.WithCancel(parent)函数，创建一个可取消的子Context，然后当作参数传给goroutine使用，这样就可以使用这个子Context跟踪这个goroutine。
	ctx, cancel := context.WithCancel(context.Background())
	go func(ctx context.Context) {
		for {
			select {
			// 在goroutine中，使用select调用<-ctx.Done()判断是否要结束，如果接受到值的话，就可以返回结束goroutine了；如果接收不到，就会继续进行监控。
			case <-ctx.Done():
				fmt.Println("监控停止, exiting...")
				return
			default:
				fmt.Println("监控中...")
				time.Sleep(2 * time.Second)
			}
		}
	}(ctx)

	time.Sleep(10 * time.Second)
	fmt.Println("通知监控停止...")
	cancel()
	time.Sleep(1 * time.Second)
}
