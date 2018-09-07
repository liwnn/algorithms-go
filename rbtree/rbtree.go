package main

import "fmt"

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
}

func (t *rbTree) print(n *node) {
	if n != t.nil {
		t.print(n.left)
		t.print(n.right)
		color := "black"
		if n.color == RED {
			color = "red"
		}
		fmt.Printf("%d(%s)\n", n.key, color)
	}
}

func main() {
	t := newRbTree()
	a := []int{1, 0, 3, 4}
	x := t.root
	for _, v := range a {
		x = &node{
			key: v,
		}
		t.insert(x)
	}
	t.delete(x)
	t.print(t.root)
}
