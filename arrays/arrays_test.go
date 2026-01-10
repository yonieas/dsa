package arrays_test

import (
	"fmt"
	"os"
	"slices"
	"testing"

	"github.com/josestg/dsa/arrays"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func newArray[T any](t *testing.T, length int) *arrays.Array[T] {
	t.Helper()
	a := arrays.New[T](length)
	t.Cleanup(a.Free)
	return a
}

func assertPanics(t *testing.T, fn func()) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Error("expected panic")
		}
	}()
	fn()
}

func TestNew(t *testing.T) {
	t.Run("length", func(t *testing.T) {
		a := newArray[int](t, 5)
		if a.Len() != 5 {
			t.Errorf("Len() = %d, want 5", a.Len())
		}
	})

	t.Run("zero values", func(t *testing.T) {
		a := newArray[int](t, 3)
		for i := range a.Len() {
			if v := a.Get(i); v != 0 {
				t.Errorf("Get(%d) = %d, want 0", i, v)
			}
		}
	})

	t.Run("struct type", func(t *testing.T) {
		type point struct{ x, y int }
		a := newArray[point](t, 2)
		if v := a.Get(0); v != (point{}) {
			t.Errorf("Get(0) = %v, want zero", v)
		}
	})

	t.Run("panics on zero-sized type", func(t *testing.T) {
		assertPanics(t, func() { arrays.New[struct{}](1) })
	})
}

func TestGetSet(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		a := newArray[int](t, 5)
		a.Set(2, 42)

		if v := a.Get(2); v != 42 {
			t.Errorf("Get(2) = %d, want 42", v)
		}

		for i := range a.Len() {
			if i != 2 && a.Get(i) != 0 {
				t.Errorf("Get(%d) = %d, want 0", i, a.Get(i))
			}
		}
	})

	t.Run("value semantics", func(t *testing.T) {
		type box struct{ val int }
		a := newArray[box](t, 1)
		a.Set(0, box{10})

		got := a.Get(0)
		got.val = 999

		if v := a.Get(0); v.val != 10 {
			t.Errorf("array modified through copy: got %d, want 10", v.val)
		}
	})

	t.Run("bounds", func(t *testing.T) {
		a := newArray[int](t, 5)

		assertPanics(t, func() { a.Get(-1) })
		assertPanics(t, func() { a.Get(5) })
		assertPanics(t, func() { a.Set(-1, 0) })
		assertPanics(t, func() { a.Set(5, 0) })
	})
}

func TestIter(t *testing.T) {
	setup := func(t *testing.T) *arrays.Array[int] {
		a := newArray[int](t, 5)
		for i := range 5 {
			a.Set(i, (i+1)*10)
		}
		return a
	}

	t.Run("forward", func(t *testing.T) {
		a := setup(t)
		got := slices.Collect(a.Iter)
		want := []int{10, 20, 30, 40, 50}
		if !slices.Equal(got, want) {
			t.Errorf("Iter = %v, want %v", got, want)
		}
	})

	t.Run("backward", func(t *testing.T) {
		a := setup(t)
		got := slices.Collect(a.IterBackward)
		want := []int{50, 40, 30, 20, 10}
		if !slices.Equal(got, want) {
			t.Errorf("IterBackward = %v, want %v", got, want)
		}
	})

	t.Run("break", func(t *testing.T) {
		a := setup(t)
		count := 0
		for range a.Iter {
			count++
			if count == 3 {
				break
			}
		}
		if count != 3 {
			t.Errorf("count = %d, want 3", count)
		}
	})
}

func TestEnum(t *testing.T) {
	setup := func(t *testing.T) *arrays.Array[string] {
		a := newArray[string](t, 4)
		for i, v := range []string{"a", "b", "c", "d"} {
			a.Set(i, v)
		}
		return a
	}

	t.Run("forward", func(t *testing.T) {
		a := setup(t)
		var indices []int
		var values []string
		for i, v := range a.Enum {
			indices = append(indices, i)
			values = append(values, v)
		}

		if !slices.Equal(indices, []int{0, 1, 2, 3}) {
			t.Errorf("indices = %v, want [0 1 2 3]", indices)
		}
		if !slices.Equal(values, []string{"a", "b", "c", "d"}) {
			t.Errorf("values = %v, want [a b c d]", values)
		}
	})

	t.Run("backward", func(t *testing.T) {
		a := setup(t)
		var indices []int
		var values []string
		for i, v := range a.EnumBackward {
			indices = append(indices, i)
			values = append(values, v)
		}

		if !slices.Equal(indices, []int{3, 2, 1, 0}) {
			t.Errorf("indices = %v, want [3 2 1 0]", indices)
		}
		if !slices.Equal(values, []string{"d", "c", "b", "a"}) {
			t.Errorf("values = %v, want [d c b a]", values)
		}
	})
}

func TestString(t *testing.T) {
	tests := []struct {
		name   string
		values []int
		want   string
	}{
		{"multiple", []int{1, 2, 3, 4, 5}, "[1 2 3 4 5]"},
		{"single", []int{42}, "[42]"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := newArray[int](t, len(tt.values))
			for i, v := range tt.values {
				a.Set(i, v)
			}
			if got := a.String(); got != tt.want {
				t.Errorf("String() = %q, want %q", got, tt.want)
			}
		})
	}

	t.Run("matches Go format", func(t *testing.T) {
		a := newArray[int](t, 3)
		a.Set(0, 10)
		a.Set(1, 20)
		a.Set(2, 30)

		want := fmt.Sprint([3]int{10, 20, 30})
		if got := a.String(); got != want {
			t.Errorf("String() = %q, want %q", got, want)
		}
	})
}

func TestEmpty(t *testing.T) {
	t.Run("non-empty array", func(t *testing.T) {
		a := newArray[int](t, 5)
		if a.Empty() {
			t.Error("Empty() = true, want false for non-empty array")
		}
		if a.Size() != 5 {
			t.Errorf("Size() = %d, want 5", a.Size())
		}
	})

	t.Run("single element", func(t *testing.T) {
		a := newArray[int](t, 1)
		if a.Empty() {
			t.Error("Empty() = true, want false for single-element array")
		}
	})
}

func TestFree(t *testing.T) {
	t.Run("idempotent", func(t *testing.T) {
		a := arrays.New[int](5)
		a.Free()
		a.Free()
		a.Free()
	})
}

func TestGarbageCollected(t *testing.T) {
	t.Run("enabled", func(t *testing.T) {
		a := arrays.NewGarbageCollected[int](5, true)
		a.Set(0, 100)
		if v := a.Get(0); v != 100 {
			t.Errorf("Get(0) = %d, want 100", v)
		}
	})

	t.Run("disabled", func(t *testing.T) {
		a := arrays.NewGarbageCollected[int](5, false)
		defer a.Free()

		a.Set(0, 200)
		if v := a.Get(0); v != 200 {
			t.Errorf("Get(0) = %d, want 200", v)
		}
	})
}

func TestComplexStruct(t *testing.T) {
	type record struct {
		active bool
		count  int
		name   string
	}

	data := []record{
		{true, 1, "one"},
		{false, 2, "two"},
		{true, 4, "three"},
	}

	a := newArray[record](t, len(data))
	for i, v := range data {
		a.Set(i, v)
	}

	for i, want := range data {
		if got := a.Get(i); got != want {
			t.Errorf("Get(%d) = %v, want %v", i, got, want)
		}
	}
}
