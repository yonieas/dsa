package linkedlist

import (
	"github.com/josestg/dsa/internal/generics"
	"github.com/josestg/dsa/sequence"
)

// CircularDoublyLinkedList is a doubly linked list where tail.next points to head
// and head.prev points to tail.
//
//	  ┌────────────────────────────────────────────┐
//	  │                                            │
//	  ▼                                            │
//	┌───┐     ┌───┐     ┌───┐     ┌───┐            │
//	│ A │◄───►│ B │◄───►│ C │◄───►│ D │────────────┘
//	└───┘     └───┘     └───┘     └───┘
//	  ↑                             ↑
//	  └─────────────────────────────┘
//	 head                          tail
//
// All operations are O(1) at both ends:
//   - Append: O(1)
//   - Prepend: O(1)
//   - Pop: O(1)
//   - Shift: O(1)
//
// Use cases:
//   - LRU cache implementation
//   - Undo/redo with wraparound
//   - Circular navigation (photo galleries, carousels)
//   - Deque with circular iteration
type CircularDoublyLinkedList[T any] struct {
	head *BinaryNode[T]
	size int
}

// NewCircularDoublyLinkedList creates an empty circular doubly linked list.
func NewCircularDoublyLinkedList[T any]() *CircularDoublyLinkedList[T] {
	return &CircularDoublyLinkedList[T]{}
}

// Size returns the number of elements.
func (l *CircularDoublyLinkedList[T]) Size() int {
	return l.size
}

// Empty returns true if the list has no elements.
func (l *CircularDoublyLinkedList[T]) Empty() bool {
	return l.size == 0
}

// Head returns the first element without removing it.
// Panics if the list is empty.
func (l *CircularDoublyLinkedList[T]) Head() T {
	if v, ok := l.TryHead(); !ok {
		panic("CircularDoublyLinkedList.Head: list is empty")
	} else {
		return v
	}
}

// TryHead attempts to return the first element.
// Returns (value, true) on success, or (zero, false) if empty.
func (l *CircularDoublyLinkedList[T]) TryHead() (T, bool) {
	if l.Empty() {
		return generics.ZeroValue[T](), false
	}
	return l.head.data, true
}

// Tail returns the last element without removing it.
// Panics if the list is empty.
func (l *CircularDoublyLinkedList[T]) Tail() T {
	if v, ok := l.TryTail(); !ok {
		panic("CircularDoublyLinkedList.Tail: list is empty")
	} else {
		return v
	}
}

// TryTail attempts to return the last element.
// Returns (value, true) on success, or (zero, false) if empty.
func (l *CircularDoublyLinkedList[T]) TryTail() (T, bool) {
	if l.Empty() {
		return generics.ZeroValue[T](), false
	}
	return l.head.prev.data, true
}

// Append adds an element to the back of the list.
//
// complexity: O(1)
//
// SCORE: 15
func (l *CircularDoublyLinkedList[T]) Append(data T) {
	// hint: 1) create new node with next=nil, prev=nil
	//       2) if empty: set node.next = node, node.prev = node, head = node
	//       3) else: tail = head.prev; node.next = head; node.prev = tail
	//                tail.next = node; head.prev = node
	//       4) increment size
	node := NewBinaryNode(data, nil, nil)
	if l.Empty() {
		node.next = node
		node.prev = node
		l.head = node
	} else {
		tail := l.head.prev
		node.next = l.head
		node.prev = tail
		tail.next = node
		l.head.prev = node
	}
	l.size++
}

// Prepend adds an element to the front of the list.
//
// complexity: O(1)
//
// SCORE: 10
func (l *CircularDoublyLinkedList[T]) Prepend(data T) {
	// hint: 1) create new node
	//       2) if empty: set node.next = node, node.prev = node
	//       3) else: tail = head.prev; node.next = head; node.prev = tail
	//                head.prev = node; tail.next = node
	//       4) set head = node
	//       5) increment size
	node := NewBinaryNode(data, nil, nil)
	if l.Empty() {
		node.next = node
		node.prev = node
	} else {
		tail := l.head.prev
		node.next = l.head
		node.prev = tail
		l.head.prev = node
		tail.next = node
	}
	l.head = node
	l.size++
}

// Pop removes and returns the last element.
//
// complexity: O(1)
//
// Panics if the list is empty.
func (l *CircularDoublyLinkedList[T]) Pop() T {
	if data, ok := l.TryPop(); !ok {
		panic("CircularDoublyLinkedList.Pop: list is empty")
	} else {
		return data
	}
}

// TryPop attempts to remove and return the last element.
// Returns (value, true) on success, or (zero, false) if empty.
//
// This is a non-panicking version of Pop.
//
// complexity: O(1)
//
// SCORE: 15
func (l *CircularDoublyLinkedList[T]) TryPop() (T, bool) {
	// hint: 1) if empty, return (zero, false)
	//       2) tail = head.prev; save tail.data
	//       3) if size == 1: call reset()
	//       4) else: newTail = tail.prev; newTail.next = head; head.prev = newTail
	//       5) decrement size, return (saved, true)
	if l.Empty() {
		return generics.ZeroValue[T](), false
	}
	tail := l.head.prev
	popData := tail.data
	if l.size == 1 {
		l.reset()
		return popData, true
	} else {
		newTail := tail.prev
		newTail.next = l.head
		l.head.prev = newTail
	}
	l.size--
	return popData, true
}

// Shift removes and returns the first element.
//
// complexity: O(1)
//
// Panics if the list is empty.
func (l *CircularDoublyLinkedList[T]) Shift() T {
	if data, ok := l.TryShift(); !ok {
		panic("CircularDoublyLinkedList.Shift: list is empty")
	} else {
		return data
	}
}

// TryShift attempts to remove and return the first element.
// Returns (value, true) on success, or (zero, false) if empty.
//
// This is a non-panicking version of Shift.
//
// complexity: O(1)
//
// SCORE: 10
func (l *CircularDoublyLinkedList[T]) TryShift() (T, bool) {
	// hint: 1) if empty, return (zero, false)
	//       2) save head.data
	//       3) if size == 1: call reset()
	//       4) else: tail = head.prev; newHead = head.next
	//                tail.next = newHead; newHead.prev = tail; head = newHead
	//       5) decrement size, return (saved, true)
	if l.Empty() {
		return generics.ZeroValue[T](), false
	}
	shiftData := l.head.data
	if l.size == 1 {
		l.reset()
		return shiftData, true
	} else {
		tail := l.head.prev
		newHead := l.head.next
		newHead.prev = tail
		l.head = newHead
	}
	l.size--
	return shiftData, true
}

// Rotate moves the head pointer n positions.
// Positive n rotates forward (head moves to next), negative rotates backward.
//
// complexity: O(min(|n|, size - |n|))
func (l *CircularDoublyLinkedList[T]) Rotate(n int) {
	if l.Empty() || l.size == 1 {
		return
	}
	n = n % l.size
	if n == 0 {
		return
	}
	if n > 0 {
		for range n {
			l.head = l.head.next
		}
	} else {
		for range -n {
			l.head = l.head.prev
		}
	}
}

// Get retrieves the element at the given index.
// Traverses from whichever end is closer.
//
// complexity: O(min(k, size-k))
func (l *CircularDoublyLinkedList[T]) Get(index int) T {
	if v, ok := l.TryGet(index); !ok {
		if l.Empty() {
			panic("CircularDoublyLinkedList.Get: list is empty")
		}
		panic("CircularDoublyLinkedList.Get: index out of range")
	} else {
		return v
	}
}

// TryGet attempts to retrieve the element at the given index.
// Returns (value, true) on success, or (zero, false) if index is out of bounds.
//
// SCORE: 10
func (l *CircularDoublyLinkedList[T]) TryGet(index int) (T, bool) {
	// hint: 1) check bounds (empty || index < 0 || index >= size)
	//       2) use nodeAt(index) helper to find the node
	//       3) return (node.data, true)
	if l.Empty() || index < 0 || index >= l.size {
		return generics.ZeroValue[T](), false
	}
	node := l.nodeAt(index)
	return node.data, true
}

// Set updates the element at the given index.
//
// complexity: O(min(k, size-k))
func (l *CircularDoublyLinkedList[T]) Set(index int, data T) {
	if !l.TrySet(index, data) {
		if l.Empty() {
			panic("CircularDoublyLinkedList.Set: list is empty")
		}
		panic("CircularDoublyLinkedList.Set: index out of range")
	}
}

// TrySet attempts to update the element at the given index.
// Returns true on success, false if index is out of bounds.
//
// SCORE: 5
func (l *CircularDoublyLinkedList[T]) TrySet(index int, data T) bool {
	// hint: 1) check bounds, return false if invalid
	//       2) use nodeAt(index) to find node
	//       3) node.data = data
	//       4) return true
	if l.Empty() || index < 0 || index >= l.size {
		return false
	}
	node := l.nodeAt(index)
	node.data = data
	return true
}

// Insert adds an element at the given index.
//
// complexity: O(min(k, size-k))
//
// SCORE: 5
func (l *CircularDoublyLinkedList[T]) Insert(index int, data T) {
	// hint: 1) if index == 0: Prepend(data); return
	//       2) if index == size: Append(data); return
	//       3) checkBounds(index)
	//       4) curr = nodeAt(index)
	//       5) create new node with next=curr, prev=curr.prev
	//       6) curr.prev.next = node; curr.prev = node
	//       7) increment size
	if index == 0 {
		l.Prepend(data)
		return
	}
	if index == l.size {
		l.Append(data)
		return
	}
	l.checkBounds(index)
	// chain now like : previous node <-> The occupant (cur)
	cur := l.nodeAt(index)
	// Create new node that will be inserted
	newNode := NewBinaryNode(data, cur, cur.prev)
	// Re-wire to this : previous node <-> newNode <-> the occupant (cur)
	cur.prev.next = newNode
	cur.prev = newNode
	l.size++
}

// Remove deletes and returns the element at the given index.
//
// complexity: O(min(k, size-k))
func (l *CircularDoublyLinkedList[T]) Remove(index int) T {
	if v, ok := l.TryRemove(index); !ok {
		panic("CircularDoublyLinkedList.Remove: index out of range")
	} else {
		return v
	}
}

// TryRemove attempts to remove the element at the given index.
// Returns (value, true) on success, or (zero, false) if index is out of bounds.
//
// SCORE: 10
func (l *CircularDoublyLinkedList[T]) TryRemove(index int) (T, bool) {
	// hint: 1) check bounds, return (zero, false) if invalid
	//       2) if index == 0: return TryShift()
	//       3) if index == size-1: return TryPop()
	//       4) curr = nodeAt(index)
	//       5) rewire: curr.prev.next = curr.next; curr.next.prev = curr.prev
	//       6) save curr.data, clear pointers, decrement size
	//       7) return (saved, true)
	if l.Empty() || index < 0 || index >= l.size {
		return generics.ZeroValue[T](), false
	}
	if index == 0 {
		return l.TryShift()
	}
	if index == l.size-1 {
		return l.TryPop()
	}
	cur := l.nodeAt(index)
	cur.prev.next = cur.next
	cur.next.prev = cur.prev
	removeData := cur.data
	cur.next = nil
	cur.prev = nil
	l.size--
	return removeData, true
}

// Iter iterates over all elements starting from head.
//
// SCORE: 10
func (l *CircularDoublyLinkedList[T]) Iter(yield func(T) bool) {
	// hint: 1) if empty, return
	//       2) p := head; loop exactly 'size' times
	//       3) yield(p.data); p = p.next
	if l.Empty() {
		return
	}
	p := l.head
	for range l.Size() {
		if !yield(p.data) {
			return
		}
		p = p.next
	}
}

// IterBackward iterates over all elements from tail to head.
//
// SCORE: 10
func (l *CircularDoublyLinkedList[T]) IterBackward(yield func(T) bool) {
	// hint: 1) if empty, return
	//       2) p := head.prev (tail); loop exactly 'size' times
	//       3) yield(p.data); p = p.prev
	if l.Empty() {
		return
	}
	p := l.head.prev
	for range l.Size() {
		if !yield(p.data) {
			return
		}
		p = p.prev
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
func (l *CircularDoublyLinkedList[T]) Enum(yield func(int, T) bool) {
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
//   - space: O(1)
func (l *CircularDoublyLinkedList[T]) EnumBackward(yield func(int, T) bool) {
	i := l.Size() - 1
	for v := range l.IterBackward {
		if !yield(i, v) {
			break
		}
		i--
	}
}

// String returns the string representation.
func (l *CircularDoublyLinkedList[T]) String() string {
	return sequence.String(l.Iter)
}

func (l *CircularDoublyLinkedList[T]) nodeAt(index int) *BinaryNode[T] {
	if index < l.size/2 {
		p := l.head
		for range index {
			p = p.next
		}
		return p
	}
	p := l.head.prev
	for i := l.size - 1; i > index; i-- {
		p = p.prev
	}
	return p
}

func (l *CircularDoublyLinkedList[T]) checkBounds(index int) {
	if index < 0 || index >= l.size {
		panic("CircularDoublyLinkedList: index out of range")
	}
}

func (l *CircularDoublyLinkedList[T]) reset() {
	if l.head != nil {
		l.head.next = nil
		l.head.prev = nil
	}
	l.head = nil
	l.size = 0
}

// CircularIterator returns an infinite iterator that cycles through elements.
func (l *CircularDoublyLinkedList[T]) CircularIterator(yield func(T) bool) {
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

// CircularIteratorBackward returns an infinite backward iterator.
func (l *CircularDoublyLinkedList[T]) CircularIteratorBackward(yield func(T) bool) {
	if l.Empty() {
		return
	}
	p := l.head.prev
	for {
		if !yield(p.data) {
			break
		}
		p = p.prev
	}
}

// Cycle advances the head forward and returns the old head's data.
// Wraps around to the beginning when reaching the end.
// Useful for round-robin iteration.
//
// Panics if the list is empty.
func (l *CircularDoublyLinkedList[T]) Cycle() T {
	if v, ok := l.TryCycle(); !ok {
		panic("CircularDoublyLinkedList.Cycle: list is empty")
	} else {
		return v
	}
}

// TryCycle attempts to return the current head and advance to the next element.
// Returns (value, true) on success, or (zero, false) if empty.
func (l *CircularDoublyLinkedList[T]) TryCycle() (T, bool) {
	if l.Empty() {
		return generics.ZeroValue[T](), false
	}
	data := l.head.data
	l.head = l.head.next
	return data, true
}

// ReverseCycle moves the head backward and returns the new head's data.
// Wraps around to the end when reaching the beginning.
// Useful for reverse round-robin iteration.
//
// Panics if the list is empty.
func (l *CircularDoublyLinkedList[T]) ReverseCycle() T {
	if v, ok := l.TryReverseCycle(); !ok {
		panic("CircularDoublyLinkedList.ReverseCycle: list is empty")
	} else {
		return v
	}
}

// TryReverseCycle attempts to move the head backward and return the new head's data.
// Returns (value, true) on success, or (zero, false) if empty.
func (l *CircularDoublyLinkedList[T]) TryReverseCycle() (T, bool) {
	if l.Empty() {
		return generics.ZeroValue[T](), false
	}
	l.head = l.head.prev
	return l.head.data, true
}
