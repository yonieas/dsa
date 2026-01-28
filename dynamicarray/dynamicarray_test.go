package dynamicarray_test

import (
	"testing"

	"github.com/josestg/dsa/adt/prop"
	"github.com/josestg/dsa/dynamicarray"
)

func TestDynamicArray(t *testing.T) {
	f := func() *dynamicarray.DynamicArray[int] {
		return dynamicarray.New[int](1)
	}
	specs := []prop.Spec{
		prop.Append(f),
		prop.Prepend(f),
		prop.GetSet(f),
		prop.HeadTail(f),
		prop.Pop(f),
		prop.Shift(f),
		prop.TryPop(f),
		prop.TryShift(f),
		prop.TryHead(f),
		prop.TryTail(f),
		prop.TryGet(f),
		prop.TrySet(f),
		prop.TryRemove(f),
		prop.Iter(f),
		prop.IterBackward(f),
		prop.EnumBackward(f),
		prop.Insert(f),
		prop.Remove(f),
		prop.Swap(f),
		prop.Cap(func(cap int) *dynamicarray.DynamicArray[int] {
			return dynamicarray.New[int](cap)
		}),
		prop.Clip(func(cap int) *dynamicarray.DynamicArray[int] {
			return dynamicarray.New[int](cap)
		}),
	}
	for _, s := range specs {
		t.Run(s.Name, s.Test)
	}
}
