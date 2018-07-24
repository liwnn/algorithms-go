package main

import (
	"fmt"
	"math/rand"
)

// Entry e
type Entry struct {
	val int
}

type zskiplistLevel struct {
	forward *zskiplistNode
	span    uint
}

// node is an element of a skip list
type zskiplistNode struct {
	ele   *Entry
	score int
	level []zskiplistLevel
}

func zslCreateNode(level int, score int, ele *Entry) *zskiplistNode {
	zn := &zskiplistNode{
		ele:   ele,
		score: score,
		level: make([]zskiplistLevel, level),
	}
	return zn
}

// SkipList represents a skip list
type SkipList struct {
	header   *zskiplistNode
	level    int // current level count
	maxLevel int
	p        float64

	update []*zskiplistNode // for less alloc
}

// NewSkipList creates a skip list
func zslCreate() *SkipList {
	zsl := &SkipList{
		level:    1,
		maxLevel: 8,    // (1/p)^maxLevel >= maxNode
		p:        0.25, // Skiplist P = 1/4
	}
	zsl.header = zslCreateNode(zsl.maxLevel, 0, nil)
	zsl.update = make([]*zskiplistNode, zsl.maxLevel)
	return zsl
}

// Search for an element by traversing forward pointers
func (list *SkipList) Search(searchKey int) *Entry {
	x := list.header
	for i := list.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil && x.level[i].forward.score < searchKey {
			x = x.level[i].forward
		}
	}
	x = x.level[0].forward
	if x != nil && x.score == searchKey {
		return x.ele
	}
	return nil
}

// Insert element
func (list *SkipList) Insert(score int, ele *Entry) {
	var update = make([]*zskiplistNode, list.maxLevel)
	var rank = make([]uint, list.maxLevel)
	x := list.header
	for i := list.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil && x.level[i].forward.score < score {
			rank[i] += x.level[i].span
			x = x.level[i].forward
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

	x = zslCreateNode(lvl, score, ele)
	for i := 0; i < lvl; i++ {
		x.level[i].forward = update[i].level[i].forward
		update[i].level[i].forward = x
		x.level[i].span = rank[i]
	}
}

// Delete element
func (list *SkipList) Delete(score int, ele int) {
	var update = list.update // [0...list.maxLevel)
	x := list.header
	for i := list.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil && x.level[i].forward.score < score {
			x = x.level[i].forward
		}
		update[i] = x
	}
	x = x.level[0].forward
	if x != nil && x.ele.val == ele {
		for i := 0; i < list.level; i++ {
			if update[i].level[i].forward != x {
				break
			}
			update[i].level[i].forward = x.level[i].forward
		}
		for list.level > 0 && list.header.level[list.level-1].forward == nil {
			list.level--
		}
	}
}

// find the rank for an element
func (list *SkipList) zslGetRank(ele int) uint {
	var rank uint
	x := list.header
	for i := list.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil && x.level[i].forward.ele.val != ele {
			rank += x.level[i].span
			x = x.level[i].forward
		}
	}
	return rank
}

func (list *SkipList) randomLevel() int {
	lvl := 1
	for rand.Float64() < list.p && lvl < list.maxLevel {
		lvl++
	}
	return lvl
}

// ZSet zset
type ZSet struct {
	dict map[int]*Entry
	zsl  *SkipList
}

func newZSet() *ZSet {
	zs := &ZSet{
		dict: make(map[int]*Entry),
		zsl:  zslCreate(),
	}
	return zs
}

// Add a new element or update the score of an existing element
func (zs *ZSet) Add(score int, ele int) {
	de := zs.dict[ele]
	if de != nil {
		curscore := de.val
		// remove and re-insert when score changes
		if score != curscore {
			zs.zsl.Delete(score, ele)
			zs.zsl.Insert(score, de)
			de.val = ele
		}
	} else {
		de = &Entry{
			val: ele,
		}
		zs.zsl.Insert(score, de)
		zs.dict[ele] = de
	}
}

// Delete the element 'ele' from the sorted set,
// return 1 if the element existed and was deleted, 0 otherwise
func (zs *ZSet) Delete(ele int) int {
	de := zs.dict[ele]
	if de == nil {
		return 0
	}
	delete(zs.dict, ele)
	zs.zsl.Delete(de.val, ele)
	return 1
}

// Rank return 0-based rank or -1 if not exist
func (zs *ZSet) Rank(ele int) int {
	de := zs.dict[ele]
	if de != nil {
		rank := zs.zsl.zslGetRank(ele)
		return int(rank)
	}
	return -1
}

func main() {
	zs := newZSet()
	zs.Add(1, 2)
	zs.Delete(2)
	b := zs.zsl.Search(1)
	if b != nil {
		fmt.Println(b.val)
	}
	fmt.Println(zs.Rank(2))
}
