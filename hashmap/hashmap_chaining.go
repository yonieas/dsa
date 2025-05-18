package hashmap

import (
	"cmp"
	"fmt"
	"strings"

	"github.com/josestg/dsa/internal/generics"
	"github.com/josestg/dsa/linkedlist"
	"github.com/josestg/dsa/sequence"
)

type HashMap[K comparable, V any] struct {
	size          int
	buckets       []*linkedlist.SinglyLinkedList[*Entry[K, V]]
	hashFunction  func(K) int
	loadThreshold float64
}

type Options[K comparable] struct {
	Capacity      int
	LoadThreshold float64
	HashFunction  func(key K) int
}

func New[K comparable, V any]() *HashMap[K, V] {
	return NewWith[K, V](Options[K]{
		Capacity:      DefaultCapacity,
		LoadThreshold: DefaultLoadThreshold,
		HashFunction:  DefaultHashFunction[K],
	})
}

func NewWith[K comparable, V any](opts Options[K]) *HashMap[K, V] {
	capacity := cmp.Or(opts.Capacity, DefaultCapacity)
	loadFactor := cmp.Or(opts.LoadThreshold, DefaultLoadThreshold)
	if loadFactor <= 0 || loadFactor >= 1 {
		panic("HashMap.NewWith: load factor must be in range (0,1) exclusive")
	}
	hashFunction := opts.HashFunction
	if hashFunction == nil {
		hashFunction = DefaultHashFunction
	}

	buckets := make([]*linkedlist.SinglyLinkedList[*Entry[K, V]], capacity)
	for i := range buckets {
		buckets[i] = linkedlist.NewSinglyLinkedList[*Entry[K, V]]()
	}

	return &HashMap[K, V]{
		size:          0,
		buckets:       buckets,
		loadThreshold: loadFactor,
		hashFunction:  hashFunction,
	}
}

func (h *HashMap[K, V]) Put(key K, value V) {
	if h.LoadFactor() >= h.loadThreshold {
		h.growAndRehash()
	}
	h.put(key, value)
}

func (h *HashMap[K, V]) Del(key K) {
	index := h.bucketIndex(key)
	entries := h.buckets[index]
	for i, v := range sequence.Enum(entries.Iter) {
		if v.key == key {
			_ = entries.Remove(i)
			h.size--
			break
		}
	}
}

func (h *HashMap[K, V]) put(key K, value V) {
	index := h.bucketIndex(key)
	entries := h.buckets[index]

	for v := range entries.Iter {
		if v.key == key {
			v.val = value
			return
		}
	}

	entries.Append(NewEntry(key, value))
	h.size++
}

func (h *HashMap[K, V]) Iter(yield func(*Entry[K, V]) bool) {
	for _, entries := range h.buckets {
		for v := range entries.Iter {
			if !yield(v) {
				return
			}
		}
	}
}

func (h *HashMap[K, V]) Keys(yield func(K) bool) {
	for e := range h.Iter {
		if !yield(e.Key()) {
			break
		}
	}
}

func (h *HashMap[K, V]) Get(key K) (V, bool) {
	if !h.Empty() {
		index := h.bucketIndex(key)
		entries := h.buckets[index]
		for v := range entries.Iter {
			if v.key == key {
				return v.val, true
			}
		}
	}
	return generics.ZeroValue[V](), false
}

func (h *HashMap[K, V]) Exists(key K) bool {
	_, found := h.Get(key)
	return found
}

func (h *HashMap[K, V]) String() string {
	var buf strings.Builder
	buf.WriteRune('[')
	for i, e := range sequence.Enum(h.Iter) {
		if i > 0 {
			buf.WriteRune(' ')
		}
		_, _ = fmt.Fprintf(&buf, "%v:%v", e.Key(), e.Value())
	}
	buf.WriteRune(']')
	return buf.String()
}

func (h *HashMap[K, V]) bucketIndex(key K) int {
	hash := h.hashFunction(key)
	if hash < 0 {
		hash = -hash
	}
	return hash % len(h.buckets)
}

func (h *HashMap[K, V]) growAndRehash() {
	h2 := NewWith[K, V](Options[K]{
		Capacity:      2 * len(h.buckets),
		LoadThreshold: h.loadThreshold,
		HashFunction:  h.hashFunction,
	})

	for e := range h.Iter {
		h2.put(e.Key(), e.Value())
	}

	h.size = h2.size
	h.buckets = h2.buckets
	h2 = nil
}

func (h *HashMap[K, V]) Size() int {
	return h.size
}

func (h *HashMap[K, V]) Capacity() int {
	return len(h.buckets)
}

func (h *HashMap[K, V]) LoadFactor() float64 {
	return float64(h.Size()) / float64(h.Capacity())
}

func (h *HashMap[K, V]) Empty() bool {
	return h.Size() == 0
}
