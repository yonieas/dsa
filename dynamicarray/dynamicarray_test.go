package dynamicarray_test

import (
	"math/rand"
	"testing"

	"github.com/josestg/dsa/adt/adttest"
	"github.com/josestg/dsa/dynamicarray"
	"github.com/stretchr/testify/assert"
)

func TestDynamicArray(t *testing.T) {
	c := func() *dynamicarray.DynamicArray[int] {
		return dynamicarray.New[int](1)
	}

	g := func() int {
		return rand.Intn(128)
	}

	d := func(a *dynamicarray.DynamicArray[int]) {
		a.Free()
	}

	tests := []struct {
		name      string
		simulator adttest.Runner
	}{
		{name: "append", simulator: adttest.AppendSimulator(c, g, d)},
		{name: "prepend", simulator: adttest.PrependSimulator(c, g, d)},
		{name: "pop", simulator: adttest.PopSimulator(c, g, d)},
		{name: "shift", simulator: adttest.ShiftSimulator(c, g, d)},
		{name: "get", simulator: adttest.GetSimulator(c, g, d)},
		{name: "set", simulator: adttest.SetSimulator(c, g, d)},
		{name: "iter", simulator: adttest.IterSimulator(c, g, d)},
		{name: "backward iter", simulator: adttest.IterBackwardSimulator(c, g, d)},
		{name: "to string", simulator: adttest.BracketStringSimulator(c, g, d)},
		{name: "sort", simulator: adttest.SortSimulator(c, g, d)},
		{name: "doubling grow", simulator: adttest.DoublingGrowSimulator(c, g, d)},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.simulator)
	}
}

func TestDynamicArray_Clip(t *testing.T) {
	t.Run("cap > len", func(t *testing.T) {
		a := dynamicarray.New[int](1)
		t.Cleanup(a.Free)

		n := 10 // the threshold - 1.
		for i := 0; i < n; i++ {
			v := 2*i + 1 // generating odd number.
			a.Append(v)
		}

		assert.NotEqual(t, a.Cap(), a.Size())
		a.Clip()
		assert.Equal(t, a.Cap(), a.Size())
	})

	t.Run("cap == len", func(t *testing.T) {
		a := dynamicarray.New[int](1)
		t.Cleanup(a.Free)

		n := 4 // the threshold - 1.
		for i := 0; i < n; i++ {
			v := 2*i + 1 // generating odd number.
			a.Append(v)
		}

		assert.Equal(t, a.Cap(), a.Size())
		a.Clip()
		assert.Equal(t, a.Cap(), a.Size())
	})
}
