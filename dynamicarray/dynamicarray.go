// Package dynamicarray provides a resizable array implementation.
//
// # What is a Dynamic Array?
//
// A dynamic array solves the biggest limitation of fixed arrays: you don't
// need to know the size upfront. It grows automatically as you add elements.
//
// The clever trick is "amortized doubling": when the array fills up, we
// allocate a new array with 2x the capacity and copy everything over.
// While a single resize is expensive O(n), it happens so rarely that the
// average cost per insertion is still O(1). This is called amortized analysis.
//
// # How It Works
//
//  1. Start with initial capacity (e.g., 4 elements)
//  2. Add elements until full
//  3. When full: allocate 2x capacity, copy elements, free old array
//  4. Continue adding elements
//
// Example of growth: capacity 4, 8, 16, 32, 64, ...
//
// # Why Doubling?
//
// Why double instead of adding a fixed amount (like +10 each time)?
//
// If we add +10 each time and insert n elements, we'd resize n/10 times,
// doing O(n^2) total work. With doubling, we resize log2(n) times, doing
// O(n) total work. The math works out to O(1) amortized per insertion.
//
// Go's built-in slices use a similar strategy (with optimizations for
// large slices to reduce memory waste).
//
// # When to Use
//
// Use dynamic arrays when you need array-like access but unknown final size,
// most operations are append (add to end), or you want cache-friendly
// sequential access.
//
// Consider linked lists when you have frequent insertions/deletions at front
// or need stable pointers to elements.
//
// # Complexity
//
//	Access:    O(1)
//	Append:    O(1) amortized, O(n) worst case on resize
//	Prepend:   O(n)
//	Insert:    O(n)
//	Delete:    O(n)
//
// # Further Reading
//
// CLRS "Introduction to Algorithms", Chapter 17 (Amortized Analysis).
// Go Blog: "Go Slices: usage and internals".
// https://en.wikipedia.org/wiki/Dynamic_array
package dynamicarray

import (
	"github.com/josestg/dsa/arrays"
	"github.com/josestg/dsa/sequence"
)

// DynamicArray is a resizable array that grows automatically.
//
//	capacity = 8
//	┌───┬───┬───┬───┬───┬───┬───┬───┐
//	│ A │ B │ C │ D │ E │   │   │   │  <- 3 empty slots
//	└───┴───┴───┴───┴───┴───┴───┴───┘
//	          size = 5
//
// When size reaches capacity, a resize doubles the capacity.
type DynamicArray[T any] struct {
	backend *arrays.Array[T]
	size    int
}

// New creates an empty DynamicArray with given initial capacity.
//
//	capacity = 4
//	┌───┬───┬───┬───┐
//	│   │   │   │   │  ← all empty
//	└───┴───┴───┴───┘
//	     size = 0
//
// complexity:
//   - time : O(capacity)
//   - space: O(capacity)
//
// Panics if capacity <= 0.
func New[T any](capacity int) *DynamicArray[T] {
	if capacity <= 0 {
		panic("dynamicarray.New: must have at minimum 1 capacity")
	}
	b := arrays.NewGarbageCollected[T](capacity, true)
	return &DynamicArray[T]{
		backend: b,
		size:    0,
	}
}

// Empty returns true if the array has no elements.
//
// complexity:
//   - time : O(1)
//   - space: O(1)
func (d *DynamicArray[T]) Empty() bool { return d.Size() == 0 }

// Size returns the number of elements in the array.
//
//	┌───┬───┬───┬───┬───┬───┬───┬───┐
//	│ A │ B │ C │ D │ E │   │   │   │
//	└───┴───┴───┴───┴───┴───┴───┴───┘
//	     size = 5, cap = 8
//
// complexity:
//   - time : O(1)
//   - space: O(1)
//
// SCORE: 5
func (d *DynamicArray[T]) Size() int {
	// hint: return the size field
	return d.size
}

// Cap returns the current capacity.
//
// complexity:
//   - time : O(1)
//   - space: O(1)
//
// SCORE: 5
func (d *DynamicArray[T]) Cap() int {
	// hint: use d.backend.Size() - the backend array size is the capacity
	return d.backend.Size()
}

// Tail returns the last element without removing it.
//
//	┌───┬───┬───┬───┬───┐
//	│ A │ B │ C │ D │ E │
//	└───┴───┴───┴───┴───┘
//	                  ↑
//	             Tail() -> E
//
// complexity:
//   - time : O(1)
//   - space: O(1)
//
// Panics if the array is empty.
func (d *DynamicArray[T]) Tail() T {
	if v, ok := d.TryTail(); !ok {
		panic("DynamicArray.Tail: is empty array")
	} else {
		return v
	}
}

// TryTail attempts to return the last element.
// Returns (value, true) on success, or (zero, false) if empty.
//
// complexity:
//   - time : O(1)
//   - space: O(1)
//
// SCORE: 10
func (d *DynamicArray[T]) TryTail() (T, bool) {
	// hint: check if empty, then return element at index (size - 1)
	if d.Empty() {
		var arr T
		return arr, false
	}
	return d.Get(d.Size() - 1), true
}

// Head returns the first element without removing it.
//
//	┌───┬───┬───┬───┬───┐
//	│ A │ B │ C │ D │ E │
//	└───┴───┴───┴───┴───┘
//	  ↑
//	Head() -> A
//
// complexity:
//   - time : O(1)
//   - space: O(1)
//
// Panics if the array is empty.
func (d *DynamicArray[T]) Head() T {
	v, ok := d.TryHead()
	if !ok {
		panic("DynamicArray.Head: is empty array")
	}
	return v
}

// TryHead attempts to return the first element.
// Returns (value, true) on success, or (zero, false) if empty.
//
// complexity:
//   - time : O(1)
//   - space: O(1)
func (d *DynamicArray[T]) TryHead() (T, bool) {
	if d.Empty() {
		var zero T
		return zero, false
	}
	return d.backend.Get(0), true
}

// Get retrieves the element at the given index.
//
//	┌───┬───┬───┬───┬───┐
//	│ A │ B │ C │ D │ E │
//	└───┴───┴───┴───┴───┘
//	  0   1   2   3   4
//
//	Get(2) -> C
//
// complexity:
//   - time : O(1) - direct memory access
//   - space: O(1)
//
// Panics if index < 0 or index >= Size().
func (d *DynamicArray[T]) Get(index int) T {
	d.checkBounds(index)
	return d.backend.Get(index)
}

// TryGet attempts to retrieve the element at the given index.
// Returns (value, true) on success, or (zero, false) if index is out of bounds.
//
// complexity:
//   - time : O(1)
//   - space: O(1)
func (d *DynamicArray[T]) TryGet(index int) (T, bool) {
	if index < 0 || index >= d.Size() {
		var zero T
		return zero, false
	}
	return d.backend.Get(index), true
}

// Set updates the element at the given index.
//
//	Before Set(2, X):          After Set(2, X):
//	┌───┬───┬───┬───┬───┐      ┌───┬───┬───┬───┬───┐
//	│ A │ B │ C │ D │ E │      │ A │ B │ X │ D │ E │
//	└───┴───┴───┴───┴───┘      └───┴───┴───┴───┴───┘
//	  0   1   2   3   4          0   1   2   3   4
//
// complexity:
//   - time : O(1)
//   - space: O(1)
//
// Panics if index < 0 or index >= Size().
func (d *DynamicArray[T]) Set(index int, value T) {
	d.checkBounds(index)
	d.backend.Set(index, value)
}

// TrySet attempts to update the element at the given index.
// Returns true on success, false if index is out of bounds.
//
// complexity:
//   - time : O(1)
//   - space: O(1)
func (d *DynamicArray[T]) TrySet(index int, value T) bool {
	if index < 0 || index >= d.Size() {
		return false
	}
	d.backend.Set(index, value)
	return true
}

// Prepend adds an element to the front of the array.
//
//	Before Prepend(Z):         After Prepend(Z):
//	┌───┬───┬───┬───┬───┐      ┌───┬───┬───┬───┬───┬───┐
//	│ A │ B │ C │ D │ E │      │ Z │ A │ B │ C │ D │ E │
//	└───┴───┴───┴───┴───┘      └───┴───┴───┴───┴───┴───┘
//
// complexity:
//   - time : O(n) - must shift all elements
//   - space: O(1)
//
// Note: Use Append for O(1) amortized insertion.
//
// SCORE: 15
func (d *DynamicArray[T]) Prepend(value T) {
	// hint: 1) grow if size >= capacity
	//       2) shift all elements right by 1 (from end to start)
	//       3) set value at index 0
	//       4) increment size
	if d.Size() >= d.Cap() {
		d.grow()
	}
	d.Append(value)
	for i := d.size - 1; i > 0; i-- {
		d.Swap(i, i-1)
	}
	d.backend.Set(0, value)
	// d.size++
}

// Shift removes and returns the first element.
//
//	Before Shift():            After Shift():
//	┌───┬───┬───┬───┬───┐      ┌───┬───┬───┬───┐
//	│ A │ B │ C │ D │ E │      │ B │ C │ D │ E │
//	└───┴───┴───┴───┴───┘      └───┴───┴───┴───┘
//
//	Shift() -> A
//
// complexity:
//   - time : O(n) - must shift all elements
//   - space: O(1)
//
// Panics if the array is empty.
func (d *DynamicArray[T]) Shift() T {
	if v, ok := d.TryShift(); !ok {
		panic("DynamicArray.Shift: array is empty")
	} else {
		return v
	}
}

// TryShift attempts to remove and return the first element.
// Returns (value, true) on success, or (zero, false) if empty.
//
// complexity:
//   - time : O(n)
//   - space: O(1)
//
// SCORE: 15
func (d *DynamicArray[T]) TryShift() (T, bool) {
	// hint: 1) check if empty, return (zero, false)
	//       2) save element at index 0
	//       3) shift all elements left by 1
	//       4) clear last slot (for GC), decrement size
	//       5) return (saved, true)
	var arr T
	if d.Empty() {
		return arr, false
	}
	e := d.Get(0)
	for i := 0; i < d.size-1; i++ {
		d.Swap(i, i+1)
	}
	d.size--
	return e, true
}

// Swap exchanges elements at two indices.
//
//	Before Swap(1, 3):         After Swap(1, 3):
//	┌───┬───┬───┬───┬───┐      ┌───┬───┬───┬───┬───┐
//	│ A │ B │ C │ D │ E │      │ A │ D │ C │ B │ E │
//	└───┴───┴───┴───┴───┘      └───┴───┴───┴───┴───┘
//	      ↑       ↑                  ↑       ↑
//	    swapped positions
//
// complexity:
//   - time : O(1)
//   - space: O(1)
func (d *DynamicArray[T]) Swap(i, j int) {
	if i != j {
		x, y := d.Get(i), d.Get(j)
		d.Set(i, y)
		d.Set(j, x)
	}
}

// Append adds an element to the end of the array.
//
//	Before Append(F):          After Append(F):
//	┌───┬───┬───┬───┬───┬───┐  ┌───┬───┬───┬───┬───┬───┐
//	│ A │ B │ C │ D │ E │   │  │ A │ B │ C │ D │ E │ F │
//	└───┴───┴───┴───┴───┴───┘  └───┴───┴───┴───┴───┴───┘
//	     size=5, cap=6              size=6, cap=6
//
// If size == capacity, the array doubles in size first:
//
//	Before Append(G) (full):   After resize + Append(G):
//	┌───┬───┬───┬───┬───┬───┐  ┌───┬───┬───┬───┬───┬───┬───┬─...─┐
//	│ A │ B │ C │ D │ E │ F │  │ A │ B │ C │ D │ E │ F │ G │     │
//	└───┴───┴───┴───┴───┴───┘  └───┴───┴───┴───┴───┴───┴───┴─...─┘
//	     size=6, cap=6              size=7, cap=12
//
// complexity:
//   - time : O(1) amortized
//   - space: O(1) amortized
//
// SCORE: 25
func (d *DynamicArray[T]) Append(value T) {
	// hint: 1) if size >= capacity, call d.grow()
	//       2) set value at index = size
	//       3) increment size
	if d.Size() >= d.Cap() {
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
	//
	// https://www.josestg.com/posts/math/how-can-be-adding-a-new-item-to-a-dynamic-array-achieved-in-constant-time/
	const threshold = 256
	capacity := d.Cap()
	newCapacity := 2 * capacity
	if capacity >= threshold {
		newCapacity = capacity + (capacity+3*threshold)/4
	}
	newBackend := arrays.NewGarbageCollected[T](newCapacity, true)
	for i, v := range d.backend.Enum {
		newBackend.Set(i, v)
	}
	d.backend = newBackend
}

// Pop removes and returns the last element.
//
//	Before Pop():              After Pop():
//	┌───┬───┬───┬───┬───┐      ┌───┬───┬───┬───┐
//	│ A │ B │ C │ D │ E │      │ A │ B │ C │ D │
//	└───┴───┴───┴───┴───┘      └───┴───┴───┴───┘
//
//	Pop() -> E
//
// complexity:
//   - time : O(1)
//   - space: O(1)
//
// Panics if the array is empty.
func (d *DynamicArray[T]) Pop() T {
	v, ok := d.TryPop()
	if !ok {
		panic("DynamicArray.Pop: array is empty")
	}
	return v
}

// TryPop attempts to remove and return the last element.
// Returns (value, true) on success, or (zero, false) if empty.
//
// complexity:
//   - time : O(1)
//   - space: O(1)
func (d *DynamicArray[T]) TryPop() (T, bool) {
	var zero T
	if d.Size() == 0 {
		return zero, false
	}
	val := d.backend.Get(d.size - 1)
	d.backend.Set(d.size-1, zero)
	d.size--
	return val, true
}

// Clip reduces capacity to match size.
//
//	Before Clip():             After Clip():
//	┌───┬───┬───┬───┬───┬───┐  ┌───┬───┬───┐
//	│ A │ B │ C │   │   │   │  │ A │ B │ C │
//	└───┴───┴───┴───┴───┴───┘  └───┴───┴───┘
//	     size=3, cap=6              size=3, cap=3
//
// complexity:
//   - time : O(n)
//   - space: O(n)
//
// Panics if the array is empty.
func (d *DynamicArray[T]) Clip() {
	if d.Empty() {
		panic("DynamicArray.Clip: array is empty")
	}

	if d.size == d.Cap() {
		return
	}

	newBackend := arrays.NewGarbageCollected[T](d.size, true)
	for i := range d.Size() {
		newBackend.Set(i, d.backend.Get(i))
	}
	d.backend = newBackend
}

// Iter iterates over elements from front to back.
//
//	┌───┬───┬───┬───┬───┐
//	│ A │ B │ C │ D │ E │
//	└───┴───┴───┴───┴───┘
//	  ↓   ↓   ↓   ↓   ↓
//	  1   2   3   4   5   ← iteration order
//
// complexity:
//   - time : O(n)
//   - space: O(1)
func (d *DynamicArray[T]) Iter(yield func(T) bool) {
	for i := range d.Size() {
		if !yield(d.Get(i)) {
			break
		}
	}
}

// IterBackward iterates over elements from back to front.
//
//	┌───┬───┬───┬───┬───┐
//	│ A │ B │ C │ D │ E │
//	└───┴───┴───┴───┴───┘
//	  ↓   ↓   ↓   ↓   ↓
//	  5   4   3   2   1   ← iteration order
//
// complexity:
//   - time : O(n)
//   - space: O(1)
//
// SCORE: 5
func (d *DynamicArray[T]) IterBackward(yield func(T) bool) {
	// hint: loop from (size-1) down to 0, call yield(Get(i))
	for i := d.size - 1; i >= 0; i-- {
		val := d.Get(i)
		if !yield(val) {
			break
		}
	}
}

// Enum iterates over elements with their indices from front to back.
//
//	┌───┬───┬───┬───┬───┐
//	│ A │ B │ C │ D │ E │
//	└───┴───┴───┴───┴───┘
//	  0   1   2   3   4    <- indices
//	  ↓   ↓   ↓   ↓   ↓
//	 1st 2nd 3rd 4th 5th   <- enumeration order
//
// Example:
//
//	for index, value := range arr.Enum {
//	    fmt.Printf("%d: %v\n", index, value)
//	}
//
// complexity:
//   - time : O(n)
//   - space: O(1)
func (d *DynamicArray[T]) Enum(yield func(int, T) bool) {
	for i := range d.Size() {
		if !yield(i, d.Get(i)) {
			break
		}
	}
}

// EnumBackward iterates over elements with their indices from back to front.
//
//	┌───┬───┬───┬───┬───┐
//	│ A │ B │ C │ D │ E │
//	└───┴───┴───┴───┴───┘
//	  0   1   2   3   4    <- indices
//	  ↓   ↓   ↓   ↓   ↓
//	 5th 4th 3rd 2nd 1st   <- enumeration order
//
// Example:
//
//	for index, value := range arr.EnumBackward {
//	    fmt.Printf("%d: %v\n", index, value)  // prints 4:E, 3:D, 2:C, 1:B, 0:A
//	}
//
// complexity:
//   - time : O(n)
//   - space: O(1)
//
// SCORE: 5
func (d *DynamicArray[T]) EnumBackward(yield func(int, T) bool) {
	// hint: loop from (size-1) down to 0, call yield(i, Get(i))
	for i := d.size - 1; i >= 0; i-- {
		val := d.Get(i)
		if !yield(i, val) {
			break
		}
	}
}

// String returns the string representation.
//
//	┌───┬───┬───┬───┬───┐
//	│ 1 │ 2 │ 3 │ 4 │ 5 │
//	└───┴───┴───┴───┴───┘
//
//	String() -> "[1 2 3 4 5]"
//
// complexity:
//   - time : O(n)
//   - space: O(n)
func (d *DynamicArray[T]) String() string {
	return sequence.String(d.Iter)
}

// Insert adds an element at the given index.
//
//	Before Insert(2, X):       After Insert(2, X):
//	┌───┬───┬───┬───┬───┐      ┌───┬───┬───┬───┬───┬───┐
//	│ A │ B │ C │ D │ E │      │ A │ B │ X │ C │ D │ E │
//	└───┴───┴───┴───┴───┘      └───┴───┴───┴───┴───┴───┘
//	  0   1   2   3   4          0   1   2   3   4   5
//	          ↑                          ↑
//	    insert position            new element
//
// complexity:
//   - time : O(n) - must shift elements
//   - space: O(1)
//
// Panics if index < 0 or index > Size().
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

// Remove deletes and returns the element at the given index.
//
//	Before Remove(2):          After Remove(2):
//	┌───┬───┬───┬───┬───┐      ┌───┬───┬───┬───┐
//	│ A │ B │ C │ D │ E │      │ A │ B │ D │ E │
//	└───┴───┴───┴───┴───┘      └───┴───┴───┴───┘
//	  0   1   2   3   4          0   1   2   3
//	          ↑
//	       removed
//
//	Remove(2) -> C
//
// complexity:
//   - time : O(n) - must shift elements
//   - space: O(1)
//
// Panics if index < 0 or index >= Size().
func (d *DynamicArray[T]) Remove(index int) T {
	if v, ok := d.TryRemove(index); !ok {
		panic("dynamicarray: index out of range")
	} else {
		return v
	}
}

// TryRemove attempts to remove the element at the given index.
// Returns (value, true) on success, or (zero, false) if index is out of bounds.
//
// complexity:
//   - time : O(n) - shifts elements
//   - space: O(1)
//
// SCORE: 15
func (d *DynamicArray[T]) TryRemove(index int) (T, bool) {
	// hint: 1) check bounds, return (zero, false) if invalid
	//       2) save element at index
	//       3) shift elements left from index+1 to size-1
	//       4) clear last slot, decrement size
	//       5) return (saved, true)
	// Check bounds
	if index < 0 || index >= d.Size() {
		var arr T
		return arr, false
	}
	// Save element at index
	val := d.backend.Get(index)
	// Shift elements
	for i := index; i < d.size-1; i++ {
		d.backend.Set(i, d.backend.Get(i+1))
	}
	// decrement size
	d.size--
	return val, true
}

func (d *DynamicArray[T]) checkBounds(index int) {
	if index < 0 || index >= d.Size() {
		panic("dynamicarray: index out of range")
	}
}
