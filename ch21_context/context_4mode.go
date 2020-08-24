package main

import (
	"context"
	"fmt"
	"time"
)

type myKey int

/*
1. with value
*/
func main() {
	keys := []myKey{
		myKey(20),
		myKey(30),
		myKey(60),
		myKey(61),
	}
	values := []string{
		"value in node2",
		"value in node3",
		"value in node6",
		"value in node6_branch",
	}

	rootNode := context.Background()
	node1, cancelFunc1 := context.WithCancel(rootNode)
	defer cancelFunc1()
	// ------------------------- demo1-WithValue -------------------------
	fmt.Println("------------------------- demo1-WithValue -------------------------")
	node2 := context.WithValue(node1, keys[0], values[0])
	node3 := context.WithValue(node2, keys[1], values[1])
	fmt.Printf("The value of the key %v found in the node3: %v\n", keys[0], node3.Value(keys[0]))
	fmt.Printf("The value of the key %v found in the node3: %v\n", keys[1], node3.Value(keys[1]))
	fmt.Printf("The value of the key %v found in the node3: %v\n", keys[2], node3.Value(keys[2]))

	// 示例2。
	fmt.Println("------------------------- demo2-WithTimeout -------------------------")
	node4, _ := context.WithCancel(node3)
	node5, _ := context.WithTimeout(node4, time.Second*5)
	fmt.Printf("The value of the key %v found in the node5: %v\n", keys[0], node5.Value(keys[0]))
	fmt.Printf("The value of the key %v found in the node5: %v\n", keys[1], node5.Value(keys[1]))
	fmt.Println("------------------------- demo2-WithTimeout -------------------------")
}
