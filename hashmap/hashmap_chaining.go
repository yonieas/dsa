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
//
// SCORE: 15
func (h *HashMap[K, V]) Put(key K, value V) {
	// hint: 1) if LoadFactor() >= threshold, call growAndRehash()
	//       2) call h.put(key, value)
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
//
// SCORE: 15
func (h *HashMap[K, V]) Del(key K) {
	// hint: 1) get bucket index using bucketIndex(key)
	//       2) iterate bucket (entries) with sequence.Enum to track index
	//       3) if entry.key == key: entries.Remove(i), decrement size, break
	bucketIdx := h.bucketIndex(key)
	entry := h.buckets[bucketIdx]
	idx := 0
	entry.Iter(func(e *Entry[K, V]) bool {
		if e.Key() == key {
			entry.Remove(idx)
			h.size--
			return false
		}
		idx++
		return true
	})
}

// SCORE: 10
func (h *HashMap[K, V]) put(key K, value V) {
	// hint: 1) get bucket index using bucketIndex(key)
	//       2) iterate bucket to check if key exists, update value if found
	//       3) if not found: entries.Append(NewEntry(key, value)), increment size
	bucketIdx := h.bucketIndex(key)
	entry := h.buckets[bucketIdx]
	found := false
	entry.Iter(func(e *Entry[K, V]) bool {
		if e.Key() == key {
			e.val = value
			found = true
			return false // Stop searching immediately
		}
		return true
	})
	// Append if the key of entry not found in the bucket
	if !found {
		entry.Append(NewEntry(key, value))
		h.size++
	}
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
//
// SCORE: 10
func (h *HashMap[K, V]) Iter(yield func(*Entry[K, V]) bool) {
	// hint: for each bucket in h.buckets, iterate entries and yield each
	for _, bucket := range h.buckets {
		for v := range bucket.Iter {
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
//
// SCORE: 5
func (h *HashMap[K, V]) Keys(yield func(K) bool) {
	// hint: use h.Iter and yield e.Key() for each entry
	h.Iter(func(e *Entry[K, V]) bool {
		if !yield(e.Key()) {
			return false
		}
		return true
	})
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
//
// SCORE: 15
func (h *HashMap[K, V]) Get(key K) (V, bool) {
	// hint: 1) if Empty(), return (zero, false)
	//       2) get bucket index, iterate bucket entries
	//       3) if entry.key == key, return (entry.val, true)
	//       4) return (zero, false)
	if h.Empty() {
		return generics.ZeroValue[V](), false
	}
	bucketIdx := h.bucketIndex(key)
	entry := h.buckets[bucketIdx]
	if entry == nil {
		return generics.ZeroValue[V](), false
	}
	found := false
	value := generics.ZeroValue[V]()
	entry.Iter(func(e *Entry[K, V]) bool {
		if e.Key() == key {
			value = e.Value()
			found = true
			return false // Stop searching
		}
		return true // Continue to next entry
	})

	if found {
		return value, true
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
//
// SCORE: 5
func (h *HashMap[K, V]) String() string {
	// hint: use strings.Builder, iterate with sequence.Enum(h.Iter)
	//       format each entry as "key:value"
	var sb strings.Builder
	sb.WriteByte('[')
	first := true
	h.Iter(func(e *Entry[K, V]) bool {
		// Give a space in every entry except the first
		if !first {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%v:%v", e.Key(), e.Value()))
		first = false
		return true
	})
	sb.WriteByte(']')
	return sb.String()
}

// SCORE: 10
func (h *HashMap[K, V]) bucketIndex(key K) int {
	// hint: 1) hash := h.hashFunction(key)
	//       2) if hash < 0, hash = -hash (make positive)
	//       3) return hash % len(h.buckets)
	hash := h.hashFunction(key)
	if hash < 0 {
		hash = hash * -1
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
//
// SCORE: 15
func (h *HashMap[K, V]) growAndRehash() {
	// hint: 1) create new HashMap with 2x capacity using NewWith
	//       2) iterate all entries with h.Iter, call h2.put(key, val)
	//       3) replace h.size and h.buckets with h2's values
	nh := NewWith[K, V](Options[K]{
		Capacity:      h.Cap() * 2,
		LoadThreshold: h.loadThreshold,
		HashFunction:  h.hashFunction,
	})
	h.Iter(func(e *Entry[K, V]) bool {
		nh.put(e.Key(), e.Value())
		return true
	})
	h.size = nh.size
	h.buckets = nh.buckets
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
