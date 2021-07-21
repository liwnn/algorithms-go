package skiplist

import (
	"math/rand"
)

const (
	kMaxLevel = 12   // (1/p)^MaxLevel >= maxNode
	kP        = 0.25 // Skiplist P = 1/4
)

// skipListNode is an element of a skip list
type skipListNode struct {
	key     int
	value   int
	forward []*skipListNode
}

func newNode(level int, key int, value int) *skipListNode {
	return &skipListNode{
		key:     key,
		value:   value,
		forward: make([]*skipListNode, level),
	}
}

// SkipList implemente "Skip Lists: A Probabilistic Alternative to Balanced Trees"
type SkipList struct {
	header *skipListNode
	level  int // current max level
}

// NewSkipList creates a skip list
func NewSkipList() *SkipList {
	sl := &SkipList{
		level:  1,
		header: newNode(kMaxLevel, 0, 0),
	}
	return sl
}

// Search for an element by traversing forward pointers
func (sl *SkipList) Search(searchKey int) (int, bool) {
	x := sl.header
	// loop : x→key < searchKey <= x→forward[i]→key
	for i := sl.level - 1; i >= 0; i-- {
		for x.forward[i] != nil && x.forward[i].key < searchKey {
			x = x.forward[i]
		}
	}
	x = x.forward[0]
	if x != nil && x.key == searchKey {
		return x.value, true
	}
	return -1, false
}

// Insert element
func (sl *SkipList) Insert(searchKey int, newValue int) {
	var prev [kMaxLevel]*skipListNode // [0...list.maxLevel)
	x := sl.header
	for i := sl.level - 1; i >= 0; i-- {
		for x.forward[i] != nil && x.forward[i].key < searchKey {
			x = x.forward[i]
		}
		prev[i] = x
	}
	x = x.forward[0]
	if x != nil && x.key == searchKey {
		x.value = newValue
	} else {
		lvl := sl.randomLevel()
		if lvl > sl.level {
			for i := sl.level; i < lvl; i++ {
				prev[i] = sl.header
			}
			sl.level = lvl
		}

		x = newNode(lvl, searchKey, newValue)
		for i := 0; i < lvl; i++ {
			x.forward[i] = prev[i].forward[i]
			prev[i].forward[i] = x
		}
	}
}

// Delete element
func (sl *SkipList) Delete(searchKey int) {
	var prev [kMaxLevel]*skipListNode // [0...list.maxLevel)
	x := sl.header
	for i := sl.level - 1; i >= 0; i-- {
		for x.forward[i] != nil && x.forward[i].key < searchKey {
			x = x.forward[i]
		}
		prev[i] = x
	}
	x = x.forward[0]
	if x != nil && x.key == searchKey {
		for i := 0; i < sl.level; i++ {
			if prev[i].forward[i] != x {
				break
			}
			prev[i].forward[i] = x.forward[i]
		}
		for sl.level > 1 && sl.header.forward[sl.level-1] == nil {
			sl.level--
		}
	}
}

func (sl *SkipList) randomLevel() int {
	lvl := 1
	for lvl < kMaxLevel && rand.Float64() < kP {
		lvl++
	}
	return lvl
}
