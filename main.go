package main

import (
	"fmt"

	"algorithms-go/badwords"
)

func main() {
	b := badwords.NewBadWords()
	b.AddBadWord("Sb")
	if b.ContainsBadWord("Sb") {
		fmt.Println("has")
	}
	c := b.ReplaceBadWord("sB", '*')
	fmt.Println(c)
}
