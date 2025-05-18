package sets_test

import (
	"math/rand"
	"testing"

	"github.com/josestg/dsa/adt/adttest"
	"github.com/josestg/dsa/sets"
)

func TestHashSet(t *testing.T) {
	c := sets.New[int]
	g := func() int {
		return rand.Intn(128)
	}

	tests := []struct {
		name      string
		simulator adttest.Runner
	}{
		{name: "sets", simulator: adttest.HashSetSimulator(c, g)},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.simulator)
	}
}
