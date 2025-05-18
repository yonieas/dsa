package stack

import (
	"github.com/josestg/dsa/adt"
	"github.com/josestg/dsa/linkedlist"
)

type Backend[E any] interface {
	adt.Sizer
	adt.Emptier
	adt.Tailer[E]
	adt.Popper[E]
	adt.Appender[E]
	adt.Stringer
}

type Stack[E any] struct {
	b Backend[E]
}

func New[E any]() *Stack[E] {
	return NewWith[E](linkedlist.NewDoublyLinkedList[E]())
}

func NewWith[E any](b Backend[E]) *Stack[E] {
	return &Stack[E]{b: b}
}

func (s *Stack[E]) Empty() bool {
	s.ensureBackend()
	return s.b.Empty()
}

func (s *Stack[E]) Size() int {
	s.ensureBackend()
	return s.b.Size()
}

func (s *Stack[E]) Peek() E {
	s.ensureBackend()
	if s.Empty() {
		panic("stack.Peek: stack is empty")
	}
	return s.b.Tail()
}

func (s *Stack[E]) Push(data E) {
	s.ensureBackend()
	s.b.Append(data)
}

func (s *Stack[E]) Pop() E {
	s.ensureBackend()
	if s.Empty() {
		panic("stack.Pop: stack underflow")
	}
	return s.b.Pop()
}

func (s *Stack[E]) String() string {
	s.ensureBackend()
	return s.b.String()
}

func (s *Stack[E]) ensureBackend() {
	if s.b == nil {
		s.b = linkedlist.NewDoublyLinkedList[E]()
	}
}
