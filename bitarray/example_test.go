package bitarray_test

import (
	"fmt"

	"algorithms-go/bitarray"
)

// ExampleBitArray example
func ExampleBitArray() {
	b := bitarray.New(8)

	b.Set(1, true)
	fmt.Println(b.Get(1))
	b.Set(1, false)
	fmt.Println(b.Get(1))

	// Output:
	// true
	// false
}
