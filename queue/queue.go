package queue

import (
	"github.com/josestg/dsa/adt"
	"github.com/josestg/dsa/linkedlist"
)

type Backend[E any] interface {
	adt.Sizer
	adt.Emptier
	adt.Header[E]
	adt.Tailer[E]
	adt.Shifter[E]
	adt.Appender[E]
	adt.Stringer
}

type Queue[E any] struct {
	b Backend[E]
}

func New[E any]() *Queue[E] {
	return NewWith[E](linkedlist.NewDoublyLinkedList[E]())
}

func NewWith[E any](b Backend[E]) *Queue[E] {
	return &Queue[E]{b: b}
}

func (q *Queue[E]) Empty() bool {
	q.ensureBackend()
	return q.b.Empty()
}

func (q *Queue[E]) Size() int {
	q.ensureBackend()
	return q.b.Size()
}

func (q *Queue[E]) Peek() E {
	q.ensureBackend()
	return q.b.Head()
}

func (q *Queue[E]) Enqueue(data E) {
	q.ensureBackend()
	q.b.Append(data)
}

func (q *Queue[E]) Dequeue() E {
	q.ensureBackend()
	return q.b.Shift()
}

func (q *Queue[E]) String() string {
	q.ensureBackend()
	return q.b.String()
}

func (q *Queue[E]) ensureBackend() {
	if q.b == nil {
		q.b = linkedlist.NewSinglyLinkedList[E]()
	}
}
