package skiplist

import (
	"math/rand"
	"testing"
	"time"
)

func TestExample(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	var insertNum = 2
	sl := NewSkipList()
	for i := 0; i < insertNum; i++ {
		sl.Insert(i, i)
	}

	var i int
	for node := sl.header.forward[0]; node != nil; node = node.forward[0] {
		if node.key != i {
			t.Error("insert", node)
		}
		i++
	}

	v, ok := sl.Search(1)
	if !ok || v != 1 {
		t.Error("search")
	}

	sl.Delete(1)
}

func TestInsert(t *testing.T) {
	l := NewSkipList()
	l.Insert(0, 0)
}

func BenchmarkInsert(b *testing.B) {
	l := NewSkipList()
	for i := 0; i < b.N; i++ {
		l.Insert(i, i)
	}
}
