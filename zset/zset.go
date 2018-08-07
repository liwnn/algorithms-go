package zset

import (
	"math/rand"
)

// Entry e
type Entry struct {
	val   int
	score int
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
	header *zskiplistNode
	length uint
	level  int // current level count

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
		if i == list.level-1 {
			rank[i] = 0
		} else {
			rank[i] = rank[i+1]
		}
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
			update[i].level[i].span = list.length
		}
		list.level = lvl
	}

	x = zslCreateNode(lvl, score, ele)
	for i := 0; i < lvl; i++ {
		x.level[i].forward = update[i].level[i].forward
		update[i].level[i].forward = x
		x.level[i].span = (update[i].level[i].span - (rank[0] - rank[i]))
		update[i].level[i].span = (rank[0] - rank[i]) + 1
	}

	for i := lvl; i < list.level; i++ {
		update[i].level[i].span++
	}
	list.length++
}

// delete element
func (list *SkipList) delete(score int, ele int) *zskiplistNode {
	var update = list.update // [0...list.maxLevel)
	x := list.header
	for i := list.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil && x.level[i].forward.score < score {
			x = x.level[i].forward
		}
		update[i] = x
	}
	x = x.level[0].forward
	if x != nil && x.score == score && x.ele.val == ele {
		for i := 0; i < list.level; i++ {
			if update[i].level[i].forward != x {
				break
			}
			update[i].level[i].forward = x.level[i].forward
		}
		for list.level > 0 && list.header.level[list.level-1].forward == nil {
			list.level--
		}
		for i := 1; i < list.level && i < len(x.level); i++ {
			update[i].level[i].span += (x.level[i].span - 1)
		}
		list.length--
		return x
	}
	return nil
}

// Find the rank for an element.
// Return 0-based rank
func (list *SkipList) zslGetRank(score int, ele int) uint {
	var rank uint
	x := list.header
	for i := list.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil &&
			(x.level[i].forward.score < score ||
				(x.level[i].forward.score == score && x.level[i].forward.ele.val != ele)) {
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

// NewZSet create ZSet
func NewZSet() *ZSet {
	zs := &ZSet{
		dict: make(map[int]*Entry),
		zsl:  zslCreate(),
	}
	return zs
}

// Add a new element or update the score of an existing element
func (zs *ZSet) Add(score int, ele int) {
	if de := zs.dict[ele]; de != nil {
		curscore := de.val
		// remove and re-insert when score changes
		if score != curscore {
			zs.zsl.delete(score, ele)
			zs.zsl.Insert(score, de)
			de.val = ele
			de.score = score
		}
	} else {
		de = &Entry{
			val:   ele,
			score: score,
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
	zs.zsl.delete(de.score, ele)
	delete(zs.dict, ele)
	return 1
}

// Rank return 0-based rank or -1 if not exist
func (zs *ZSet) Rank(ele int) int {
	de := zs.dict[ele]
	if de != nil {
		rank := zs.zsl.zslGetRank(de.score, ele)
		return int(rank)
	}
	return -1
}

// Length return the element count
func (zs *ZSet) Length() uint {
	return zs.zsl.length
}

func (zs *ZSet) MinScore() int {
	return zs.zsl.header.level[0].forward.score
}

func (zs *ZSet) DeleteHeader() {
	zs.Delete(zs.zsl.header.level[0].forward.ele.val)
}
