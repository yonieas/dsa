// Package hashmap provides hash-based key-value storage implementations.
//
// A HashMap (also called Hash Table or Dictionary) stores key-value pairs
// and provides fast lookup, insertion, and deletion operations.
//
// How it works:
//  1. Compute hash(key) to get an integer
//  2. Use modulo to find bucket index: hash(key) % num_buckets
//  3. Store/retrieve the key-value pair in that bucket
//
// Collision handling:
// When two keys hash to the same bucket, this implementation uses
// "separate chaining" - each bucket holds a linked list of entries.
//
//	bucket 0: → [key1:val1] → [key5:val5] → nil
//	bucket 1: → [key2:val2] → nil
//	bucket 2: → nil
//	bucket 3: → [key3:val3] → [key4:val4] → nil
package hashmap

import (
	"fmt"
	"hash/fnv"
)

// Default configuration values.
const (
	// DefaultCapacity is the initial number of buckets.
	DefaultCapacity = 16

	// DefaultLoadThreshold triggers resize when exceeded.
	// When size/capacity > 0.75, the map doubles in size.
	DefaultLoadThreshold = 0.75
)

// Entry is a key-value pair stored in the HashMap.
//
//	┌─────────────────────────┐
//	│  key   │     value      │
//	│ "alice"│      30        │
//	└─────────────────────────┘
type Entry[K comparable, V any] struct {
	key K
	val V
}

// Key returns the key of this entry.
func (e Entry[K, V]) Key() K { return e.key }

// Value returns the value of this entry.
func (e Entry[K, V]) Value() V { return e.val }

// NewEntry creates a new Entry with the given key and value.
func NewEntry[K comparable, V any](key K, val V) *Entry[K, V] {
	return &Entry[K, V]{
		key: key,
		val: val,
	}
}

// DefaultHashFunction computes a hash value for any comparable key.
//
// How it works:
//  1. Convert key to string representation
//  2. Apply FNV-1a hash algorithm
//  3. Return as integer
//
// Example:
//
//	DefaultHashFunction("hello") → 1335831723
//	bucketIndex = 1335831723 % 16 = 11
func DefaultHashFunction[K comparable](key K) int {
	h := fnv.New32()
	_, _ = fmt.Fprint(h, key)
	return int(h.Sum32())
}
