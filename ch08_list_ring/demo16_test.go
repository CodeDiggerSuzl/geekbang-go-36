package ch08

import (
	"container/list"
	"container/ring"
	"log"
	"testing"
)

func Test_List(t *testing.T) {
	var l list.List
	log.Printf("%v", l)
	var r ring.Ring
	log.Printf("%v", r)
}
