package sets_test

import (
	"testing"

	"github.com/josestg/dsa/adt/prop"
	"github.com/josestg/dsa/sets"
)

func TestHashSet(t *testing.T) {
	specs := []prop.Spec{
		prop.Set(sets.New[int]),
		prop.Union(sets.New[int]),
		prop.Intersection(sets.New[int]),
		prop.Disjoint(sets.New[int]),
	}

	for _, spec := range specs {
		t.Run(spec.Name, spec.Test)
	}
}
