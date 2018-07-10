package skiplist

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestExample(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	sl := NewSkipList()
	sl.Insert(1, 1)
	sl.Insert(1, 2)
	sl.Insert(1, 3)

	for i := sl.header.forward[0]; i != nil; i = i.forward[0] {
		fmt.Println(i.key, i.value)
	}

	sl.Delete(1)
}

func TestInsert(t *testing.T) {
	l := NewSkipList()
	l.Insert(0, 0)
}

func BenchmarkInsert(b *testing.B) {
	for i := 0; i < b.N; i++ {
		l := NewSkipList()
		l.maxLevel = 8
		for j := 0; j < 5000*3; j++ {
			l.Insert(i/5000, i)
		}
	}
}
