package badwords_test

import (
	"fmt"
	"testing"

	"algorithms-go/badwords"
)

func TestBadWords(t *testing.T) {
	b := badwords.NewBadWords()
	b.AddBadWord("sb")
	b.AddBadWord("bt")
	c := b.ReplaceBadWord("sbbbtttb", '*')
	fmt.Println(c)
}
