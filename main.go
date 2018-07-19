package main

import (
	"fmt"

	"algorithms-go/badwords"
)

func main() {
	b := badwords.NewBadWords()
	b.AddBadWord("SB")
	if b.ContainsBadWord("Sb") {
		fmt.Println("has")
	}
	c := b.ReplaceBadWord("sB", '*')
	fmt.Println(c)
}
