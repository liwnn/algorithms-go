package zset

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestAdd(t *testing.T) {
	zs := NewZSet()
	zs.Add(1, 1)
	zs.Add(10, 2)
	zs.Add(100, 3)
	zs.Add(1000, 4)
	zs.Add(2, 3)
	r := zs.Range(1, 5, true)
	for _, v := range r {
		fmt.Printf("%d:%d ", v.Key(), v.score)
	}
	fmt.Println()
}

func TestDelete(t *testing.T) {
	zs := NewZSet()
	zs.Add(1, 1)
	zs.Add(1, 2)
	zs.Add(1, 3)
	zs.Add(1, 4)
	zs.Add(1, 5)
	zs.Delete(3)
	r := zs.Range(1, 5, true)
	for _, v := range r {
		fmt.Printf("%v ", v.Key())
	}
	fmt.Println()
}

func TestGetRank(t *testing.T) {
	zs := NewZSet()
	zs.Add(1, 1)
	zs.Add(1, 2)
	zs.Add(1, 3)
	zs.Add(1, 4)
	zs.Add(1, 5)
	fmt.Println(zs.Rank(3, true))
}

func TestGetElementByRank(t *testing.T) {
	zs := NewZSet()
	for i := 0; i < 10; i++ {
		zs.Add(uint32(i), uint64(i))
	}
	node := zs.zsl.getElementByRank(6)
	fmt.Println(node.ele)
}

func TestRange(t *testing.T) {
	zs := NewZSet()
	for i := 0; i < 1000; i++ {
		zs.Add(uint32(i), 10000+uint64(i))
	}
	m := zs.Range(1, 20, true)
	for _, v := range m {
		fmt.Println(v)
	}
}

func BenchmarkInsert(b *testing.B) {
	sl := zslCreate()
	for i := 0; i < b.N; i++ {
		randData := &Element{
			key:   uint64(i),
			score: uint32(i),
		}
		sl.insert(randData)
	}
}

func BenchmarkAdd(b *testing.B) {
	r := NewZSet()
	for i := 0; i < b.N; i++ {
		r.Add(uint32(i), uint64(i)%20000)
	}
}

func BenchmarkChange(b *testing.B) {
	r := NewZSet()
	for i := 0; i < 5000; i++ {
		r.Add(uint32(6*i), uint64(i))
	}

	for i := 0; i < b.N; i++ {
		if r.Length() >= 5000 && 10*uint32(i) < r.MinScore() {
			continue
		}
		r.Add(10*uint32(i), uint64(i)%5000)
		if r.Length() > 5000 {
			r.DeleteFirst()
		}
	}
	if r.Length() != 5000 {
		b.Error("ll")
	}
}

func BenchmarkRange(b *testing.B) {
	zs := NewZSet()
	for i := 0; i < 5000; i++ {
		zs.Add(uint32(i), 10000+uint64(i))
	}
	for i := 0; i < b.N; i++ {
		zs.Range(850, 870, true)
	}
}

var s *ZSet

func init() {
	s = NewZSet()
}
func BenchmarkSortedSet_Add(b *testing.B) {
	b.StopTimer()
	// data initialization
	scores := make([]uint32, b.N)
	IDs := make([]uint64, b.N)
	for i := range IDs {
		scores[i] = rand.Uint32() + uint32(rand.Int31n(99))
		IDs[i] = uint64(i) + 100000
	}
	// BCE
	_ = scores[:b.N]
	_ = IDs[:b.N]

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		s.Add(scores[i], IDs[i])
	}
}

func BenchmarkSortedSet_GetRank(b *testing.B) {
	l := (uint64)(s.Length())
	for i := 0; i < b.N; i++ {
		s.Rank(100000+uint64(i)%l, true)
	}
}

func BenchmarkSortedSet_GetDataByRank(b *testing.B) {
	l := s.Length()
	for i := 0; i < b.N; i++ {
		s.zsl.getElementByRank(uint32(i) % l)
	}
}
