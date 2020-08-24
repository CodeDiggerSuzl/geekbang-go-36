package main

import (
	"context"
	"fmt"
	"time"
)

var key string = "name"

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	valueCtx := context.WithValue(ctx, key, "flag")
	go watch2(valueCtx)
	time.Sleep(time.Second * 5)
	fmt.Println("ending ")
	cancel()
	time.Sleep(5 * time.Second)

}

func watch2(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println(ctx.Value(key), "监控停止...")
			return
		default:
			fmt.Println(ctx.Value(key), "watching")
			time.Sleep(2 * time.Second)
		}
	}
}
