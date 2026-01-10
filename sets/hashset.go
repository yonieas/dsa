// Package sets provides Set data structure implementations.
//
// # What is a Set?
//
// A Set is a collection of unique elements with no duplicates and no defined
// ordering. It answers one question fast: "Is this element present?" The
// answer is O(1) on average thanks to hashing.
//
// # Why Use Sets?
//
// Sets are perfect when you care about membership but not order or count.
// Deduplicating a list? Add everything to a set. Finding common elements
// between two groups? Intersect their sets. Checking if a username is taken?
// One set lookup.
//
// # Set Operations
//
//	Add:          Insert an element (no-op if already present)
//	Del:          Remove an element (no-op if not present)
//	Exists:       Check if element is present
//	Union:        All elements from both sets (A or B)
//	Intersection: Elements present in both sets (A and B)
//	Disjoint:     True if sets share no common elements
//
// # Implementation
//
// This HashSet is built on a HashMap with empty struct values. The keys are
// the set elements. Since struct{} uses zero bytes, we get the full power of
// hash table lookups with minimal overhead.
//
// # Complexity
//
//	Add/Del/Exists:   O(1) average
//	Union/Intersect:  O(n + m) where n and m are set sizes
//	Space:            O(n)
//
// # Further Reading
//
// CLRS "Introduction to Algorithms", Chapter 11 (Hash Tables).
// https://en.wikipedia.org/wiki/Set_(abstract_data_type)
package sets

import (
	"fmt"
	"strings"

	"github.com/josestg/dsa/hashmap"
	"github.com/josestg/dsa/sequence"
)

// none is an empty struct used as placeholder value.
// struct{} uses 0 bytes of memory - we only care about keys.
type none = struct{}

// Options configures HashSet behavior (same as HashMap options).
type Options[E comparable] = hashmap.Options[E]

// HashSet is a set implementation backed by a HashMap.
//
// Internally stores elements as HashMap keys with empty values:
//
//	┌───────────────────────────────┐
//	│  "apple"  → {}               │
//	│  "banana" → {}               │
//	│  "cherry" → {}               │
//	└───────────────────────────────┘
//
// The value {} (empty struct) uses zero memory.
type HashSet[E comparable] struct {
	backend *hashmap.HashMap[E, none]
}

// New creates an empty HashSet.
//
//	{}  ← empty set
//
// complexity:
//   - time : O(1)
//   - space: O(1)
func New[E comparable]() *HashSet[E] {
	return &HashSet[E]{
		backend: hashmap.New[E, none](),
	}
}

// NewWith creates a HashSet with custom configuration.
func NewWith[E comparable](opts Options[E]) *HashSet[E] {
	return &HashSet[E]{
		backend: hashmap.NewWith[E, none](hashmap.Options[E](opts)),
	}
}

// Add inserts an element into the set.
//
//	Before Add("cherry"):
//
//	{ "apple", "banana" }
//
//	After Add("cherry"):
//
//	{ "apple", "banana", "cherry" }
//
//	After Add("apple") again:
//
//	{ "apple", "banana", "cherry" }  ← no change (already exists)
//
// complexity:
//   - time : O(1) average
//   - space: O(1)
func (s *HashSet[E]) Add(data E) {
	s.backend.Put(data, none{})
}

// Del removes an element from the set.
//
//	Before Del("banana"):
//
//	{ "apple", "banana", "cherry" }
//
//	After Del("banana"):
//
//	{ "apple", "cherry" }
//
//	After Del("grape"):
//
//	{ "apple", "cherry" }  ← no change (wasn't in set)
//
// complexity:
//   - time : O(1) average
//   - space: O(1)
func (s *HashSet[E]) Del(data E) {
	s.backend.Del(data)
}

// Exists checks if an element is in the set.
//
//	{ "apple", "banana", "cherry" }
//
//	Exists("banana") → true
//	Exists("grape")  → false
//
// complexity:
//   - time : O(1) average
//   - space: O(1)
func (s *HashSet[E]) Exists(data E) bool {
	return s.backend.Exists(data)
}

// Size returns the number of elements in the set.
//
//	{ "apple", "banana", "cherry" }
//
//	Size() → 3
//
// complexity:
//   - time : O(1)
//   - space: O(1)
func (s *HashSet[E]) Size() int {
	return s.backend.Size()
}

// Empty returns true if the set has no elements.
//
//	{}                     { "apple" }
//	Empty() → true         Empty() → false
//
// complexity:
//   - time : O(1)
//   - space: O(1)
func (s *HashSet[E]) Empty() bool {
	return s.backend.Empty()
}

// String returns the string representation of the set.
//
//	{ 1, 2, 3 }
//
//	String() → "{1 2 3}"
//
// Note: Order is not guaranteed.
//
// complexity:
//   - time : O(n)
//   - space: O(n)
func (s *HashSet[E]) String() string {
	var buf strings.Builder
	buf.WriteRune('{')
	for i, v := range sequence.Enum(s.Iter) {
		if i > 0 {
			buf.WriteRune(' ')
		}
		_, _ = fmt.Fprint(&buf, v)
	}
	buf.WriteRune('}')
	return buf.String()
}

// Iter iterates over all elements in the set.
//
//	{ "apple", "banana", "cherry" }
//
//	for elem := range set.Iter {
//	    fmt.Println(elem)  // prints each element
//	}
//
// Order is not guaranteed.
//
// complexity:
//   - time : O(n)
//   - space: O(1)
func (s *HashSet[E]) Iter(yield func(E) bool) {
	for e := range s.backend.Iter {
		if !yield(e.Key()) {
			break
		}
	}
}

// Union returns a new set with elements from both sets.
//
//	A = { 1, 2, 3 }
//	B = { 3, 4, 5 }
//
//	A ∪ B (union):
//
//	    A           B
//	 ┌─────┐    ┌─────┐
//	 │1 2 3│────│3 4 5│
//	 └─────┘    └─────┘
//	     └───┬───┘
//	         ↓
//	   { 1, 2, 3, 4, 5 }
//
// Neither original set is modified.
//
// complexity:
//   - time : O(n + m) where n, m are sizes of the two sets
//   - space: O(n + m) for the result
func (s *HashSet[E]) Union(s2 *HashSet[E]) *HashSet[E] {
	union := New[E]()
	for v := range s.Iter {
		union.Add(v)
	}
	for v := range s2.Iter {
		union.Add(v)
	}
	return union
}

// Intersection returns a new set with elements in both sets.
//
//	A = { 1, 2, 3, 4 }
//	B = { 3, 4, 5, 6 }
//
//	A ∩ B (intersection):
//
//	      A           B
//	   ┌─────┐    ┌─────┐
//	   │1 2│3 4│──│3 4│5 6│
//	   └─────┘    └─────┘
//	        └──┬──┘
//	           ↓
//	       { 3, 4 }
//
// Neither original set is modified.
//
// complexity:
//   - time : O(min(n, m)) - iterates smaller set
//   - space: O(min(n, m)) for the result
func (s *HashSet[E]) Intersection(s2 *HashSet[E]) *HashSet[E] {
	intersection := New[E]()
	var left, right *HashSet[E]
	if s.Size() < s2.Size() {
		left, right = s, s2
	} else {
		left, right = s2, s
	}
	for v := range left.Iter {
		if right.Exists(v) {
			intersection.Add(v)
		}
	}
	return intersection
}

// Disjoint returns true if the two sets have no common elements.
//
//	A = { 1, 2, 3 }
//	B = { 4, 5, 6 }
//	A.Disjoint(B) → true (no overlap)
//
//	A = { 1, 2, 3 }
//	C = { 3, 4, 5 }
//	A.Disjoint(C) → false (3 is in both)
//
// complexity:
//   - time : O(min(n, m))
//   - space: O(1)
func (s *HashSet[E]) Disjoint(s2 *HashSet[E]) bool {
	var left, right *HashSet[E]
	if s.Size() < s2.Size() {
		left, right = s, s2
	} else {
		left, right = s2, s
	}
	for v := range left.Iter {
		if right.Exists(v) {
			return false
		}
	}
	return true
}
