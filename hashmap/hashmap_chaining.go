// Package hashmap provides a HashMap implementation using separate chaining.
//
// # What is a HashMap?
//
// A HashMap (also called hash table or dictionary) stores key-value pairs
// with O(1) average time for get, put, and delete operations. The magic is
// the hash function: it converts any key into an array index, giving direct
// access without searching.
//
// # How It Works
//
// When you insert a key-value pair:
//
//  1. Compute hash(key) to get an integer
//  2. Map that integer to a bucket index: index = hash % numBuckets
//  3. Store the pair in that bucket
//
// Retrieval reverses the process: hash the key, find the bucket, return the
// value.
//
// # Handling Collisions
//
// Different keys can hash to the same bucket (collision). This implementation
// uses separate chaining: each bucket holds a linked list of all entries that
// hash to that index. On collision, we just append to the list.
//
// Alternative strategies include open addressing (probing for next empty slot)
// and Robin Hood hashing (minimize probe lengths).
//
// # Load Factor
//
// As the table fills up, collisions increase and performance degrades. The
// load factor (size / numBuckets) measures this. When it exceeds a threshold
// (typically 0.75), we resize: allocate a larger bucket array and rehash all
// entries.
//
// # Complexity
//
//	Get/Put/Delete: O(1) average, O(n) worst case (all keys collide)
//	Space:          O(n)
//
// The O(n) worst case is rare with a good hash function. In practice, hash
// tables are the go-to structure for fast lookups.
//
// # Further Reading
//
// CLRS "Introduction to Algorithms", Chapter 11.
// Sedgewick "Algorithms", Section 3.4.
// https://en.wikipedia.org/wiki/Hash_table
package hashmap

import (
	"cmp"
	"fmt"
	"strings"

	"github.com/josestg/dsa/internal/generics"
	"github.com/josestg/dsa/linkedlist"
	"github.com/josestg/dsa/sequence"
)

// HashMap is a hash table using separate chaining for collision resolution.
//
//	                  buckets array
//	            ┌─────────────────────────┐
//	index 0 ──► │ [a:1] → [e:5] → nil     │
//	index 1 ──► │ nil                     │
//	index 2 ──► │ [b:2] → nil             │
//	index 3 ──► │ [c:3] → [d:4] → nil     │
//	            └─────────────────────────┘
//
// Each bucket is a linked list of entries that hash to the same index.
type HashMap[K comparable, V any] struct {
	size          int
	buckets       []*linkedlist.SinglyLinkedList[*Entry[K, V]]
	hashFunction  func(K) int
	loadThreshold float64
}

// Options configures HashMap behavior.
type Options[K comparable] struct {
	Capacity      int
	LoadThreshold float64
	HashFunction  func(key K) int
}

// New creates a HashMap with default settings.
//
//	capacity = 16, loadThreshold = 0.75
//
//	┌───┬───┬───┬───┬───┬───┬───┬───┬───┬───┬───┬───┬───┬───┬───┬───┐
//	│ 0 │ 1 │ 2 │ 3 │ 4 │ 5 │ 6 │ 7 │ 8 │ 9 │10 │11 │12 │13 │14 │15 │
//	├───┼───┼───┼───┼───┼───┼───┼───┼───┼───┼───┼───┼───┼───┼───┼───┤
//	│nil│nil│nil│nil│nil│nil│nil│nil│nil│nil│nil│nil│nil│nil│nil│nil│
//	└───┴───┴───┴───┴───┴───┴───┴───┴───┴───┴───┴───┴───┴───┴───┴───┘
//
// complexity:
//   - time : O(capacity)
//   - space: O(capacity)
func New[K comparable, V any]() *HashMap[K, V] {
	return NewWith[K, V](Options[K]{
		Capacity:      DefaultCapacity,
		LoadThreshold: DefaultLoadThreshold,
		HashFunction:  DefaultHashFunction[K],
	})
}

// NewWith creates a HashMap with custom configuration.
//
// Panics if LoadThreshold is not in range (0, 1).
func NewWith[K comparable, V any](opts Options[K]) *HashMap[K, V] {
	capacity := cmp.Or(opts.Capacity, DefaultCapacity)
	loadFactor := cmp.Or(opts.LoadThreshold, DefaultLoadThreshold)
	if loadFactor <= 0 || loadFactor >= 1 {
		panic("HashMap.NewWith: load factor must be in range (0,1) exclusive")
	}
	hashFunction := opts.HashFunction
	if hashFunction == nil {
		hashFunction = DefaultHashFunction
	}

	buckets := make([]*linkedlist.SinglyLinkedList[*Entry[K, V]], capacity)
	for i := range buckets {
		buckets[i] = linkedlist.NewSinglyLinkedList[*Entry[K, V]]()
	}

	return &HashMap[K, V]{
		size:          0,
		buckets:       buckets,
		loadThreshold: loadFactor,
		hashFunction:  hashFunction,
	}
}

// Put inserts or updates a key-value pair.
//
//	Before Put("c", 3):
//
//	bucket 0 ──► [a:1] → nil
//	bucket 1 ──► nil
//	bucket 2 ──► [b:2] → nil
//	bucket 3 ──► nil
//
//	hash("c") % 4 = 2  (same bucket as "b")
//
//	After Put("c", 3):
//
//	bucket 0 ──► [a:1] → nil
//	bucket 1 ──► nil
//	bucket 2 ──► [b:2] → [c:3] → nil  ← chained!
//	bucket 3 ──► nil
//
// If key already exists, its value is updated.
// If load factor exceeds threshold, the map resizes automatically.
//
// complexity:
//   - time : O(1) average, O(n) worst case
//   - space: O(1)
func (h *HashMap[K, V]) Put(key K, value V) {
	if h.LoadFactor() >= h.loadThreshold {
		h.growAndRehash()
	}
	h.put(key, value)
}

// Del removes a key-value pair from the map.
//
//	Before Del("b"):
//
//	bucket 0 ──► [a:1] → nil
//	bucket 1 ──► nil
//	bucket 2 ──► [b:2] → [c:3] → nil
//	bucket 3 ──► nil
//
//	After Del("b"):
//
//	bucket 0 ──► [a:1] → nil
//	bucket 1 ──► nil
//	bucket 2 ──► [c:3] → nil  ← "b" removed
//	bucket 3 ──► nil
//
// If key doesn't exist, no action is taken.
//
// complexity:
//   - time : O(1) average, O(n) worst case
//   - space: O(1)
func (h *HashMap[K, V]) Del(key K) {
	index := h.bucketIndex(key)
	entries := h.buckets[index]
	for i, v := range sequence.Enum(entries.Iter) {
		if v.key == key {
			_ = entries.Remove(i)
			h.size--
			break
		}
	}
}

func (h *HashMap[K, V]) put(key K, value V) {
	index := h.bucketIndex(key)
	entries := h.buckets[index]

	for v := range entries.Iter {
		if v.key == key {
			v.val = value
			return
		}
	}

	entries.Append(NewEntry(key, value))
	h.size++
}

// Iter iterates over all entries in the map.
//
//	bucket 0 ──► [a:1] → nil
//	bucket 1 ──► nil
//	bucket 2 ──► [b:2] → [c:3] → nil
//
//	for entry := range m.Iter {
//	    // yields: [a:1], [b:2], [c:3]
//	}
//
// Iteration order is not guaranteed.
//
// complexity:
//   - time : O(n + capacity)
//   - space: O(1)
func (h *HashMap[K, V]) Iter(yield func(*Entry[K, V]) bool) {
	for _, entries := range h.buckets {
		for v := range entries.Iter {
			if !yield(v) {
				return
			}
		}
	}
}

// Keys iterates over all keys in the map.
//
//	for key := range m.Keys {
//	    // yields: "a", "b", "c"
//	}
//
// complexity:
//   - time : O(n + capacity)
//   - space: O(1)
func (h *HashMap[K, V]) Keys(yield func(K) bool) {
	for e := range h.Iter {
		if !yield(e.Key()) {
			break
		}
	}
}

// Get retrieves the value for a key.
//
//	bucket 2 ──► [b:2] → [c:3] → nil
//
//	Get("c"):
//	  1. hash("c") % capacity = 2
//	  2. Search bucket 2: [b:2] → no, [c:3] → yes!
//	  3. Return (3, true)
//
//	Get("z"):
//	  1. hash("z") % capacity = 1
//	  2. Search bucket 1: nil
//	  3. Return (0, false)
//
// complexity:
//   - time : O(1) average, O(n) worst case
//   - space: O(1)
func (h *HashMap[K, V]) Get(key K) (V, bool) {
	if !h.Empty() {
		index := h.bucketIndex(key)
		entries := h.buckets[index]
		for v := range entries.Iter {
			if v.key == key {
				return v.val, true
			}
		}
	}
	return generics.ZeroValue[V](), false
}

// Exists checks if a key is present in the map.
//
//	Exists("b") → true
//	Exists("z") → false
//
// complexity:
//   - time : O(1) average
//   - space: O(1)
func (h *HashMap[K, V]) Exists(key K) bool {
	_, found := h.Get(key)
	return found
}

// String returns the string representation.
//
//	{a:1, b:2, c:3}
//
//	String() → "[a:1 b:2 c:3]"
//
// complexity:
//   - time : O(n)
//   - space: O(n)
func (h *HashMap[K, V]) String() string {
	var buf strings.Builder
	buf.WriteRune('[')
	for i, e := range sequence.Enum(h.Iter) {
		if i > 0 {
			buf.WriteRune(' ')
		}
		_, _ = fmt.Fprintf(&buf, "%v:%v", e.Key(), e.Value())
	}
	buf.WriteRune(']')
	return buf.String()
}

func (h *HashMap[K, V]) bucketIndex(key K) int {
	hash := h.hashFunction(key)
	if hash < 0 {
		hash = -hash
	}
	return hash % len(h.buckets)
}

// growAndRehash doubles the capacity and redistributes all entries.
//
//	Before (capacity=4, size=3):
//
//	bucket 0 ──► [a:1] → nil
//	bucket 1 ──► nil
//	bucket 2 ──► [b:2] → [c:3] → nil
//	bucket 3 ──► nil
//
//	After (capacity=8, rehashed):
//
//	bucket 0 ──► nil
//	bucket 1 ──► nil
//	bucket 2 ──► [b:2] → nil   ← entries moved to new positions
//	bucket 3 ──► nil
//	bucket 4 ──► nil
//	bucket 5 ──► [a:1] → nil
//	bucket 6 ──► [c:3] → nil
//	bucket 7 ──► nil
//
// Rehashing is needed because bucket indices depend on capacity:
// hash(key) % old_capacity ≠ hash(key) % new_capacity
//
// complexity:
//   - time : O(n)
//   - space: O(n)
func (h *HashMap[K, V]) growAndRehash() {
	h2 := NewWith[K, V](Options[K]{
		Capacity:      2 * len(h.buckets),
		LoadThreshold: h.loadThreshold,
		HashFunction:  h.hashFunction,
	})

	for e := range h.Iter {
		h2.put(e.Key(), e.Value())
	}

	h.size = h2.size
	h.buckets = h2.buckets
	h2 = nil
}

// Size returns the number of key-value pairs.
//
// complexity:
//   - time : O(1)
//   - space: O(1)
func (h *HashMap[K, V]) Size() int {
	return h.size
}

// Cap returns the number of buckets.
// This satisfies the adt.Caper interface.
//
// complexity:
//   - time : O(1)
//   - space: O(1)
func (h *HashMap[K, V]) Cap() int {
	return len(h.buckets)
}

// LoadFactor returns size/capacity.
//
//	size = 12, capacity = 16
//	LoadFactor() → 12/16 = 0.75
//
// When LoadFactor >= LoadThreshold, the map resizes.
//
// complexity:
//   - time : O(1)
//   - space: O(1)
func (h *HashMap[K, V]) LoadFactor() float64 {
	return float64(h.Size()) / float64(h.Cap())
}

// Empty returns true if the map has no entries.
//
// complexity:
//   - time : O(1)
//   - space: O(1)
func (h *HashMap[K, V]) Empty() bool {
	return h.Size() == 0
}
