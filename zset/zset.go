package zset

import (
	"math/rand"
)

const (
	skipListMaxLevel = 32
)

// Entry e
type Entry struct {
	val   int
	score int
	order int
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
		maxLevel: skipListMaxLevel, // (1/p)^maxLevel >= maxNode
		p:        0.25,             // Skiplist P = 1/4
	}
	zsl.header = zslCreateNode(zsl.maxLevel, 0, nil)
	zsl.update = make([]*zskiplistNode, zsl.maxLevel)
	return zsl
}

// Insert element
func (list *SkipList) Insert(score int, ele *Entry) {
	var update = list.update
	var rank [skipListMaxLevel]uint
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
	if x.level[0].forward != nil && x.score == x.level[0].forward.score {
		x.ele.order = x.level[0].forward.ele.order + 1
	}

	for i := lvl; i < list.level; i++ {
		update[i].level[i].span++
	}
	list.length++
}

// delete element
func (list *SkipList) delete(score int, ele *Entry) *zskiplistNode {
	var update = list.update // [0...list.maxLevel)
	x := list.header
	for i := list.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil &&
			(x.level[i].forward.score < score ||
				x.level[i].forward.score == score && ele.order < x.level[i].forward.ele.order) {
			x = x.level[i].forward
		}
		update[i] = x
	}
	x = x.level[0].forward
	if x != nil && x.score == score && x.ele.val == ele.val {
		for i := 0; i < list.level; i++ {
			if update[i].level[i].forward == x {
				update[i].level[i].forward = x.level[i].forward
				update[i].level[i].span += x.level[i].span - 1
			} else {
				update[i].level[i].span--
			}
		}
		for list.level > 0 && list.header.level[list.level-1].forward == nil {
			list.level--
		}
		list.length--
		return x
	}
	return nil
}

// Find the rank for an element.
// Returns 0 when the element cannot be found, rank otherwise.
// Note that the rank is 1-based
func (list *SkipList) zslGetRank(score int, ele *Entry) uint {
	var rank uint
	x := list.header
	for i := list.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil &&
			(x.level[i].forward.score < score ||
				x.level[i].forward.score == score && ele.order <= x.level[i].forward.ele.order) {
			rank += x.level[i].span
			x = x.level[i].forward
		}
		if x.ele != nil && x.ele.val == ele.val {
			return rank
		}
	}
	return 0
}

func (list *SkipList) randomLevel() int {
	lvl := 1
	for rand.Float64() < list.p && lvl < list.maxLevel {
		lvl++
	}
	return lvl
}

// Finds an element by its rank. The rank argument needs to be 1-based.
func (list *SkipList) getElementByRank(rank uint) *zskiplistNode {
	var traversed uint
	x := list.header
	for i := list.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil && traversed+x.level[i].span <= rank {
			traversed += x.level[i].span
			x = x.level[i].forward
		}
		if traversed == rank {
			return x
		}
	}
	return nil
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
		// remove and re-insert when score changes
		if score != de.score {
			zs.zsl.delete(de.score, de)
			de.score = score
			zs.zsl.Insert(score, de)
		}
	} else {
		de := &Entry{
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
	zs.zsl.delete(de.score, de)
	delete(zs.dict, ele)
	return 1
}

// Rank return 1-based rank or 0 if not exist
func (zs *ZSet) Rank(ele int, reverse bool) uint {
	de := zs.dict[ele]
	if de != nil {
		rank := zs.zsl.zslGetRank(de.score, de)
		if rank > 0 {
			if reverse {
				return zs.zsl.length - rank + 1
			}
			return rank
		}
	}
	return 0
}

// Range return 1-based elements in [start, end]
func (zs *ZSet) Range(start uint, end uint) []*Entry {
	if start == 0 {
		start = 1
	}
	if end == 0 {
		end = zs.zsl.length
	}
	if start > end || start > zs.zsl.length {
		return nil
	}
	if end > zs.zsl.length {
		end = zs.zsl.length
	}
	rangelen := end - start + 1
	node := zs.zsl.getElementByRank(start)
	var ret = make([]*Entry, 0, rangelen)
	for i := uint(0); i < rangelen; i++ {
		ret = append(ret, node.ele)
		node = node.level[0].forward
	}
	return ret
}

// Length return the element count
func (zs *ZSet) Length() uint {
	return zs.zsl.length
}

// MinScore 最小积分
func (zs *ZSet) MinScore() int {
	if zs.zsl.header.level[0].forward != nil {
		return zs.zsl.header.level[0].forward.score
	}
	return 0
}

func (zs *ZSet) DeleteHeader() {
	zs.zsl.delete(zs.zsl.header.level[0].forward.score, zs.zsl.header.level[0].forward.ele)
	delete(zs.dict, zs.zsl.header.level[0].forward.ele.val)
}
