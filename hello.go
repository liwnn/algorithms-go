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
	p := z.p
	if z.left == nil && z.right == nil {
		if p.right == z {
			p.right = nil
		} else {
			p.left = nil
		}
	}

	if z.right != nil && z.left == nil {
		p.right = z.right
	}

	if z.left != nil && z.right != nil {
		x := treeMinimum(z.right)
		if x == z.right {
			x.p = z.p
			x.left = z.left

			if z.p.left == z {
				x.p.left = x
			} else {
				x.p.right = x
			}
		} else {
			x.right = z.right
			z.right.p = x

			z.left.p = x

			z.right.p = x

			if x == x.p.left {
				x.p.left = nil
			} else {
				x.p.right = nil
			}

			x.left = z.left
			z.left.p = x
			x.p = z.p

			if z.p.left == z {
				z.p.left = x
			} else {
				z.p.right = x
			}
		}
	}
}

func main() {
	var tree = &Tree{}
	treeInsert(tree, &Node{
		key: 7,
	})
	treeInsert(tree, &Node{
		key: 3,
	})

	treeInsert(tree, &Node{
		key: 2,
	})

	treeInsert(tree, &Node{
		key: 5,
	})

	treeInsert(tree, &Node{
		key: 4,
	})

	treeInsert(tree, &Node{
		key: 8,
	})

	treeInsert(tree, &Node{
		key: 9,
	})

	treeInsert(tree, &Node{
		key: 6,
	})

	node := treeSearch(tree.root, 7)
	_ = node
	treeDelete(tree, node)
	inorderTreeWalk(tree.root)
}
