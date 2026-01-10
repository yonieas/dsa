package stack_test

import (
	"testing"

	"github.com/josestg/dsa/adt/prop"
	"github.com/josestg/dsa/dynamicarray"
	"github.com/josestg/dsa/stack"
)

func TestStack(t *testing.T) {
	specs := []prop.Spec{
		prop.Stack(stack.New[int]),
		prop.TryPeekStack(stack.New[int]),
		prop.TryPopStack(stack.New[int]),
		prop.IterStack(stack.New[int]),
	}

	for _, spec := range specs {
		t.Run(spec.Name, spec.Test)
	}
}

func TestStackWithArrayBackend(t *testing.T) {
	f := func() *stack.Stack[int] {
		return stack.NewWith(dynamicarray.New[int](4))
	}
	specs := []prop.Spec{
		prop.Stack(f),
		prop.TryPeekStack(f),
		prop.TryPopStack(f),
		prop.IterStack(f),
	}

	for _, spec := range specs {
		t.Run(spec.Name, spec.Test)
	}
}
