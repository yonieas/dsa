// Package stack provides a Stack (LIFO) data structure implementation.
//
// # What is a Stack?
//
// A Stack is a Last-In-First-Out (LIFO) collection. Elements are added and
// removed from the same end, called the "top." Think of a stack of plates:
// you add plates on top and take plates from the top.
//
// # Why Use Stacks?
//
// Stacks are fundamental to how computers work. Every function call pushes a
// frame onto the call stack; returning pops it. Many algorithms naturally
// express themselves as "do something, then undo it" which is exactly what
// stacks model.
//
// Common use cases: parsing expressions and checking balanced brackets,
// implementing undo/redo, converting recursion to iteration, depth-first
// search traversal, and backtracking algorithms.
//
// # Operations
//
//	Push:  Add element to top
//	Pop:   Remove and return top element
//	Peek:  Return top element without removing
//	Empty: Check if stack has no elements
//
// # Implementation
//
// This stack is backend-agnostic. By default it uses a linked list (O(1)
// push/pop, no resizing). You can provide a dynamic array for better cache
// locality at the cost of occasional resize operations.
//
// # Complexity
//
//	Push/Pop/Peek: O(1)
//	Space:         O(n)
//
// # Further Reading
//
// CLRS "Introduction to Algorithms", Chapter 10.1.
// Sedgewick "Algorithms", Section 1.3.
// https://en.wikipedia.org/wiki/Stack_(abstract_data_type)
package stack

import (
	"github.com/josestg/dsa/adt"
	"github.com/josestg/dsa/linkedlist"
)

// Backend defines the interface required by the stack's underlying storage.
type Backend[E any] interface {
	adt.Sizer
	adt.Emptier
	adt.Tailer[E]
	adt.Popper[E]
	adt.Appender[E]
	adt.Iterator[E]
	adt.Stringer
}

// Stack is a LIFO (Last-In-First-Out) data structure.
//
//	   │ Push(D)
//	   ↓
//	┌─────┐
//	│  D  │ ← top (Peek returns D, Pop removes D)
//	├─────┤
//	│  C  │
//	├─────┤
//	│  B  │
//	├─────┤
//	│  A  │
//	└─────┘
//	(bottom)
//
// Operations:
//   - Push: Add element to the top
//   - Pop: Remove and return the top element
//   - Peek: View the top element without removing
type Stack[E any] struct {
	b Backend[E]
}

// New creates an empty Stack using DoublyLinkedList as backend.
//
//	┌─────┐
//	│     │  ← empty stack
//	└─────┘
//
// complexity:
//   - time : O(1)
//   - space: O(1)
func New[E any]() *Stack[E] {
	return NewWith[E](linkedlist.NewDoublyLinkedList[E]())
}

// NewWith creates a Stack using a custom backend.
// This allows using different underlying data structures.
func NewWith[E any](b Backend[E]) *Stack[E] {
	return &Stack[E]{b: b}
}

// Empty returns true if the stack has no elements.
//
//	Empty stack:             Non-empty stack:
//	    ┌─────┐                 ┌─────┐
//	    │     │                 │  C  │
//	    └─────┘                 ├─────┤
//	                            │  B  │
//	                            ├─────┤
//	                            │  A  │
//	                            └─────┘
//
//	Empty() → true           Empty() → false
//
// complexity:
//   - time : O(1)
//   - space: O(1)
func (s *Stack[E]) Empty() bool {
	s.ensureBackend()
	return s.b.Empty()
}

// Size returns the number of elements in the stack.
//
//	    ┌─────┐
//	    │  C  │
//	    ├─────┤
//	    │  B  │
//	    ├─────┤
//	    │  A  │
//	    └─────┘
//
//	Size() → 3
//
// complexity:
//   - time : O(1)
//   - space: O(1)
func (s *Stack[E]) Size() int {
	s.ensureBackend()
	return s.b.Size()
}

// Peek returns the top element without removing it.
//
//	┌─────┐
//	│  C  │ ← Peek() → C
//	├─────┤
//	│  B  │
//	├─────┤
//	│  A  │
//	└─────┘
//
// The top element is the most recently pushed element.
// The stack remains unchanged after Peek.
//
// complexity:
//   - time : O(1)
//   - space: O(1)
//
// Panics if the stack is empty.
func (s *Stack[E]) Peek() E {
	if v, ok := s.TryPeek(); !ok {
		panic("stack.Peek: stack is empty")
	} else {
		return v
	}
}

// TryPeek attempts to return the top element without removing it.
// Returns (value, true) on success, or (zero, false) if empty.
//
// complexity:
//   - time : O(1)
//   - space: O(1)
func (s *Stack[E]) TryPeek() (E, bool) {
	s.ensureBackend()
	if s.Empty() {
		var zero E
		return zero, false
	}
	return s.b.Tail(), true
}

// Push adds an element to the top of the stack.
//
//	Before Push(D):          After Push(D):
//
//	    ┌─────┐                 ┌─────┐
//	    │  C  │ ← top           │  D  │ ← new top
//	    ├─────┤                 ├─────┤
//	    │  B  │                 │  C  │
//	    ├─────┤                 ├─────┤
//	    │  A  │                 │  B  │
//	    └─────┘                 ├─────┤
//	                            │  A  │
//	                            └─────┘
//
// complexity:
//   - time : O(1)
//   - space: O(1)
func (s *Stack[E]) Push(data E) {
	s.ensureBackend()
	s.b.Append(data)
}

// Pop removes and returns the top element.
//
//	Before Pop():            After Pop():
//
//	    ┌─────┐                 ┌─────┐
//	    │  C  │ ← removed       │  B  │ ← new top
//	    ├─────┤                 ├─────┤
//	    │  B  │                 │  A  │
//	    ├─────┤                 └─────┘
//	    │  A  │
//	    └─────┘
//
//	Pop() → C (the removed element)
//
// complexity:
//   - time : O(1)
//   - space: O(1)
//
// Panics if the stack is empty (stack underflow).
func (s *Stack[E]) Pop() E {
	if v, ok := s.TryPop(); !ok {
		panic("stack.Pop: stack underflow")
	} else {
		return v
	}
}

// TryPop attempts to remove and return the top element.
// Returns (value, true) on success, or (zero, false) if empty.
//
// complexity:
//   - time : O(1)
//   - space: O(1)
func (s *Stack[E]) TryPop() (E, bool) {
	s.ensureBackend()
	if s.Empty() {
		var zero E
		return zero, false
	}
	return s.b.Pop(), true
}

// String returns the string representation of the stack.
//
//	    ┌─────┐
//	    │  3  │ ← top
//	    ├─────┤
//	    │  2  │
//	    ├─────┤
//	    │  1  │
//	    └─────┘
//
//	String() → "[1 2 3]"
//
// Note: The rightmost element is the top of the stack.
//
// complexity:
//   - time : O(n)
//   - space: O(n)
func (s *Stack[E]) String() string {
	s.ensureBackend()
	return s.b.String()
}

// Iter iterates over all elements from bottom to top.
//
//	    ┌─────┐
//	    │  C  │ ← top (visited last)
//	    ├─────┤
//	    │  B  │
//	    ├─────┤
//	    │  A  │ ← bottom (visited first)
//	    └─────┘
//
//	for v := range stack.Iter {
//	    fmt.Println(v)  // prints A, B, C
//	}
//
// complexity:
//   - time : O(n)
//   - space: O(1)
func (s *Stack[E]) Iter(yield func(E) bool) {
	s.ensureBackend()
	s.b.Iter(yield)
}

func (s *Stack[E]) ensureBackend() {
	if s.b == nil {
		s.b = linkedlist.NewDoublyLinkedList[E]()
	}
}
