package bitarray_test

import (
	"testing"

	"algorithms-go/bitarray"
)

func TestBitArray(t *testing.T) {
	l := 8
	myBit1 := bitarray.NewBitArray(l)
	myBit1.Set(1, true)
	if !myBit1.Get(1) {
		t.Error("error get")
	}
}
