package ch07

import (
	"log"
	"testing"
)

func Test_Slice(t *testing.T) {
	s := make([]int, 5)
	log.Println(len(s))
	log.Println(cap(s))
	log.Printf("The value of s %d", s)
	s2 := make([]int, 5, 8)
	log.Println(len(s2))
	log.Printf("The value of s2 %d", s2)
	log.Println(len(s2))
	log.Println(cap(s2))
}
