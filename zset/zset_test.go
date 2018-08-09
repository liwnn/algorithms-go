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
	for i := 0; i < 1000; i++ {
		zs.Add(i, 10000+i)
	}
	if zs.zsl.length != 1000 {
		t.Error("Add error")
	}
}

func TestDelete(t *testing.T) {
	zs := NewZSet()
	zs.Add(1, 1)
	zs.Add(1, 2)
	zs.Add(1, 3)
	zs.Add(1, 4)
	zs.Add(1, 5)
	zs.Delete(3)
	r := zs.Range(1, 5)
	for _, v := range r {
		fmt.Printf("%v ", v.val)
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
		zs.Add(i, i)
	}
	node := zs.zsl.getElementByRank(6)
	fmt.Println(node.ele, node.score)
}

func TestRange(t *testing.T) {
	zs := NewZSet()
	for i := 0; i < 1000; i++ {
		zs.Add(i, 10000+i)
	}
	m := zs.Range(1, 20)
	for _, v := range m {
		fmt.Println(v)
	}
}

func BenchmarkInsert(b *testing.B) {
	sl := zslCreate()
	for i := 0; i < b.N; i++ {
		randData := &Entry{
			val:   i,
			score: i % 20000,
		}
		sl.Insert(i, randData)
	}
}

func BenchmarkAdd(b *testing.B) {
	r := NewZSet()
	for i := 0; i < 5000; i++ {
		r.Add(6*i, i%20000)
	}
	for i := 0; i < b.N; i++ {
		r.Add(10*i, i%20000)
	}
}

func BenchmarkChange(b *testing.B) {
	r := NewZSet()
	for i := 0; i < 5000; i++ {
		r.Add(6*i, i%20000)
	}

	for i := 0; i < b.N; i++ {
		if 10*i < r.MinScore() {
			continue
		}
		r.Add(10*i, i%20000)
		if r.Length() > 5000 {
			r.DeleteHeader()
		}
	}
	if r.Length() != 5000 {
		b.Error("ll")
	}
}

func BenchmarkRange(b *testing.B) {
	zs := NewZSet()
	for i := 0; i < 5000; i++ {
		zs.Add(i, 10000+i)
	}
	for i := 0; i < b.N; i++ {
		zs.Range(850, 870)
	}
}
