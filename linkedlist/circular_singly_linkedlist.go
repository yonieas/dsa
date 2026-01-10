package linkedlist

import (
	"github.com/josestg/dsa/internal/generics"
	"github.com/josestg/dsa/sequence"
)

// CircularSinglyLinkedList is a singly linked list where the tail points back to the head.
// This creates a circular structure with no nil terminators.
//
//	  ┌──────────────────────────────────┐
//	  │                                  │
//	  ▼                                  │
//	┌───┐     ┌───┐     ┌───┐     ┌───┐  │
//	│ A │────►│ B │────►│ C │────►│ D │──┘
//	└───┘     └───┘     └───┘     └───┘
//	  ↑
//	 head
//
// Use cases:
//   - Round-robin scheduling
//   - Circular buffers
//   - Repeating playlists
//   - Turn-based games
//
// With only a head pointer:
//   - Append: O(n) - must traverse to find tail
//   - Prepend: O(n) - must update tail's next pointer
//   - Shift: O(n) - must update tail's next pointer
//
// Note: For O(1) operations at both ends, use CircularDoublyLinkedList.
type CircularSinglyLinkedList[T any] struct {
	head *UnaryNode[T]
	size int
}

// NewCircularSinglyLinkedList creates an empty circular singly linked list.
func NewCircularSinglyLinkedList[T any]() *CircularSinglyLinkedList[T] {
	return &CircularSinglyLinkedList[T]{}
}

// Size returns the number of elements.
func (l *CircularSinglyLinkedList[T]) Size() int {
	return l.size
}

// Empty returns true if the list has no elements.
func (l *CircularSinglyLinkedList[T]) Empty() bool {
	return l.size == 0
}

// Head returns the first element without removing it.
// Panics if the list is empty.
func (l *CircularSinglyLinkedList[T]) Head() T {
	if v, ok := l.TryHead(); !ok {
		panic("CircularSinglyLinkedList.Head: list is empty")
	} else {
		return v
	}
}

// TryHead attempts to return the first element.
// Returns (value, true) on success, or (zero, false) if empty.
func (l *CircularSinglyLinkedList[T]) TryHead() (T, bool) {
	if l.Empty() {
		return generics.ZeroValue[T](), false
	}
	return l.head.data, true
}

// Tail returns the last element without removing it.
// Panics if the list is empty.
func (l *CircularSinglyLinkedList[T]) Tail() T {
	if v, ok := l.TryTail(); !ok {
		panic("CircularSinglyLinkedList.Tail: list is empty")
	} else {
		return v
	}
}

// TryTail attempts to return the last element.
// Returns (value, true) on success, or (zero, false) if empty.
func (l *CircularSinglyLinkedList[T]) TryTail() (T, bool) {
	if l.Empty() {
		return generics.ZeroValue[T](), false
	}
	return l.getTail().data, true
}

// Append adds an element to the back of the list.
//
//	Before Append(E):
//	        ┌──────────────────────────┐
//	        ▼                          │
//	      ┌───┐     ┌───┐     ┌───┐    │
//	      │ A │────►│ B │────►│ C │────┘
//	      └───┘     └───┘     └───┘
//
//	After Append(E):
//	        ┌────────────────────────────────────┐
//	        ▼                                    │
//	      ┌───┐     ┌───┐     ┌───┐     ┌───┐    │
//	      │ A │────►│ B │────►│ C │────►│ E │────┘
//	      └───┘     └───┘     └───┘     └───┘
//
// complexity: O(n)
func (l *CircularSinglyLinkedList[T]) Append(data T) {
	n := NewUnaryNode(data, nil)
	if l.Empty() {
		n.next = n
		l.head = n
	} else {
		tail := l.getTail()
		n.next = l.head
		tail.next = n
	}
	l.size++
}

// Prepend adds an element to the front of the list.
//
// complexity: O(n)
func (l *CircularSinglyLinkedList[T]) Prepend(data T) {
	n := NewUnaryNode(data, nil)
	if l.Empty() {
		n.next = n
		l.head = n
	} else {
		tail := l.getTail()
		n.next = l.head
		tail.next = n
		l.head = n
	}
	l.size++
}

// Pop removes and returns the last element.
//
// complexity: O(n)
//
// Panics if the list is empty.
func (l *CircularSinglyLinkedList[T]) Pop() T {
	if data, ok := l.TryPop(); !ok {
		panic("CircularSinglyLinkedList.Pop: list is empty")
	} else {
		return data
	}
}

// TryPop attempts to remove and return the last element.
// Returns (value, true) on success, or (zero, false) if empty.
//
// This is a non-panicking version of Pop.
//
// complexity: O(n)
func (l *CircularSinglyLinkedList[T]) TryPop() (T, bool) {
	if l.Empty() {
		return generics.ZeroValue[T](), false
	}
	if l.size == 1 {
		data := l.head.data
		l.reset()
		return data, true
	}
	prev := l.head
	for prev.next.next != l.head {
		prev = prev.next
	}
	data := prev.next.data
	prev.next = l.head
	l.size--
	return data, true
}

// Shift removes and returns the first element.
//
// complexity: O(n)
//
// Panics if the list is empty.
func (l *CircularSinglyLinkedList[T]) Shift() T {
	if data, ok := l.TryShift(); !ok {
		panic("CircularSinglyLinkedList.Shift: list is empty")
	} else {
		return data
	}
}

// TryShift attempts to remove and return the first element.
// Returns (value, true) on success, or (zero, false) if empty.
//
// This is a non-panicking version of Shift.
//
// complexity: O(n)
func (l *CircularSinglyLinkedList[T]) TryShift() (T, bool) {
	if l.Empty() {
		return generics.ZeroValue[T](), false
	}
	data := l.head.data
	if l.size == 1 {
		l.reset()
		return data, true
	}
	tail := l.getTail()
	oldHead := l.head
	l.head = l.head.next
	tail.next = l.head
	oldHead.next = nil
	l.size--
	return data, true
}

// Rotate moves the head pointer n positions forward (positive) or backward (negative).
//
//	Before Rotate(1):
//	        ┌──────────────────────────┐
//	        ▼                          │
//	      ┌───┐     ┌───┐     ┌───┐    │
//	      │ A │────►│ B │────►│ C │────┘
//	      └───┘     └───┘     └───┘
//	        ↑
//	       head
//
//	After Rotate(1):
//	        ┌──────────────────────────┐
//	        │                          ▼
//	        │     ┌───┐     ┌───┐     ┌───┐
//	        └─────│ A │────►│ B │────►│ C │
//	              └───┘     └───┘     └───┘
//	                ↑                   ↑
//	               tail               head
//
// complexity: O(n) for negative rotation, O(k) for positive where k = n mod size
func (l *CircularSinglyLinkedList[T]) Rotate(n int) {
	if l.Empty() || l.size == 1 {
		return
	}
	n = n % l.size
	if n == 0 {
		return
	}
	if n < 0 {
		n = l.size + n
	}
	for range n {
		l.head = l.head.next
	}
}

// Get retrieves the element at the given index.
//
// complexity: O(k) where k is the index
func (l *CircularSinglyLinkedList[T]) Get(index int) T {
	if v, ok := l.TryGet(index); !ok {
		if l.Empty() {
			panic("CircularSinglyLinkedList.Get: list is empty")
		}
		panic("CircularSinglyLinkedList.Get: index out of range")
	} else {
		return v
	}
}

// TryGet attempts to retrieve the element at the given index.
// Returns (value, true) on success, or (zero, false) if index is out of bounds.
func (l *CircularSinglyLinkedList[T]) TryGet(index int) (T, bool) {
	if l.Empty() || index < 0 || index >= l.Size() {
		return generics.ZeroValue[T](), false
	}
	data, ok := sequence.ValueAt(l.Iter, index)
	if !ok {
		return generics.ZeroValue[T](), false
	}
	return data, true
}

// Set updates the element at the given index.
//
// complexity: O(k) where k is the index
func (l *CircularSinglyLinkedList[T]) Set(index int, data T) {
	if !l.TrySet(index, data) {
		if l.Empty() {
			panic("CircularSinglyLinkedList.Set: list is empty")
		}
		panic("CircularSinglyLinkedList.Set: index out of range")
	}
}

// TrySet attempts to update the element at the given index.
// Returns true on success, false if index is out of bounds.
func (l *CircularSinglyLinkedList[T]) TrySet(index int, data T) bool {
	if l.Empty() || index < 0 || index >= l.Size() {
		return false
	}
	p := l.head
	for range index {
		p = p.next
	}
	p.data = data
	return true
}

// Insert adds an element at the given index.
//
// complexity: O(k) where k is the index
func (l *CircularSinglyLinkedList[T]) Insert(index int, data T) {
	if index == 0 {
		l.Prepend(data)
		return
	}
	if index == l.size {
		l.Append(data)
		return
	}
	l.checkBounds(index)
	p := l.head
	for i := 0; i < index-1; i++ {
		p = p.next
	}
	n := NewUnaryNode(data, p.next)
	p.next = n
	l.size++
}

// Remove deletes and returns the element at the given index.
//
// complexity: O(k) where k is the index
func (l *CircularSinglyLinkedList[T]) Remove(index int) T {
	if v, ok := l.TryRemove(index); !ok {
		panic("CircularSinglyLinkedList.Remove: index out of range")
	} else {
		return v
	}
}

// TryRemove attempts to remove the element at the given index.
// Returns (value, true) on success, or (zero, false) if index is out of bounds.
func (l *CircularSinglyLinkedList[T]) TryRemove(index int) (T, bool) {
	if index < 0 || index >= l.Size() {
		return generics.ZeroValue[T](), false
	}
	if index == 0 {
		return l.TryShift()
	}
	if index == l.size-1 {
		return l.TryPop()
	}
	p := l.head
	for i := 0; i < index-1; i++ {
		p = p.next
	}
	target := p.next
	p.next = target.next
	target.next = nil
	l.size--
	return target.data, true
}

// Iter iterates over all elements starting from head.
// Stops after visiting all elements once (does not loop infinitely).
func (l *CircularSinglyLinkedList[T]) Iter(yield func(T) bool) {
	if l.Empty() {
		return
	}
	p := l.head
	for range l.size {
		if !yield(p.data) {
			break
		}
		p = p.next
	}
}

// IterBackward iterates over all elements from tail to head.
//
// complexity:
//   - time : O(n)
//   - space: O(n) - creates a temporary reversed copy
//
// Note: For singly linked lists, backward iteration requires extra work
// since nodes don't have prev pointers. Use CircularDoublyLinkedList for
// efficient backward iteration.
func (l *CircularSinglyLinkedList[T]) IterBackward(yield func(T) bool) {
	if l.Empty() {
		return
	}
	// Collect values in reverse order
	values := make([]T, l.size)
	p := l.head
	for i := range l.size {
		values[i] = p.data
		p = p.next
	}
	// Yield in reverse order
	for i := l.size - 1; i >= 0; i-- {
		if !yield(values[i]) {
			break
		}
	}
}

// Enum iterates over all elements with their indices from front to back.
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
func (l *CircularSinglyLinkedList[T]) Enum(yield func(int, T) bool) {
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
// complexity:
//   - time : O(n)
//   - space: O(n) - creates a temporary reversed copy
func (l *CircularSinglyLinkedList[T]) EnumBackward(yield func(int, T) bool) {
	i := l.Size() - 1
	for v := range l.IterBackward {
		if !yield(i, v) {
			break
		}
		i--
	}
}

// String returns the string representation.
func (l *CircularSinglyLinkedList[T]) String() string {
	return sequence.String(l.Iter)
}

func (l *CircularSinglyLinkedList[T]) getTail() *UnaryNode[T] {
	if l.Empty() {
		return nil
	}
	p := l.head
	for p.next != l.head {
		p = p.next
	}
	return p
}

func (l *CircularSinglyLinkedList[T]) checkBounds(index int) {
	if index < 0 || index >= l.size {
		panic("CircularSinglyLinkedList: index out of range")
	}
}

func (l *CircularSinglyLinkedList[T]) reset() {
	if l.head != nil {
		l.head.next = nil
	}
	l.head = nil
	l.size = 0
}

// CircularIterator returns an infinite iterator that keeps cycling through elements.
// Use with caution - must break out of the loop manually.
//
//	count := 0
//	for v := range list.CircularIterator {
//	    fmt.Println(v)
//	    count++
//	    if count >= 10 { break }
//	}
func (l *CircularSinglyLinkedList[T]) CircularIterator(yield func(T) bool) {
	if l.Empty() {
		return
	}
	p := l.head
	for {
		if !yield(p.data) {
			break
		}
		p = p.next
	}
}

// Cycle advances the head to the next element and returns the old head's data.
// Wraps around to the beginning when reaching the end.
// Useful for round-robin iteration.
//
// complexity: O(1)
//
// Panics if the list is empty.
func (l *CircularSinglyLinkedList[T]) Cycle() T {
	if v, ok := l.TryCycle(); !ok {
		panic("CircularSinglyLinkedList.Cycle: list is empty")
	} else {
		return v
	}
}

// TryCycle attempts to return the current head and advance to the next element.
// Returns (value, true) on success, or (zero, false) if empty.
func (l *CircularSinglyLinkedList[T]) TryCycle() (T, bool) {
	if l.Empty() {
		return generics.ZeroValue[T](), false
	}
	data := l.head.data
	l.head = l.head.next
	return data, true
}
