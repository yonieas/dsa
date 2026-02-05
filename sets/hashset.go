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
//
// SCORE: 10
func (s *HashSet[E]) Add(data E) {
	// hint: call s.backend.Put(data, none{})
	//       - the key is the element, value is empty struct
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
//
// SCORE: 10
func (s *HashSet[E]) Del(data E) {
	// hint: call s.backend.Del(data)
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
//
// SCORE: 10
func (s *HashSet[E]) Exists(data E) bool {
	// hint: return s.backend.Exists(data)
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
//
// SCORE: 10
func (s *HashSet[E]) String() string {
	// hint: use strings.Builder, iterate with sequence.Enum(s.Iter)
	//       format as "{elem1 elem2 ...}"
	var sb strings.Builder
	sb.WriteByte('{')
	for i, e := range sequence.Enum(s.Iter) {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%v", e))
	}
	sb.WriteByte('}')
	return sb.String()
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
//
// SCORE: 10
func (s *HashSet[E]) Iter(yield func(E) bool) {
	// hint: iterate s.backend.Iter, yield e.Key() for each entry
	for e := range s.backend.Iter {
		if !yield(e.Key()) {
			return
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
//
// SCORE: 15
func (s *HashSet[E]) Union(s2 *HashSet[E]) *HashSet[E] {
	// hint: 1) create new empty set: union := New[E]()
	//       2) iterate s, add each element to union
	//       3) iterate s2, add each element to union
	//       4) return union
	union := New[E]()
	for e := range s.backend.Keys {
		union.Add(e)
	}
	for e := range s2.backend.Keys {
		union.Add(e)
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
//
// SCORE: 20
func (s *HashSet[E]) Intersection(s2 *HashSet[E]) *HashSet[E] {
	// hint: 1) create new empty set: intersection := New[E]()
	//       2) choose smaller set to iterate (efficiency)
	//       3) for each element in smaller set, if other.Exists(v), add to result
	//       4) return intersection
	intersection := New[E]()
	if s.Size() < s2.Size() {
		for e := range s.backend.Keys {
			if s2.Exists(e) {
				intersection.Add(e)
			}
		}
	} else {
		for e := range s2.backend.Keys {
			if s.Exists(e) {
				intersection.Add(e)
			}
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
//
// SCORE: 15
func (s *HashSet[E]) Disjoint(s2 *HashSet[E]) bool {
	// hint: 1) choose smaller set to iterate (efficiency)
	//       2) for each element in smaller set, if other.Exists(v), return false
	//       3) return true (no common elements found)
	if s.Size() < s2.Size() {
		for v := range s.backend.Keys {
			if s2.backend.Exists(v) {
				return false
			}
		}
	} else {
		for v := range s2.backend.Keys {
			if s.backend.Exists(v) {
				return false
			}
		}
	}
	return true
}
