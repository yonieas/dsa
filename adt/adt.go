package adt

import "fmt"

type Sizer interface {
	Size() int
}

type Caper interface {
	Cap() int
}

type Emptier interface {
	Empty() bool
}

type Getter[T any] interface {
	Get(int) T
}

type Setter[T any] interface {
	Set(int, T)
}

type Appender[T any] interface {
	Append(T)
}

type Prepender[T any] interface {
	Prepend(T)
}

type Tailer[T any] interface {
	Tail() T
}

type Header[T any] interface {
	Head() T
}

type Popper[T any] interface {
	Pop() T
}

type Shifter[T any] interface {
	Shift() T
}

type Iterator[T any] interface {
	Iter(func(T) bool)
}

type BackwordIterator[T any] interface {
	IterBackward(func(T) bool)
}

type Peeker[T any] interface {
	Peek() T
}

type Stringer = fmt.Stringer

type Stack[T any] interface {
	Sizer
	Emptier
	Peeker[T]
	Pop() T
	Push(T)
}

type Queue[T any] interface {
	Sizer
	Emptier
	Peeker[T]
	Enqueue(T)
	Dequeue() T
}
