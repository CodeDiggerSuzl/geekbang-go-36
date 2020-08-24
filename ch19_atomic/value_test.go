package ch19

import (
	"fmt"
	"sync/atomic"
	"testing"
)

func Test_Atomic(t *testing.T) {
	var box atomic.Value
	fmt.Println("Copy box to box2.")
	box2 := box // 原子值在真正使用前可以被复制。
	v1 := [...]int{1, 2, 3}
	fmt.Printf("Store %v to box.\n", v1)
	box.Store(v1)
	fmt.Printf("The value load from box is %v.\n", box.Load())
	fmt.Printf("The value load from box2 is %v.\n", box2.Load())
	fmt.Println()

	v2 := "123"
	fmt.Printf("Store %q to box2.\n", v2)
	box2.Store(v2) // 这里并不会引发panic。
	fmt.Printf("The value load from box is %v.\n", box.Load())
	fmt.Printf("The value load from box2 is %q.\n", box2.Load())
	fmt.Println()

	fmt.Println("Copy box to box3.")
	box3 := box // 原子值在真正使用后不应该被复制！
	fmt.Printf("The value load from box3 is %v.\n", box3.Load())
	v3 := 123
	fmt.Printf("Store %d to box2.\n", v3)
	// box3.Store(v3) // 这里会引发一个panic，报告存储值的类型不一致。
	_ = box3
	fmt.Println()
}
