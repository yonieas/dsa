// Package sequence provides utility functions for working with iterators.
//
// These helpers work with Go 1.23+ iter.Seq sequences.
//
// Key concepts:
//   - Iterator: A function that yields elements one at a time
//   - Lazy evaluation: Elements are generated on demand
//   - Composable: Can chain operations
package sequence

import (
	"fmt"
	"iter"
	"strings"
)

// Enum wraps a sequence and adds indices, returning (index, value) pairs.
//
// Similar to Python's enumerate() or Rust's enumerate().
//
//	input:  [A] → [B] → [C]
//	output: (0,A) → (1,B) → (2,C)
//
// Example:
//
//	for i, v := range sequence.Enum(list.Iter) {
//	    fmt.Printf("%d: %v\n", i, v)
//	}
//
// complexity:
//   - time : O(1) per element
//   - space: O(1)
func Enum[E any](s iter.Seq[E]) iter.Seq2[int, E] {
	return func(yield func(int, E) bool) {
		var i int
		for v := range s {
			if !yield(i, v) {
				break
			}
			i++
		}
	}
}

// ValueAt retrieves the element at a specific index.
//
//	sequence: [A] → [B] → [C] → [D]
//	          0      1      2      3
//
//	ValueAt(seq, 2) → (C, true)
//	ValueAt(seq, 5) → (zero, false)
//
// Note: This requires iterating up to the index.
// For random access, use array-based structures.
//
// complexity:
//   - time : O(index)
//   - space: O(1)
func ValueAt[E any](s iter.Seq[E], index int) (E, bool) {
	for i, v := range Enum(s) {
		if i == index {
			return v, true
		}
	}
	var zeroValue E
	return zeroValue, false
}

// String converts a sequence to a bracketed string.
//
//	sequence: [1] → [2] → [3]
//
//	String(seq) → "[1 2 3]"
//
// complexity:
//   - time : O(n)
//   - space: O(n)
func String[E any](s iter.Seq[E]) string {
	return Format(s, " ")
}

// Format converts a sequence to a string with custom separator.
//
//	sequence: [1] → [2] → [3]
//
//	Format(seq, " ")   → "[1 2 3]"
//	Format(seq, ", ")  → "[1, 2, 3]"
//	Format(seq, "->")  → "[1->2->3]"
//
// complexity:
//   - time : O(n)
//   - space: O(n)
func Format[E any](s iter.Seq[E], sep string) string {
	var buf strings.Builder
	buf.WriteRune('[')
	for i, v := range Enum(s) {
		if i > 0 {
			buf.WriteString(sep)
		}
		if _, err := fmt.Fprint(&buf, v); err != nil {
			panic(fmt.Errorf("sequence.Format: write value: %w", err))
		}
	}
	buf.WriteRune(']')
	return buf.String()
}
