package bitsets_test

import (
	"testing"

	"github.com/josestg/dsa/adt/prop"
	"github.com/josestg/dsa/bitsets"
)

func TestBitSet(t *testing.T) {
	newBitSet := func() *bitsets.BitSet { return bitsets.New(64) }

	specs := []prop.Spec{
		prop.Set(newBitSet),
		prop.Union(newBitSet),
		prop.Intersection(newBitSet),
		prop.Disjoint(newBitSet),
		prop.BitAddExistsDel(bitsets.New),
		prop.BitToggle(bitsets.New),
		prop.BitReset(bitsets.New),
		prop.BitCount(bitsets.New),
		prop.BitLen(bitsets.New),
		prop.BitBounds(bitsets.New),
		prop.BitString(bitsets.New),
	}

	for _, spec := range specs {
		t.Run(spec.Name, spec.Test)
	}
}
