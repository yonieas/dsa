package hashmap

import (
	"fmt"
	"hash/fnv"
)

const (
	DefaultCapacity      = 16
	DefaultLoadThreshold = 0.75
)

type Entry[K comparable, V any] struct {
	key K
	val V
}

func (e Entry[K, V]) Key() K   { return e.key }
func (e Entry[K, V]) Value() V { return e.val }

func NewEntry[K comparable, V any](key K, val V) *Entry[K, V] {
	return &Entry[K, V]{
		key: key,
		val: val,
	}
}

func DefaultHashFunction[K comparable](key K) int {
	h := fnv.New32()
	_, _ = fmt.Fprint(h, key)
	return int(h.Sum32())
}
