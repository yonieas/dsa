#include "../array.h"
#include "greatest.h"

const size_t int_elem_size = sizeof(int);

TEST test_array_init_and_free(void) {
  Array arr;
  ASSERT_EQ(array_init(&arr, 10, int_elem_size), S_OK);
  ASSERT(arr.head != NULL);
  ASSERT_EQ(arr.length, 10);

  ASSERT_EQ(array_free(&arr), S_OK);
  ASSERT_EQ(arr.head, NULL);
  ASSERT_EQ(arr.length, 0);
  PASS();
}

TEST test_array_zero_values(void) {
  Array arr;
  ASSERT_EQ(array_init(&arr, 10, int_elem_size), S_OK);
  for (int i = 0; i < 10; i++) {
    int val = -1;
    ASSERT_EQ(array_get(&arr, i, &val, int_elem_size), S_OK);
    ASSERT_EQ(val, 0);
  }

  ASSERT_EQ(array_free(&arr), S_OK);
  ASSERT_EQ(arr.head, NULL);
  ASSERT_EQ(arr.length, 0);
  PASS();
}

TEST test_array_init_invalid_inputs(void) {
  ASSERT_EQ(array_init(NULL, 10, int_elem_size), S_ERR_SELF_IS_NULL);

  Array arr;
  ASSERT_EQ(array_init(&arr, 0, int_elem_size), S_ERR_OUT_OF_RANGE);
  PASS();
}

TEST test_array_set_and_get(void) {
  Array arr;
  array_init(&arr, 5, int_elem_size);

  int val = 42;
  ASSERT_EQ(array_set(&arr, 2, &val, int_elem_size), S_OK);

  int value = 0;
  ASSERT_EQ(array_get(&arr, 2, &value, int_elem_size), S_OK);
  ASSERT_EQ(value, val);

  array_free(&arr);
  PASS();
}

TEST test_array_bounds_check(void) {
  Array arr;
  array_init(&arr, 3, int_elem_size);

  ASSERT_EQ(array_set(&arr, 3, (void *)1, int_elem_size), S_ERR_OUT_OF_RANGE);
  ASSERT_EQ(array_get(&arr, 3, NULL, int_elem_size),
            S_ERR_RETURN_PARAMS_IS_NULL);
  ASSERT_EQ(array_get(NULL, 1, NULL, int_elem_size), S_ERR_SELF_IS_NULL);

  array_free(&arr);
  PASS();
}

TEST test_array_len(void) {
  Array arr;
  array_init(&arr, 7, int_elem_size);

  size_t len = 0;
  ASSERT_EQ(array_len(&arr, &len), S_OK);
  ASSERT_EQ(len, 7);

  ASSERT_EQ(array_len(NULL, &len), S_ERR_SELF_IS_NULL);
  ASSERT_EQ(array_len(&arr, NULL), S_ERR_RETURN_PARAMS_IS_NULL);

  array_free(&arr);
  PASS();
}

TEST test_array_int_simulation(void) {
  Array arr;
  const size_t length = 100;

  printf("creating array of length %zu...\n", length);
  ASSERT_EQ(array_init(&arr, length, int_elem_size), S_OK);

  printf("filling array with i * 2...\n");
  for (size_t i = 0; i < length; i++) {
    int val = (int)i * 2;
    ASSERT_EQ(array_set(&arr, i, &val, int_elem_size), S_OK);
  }

  printf("verifying values...\n");
  for (size_t i = 0; i < length; i++) {
    int val = 0;
    ASSERT_EQ(array_get(&arr, i, &val, int_elem_size), S_OK);
    ASSERT_EQ(val, (i * 2));
  }

  printf("adding 1 to values at index 10 to 19...\n");
  for (size_t i = 10; i < 20; i++) {
    int val = 0;
    ASSERT_EQ(array_get(&arr, i, &val, int_elem_size), S_OK);
    val++;
    ASSERT_EQ(array_set(&arr, i, &val, int_elem_size), S_OK);
  }

  printf("confirming updated values at index 10 to 19...\n");
  for (size_t i = 10; i < 20; i++) {
    int val = 0;
    ASSERT_EQ(array_get(&arr, i, &val, int_elem_size), S_OK);
    ASSERT_EQ(val, (int)(i * 2 + 1));
  }

  printf("cleaning up...\n");
  array_free(&arr);

  printf("array simulation test passed.\n");
  PASS();
}

typedef struct {
  int x, y, z;
} Vec3;

TEST test_array_struct_simulation(void) {
  const Vec3 a = {1, 2, 3};
  const Vec3 b = {4, 5, 6};
  const Vec3 c = {7, 8, 9};

  Array arr;
  const size_t length = 3;
  array_init(&arr, length, sizeof(Vec3));

  ASSERT_EQ(array_set(&arr, 0, &a, sizeof(Vec3)), S_OK);
  ASSERT_EQ(array_set(&arr, 1, &b, sizeof(Vec3)), S_OK);
  ASSERT_EQ(array_set(&arr, 2, &c, sizeof(Vec3)), S_OK);

  Vec3 a2, b2, c2;
  ASSERT_EQ(array_get(&arr, 0, &a2, sizeof(Vec3)), S_OK);
  ASSERT_EQ(array_get(&arr, 1, &b2, sizeof(Vec3)), S_OK);
  ASSERT_EQ(array_get(&arr, 2, &c2, sizeof(Vec3)), S_OK);

  ASSERT_EQ(a.x, a2.x);
  ASSERT_EQ(a.y, a2.y);
  ASSERT_EQ(a.z, a2.z);

  ASSERT_EQ(b.x, b2.x);
  ASSERT_EQ(b.y, b2.y);
  ASSERT_EQ(b.z, b2.z);

  ASSERT_EQ(c.x, c2.x);
  ASSERT_EQ(c.y, c2.y);
  ASSERT_EQ(c.z, c2.z);

  array_free(&arr);
  PASS();
}

SUITE(array_tests) {
  RUN_TEST(test_array_init_and_free);
  RUN_TEST(test_array_init_invalid_inputs);
  RUN_TEST(test_array_set_and_get);
  RUN_TEST(test_array_bounds_check);
  RUN_TEST(test_array_len);
  RUN_TEST(test_array_zero_values);
}

SUITE(array_simulation) {
  RUN_TEST(test_array_int_simulation);
  RUN_TEST(test_array_struct_simulation);
}

GREATEST_MAIN_DEFS();

int main(const int argc, char **argv) {
  GREATEST_MAIN_BEGIN();
  RUN_SUITE(array_tests);
  RUN_SUITE(array_simulation);
  GREATEST_MAIN_END();
}
