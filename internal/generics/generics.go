// Package generics provides utility functions for working with Go generics.
// These helpers make it easier to write generic code that needs common operations
// like obtaining zero values.
package generics

// ZeroValue returns the zero value for any type T.
// This is useful in generic code where you need to return a "null" or default value.
//
// The zero value depends on the type:
//   - Numeric types (int, float64, etc.): 0
//   - string: ""
//   - bool: false
//   - Pointers, slices, maps, channels, functions: nil
//   - Structs: All fields set to their zero values
//
// Time Complexity: O(1)
//
// Example usage:
//
//	func Pop[T any](slice []T) (T, bool) {
//	    if len(slice) == 0 {
//	        return generics.ZeroValue[T](), false
//	    }
//	    return slice[len(slice)-1], true
//	}
//
// This is particularly useful when implementing data structures that need
// to return a value even when the structure is empty (paired with a boolean).
func ZeroValue[T any]() (zeroValue T) {
	return zeroValue
}
