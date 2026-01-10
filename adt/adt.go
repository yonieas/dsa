// Package adt defines Abstract Data Type (ADT) interfaces.
//
// # What is an Abstract Data Type?
//
// An Abstract Data Type describes WHAT operations a data structure supports,
// but not HOW those operations are implemented. It is a contract: any type
// that satisfies the interface can be used interchangeably.
//
// For example, a Stack ADT specifies Push, Pop, and Peek operations. Whether
// the Stack uses an array or linked list internally is hidden. You can swap
// implementations without changing the code that uses them.
//
// # Why Use ADTs?
//
// ADTs let you think at the right level of abstraction. When solving a problem,
// you focus on what operations you need (a queue for BFS, a stack for DFS)
// rather than implementation details. Later, you can choose or swap the
// concrete implementation based on performance needs.
//
// This package provides small, composable interfaces. Rather than one large
// interface, we define Sizer, Emptier, Adder, and so on. Types implement only
// what they need, and generic algorithms constrain on only what they use.
//
// # Design Philosophy
//
// Go favors small interfaces. The standard library's io.Reader and io.Writer
// are single-method interfaces that compose beautifully. This package follows
// that pattern: each interface captures one capability. This makes interfaces
// easy to implement and algorithms maximally reusable.
//
// # Further Reading
//
// CLRS "Introduction to Algorithms", Chapter 10 (Elementary Data Structures).
// Go Blog: "Go Proverbs" for interface design philosophy.
// https://en.wikipedia.org/wiki/Abstract_data_type
package adt

import (
	"cmp"
	"fmt"
)

// Sizer describes a data structure that tracks its element count.
//
//	┌───┬───┬───┬───┬───┐
//	│ A │ B │ C │ D │ E │
//	└───┴───┴───┴───┴───┘
//	Size() -> 5
//
// The Size method returns the number of elements currently stored.
// An empty structure returns 0.
type Sizer interface {
	Size() int
}

// Caper describes a data structure with a capacity limit.
//
//	capacity = 8
//	┌───┬───┬───┬───┬───┬───┬───┬───┐
//	│ A │ B │ C │ D │ E │   │   │   │  <- 3 empty slots
//	└───┴───┴───┴───┴───┴───┴───┴───┘
//	Size() -> 5    Cap() -> 8
//
// The Cap method returns the maximum number of elements that can be stored
// without requiring the structure to resize or reallocate memory.
type Caper interface {
	Cap() int
}

// Emptier describes a data structure that can report if it has no elements.
//
//	Empty structure:         Non-empty structure:
//	┌───┐                    ┌───┬───┬───┐
//	│   │ (no elements)      │ A │ B │ C │
//	└───┘                    └───┴───┴───┘
//	Empty() -> true           Empty() -> false
//
// The Empty method returns true if and only if Size() equals 0.
type Emptier interface {
	Empty() bool
}

// Getter describes a data structure that supports index-based element access.
//
//	┌───┬───┬───┬───┬───┐
//	│ A │ B │ C │ D │ E │
//	└───┴───┴───┴───┴───┘
//	  0   1   2   3   4   <- indices
//
//	Get(0) -> A    (first element)
//	Get(2) -> C    (middle element)
//	Get(4) -> E    (last element)
//
// The Get method retrieves the element at the specified index.
// Valid indices are in range [0, Size()-1].
//
// Panics:
//   - If index is negative
//   - If index >= Size()
//   - If the structure is empty
type Getter[T any] interface {
	Get(int) T
}

// Setter describes a data structure that supports updating elements by index.
//
//	Before Set(2, X):        After Set(2, X):
//	┌───┬───┬───┬───┬───┐    ┌───┬───┬───┬───┬───┐
//	│ A │ B │ C │ D │ E │    │ A │ B │ X │ D │ E │
//	└───┴───┴───┴───┴───┘    └───┴───┴───┴───┴───┘
//	  0   1   2   3   4        0   1   2   3   4
//	          ↑                        ↑
//	       updated                  updated
//
// The Set method replaces the element at the given index with a new value.
// This does NOT change the size of the structure.
// Valid indices are in range [0, Size()-1].
//
// Panics:
//   - If index is negative
//   - If index >= Size()
type Setter[T any] interface {
	Set(int, T)
}

// Appender describes a data structure that supports adding elements at the end.
//
//	Before Append(F):        After Append(F):
//	┌───┬───┬───┬───┬───┐    ┌───┬───┬───┬───┬───┬───┐
//	│ A │ B │ C │ D │ E │    │ A │ B │ C │ D │ E │ F │
//	└───┴───┴───┴───┴───┘    └───┴───┴───┴───┴───┴───┘
//	                  ↑                            ↑
//	              old tail                     new tail
//
// The Append method inserts a new element after the last element.
// The appended element becomes the new Tail().
// Size increases by 1.
type Appender[T any] interface {
	Append(T)
}

// Prepender describes a data structure that supports adding elements at the front.
//
//	Before Prepend(Z):       After Prepend(Z):
//	┌───┬───┬───┬───┬───┐    ┌───┬───┬───┬───┬───┬───┐
//	│ A │ B │ C │ D │ E │    │ Z │ A │ B │ C │ D │ E │
//	└───┴───┴───┴───┴───┘    └───┴───┴───┴───┴───┴───┘
//	  ↑                        ↑
//	old head               new head
//
// The Prepend method inserts a new element before the first element.
// The prepended element becomes the new Head().
// Size increases by 1.
type Prepender[T any] interface {
	Prepend(T)
}

// Tailer describes a data structure that provides access to its last element.
//
//	┌───┬───┬───┬───┬───┐
//	│ A │ B │ C │ D │ E │
//	└───┴───┴───┴───┴───┘
//	                  ↑
//	                tail
//
//	Tail() -> E
//
// The Tail method returns the last element without removing it.
// This is equivalent to Get(Size()-1).
//
// Panics:
//   - If the structure is empty
type Tailer[T any] interface {
	Tail() T
}

// Header describes a data structure that provides access to its first element.
//
//	┌───┬───┬───┬───┬───┐
//	│ A │ B │ C │ D │ E │
//	└───┴───┴───┴───┴───┘
//	  ↑
//	head
//
//	Head() -> A
//
// The Head method returns the first element without removing it.
// This is equivalent to Get(0).
//
// Panics:
//   - If the structure is empty
type Header[T any] interface {
	Head() T
}

// Popper describes a data structure that supports removing elements from the end.
//
//	Before Pop():            After Pop():
//	┌───┬───┬───┬───┬───┐    ┌───┬───┬───┬───┐
//	│ A │ B │ C │ D │ E │    │ A │ B │ C │ D │
//	└───┴───┴───┴───┴───┘    └───┴───┴───┴───┘
//	                  ↑                    ↑
//	              removed              new tail
//
//	Pop() -> E (the removed element)
//
// The Pop method removes and returns the last element.
// Size decreases by 1.
//
// This is a fundamental operation for Stack (LIFO - Last In, First Out).
//
// Panics:
//   - If the structure is empty
type Popper[T any] interface {
	Pop() T
}

// Pusher describes a data structure that supports adding elements to the top.
//
//	Before Push(E):          After Push(E):
//	┌─────┐                  ┌─────┐
//	│  D  │ <- top           │  E  │ <- new top
//	├─────┤                  ├─────┤
//	│  C  │                  │  D  │
//	├─────┤                  ├─────┤
//	│  B  │                  │  C  │
//	├─────┤                  ├─────┤
//	│  A  │                  │  B  │
//	└─────┘                  ├─────┤
//	                         │  A  │
//	                         └─────┘
//
// The Push method adds an element to the top of a stack.
// The pushed element becomes the new top (returned by Peek, removed by Pop).
// Size increases by 1.
//
// This is a fundamental operation for Stack (LIFO - Last In, First Out).
type Pusher[T any] interface {
	Push(T)
}

// Shifter describes a data structure that supports removing elements from the front.
//
//	Before Shift():          After Shift():
//	┌───┬───┬───┬───┬───┐    ┌───┬───┬───┬───┐
//	│ A │ B │ C │ D │ E │    │ B │ C │ D │ E │
//	└───┴───┴───┴───┴───┘    └───┴───┴───┴───┘
//	  ↑                        ↑
//	removed                new head
//
//	Shift() -> A (the removed element)
//
// The Shift method removes and returns the first element.
// All remaining elements shift to lower indices.
// Size decreases by 1.
//
// This is a fundamental operation for Queue (FIFO - First In, First Out).
//
// Panics:
//   - If the structure is empty
type Shifter[T any] interface {
	Shift() T
}

// Iterator describes a data structure that can be traversed element by element.
//
//	┌───┬───┬───┬───┬───┐
//	│ A │ B │ C │ D │ E │
//	└───┴───┴───┴───┴───┘
//	  ↓   ↓   ↓   ↓   ↓
//	  1   2   3   4   5    <- iteration order (forward)
//
// The Iter method accepts a yield function called for each element.
//
// Example using Go 1.23+ range-over-func:
//
//	for value := range structure.Iter {
//	    fmt.Println(value)
//	}
type Iterator[T any] interface {
	Iter(func(T) bool)
}

// Enumerator describes a data structure that can be traversed with index-value pairs.
//
//	┌───┬───┬───┬───┬───┐
//	│ A │ B │ C │ D │ E │
//	└───┴───┴───┴───┴───┘
//	  0   1   2   3   4    <- indices
//	  ↓   ↓   ↓   ↓   ↓
//	 1st 2nd 3rd 4th 5th   <- enumeration order (forward)
//
// The Enum method accepts a yield function called for each element with its index.
// This is useful when you need both position and value during iteration.
//
// Example using Go 1.23+ range-over-func:
//
//	for index, value := range structure.Enum {
//	    fmt.Printf("%d: %v\n", index, value)
//	}
type Enumerator[T any] interface {
	Enum(func(index int, value T) bool)
}

// BackwardEnumerator describes a data structure that can be traversed in reverse
// with index-value pairs.
//
//	┌───┬───┬───┬───┬───┐
//	│ A │ B │ C │ D │ E │
//	└───┴───┴───┴───┴───┘
//	  0   1   2   3   4    <- indices
//	  ↓   ↓   ↓   ↓   ↓
//	 5th 4th 3rd 2nd 1st   <- enumeration order (backward)
//
// The EnumBackward method visits elements from tail to head,
// providing both index and value for each element.
//
// Example using Go 1.23+ range-over-func:
//
//	for index, value := range structure.EnumBackward {
//	    fmt.Printf("%d: %v\n", index, value)  // prints 4:E, 3:D, 2:C, 1:B, 0:A
//	}
type BackwardEnumerator[T any] interface {
	EnumBackward(func(index int, value T) bool)
}

// Entrier describes a map-like data structure that can be traversed with key-value pairs.
//
//	┌─────────────────────────────────────┐
//	│  "alice" -> 30                      │
//	│  "bob"   -> 25                      │
//	│  "carol" -> 28                      │
//	└─────────────────────────────────────┘
//
// The Entries method accepts a yield function called for each key-value pair.
// This is useful for iterating over all entries in a map or dictionary.
//
// Example using Go 1.23+ range-over-func:
//
//	for key, value := range m.Entries {
//	    fmt.Printf("%v: %v\n", key, value)
//	}
//
// Note: Iteration order depends on the implementation.
// For ordered maps, entries are visited in key order.
// For hash maps, order may not be deterministic.
type Entrier[K cmp.Ordered, V any] interface {
	Entries(func(key K, value V) bool)
}

// BackwardIterator describes a data structure that can be traversed in reverse.
//
//	┌───┬───┬───┬───┬───┐
//	│ A │ B │ C │ D │ E │
//	└───┴───┴───┴───┴───┘
//	  ↓   ↓   ↓   ↓   ↓
//	  5   4   3   2   1    <- iteration order (backward)
//
// The IterBackward method visits elements from tail to head.
type BackwardIterator[T any] interface {
	IterBackward(func(T) bool)
}

// Peeker describes a data structure that allows viewing the next element to be removed.
//
//	Stack (LIFO):            Queue (FIFO):
//	┌───┐                    ┌───┬───┬───┬───┬───┐
//	│ E │ <- Peek()          │ A │ B │ C │ D │ E │
//	├───┤                    └───┴───┴───┴───┴───┘
//	│ D │                      ↑
//	├───┤                    Peek() (front)
//	│ C │
//	└───┘
//
// The Peek method returns the element that would be removed by Pop (Stack)
// or Dequeue (Queue), without actually removing it.
//
// Panics:
//   - If the structure is empty
type Peeker[T any] interface {
	Peek() T
}

// Enqueuer describes a data structure that supports adding elements to the rear.
//
//	Before Enqueue(F):       After Enqueue(F):
//	front             rear   front                 rear
//	  ↓                 ↓      ↓                   ↓
//	┌───┬───┬───┬───┬───┐    ┌───┬───┬───┬───┬───┬───┐
//	│ A │ B │ C │ D │ E │    │ A │ B │ C │ D │ E │ F │
//	└───┴───┴───┴───┴───┘    └───┴───┴───┴───┴───┴───┘
//	                                               ↑
//	                                            added
//
// The Enqueue method adds an element to the rear of a queue.
// The enqueued element becomes the new rear.
// Size increases by 1.
//
// This is a fundamental operation for Queue (FIFO - First In, First Out).
type Enqueuer[T any] interface {
	Enqueue(T)
}

// Dequeuer describes a data structure that supports removing elements from the front.
//
//	Before Dequeue():        After Dequeue():
//	front             rear   front         rear
//	  ↓               ↓       ↓             ↓
//	┌───┬───┬───┬───┬───┐    ┌───┬───┬───┬───┐
//	│ A │ B │ C │ D │ E │    │ B │ C │ D │ E │
//	└───┴───┴───┴───┴───┘    └───┴───┴───┴───┘
//	  ↑
//	removed
//
//	Dequeue() -> A (the removed element)
//
// The Dequeue method removes and returns the front element.
// Size decreases by 1.
//
// This is a fundamental operation for Queue (FIFO - First In, First Out).
//
// Panics:
//   - If the structure is empty
type Dequeuer[T any] interface {
	Dequeue() T
}

// Stringer is an alias for fmt.Stringer from the standard library.
// Data structures implementing this interface provide a human-readable
// string representation, useful for debugging and logging.
type Stringer = fmt.Stringer

// Stack defines the interface for a Last-In-First-Out (LIFO) data structure.
//
//	   Push(D)
//	   ↓
//	┌─────┐
//	│  D  │ <- top (Peek returns D, Pop removes D)
//	├─────┤
//	│  C  │
//	├─────┤
//	│  B  │
//	├─────┤
//	│  A  │
//	└─────┘
//
// The last element pushed is the first to be popped (LIFO).
type Stack[T any] interface {
	Sizer
	Emptier
	Peeker[T]
	Popper[T]
	Pusher[T]
}

// Queue defines the interface for a First-In-First-Out (FIFO) data structure.
//
//	                  front                        rear
//	                    │                            │
//	                    ↓                            ↓
//	Dequeue() <-   ┌───┬───┬───┬───┬───┐   <- Enqueue(F)
//	               │ A │ B │ C │ D │ E │
//	               └───┴───┴───┴───┴───┘
//	                 ↑
//	            Peek() returns A
//
// The first element enqueued is the first to be dequeued (FIFO).
type Queue[T any] interface {
	Sizer
	Emptier
	Peeker[T]
	Enqueuer[T]
	Dequeuer[T]
}

// Remover describes a data structure that supports removing elements by index.
//
//	Before Remove(2):        After Remove(2):
//	┌───┬───┬───┬───┬───┐    ┌───┬───┬───┬───┐
//	│ A │ B │ C │ D │ E │    │ A │ B │ D │ E │
//	└───┴───┴───┴───┴───┘    └───┴───┴───┴───┘
//	  0   1   2   3   4        0   1   2   3
//	          ↑
//	       removed
//
//	Remove(2) -> C (the removed element)
//
// The Remove method deletes the element at the given index and returns it.
// Elements after the removed index shift to lower indices.
// Size decreases by 1.
//
// Panics:
//   - If index is negative
//   - If index >= Size()
type Remover[E any] interface {
	Remove(index int) E
}

// Inserter describes a data structure that supports inserting elements at any position.
//
//	Before Insert(2, X):     After Insert(2, X):
//	┌───┬───┬───┬───┬───┐    ┌───┬───┬───┬───┬───┬───┐
//	│ A │ B │ C │ D │ E │    │ A │ B │ X │ C │ D │ E │
//	└───┴───┴───┴───┴───┘    └───┴───┴───┴───┴───┴───┘
//	  0   1   2   3   4        0   1   2   3   4   5
//	          ↑                        ↑
//	    insert here              new element
//
// The Insert method adds a new element at the specified index.
// Elements at and after that index shift to higher indices.
// Size increases by 1.
// Valid indices are in range [0, Size()] (can insert at the end).
//
// Panics:
//   - If index is negative
//   - If index > Size()
type Inserter[E any] interface {
	Insert(index int, data E)
}

// Adder describes a data structure that accepts elements without specifying position.
//
// The placement of the element depends on the data structure:
//   - Set: Added if not already present (maintains uniqueness)
//   - BST: Placed according to ordering rules
//   - Bag: Always added (allows duplicates)
//
// This is different from Append/Prepend which specify position.
type Adder[E any] interface {
	Add(E)
}

// Deleter describes a data structure that supports removing elements by value.
//
//	Before Del(C):           After Del(C):
//	┌───┬───┬───┬───┬───┐    ┌───┬───┬───┬───┐
//	│ A │ B │ C │ D │ E │    │ A │ B │ D │ E │
//	└───┴───┴───┴───┴───┘    └───┴───┴───┴───┘
//	          ↑
//	     value to delete
//
// The Del method removes the specified element if it exists.
// If the element does not exist, no action is taken (no error).
// Size decreases by 1 if an element was removed.
type Deleter[E any] interface {
	Del(E)
}

// Exister describes a data structure that can check for element membership.
//
//	┌───┬───┬───┬───┬───┐
//	│ A │ B │ C │ D │ E │
//	└───┴───┴───┴───┴───┘
//
//	Exists(C) -> true
//	Exists(Z) -> false
//
// The Exists method returns true if the element is present in the structure.
// This is also known as "Contains" in some libraries.
type Exister[E any] interface {
	Exists(E) bool
}

// Keys describes a data structure that can iterate over its keys.
// This is used by Map-like structures where elements are key-value pairs.
//
//	┌─────────────────────────────────────┐
//	│  "alice" -> 30                      │
//	│  "bob"   -> 25                      │
//	│  "carol" -> 28                      │
//	└─────────────────────────────────────┘
//
//	for key := range m.Keys {
//	    // yields: "alice", "bob", "carol"
//	}
//
// The Keys method provides an iterator over all keys.
// Iteration order depends on the implementation (may not be deterministic).
type Keys[K comparable] interface {
	Keys(yield func(K) bool)
}

// Rotator describes a circular data structure that can rotate its elements.
//
//	Before Rotate(1):           After Rotate(1):
//	┌───┬───┬───┐               ┌───┬───┬───┐
//	│ A │ B │ C │──► back to A  │ A │ B │ C │──► back to A
//	└───┴───┴───┘               └───┴───┴───┘
//	  ↑                               ↑
//	 head                            head
//
// The Rotate method moves the head pointer n positions.
// Positive n rotates forward, negative n rotates backward.
// Rotating by the size of the structure returns to the original position.
type Rotator interface {
	Rotate(n int)
}

// Cycler describes a circular data structure that can cycle through elements.
//
//	Before Cycle():             After Cycle():
//	┌───┬───┬───┐               ┌───┬───┬───┐
//	│ A │ B │ C │──► back to A  │ A │ B │ C │──► back to A
//	└───┴───┴───┘               └───┴───┴───┘
//	  ↑                               ↑
//	 head                            head
//
//	Cycle() -> A (returns current head, advances to next)
//
// The Cycle method returns the current head and advances to the next element.
// When reaching the end, it wraps around to the beginning.
// Useful for round-robin iteration.
type Cycler[T any] interface {
	Cycle() T
}

// ReverseCycler describes a circular data structure that can cycle backward.
//
//	Before ReverseCycle():      After ReverseCycle():
//	┌───┬───┬───┐               ┌───┬───┬───┐
//	│ A │ B │ C │──► back to A  │ A │ B │ C │──► back to A
//	└───┴───┴───┘               └───┴───┴───┘
//	  ↑                                   ↑
//	 head                                head
//
//	ReverseCycle() -> C (moves head backward, returns new head)
//
// The ReverseCycle method moves the head backward and returns the new head.
// When reaching the beginning, it wraps around to the end.
// Useful for reverse round-robin iteration.
type ReverseCycler[T any] interface {
	ReverseCycle() T
}

// Unioner describes a set that can compute the union with another set.
//
//	A = { 1, 2, 3 }
//	B = { 3, 4, 5 }
//
//	A ∪ B (union):
//	┌─────────────────┐
//	│ 1, 2, 3, 4, 5   │
//	└─────────────────┘
//
// The Union method returns a new set containing all elements from both sets.
// Neither original set is modified.
type Unioner[Self any] interface {
	Union(Self) Self
}

// Intersecter describes a set that can compute the intersection with another set.
//
//	A = { 1, 2, 3, 4 }
//	B = { 3, 4, 5, 6 }
//
//	A ∩ B (intersection):
//	┌─────────┐
//	│  3, 4   │
//	└─────────┘
//
// The Intersection method returns a new set containing only elements in both sets.
// Neither original set is modified.
type Intersecter[Self any] interface {
	Intersection(Self) Self
}

// Disjointer describes a set that can check if it has no common elements with another.
//
//	A = { 1, 2, 3 }
//	B = { 4, 5, 6 }
//	A.Disjoint(B) → true (no overlap)
//
//	A = { 1, 2, 3 }
//	C = { 3, 4, 5 }
//	A.Disjoint(C) → false (3 is in both)
//
// The Disjoint method returns true if the intersection is empty.
type Disjointer[Self any] interface {
	Disjoint(Self) bool
}
