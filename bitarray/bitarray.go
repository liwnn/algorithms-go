package bitarray

// BitArray manages a compact array of bit values, which are represented as bool,
// where true indicates that the bit is on (1) and false indicates the bit is off (0).
type BitArray struct {
	values []byte
}

// New initializes BitArray b.
func New(length uint32) *BitArray {
	size := length >> 3
	if length&7 > 0 {
		size++
	}
	b := &BitArray{
		values: make([]byte, size),
	}
	return b
}

// Set index to 1 if value is true, or set index to 0 if value is false.
func (b *BitArray) Set(index uint32, value bool) {
	if index >= uint32(len(b.values))<<3 {
		if !value {
			return
		}
		length := uint64(index) + 1
		size := length >> 3
		if length&7 > 0 {
			size++
		}

		if b.values == nil {
			b.values = make([]byte, size)
		} else if int(size) <= cap(b.values) {
			b.values = b.values[:size]
		} else {
			capacity := size
			if size >= 1024 {
				capacity += size / 4
			} else {
				capacity += size // double cap
			}
			v := make([]byte, size, capacity)
			copy(v, b.values)
			b.values = v
		}
	}

	if value {
		b.values[index>>3] |= 1 << (index & 7)
	} else {
		b.values[index>>3] &^= (1 << (index & 7))
	}
}

// Get true if index is set 1, or return false.
func (b *BitArray) Get(index uint32) bool {
	if uint64(index) >= uint64(len(b.values))<<3 {
		return false
	}
	return (b.values[index>>3] & (1 << (index & 7))) != 0
}
