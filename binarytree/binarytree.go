package binarytree

import (
	"fmt"
)

// Node node
type Node struct {
	key int32

	left  *Node
	right *Node
	p     *Node
}

// Tree tree
type Tree struct {
	root *Node
}

// Search search
func (t *Tree) Search(k int32) *Node {
	var x = t.root
	for x != nil && x.key != k {
		if k < x.key {
			x = x.left
		} else {
			x = x.right
		}
	}
	return x
}

// Minimum return min node.
func (t *Tree) Minimum(x *Node) *Node {
	for x.left != nil {
		x = x.left
	}
	return x
}

// Maximum return max node.
func (t *Tree) Maximum(x *Node) *Node {
	for x.right != nil {
		x = x.right
	}
	return x
}

// Successor return successor.
func (t *Tree) Successor(x *Node) *Node {
	if x.right != nil {
		return t.Minimum(x.right)
	}

	y := x.p
	for y != nil && y.right == x {
		x = y
		y = x.p
	}
	return y
}

// Predecessor return predecessor
func (t *Tree) Predecessor(x *Node) *Node {
	if x.left != nil {
		return t.Maximum(x.left)
	}

	y := x.p
	for y != nil && y.left == x {
		x = y
		y = x.p
	}
	return y
}

// Insert node z
func (t *Tree) Insert(z *Node) {
	x := t.root
	var y *Node
	for x != nil {
		y = x
		if z.key < x.key {
			x = x.left
		} else {
			x = x.right
		}
	}

	z.p = y

	if y == nil {
		t.root = z
	} else if z.key < y.key {
		y.left = z
	} else {
		y.right = z
	}
}

// Transplant replace u with v.
func (t *Tree) Transplant(u, v *Node) {
	if u.p == nil {
		t.root = v
	} else if u.p.left == u {
		u.p.left = v
	} else {
		u.p.right = v
	}

	if v != nil {
		v.p = u.p
	}
}

// Delete node z.
func (t *Tree) Delete(z *Node) {
	if z.left == nil { // 左节点为空，右节点存在或空
		t.Transplant(z, z.right)
	} else if z.right == nil { // 右节点为空，左节点不空
		t.Transplant(z, z.left)
	} else {
		y := t.Minimum(z.right)
		if y.p != z { // 后继y不为z的右子节点
			t.Transplant(y, y.right) // 用y的右孩子替换y
			y.right = z.right        // y 指向z的右孩子，这样之后就跟y是z的右孩子情况一致
			y.right.p = y
		}
		t.Transplant(z, y) // y是z的有孩子，直接y替换z
		y.left = z.left    // 左孩子设置
		y.left.p = y
	}
}

// PrintTree 打印树
func PrintTree(t *Tree) {
	levelNode := make(map[int][]*Node)
	levelNode[0] = []*Node{t.root}
	for level := 0; ; level++ {
		var nodes = levelNode[level]
		var next []*Node
		for _, node := range nodes {
			if node != nil {
				next = append(next, node.left, node.right)
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
			node := levelNode[j][i]
			if node != nil {
				fmt.Printf("%*c(%d)%*c", w-2, ' ', node.key, w-1, ' ') // (key)
			} else {
				fmt.Printf("%*c%cl%*c", w-1, 'n', 'i', w-1, ' ') // (key)
			}
		}
		fmt.Printf("\n")
	}
}
