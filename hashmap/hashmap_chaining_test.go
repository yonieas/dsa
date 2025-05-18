package hashmap_test

import (
	"math/rand"
	"testing"

	"github.com/josestg/dsa/adt/adttest"
	"github.com/josestg/dsa/hashmap"
)

func TestHashMap(t *testing.T) {
	c := hashmap.New[int, int]
	kg := func() int {
		return rand.Intn(128)
	}
	vg := func() int {
		return rand.Intn(128)
	}

	tests := []struct {
		name      string
		simulator adttest.Runner
	}{
		{name: "hashmap", simulator: adttest.HashMapSimulator(c, kg, vg)},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.simulator)
	}
}
