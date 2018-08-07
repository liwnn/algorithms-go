package main

import (
	"fmt"
	"math/rand"
	"time"

	"algorithms-go/badwords"
	"algorithms-go/zset"
)

func main() {
	b := badwords.NewBadWords()
	b.AddBadWord("Sb")
	if b.ContainsBadWord("Sb") {
		fmt.Println("has")
	}
	c := b.ReplaceBadWord("sB", '*')
	fmt.Println(c)

	rand.Seed(time.Now().UnixNano())
	zs := zset.NewZSet()
	zs.Add(1, 10001)
	zs.Add(2, 10002)
	zs.Add(3, 10003)
	zs.Add(2, 10004)
	fmt.Println(zs.Rank(10001))
	zs.Delete(10001)
	fmt.Println(zs.Rank(10002))
}
