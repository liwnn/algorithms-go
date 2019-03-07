package quicksort

import (
	"fmt"
	"testing"
)

func Test_sort(t *testing.T) {
	var a = []int{3, 1, 8, 9, 0, 2, 2, 6, 100, 123, 23, 4}
	// var a = []int{1, 2, 3, 4, 5}
	quickSort(a, 0, len(a)-1)
	fmt.Println(a)
}

func Benchmark_sort(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var a = []int{3, 1, 8, 9, 0, 2, 6, 100, 123, 23, 4}
		// var a = []int{1, 2, 3, 4, 5}
		quickSort(a, 0, len(a)-1)
	}
}
