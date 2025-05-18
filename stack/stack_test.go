package stack_test

import (
	"math/rand"
	"testing"

	"github.com/josestg/dsa/adt/adttest"
	"github.com/josestg/dsa/stack"
)

func TestStack(t *testing.T) {
	c := stack.New[int]
	g := func() int {
		return rand.Intn(128)
	}

	simulator := adttest.StackSimulator(c, g)
	simulator.Run(t)
}
