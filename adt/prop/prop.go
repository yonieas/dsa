package prop

import (
	"iter"
	"slices"
	"testing"

	"github.com/josestg/dsa/adt"
)

const numSample = 5

func odds(n int) iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := range n {
			if !yield(2*i + 1) {
				break
			}
		}
	}
}

type Spec struct {
	Name string
	Test func(t *testing.T)
}

func Append[Abstract interface {
	adt.Sizer
	adt.Emptier
	adt.Tailer[int]
	adt.Appender[int]
}](f func() Abstract) Spec {
	return Spec{
		Name: "Append",
		Test: func(t *testing.T) {
			s := f()
			Empty(t, s)
			mustPanic(t, func() { _ = s.Tail() })

			s.Append(42)
			eq(t, s.Size(), 1)
			eq(t, s.Tail(), 42)

			s.Append(99)
			eq(t, s.Size(), 2)
			eq(t, s.Tail(), 99)

			for x := range odds(numSample) {
				s.Append(x)
				eq(t, s.Tail(), x)
			}
			eq(t, s.Size(), 2+numSample)
		},
	}
}

func Prepend[Abstract interface {
	adt.Sizer
	adt.Emptier
	adt.Header[int]
	adt.Prepender[int]
}](f func() Abstract) Spec {
	return Spec{
		Name: "Prepend",
		Test: func(t *testing.T) {
			s := f()
			Empty(t, s)
			mustPanic(t, func() { _ = s.Head() })

			s.Prepend(42)
			eq(t, s.Size(), 1)
			eq(t, s.Head(), 42)

			s.Prepend(99)
			eq(t, s.Size(), 2)
			eq(t, s.Head(), 99)

			for x := range odds(numSample) {
				s.Prepend(x)
				eq(t, s.Head(), x)
			}
			eq(t, s.Size(), 2+numSample)
		},
	}
}

func GetSet[Abstract interface {
	adt.Sizer
	adt.Getter[int]
	adt.Setter[int]
	adt.Appender[int]
}](f func() Abstract) Spec {
	return Spec{
		Name: "GetSet",
		Test: func(t *testing.T) {
			s := f()

			mustPanic(t, func() { s.Get(0) })
			mustPanic(t, func() { s.Get(-1) })
			mustPanic(t, func() { s.Get(100) })
			mustPanic(t, func() { s.Set(0, 1) })
			mustPanic(t, func() { s.Set(-1, 1) })

			s.Append(10)
			eq(t, s.Get(0), 10)
			s.Set(0, 20)
			eq(t, s.Get(0), 20)
			mustPanic(t, func() { s.Get(1) })
			mustPanic(t, func() { s.Set(1, 0) })

			s.Append(30)
			s.Append(40)
			eq(t, s.Get(0), 20)
			eq(t, s.Get(1), 30)
			eq(t, s.Get(2), 40)

			s.Set(1, 999)
			eq(t, s.Get(0), 20)
			eq(t, s.Get(1), 999)
			eq(t, s.Get(2), 40)

			mustPanic(t, func() { s.Get(s.Size()) })
			mustPanic(t, func() { s.Set(s.Size(), 0) })
			mustPanic(t, func() { s.Get(-1) })
			mustPanic(t, func() { s.Set(-1, 0) })
		},
	}
}

func HeadTail[Abstract interface {
	adt.Sizer
	adt.Emptier
	adt.Header[int]
	adt.Tailer[int]
	adt.Appender[int]
}](f func() Abstract) Spec {
	return Spec{
		Name: "HeadTail",
		Test: func(t *testing.T) {
			s := f()
			Empty(t, s)
			mustPanic(t, func() { _ = s.Head() })
			mustPanic(t, func() { _ = s.Tail() })

			s.Append(1)
			eq(t, s.Head(), 1)
			eq(t, s.Tail(), 1)

			s.Append(2)
			eq(t, s.Head(), 1)
			eq(t, s.Tail(), 2)

			s.Append(3)
			eq(t, s.Head(), 1)
			eq(t, s.Tail(), 3)

			for range 10 {
				_ = s.Head()
				_ = s.Tail()
			}
			eq(t, s.Size(), 3)
		},
	}
}

func Pop[Abstract interface {
	adt.Sizer
	adt.Emptier
	adt.Tailer[int]
	adt.Popper[int]
	adt.Appender[int]
}](f func() Abstract) Spec {
	return Spec{
		Name: "Pop",
		Test: func(t *testing.T) {
			s := f()
			mustPanic(t, func() { s.Pop() })

			s.Append(1)
			eq(t, s.Tail(), 1)
			eq(t, s.Pop(), 1)
			Empty(t, s)
			mustPanic(t, func() { s.Pop() })

			s.Append(10)
			s.Append(20)
			s.Append(30)
			eq(t, s.Pop(), 30)
			eq(t, s.Size(), 2)
			eq(t, s.Pop(), 20)
			eq(t, s.Size(), 1)
			eq(t, s.Tail(), 10)
			eq(t, s.Pop(), 10)
			Empty(t, s)

			var want []int
			for x := range odds(numSample) {
				s.Append(x)
				want = append(want, x)
			}
			for i := len(want) - 1; i >= 0; i-- {
				eq(t, s.Tail(), want[i])
				eq(t, s.Pop(), want[i])
			}
			Empty(t, s)
		},
	}
}

func Shift[Abstract interface {
	adt.Sizer
	adt.Emptier
	adt.Header[int]
	adt.Shifter[int]
	adt.Appender[int]
}](f func() Abstract) Spec {
	return Spec{
		Name: "Shift",
		Test: func(t *testing.T) {
			s := f()
			mustPanic(t, func() { s.Shift() })

			s.Append(1)
			eq(t, s.Head(), 1)
			eq(t, s.Shift(), 1)
			Empty(t, s)
			mustPanic(t, func() { s.Shift() })

			s.Append(10)
			s.Append(20)
			s.Append(30)
			eq(t, s.Shift(), 10)
			eq(t, s.Size(), 2)
			eq(t, s.Shift(), 20)
			eq(t, s.Size(), 1)
			eq(t, s.Head(), 30)
			eq(t, s.Shift(), 30)
			Empty(t, s)

			var want []int
			for x := range odds(numSample) {
				s.Append(x)
				want = append(want, x)
			}
			for _, w := range want {
				eq(t, s.Head(), w)
				eq(t, s.Shift(), w)
			}
			Empty(t, s)
		},
	}
}

func TryPop[Abstract interface {
	adt.Sizer
	adt.Emptier
	adt.Appender[int]
	TryPop() (int, bool)
}](f func() Abstract) Spec {
	return Spec{
		Name: "TryPop",
		Test: func(t *testing.T) {
			s := f()

			v, found := s.TryPop()
			eq(t, v, 0)
			eq(t, found, false)

			s.Append(10)
			s.Append(20)
			s.Append(30)

			v, found = s.TryPop()
			eq(t, v, 30)
			eq(t, found, true)
			eq(t, s.Size(), 2)

			v, found = s.TryPop()
			eq(t, v, 20)
			eq(t, found, true)
			eq(t, s.Size(), 1)

			v, found = s.TryPop()
			eq(t, v, 10)
			eq(t, found, true)
			Empty(t, s)

			v, found = s.TryPop()
			eq(t, v, 0)
			eq(t, found, false)
		},
	}
}

func TryShift[Abstract interface {
	adt.Sizer
	adt.Emptier
	adt.Appender[int]
	TryShift() (int, bool)
}](f func() Abstract) Spec {
	return Spec{
		Name: "TryShift",
		Test: func(t *testing.T) {
			s := f()

			v, found := s.TryShift()
			eq(t, v, 0)
			eq(t, found, false)

			s.Append(10)
			s.Append(20)
			s.Append(30)

			v, found = s.TryShift()
			eq(t, v, 10)
			eq(t, found, true)
			eq(t, s.Size(), 2)

			v, found = s.TryShift()
			eq(t, v, 20)
			eq(t, found, true)
			eq(t, s.Size(), 1)

			v, found = s.TryShift()
			eq(t, v, 30)
			eq(t, found, true)
			Empty(t, s)

			v, found = s.TryShift()
			eq(t, v, 0)
			eq(t, found, false)
		},
	}
}

func Stack[Abstract interface {
	adt.Stack[int]
}](f func() Abstract) Spec {
	return Spec{
		Name: "Stack",
		Test: func(t *testing.T) {
			s := f()
			Empty(t, s)
			mustPanic(t, func() { s.Peek() })
			mustPanic(t, func() { s.Pop() })

			s.Push(1)
			eq(t, s.Peek(), 1)
			eq(t, s.Size(), 1)

			s.Push(2)
			eq(t, s.Peek(), 2)

			s.Push(3)
			eq(t, s.Peek(), 3)
			eq(t, s.Size(), 3)

			eq(t, s.Pop(), 3)
			eq(t, s.Peek(), 2)
			eq(t, s.Pop(), 2)
			eq(t, s.Peek(), 1)
			eq(t, s.Pop(), 1)
			Empty(t, s)

			for i := range 10 {
				s.Push(i)
			}
			for i := 9; i >= 0; i-- {
				eq(t, s.Pop(), i)
			}
			Empty(t, s)
		},
	}
}

func Queue[Abstract interface {
	adt.Queue[int]
}](f func() Abstract) Spec {
	return Spec{
		Name: "Queue",
		Test: func(t *testing.T) {
			s := f()
			Empty(t, s)
			mustPanic(t, func() { s.Peek() })
			mustPanic(t, func() { s.Dequeue() })

			s.Enqueue(1)
			eq(t, s.Peek(), 1)
			eq(t, s.Size(), 1)

			s.Enqueue(2)
			eq(t, s.Peek(), 1)

			s.Enqueue(3)
			eq(t, s.Peek(), 1)
			eq(t, s.Size(), 3)

			eq(t, s.Dequeue(), 1)
			eq(t, s.Peek(), 2)
			eq(t, s.Dequeue(), 2)
			eq(t, s.Peek(), 3)
			eq(t, s.Dequeue(), 3)
			Empty(t, s)

			for i := range 10 {
				s.Enqueue(i)
			}
			for i := range 10 {
				eq(t, s.Dequeue(), i)
			}
			Empty(t, s)

			s.Enqueue(100)
			eq(t, s.Dequeue(), 100)
			s.Enqueue(200)
			s.Enqueue(300)
			eq(t, s.Dequeue(), 200)
			s.Enqueue(400)
			eq(t, s.Dequeue(), 300)
			eq(t, s.Dequeue(), 400)
			Empty(t, s)
		},
	}
}

func Iter[Abstract interface {
	adt.Sizer
	adt.Iterator[int]
	adt.Appender[int]
}](f func() Abstract) Spec {
	return Spec{
		Name: "Iter",
		Test: func(t *testing.T) {
			s := f()

			got := slices.Collect(s.Iter)
			eq(t, len(got), 0)

			s.Append(42)
			got = slices.Collect(s.Iter)
			ok(t, slices.Equal(got, []int{42}))

			s.Append(43)
			s.Append(44)
			got = slices.Collect(s.Iter)
			ok(t, slices.Equal(got, []int{42, 43, 44}))

			count := 0
			for range s.Iter {
				count++
				if count == 1 {
					break
				}
			}
			eq(t, count, 1)

			count = 0
			for range s.Iter {
				count++
				if count == 2 {
					break
				}
			}
			eq(t, count, 2)

			eq(t, s.Size(), 3)
		},
	}
}

func IterBackward[Abstract interface {
	adt.Sizer
	adt.BackwardIterator[int]
	adt.Appender[int]
}](f func() Abstract) Spec {
	return Spec{
		Name: "IterBackward",
		Test: func(t *testing.T) {
			s := f()

			got := slices.Collect(s.IterBackward)
			eq(t, len(got), 0)

			s.Append(1)
			got = slices.Collect(s.IterBackward)
			ok(t, slices.Equal(got, []int{1}))

			s.Append(2)
			s.Append(3)
			got = slices.Collect(s.IterBackward)
			ok(t, slices.Equal(got, []int{3, 2, 1}))

			count := 0
			for range s.IterBackward {
				count++
				if count == 2 {
					break
				}
			}
			eq(t, count, 2)

			eq(t, s.Size(), 3)
		},
	}
}

func Enum[Abstract interface {
	adt.Sizer
	adt.Enumerator[int]
	adt.Appender[int]
}](f func() Abstract) Spec {
	return Spec{
		Name: "Enum",
		Test: func(t *testing.T) {
			s := f()

			count := 0
			for range s.Enum {
				count++
			}
			eq(t, count, 0)

			s.Append(10)
			for i, v := range s.Enum {
				eq(t, i, 0)
				eq(t, v, 10)
			}

			s.Append(20)
			s.Append(30)
			want := []int{10, 20, 30}
			for i, v := range s.Enum {
				eq(t, v, want[i])
			}

			count = 0
			for i := range s.Enum {
				eq(t, i, count)
				count++
				if count == 2 {
					break
				}
			}
			eq(t, count, 2)
		},
	}
}

func EnumBackward[Abstract interface {
	adt.Sizer
	adt.BackwardEnumerator[int]
	adt.Appender[int]
}](f func() Abstract) Spec {
	return Spec{
		Name: "EnumBackward",
		Test: func(t *testing.T) {
			s := f()

			count := 0
			for range s.EnumBackward {
				count++
			}
			eq(t, count, 0)

			s.Append(10)
			for i, v := range s.EnumBackward {
				eq(t, i, 0)
				eq(t, v, 10)
			}

			s.Append(20)
			s.Append(30)
			wantIdx := []int{2, 1, 0}
			wantVal := []int{30, 20, 10}
			pos := 0
			for i, v := range s.EnumBackward {
				eq(t, i, wantIdx[pos])
				eq(t, v, wantVal[pos])
				pos++
			}

			count = 0
			for range s.EnumBackward {
				count++
				if count == 2 {
					break
				}
			}
			eq(t, count, 2)
		},
	}
}

func Insert[Abstract interface {
	adt.Sizer
	adt.Getter[int]
	adt.Inserter[int]
}](f func() Abstract) Spec {
	return Spec{
		Name: "Insert",
		Test: func(t *testing.T) {
			s := f()
			mustPanic(t, func() { s.Insert(-1, 0) })
			mustPanic(t, func() { s.Insert(1, 0) })

			s.Insert(0, 10)
			eq(t, s.Get(0), 10)
			eq(t, s.Size(), 1)

			s.Insert(0, 5)
			eq(t, s.Get(0), 5)
			eq(t, s.Get(1), 10)
			eq(t, s.Size(), 2)

			s.Insert(2, 20)
			eq(t, s.Get(0), 5)
			eq(t, s.Get(1), 10)
			eq(t, s.Get(2), 20)
			eq(t, s.Size(), 3)

			s.Insert(1, 7)
			eq(t, s.Get(0), 5)
			eq(t, s.Get(1), 7)
			eq(t, s.Get(2), 10)
			eq(t, s.Get(3), 20)
			eq(t, s.Size(), 4)

			mustPanic(t, func() { s.Insert(-1, 0) })
			mustPanic(t, func() { s.Insert(s.Size()+1, 0) })
		},
	}
}

func Remove[Abstract interface {
	adt.Sizer
	adt.Emptier
	adt.Getter[int]
	adt.Remover[int]
	adt.Appender[int]
}](f func() Abstract) Spec {
	return Spec{
		Name: "Remove",
		Test: func(t *testing.T) {
			s := f()
			mustPanic(t, func() { s.Remove(0) })
			mustPanic(t, func() { s.Remove(-1) })

			s.Append(10)
			eq(t, s.Remove(0), 10)
			Empty(t, s)
			mustPanic(t, func() { s.Remove(0) })

			s.Append(1)
			s.Append(2)
			s.Append(3)
			s.Append(4)
			s.Append(5)

			eq(t, s.Remove(2), 3)
			eq(t, s.Size(), 4)
			eq(t, s.Get(0), 1)
			eq(t, s.Get(1), 2)
			eq(t, s.Get(2), 4)
			eq(t, s.Get(3), 5)

			eq(t, s.Remove(0), 1)
			eq(t, s.Get(0), 2)

			eq(t, s.Remove(s.Size()-1), 5)
			eq(t, s.Get(s.Size()-1), 4)

			mustPanic(t, func() { s.Remove(s.Size()) })
			mustPanic(t, func() { s.Remove(-1) })
		},
	}
}

func Swap[Abstract interface {
	adt.Sizer
	adt.Getter[int]
	adt.Appender[int]
	Swap(i, j int)
}](f func() Abstract) Spec {
	return Spec{
		Name: "Swap",
		Test: func(t *testing.T) {
			s := f()

			s.Append(1)
			s.Swap(0, 0)
			eq(t, s.Get(0), 1)

			s.Append(2)
			s.Swap(0, 1)
			eq(t, s.Get(0), 2)
			eq(t, s.Get(1), 1)

			s.Swap(1, 0)
			eq(t, s.Get(0), 1)
			eq(t, s.Get(1), 2)

			s.Append(3)
			s.Append(4)
			s.Append(5)

			s.Swap(0, 4)
			eq(t, s.Get(0), 5)
			eq(t, s.Get(4), 1)

			s.Swap(1, 3)
			eq(t, s.Get(1), 4)
			eq(t, s.Get(3), 2)

			s.Swap(2, 2)
			eq(t, s.Get(2), 3)
		},
	}
}

func AddExistsDel[Abstract interface {
	adt.Sizer
	adt.Emptier
	adt.Adder[int]
	adt.Exister[int]
	adt.Deleter[int]
}](f func() Abstract) Spec {
	return Spec{
		Name: "AddExistsDel",
		Test: func(t *testing.T) {
			s := f()
			Empty(t, s)

			ok(t, !s.Exists(1))
			ok(t, !s.Exists(0))
			ok(t, !s.Exists(-1))

			s.Add(1)
			ok(t, s.Exists(1))
			ok(t, !s.Exists(2))
			eq(t, s.Size(), 1)

			s.Add(1)
			ok(t, s.Exists(1))
			ok(t, s.Size() >= 1)

			s.Add(2)
			s.Add(3)
			s.Add(4)
			s.Add(5)
			ok(t, s.Exists(1))
			ok(t, s.Exists(2))
			ok(t, s.Exists(3))
			ok(t, s.Exists(4))
			ok(t, s.Exists(5))
			ok(t, !s.Exists(99))
			ok(t, !s.Exists(0))

			s.Del(3)
			ok(t, !s.Exists(3))
			ok(t, s.Exists(1))
			ok(t, s.Exists(2))
			ok(t, s.Exists(4))
			ok(t, s.Exists(5))

			s.Del(99)

			s.Del(1)
			ok(t, !s.Exists(1))

			s.Del(5)
			ok(t, !s.Exists(5))

			s.Del(2)
			s.Del(4)
			Empty(t, s)

			s.Del(1)
			Empty(t, s)
		},
	}
}

func Cap[Abstract interface {
	adt.Sizer
	adt.Caper
	adt.Appender[int]
}](f func(capacity int) Abstract) Spec {
	return Spec{
		Name: "Cap",
		Test: func(t *testing.T) {
			s := f(10)
			ok(t, s.Cap() >= 10)
			eq(t, s.Size(), 0)

			s.Append(1)
			ok(t, s.Cap() >= s.Size())

			for i := range 20 {
				s.Append(i)
				ok(t, s.Cap() >= s.Size())
			}
		},
	}
}

func Clip[Abstract interface {
	adt.Sizer
	adt.Caper
	adt.Appender[int]
	Clip()
}](f func(capacity int) Abstract) Spec {
	return Spec{
		Name: "Clip",
		Test: func(t *testing.T) {
			s := f(100)
			ok(t, s.Cap() >= 100)

			for i := range 5 {
				s.Append(i)
			}
			eq(t, s.Size(), 5)
			ok(t, s.Cap() >= 100)

			s.Clip()
			eq(t, s.Size(), 5)
			eq(t, s.Cap(), 5)
		},
	}
}

func BitAddExistsDel[Abstract interface {
	Add(pos int)
	Del(pos int)
	Exists(pos int) bool
	Len() int
	Count() int
}](f func(numOfBits int) Abstract) Spec {
	return Spec{
		Name: "BitAddExistsDel",
		Test: func(t *testing.T) {
			b := f(128)

			for i := range b.Len() {
				if b.Exists(i) {
					t.Errorf("Exists(%d) = true, want false", i)
				}
			}

			b.Add(0)
			ok(t, b.Exists(0))
			ok(t, !b.Exists(1))
			eq(t, b.Count(), 1)

			b.Add(0)
			ok(t, b.Exists(0))
			eq(t, b.Count(), 1)

			b.Add(1)
			b.Add(63)
			b.Add(64)
			b.Add(127)

			ok(t, b.Exists(0))
			ok(t, b.Exists(1))
			ok(t, b.Exists(63))
			ok(t, b.Exists(64))
			ok(t, b.Exists(127))
			ok(t, !b.Exists(2))
			ok(t, !b.Exists(62))
			ok(t, !b.Exists(65))
			eq(t, b.Count(), 5)

			b.Del(63)
			ok(t, !b.Exists(63))
			ok(t, b.Exists(64))
			eq(t, b.Count(), 4)

			b.Del(63)
			ok(t, !b.Exists(63))
			eq(t, b.Count(), 4)

			b.Del(100)
			eq(t, b.Count(), 4)

			b.Del(0)
			b.Del(1)
			b.Del(64)
			b.Del(127)
			eq(t, b.Count(), 0)
		},
	}
}

func BitToggle[Abstract interface {
	Add(pos int)
	Toggle(pos int)
	Exists(pos int) bool
	Count() int
}](f func(numOfBits int) Abstract) Spec {
	return Spec{
		Name: "BitToggle",
		Test: func(t *testing.T) {
			b := f(128)

			ok(t, !b.Exists(10))
			b.Toggle(10)
			ok(t, b.Exists(10))
			eq(t, b.Count(), 1)

			b.Toggle(10)
			ok(t, !b.Exists(10))
			eq(t, b.Count(), 0)

			b.Add(50)
			ok(t, b.Exists(50))
			b.Toggle(50)
			b.Toggle(50)
			ok(t, b.Exists(50))

			b.Toggle(0)
			b.Toggle(63)
			b.Toggle(64)
			b.Toggle(127)
			ok(t, b.Exists(0))
			ok(t, b.Exists(63))
			ok(t, b.Exists(64))
			ok(t, b.Exists(127))
		},
	}
}

func BitReset[Abstract interface {
	Add(pos int)
	Exists(pos int) bool
	Reset()
	Len() int
	Count() int
}](f func(numOfBits int) Abstract) Spec {
	return Spec{
		Name: "BitReset",
		Test: func(t *testing.T) {
			b := f(256)

			b.Reset()
			eq(t, b.Count(), 0)

			for i := 0; i < b.Len(); i += 7 {
				b.Add(i)
			}
			ok(t, b.Count() > 0)

			b.Reset()
			eq(t, b.Count(), 0)

			for i := range b.Len() {
				if b.Exists(i) {
					t.Errorf("Exists(%d) = true after Reset", i)
				}
			}

			b.Add(100)
			ok(t, b.Exists(100))
			eq(t, b.Count(), 1)
		},
	}
}

func BitCount[Abstract interface {
	Add(pos int)
	Del(pos int)
	Toggle(pos int)
	Reset()
	Count() int
}](f func(numOfBits int) Abstract) Spec {
	return Spec{
		Name: "BitCount",
		Test: func(t *testing.T) {
			b := f(128)

			eq(t, b.Count(), 0)

			b.Add(0)
			eq(t, b.Count(), 1)

			b.Add(1)
			eq(t, b.Count(), 2)

			b.Add(1)
			eq(t, b.Count(), 2)

			for i := 10; i < 20; i++ {
				b.Add(i)
			}
			eq(t, b.Count(), 12)

			b.Del(15)
			eq(t, b.Count(), 11)

			b.Toggle(100)
			eq(t, b.Count(), 12)
			b.Toggle(100)
			eq(t, b.Count(), 11)

			b.Reset()
			eq(t, b.Count(), 0)
		},
	}
}

func BitLen[Abstract interface {
	Add(pos int)
	Reset()
	Len() int
	Count() int
}](f func(numOfBits int) Abstract) Spec {
	return Spec{
		Name: "BitLen",
		Test: func(t *testing.T) {
			b64 := f(64)
			eq(t, b64.Len(), 64)

			b128 := f(128)
			eq(t, b128.Len(), 128)

			b256 := f(256)
			eq(t, b256.Len(), 256)

			b128.Add(0)
			b128.Add(127)
			eq(t, b128.Len(), 128)
			eq(t, b128.Count(), 2)

			b128.Reset()
			eq(t, b128.Len(), 128)
		},
	}
}

func BitBounds[Abstract interface {
	Add(pos int)
	Del(pos int)
	Toggle(pos int)
	Exists(pos int) bool
}](f func(numOfBits int) Abstract) Spec {
	return Spec{
		Name: "BitBounds",
		Test: func(t *testing.T) {
			b := f(64)

			mustPanic(t, func() { b.Add(-1) })
			mustPanic(t, func() { b.Del(-1) })
			mustPanic(t, func() { b.Toggle(-1) })
			mustPanic(t, func() { b.Exists(-1) })

			mustPanic(t, func() { b.Add(64) })
			mustPanic(t, func() { b.Del(64) })
			mustPanic(t, func() { b.Toggle(64) })
			mustPanic(t, func() { b.Exists(64) })

			mustPanic(t, func() { b.Add(100) })
			mustPanic(t, func() { b.Exists(1000) })

			b.Add(63)
			ok(t, b.Exists(63))
			b.Del(63)
			ok(t, !b.Exists(63))
		},
	}
}

func BitString[Abstract interface {
	Add(pos int)
	Reset()
	Len() int
	String() string
}](f func(numOfBits int) Abstract) Spec {
	return Spec{
		Name: "BitString",
		Test: func(t *testing.T) {
			b := f(64)

			s := b.String()
			eq(t, len(s), b.Len())
			for _, c := range s {
				if c != '0' {
					t.Errorf("String() contains %c, want all '0'", c)
				}
			}

			b.Add(0)
			s = b.String()
			if s[0] != '1' {
				t.Errorf("String()[0] = %c after Add(0), want '1'", s[0])
			}

			b.Add(5)
			s = b.String()
			if s[5] != '1' {
				t.Errorf("String()[5] = %c after Add(5), want '1'", s[5])
			}

			b.Reset()
			b.Add(0)
			b.Add(2)
			b.Add(4)
			s = b.String()
			want := "10101"
			if s[:5] != want {
				t.Errorf("String()[:5] = %q, want %q", s[:5], want)
			}
		},
	}
}

func MapPutGetDel[Abstract interface {
	Put(int, int)
	Get(int) (int, bool)
	Del(int)
	Exists(int) bool
	Size() int
	Empty() bool
}](f func() Abstract) Spec {
	return Spec{
		Name: "MapPutGetDel",
		Test: func(t *testing.T) {
			m := f()
			ok(t, m.Empty())
			eq(t, m.Size(), 0)

			_, found := m.Get(1)
			ok(t, !found)
			ok(t, !m.Exists(1))

			m.Put(1, 100)
			v, found := m.Get(1)
			ok(t, found)
			eq(t, v, 100)
			ok(t, m.Exists(1))
			eq(t, m.Size(), 1)

			m.Put(1, 200)
			v, found = m.Get(1)
			ok(t, found)
			eq(t, v, 200)
			eq(t, m.Size(), 1)

			m.Put(2, 10)
			m.Put(3, 20)
			m.Put(4, 30)
			eq(t, m.Size(), 4)

			ok(t, m.Exists(1))
			ok(t, m.Exists(2))
			ok(t, m.Exists(3))
			ok(t, m.Exists(4))
			ok(t, !m.Exists(99))

			m.Del(2)
			ok(t, !m.Exists(2))
			eq(t, m.Size(), 3)

			m.Del(99)
			eq(t, m.Size(), 3)

			m.Del(1)
			m.Del(3)
			m.Del(4)
			ok(t, m.Empty())
		},
	}
}

func MapKeys[Abstract interface {
	Put(int, int)
	Keys(func(int) bool)
	Size() int
}](f func() Abstract) Spec {
	return Spec{
		Name: "MapKeys",
		Test: func(t *testing.T) {
			m := f()

			count := 0
			for range m.Keys {
				count++
			}
			eq(t, count, 0)

			m.Put(1, 10)
			m.Put(2, 20)
			m.Put(3, 30)

			keys := make(map[int]bool)
			for k := range m.Keys {
				keys[k] = true
			}
			eq(t, len(keys), 3)
			ok(t, keys[1])
			ok(t, keys[2])
			ok(t, keys[3])

			count = 0
			for range m.Keys {
				count++
				if count == 2 {
					break
				}
			}
			eq(t, count, 2)
		},
	}
}

func MapLoadFactor[Abstract interface {
	Put(int, int)
	Cap() int
	LoadFactor() float64
	Size() int
}](f func() Abstract) Spec {
	return Spec{
		Name: "MapLoadFactor",
		Test: func(t *testing.T) {
			m := f()
			eq(t, m.LoadFactor(), 0.0)

			cap := m.Cap()
			ok(t, cap > 0)

			m.Put(1, 10)
			lf := m.LoadFactor()
			ok(t, lf > 0)
			ok(t, lf == float64(m.Size())/float64(m.Cap()))

			for i := 0; i < 100; i++ {
				m.Put(i, i*10)
			}
			ok(t, m.LoadFactor() <= 1.0)
			ok(t, m.Cap() >= cap)
		},
	}
}

func MapString[Abstract interface {
	Put(int, int)
	String() string
	Size() int
}](f func() Abstract) Spec {
	return Spec{
		Name: "MapString",
		Test: func(t *testing.T) {
			m := f()
			s := m.String()
			eq(t, s, "[]")

			m.Put(1, 42)
			s = m.String()
			ok(t, len(s) > 2)
			ok(t, s[0] == '[')
			ok(t, s[len(s)-1] == ']')
		},
	}
}

func Rotate[Abstract interface {
	adt.Sizer
	adt.Header[int]
	adt.Tailer[int]
	adt.Appender[int]
	adt.Rotator
}](f func() Abstract) Spec {
	return Spec{
		Name: "Rotate",
		Test: func(t *testing.T) {
			c := f()

			c.Rotate(1)
			c.Rotate(-1)

			c.Append(1)
			c.Rotate(0)
			eq(t, c.Head(), 1)
			eq(t, c.Tail(), 1)

			c.Rotate(1)
			eq(t, c.Head(), 1)

			c.Append(2)
			c.Append(3)
			c.Append(4)

			eq(t, c.Head(), 1)
			eq(t, c.Tail(), 4)

			c.Rotate(1)
			eq(t, c.Head(), 2)
			eq(t, c.Tail(), 1)

			c.Rotate(2)
			eq(t, c.Head(), 4)
			eq(t, c.Tail(), 3)

			c.Rotate(-1)
			eq(t, c.Head(), 3)
			eq(t, c.Tail(), 2)

			c.Rotate(-2)
			eq(t, c.Head(), 1)
			eq(t, c.Tail(), 4)

			c.Rotate(4)
			eq(t, c.Head(), 1)
			eq(t, c.Tail(), 4)

			c.Rotate(-4)
			eq(t, c.Head(), 1)
			eq(t, c.Tail(), 4)

			c.Rotate(5)
			eq(t, c.Head(), 2)

			c.Rotate(-6)
			eq(t, c.Head(), 4)
		},
	}
}

func Cycle[Abstract interface {
	adt.Sizer
	adt.Emptier
	adt.Header[int]
	adt.Appender[int]
	adt.Cycler[int]
}](f func() Abstract) Spec {
	return Spec{
		Name: "Cycle",
		Test: func(t *testing.T) {
			c := f()

			mustPanic(t, func() { c.Cycle() })

			c.Append(1)
			c.Append(2)
			c.Append(3)

			eq(t, c.Cycle(), 1)
			eq(t, c.Head(), 2)

			eq(t, c.Cycle(), 2)
			eq(t, c.Head(), 3)

			eq(t, c.Cycle(), 3)
			eq(t, c.Head(), 1)

			eq(t, c.Cycle(), 1)
			eq(t, c.Head(), 2)

			eq(t, c.Size(), 3)
		},
	}
}

func ReverseCycle[Abstract interface {
	adt.Sizer
	adt.Emptier
	adt.Header[int]
	adt.Appender[int]
	adt.ReverseCycler[int]
}](f func() Abstract) Spec {
	return Spec{
		Name: "ReverseCycle",
		Test: func(t *testing.T) {
			c := f()

			mustPanic(t, func() { c.ReverseCycle() })

			c.Append(1)
			c.Append(2)
			c.Append(3)

			eq(t, c.ReverseCycle(), 3)
			eq(t, c.Head(), 3)

			eq(t, c.ReverseCycle(), 2)
			eq(t, c.Head(), 2)

			eq(t, c.ReverseCycle(), 1)
			eq(t, c.Head(), 1)

			eq(t, c.ReverseCycle(), 3)
			eq(t, c.Head(), 3)

			eq(t, c.Size(), 3)
		},
	}
}

func CircularIter[Abstract interface {
	adt.Appender[int]
	CircularIterator(func(int) bool)
}](f func() Abstract) Spec {
	return Spec{
		Name: "CircularIter",
		Test: func(t *testing.T) {
			c := f()

			count := 0
			for range c.CircularIterator {
				count++
				if count > 0 {
					break
				}
			}
			eq(t, count, 0)

			c.Append(1)
			c.Append(2)
			c.Append(3)

			var collected []int
			count = 0
			for v := range c.CircularIterator {
				collected = append(collected, v)
				count++
				if count >= 7 {
					break
				}
			}
			eq(t, len(collected), 7)
			eq(t, collected[0], 1)
			eq(t, collected[1], 2)
			eq(t, collected[2], 3)
			eq(t, collected[3], 1)
			eq(t, collected[4], 2)
			eq(t, collected[5], 3)
			eq(t, collected[6], 1)
		},
	}
}

func Size(t *testing.T, s interface {
	adt.Sizer
	adt.Emptier
}, n int) {
	t.Helper()
	ok(t, !s.Empty())
	eq(t, s.Size(), n)
}

func Empty(t *testing.T, s interface {
	adt.Sizer
	adt.Emptier
}) {
	t.Helper()
	ok(t, s.Empty())
	eq(t, s.Size(), 0)
}

func mustPanic(t *testing.T, fn func()) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Error("expected panic")
		}
	}()
	fn()
}

func Set[Abstract interface {
	adt.Sizer
	adt.Emptier
	adt.Adder[int]
	adt.Deleter[int]
	adt.Exister[int]
	adt.Iterator[int]
	adt.Stringer
}](f func() Abstract) Spec {
	return Spec{
		Name: "Set",
		Test: func(t *testing.T) {
			s := f()
			Empty(t, s)

			s.Add(1)
			eq(t, s.Size(), 1)
			ok(t, s.Exists(1))
			ok(t, !s.Exists(2))

			s.Add(1)
			eq(t, s.Size(), 1)

			s.Add(2)
			s.Add(3)
			eq(t, s.Size(), 3)
			ok(t, s.Exists(1))
			ok(t, s.Exists(2))
			ok(t, s.Exists(3))

			s.Del(2)
			eq(t, s.Size(), 2)
			ok(t, !s.Exists(2))
			ok(t, s.Exists(1))
			ok(t, s.Exists(3))

			var collected []int
			for v := range s.Iter {
				collected = append(collected, v)
			}
			eq(t, len(collected), 2)

			s.Del(1)
			s.Del(3)
			Empty(t, s)
		},
	}
}

func Union[Abstract interface {
	adt.Adder[int]
	adt.Exister[int]
	adt.Unioner[Abstract]
}](f func() Abstract) Spec {
	return Spec{
		Name: "Union",
		Test: func(t *testing.T) {
			a := f()
			a.Add(1)
			a.Add(2)
			a.Add(3)

			b := f()
			b.Add(3)
			b.Add(4)
			b.Add(5)

			union := a.Union(b)
			ok(t, union.Exists(1))
			ok(t, union.Exists(2))
			ok(t, union.Exists(3))
			ok(t, union.Exists(4))
			ok(t, union.Exists(5))
			ok(t, !union.Exists(6))
		},
	}
}

func Intersection[Abstract interface {
	adt.Adder[int]
	adt.Exister[int]
	adt.Intersecter[Abstract]
}](f func() Abstract) Spec {
	return Spec{
		Name: "Intersection",
		Test: func(t *testing.T) {
			a := f()
			a.Add(1)
			a.Add(2)
			a.Add(3)
			a.Add(4)

			b := f()
			b.Add(3)
			b.Add(4)
			b.Add(5)
			b.Add(6)

			intersection := a.Intersection(b)
			ok(t, !intersection.Exists(1))
			ok(t, !intersection.Exists(2))
			ok(t, intersection.Exists(3))
			ok(t, intersection.Exists(4))
			ok(t, !intersection.Exists(5))
			ok(t, !intersection.Exists(6))
		},
	}
}

func Disjoint[Abstract interface {
	adt.Adder[int]
	adt.Disjointer[Abstract]
}](f func() Abstract) Spec {
	return Spec{
		Name: "Disjoint",
		Test: func(t *testing.T) {
			a := f()
			a.Add(1)
			a.Add(2)
			a.Add(3)

			b := f()
			b.Add(4)
			b.Add(5)
			b.Add(6)

			ok(t, a.Disjoint(b))

			c := f()
			c.Add(3)
			c.Add(4)
			c.Add(5)

			ok(t, !a.Disjoint(c))
		},
	}
}

func BSTMinMax[Abstract interface {
	adt.Emptier
	adt.Adder[int]
	Min() (int, bool)
	Max() (int, bool)
}](f func() Abstract) Spec {
	return Spec{
		Name: "BSTMinMax",
		Test: func(t *testing.T) {
			tree := f()

			_, found := tree.Min()
			ok(t, !found)
			_, found = tree.Max()
			ok(t, !found)

			tree.Add(5)
			min, found := tree.Min()
			ok(t, found)
			eq(t, min, 5)
			max, found := tree.Max()
			ok(t, found)
			eq(t, max, 5)

			tree.Add(3)
			tree.Add(7)
			tree.Add(1)
			tree.Add(9)

			min, _ = tree.Min()
			eq(t, min, 1)
			max, _ = tree.Max()
			eq(t, max, 9)

			tree.Add(0)
			min, _ = tree.Min()
			eq(t, min, 0)

			tree.Add(10)
			max, _ = tree.Max()
			eq(t, max, 10)
		},
	}
}

func BSTInOrder[Abstract interface {
	adt.Adder[int]
	InOrder(func(int) bool)
}](f func() Abstract) Spec {
	return Spec{
		Name: "BSTInOrder",
		Test: func(t *testing.T) {
			tree := f()

			var collected []int
			tree.InOrder(func(v int) bool {
				collected = append(collected, v)
				return true
			})
			eq(t, len(collected), 0)

			tree.Add(5)
			tree.Add(3)
			tree.Add(7)
			tree.Add(1)
			tree.Add(4)
			tree.Add(6)
			tree.Add(8)

			collected = nil
			tree.InOrder(func(v int) bool {
				collected = append(collected, v)
				return true
			})
			ok(t, slices.Equal(collected, []int{1, 3, 4, 5, 6, 7, 8}))

			count := 0
			tree.InOrder(func(v int) bool {
				count++
				return count < 3
			})
			eq(t, count, 3)
		},
	}
}

func BSTPreOrder[Abstract interface {
	adt.Adder[int]
	PreOrder(func(int) bool)
}](f func() Abstract) Spec {
	return Spec{
		Name: "BSTPreOrder",
		Test: func(t *testing.T) {
			tree := f()

			tree.Add(5)
			tree.Add(3)
			tree.Add(7)
			tree.Add(1)
			tree.Add(4)

			var collected []int
			tree.PreOrder(func(v int) bool {
				collected = append(collected, v)
				return true
			})
			ok(t, slices.Equal(collected, []int{5, 3, 1, 4, 7}))
		},
	}
}

func BSTPostOrder[Abstract interface {
	adt.Adder[int]
	PostOrder(func(int) bool)
}](f func() Abstract) Spec {
	return Spec{
		Name: "BSTPostOrder",
		Test: func(t *testing.T) {
			tree := f()

			tree.Add(5)
			tree.Add(3)
			tree.Add(7)
			tree.Add(1)
			tree.Add(4)

			var collected []int
			tree.PostOrder(func(v int) bool {
				collected = append(collected, v)
				return true
			})
			ok(t, slices.Equal(collected, []int{1, 4, 3, 7, 5}))
		},
	}
}

func BSTString[Abstract interface {
	adt.Emptier
	adt.Adder[int]
	adt.Stringer
}](f func() Abstract) Spec {
	return Spec{
		Name: "BSTString",
		Test: func(t *testing.T) {
			tree := f()
			ok(t, tree.Empty())
			eq(t, tree.String(), "[]")

			tree.Add(5)
			eq(t, tree.String(), "[5]")

			tree.Add(3)
			tree.Add(7)
			tree.Add(1)
			tree.Add(4)
			// In-order traversal should give sorted order
			eq(t, tree.String(), "[1 3 4 5 7]")
		},
	}
}

func TryHead[Abstract interface {
	adt.Sizer
	adt.Emptier
	adt.Appender[int]
	TryHead() (int, bool)
}](f func() Abstract) Spec {
	return Spec{
		Name: "TryHead",
		Test: func(t *testing.T) {
			s := f()

			v, found := s.TryHead()
			eq(t, v, 0)
			eq(t, found, false)

			s.Append(42)
			v, found = s.TryHead()
			eq(t, v, 42)
			eq(t, found, true)

			s.Append(99)
			v, found = s.TryHead()
			eq(t, v, 42)
			eq(t, found, true)
		},
	}
}

func TryTail[Abstract interface {
	adt.Sizer
	adt.Emptier
	adt.Appender[int]
	TryTail() (int, bool)
}](f func() Abstract) Spec {
	return Spec{
		Name: "TryTail",
		Test: func(t *testing.T) {
			s := f()

			v, found := s.TryTail()
			eq(t, v, 0)
			eq(t, found, false)

			s.Append(42)
			v, found = s.TryTail()
			eq(t, v, 42)
			eq(t, found, true)

			s.Append(99)
			v, found = s.TryTail()
			eq(t, v, 99)
			eq(t, found, true)
		},
	}
}

func TryGet[Abstract interface {
	adt.Sizer
	adt.Emptier
	adt.Appender[int]
	TryGet(int) (int, bool)
}](f func() Abstract) Spec {
	return Spec{
		Name: "TryGet",
		Test: func(t *testing.T) {
			s := f()

			v, found := s.TryGet(0)
			eq(t, v, 0)
			eq(t, found, false)

			v, found = s.TryGet(-1)
			eq(t, found, false)

			s.Append(10)
			s.Append(20)
			s.Append(30)

			v, found = s.TryGet(0)
			eq(t, v, 10)
			eq(t, found, true)

			v, found = s.TryGet(1)
			eq(t, v, 20)
			eq(t, found, true)

			v, found = s.TryGet(2)
			eq(t, v, 30)
			eq(t, found, true)

			v, found = s.TryGet(3)
			eq(t, found, false)

			v, found = s.TryGet(-1)
			eq(t, found, false)
		},
	}
}

func TrySet[Abstract interface {
	adt.Sizer
	adt.Emptier
	adt.Getter[int]
	adt.Appender[int]
	TrySet(int, int) bool
}](f func() Abstract) Spec {
	return Spec{
		Name: "TrySet",
		Test: func(t *testing.T) {
			s := f()

			ok := s.TrySet(0, 42)
			eq(t, ok, false)

			ok = s.TrySet(-1, 42)
			eq(t, ok, false)

			s.Append(10)
			s.Append(20)
			s.Append(30)

			ok = s.TrySet(1, 99)
			eq(t, ok, true)
			eq(t, s.Get(1), 99)

			ok = s.TrySet(0, 88)
			eq(t, ok, true)
			eq(t, s.Get(0), 88)

			ok = s.TrySet(3, 100)
			eq(t, ok, false)

			ok = s.TrySet(-1, 100)
			eq(t, ok, false)
		},
	}
}

func TryRemove[Abstract interface {
	adt.Sizer
	adt.Emptier
	adt.Getter[int]
	adt.Appender[int]
	TryRemove(int) (int, bool)
}](f func() Abstract) Spec {
	return Spec{
		Name: "TryRemove",
		Test: func(t *testing.T) {
			s := f()

			v, found := s.TryRemove(0)
			eq(t, found, false)

			s.Append(10)
			s.Append(20)
			s.Append(30)

			v, found = s.TryRemove(1)
			eq(t, v, 20)
			eq(t, found, true)
			eq(t, s.Size(), 2)

			v, found = s.TryRemove(0)
			eq(t, v, 10)
			eq(t, found, true)
			eq(t, s.Size(), 1)

			v, found = s.TryRemove(0)
			eq(t, v, 30)
			eq(t, found, true)
			Empty(t, s)

			v, found = s.TryRemove(0)
			eq(t, found, false)
		},
	}
}

func TryCycle[Abstract interface {
	adt.Sizer
	adt.Emptier
	adt.Appender[int]
	TryCycle() (int, bool)
	TryHead() (int, bool)
}](f func() Abstract) Spec {
	return Spec{
		Name: "TryCycle",
		Test: func(t *testing.T) {
			c := f()

			v, found := c.TryCycle()
			eq(t, found, false)

			c.Append(1)
			c.Append(2)
			c.Append(3)

			v, found = c.TryCycle()
			eq(t, v, 1)
			eq(t, found, true)

			head, _ := c.TryHead()
			eq(t, head, 2)

			v, found = c.TryCycle()
			eq(t, v, 2)
			eq(t, found, true)

			v, found = c.TryCycle()
			eq(t, v, 3)
			eq(t, found, true)

			// Wraps around
			v, found = c.TryCycle()
			eq(t, v, 1)
			eq(t, found, true)
		},
	}
}

func TryReverseCycle[Abstract interface {
	adt.Sizer
	adt.Emptier
	adt.Appender[int]
	TryReverseCycle() (int, bool)
	TryHead() (int, bool)
}](f func() Abstract) Spec {
	return Spec{
		Name: "TryReverseCycle",
		Test: func(t *testing.T) {
			c := f()

			v, found := c.TryReverseCycle()
			eq(t, found, false)

			c.Append(1)
			c.Append(2)
			c.Append(3)

			v, found = c.TryReverseCycle()
			eq(t, v, 3)
			eq(t, found, true)

			v, found = c.TryReverseCycle()
			eq(t, v, 2)
			eq(t, found, true)

			v, found = c.TryReverseCycle()
			eq(t, v, 1)
			eq(t, found, true)

			// Wraps around
			v, found = c.TryReverseCycle()
			eq(t, v, 3)
			eq(t, found, true)
		},
	}
}

func TryPeekStack[Abstract interface {
	adt.Sizer
	adt.Emptier
	adt.Pusher[int]
	TryPeek() (int, bool)
}](f func() Abstract) Spec {
	return Spec{
		Name: "TryPeekStack",
		Test: func(t *testing.T) {
			s := f()

			v, found := s.TryPeek()
			eq(t, found, false)

			s.Push(10)
			v, found = s.TryPeek()
			eq(t, v, 10)
			eq(t, found, true)

			s.Push(20)
			v, found = s.TryPeek()
			eq(t, v, 20)
			eq(t, found, true)
		},
	}
}

func TryPopStack[Abstract interface {
	adt.Sizer
	adt.Emptier
	adt.Pusher[int]
	TryPop() (int, bool)
}](f func() Abstract) Spec {
	return Spec{
		Name: "TryPopStack",
		Test: func(t *testing.T) {
			s := f()

			v, found := s.TryPop()
			eq(t, found, false)

			s.Push(10)
			s.Push(20)
			s.Push(30)

			v, found = s.TryPop()
			eq(t, v, 30)
			eq(t, found, true)

			v, found = s.TryPop()
			eq(t, v, 20)
			eq(t, found, true)

			v, found = s.TryPop()
			eq(t, v, 10)
			eq(t, found, true)

			v, found = s.TryPop()
			eq(t, found, false)
		},
	}
}

func TryPeekQueue[Abstract interface {
	adt.Sizer
	adt.Emptier
	adt.Enqueuer[int]
	TryPeek() (int, bool)
}](f func() Abstract) Spec {
	return Spec{
		Name: "TryPeekQueue",
		Test: func(t *testing.T) {
			q := f()

			v, found := q.TryPeek()
			eq(t, found, false)

			q.Enqueue(10)
			v, found = q.TryPeek()
			eq(t, v, 10)
			eq(t, found, true)

			q.Enqueue(20)
			v, found = q.TryPeek()
			eq(t, v, 10)
			eq(t, found, true)
		},
	}
}

func TryDequeue[Abstract interface {
	adt.Sizer
	adt.Emptier
	adt.Enqueuer[int]
	TryDequeue() (int, bool)
}](f func() Abstract) Spec {
	return Spec{
		Name: "TryDequeue",
		Test: func(t *testing.T) {
			q := f()

			v, found := q.TryDequeue()
			eq(t, found, false)

			q.Enqueue(10)
			q.Enqueue(20)
			q.Enqueue(30)

			v, found = q.TryDequeue()
			eq(t, v, 10)
			eq(t, found, true)

			v, found = q.TryDequeue()
			eq(t, v, 20)
			eq(t, found, true)

			v, found = q.TryDequeue()
			eq(t, v, 30)
			eq(t, found, true)

			v, found = q.TryDequeue()
			eq(t, found, false)
		},
	}
}

func IterStack[Abstract interface {
	adt.Sizer
	adt.Emptier
	adt.Pusher[int]
	adt.Iterator[int]
}](f func() Abstract) Spec {
	return Spec{
		Name: "IterStack",
		Test: func(t *testing.T) {
			s := f()

			got := slices.Collect(s.Iter)
			eq(t, len(got), 0)

			s.Push(1)
			s.Push(2)
			s.Push(3)

			got = slices.Collect(s.Iter)
			ok(t, slices.Equal(got, []int{1, 2, 3}))

			count := 0
			for range s.Iter {
				count++
				if count == 2 {
					break
				}
			}
			eq(t, count, 2)

			eq(t, s.Size(), 3)
		},
	}
}

func IterQueue[Abstract interface {
	adt.Sizer
	adt.Emptier
	adt.Enqueuer[int]
	adt.Iterator[int]
}](f func() Abstract) Spec {
	return Spec{
		Name: "IterQueue",
		Test: func(t *testing.T) {
			q := f()

			got := slices.Collect(q.Iter)
			eq(t, len(got), 0)

			q.Enqueue(1)
			q.Enqueue(2)
			q.Enqueue(3)

			got = slices.Collect(q.Iter)
			ok(t, slices.Equal(got, []int{1, 2, 3}))

			count := 0
			for range q.Iter {
				count++
				if count == 2 {
					break
				}
			}
			eq(t, count, 2)

			eq(t, q.Size(), 3)
		},
	}
}

func BSTIterBackward[Abstract interface {
	adt.Emptier
	adt.Adder[int]
	adt.BackwardIterator[int]
}](f func() Abstract) Spec {
	return Spec{
		Name: "BSTIterBackward",
		Test: func(t *testing.T) {
			tree := f()

			got := slices.Collect(tree.IterBackward)
			eq(t, len(got), 0)

			tree.Add(5)
			tree.Add(3)
			tree.Add(7)
			tree.Add(1)
			tree.Add(4)

			got = slices.Collect(tree.IterBackward)
			ok(t, slices.Equal(got, []int{7, 5, 4, 3, 1}))

			count := 0
			for range tree.IterBackward {
				count++
				if count == 2 {
					break
				}
			}
			eq(t, count, 2)
		},
	}
}

func ok(t *testing.T, cond bool) {
	t.Helper()
	if !cond {
		t.Error("condition failed")
	}
}

func eq[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}
