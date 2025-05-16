#ifndef FONDASI_DEV_ARRAYLIB_ARRAY_H
#define FONDASI_DEV_ARRAYLIB_ARRAY_H

#include <stdlib.h>

/// A basic, fixed-size array of integers.
///
/// This array is manually managed. Memory is allocated on initialization
/// and must be released with `array_free`. All operations return a status
/// to indicate success or failure.
typedef struct {
  void *head;
  size_t length;
  size_t elem_size;
} Array;

/// Status values returned by all array operations.
typedef enum {
  S_OK = 0,                         ///< Operation completed successfully.
  S_ERR_SELF_IS_NULL = 1,           ///< Provided `self` pointer is NULL.
  S_ERR_RETURN_PARAMS_IS_NULL = 2,  ///< Output pointer is NULL.
  S_ERR_OUT_OF_MEMORY = 4,          ///< Memory allocation failed.
  S_ERR_OUT_OF_RANGE = 5,           ///< Index is outside array bounds.
  S_ERR_ELEMENT_SIZE_MISMATCH =6    ///< Element size is different.
} status_t;

/// Initializes a new array of the given length.
///
/// This allocates memory for `length` integers and sets all values to 0.
/// The array must be freed later with `array_free`.
///
/// \param self       Pointer to the Array to initialize.
/// \param length     Number of elements to allocate.
/// \param elem_size  Size of an element.
///
/// \retval S_OK                 Success.
/// \retval S_ERR_SELF_IS_NULL  `self` is NULL.
/// \retval S_ERR_OUT_OF_RANGE  `length` is 0 or too large to allocate.
/// \retval S_ERR_OUT_OF_MEMORY Allocation failed.
status_t array_init(Array *self, size_t length, size_t elem_size);

/// Frees any memory associated with the array and resets its state.
///
/// Safe to call multiple times. After this call, the array is empty.
///
/// \param self Pointer to the Array to clean up.
///
/// \retval S_OK                Success.
/// \retval S_ERR_SELF_IS_NULL `self` is NULL.
status_t array_free(Array *self);

/// Writes a value into the array at the specified index.
///
/// \param self       Pointer to the Array.
/// \param index      Index to write to (0-based).
/// \param value_in   Value to store.
/// \param value_size  Size of the value to store.
///
/// \retval S_OK                Success.
/// \retval S_ERR_SELF_IS_NULL `self` is NULL.
/// \retval S_ERR_OUT_OF_RANGE Index is invalid.
status_t array_set(const Array *self, size_t index, const void* value_in, size_t value_size);

/// Reads a value from the array at the given index.
///
/// \param self         Pointer to the Array.
/// \param index        Index to read from.
/// \param value_out    Pointer to store the result.
/// \param value_size   Size of the value to read.
///
/// \retval S_OK                         Success.
/// \retval S_ERR_SELF_IS_NULL           `self` is NULL.
/// \retval S_ERR_RETURN_PARAMS_IS_NULL `value_out` is NULL.
/// \retval S_ERR_OUT_OF_RANGE           Index is invalid.
status_t array_get(const Array *self, size_t index, void *value_out, size_t value_size);

/// Returns the current length of the array.
///
/// \param self   Pointer to the Array.
/// \param length Pointer to store the result.
///
/// \retval S_OK                         Success.
/// \retval S_ERR_SELF_IS_NULL           `self` is NULL.
/// \retval S_ERR_RETURN_PARAMS_IS_NULL `length` is NULL.
status_t array_len(const Array *self, size_t *length);

#endif // FONDASI_DEV_ARRAYLIB_ARRAY_H
