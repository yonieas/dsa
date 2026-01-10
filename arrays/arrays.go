// Package arrays provides a fixed-size array backed by C memory.
//
// # What is an Array?
//
// An array is the simplest and most fundamental data structure. It stores
// elements in contiguous memory locations, like boxes lined up in a row.
// The magic of arrays is that you can jump directly to any element using
// its index, no need to traverse through other elements first.
//
// Think of it like a row of mailboxes in an apartment building. If you know
// the apartment number (index), you can go directly to that mailbox without
// checking others. This is why array access is O(1), constant time,
// regardless of how many elements exist.
//
// # Why Use This Implementation?
//
// This package uses CGO to interface with C, which is useful for learning
// how arrays work at the memory level, interfacing with C libraries that
// expect raw memory, and understanding memory management before Go's GC
// hides it from you.
//
// # When to Use Arrays
//
// Use arrays when you know the size upfront and it won't change, you need
// fast random access by index, memory efficiency matters (no pointer
// overhead), or you're building other data structures like hash tables.
//
// Avoid arrays when size is unknown or changes frequently (use dynamic
// array/slice instead), or you need frequent insertions/deletions in
// the middle.
//
// # Complexity
//
//	Access by index: O(1)
//	Search:          O(n)
//	Insert/Delete:   O(n)
//
// # Further Reading
//
// CLRS "Introduction to Algorithms", Chapter 10.
// Sedgewick "Algorithms", Section 1.1.
// https://en.wikipedia.org/wiki/Array_data_structure
package arrays

// #cgo CFLAGS: -I${SRCDIR}
// #include "array.h"
import "C"
import (
    "fmt"
    "reflect"
    "runtime"
    "strings"
    "unsafe"
)

// Array is a fixed-size generic array backed by C memory.
//
//	┌───┬───┬───┬───┬───┐
//	│ 0 │ 1 │ 2 │ 3 │ 4 │  <- fixed length, cannot grow
//	└───┴───┴───┴───┴───┘
//	        C memory
type Array[T any] struct {
    backend   C.Array
    gcEnabled bool
    cleanup   runtime.Cleanup
}

// New creates a fixed-size array with the given length.
// All elements are initialized to their zero value.
//
//	New[int](5) creates:
//	┌───┬───┬───┬───┬───┐
//	│ 0 │ 0 │ 0 │ 0 │ 0 │
//	└───┴───┴───┴───┴───┘
//	  0   1   2   3   4
//
// IMPORTANT: You must call Free() when done to prevent memory leaks.
//
// complexity:
//   - time : O(length)
//   - space: O(length)
//
// Panics if:
//   - T is a zero-sized type (like struct{})
//   - T is a pointer type
func New[T any](length int) *Array[T] {
    return NewGarbageCollected[T](length, false)
}

// NewGarbageCollected creates an array with optional automatic cleanup.
//
// If enabled=true:
//   - Go's finalizer calls Free() when array becomes unreachable
//   - No need to manually call Free()
//   - Cleanup may be delayed (depends on GC)
//
// If enabled=false:
//   - You MUST call Free() manually to prevent memory leaks
//   - Use for performance-critical code with controlled lifecycle
func NewGarbageCollected[T any](length int, enabled bool) *Array[T] {
    var zero T
    elemSize := unsafe.Sizeof(zero)
    if elemSize == 0 {
        panic("arrays: zero-sized types are not supported")
    }

    t := reflect.TypeOf(zero)
    if t.Kind() == reflect.Ptr {
        panic("arrays: pointer types are not supported as array elements")
    }

    var backend C.Array
    s := C.array_init(
        &backend,
        C.size_t(length),
        C.size_t(elemSize),
    )
    mustOk(s)
    a := &Array[T]{
        backend:   backend,
        gcEnabled: enabled,
    }
    if a.gcEnabled {
        a.cleanup = runtime.AddCleanup(a, func(backend C.Array) {
            b := backend
            C.array_free(&b)
        }, a.backend)
    }
    return a
}

// Free releases the C memory backing this array.
//
//		arr := arrays.New[int](100)
//	 defer arr.Free() // Release memory
//		// Don't use arr after this!
//
// Safe to call multiple times (idempotent).
//
// complexity:
//   - time : O(1)
//   - space: O(1)
func (a *Array[T]) Free() {
    if a.backend.head != nil {
        s := C.array_free(&a.backend)
        mustOk(s)
        if a.gcEnabled {
            a.cleanup.Stop()
        }
    }
}

// Len returns the number of elements in the array.
//
//	┌───┬───┬───┬───┬───┐
//	│ A │ B │ C │ D │ E │
//	└───┴───┴───┴───┴───┘
//
//	Len() -> 5
//
// complexity:
//   - time : O(1)
//   - space: O(1)
func (a *Array[T]) Len() int {
    var length C.size_t
    mustOk(C.array_len(&a.backend, &length))
    return int(length)
}

// Size returns the number of elements in the array.
// This is an alias for Len() to satisfy the adt.Sizer interface.
//
// complexity:
//   - time : O(1)
//   - space: O(1)
func (a *Array[T]) Size() int {
    return a.Len()
}

// Empty returns true if the array has no elements.
//
// complexity:
//   - time : O(1)
//   - space: O(1)
func (a *Array[T]) Empty() bool {
    return a.Len() == 0
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
func (a *Array[T]) Head() T {
    if v, ok := a.TryHead(); !ok {
        panic("Array.Head: array is empty")
    } else {
        return v
    }
}

// TryHead attempts to return the first element.
// Returns (value, true) on success, or (zero, false) if empty.
//
// complexity:
//   - time : O(1)
//   - space: O(1)
func (a *Array[T]) TryHead() (T, bool) {
    if a.Len() == 0 {
        var zero T
        return zero, false
    }
    return a.Get(0), true
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
func (a *Array[T]) Tail() T {
    if v, ok := a.TryTail(); !ok {
        panic("Array.Tail: array is empty")
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
func (a *Array[T]) TryTail() (T, bool) {
    if a.Len() == 0 {
        var zero T
        return zero, false
    }
    return a.Get(a.Len() - 1), true
}

// Set updates the element at the given index.
//
//	Before Set(2, X):          After Set(2, X):
//	┌───┬───┬───┬───┬───┐      ┌───┬───┬───┬───┬───┐
//	│ A │ B │ C │ D │ E │      │ A │ B │ X │ D │ E │
//	└───┴───┴───┴───┴───┘      └───┴───┴───┴───┴───┘
//	  0   1   2   3   4          0   1   2   3   4
//	          ↑                          ↑
//	       changed                   updated
//
// complexity:
//   - time : O(1) - direct memory access
//   - space: O(1)
//
// Panics if index < 0 or index >= Len().
func (a *Array[T]) Set(index int, value T) {
    a.boundCheck(index)
    a.setUnchecked(index, value)
}

func (a *Array[T]) setUnchecked(index int, value T) {
    s := C.array_set(
        &a.backend,
        C.size_t(index),
        unsafe.Pointer(&value),
        C.size_t(unsafe.Sizeof(value)),
    )
    mustOk(s)
}

// TrySet attempts to update the element at the given index.
// Returns true on success, false if index is out of bounds.
//
// complexity:
//   - time : O(1)
//   - space: O(1)
func (a *Array[T]) TrySet(index int, value T) bool {
    if index < 0 || index >= a.Len() {
        return false
    }
    a.setUnchecked(index, value)
    return true
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
// This is the key advantage of arrays: O(1) random access.
//
// complexity:
//   - time : O(1) - direct memory access
//   - space: O(1)
//
// Panics if index < 0 or index >= Len().
func (a *Array[T]) Get(index int) T {
    a.boundCheck(index)
    return a.getUnchecked(index)
}

func (a *Array[T]) getUnchecked(index int) T {
    var out T
    s := C.array_get(
        &a.backend,
        C.size_t(index),
        unsafe.Pointer(&out),
        C.size_t(unsafe.Sizeof(out)),
    )
    mustOk(s)
    return out
}

// TryGet attempts to retrieve the element at the given index.
// Returns (value, true) on success, or (zero, false) if index is out of bounds.
//
// complexity:
//   - time : O(1)
//   - space: O(1)
func (a *Array[T]) TryGet(index int) (T, bool) {
    if index < 0 || index >= a.Len() {
        var zero T
        return zero, false
    }
    return a.getUnchecked(index), true
}

// Iter iterates over all elements from front to back.
//
//	┌───┬───┬───┬───┬───┐
//	│ A │ B │ C │ D │ E │
//	└───┴───┴───┴───┴───┘
//	  ↓   ↓   ↓   ↓   ↓
//	 1st 2nd 3rd 4th 5th   <- iteration order
//
// Example:
//
//	for value := range arr.Iter {
//	    fmt.Println(value)  // prints A, B, C, D, E
//	}
//
// complexity:
//   - time : O(n)
//   - space: O(1)
func (a *Array[T]) Iter(yield func(T) bool) {
    for i := range a.Len() {
        if !yield(a.Get(i)) {
            break
        }
    }
}

// IterBackward iterates over all elements from back to front.
//
//	┌───┬───┬───┬───┬───┐
//	│ A │ B │ C │ D │ E │
//	└───┴───┴───┴───┴───┘
//	  ↓   ↓   ↓   ↓   ↓
//	 5th 4th 3rd 2nd 1st   <- iteration order
//
// Example:
//
//	for value := range arr.IterBackward {
//	    fmt.Println(value)  // prints E, D, C, B, A
//	}
//
// complexity:
//   - time : O(n)
//   - space: O(1)
func (a *Array[T]) IterBackward(yield func(T) bool) {
    for i := a.Len() - 1; i >= 0; i-- {
        if !yield(a.Get(i)) {
            break
        }
    }
}

// Enum iterates over all elements with their indices from front to back.
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
//	    fmt.Printf("%d: %v\n", index, value)  // prints 0:A, 1:B, 2:C, 3:D, 4:E
//	}
//
// complexity:
//   - time : O(n)
//   - space: O(1)
func (a *Array[T]) Enum(yield func(int, T) bool) {
    for i := range a.Len() {
        if !yield(i, a.Get(i)) {
            break
        }
    }
}

// EnumBackward iterates over all elements with their indices from back to front.
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
func (a *Array[T]) EnumBackward(yield func(int, T) bool) {
    for i := a.Len() - 1; i >= 0; i-- {
        if !yield(i, a.Get(i)) {
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
func (a *Array[T]) String() string {
    var sb strings.Builder
    sb.WriteString("[")
    for i, v := range a.Enum {
        if i > 0 {
            sb.WriteRune(' ')
        }
        _, err := fmt.Fprintf(&sb, "%v", v)
        if err != nil {
            panic(fmt.Errorf("arrays: to string at index %d: %v", i, err))
        }
    }
    sb.WriteString("]")
    return sb.String()
}

func (a *Array[T]) boundCheck(index int) {
    n := a.Len()
    if index < 0 || index >= n {
        panic("index out of range")
    }
}

func mustOk(s C.status_t) {
    if err := errorOf(s); err != nil {
        panic(err)
    }
}

func errorOf(s C.status_t) error {
    switch s {
    default:
        return fmt.Errorf("arrays: status_t(%v): unrecognized status", s)
    case C.S_OK:
        return nil
    case C.S_ERR_SELF_IS_NULL:
        return fmt.Errorf("arrays: status_t(%v): self is null", s)
    case C.S_ERR_RETURN_PARAMS_IS_NULL:
        return fmt.Errorf("arrays: status_t(%v): out params is missing", s)
    case C.S_ERR_OUT_OF_MEMORY:
        return fmt.Errorf("arrays: status_t(%v): out of memory", s)
    case C.S_ERR_OUT_OF_RANGE:
        return fmt.Errorf("arrays: status_t(%v): index out of range", s)
    case C.S_ERR_ELEMENT_SIZE_MISMATCH:
        return fmt.Errorf("arrays: status_t(%v): type size mismatched", s)
    }
}
