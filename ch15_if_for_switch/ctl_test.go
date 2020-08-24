package ch15

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_Range(t *testing.T) {

	num := []int{1, 2, 3, 4, 5, 6} // slice
	for i := range num {
		if i == 3 {
			num[i] |= i
		}
	}
	fmt.Println(num)
	fmt.Println(4 | 1)
}
func Test_ArrRange(t *testing.T) {
	numbers2 := []int{1, 2, 3} // array
	maxIndex2 := len(numbers2) - 1
	fmt.Println(maxIndex2)
	for i, e := range numbers2 {
		if i == maxIndex2 {
			numbers2[0] = e
		} else {
			numbers2[i+1] = e
		}
	}
	fmt.Println(numbers2)
}

func Test_Range_Cpy(t *testing.T) {
	b := []int{1, 2, 3}
	for i, v := range b { //i,v从a复制的对象里提取出
		if i == 0 {
			b[1], b[2] = 200, 300
			fmt.Println(b) //输出[1 200 300]
		}
		b[i] = v + 100
	}
	fmt.Println(b) // [101 300 400]
}

func Test_Type(t *testing.T) {
	var arr = [3]int{1, 2, 3}
	fmt.Println(reflect.TypeOf(arr))

	var slice = make([]int, 3, 4)
	fmt.Println(reflect.TypeOf(slice))
	fmt.Println(len(arr), len(slice), cap(arr), cap(slice))
}
