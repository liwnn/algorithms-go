package skiplist

import (
	"math/rand"
)

// node is an element of a skip list
type node struct {
	key     int
	value   int
	forward []*node
}

func makeNode(level int, key int, value int) *node {
	return &node{
		key:     key,
		value:   value,
		forward: make([]*node, level),
	}
}

// SkipList represents a skip list
type SkipList struct {
	header   *node
	level    int // current level count
	maxLevel int
	p        float64

	update []*node // for less alloc
}

// NewSkipList creates a skip list
func NewSkipList() *SkipList {
	sl := &SkipList{
		level:    1,
		maxLevel: 8,    // (1/p)^maxLevel >= maxNode
		p:        0.25, // Skiplist P = 1/4
	}
	sl.header = makeNode(sl.maxLevel, 0, 0)
	sl.update = make([]*node, sl.maxLevel)
	return sl
}

// Search for an element by traversing forward pointers
func (list *SkipList) Search(searchKey int) (int, bool) {
	x := list.header
	for i := list.level - 1; i >= 0; i-- {
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
func (list *SkipList) Insert(searchKey int, newValue int) {
	var update = list.update // [0...list.maxLevel)
	x := list.header
	for i := list.level - 1; i >= 0; i-- {
		for x.forward[i] != nil && x.forward[i].key < searchKey {
			x = x.forward[i]
		}
		update[i] = x
	}

	lvl := list.randomLevel()
	if lvl > list.level {
		for i := list.level; i < lvl; i++ {
			update[i] = list.header
		}
		list.level = lvl
	}

	x = makeNode(lvl, searchKey, newValue)
	for i := 0; i < lvl; i++ {
		x.forward[i] = update[i].forward[i]
		update[i].forward[i] = x
	}
}

// Delete element
func (list *SkipList) Delete(searchKey int) {
	var update = list.update // [0...list.maxLevel)
	x := list.header
	for i := list.level - 1; i >= 0; i-- {
		for x.forward[i] != nil && x.forward[i].key < searchKey {
			x = x.forward[i]
		}
		update[i] = x
	}
	x = x.forward[0]
	if x != nil && x.key == searchKey {
		for i := 0; i < list.level; i++ {
			if update[i].forward[i] != x {
				break
			}
			update[i].forward[i] = x.forward[i]
		}
		for list.level > 0 && list.header.forward[list.level-1] == nil {
			list.level--
		}
	}
}

func (list *SkipList) randomLevel() int {
	lvl := 1
	for rand.Float64() < list.p && lvl < list.maxLevel {
		lvl++
	}
	return lvl
}
