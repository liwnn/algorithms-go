# bitarray
A simple bit array implement in go

## Example
``` go
package main

import (
	"fmt"

	"github.com/liwnn/bitarray"
)

func main() {
	b := bitarray.New(8)
	b.Set(1, true)
	if b.Get(1) {
		fmt.Println("1 is set!")
	}
	b.Set(1, false)
}
```

## Installation
```
    go get github.com/liwnn/bitarray
```
