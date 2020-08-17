package ch09

import (
	"fmt"
	"testing"
)

func Test_Map(t *testing.T) {
	var m map[string]int
	key := "second"
	// get from a nil map. OK
	elem, ok := m["secode"]
	fmt.Printf("The element paired with key %q in nil map: %d (%v)\n", key, elem, ok)
	fmt.Printf("The length of nil map: %d\n", len(m))
	fmt.Printf("Delete the key-element pair by key %q...\n", key)
	// del from a nil map. OK
	delete(m, key)

	elem = 2
	fmt.Println("Add a key-element pair to a null map...")
	// set in a nil map. NOT OK: panic: assignment to entry in nil map [recovered]
	m["second"] = elem
}
