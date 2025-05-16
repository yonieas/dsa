#include "array.h"
#include <limits.h>
#include <stdint.h>
#include <stdlib.h>

#include <string.h>

status_t array_init(Array *self, const size_t length, const size_t elem_size) {
  if (self == NULL)
    return S_ERR_SELF_IS_NULL;

  const size_t max_length = SIZE_MAX / elem_size;
  if (length == 0 || length > max_length) {
    return S_ERR_OUT_OF_RANGE;
  }

  self->head = calloc(length, elem_size);
  if (self->head == NULL) {
    return S_ERR_OUT_OF_MEMORY;
  }

  self->length = length;
  self->elem_size = elem_size;
  return S_OK;
}

status_t array_free(Array *self) {
  if (self == NULL)
    return S_ERR_SELF_IS_NULL;

  free(self->head);
  self->head = NULL;
  self->length = 0;
  self->elem_size = 0;

  return S_OK;
}

status_t array_set(const Array *self, const size_t index, const void *value_in,
                   const size_t value_size) {
  if (self == NULL)
    return S_ERR_SELF_IS_NULL;
  if (index >= self->length)
    return S_ERR_OUT_OF_RANGE;
  if (self->elem_size != value_size)
    return S_ERR_ELEMENT_SIZE_MISMATCH;

  memcpy((char *)self->head + index * self->elem_size, value_in, value_size);
  return S_OK;
}

status_t array_get(const Array *self, const size_t index, void *value_out,
                   const size_t value_size) {
  if (self == NULL)
    return S_ERR_SELF_IS_NULL;
  if (value_out == NULL)
    return S_ERR_RETURN_PARAMS_IS_NULL;
  if (index >= self->length)
    return S_ERR_OUT_OF_RANGE;
  if (self->elem_size != value_size)
    return S_ERR_ELEMENT_SIZE_MISMATCH;

  memcpy(value_out, (char *)self->head + index * self->elem_size, value_size);
  return S_OK;
}

status_t array_len(const Array *self, size_t *length) {
  if (self == NULL)
    return S_ERR_SELF_IS_NULL;
  if (length == NULL)
    return S_ERR_RETURN_PARAMS_IS_NULL;

  *length = self->length;
  return S_OK;
}