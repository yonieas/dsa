// Package queue provides a Queue (FIFO) data structure implementation.
//
// # What is a Queue?
//
// A Queue is a First-In-First-Out (FIFO) collection. Elements are added at
// the rear and removed from the front, like a line at a ticket counter:
// first person in line gets served first.
//
// # Why Use Queues?
//
// Queues model any situation where requests are processed in order of arrival.
// Operating systems use queues to schedule processes. Network routers queue
// packets. Web servers queue incoming requests. Event-driven programs queue
// events.
//
// Breadth-First Search (BFS) uses a queue to explore level by level. This
// naturally finds shortest paths in unweighted graphs.
//
// # Operations
//
//	Enqueue: Add element to the rear
//	Dequeue: Remove and return element from the front
//	Front:   Return front element without removing
//	Empty:   Check if queue has no elements
//
// # Implementation
//
// This queue is backend-agnostic. By default it uses a doubly linked list for
// O(1) enqueue and dequeue. You can provide a dynamic array backend, though
// naive array-based queues have O(n) dequeue. Consider a circular buffer for
// array-based O(1) operations.
//
// # Complexity
//
//	Enqueue/Dequeue/Front: O(1)
//	Space:                 O(n)
//
// # Further Reading
//
// CLRS "Introduction to Algorithms", Chapter 10.1.
// Sedgewick "Algorithms", Section 1.3.
// https://en.wikipedia.org/wiki/Queue_(abstract_data_type)
package queue

import (
	"github.com/josestg/dsa/adt"
	"github.com/josestg/dsa/linkedlist"
)

// Backend defines the interface required by the queue's underlying storage.
type Backend[E any] interface {
	adt.Sizer
	adt.Emptier
	adt.Header[E]
	adt.Tailer[E]
	adt.Shifter[E]
	adt.Appender[E]
	adt.Iterator[E]
	adt.Stringer
}

// Queue is a FIFO (First-In-First-Out) data structure.
//
//	           front                 rear
//	             ↓                     ↓
//	Dequeue ←  ┌───┬───┬───┬───┬───┐  ← Enqueue
//	           │ A │ B │ C │ D │ E │
//	           └───┴───┴───┴───┴───┘
//	             ↑
//	        Peek() → A
//
// Operations:
//   - Enqueue: Add element to the rear
//   - Dequeue: Remove and return the front element
//   - Peek: View the front element without removing
type Queue[E any] struct {
	b Backend[E]
}

// New creates an empty Queue using DoublyLinkedList as backend.
//
//	front     rear
//	  ↓        ↓
//	┌──────────────┐
//	│    empty     │
//	└──────────────┘
//
// complexity:
//   - time : O(1)
//   - space: O(1)
func New[E any]() *Queue[E] {
	return NewWith[E](linkedlist.NewDoublyLinkedList[E]())
}

// NewWith creates a Queue using a custom backend.
// This allows using different underlying data structures.
func NewWith[E any](b Backend[E]) *Queue[E] {
	return &Queue[E]{b: b}
}

// Empty returns true if the queue has no elements.
//
//	Empty queue:              Non-empty queue:
//	front   rear              front             rear
//	  ↓       ↓                 ↓                 ↓
//	┌───────────┐             ┌───┬───┬───┬───┐
//	│   empty   │             │ A │ B │ C │ D │
//	└───────────┘             └───┴───┴───┴───┘
//
//	Empty() → true            Empty() → false
//
// complexity:
//   - time : O(1)
//   - space: O(1)
func (q *Queue[E]) Empty() bool {
	q.ensureBackend()
	return q.b.Empty()
}

// Size returns the number of elements in the queue.
//
//	front                 rear
//	  ↓                     ↓
//	┌───┬───┬───┬───┬───┐
//	│ A │ B │ C │ D │ E │
//	└───┴───┴───┴───┴───┘
//
//	Size() → 5
//
// complexity:
//   - time : O(1)
//   - space: O(1)
func (q *Queue[E]) Size() int {
	q.ensureBackend()
	return q.b.Size()
}

// Peek returns the front element without removing it.
//
//	front                 rear
//	  ↓                     ↓
//	┌───┬───┬───┬───┬───┐
//	│ A │ B │ C │ D │ E │
//	└───┴───┴───┴───┴───┘
//	  ↑
//	Peek() → A
//
// The front element is the oldest element (first to be dequeued).
// The queue remains unchanged after Peek.
//
// complexity:
//   - time : O(1)
//   - space: O(1)
//
// Panics if the queue is empty.
func (q *Queue[E]) Peek() E {
	if v, ok := q.TryPeek(); !ok {
		panic("queue.Peek: queue is empty")
	} else {
		return v
	}
}

// TryPeek attempts to return the front element without removing it.
// Returns (value, true) on success, or (zero, false) if empty.
//
// complexity:
//   - time : O(1)
//   - space: O(1)
func (q *Queue[E]) TryPeek() (E, bool) {
	q.ensureBackend()
	if q.Empty() {
		var zero E
		return zero, false
	}
	return q.b.Head(), true
}

// Enqueue adds an element to the rear of the queue.
//
//	Before Enqueue(F):
//
//	front                 rear
//	  ↓                     ↓
//	┌───┬───┬───┬───┬───┐
//	│ A │ B │ C │ D │ E │
//	└───┴───┴───┴───┴───┘
//
//	After Enqueue(F):
//
//	front                       rear
//	  ↓                           ↓
//	┌───┬───┬───┬───┬───┬───┐
//	│ A │ B │ C │ D │ E │ F │
//	└───┴───┴───┴───┴───┴───┘
//	                      ↑
//	                   added
//
// complexity:
//   - time : O(1)
//   - space: O(1)
func (q *Queue[E]) Enqueue(data E) {
	q.ensureBackend()
	q.b.Append(data)
}

// Dequeue removes and returns the front element.
//
//	Before Dequeue():
//
//	front                 rear
//	  ↓                     ↓
//	┌───┬───┬───┬───┬───┐
//	│ A │ B │ C │ D │ E │
//	└───┴───┴───┴───┴───┘
//
//	After Dequeue():
//
//	front             rear
//	  ↓                 ↓
//	┌───┬───┬───┬───┐
//	│ B │ C │ D │ E │
//	└───┴───┴───┴───┘
//
//	Dequeue() → A (the removed element)
//
// The oldest element (first enqueued) is removed first (FIFO).
//
// complexity:
//   - time : O(1)
//   - space: O(1)
//
// Panics if the queue is empty (queue underflow).
func (q *Queue[E]) Dequeue() E {
	if v, ok := q.TryDequeue(); !ok {
		panic("queue.Dequeue: queue underflow")
	} else {
		return v
	}
}

// TryDequeue attempts to remove and return the front element.
// Returns (value, true) on success, or (zero, false) if empty.
//
// complexity:
//   - time : O(1)
//   - space: O(1)
func (q *Queue[E]) TryDequeue() (E, bool) {
	q.ensureBackend()
	if q.Empty() {
		var zero E
		return zero, false
	}
	return q.b.Shift(), true
}

// String returns the string representation of the queue.
//
//	front                 rear
//	  ↓                     ↓
//	┌───┬───┬───┬───┬───┐
//	│ 1 │ 2 │ 3 │ 4 │ 5 │
//	└───┴───┴───┴───┴───┘
//
//	String() → "[1 2 3 4 5]"
//
// The leftmost element is the front of the queue (next to be dequeued).
//
// complexity:
//   - time : O(n)
//   - space: O(n)
func (q *Queue[E]) String() string {
	q.ensureBackend()
	return q.b.String()
}

// Iter iterates over all elements from front to rear.
//
//	front                 rear
//	  ↓                     ↓
//	┌───┬───┬───┬───┬───┐
//	│ A │ B │ C │ D │ E │
//	└───┴───┴───┴───┴───┘
//	  ↓   ↓   ↓   ↓   ↓
//	 1st 2nd 3rd 4th 5th   ← iteration order
//
//	for v := range queue.Iter {
//	    fmt.Println(v)  // prints A, B, C, D, E
//	}
//
// complexity:
//   - time : O(n)
//   - space: O(1)
func (q *Queue[E]) Iter(yield func(E) bool) {
	q.ensureBackend()
	q.b.Iter(yield)
}

func (q *Queue[E]) ensureBackend() {
	if q.b == nil {
		q.b = linkedlist.NewDoublyLinkedList[E]()
	}
}
