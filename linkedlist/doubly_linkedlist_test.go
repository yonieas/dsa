package linkedlist_test

import (
	"math/rand"
	"testing"

	"github.com/josestg/dsa/adt/adttest"
	"github.com/josestg/dsa/linkedlist"
)

func TestDoublyLinkedList(t *testing.T) {
	c := linkedlist.NewDoublyLinkedList[int]
	g := func() int {
		return rand.Intn(128)
	}

	tests := []struct {
		name      string
		simulator adttest.Runner
	}{
		{name: "append", simulator: adttest.AppendSimulator(c, g)},
		{name: "prepend", simulator: adttest.PrependSimulator(c, g)},
		{name: "pop", simulator: adttest.PopSimulator(c, g)},
		{name: "shift", simulator: adttest.ShiftSimulator(c, g)},
		{name: "get", simulator: adttest.GetSimulator(c, g)},
		{name: "set", simulator: adttest.SetSimulator(c, g)},
		{name: "iter", simulator: adttest.IterSimulator(c, g)},
		{name: "backward iter", simulator: adttest.IterBackwardSimulator(c, g)},
		{name: "to string", simulator: adttest.BracketStringSimulator(c, g)},
		{name: "sort", simulator: adttest.SortSimulator(c, g)},
		{name: "insert and remove", simulator: adttest.InsertRemoveSimulator(c, g)},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.simulator)
	}
}
