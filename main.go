package main

import (
	"fmt"

	"algorithms-go/badwords"
)

func main() {
	b := badwords.NewBadWords()
	b.AddBadWord("s")
	b.AddBadWord("bt")
	c := b.ReplaceBadWord("s bsbb tstb", '*')
	fmt.Println(c)
}
