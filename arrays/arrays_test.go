package arrays_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/josestg/dsa/arrays"
)

func TestSimulatorArray(t *testing.T) {
	type elem struct {
		b bool
		i int
		s string
	}

	var zero elem

	a := arrays.New[elem](5)
	t.Cleanup(a.Free)

	for _, v := range a.Iter(false) {
		if v != zero {
			t.Errorf("all elem must be zero-value")
		}
	}

	e := a.Get(0)
	e.s = "xyz"

	if a.Get(0) != zero {
		t.Error("changing the copy should not affects the original value")
	}

	a.Set(0, e)
	if a.Get(0) != e {
		t.Errorf("set should be replaced the value")
	}

	ga := [5]elem{}
	for i := range a.Len() {
		v := a.Get(i)

		v.i = 1 << i
		v.b = i%2 == 0
		v.s = strconv.Itoa(i)

		a.Set(i, v)
		ga[i] = v
	}

	as := a.String()
	gas := fmt.Sprint(ga)
	if as != gas {
		t.Errorf("string representation of array and Go's array should be equal. as: %q, gas: %q", as, gas)
	}

	for i, v := range a.Iter(true) {
		if v != ga[i] {
			t.Errorf("ga[%d] != a[%d]", i, i)
		}
	}
}
