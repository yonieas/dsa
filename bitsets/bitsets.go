package bitsets

import (
	"fmt"
	"math/bits"
	"strings"
)

type BitSet struct {
	bitfields []uint64
}

// New creates a new bitset that can store numOfBits.
// numOfBits must be a multiple of 8.
func New(numOfBits int) *BitSet {
	if numOfBits%8 != 0 {
		panic("BitSet: numOfBits must be a multiple of 8")
	}
	n := max(1, (numOfBits+63)/64) // divide into 64-bit chunks
	return &BitSet{
		bitfields: make([]uint64, n),
	}
}

func (b *BitSet) Add(pos int) {
	idx, offset := b.index(pos)
	b.bitfields[idx] |= 1 << offset
}

func (b *BitSet) Del(pos int) {
	idx, offset := b.index(pos)
	b.bitfields[idx] &^= 1 << offset
}

func (b *BitSet) Toggle(pos int) {
	idx, offset := b.index(pos)
	b.bitfields[idx] ^= 1 << offset
}

func (b *BitSet) Exists(pos int) bool {
	idx, offset := b.index(pos)
	return (b.bitfields[idx] & (1 << offset)) != 0
}

func (b *BitSet) Reset() {
	for i := range b.bitfields {
		b.bitfields[i] = 0
	}
}

func (b *BitSet) Len() int {
	return len(b.bitfields) * 64
}

func (b *BitSet) String() string {
	var sb strings.Builder
	for i := 0; i < b.Len(); i++ {
		if b.Exists(i) {
			sb.WriteByte('1')
		} else {
			sb.WriteByte('0')
		}
	}
	return sb.String()
}

func (b *BitSet) Count() int {
	count := 0
	for _, field := range b.bitfields {
		count += bits.OnesCount64(field)
	}
	return count
}

func (b *BitSet) index(pos int) (int, int) {
	if pos < 0 || pos >= b.Len() {
		panic(fmt.Sprintf("BitSet: index out of range [%d]", pos))
	}
	return pos / 64, pos % 64
}
