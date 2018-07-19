package badwords

import (
	"strings"

	"algorithms-go/bitarray"
)

// BadWords 屏蔽字
type BadWords struct {
	hashSet        map[string]bool // 脏字集合
	fastCharCheck  [0xffff]byte
	fastCharLength [0xffff]byte
	lastCharCheck  *bitarray.BitArray
	oneCharCheck   *bitarray.BitArray
	maxLength      int
	temp           []rune
}

// NewBadWords new
func NewBadWords() *BadWords {
	return &BadWords{
		hashSet:       make(map[string]bool),
		lastCharCheck: bitarray.NewBitArray(0xffff),
		oneCharCheck:  bitarray.NewBitArray(0xffff),
		maxLength:     0,
	}
}

// AddBadWord 增加屏蔽字
func (b *BadWords) AddBadWord(word string) {
	word = strings.ToLower(word)
	runeWord := []rune(word)
	if len(runeWord) == 1 {
		if b.oneCharCheck.Get(int(runeWord[0])) {
			return
		}
		b.oneCharCheck.Set(int(runeWord[0]), true)
		b.fastCharCheck[int(runeWord[0])] |= 1
	} else {
		if _, ok := b.hashSet[word]; ok {
			return
		}
		b.hashSet[word] = true

		for i, c := range runeWord {
			if i == len(runeWord)-1 {
				b.lastCharCheck.Set(int(c), true)
			} else if i < 7 {
				b.fastCharCheck[int(c)] |= 1 << uint32(i)
			} else {
				b.fastCharCheck[int(c)] |= 0x80
			}
		}
		m := uint32(len(runeWord) - 2)
		if m > 7 {
			m = 7
		}
		b.fastCharLength[runeWord[0]] |= 1 << m
	}

	if len(runeWord) > b.maxLength {
		b.maxLength = len(runeWord)
	}
}

func toLower(c rune) rune {
	if c >= 'A' && c <= 'Z' {
		return c + 32
	}
	return c
}

// ReplaceBadWord 替换屏蔽字为*
func (b *BadWords) ReplaceBadWord(text string, replaceChar rune) string {
	var runeText = []rune(text)
	var charCount = len(runeText)
	sub := make([]rune, 0, b.maxLength)
	for index := 0; index < charCount; index++ {
		firstChar := toLower(runeText[index])
		if b.fastCharCheck[int(firstChar)]&1 == 0 {
			continue
		}

		if b.oneCharCheck.Get(int(firstChar)) {
			runeText[index] = replaceChar
			continue
		}

		sub = sub[:0]
		sub = append(sub, firstChar)
		spaceCount := 0
		for j := 1; j < (b.maxLength+spaceCount) && j < charCount-index; j++ {
			currentChar := toLower(runeText[index+j])
			if b.isJumpChar(currentChar) {
				spaceCount++
				continue
			}

			m := uint32(j - spaceCount - 1)
			if m > 7 {
				m = 7
			}
			sub = append(sub, currentChar)
			if b.fastCharLength[firstChar]>>m&1 == 1 && b.lastCharCheck.Get(int(currentChar)) {
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

			k := uint32(j - spaceCount)
			if k > 7 {
				k = 7
			}
			if b.fastCharCheck[int(currentChar)]&(1<<k) == 0 {
				break
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
		firstChar := toLower(runeText[index])
		if b.fastCharCheck[int(firstChar)]&1 == 0 {
			continue
		}

		if b.oneCharCheck.Get(int(firstChar)) {
			return true
		}
		sub = sub[:0]
		sub = append(sub, firstChar)
		spaceCount := 0
		for j := 1; j < b.maxLength+spaceCount && j < charCount-index; j++ {
			currentChar := toLower(runeText[index+j])
			if b.isJumpChar(currentChar) {
				spaceCount++
				continue
			}

			m := uint32(j - spaceCount - 1)
			if m > 7 {
				m = 7
			}
			sub = append(sub, currentChar)
			if b.fastCharLength[firstChar]>>m&1 == 1 && b.lastCharCheck.Get(int(currentChar)) {
				if _, ok := b.hashSet[string(sub)]; ok {
					return true
				}
			}

			k := uint32(j - spaceCount)
			if k > 7 {
				k = 7
			}
			if b.fastCharCheck[int(currentChar)]&(1<<k) == 0 {
				break
			}
		}
	}
	return false
}

func (b *BadWords) isJumpChar(c rune) bool {
	return c == ' ' || c == '\t'
}
