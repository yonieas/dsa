package graph

import (
	"iter"

	"github.com/josestg/dsa/hashmap"
	"github.com/josestg/dsa/linkedlist"
	"github.com/josestg/dsa/sequence"
)

type Graph[V comparable] struct {
	directed  bool
	adjacency *hashmap.HashMap[V, *linkedlist.SinglyLinkedList[V]]
}

func New[V comparable](directed bool) *Graph[V] {
	return &Graph[V]{
		directed:  directed,
		adjacency: hashmap.New[V, *linkedlist.SinglyLinkedList[V]](),
	}
}

func (g *Graph[V]) ensureNode(v V) *linkedlist.SinglyLinkedList[V] {
	neighbors, ok := g.adjacency.Get(v)
	if !ok {
		neighbors = linkedlist.NewSinglyLinkedList[V]()
		g.adjacency.Put(v, neighbors)
	}
	return neighbors
}

func (g *Graph[V]) Size() int {
	return g.adjacency.Size()
}

func (g *Graph[V]) Empty() bool {
	return g.adjacency.Empty()
}

func (g *Graph[V]) AddEdge(from, to V) {
	list := g.ensureNode(from)
	for v := range list.Iter {
		if v == to {
			return
		}
	}
	list.Append(to)
	if !g.directed {
		rev := g.ensureNode(to)
		for v := range rev.Iter {
			if v == from {
				return
			}
		}
		rev.Append(from)
	}
}

func (g *Graph[V]) DelEdge(from, to V) {
	if list, ok := g.adjacency.Get(from); ok {
		for i, v := range sequence.Enum(list.Iter) {
			if v == to {
				_ = list.Remove(i)
				break
			}
		}
	}
	if !g.directed {
		if list, ok := g.adjacency.Get(to); ok {
			for i, v := range sequence.Enum(list.Iter) {
				if v == from {
					_ = list.Remove(i)
					break
				}
			}
		}
	}
}

func (g *Graph[V]) HasEdge(from, to V) bool {
	if list, ok := g.adjacency.Get(from); ok {
		for v := range list.Iter {
			if v == to {
				return true
			}
		}
	}
	return false
}

func (g *Graph[V]) HasVertex(v V) bool {
	return g.adjacency.Exists(v)
}

func (g *Graph[V]) Vertex(yield func(V) bool) {
	for v := range g.adjacency.Iter {
		if !yield(v.Key()) {
			break
		}
	}
}

func (g *Graph[V]) Neighbors(v V) iter.Seq[V] {
	return func(yield func(V) bool) {
		if list, ok := g.adjacency.Get(v); ok {
			for val := range list.Iter {
				if !yield(val) {
					break
				}
			}
		}
	}
}
