package zset

import (
	"math/rand"
)

const (
	skipListMaxLevel = 32
)

// Element e
type Element struct {
	score uint32
	key   uint64
}

// Score return score
func (e *Element) Score() uint32 {
	return e.score
}

// SetScore set score
func (e *Element) SetScore(score uint32) {
	e.score = score
}

// Key return key
func (e *Element) Key() uint64 {
	return e.key
}

type zSkipListLevel struct {
	forward *zSkipListNode
	span    uint32
}

// node is an element of a skip list
type zSkipListNode struct {
	ele   *Element
	level []zSkipListLevel
	order int
}

func zslCreateNode(level int, ele *Element) *zSkipListNode {
	zn := &zSkipListNode{
		ele:   ele,
		level: make([]zSkipListLevel, level),
	}
	return zn
}

// zSkipList represents a skip list
type zSkipList struct {
	header *zSkipListNode
	length uint32
	level  int // current level count

	maxLevel int
	p        float64

	update []*zSkipListNode // for less alloc
}

// zslCreate creates a skip list
func zslCreate() *zSkipList {
	zsl := &zSkipList{
		level:    1,
		maxLevel: skipListMaxLevel, // (1/p)^maxLevel >= maxNode
		p:        0.25,             // SkipList P = 1/4
	}
	zsl.header = zslCreateNode(zsl.maxLevel, nil)
	zsl.update = make([]*zSkipListNode, zsl.maxLevel)
	return zsl
}

// insert element
func (list *zSkipList) insert(ele *Element) *zSkipListNode {
	var update = list.update
	var rank [skipListMaxLevel]uint32
	x := list.header
	for i := list.level - 1; i >= 0; i-- {
		if i == list.level-1 {
			rank[i] = 0
		} else {
			rank[i] = rank[i+1]
		}
		for x.level[i].forward != nil && x.level[i].forward.ele.Score() < ele.Score() {
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

	x = zslCreateNode(lvl, ele)
	for i := 0; i < lvl; i++ {
		x.level[i].forward = update[i].level[i].forward
		update[i].level[i].forward = x
		x.level[i].span = update[i].level[i].span - (rank[0] - rank[i])
		update[i].level[i].span = (rank[0] - rank[i]) + 1
	}
	if x.level[0].forward != nil && x.ele.Score() == x.level[0].forward.ele.Score() {
		x.order = x.level[0].forward.order + 1
	}

	for i := lvl; i < list.level; i++ {
		update[i].level[i].span++
	}
	list.length++
	return x
}

// delete element
func (list *zSkipList) delete(node *zSkipListNode) *zSkipListNode {
	var update = list.update // [0...list.maxLevel)
	x := list.header
	for i := list.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil &&
			(x.level[i].forward.ele.Score() < node.ele.Score() ||
				x.level[i].forward.ele.Score() == node.ele.Score() && node.order < x.level[i].forward.order) {
			x = x.level[i].forward
		}
		update[i] = x
	}
	x = x.level[0].forward
	if x != nil && x.ele.Score() == node.ele.Score() && x.ele.Key() == node.ele.Key() {
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

// update score
func (list *zSkipList) updateScore(node *zSkipListNode, score uint32) {
	var update = list.update // [0...list.maxLevel)
	x := list.header
	for i := list.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil &&
			(x.level[i].forward.ele.Score() < node.ele.Score() ||
				x.level[i].forward.ele.Score() == node.ele.Score() && node.order < x.level[i].forward.order) {
			x = x.level[i].forward
		}
		update[i] = x
	}

	x = x.level[0].forward
	if x != nil && x.ele.Score() == node.ele.Score() && x.ele.Key() == node.ele.Key() {
		if score > node.ele.Score() && (x.level[0].forward == nil || score < x.level[0].forward.ele.Score()) {
			node.ele.SetScore(score)
			return
		}

		// remove and re-insert when score changes
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

		node.ele.SetScore(score)
		list.insert(node.ele)
	}
}

// Find the rank for an element.
// Returns 0 when the element cannot be found, rank otherwise.
// Note that the rank is 1-based
func (list *zSkipList) zslGetRank(node *zSkipListNode) uint32 {
	var rank uint32
	x := list.header
	for i := list.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil &&
			(x.level[i].forward.ele.Score() < node.ele.Score() ||
				x.level[i].forward.ele.Score() == node.ele.Score() && node.order <= x.level[i].forward.order) {
			rank += x.level[i].span
			x = x.level[i].forward
		}
		if x.ele != nil && x.ele.Key() == node.ele.Key() {
			return rank
		}
	}
	return 0
}

func (list *zSkipList) randomLevel() int {
	lvl := 1
	for rand.Float64() < list.p && lvl < list.maxLevel {
		lvl++
	}
	return lvl
}

// Finds an element by its rank. The rank argument needs to be 1-based.
func (list *zSkipList) getElementByRank(rank uint32) *zSkipListNode {
	var traversed uint32
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

// ZSet set
type ZSet struct {
	dict map[uint64]*zSkipListNode
	zsl  *zSkipList
}

// NewZSet create ZSet
func NewZSet() *ZSet {
	zs := &ZSet{
		dict: make(map[uint64]*zSkipListNode),
		zsl:  zslCreate(),
	}
	return zs
}

// Add a new element or update the score of an existing element
func (zs *ZSet) Add(score uint32, key uint64) {
	if node := zs.dict[key]; node != nil {
		if score != node.ele.Score() {
			zs.zsl.updateScore(node, score)
		}
	} else {
		ele := &Element{
			key:   key,
			score: score,
		}
		node := zs.zsl.insert(ele)
		zs.dict[key] = node
	}
}

// Delete the element 'ele' from the sorted set,
// return 1 if the element existed and was deleted, 0 otherwise
func (zs *ZSet) Delete(id uint64) int {
	node := zs.dict[id]
	if node == nil {
		return 0
	}
	zs.zsl.delete(node)
	delete(zs.dict, id)
	return 1
}

// Rank return 1-based rank or 0 if not exist
func (zs *ZSet) Rank(id uint64, reverse bool) uint32 {
	node := zs.dict[id]
	if node != nil {
		rank := zs.zsl.zslGetRank(node)
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
func (zs *ZSet) Range(start uint32, end uint32, reverse bool) []*Element {
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
	rangeLen := end - start + 1
	if reverse {
		node := zs.zsl.getElementByRank(zs.zsl.length - end + 1)
		var ret = make([]*Element, rangeLen)
		for i := uint32(0); i < rangeLen; i++ {
			ret[rangeLen-i-1] = node.ele
			node = node.level[0].forward
		}
		return ret
	}
	node := zs.zsl.getElementByRank(start)
	var ret = make([]*Element, 0, rangeLen)
	for i := uint32(0); i < rangeLen; i++ {
		ret = append(ret, node.ele)
		node = node.level[0].forward
	}
	return ret
}

// Length return the element count
func (zs *ZSet) Length() uint32 {
	return zs.zsl.length
}

// MinScore 最小积分
func (zs *ZSet) MinScore() uint32 {
	if zs.zsl.header.level[0].forward != nil {
		return zs.zsl.header.level[0].forward.ele.Score()
	}
	return 0
}

// DeleteFirst 删除第一个元素
func (zs *ZSet) DeleteFirst() {
	zs.zsl.delete(zs.zsl.header.level[0].forward)
	delete(zs.dict, zs.zsl.header.level[0].forward.ele.Key())
}
