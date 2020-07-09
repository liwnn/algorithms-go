package binarytree

import "testing"

func Test_Print(t *testing.T) {
	var tree = &Tree{}
	var k = []int32{7, 3, 2, 5, 4, 6, 8, 10, 9}
	for _, v := range k {
		tree.Insert(&Node{
			key: v,
		})
	}
	PrintTree(tree)
}
