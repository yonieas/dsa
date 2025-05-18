package graph

import (
	"github.com/josestg/dsa/queue"
	"github.com/josestg/dsa/sets"
)

type WalkAlgorithm int

const (
	BFS WalkAlgorithm = iota
	DFSPreOrder
	DFSPostOrder
)

type Walker[T comparable] struct {
	alg     WalkAlgorithm
	graph   *Graph[T]
	visited *sets.HashSet[T]
}

func NewWalker[T comparable](g *Graph[T], alg WalkAlgorithm) *Walker[T] {
	return &Walker[T]{
		alg:     alg,
		graph:   g,
		visited: sets.New[T](),
	}
}

func (w *Walker[T]) Visited(node T) bool {
	return w.visited.Exists(node)
}

func (w *Walker[T]) Explored() bool {
	for n := range w.graph.Vertex {
		if !w.Visited(n) {
			return false
		}
	}
	return true
}

func (w *Walker[T]) WalkAll(visit func(T)) {
	for n := range w.graph.Vertex {
		if !w.Visited(n) {
			w.Walk(n, visit)
		}
	}
}

func (w *Walker[T]) Walk(start T, visit func(T)) {
	switch w.alg {
	case BFS:
		w.bfs(start, visit)
	case DFSPreOrder, DFSPostOrder:
		w.dfs(start, visit)
	default:
		panic("unknown walk algorithm")
	}
}

func (w *Walker[T]) dfs(start T, visit func(T)) {
	var traverse func(T)
	traverse = func(n T) {
		if w.Visited(n) {
			return
		}
		if w.alg == DFSPreOrder {
			visit(n)
		}
		w.visited.Add(n)
		for adj := range w.graph.Neighbors(n) {
			traverse(adj)
		}
		if w.alg == DFSPostOrder {
			visit(n)
		}
	}
	traverse(start)
}

func (w *Walker[T]) bfs(start T, visit func(T)) {
	q := queue.New[T]()
	q.Enqueue(start)
	w.visited.Add(start)
	for !q.Empty() {
		node := q.Dequeue()
		visit(node)
		for neighbor := range w.graph.Neighbors(node) {
			if !w.Visited(neighbor) {
				w.visited.Add(neighbor)
				q.Enqueue(neighbor)
			}
		}
	}
}
