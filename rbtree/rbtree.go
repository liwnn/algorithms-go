package main

import (
	"container/list"
	"fmt"
)

// enum
const (
	RED   = 0
	BLACK = 1
)

type node struct {
	color int
	key   int
	left  *node
	right *node
	p     *node
}

type rbTree struct {
	root *node
	nil  *node
}

func newRbTree() *rbTree {
	t := &rbTree{
		nil: &node{
			color: BLACK,
		},
	}
	t.root = t.nil
	return t
}

func (t *rbTree) insert(z *node) {
	y := t.nil
	x := t.root
	for x != t.nil {
		y = x
		if z.key < x.key {
			x = x.left
		} else {
			x = x.right
		}
	}

	if y == t.nil {
		t.root = z
	} else if z.key < y.key {
		y.left = z
	} else {
		y.right = z
	}
	z.p = y
	z.left = t.nil
	z.right = t.nil
	z.color = RED

	t.insertFixup(z)
}

func (t *rbTree) insertFixup(z *node) {
	for z.p != t.nil && z.p.color == RED {
		if z.p.p.left == z.p {
			if z.p.p.right.color == RED {
				z.p.color = BLACK
				z.p.p.color = RED
				z.p.p.right.color = BLACK
				z = z.p.p
			} else {
				if z == z.p.right {
					z = z.p
					t.leftRotate(z)
				} else {
					z.p.color = BLACK
					z.p.p.color = RED
					t.rightRotate(z.p.p)
				}
			}
		} else {
			if z.p.p.left.color == RED {
				z.p.color = BLACK
				z.p.p.color = RED
				z.p.p.left.color = BLACK
				z = z.p.p
			} else {
				if z == z.p.right {
					z.p.color = BLACK
					z.p.p.color = RED
					t.leftRotate(z.p.p)
				} else {
					z = z.p
					t.rightRotate(z)
				}
			}
		}
	}
	t.root.color = BLACK
}

func (t *rbTree) leftRotate(x *node) {
	y := x.right
	x.right = y.left
	y.left.p = x
	y.p = x.p
	if x == t.root {
		t.root = y
	} else if x == x.p.left {
		x.p.left = y
	} else if x == x.p.right {
		x.p.right = y
	}
	y.left = x
	x.p = y
}

func (t *rbTree) rightRotate(x *node) {
	y := x.left
	x.left = y.right
	y.right.p = x
	if x == t.root {
		t.root = y
	} else if x.p.left == x {
		x.p.left = y
	} else if x.p.right == x {
		x.p.right = y
	}
	y.p = x.p
	y.right = x
	x.p = y
}

func (t *rbTree) delete(z *node) {
	y := z
	var x *node
	yOriginalColor := y.color
	if z.left == t.nil {
		x = z.right
		t.transplant(z, z.right)
	} else if z.right == t.nil {
		x = z.left
		t.transplant(z, z.left)
	} else {
		y = t.minimum(z.right)
		yOriginalColor = y.color
		x := y.right
		if y.p == z { // y == z.right
			x.p = y
		} else {
			t.transplant(y, y.right)
			y.right = z.right
			y.right.p = y
		}
		t.transplant(z, y)
		y.left = z.left
		y.left.p = y
		y.color = z.color
	}

	if yOriginalColor == BLACK {
		t.deleteFixup(x)
	}
}

// v 替换u
func (t *rbTree) transplant(u *node, v *node) {
	if u.p == t.nil {
		t.root = v
	} else if u == u.p.left {
		u.p.left = v
	} else if u == u.p.right {
		u.p.right = v
	}
	v.p = u.p
}

func (t *rbTree) minimum(x *node) *node {
	for x.left != t.nil {
		x = x.left
	}
	return x
}

func (t *rbTree) deleteFixup(x *node) {
	for x != t.root && x.color == BLACK {
		if x == x.p.left {
			w := x.p.right
			if w.color == RED {
				w.color = BLACK
				x.p.color = RED
				t.leftRotate(x.p)
				w = x.p.right
			} else {
				if w.left.color == BLACK && w.right.color == BLACK {
					w.color = RED
					x = x.p
				} else if w.right.color == BLACK {
					w.left.color = BLACK
					w.color = RED
					t.rightRotate(w)
					w = x.p.right
				} else {
					w.color = x.p.color
					x.p.color = BLACK
					w.right.color = BLACK
					t.leftRotate(x.p)
					x = t.root
				}
			}
		} else {
			w := x.p.left
			if w.color == RED {
				w.color = BLACK
				x.p.color = RED
				t.rightRotate(x.p)
				w = x.p.left
			} else {
				if w.left.color == BLACK && w.right.color == BLACK {
					w.color = RED
					x = x.p
				} else if w.left.color == BLACK {
					w.right.color = BLACK
					w.color = RED
					t.leftRotate(w)
					w = x.p.left
				} else {
					w.color = x.p.color
					x.p.color = BLACK
					w.left.color = BLACK
					t.rightRotate(x.p)
					x = t.root
				}
			}
		}
	}
	x.color = BLACK
}

func main() {
	t := newRbTree()
	a := []int{41, 38, 31, 12, 19, 8}
	b := []*node{nil, nil, nil, nil, nil, nil}
	for i, v := range a {
		b[i] = &node{
			key: v,
		}
		t.insert(b[i])
	}
	print(t)
	t.delete(b[5])
	t.delete(b[3])
	t.delete(b[4])
	t.delete(b[2])
	t.delete(b[1])
	t.delete(b[0])
}

func print(t *rbTree) {
	type nodeEx struct {
		*node
		lvl  int
		text string
	}
	nodeList := list.New()
	nodeList.PushBack(&nodeEx{t.root, 0, "root"})
	k := 0
	for nodeList.Len() > 0 {
		e := nodeList.Front()
		n := e.Value.(*nodeEx)
		nodeList.Remove(e)
		color := "black"
		if n.color == RED {
			color = "red"
		}
		if n.lvl != k {
			k = n.lvl
			fmt.Println()
		}
		fmt.Printf("%d(%s:%d:%s) ", n.key, n.text, n.lvl, color)
		if n.left != t.nil {
			nodeList.PushBack(&nodeEx{n.left, n.lvl + 1, "left"})
		}
		if n.right != t.nil {
			nodeList.PushBack(&nodeEx{n.right, n.lvl + 1, "right"})
		}
	}
}
