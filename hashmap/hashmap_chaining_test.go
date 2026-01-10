package hashmap_test

import (
	"testing"

	"github.com/josestg/dsa/adt/prop"
	"github.com/josestg/dsa/hashmap"
)

func TestHashMap(t *testing.T) {
	specs := []prop.Spec{
		prop.MapPutGetDel(hashmap.New[int, int]),
		prop.MapKeys(hashmap.New[int, int]),
		prop.MapLoadFactor(hashmap.New[int, int]),
		prop.MapString(hashmap.New[int, int]),
	}

	for _, spec := range specs {
		t.Run(spec.Name, spec.Test)
	}
}
