package quicksort

func partion(a []int, p, r int) int {
	x := a[r]
	j := p
	for i := p; i < r; i++ {
		if a[i] < x {
			a[i], a[j] = a[j], a[i]
			j++
		}
	}
	a[j], a[r] = a[r], a[j]
	return j
}

func quickSort(a []int, p, r int) {
	if p < r {
		q := partion(a, p, r)
		quickSort(a, p, q-1)
		quickSort(a, q+1, r)
	}
}
