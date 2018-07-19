package main

import (
	"fmt"

	"algorithms-go/badwords"
)

func main() {
	b := badwords.NewBadWords()
	b.AddBadWord("sb")
	b.AddBadWord("bt")
	c := b.ReplaceBadWord("s bb btstb", '*')
	fmt.Println(c)
}
