package zset_test

import (
	"fmt"
	"math/rand"
	"time"

	"algorithms-go/zset"
)

func Example() {
	rand.Seed(time.Now().UnixNano())
	zs := zset.NewZSet()
	zs.Add(1, 10001)
	zs.Add(2, 10002)
	zs.Add(3, 10003)
	zs.Add(2, 10004)
	fmt.Println(zs.Rank(10002))
	fmt.Println(zs.Rank(10003))
	zs.Delete(10001)
	fmt.Println(zs.Rank(10002))
	// Output:
	// 2
	// 3
	// 1
}
