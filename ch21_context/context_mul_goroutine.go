package main

import (
	"context"
	"fmt"
	"time"
)

// ctx 控制多个 goroutine
func main() {
	ctx, cancelFunc := context.WithCancel(context.Background())
	go watch(ctx, "first")
	go watch(ctx, "second")
	go watch(ctx, "third")
	time.Sleep(time.Second * 3)
	cancelFunc()
	fmt.Println("all done,exited")
}

func watch(ctx context.Context, name string) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("done watching, exiting...")
			return
		default:
			fmt.Println(name, "goroutine watching....")
			time.Sleep(time.Second * 2)
		}
	}
}
