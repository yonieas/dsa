package linkedlist

import (
	"github.com/josestg/dsa/internal/generics"
	"github.com/josestg/dsa/sequence"
)

type UnaryNode[T any] struct {
	data T
	next *UnaryNode[T]
}

func NewUnaryNode[T any](data T, next *UnaryNode[T]) *UnaryNode[T] {
	return &UnaryNode[T]{
		data: data,
		next: next,
	}
}

type SinglyLinkedList[T any] struct {
	head *UnaryNode[T]
	tail *UnaryNode[T]
	size int
}

func NewSinglyLinkedList[T any]() *SinglyLinkedList[T] {
	return &SinglyLinkedList[T]{}
}

func (l *SinglyLinkedList[T]) Empty() bool {
	return l.size == 0 && l.head == nil && l.tail == nil
}

func (l *SinglyLinkedList[T]) Size() int {
	return l.size
}

func (l *SinglyLinkedList[T]) Head() T {
	if l.Empty() {
		panic("SinglyLinkedList.Head: is empty list")
	}
	return l.head.data
}

func (l *SinglyLinkedList[T]) Tail() T {
	if l.Empty() {
		panic("SinglyLinkedList.Tail: is empty list")
	}
	return l.tail.data
}

func (l *SinglyLinkedList[T]) Append(data T) {
	n := NewUnaryNode(data, nil)
	if l.Empty() {
		l.head = n
		l.tail = n
	} else {
		l.tail.next = n
		l.tail = n
	}
	l.size++
}

func (l *SinglyLinkedList[T]) Prepend(data T) {
	if l.Empty() {
		l.Append(data)
	} else {
		n := NewUnaryNode(data, l.head)
		l.head = n
		l.size++
	}
}

func (l *SinglyLinkedList[T]) Pop() T {
	if data, ok := l.TryPop(); !ok {
		panic("SinglyLinkedList.Pop: is empty list")
	} else {
		return data
	}
}

func (l *SinglyLinkedList[T]) TryPop() (T, bool) {
	if l.Empty() {
		return generics.ZeroValue[T](), false
	}
	data := l.tail.data
	if l.Size() == 1 {
		l.reset()
	} else {
		p := l.head
		for p.next != l.tail {
			p = p.next
		}
		p.next = nil
		l.tail = p
		l.size--
	}
	return data, true
}

func (l *SinglyLinkedList[T]) Shift() T {
	if data, ok := l.TryShift(); !ok {
		panic("SinglyLinkedList.Shift: is empty list")
	} else {
		return data
	}
}

func (l *SinglyLinkedList[T]) TryShift() (T, bool) {
	if l.Empty() {
		return generics.ZeroValue[T](), false
	}
	head := l.head
	if l.Size() == 1 {
		l.reset()
	} else {
		l.head = l.head.next
		head.next = nil
		l.size--
	}
	return head.data, true
}

func (l *SinglyLinkedList[T]) Iter(yield func(T) bool) {
	l.iterForward(func(u *UnaryNode[T]) bool { return yield(u.data) })
}

func (l *SinglyLinkedList[T]) IterBackward(yield func(T) bool) {
	l.iterBackward(func(u *UnaryNode[T]) bool { return yield(u.data) })
}

func (l *SinglyLinkedList[T]) iterForward(yield func(*UnaryNode[T]) bool) {
	p := l.head
	for p != nil {
		if !yield(p) {
			break
		}
		p = p.next
	}
}

func (l *SinglyLinkedList[T]) iterBackward(yield func(*UnaryNode[T]) bool) {
	l2 := NewSinglyLinkedList[T]()
	for v := range l.iterForward {
		l2.Prepend(v.data)
	}
	l2.iterForward(yield)
}

func (l *SinglyLinkedList[T]) Get(index int) T {
	if l.Empty() {
		panic("SinglyLinkedList.Get: is empty list")
	}
	l.checkBounds(index)
	data, ok := sequence.ValueAt(l.Iter, index)
	if !ok {
		panic("SinglyLinkedList.Get: should be unreachable")
	}
	return data
}

func (l *SinglyLinkedList[T]) Set(index int, data T) {
	if l.Empty() {
		panic("SinglyLinkedList.Set: is empty list")
	}
	l.checkBounds(index)
	n, ok := sequence.ValueAt(l.iterForward, index)
	if !ok {
		panic("SinglyLinkedList.Set: should be unreachable")
	}
	n.data = data
}
func (l *SinglyLinkedList[T]) String() string {
	return sequence.String(l.Iter)
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
