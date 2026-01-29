// Package linkedlist provides linked list implementations.
//
// # What is a Linked List?
//
// A linked list is a sequence of nodes where each node holds data and a
// pointer to the next node. Unlike arrays, linked list elements are scattered
// in memory and connected through these pointers.
//
// The key insight is that insertion and deletion become O(1) operations when
// you have a reference to the target location. You just rewire the pointers,
// no shifting needed. However, you lose the ability to jump to any position
// directly since you must traverse from the head.
//
// # Variants
//
// This package provides four variants:
//
//	SinglyLinkedList:  Each node points to the next node only.
//	DoublyLinkedList:  Each node points to both next and previous.
//	CircularSingly:    Like singly, but tail connects back to head.
//	CircularDoubly:    Like doubly, but forms a complete loop.
//
// # Trade-offs vs Arrays
//
// Linked lists excel at insertion/deletion at known positions (O(1) vs O(n))
// and dynamic sizing without reallocation. Arrays win at random access
// (O(1) vs O(n)) and cache performance due to contiguous memory layout.
//
// Choose linked lists when you frequently insert/delete at the front, need
// stable references to nodes, or want to avoid resize overhead. Choose arrays
// when you need random access or iterate sequentially (cache-friendly).
//
// # Complexity
//
//	Prepend:       O(1)
//	Append:        O(n) singly, O(1) doubly with tail pointer
//	Access by idx: O(n)
//	Insert at pos: O(n) to find position, O(1) to insert
//	Delete at pos: O(n) to find position, O(1) to delete
//
// # Further Reading
//
// CLRS "Introduction to Algorithms", Chapter 10.2.
// Sedgewick "Algorithms", Section 1.3.
// https://en.wikipedia.org/wiki/Linked_list
package linkedlist

import (
	"github.com/josestg/dsa/internal/generics"
	"github.com/josestg/dsa/sequence"
)

// UnaryNode is a node in a singly linked list.
// Each node holds data and a pointer to the next node.
//
//	┌──────────────────┐
//	│  data  │  next ──┼───► (next node or nil)
//	└──────────────────┘
type UnaryNode[T any] struct {
	data T
	next *UnaryNode[T]
}

// NewUnaryNode creates a new node with the given data and next pointer.
//
//	NewUnaryNode(42, nil) creates:
//
//	┌──────────────────┐
//	│   42   │  nil    │
//	└──────────────────┘
//
//	NewUnaryNode(42, existingNode) creates:
//
//	┌──────────────────┐     ┌──────────────────┐
//	│   42   │  next ──┼────►│  existingNode    │
//	└──────────────────┘     └──────────────────┘
func NewUnaryNode[T any](data T, next *UnaryNode[T]) *UnaryNode[T] {
	return &UnaryNode[T]{
		data: data,
		next: next,
	}
}

// SinglyLinkedList is a linked list where each node points only to the next node.
// It maintains pointers to both head (front) and tail (back) for efficient operations.
//
//	head                 tail
//	  ↓                    ↓
//	┌───┐     ┌───┐     ┌───┐
//	│ A │────►│ B │────►│ C │────► nil
//	└───┘     └───┘     └───┘
//
// With head and tail pointers:
//   - Append (add to back): O(1)
//   - Prepend (add to front): O(1)
//   - Pop (remove from back): O(n) - must traverse to find second-to-last
//   - Shift (remove from front): O(1)
type SinglyLinkedList[T any] struct {
	head *UnaryNode[T]
	tail *UnaryNode[T]
	size int
}

// NewSinglyLinkedList creates an empty singly linked list.
//
//	head     tail
//	  ↓        ↓
//	 nil      nil      size = 0
func NewSinglyLinkedList[T any]() *SinglyLinkedList[T] {
	return &SinglyLinkedList[T]{}
}

// Empty returns true if the list has no elements.
//
//	Empty list:              Non-empty list:
//	head = nil               head ──► [A] ──► [B] ──► nil
//	tail = nil               tail ──────────────┘
//	size = 0                 size = 2
//
//	Empty() → true           Empty() → false
func (l *SinglyLinkedList[T]) Empty() bool {
	return l.size == 0 && l.head == nil && l.tail == nil
}

// Size returns the number of elements in the list.
//
//	head                 tail
//	  ↓                    ↓
//	┌───┐     ┌───┐     ┌───┐
//	│ A │────►│ B │────►│ C │────► nil
//	└───┘     └───┘     └───┘
//
//	Size() → 3
//
// complexity:
//   - time : O(1)
//   - space: O(1)
func (l *SinglyLinkedList[T]) Size() int {
	return l.size
}

// Head returns the first element without removing it.
//
//	head                 tail
//	  ↓                    ↓
//	┌───┐     ┌───┐     ┌───┐
//	│ A │────►│ B │────►│ C │────► nil
//	└───┘     └───┘     └───┘
//	  ↑
//	Head() → A
//
// complexity:
//   - time : O(1)
//   - space: O(1)
//
// Panics if the list is empty.
func (l *SinglyLinkedList[T]) Head() T {
	if v, ok := l.TryHead(); !ok {
		panic("SinglyLinkedList.Head: is empty list")
	} else {
		return v
	}
}

// TryHead attempts to return the first element.
// Returns (value, true) on success, or (zero, false) if empty.
//
// complexity:
//   - time : O(1)
//   - space: O(1)
func (l *SinglyLinkedList[T]) TryHead() (T, bool) {
	if l.Empty() {
		return generics.ZeroValue[T](), false
	}
	return l.head.data, true
}

// Tail returns the last element without removing it.
//
//	head                 tail
//	  ↓                    ↓
//	┌───┐     ┌───┐     ┌───┐
//	│ A │────►│ B │────►│ C │────► nil
//	└───┘     └───┘     └───┘
//	                     ↑
//	                Tail() → C
//
// complexity:
//   - time : O(1)
//   - space: O(1)
//
// Panics if the list is empty.
func (l *SinglyLinkedList[T]) Tail() T {
	if v, ok := l.TryTail(); !ok {
		panic("SinglyLinkedList.Tail: is empty list")
	} else {
		return v
	}
}

// TryTail attempts to return the last element.
// Returns (value, true) on success, or (zero, false) if empty.
//
// complexity:
//   - time : O(1)
//   - space: O(1)
func (l *SinglyLinkedList[T]) TryTail() (T, bool) {
	if l.Empty() {
		return generics.ZeroValue[T](), false
	}
	return l.tail.data, true
}

// Append adds an element to the back of the list.
//
//	Before Append(D):
//
//	head                 tail
//	  ↓                    ↓
//	┌───┐     ┌───┐     ┌───┐
//	│ A │────►│ B │────►│ C │────► nil
//	└───┘     └───┘     └───┘
//
//	After Append(D):
//
//	head                           tail
//	  ↓                              ↓
//	┌───┐     ┌───┐     ┌───┐     ┌───┐
//	│ A │────►│ B │────►│ C │────►│ D │────► nil
//	└───┘     └───┘     └───┘     └───┘
//
// complexity:
//   - time : O(1)
//   - space: O(1)
//
// SCORE: 15
func (l *SinglyLinkedList[T]) Append(data T) {
	// hint: 1) create new node with NewUnaryNode(data, nil)
	//       2) if empty: set head and tail to new node
	//       3) else: set tail.next = new node, then tail = new node
	//       4) increment size

	newNode := NewUnaryNode(data, nil)
	if l.Empty() {
		l.head = newNode
		l.tail = newNode
	} else {
		l.tail.next = newNode
		l.tail = newNode
	}
	l.size++
}

// Prepend adds an element to the front of the list.
//
//	Before Prepend(Z):
//
//	head                 tail
//	  ↓                    ↓
//	┌───┐     ┌───┐     ┌───┐
//	│ A │────►│ B │────►│ C │────► nil
//	└───┘     └───┘     └───┘
//
//	After Prepend(Z):
//
//	head                           tail
//	  ↓                              ↓
//	┌───┐     ┌───┐     ┌───┐     ┌───┐
//	│ Z │────►│ A │────►│ B │────►│ C │────► nil
//	└───┘     └───┘     └───┘     └───┘
//
// complexity:
//   - time : O(1)
//   - space: O(1)
//
// SCORE: 10
func (l *SinglyLinkedList[T]) Prepend(data T) {
	// hint: 1) create new node with NewUnaryNode(data, l.head)
	//       2) if empty: set tail = new node
	//       3) set head = new node
	//       4) increment size

	newNode := NewUnaryNode(data, l.head)
	if l.Empty() {
		l.tail = newNode
	}
	l.head = newNode
	l.size++
}

// Pop removes and returns the last element.
//
//	Before Pop():
//
//	head                 tail
//	  ↓                    ↓
//	┌───┐     ┌───┐     ┌───┐
//	│ A │────►│ B │────►│ C │────► nil
//	└───┘     └───┘     └───┘
//
//	After Pop():
//
//	head         tail
//	  ↓            ↓
//	┌───┐     ┌───┐
//	│ A │────►│ B │────► nil
//	└───┘     └───┘
//
//	Pop() → C (removed element)
//
// complexity:
//   - time : O(n) - must traverse to find second-to-last node
//   - space: O(1)
//     where n is the number of elements.
//
// Note: This is O(n) because singly linked lists cannot go backward.
// Use DoublyLinkedList if you need O(1) Pop.
//
// Panics if the list is empty.
func (l *SinglyLinkedList[T]) Pop() T {
	if data, ok := l.TryPop(); !ok {
		panic("SinglyLinkedList.Pop: is empty list")
	} else {
		return data
	}
}

// TryPop attempts to remove and return the last element.
// Returns (value, true) on success, or (zero, false) if empty.
//
// This is a non-panicking version of Pop.
//
// complexity:
//   - time : O(n)
//   - space: O(1)
//
// SCORE: 15
func (l *SinglyLinkedList[T]) TryPop() (T, bool) {
	// hint: 1) if empty, return (zero, false)
	//       2) if size == 1, save data, call reset(), return (data, true)
	//       3) else: traverse to find second-to-last node (node.next == tail)
	//       4) save tail.data, set node.next = nil, tail = node
	//       5) decrement size, return (saved, true)

	if l.Empty() {
		var sl T
		return sl, false
	}
	data := l.tail.data
	if l.size == 1 {
		l.reset()
		return data, true
	} else {
		cur := l.head
		for cur.next != l.tail {
			cur = cur.next
		}
		l.tail = cur
		cur.next = nil
		l.size--
	}
	return data, true
}

// Shift removes and returns the first element.
//
//	Before Shift():
//
//	head                 tail
//	  ↓                    ↓
//	┌───┐     ┌───┐     ┌───┐
//	│ A │────►│ B │────►│ C │────► nil
//	└───┘     └───┘     └───┘
//
//	After Shift():
//
//	head         tail
//	  ↓            ↓
//	┌───┐     ┌───┐
//	│ B │────►│ C │────► nil
//	└───┘     └───┘
//
//	Shift() → A (removed element)
//
// complexity:
//   - time : O(1)
//   - space: O(1)
//
// Panics if the list is empty.
func (l *SinglyLinkedList[T]) Shift() T {
	if data, ok := l.TryShift(); !ok {
		panic("SinglyLinkedList.Shift: is empty list")
	} else {
		return data
	}
}

// TryShift attempts to remove and return the first element.
// Returns (value, true) on success, or (zero, false) if empty.
//
// This is a non-panicking version of Shift.
//
// complexity:
//   - time : O(1)
//   - space: O(1)
//
// SCORE: 10
func (l *SinglyLinkedList[T]) TryShift() (T, bool) {
	// hint: 1) if empty, return (zero, false)
	//       2) save head.data
	//       3) if size == 1: call reset()
	//       4) else: head = head.next, decrement size
	//       5) return (saved, true)

	if l.Empty() {
		var sl T
		return sl, false
	}
	data := l.head.data
	if l.size == 1 {
		l.reset()
	} else {
		l.head = l.head.next
		l.size--
	}
	return data, true
}

// Iter iterates over all elements from front to back.
//
//	head                 tail
//	  ↓                    ↓
//	┌───┐     ┌───┐     ┌───┐
//	│ A │────►│ B │────►│ C │────► nil
//	└───┘     └───┘     └───┘
//	  ↓         ↓         ↓
//	 1st       2nd       3rd       ← iteration order
//
// Example:
//
//	for value := range list.Iter {
//	    fmt.Println(value)  // prints A, B, C
//	}
//
// complexity:
//   - time : O(n)
//   - space: O(1)
func (l *SinglyLinkedList[T]) Iter(yield func(T) bool) {
	l.iterForward(func(u *UnaryNode[T]) bool { return yield(u.data) })
}

// IterBackward iterates over all elements from back to front.
//
//	head                 tail
//	  ↓                    ↓
//	┌───┐     ┌───┐     ┌───┐
//	│ A │────►│ B │────►│ C │────► nil
//	└───┘     └───┘     └───┘
//	  ↓         ↓         ↓
//	 3rd       2nd       1st       ← iteration order
//
// complexity:
//   - time : O(n)
//   - space: O(n) - creates a temporary reversed copy
//
// Note: For singly linked lists, backward iteration requires extra work
// since nodes don't have prev pointers. Use DoublyLinkedList for efficient
// backward iteration.
func (l *SinglyLinkedList[T]) IterBackward(yield func(T) bool) {
	l.iterBackward(func(u *UnaryNode[T]) bool { return yield(u.data) })
}

// Enum iterates over all elements with their indices from front to back.
//
//	head                 tail
//	  ↓                    ↓
//	┌───┐     ┌───┐     ┌───┐
//	│ A │────►│ B │────►│ C │────► nil
//	└───┘     └───┘     └───┘
//	  0         1         2        ← indices
//
// Example:
//
//	for index, value := range list.Enum {
//	    fmt.Printf("%d: %v\n", index, value)
//	}
//
// complexity:
//   - time : O(n)
//   - space: O(1)
func (l *SinglyLinkedList[T]) Enum(yield func(int, T) bool) {
	i := 0
	for v := range l.Iter {
		if !yield(i, v) {
			break
		}
		i++
	}
}

// EnumBackward iterates over all elements with their indices from back to front.
//
//	head                 tail
//	  ↓                    ↓
//	┌───┐     ┌───┐     ┌───┐
//	│ A │────►│ B │────►│ C │────► nil
//	└───┘     └───┘     └───┘
//	  0         1         2        ← indices
//	  ↓         ↓         ↓
//	 3rd       2nd       1st       ← enumeration order
//
// complexity:
//   - time : O(n)
//   - space: O(n) - creates a temporary reversed copy
func (l *SinglyLinkedList[T]) EnumBackward(yield func(int, T) bool) {
	i := l.Size() - 1
	for v := range l.IterBackward {
		if !yield(i, v) {
			break
		}
		i--
	}
}

// SCORE: 10
func (l *SinglyLinkedList[T]) iterForward(yield func(*UnaryNode[T]) bool) {
	// hint: p := l.head; loop while p != nil; call yield(p); p = p.next

	p := l.head
	for p != nil {
		if !yield(p) {
			return
		}
		p = p.next
	}
}

// SCORE: 10
func (l *SinglyLinkedList[T]) iterBackward(yield func(*UnaryNode[T]) bool) {
	// hint: singly linked list has no prev pointer, so:
	//       option 1: collect all nodes into a slice, iterate in reverse
	//       option 2: use recursion (traverse to end, yield on the way back)

	var node []*UnaryNode[T]
	// Traverse
	for cur := l.head; cur != nil; cur = cur.next {
		node = append(node, cur)
	}
	// Iterate reverse
	for i := len(node) - 1; i >= 0; i-- {
		if !yield(node[i]) {
			return
		}
	}
}

// Get retrieves the element at the given index.
//
//	head                           tail       ← cursors
//	  ↓                              ↓
//	┌───┐     ┌───┐     ┌───┐     ┌───┐
//	│ A │────►│ B │────►│ C │────►│ D │────► nil           ← values
//	└───┘     └───┘     └───┘     └───┘
//	  0         1         2         3                      ← indices
//
//	Get(0) → A
//	Get(2) → C
//	Get(3) → D
//
// complexity:
//   - time : O(k) where k is the given index
//   - space: O(1)
//
// Panics if:
//   - The list is empty
//   - index < 0 or index >= Size()
func (l *SinglyLinkedList[T]) Get(index int) T {
	if v, ok := l.TryGet(index); !ok {
		if l.Empty() {
			panic("SinglyLinkedList.Get: is empty list")
		}
		panic("SinglyLinkedList.Get: index out of range")
	} else {
		return v
	}
}

// TryGet attempts to retrieve the element at the given index.
// Returns (value, true) on success, or (zero, false) if index is out of bounds.
//
// complexity:
//   - time : O(n)
//   - space: O(1)
//
// SCORE: 10
func (l *SinglyLinkedList[T]) TryGet(index int) (T, bool) {
	// hint: 1) check bounds (index < 0 || index >= size), return false
	//       2) traverse from head, counting until you reach index
	//       3) return (node.data, true)

	// Check bounds
	if index < 0 || index >= l.size {
		var sl T
		return sl, false
	}
	// Traverse until reach index
	cur := l.head
	for i := 0; i < index; i++ {
		cur = cur.next
	}
	return cur.data, true
}

// Set updates the element at the given index.
//
//	Before Set(2, X):
//
//	head                           tail       ← cursors
//	  ↓                              ↓
//	┌───┐     ┌───┐     ┌───┐     ┌───┐
//	│ A │────►│ B │────►│ C │────►│ D │────► nil           ← values
//	└───┘     └───┘     └───┘     └───┘
//	  0         1         2         3                      ← indices
//
//	After Set(2, X):
//
//	head                           tail       ← cursors
//	  ↓                              ↓
//	┌───┐     ┌───┐     ┌───┐     ┌───┐
//	│ A │────►│ B │────►│ X │────►│ D │────► nil           ← values
//	└───┘     └───┘     └───┘     └───┘
//	  0         1         2         3                      ← indices
//	                      ↑
//	                  updated
//
// complexity:
//   - time : O(k) where k is the given index
//   - space: O(1)
//
// Panics if:
//   - The list is empty
//   - index < 0 or index >= Size()
func (l *SinglyLinkedList[T]) Set(index int, data T) {
	if !l.TrySet(index, data) {
		if l.Empty() {
			panic("SinglyLinkedList.Set: is empty list")
		}
		panic("SinglyLinkedList.Set: index out of range")
	}
}

// TrySet attempts to update the element at the given index.
// Returns true on success, false if index is out of bounds.
//
// complexity:
//   - time : O(n)
//   - space: O(1)
//
// SCORE: 5
func (l *SinglyLinkedList[T]) TrySet(index int, data T) bool {
	// hint: 1) check bounds, return false if invalid
	//       2) traverse to node at index
	//       3) update node.data = data
	//       4) return true

	// Bound check
	if l.Empty() || index < 0 || index >= l.Size() {
		return false
	}
	// if index is 0
	if index == 0 {
		l.head.data = data
		return true
	}
	// Traverse index
	cur := l.head
	for i := 0; i < index; i++ {
		if cur == nil {
			return false
		}
		cur = cur.next
	}
	// update data
	cur.data = data
	return true
}

// String returns the string representation of the list.
//
//	head                           tail
//	  ↓                              ↓
//	┌───┐     ┌───┐     ┌───┐     ┌───┐
//	│ 1 │────►│ 2 │────►│ 3 │────►│ 4 │────► nil
//	└───┘     └───┘     └───┘     └───┘
//
//	String() → "[1 2 3 4]"
//
// complexity:
//   - time : O(n)
//   - space: O(n)
func (l *SinglyLinkedList[T]) String() string {
	return sequence.String(l.Iter)
}

// Remove deletes and returns the element at the given index.
//
//	Before Remove(2):
//
//	head                           tail       ← cursors
//	  ↓                              ↓
//	┌───┐     ┌───┐     ┌───┐     ┌───┐
//	│ A │────►│ B │────►│ C │────►│ D │────► nil           ← values
//	└───┘     └───┘     └───┘     └───┘
//	  0         1         2         3                      ← indices
//	                      ↑
//	                   remove
//
//	After Remove(2):
//
//	head                 tail                 ← cursors
//	  ↓                    ↓
//	┌───┐     ┌───┐     ┌───┐
//	│ A │────►│ B │────►│ D │────► nil                     ← values
//	└───┘     └───┘     └───┘
//	  0         1         2                                ← indices
//
//	Remove(2) → C (removed element)
//
// complexity:
//   - time : O(k) where k is the given index
//   - space: O(1)
//
// Panics if:
//   - index < 0 or index >= Size()
func (l *SinglyLinkedList[T]) Remove(index int) T {
	if v, ok := l.TryRemove(index); !ok {
		panic("SinglyLinkedList.Remove: index out of range")
	} else {
		return v
	}
}

// TryRemove attempts to remove the element at the given index.
// Returns (value, true) on success, or (zero, false) if index is out of bounds.
//
// complexity:
//   - time : O(n)
//   - space: O(1)
//
// SCORE: 10
func (l *SinglyLinkedList[T]) TryRemove(index int) (T, bool) {
	// hint: 1) check bounds, return (zero, false) if invalid
	//       2) if index == 0: return TryShift()
	//       3) if index == size-1: return TryPop()
	//       4) else: traverse to node at (index-1), rewire: prev.next = prev.next.next
	//       5) decrement size, return (removed.data, true)

	// Check bounds
	if index < 0 || index >= l.size {
		return generics.ZeroValue[T](), false
	}
	// Remove index if at head/tail
	if index == 0 {
		return l.TryShift()
	}
	if index == l.size-1 {
		return l.TryPop()
	}
	// Remove mid index
	cur := l.head
	for i := 0; i < index-1; i++ {
		cur = cur.next
	}
	removedNode := cur.next
	val := removedNode.data
	cur.next = removedNode.next
	l.size--
	return val, true
}

// Insert adds an element at the given index.
// Elements at and after the index shift to higher indices.
//
//	Before Insert(2, X):
//
//	head                 tail                 ← cursors
//	  ↓                    ↓
//	┌───┐     ┌───┐     ┌───┐
//	│ A │────►│ B │────►│ C │────► nil                     ← values
//	└───┘     └───┘     └───┘
//	  0         1         2                                ← indices
//	                ↑
//	          insert here
//
//	After Insert(2, X):
//
//	head                           tail       ← cursors
//	  ↓                              ↓
//	┌───┐     ┌───┐     ┌───┐     ┌───┐
//	│ A │────►│ B │────►│ X │────►│ C │────► nil           ← values
//	└───┘     └───┘     └───┘     └───┘
//	  0         1         2         3                      ← indices
//	                      ↑
//	                  inserted
//
// Special cases:
//   - Insert(0, X) is equivalent to Prepend(X)
//   - Insert(Size(), X) is equivalent to Append(X)
//
// complexity:
//   - time : O(k) where k is the given index
//   - space: O(1)
//
// Panics if:
//   - index < 0 or index > Size()
//
// SCORE: 5
func (l *SinglyLinkedList[T]) Insert(index int, data T) {
	// hint: 1) if index == 0: Prepend(data); return
	//       2) if index == Size(): Append(data); return
	//       3) check bounds, panic if invalid
	//       4) traverse to node at (index-1)
	//       5) create new node with next = prev.next
	//       6) prev.next = new node
	//       7) increment size

	if index == 0 {
		l.Prepend(data)
		return
	}
	if index == l.Size() {
		l.Append(data)
		return
	}
	if index < 0 || index >= l.size {
		panic("SinglyLinkedList.Insert: index out of range")
	}
	cur := l.head
	for i := 0; i < index-1; i++ {
		cur = cur.next
	}
	newNode := NewUnaryNode(data, cur.next)
	cur.next = newNode
	l.size++
}

func (l *SinglyLinkedList[T]) checkBounds(index int) {
	if index < 0 || index >= l.Size() {
		panic("SinglyLinkedList.checkBounds: index out of range")
	}
}

func (l *SinglyLinkedList[T]) reset() {
	l.head = nil
	l.tail = nil
	l.size = 0
}
