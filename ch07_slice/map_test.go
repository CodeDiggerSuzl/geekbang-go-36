package ch07

import (
	"fmt"
	"testing"
)

func Test_InterfaceKey(t *testing.T) {
	var mapT = map[interface{}]int{
		"1":      1,
		[]int{2}: 2,
		3:        3,
	}
	fmt.Println(mapT)
}
