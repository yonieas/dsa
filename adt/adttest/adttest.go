package adttest

import (
	"cmp"
	"fmt"
	"math"
	"math/rand/v2"
	"slices"
	"sort"
	"testing"

	"github.com/josestg/dsa/adt"
	"github.com/josestg/dsa/sequence"
	"github.com/stretchr/testify/assert"
)

type Runner func(t *testing.T)

func (r Runner) Run(t *testing.T) {
	t.Helper()
	r(t)
}

type Generator[T any] func() T

func (g Generator[T]) New() T { return g() }

func randSample() int {
	return max(8, rand.IntN(64))
}

func AppendSimulator[
	E any,
	Abstract interface {
		adt.Sizer
		adt.Tailer[E]
		adt.Appender[E]
	},
	Constructor func() Abstract,
](
	c Constructor,
	g Generator[E],
	destructors ...func(Abstract),
) Runner {
	return func(t *testing.T) {
		t.Helper()
		a := c()
		setCleanup(t, a, destructors)
		assert.Equal(t, 0, a.Size())
		assert.Panics(t, func() { _ = a.Tail() })
		n := randSample()
		for range n {
			v := g.New()
			a.Append(v)
			assert.Equal(t, v, a.Tail())
		}
		assert.Equal(t, n, a.Size())
	}
}

func PrependSimulator[
	E any,
	Abstract interface {
		adt.Sizer
		adt.Header[E]
		adt.Prepender[E]
	},
	Constructor func() Abstract,
](
	c Constructor,
	g Generator[E],
	destructors ...func(Abstract),
) Runner {
	return func(t *testing.T) {
		t.Helper()
		a := c()
		setCleanup(t, a, destructors)
		assert.Equal(t, 0, a.Size())
		assert.Panics(t, func() { _ = a.Head() })
		n := randSample()
		for range n {
			v := g.New()
			a.Prepend(v)
			assert.Equal(t, v, a.Head())
		}
		assert.Equal(t, n, a.Size())
	}
}

func PopSimulator[
	E any,
	Abstract interface {
		adt.Sizer
		adt.Emptier
		adt.Popper[E]
		adt.Appender[E]
	},
	Constructor func() Abstract,
](
	c Constructor,
	g Generator[E],
	destructors ...func(Abstract),
) Runner {
	return func(t *testing.T) {
		t.Helper()
		a := c()
		setCleanup(t, a, destructors)

		assert.True(t, a.Empty())
		assert.Panics(t, func() { _ = a.Pop() })

		n := randSample()
		s1 := make([]E, 0, n)
		for range n {
			v := g.New()
			a.Append(v)
			s1 = append(s1, v)
		}

		assert.Equal(t, n, a.Size())
		assert.False(t, a.Empty())

		s2 := make([]E, 0, n)
		for range n {
			s2 = append(s2, a.Pop())
		}

		slices.Reverse(s1)
		assert.Equal(t, s1, s2)
		assert.True(t, a.Empty())
	}
}

func ShiftSimulator[
	E any,
	Abstract interface {
		adt.Sizer
		adt.Emptier
		adt.Shifter[E]
		adt.Prepender[E]
	},
	Constructor func() Abstract,
](
	c Constructor,
	g Generator[E],
	destructors ...func(Abstract),
) Runner {
	return func(t *testing.T) {
		t.Helper()
		a := c()
		setCleanup(t, a, destructors)

		assert.True(t, a.Empty())
		assert.Panics(t, func() { _ = a.Shift() })

		n := randSample()
		s1 := make([]E, 0, n)
		for range n {
			v := g.New()
			a.Prepend(v)
			s1 = append(s1, v)
		}

		assert.Equal(t, n, a.Size())
		assert.False(t, a.Empty())

		s2 := make([]E, 0, n)
		for range n {
			s2 = append(s2, a.Shift())
		}

		slices.Reverse(s1)
		assert.Equal(t, s1, s2)
		assert.True(t, a.Empty())
	}
}

func GetSimulator[
	E any,
	Abstract interface {
		adt.Sizer
		adt.Getter[E]
		adt.Appender[E]
	},
	Constructor func() Abstract,
](
	c Constructor,
	g Generator[E],
	destructors ...func(Abstract),
) Runner {
	return func(t *testing.T) {
		t.Helper()
		a := c()
		setCleanup(t, a, destructors)

		assert.Zero(t, a.Size())
		assert.Panics(t, func() { _ = a.Get(0) })

		size := a.Size()
		assert.Panics(t, func() { _ = a.Get(-1) })
		assert.Panics(t, func() { _ = a.Get(size) })

		n := randSample()
		s := make([]E, 0, n)
		for range n {
			v := g.New()
			a.Append(v)
			s = append(s, v)
			assert.Equal(t, v, a.Get(a.Size()-1))
		}

		for i := 0; i < n; i++ {
			assert.Equal(t, s[i], a.Get(size+i))
		}
	}
}

func SetSimulator[
	E any,
	Abstract interface {
		adt.Sizer
		adt.Getter[E]
		adt.Setter[E]
		adt.Appender[E]
	},
	Constructor func() Abstract,
](
	c Constructor,
	g Generator[E],
	destructors ...func(Abstract),
) Runner {
	return func(t *testing.T) {
		t.Helper()

		a := c()
		setCleanup(t, a, destructors)

		assert.Zero(t, a.Size())
		assert.Panics(t, func() { a.Set(0, g.New()) })
		assert.Panics(t, func() { a.Set(-1, g.New()) })
		assert.Panics(t, func() { a.Set(a.Size(), g.New()) })

		n := randSample()
		s := make([]E, 0, n)
		for range n {
			v := g.New()
			a.Append(v)
			s = append(s, v)
			assert.Equal(t, v, a.Get(a.Size()-1))
		}

		for i := range n {
			assert.Equal(t, s[i], a.Get(i))
		}

		for i := range n {
			v := g.New()
			a.Set(i, v)
			s[i] = v
		}

		for i := range n {
			assert.Equal(t, s[i], a.Get(i))
		}
	}
}

func IterSimulator[
	E any,
	Abstract interface {
		adt.Sizer
		adt.Appender[E]
		adt.Iterator[E]
	},
	Constructor func() Abstract,
](
	c Constructor,
	g Generator[E],
	destructors ...func(Abstract),
) Runner {
	return func(t *testing.T) {
		t.Helper()
		a := c()
		setCleanup(t, a, destructors)
		assert.Zero(t, a.Size())

		n := randSample()
		s := make([]E, 0, n)
		for range n {
			v := g.New()
			a.Append(v)
			s = append(s, v)
		}

		assert.Equal(t, n, a.Size())
		assert.Equal(t, s, slices.Collect(a.Iter))

		// break
		for i := range sequence.Enum(a.Iter) {
			if i > n/2 {
				break
			}
		}
	}
}

func IterBackwardSimulator[
	E any,
	Abstract interface {
		adt.Sizer
		adt.Appender[E]
		adt.BackwordIterator[E]
	},
	Constructor func() Abstract,
](
	c Constructor,
	g Generator[E],
	destructors ...func(Abstract),
) Runner {
	return func(t *testing.T) {
		t.Helper()
		a := c()
		setCleanup(t, a, destructors)
		assert.Zero(t, a.Size())

		n := randSample()
		s := make([]E, 0, n)
		for range n {
			v := g.New()
			a.Append(v)
			s = append(s, v)
		}

		assert.Equal(t, n, a.Size())

		slices.Reverse(s)
		assert.Equal(t, s, slices.Collect(a.IterBackward))

		// break
		for i := range sequence.Enum(a.IterBackward) {
			if i > n/2 {
				break
			}
		}
	}
}

func BracketStringSimulator[
	E any,
	Abstract interface {
		adt.Sizer
		adt.Appender[E]
		fmt.Stringer
	},
	Constructor func() Abstract,
](
	c Constructor,
	g Generator[E],
	destructors ...func(Abstract),
) Runner {
	return func(t *testing.T) {
		t.Helper()
		a := c()
		setCleanup(t, a, destructors)

		assert.Zero(t, a.Size())
		assert.Equal(t, "[]", a.String())

		n := randSample()
		s := make([]E, 0, n)
		for range n {
			v := g.New()
			a.Append(v)
			s = append(s, v)
		}

		assert.Equal(t, n, a.Size())
		assert.Equal(t, fmt.Sprint(s), a.String())
	}
}

func SortSimulator[
	E cmp.Ordered,
	Abstract interface {
		adt.Sizer
		adt.Getter[E]
		adt.Setter[E]
		adt.Appender[E]
	},
	Constructor func() Abstract,
](
	c Constructor,
	g Generator[E],
	destructors ...func(Abstract),
) Runner {
	return func(t *testing.T) {
		t.Helper()
		a := c()
		setCleanup(t, a, destructors)
		assert.Zero(t, a.Size())

		n := randSample()
		s := make([]E, 0, n)
		for range n {
			v := g.New()
			a.Append(v)
			s = append(s, v)
		}

		sa := sortable[E]{abstract: a}
		sort.Sort(&sa)
		sort.IsSorted(&sa)
		slices.Sort(s)
		for i := range n {
			assert.Equal(t, sa.abstract.Get(i), s[i])
		}
	}
}

type sortable[E cmp.Ordered] struct {
	abstract interface {
		adt.Sizer
		adt.Getter[E]
		adt.Setter[E]
	}
}

func (s *sortable[E]) Len() int { return s.abstract.Size() }

func (s *sortable[E]) Less(i, j int) bool {
	return cmp.Less(s.abstract.Get(i), s.abstract.Get(j))
}

func (s *sortable[E]) Swap(i, j int) {
	x, y := s.abstract.Get(i), s.abstract.Get(j)
	s.abstract.Set(i, y)
	s.abstract.Set(j, x)
}

func DoublingGrowSimulator[
	E cmp.Ordered,
	Abstract interface {
		adt.Sizer
		adt.Caper
		adt.Appender[E]
		adt.Prepender[E]
		adt.Iterator[E]
	},
	Constructor func() Abstract,
](
	c Constructor,
	g Generator[E],
	destructors ...func(Abstract),
) Runner {
	return func(t *testing.T) {
		t.Helper()
		t.Run("grow by prepend", func(t *testing.T) {
			a := c()
			setCleanup(t, a, destructors)
			assert.Zero(t, a.Size())
			assert.Equal(t, 1, a.Cap())

			n := 255 // the threshold - 1.
			s := make([]E, 0, 1)
			for i := 0; i < n; i++ {
				v := g.New()
				a.Prepend(v)
				s = append(s, v)
			}

			numOfGrow := int(math.Floor(math.Log2(float64(n)))) + 1
			assert.Equal(t, 1<<numOfGrow, a.Cap())
			assert.Equal(t, cap(s), a.Cap())

			slices.Reverse(s)
			assert.Equal(t, s, slices.Collect(a.Iter))
		})

		t.Run("grow by append", func(t *testing.T) {
			a := c()
			setCleanup(t, a, destructors)
			assert.Zero(t, a.Size())
			assert.Equal(t, 1, a.Cap())

			n := 255 // the threshold - 1.
			s := make([]E, 0, 1)
			for i := 0; i < n; i++ {
				v := g.New()
				a.Append(v)
				s = append(s, v)
			}

			numOfGrow := int(math.Floor(math.Log2(float64(n)))) + 1
			assert.Equal(t, 1<<numOfGrow, a.Cap())
			assert.Equal(t, cap(s), a.Cap())
			assert.Equal(t, s, slices.Collect(a.Iter))
		})
	}
}

func StackSimulator[
	E any,
	Abstract interface {
		adt.Stack[E]
	},
	Constructor func() Abstract,
](
	c Constructor,
	g Generator[E],
	destructors ...func(Abstract),
) Runner {
	return func(t *testing.T) {
		t.Helper()
		a := c()
		setCleanup(t, a, destructors)

		assert.Zero(t, a.Size())
		assert.True(t, a.Empty())
		assert.Panics(t, func() { _ = a.Pop() })
		assert.Panics(t, func() { _ = a.Peek() })

		n := randSample()
		for range n {
			v := g.New()
			a.Push(v)
			assert.Equal(t, v, a.Peek())
		}

		assert.Equal(t, n, a.Size())
		assert.False(t, a.Empty())

		for !a.Empty() {
			peek := a.Peek()
			assert.Equal(t, peek, a.Pop())
		}

		assert.Zero(t, a.Size())
		assert.True(t, a.Empty())
	}
}

func QueueSimulator[
	E any,
	Abstract interface {
		adt.Queue[E]
	},
	Constructor func() Abstract,
](
	c Constructor,
	g Generator[E],
	destructors ...func(Abstract),
) Runner {
	return func(t *testing.T) {
		t.Helper()
		a := c()
		setCleanup(t, a, destructors)

		assert.Zero(t, a.Size())
		assert.True(t, a.Empty())
		assert.Panics(t, func() { _ = a.Dequeue() })
		assert.Panics(t, func() { _ = a.Peek() })

		n := randSample()
		s1 := make([]E, 0, n)
		for range n {
			v := g.New()
			a.Enqueue(v)
			s1 = append(s1, v)
			assert.Equal(t, s1[0], a.Peek())
		}

		assert.Equal(t, n, a.Size())
		assert.False(t, a.Empty())

		for !a.Empty() {
			peek := a.Peek()
			assert.Equal(t, peek, a.Dequeue())
		}

		assert.Zero(t, a.Size())
		assert.True(t, a.Empty())
	}
}

func InsertRemoveSimulator[
	E comparable,
	Abstract interface {
		adt.Sizer
		adt.Inserter[E]
		adt.Remover[E]
		adt.Getter[E]
	},
	Constructor func() Abstract,
](
	c Constructor,
	g Generator[E],
	destructors ...func(Abstract),
) Runner {
	return func(t *testing.T) {
		t.Helper()

		a := c()
		setCleanup(t, a, destructors)
		assert.Zero(t, a.Size())

		v1 := g.New()
		a.Insert(0, v1)
		assert.Equal(t, v1, a.Get(0))
		assert.Equal(t, 1, a.Size())

		assert.Equal(t, v1, a.Remove(0))
		assert.Zero(t, a.Size())

		v2 := g.New()
		a.Insert(0, v1)
		a.Insert(a.Size(), v2)
		assert.Equal(t, v1, a.Get(0))
		assert.Equal(t, v2, a.Get(a.Size()-1))
		assert.Equal(t, 2, a.Size())

		v3 := g.New()
		a.Insert(1, v3)
		assert.Equal(t, v1, a.Get(0))
		assert.Equal(t, v3, a.Get(1))
		assert.Equal(t, v2, a.Get(2))

		v4 := g.New()
		a.Insert(2, v4)
		assert.Equal(t, v1, a.Get(0))
		assert.Equal(t, v3, a.Get(1))
		assert.Equal(t, v4, a.Get(2))
		assert.Equal(t, v2, a.Get(3))
		assert.Equal(t, 4, a.Size())

		assert.Equal(t, v3, a.Remove(1))
		assert.Equal(t, v4, a.Remove(1))
		assert.Equal(t, v2, a.Remove(1))
		assert.Equal(t, v1, a.Remove(0))
		assert.Zero(t, a.Size())
	}
}

func HashSetSimulator[
	T comparable,
	Set interface {
		adt.Sizer
		adt.Emptier
		adt.Adder[T]
		adt.Deleter[T]
		adt.Exister[T]
		adt.Iterator[T]
	},
	Constructor func() Set,
](
	c Constructor,
	g Generator[T],
	destructors ...func(Set),
) Runner {
	return func(t *testing.T) {
		t.Helper()

		set := c()
		setCleanup(t, set, destructors)
		assert.Zero(t, set.Size())
		assert.True(t, set.Empty())

		var truth = map[T]struct{}{}
		n := randSample()

		for range n {
			v := g.New()
			set.Add(v)
			truth[v] = struct{}{}

			assert.True(t, set.Exists(v))
			assert.Equal(t, len(truth), set.Size())
		}

		for k := range truth {
			assert.True(t, set.Exists(k))
		}
		set.Iter(func(v T) bool {
			_, ok := truth[v]
			assert.True(t, ok)
			return true
		})

		half := 0
		for k := range truth {
			if half >= len(truth)/2 {
				break
			}
			set.Del(k)
			delete(truth, k)
			half++
		}

		assert.Equal(t, len(truth), set.Size())
		for k := range truth {
			assert.True(t, set.Exists(k))
		}
	}
}

func HashMapSimulator[
	K comparable,
	V comparable,
	Map interface {
		Put(K, V)
		Get(K) (V, bool)
		adt.Sizer
		adt.Emptier
		adt.Deleter[K]
		adt.Keys[K]
		adt.Exister[K]
	},
	Constructor func() Map,
](
	c Constructor,
	keyGen Generator[K],
	valGen Generator[V],
	destructors ...func(Map),
) Runner {
	return func(t *testing.T) {
		t.Helper()

		m := c()
		setCleanup(t, m, destructors)

		assert.True(t, m.Empty())
		assert.Zero(t, m.Size())

		truth := make(map[K]V)
		n := randSample()

		for range n {
			k := keyGen.New()
			v := valGen.New()
			m.Put(k, v)
			truth[k] = v
			got, ok := m.Get(k)
			assert.True(t, ok)
			assert.Equal(t, v, got)
		}

		assert.Equal(t, len(truth), m.Size())

		for k, v := range truth {
			assert.True(t, m.Exists(k))
			got, ok := m.Get(k)
			assert.True(t, ok)
			assert.Equal(t, v, got)
		}

		iterCount := 0
		for k := range m.Keys {
			expected, ok := truth[k]
			assert.True(t, ok)
			v, ok := m.Get(k)
			assert.True(t, ok)
			assert.Equal(t, expected, v)
			iterCount++
		}
		assert.Equal(t, len(truth), iterCount)

		deleted := 0
		for k := range truth {
			if deleted >= len(truth)/2 {
				break
			}
			m.Del(k)
			delete(truth, k)
			assert.False(t, m.Exists(k))
			_, ok := m.Get(k)
			assert.False(t, ok)
			deleted++
		}
		assert.Equal(t, len(truth), m.Size())
	}
}

func setCleanup[Abstract any](t *testing.T, a Abstract, destructors []func(Abstract)) {
	t.Cleanup(func() {
		if len(destructors) > 0 {
			destructors[0](a)
		}
	})
}
