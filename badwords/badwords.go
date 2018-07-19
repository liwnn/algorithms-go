package badwords

import (
	"algorithms-go/bitarray"
)

// BadWords 屏蔽字
type BadWords struct {
	hashSet        map[string]bool // 脏字集合
	firstCharCheck *bitarray.BitArray
	allCharCheck   *bitarray.BitArray
	lastCharCheck  *bitarray.BitArray
	maxLength      int
}

// NewBadWords new
func NewBadWords() *BadWords {
	return &BadWords{
		hashSet:        make(map[string]bool),
		firstCharCheck: bitarray.NewBitArray(0xffff),
		allCharCheck:   bitarray.NewBitArray(0xffff),
		lastCharCheck:  bitarray.NewBitArray(0xffff),
	}
}

// AddBadWord 增加屏蔽字
func (b *BadWords) AddBadWord(word string) {
	_, ok := b.hashSet[word]
	if ok {
		return
	}
	b.hashSet[word] = true

	runeWord := []rune(word)
	if len(runeWord) > b.maxLength {
		b.maxLength = len(runeWord)
	}
	for i, c := range runeWord {
		set := false
		if i == 0 {
			b.firstCharCheck.Set(int(c), true)
			set = true
		}

		if i == len(runeWord)-1 {
			b.lastCharCheck.Set(int(c), true)
			set = true
		}

		if set == false {
			b.allCharCheck.Set(int(c), true)
		}
	}
}

// ReplaceBadWord 替换屏蔽字为*
func (b *BadWords) ReplaceBadWord(text string, replaceChar rune) string {
	var runeText = []rune(text)
	var charCount = len(runeText)
	sub := make([]rune, 0, b.maxLength)
	for index := 0; index < charCount; index++ {
		if !b.firstCharCheck.Get(int(runeText[index])) {
			continue
		}

		sub = sub[:0]
		spaceCount := 0
		for j := 0; j < (b.maxLength+spaceCount) && j < charCount-index; j++ {
			if j > 0 {
				if b.isJumpChar(runeText[index+j]) {
					spaceCount++
					continue
				}
			}

			sub = append(sub, runeText[index+j])
			if b.lastCharCheck.Get(int(runeText[index+j])) {
				if _, ok := b.hashSet[string(sub)]; ok {
					for i := index; i <= index+j; i++ {
						if !(b.isJumpChar(runeText[i])) {
							runeText[i] = replaceChar
						}
					}
					index += j
					break
				}
			}

			if j > 0 {
				if !b.allCharCheck.Get(int(runeText[index+j])) {
					break
				}
			}
		}
	}
	return string(runeText)
}

// ContainsBadWord 是否含有屏蔽字
func (b *BadWords) ContainsBadWord(text string) bool {
	var runeText = []rune(text)
	var charCount = len(runeText)
	sub := make([]rune, 0, b.maxLength)
	for index := 0; index < charCount; index++ {
		if !b.firstCharCheck.Get(int(runeText[index])) {
			continue
		}

		sub = sub[:0]
		spaceCount := 0
		for j := 0; j < b.maxLength+spaceCount && j < charCount-index; j++ {
			if j > 0 {
				if b.isJumpChar(runeText[index+j]) {
					spaceCount++
					continue
				}
			}

			sub = append(sub, runeText[index+j])
			if b.lastCharCheck.Get(int(runeText[index+j])) {
				if _, ok := b.hashSet[string(sub)]; ok {
					return true
				}
			}

			if j > 0 {
				if !b.allCharCheck.Get(int(runeText[index+j])) {
					break
				}
			}
		}
	}
	return false
}

func (b *BadWords) isJumpChar(c rune) bool {
	return c == ' ' || c == '\t'
}
