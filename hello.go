package main

// Node 节点
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

func inorderTreeWalkRecursion(x *Node) {
	if x != nil {
		inorderTreeWalk(x.left)
		print(x.key)
		inorderTreeWalk(x.right)
	}
}

func inorderTreeWalk(x *Node) {
	var stack []*Node
	stack = append(stack, x)
	for x.left != nil {
		stack = append(stack, x.left)
		x = x.left
	}

	for len(stack) > 0 {
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		print(top.key)
		print(" ")
		if top.right != nil {
			p := top.right
			stack = append(stack, p)
			for p.left != nil {
				stack = append(stack, p.left)
				p = p.left
			}
		}
	}
}

func treeSearch(x *Node, k int32) *Node {
	var p = x
	for p != nil && p.key != k {
		if k < p.key {
			p = p.left
		} else {
			p = p.right
		}
	}
	return p
}

func treeSearchRecursion(x *Node, k int32) *Node {
	if x == nil || x.key == k {
		return x
	}

	if k < x.key {
		return treeSearch(x.left, k)
	}

	return treeSearch(x.right, k)
}

func treeMinimum(x *Node) *Node {
	for x.left != nil {
		x = x.left
	}
	return x
}

func treeMaximum(x *Node) *Node {
	for x.right != nil {
		x = x.right
	}
	return x
}

func treeSuccessor(x *Node) *Node {
	if x.right != nil {
		return treeMinimum(x.right)
	}
	y := x.p
	for y != nil && y.right == x {
		y = y.p
		x = y
	}
	return y
}

func treePreDecessor(x *Node) *Node {
	if x.left != nil {
		return treeMaximum(x.left)
	}
	y := x.p
	for y != nil && y.left == x {
		y = y.p
		x = y
	}
	return y
}

func treeInsert(t *Tree, z *Node) {
	var y *Node
	x := t.root
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

func treeDelete(t *Tree, z *Node) {
	if z == nil {
		return
	}
	p := z.p
	if z.left == nil && z.right == nil {
		if t.root == z {
			t.root = nil
		} else {
			if p.right == z {
				p.right = nil
			} else {
				p.left = nil
			}
		}
		return
	}

	if z.right != nil && z.left == nil {
		x := treeMinimum(z.right)
		if x.p.left == x {
			x.p.left = nil
		} else {
			x.p.right = nil
		}
		x.p = nil

		x.left = z.left
		if z.left != nil {
			z.left.p = x
		}

		x.right = z.right
		if z.right != nil {
			z.right.p = x
		}

		if t.root == z {
			t.root = x
		} else {
			if p.left == z {
				p.left = x
			} else {
				p.right = x
			}
			x.p = p
		}

	}

	if z.left != nil && z.right == nil {
		if t.root == z {
			t.root = z.left
			return
		}

		if p.right == z {
			p.right = z.left
		} else {
			p.left = z.left
		}
		z.left.p = p
		return
	}

	if z.left != nil && z.right != nil {
		x := treeMaximum(z.left)

		if x.p.left == x {
			x.p.left = nil
		} else {
			x.p.right = nil
		}
		x.p = nil

		x.left = z.left
		if z.left != nil {
			z.left.p = x
		}

		x.right = z.right
		if z.right != nil {
			z.right.p = x
		}

		if z == t.root {
			t.root = x
		} else {
			if z.p.left == z {
				z.p.left = x
			} else {
				z.p.right = x
			}
			x.p = z.p
		}
	}
}

func main() {
	var tree = &Tree{}

	var k = []int32{7, 3, 2, 5, 4, 6, 8, 10, 9}
	for _, v := range k {
		treeInsert(tree, &Node{
			key: v,
		})
	}

	node := treeSearch(tree.root, 6)
	_ = node
	treeDelete(tree, node)
	inorderTreeWalk(tree.root)
}
