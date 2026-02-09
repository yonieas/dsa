// Package bitsets provides a space-efficient set for storing integers.
//
// # What is a BitSet?
//
// A BitSet uses individual bits to represent set membership. Each bit position
// corresponds to an integer: if the bit is 1, that integer is in the set; if
// 0, it is not. This is incredibly compact since one byte stores 8 values.
//
// To store the numbers 0 through 63, a BitSet needs just 8 bytes (64 bits).
// A HashSet storing the same data would need hundreds of bytes due to
// pointers, buckets, and per-element overhead.
//
// # Why Use a BitSet?
//
// BitSets shine when you're working with dense sets of small integers. Set
// operations become blazingly fast because intersection is just bitwise AND,
// union is bitwise OR, and difference is AND NOT. Modern CPUs can process 64
// bits in a single instruction.
//
// Common use cases include: tracking which IDs have been seen, implementing
// Bloom filters, permission/feature flags, and finding common elements
// between datasets.
//
// # Limitations
//
// BitSets only work with non-negative integers in a bounded range (0 to
// numOfBits-1). For sparse data or large ranges, they waste space. If you
// need to store arbitrary values or negative numbers, use a HashSet instead.
//
// # Complexity
//
//	Add/Del/Exists: O(1)
//	Union/Intersect: O(n/64) where n is bit capacity
//	Count:          O(n/64) using popcount instructions
//
// # Further Reading
//
// Knuth "The Art of Computer Programming", Volume 4A, Bitwise Tricks.
// https://en.wikipedia.org/wiki/Bit_array
package bitsets

import (
	"fmt"
	"math/bits"
	"strings"
)

// BitSet represents a set of non-negative integers using a bit array.
//
// Example with 64 bits (numOfBits=64):
//
//	bit-position:  0  1  2  3  4  5  6  7 ... 63
//	bit value:     1  0  1  0  0  1  0  0 ... 0
//	               ↑     ↑        ↑
//	             in set
//
// This represents the set {0, 2, 5}.
//
// Internal structure:
//   - Uses uint64 words (64 bits each)
//   - Bit at position i is in: bitfields[i/64] at offset (i%64)
type BitSet struct {
	bitfields []uint64
}

// New creates a BitSet that can store integers in range [0, numOfBits-1].
//
//	numOfBits = 64:
//	┌────────────────────────────────────────────────────────────────┐
//	│ 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 ... 0 0 0 0 0 0 0 0 0 0 │
//	└────────────────────────────────────────────────────────────────┘
//	  0 1 2 3 4 5 6 7 8 9 ...                               ... 62 63
//
// complexity:
//   - time : O(numOfBits / 64)
//   - space: O(numOfBits / 64)
//
// Panics if numOfBits is not a multiple of 8.
func New(numOfBits int) *BitSet {
	if numOfBits%8 != 0 {
		panic("BitSet: numOfBits must be a multiple of 8")
	}
	n := max(1, (numOfBits+63)/64)
	return &BitSet{
		bitfields: make([]uint64, n),
	}
}

// Add inserts an integer into the set.
//
//	Before Add(5):
//	┌───┬───┬───┬───┬───┬───┬───┬───┐
//	│ 1 │ 0 │ 1 │ 0 │ 0 │ 0 │ 0 │ 0 │...
//	└───┴───┴───┴───┴───┴───┴───┴───┘
//	  0   1   2   3   4   5   6   7
//
//	After Add(5):
//	┌───┬───┬───┬───┬───┬───┬───┬───┐
//	│ 1 │ 0 │ 1 │ 0 │ 0 │ 1 │ 0 │ 0 │...
//	└───┴───┴───┴───┴───┴───┴───┴───┘
//	  0   1   2   3   4   5   6   7
//	                      ↑
//	                   added
//
// If already present, no change occurs.
//
// How it works:
//
//	bitfields[pos/64] |= (1 << (pos%64))
//
// complexity:
//   - time : O(1)
//   - space: O(1)
//
// Panics if pos < 0 or pos >= Len().
//
// SCORE: 10
func (b *BitSet) Add(pos int) {
	// hint: 1) idx, offset := b.index(pos)
	//       2) b.bitfields[idx] |= (1 << offset)
	idx, offset := b.index(pos)
	b.bitfields[idx] |= (1 << offset)
}

// Del removes an integer from the set.
//
//	Before Del(2):
//	┌───┬───┬───┬───┬───┬───┬───┬───┐
//	│ 1 │ 0 │ 1 │ 0 │ 0 │ 1 │ 0 │ 0 │...
//	└───┴───┴───┴───┴───┴───┴───┴───┘
//	  0   1   2   3   4   5   6   7
//
//	After Del(2):
//	┌───┬───┬───┬───┬───┬───┬───┬───┐
//	│ 1 │ 0 │ 0 │ 0 │ 0 │ 1 │ 0 │ 0 │...
//	└───┴───┴───┴───┴───┴───┴───┴───┘
//	  0   1   2   3   4   5   6   7
//	          ↑
//	       removed
//
// If not present, no change occurs.
//
// How it works:
//
//	bitfields[pos/64] &^= (1 << (pos%64))
//
// The &^ operator is Go's "bit clear" (AND NOT).
//
// complexity:
//   - time : O(1)
//   - space: O(1)
//
// Panics if pos < 0 or pos >= Len().
//
// SCORE: 10
func (b *BitSet) Del(pos int) {
	// hint: 1) idx, offset := b.index(pos)
	//       2) b.bitfields[idx] &^= (1 << offset)   // AND NOT (bit clear)
	idx, offset := b.index(pos)
	b.bitfields[idx] &^= (1 << offset)
}

// Toggle flips the membership of an integer.
//
//	Before Toggle(5):
//	┌───┬───┬───┬───┬───┬───┬───┬───┐
//	│ 1 │ 0 │ 0 │ 0 │ 0 │ 1 │ 0 │ 0 │...
//	└───┴───┴───┴───┴───┴───┴───┴───┘
//	  0   1   2   3   4   5   6   7
//
//	After Toggle(5):
//	┌───┬───┬───┬───┬───┬───┬───┬───┐
//	│ 1 │ 0 │ 0 │ 0 │ 0 │ 0 │ 0 │ 0 │...
//	└───┴───┴───┴───┴───┴───┴───┴───┘
//	  0   1   2   3   4   5   6   7
//	                      ↑
//	                  toggled (was 1, now 0)
//
// How it works:
//
//	bitfields[pos/64] ^= (1 << (pos%64))
//
// XOR with 1 flips a bit: 0→1 or 1→0.
//
// complexity:
//   - time : O(1)
//   - space: O(1)
//
// Panics if pos < 0 or pos >= Len().
//
// SCORE: 10
func (b *BitSet) Toggle(pos int) {
	// hint: 1) idx, offset := b.index(pos)
	//       2) b.bitfields[idx] ^= (1 << offset)   // XOR flips the bit
	idx, offset := b.index(pos)
	b.bitfields[idx] ^= (1 << offset)
}

// Exists checks whether an integer is in the set.
//
//	┌───┬───┬───┬───┬───┬───┬───┬───┐
//	│ 1 │ 0 │ 0 │ 0 │ 0 │ 1 │ 0 │ 0 │...
//	└───┴───┴───┴───┴───┴───┴───┴───┘
//	  0   1   2   3   4   5   6   7
//
//	Exists(0) → true
//	Exists(5) → true
//	Exists(2) → false
//
// How it works:
//
//	(bitfields[pos/64] & (1 << (pos%64))) != 0
//
// complexity:
//   - time : O(1)
//   - space: O(1)
//
// Panics if pos < 0 or pos >= Len().
//
// SCORE: 10
func (b *BitSet) Exists(pos int) bool {
	// hint: 1) idx, offset := b.index(pos)
	//       2) return (b.bitfields[idx] & (1 << offset)) != 0
	idx, offset := b.index(pos)
	return (b.bitfields[idx] & (1 << offset)) != 0
}

// Reset clears all bits, removing all elements.
//
//	Before Reset():
//	┌───┬───┬───┬───┬───┬───┬───┬───┐
//	│ 1 │ 0 │ 1 │ 0 │ 0 │ 1 │ 0 │ 0 │...
//	└───┴───┴───┴───┴───┴───┴───┴───┘
//
//	After Reset():
//	┌───┬───┬───┬───┬───┬───┬───┬───┐
//	│ 0 │ 0 │ 0 │ 0 │ 0 │ 0 │ 0 │ 0 │...
//	└───┴───┴───┴───┴───┴───┴───┴───┘
//
// complexity:
//   - time : O(Len() / 64)
//   - space: O(1)
func (b *BitSet) Reset() {
	for i := range b.bitfields {
		b.bitfields[i] = 0
	}
}

// Len returns the capacity (maximum value + 1 that can be stored).
//
// Valid values are in range [0, Len()-1].
//
// Note: This is NOT the count of elements. Use Count() for that.
//
// complexity:
//   - time : O(1)
//   - space: O(1)
func (b *BitSet) Len() int {
	return len(b.bitfields) * 64
}

// String returns a binary string representation.
//
//	Set = {0, 2, 5}
//	String() → "10100100..."
//	             ↑ ↑  ↑
//	             0 2  5
//
// complexity:
//   - time : O(Len())
//   - space: O(Len())
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

// Count returns the number of elements in the set (number of 1-bits).
//
//	┌───┬───┬───┬───┬───┬───┬───┬───┐
//	│ 1 │ 0 │ 1 │ 0 │ 0 │ 1 │ 0 │ 0 │...
//	└───┴───┴───┴───┴───┴───┴───┴───┘
//
//	Count() → 3 (positions 0, 2, 5 are set)
//
// This is also known as "population count" or "Hamming weight".
//
// complexity:
//   - time : O(Len() / 64)
//   - space: O(1)
//
// SCORE: 10
func (b *BitSet) Count() int {
	// hint: sum := 0; for each uint64 in b.bitfields:
	//       sum += bits.OnesCount64(field)  // popcount
	//       return sum
	sum := 0
	for _, field := range b.bitfields {
		sum += bits.OnesCount64(field)
	}
	return sum
}

// Size returns the number of elements in the set.
// This is an alias for Count() to satisfy adt.Sizer interface.
//
// complexity:
//   - time : O(Len() / 64)
//   - space: O(1)
func (b *BitSet) Size() int {
	return b.Count()
}

// Empty returns true if the set has no elements.
//
// complexity:
//   - time : O(Len() / 64)
//   - space: O(1)
func (b *BitSet) Empty() bool {
	return b.Count() == 0
}

// Iter iterates over all elements (set bit positions) in the set.
//
//	Set = {0, 2, 5}
//	for pos := range set.Iter {
//	    fmt.Println(pos)  // prints 0, 2, 5
//	}
//
// complexity:
//   - time : O(Len())
//   - space: O(1)
func (b *BitSet) Iter(yield func(int) bool) {
	for i := 0; i < b.Len(); i++ {
		if b.Exists(i) {
			if !yield(i) {
				return
			}
		}
	}
}

// Union returns a new BitSet with elements from both sets.
//
//	A = {0, 2, 5}     →  10100100...
//	B = {2, 3, 7}     →  00110001...
//	A ∪ B             →  10110101...
//	                     {0, 2, 3, 5, 7}
//
// Uses bitwise OR for O(n/64) performance.
//
// complexity:
//   - time : O(Len() / 64)
//   - space: O(Len() / 64)
//
// Panics if bitsets have different lengths.
//
// SCORE: 15
func (b *BitSet) Union(other *BitSet) *BitSet {
	// hint: 1) check same length, panic if different
	//       2) result := New(b.Len())
	//       3) for each i: result.bitfields[i] = b.bitfields[i] | other.bitfields[i]
	//       4) return result
	if b.Len() != other.Len() {
		panic("BitSets.Union: The length unmatch")
	}
	result := New(b.Len())
	for i, _ := range result.bitfields {
		result.bitfields[i] = b.bitfields[i] | other.bitfields[i]
	}
	return result
}

// Intersection returns a new BitSet with elements in both sets.
//
//	A = {0, 2, 5}     →  10100100...
//	B = {2, 3, 5}     →  00110100...
//	A ∩ B             →  00100100...
//	                     {2, 5}
//
// Uses bitwise AND for O(n/64) performance.
//
// complexity:
//   - time : O(Len() / 64)
//   - space: O(Len() / 64)
//
// Panics if bitsets have different lengths.
//
// SCORE: 20
func (b *BitSet) Intersection(other *BitSet) *BitSet {
	// hint: 1) check same length, panic if different
	//       2) result := New(b.Len())
	//       3) for each i: result.bitfields[i] = b.bitfields[i] & other.bitfields[i]
	//       4) return result
	if b.Len() != other.Len() {
		panic("BitSets.Intersection: The length unmatch")
	}
	result := New(b.Len())
	for i, _ := range result.bitfields {
		result.bitfields[i] = b.bitfields[i] & other.bitfields[i]
	}
	return result
}

// Disjoint returns true if the two sets have no common elements.
//
//	A = {0, 2, 5}     →  10100100...
//	B = {1, 3, 7}     →  01010001...
//	A ∩ B             →  00000000...  (empty)
//	Disjoint() → true
//
//	A = {0, 2, 5}     →  10100100...
//	C = {2, 3, 7}     →  00110001...
//	A ∩ C             →  00100000...  (has 2)
//	Disjoint() → false
//
// complexity:
//   - time : O(Len() / 64)
//   - space: O(1)
//
// Panics if bitsets have different lengths.
//
// SCORE: 15
func (b *BitSet) Disjoint(other *BitSet) bool {
	// hint: 1) check same length, panic if different
	//       2) for each i: if (b.bitfields[i] & other.bitfields[i]) != 0, return false
	//       3) return true (no common bits)
	if b.Len() != other.Len() {
		panic("BitSets.Disjoint: The length unmatch")
	}
	result := New(b.Len())
	for i, _ := range result.bitfields {
		if (b.bitfields[i] & other.bitfields[i]) != 0 {
			return false
		}
	}
	return true
}

func (b *BitSet) index(pos int) (int, int) {
	if pos < 0 || pos >= b.Len() {
		panic(fmt.Sprintf("BitSet: index out of range [%d]", pos))
	}
	return pos / 64, pos % 64
}
