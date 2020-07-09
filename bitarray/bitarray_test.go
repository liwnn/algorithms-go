package bitarray

import (
	"math"
	"math/rand"
	"testing"
)

func TestBitArraySetGet(t *testing.T) {
	b := New(8)
	var v bool
	v = b.Get(0)
	if v != false {
		t.Error("Default bit is true")
	}

	b.Set(0, true)
	v = b.Get(0)
	if v != true {
		t.Error("Set bit error")
	}

	b.Set(math.MaxUint32, true)
	if !b.Get(math.MaxUint32) {
		t.Error("Set bit math.MaxUint32 error")
	}
}

func BenchmarkSet(b *testing.B) {
	b.StopTimer()
	var size uint32 = 100000
	s := New(size)
	r := rand.New(rand.NewSource(0))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		s.Set(r.Uint32()%size, true)
	}
}

func BenchmarkGet(b *testing.B) {
	b.StopTimer()
	r := rand.New(rand.NewSource(0))
	var sz uint32 = 100000
	s := New(sz)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		s.Get(r.Uint32() % sz)
	}
}
