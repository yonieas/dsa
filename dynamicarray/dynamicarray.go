package dynamicarray

import (
	"github.com/josestg/dsa/arrays"
	"github.com/josestg/dsa/internal/generics"
	"github.com/josestg/dsa/sequence"
)

type DynamicArray[T any] struct {
	backend *arrays.Array[T]
	size    int
}

func New[T any](capacity int) *DynamicArray[T] {
	if capacity <= 0 {
		panic("DynamicArray.New: must have at minimum 1 capacity")
	}
	return &DynamicArray[T]{
		backend: arrays.New[T](capacity),
		size:    0,
	}
}

func (d *DynamicArray[T]) Free() {
	if d.backend != nil {
		d.backend.Free()
		d.backend = nil
		d.size = 0
	}
}

func (d *DynamicArray[T]) Empty() bool { return d.Size() == 0 }

func (d *DynamicArray[T]) Size() int { return d.size }

func (d *DynamicArray[T]) Cap() int { return d.backend.Len() }

func (d *DynamicArray[T]) Tail() T {
	if d.Empty() {
		panic("DynamicArray.Tail: is empty array")
	}
	return d.Get(d.Size() - 1)
}

func (d *DynamicArray[T]) Head() T {
	if d.Empty() {
		panic("DynamicArray.Head: is empty array")
	}
	return d.Get(0)
}

func (d *DynamicArray[T]) Get(index int) T {
	d.checkBounds(index)
	return d.backend.Get(index)
}

func (d *DynamicArray[T]) Set(index int, value T) {
	d.checkBounds(index)
	d.backend.Set(index, value)
}

func (d *DynamicArray[T]) Prepend(value T) {
	d.Append(value)
	for i := d.size - 1; i > 0; i-- {
		d.Swap(i, i-1)
	}
}

func (d *DynamicArray[T]) Shift() T {
	if v, ok := d.TryShift(); !ok {
		panic("DynamicArray.Shift: array is empty")
	} else {
		return v
	}
}

func (d *DynamicArray[T]) TryShift() (T, bool) {
	var zero T
	n := d.Size()
	if n == 0 {
		return zero, false
	}

	v := d.Get(0)
	d.Set(0, zero)
	for i := 0; i < n-1; i++ {
		d.Swap(i, i+1)
	}
	d.size--
	return v, true
}

func (d *DynamicArray[T]) Swap(i, j int) {
	if i != j {
		x, y := d.Get(i), d.Get(j)
		d.Set(i, y)
		d.Set(j, x)
	}
}

func (d *DynamicArray[T]) Append(value T) {
	c := d.Cap()
	if d.size >= c {
		d.grow()
	}
	d.backend.Set(d.size, value)
	d.size++
}

func (d *DynamicArray[T]) grow() {
	// Go slice implementation only doubles the capacity if the current size is less than 256.
	// See: https://cs.opensource.google/go/go/+/refs/tags/go1.24.2:src/runtime/slice.go;l=289-322
	//
	// we can do the same with some this simple approximation: oldCap + (oldCap + 3*256) / 4
	// See: https://victoriametrics.com/blog/go-slice/
	const threshold = 256
	capacity := d.Cap()
	newCapacity := 2 * capacity
	if capacity >= threshold {
		newCapacity = capacity + (capacity+3*threshold)/4
	}
	newBackend := arrays.New[T](newCapacity)
	for i, v := range d.backend.Iter(false) {
		newBackend.Set(i, v)
	}
	d.backend.Free()
	d.backend = newBackend
}

func (d *DynamicArray[T]) Pop() T {
	if v, ok := d.TryPop(); !ok {
		panic("DynamicArray.Pop: array is empty")
	} else {
		return v
	}
}

func (d *DynamicArray[T]) TryPop() (T, bool) {
	var zero T
	if d.Size() == 0 {
		return zero, false
	}
	val := d.backend.Get(d.size - 1)
	d.backend.Set(d.size-1, zero) // clear the slot.
	d.size--
	return val, true
}

func (d *DynamicArray[T]) Clip() {
	if d.Empty() {
		panic("DynamicArray.Clip: array is empty")
	}

	if d.size == d.Cap() {
		return
	}

	newBackend := arrays.New[T](d.size)
	for i := range d.Size() {
		newBackend.Set(i, d.backend.Get(i))
	}
	d.backend.Free()
	d.backend = newBackend
}

func (d *DynamicArray[T]) Iter(yield func(T) bool) {
	for i := range d.Size() {
		if !yield(d.Get(i)) {
			break
		}
	}
}

func (d *DynamicArray[T]) IterBackward(yield func(T) bool) {
	for i := d.Size() - 1; i >= 0; i-- {
		if !yield(d.Get(i)) {
			break
		}
	}
}

func (d *DynamicArray[T]) String() string {
	return sequence.String(d.Iter)
}

func (d *DynamicArray[T]) Insert(index int, value T) {
	if index < 0 || index > d.size {
		panic("DynamicArray.Insert: index out of range")
	}

	if d.size >= d.Cap() {
		d.grow()
	}

	for i := d.size; i > index; i-- {
		d.backend.Set(i, d.backend.Get(i-1))
	}

	d.backend.Set(index, value)
	d.size++
}

func (d *DynamicArray[T]) Remove(index int) T {
	d.checkBounds(index)

	val := d.backend.Get(index)

	for i := index; i < d.size-1; i++ {
		d.backend.Set(i, d.backend.Get(i+1))
	}

	d.backend.Set(d.size-1, generics.ZeroValue[T]())
	d.size--
	return val
}

func (d *DynamicArray[T]) checkBounds(index int) {
	if index < 0 || index >= d.Size() {
		panic("dynamicarray: index out of range")
	}
}
