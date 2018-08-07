package zset

import (
	"math/rand"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestAdd(t *testing.T) {
	zs := NewZSet()
	for i := 0; i < 1000; i++ {
		zs.Add(i, 10000+i)
	}
	if zs.zsl.length != 1000 {
		t.Error("Add error")
	}
}

func TestDelete(t *testing.T) {
	zs := NewZSet()
	for i := 0; i < 1000; i++ {
		zs.Add(i, 10000+i)
	}
}

func BenchmarkAdd(b *testing.B) {
	r := NewZSet()
	r.zsl.maxLevel = 32
	for i := 0; i < b.N; i++ {
		r.Add(6*i, i)
	}
	// for i := 0; i < b.N; i++ {
	// 	if i%5000 <= r.MinScore() {
	// 		continue
	// 	}
	// 	r.Add(10*i, i%5000)
	// 	if r.Length() > 5000 {
	// 		r.DeleteHeader()
	// 	}
	// }
}
