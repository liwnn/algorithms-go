package badwords

import (
	"strconv"
	"testing"
)

func BenchmarkAdd(b *testing.B) {
	bw := NewBadWords()
	for i := 0; i < b.N; i++ {
		bw.AddBadWord("a" + strconv.Itoa(i))
	}
}

func BenchmarkReplace(b *testing.B) {
	bw := NewBadWords()
	for i := 0; i < 100000; i++ {
		bw.AddBadWord("a" + strconv.Itoa(i))
	}
	str := "a1balkjd;flajdl;fjka;lskdfjlasjkdfljkljkalsdjkflsajkdfaa1231231ljljksdlajkdfkj1kjljakl1jl3j4ljk2l34jaljl12k4j31l2jk341l2jk43l;12jk43lj"
	c := '*'
	for i := 0; i < b.N; i++ {
		bw.ReplaceBadWord(str, c)
	}
}

func BenchmarkContains(b *testing.B) {
	bw := NewBadWords()
	for i := 0; i < 100000; i++ {
		bw.AddBadWord("a" + strconv.Itoa(i))
	}
	str := "a1balkjd;flajdl;fjka;lskdfjlasjkdfljkljkalsdjkflsajkdfaa1231231ljljksdlajkdfkj1kjljakl1jl3j4ljk2l34jaljl12k4j31l2jk341l2jk43l;12jk43lj"
	for i := 0; i < b.N; i++ {
		bw.ContainsBadWord(str)
	}
}
