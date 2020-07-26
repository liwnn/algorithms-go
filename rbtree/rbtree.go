package main

import (
	"fmt"
)

type nodeColor int32

// enum
const (
	RED   nodeColor = 0
	BLACK nodeColor = 1
)

type node struct {
	color nodeColor
	key   int
	left  *node
	right *node
	p     *node
}

// RBTree is red-black tree
type RBTree struct {
	root *node
	nil  *node
}

func newRBTree() *RBTree {
	t := &RBTree{
		nil: &node{
			color: BLACK,
		},
	}
	t.root = t.nil
	return t
}

func (t *RBTree) insert(z *node) {
	x := t.root
	y := t.nil
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

func (t *RBTree) insertFixup(z *node) {
	for z.p.color == RED {
		break
	}
	//for z.p != t.nil && z.p.color == RED {
	//    if z.p.p.left == z.p {
	//        if z.p.p.right.color == RED {
	//            z.p.color = BLACK
	//            z.p.p.color = RED
	//            z.p.p.right.color = BLACK
	//            z = z.p.p
	//        } else {
	//            if z == z.p.right {
	//                z = z.p
	//                t.leftRotate(z)
	//            } else {
	//                z.p.color = BLACK
	//                z.p.p.color = RED
	//                t.rightRotate(z.p.p)
	//            }
	//        }
	//    } else {
	//        if z.p.p.left.color == RED {
	//            z.p.color = BLACK
	//            z.p.p.color = RED
	//            z.p.p.left.color = BLACK
	//            z = z.p.p
	//        } else {
	//            if z == z.p.right {
	//                z.p.color = BLACK
	//                z.p.p.color = RED
	//                t.leftRotate(z.p.p)
	//            } else {
	//                z = z.p
	//                t.rightRotate(z)
	//            }
	//        }
	//    }
	//}
	//t.root.color = BLACK
}

// x                 y
//   \      ->     /
//     y          x
func (t *RBTree) leftRotate(x *node) {
	y := x.right

	// y的左节点改成x的右节点
	x.right = y.left
	if y.left != t.nil {
		y.left.p = x
	}

	// x 改成y的左节点
	y.left = x
	if x.p == t.nil {
		t.root = y
	} else if x.p.left == x {
		x.p.left = y
	} else {
		x.p.right = y
	}
	y.p = x.p
	x.p = y
}

//    x        y
//   /    ->    \
//  y            x
func (t *RBTree) rightRotate(x *node) {
	y := x.left

	// y.right -> x.left
	x.left = y.right
	if y.right != t.nil {
		y.right.p = x.left
	}

	// x -> y.right
	y.right = x
	if x.p == t.nil {
		t.root = y
	} else if x.p.left == x {
		x.p.left = y
	} else {
		x.p.right = y
	}
	y.p = x.p
	x.p = y
}

func (t *RBTree) delete(z *node) {
	//y := z
	//var x *node
	//yOriginalColor := y.color
	//if z.left == t.nil {
	//    x = z.right
	//    t.transplant(z, z.right)
	//} else if z.right == t.nil {
	//    x = z.left
	//    t.transplant(z, z.left)
	//} else {
	//    y = t.minimum(z.right)
	//    yOriginalColor = y.color
	//    x := y.right
	//    if y.p == z { // y == z.right
	//        x.p = y
	//    } else {
	//        t.transplant(y, y.right)
	//        y.right = z.right
	//        y.right.p = y
	//    }
	//    t.transplant(z, y)
	//    y.left = z.left
	//    y.left.p = y
	//    y.color = z.color
	//}
	//
	//if yOriginalColor == BLACK {
	//    t.deleteFixup(x)
	//}
}

// v 替换u
func (t *RBTree) transplant(u *node, v *node) {
	if u.p == t.nil {
		t.root = v
	} else if u == u.p.left {
		u.p.left = v
	} else if u == u.p.right {
		u.p.right = v
	}
	v.p = u.p
}

func (t *RBTree) minimum(x *node) *node {
	for x.left != t.nil {
		x = x.left
	}
	return x
}

func (t *RBTree) deleteFixup(x *node) {
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

// PrintTree 打印树
func PrintTree(t *RBTree) {
	levelNode := make(map[int][]*node)
	levelNode[0] = []*node{t.root}
	for level := 0; ; level++ {
		var nodes = levelNode[level]
		var next []*node
		for _, n := range nodes {
			if n != nil {
				next = append(next, n.left, n.right)
			} else {
				next = append(next, nil, nil)
			}
		}
		var exit = true
		for _, v := range next {
			if v != nil {
				exit = false
				break
			}
		}
		if exit {
			break
		}
		levelNode[level+1] = next
	}
	depth := len(levelNode)
	for j := 0; j < depth; j++ {
		w := 1 << (depth - j + 1)
		if j > 0 {
			for i := 0; i < 1<<(j-1); i++ {
				fmt.Printf("%*c", w+1, ' ')
				for k := 0; k < w-3; k++ {
					fmt.Printf("_")
				}
				fmt.Printf("/ \\")
				for k := 0; k < w-3; k++ {
					fmt.Printf("_")
				}
				fmt.Printf("%*c", w+2, ' ')
			}

			fmt.Printf("\n")
			for i := 0; i < 1<<(j-1); i++ {
				fmt.Printf("%*c%*c%*c", w+1, '/', w*2-2, '\\', w+1, ' ')
			}
			fmt.Printf("\n")
		}
		for i := 0; i < 1<<j; i++ {
			n := levelNode[j][i]
			if n == nil {
				continue
			}
			if n != t.nil {
				fmt.Printf("%*c", w-2, ' ') // (key)
				if n.color == RED {
					fmt.Printf("%c[1;41;37m(%d)%c[0m", 0x1B, n.key, 0x1B)
				} else {
					fmt.Printf("%c[1;40;30m(%d)%c[0m", 0x1B, n.key, 0x1B)
				}
				fmt.Printf("%*c", w-1, ' ')
			} else {
				fmt.Printf("%*c", w-2, ' ') // (key)
				fmt.Printf("%c[1;40;30m%v%c[0m", 0x1B, "nil", 0x1B)
				fmt.Printf("%*c", w-1, ' ')
			}
		}
		fmt.Printf("\n")
	}
}

func main() {
	t := newRBTree()
	a := []int{11, 7}
	b := make([]*node, 0, len(a))
	for _, v := range a {
		n := &node{
			key: v,
		}
		b = append(b, n)
		t.insert(n)
	}
	PrintTree(t)
	//t.delete(b[5])
	//t.delete(b[3])
	//t.delete(b[4])
	//t.delete(b[2])
	//t.delete(b[1])
	//t.delete(b[0])
}
