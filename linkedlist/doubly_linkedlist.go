package linkedlist

import (
	"github.com/josestg/dsa/internal/generics"
	"github.com/josestg/dsa/sequence"
)

type BinaryNode[E any] struct {
	data E
	next *BinaryNode[E]
	prev *BinaryNode[E]
}

func NewBinaryNode[E any](data E, next *BinaryNode[E], prev *BinaryNode[E]) *BinaryNode[E] {
	return &BinaryNode[E]{
		data: data,
		next: next,
		prev: prev,
	}
}

type DoublyLinkedList[E any] struct {
	head *BinaryNode[E]
	tail *BinaryNode[E]
	size int
}

func NewDoublyLinkedList[E any]() *DoublyLinkedList[E] {
	return &DoublyLinkedList[E]{}
}

func (l *DoublyLinkedList[E]) Size() int {
	return l.size
}

func (l *DoublyLinkedList[E]) Empty() bool {
	return l.size == 0 && l.head == nil && l.tail == nil
}

func (l *DoublyLinkedList[E]) Tail() E {
	if l.Empty() {
		panic("DoublyLinkedList.Tail: list is empty")
	}
	return l.tail.data
}

func (l *DoublyLinkedList[E]) Head() E {
	if l.Empty() {
		panic("DoublyLinkedList.Head: list is empty")
	}
	return l.head.data
}

func (l *DoublyLinkedList[E]) Prepend(data E) {
	n := NewBinaryNode(data, l.head, nil)
	if l.Empty() {
		l.head = n
		l.tail = n
	} else {
		l.head = n
	}
	l.size++
}

func (l *DoublyLinkedList[E]) Append(data E) {
	n := NewBinaryNode(data, nil, nil)
	if l.Empty() {
		l.head = n
		l.tail = n
	} else {
		n.prev = l.tail
		l.tail.next = n
		l.tail = n
	}
	l.size++
}

func (l *DoublyLinkedList[E]) Pop() E {
	if data, ok := l.TryPop(); !ok {
		panic("DoublyLinkedList.Pop: list is empty")
	} else {
		return data
	}
}

func (l *DoublyLinkedList[E]) TryPop() (E, bool) {
	if l.Empty() {
		return generics.ZeroValue[E](), false
	}
	data := l.tail.data
	if l.Size() == 1 {
		l.reset()
	} else {
		prev := l.tail.prev
		prev.next = nil
		l.tail.prev = nil
		l.tail = prev
		l.size--
	}
	return data, true
}

func (l *DoublyLinkedList[E]) Shift() E {
	if data, ok := l.TryShift(); !ok {
		panic("DoublyLinkedList.Shift: list if empty")
	} else {
		return data
	}
}

func (l *DoublyLinkedList[E]) TryShift() (E, bool) {
	if l.Empty() {
		return generics.ZeroValue[E](), false
	}
	head := l.head
	if l.Size() == 1 {
		l.reset()
	} else {
		l.head = l.head.next
		l.head.prev = nil
		head.next = nil
		l.size--
	}
	return head.data, true
}

func (l *DoublyLinkedList[E]) Iter(yield func(E) bool) {
	l.iterForward(func(u *BinaryNode[E]) bool { return yield(u.data) })
}

func (l *DoublyLinkedList[E]) IterBackward(yield func(E) bool) {
	l.iterBackward(func(u *BinaryNode[E]) bool { return yield(u.data) })
}

func (l *DoublyLinkedList[E]) iterForward(yield func(node *BinaryNode[E]) bool) {
	p := l.head
	for p != nil {
		if !yield(p) {
			break
		}
		p = p.next
	}
}

func (l *DoublyLinkedList[E]) iterBackward(yield func(*BinaryNode[E]) bool) {
	p := l.tail
	for p != nil {
		if !yield(p) {
			break
		}
		p = p.prev
	}
}

func (l *DoublyLinkedList[E]) Get(index int) E {
	if l.Empty() {
		panic("DoublyLinkedList.Get: list is empty")
	}
	l.checkBounds(index)
	data, ok := sequence.ValueAt(l.Iter, index)
	if !ok {
		panic("DoublyLinkedList.Get: should be unreachable")
	}
	return data
}

func (l *DoublyLinkedList[E]) Set(index int, data E) {
	if l.Empty() {
		panic("DoublyLinkedList.Set: list is empty")
	}
	l.checkBounds(index)
	n, ok := sequence.ValueAt(l.iterForward, index)
	if !ok {
		panic("DoublyLinkedList.Get: should be unreachable")
	}
	n.data = data
}

func (l *DoublyLinkedList[E]) checkBounds(index int) {
	if index < 0 || index >= l.Size() {
		panic("DoublyLinkedList.checkBounds: index out of range")
	}
}

func (l *DoublyLinkedList[E]) String() string {
	return sequence.String(l.Iter)
}

func (l *DoublyLinkedList[E]) reset() {
	l.head = nil
	l.tail = nil
	l.size = 0
}
