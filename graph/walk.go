// Package graph provides graph traversal algorithms.
package graph

import (
	"github.com/josestg/dsa/queue"
	"github.com/josestg/dsa/sets"
)

// WalkAlgorithm specifies the graph traversal algorithm.
type WalkAlgorithm int

const (
	// BFS (Breadth-First Search) explores level by level.
	// Visits all neighbors before going deeper.
	//
	//	Starting from A:
	//
	//	    A           Level 0
	//	   / \
	//	  B   C         Level 1
	//	 / \   \
	//	D   E   F       Level 2
	//
	//	Visit order: A, B, C, D, E, F
	BFS WalkAlgorithm = iota

	// DFSPreOrder (Depth-First Search, Pre-Order) goes deep first.
	// Visits node BEFORE its descendants.
	//
	//	Starting from A:
	//
	//	    A  (1)
	//	   / \
	//	  B   C  (2, 5)
	//	 / \
	//	D   E   (3, 4)
	//
	//	Visit order: A, B, D, E, C
	DFSPreOrder

	// DFSPostOrder visits node AFTER its descendants.
	//
	//	Starting from A:
	//
	//	    A  (5)
	//	   / \
	//	  B   C  (3, 4)
	//	 / \
	//	D   E   (1, 2)
	//
	//	Visit order: D, E, B, C, A
	DFSPostOrder
)

// Walker traverses a graph systematically.
// It tracks visited vertices to avoid infinite loops in cyclic graphs.
type Walker[T comparable] struct {
	alg     WalkAlgorithm
	graph   *Graph[T]
	visited *sets.HashSet[T]
}

// NewWalker creates a walker for the given graph and algorithm.
//
//	g := graph.New[string](true)
//	g.AddEdge("A", "B")
//	g.AddEdge("A", "C")
//
//	walker := graph.NewWalker(g, graph.BFS)
//
// complexity:
//   - time : O(1)
//   - space: O(1)
func NewWalker[T comparable](g *Graph[T], alg WalkAlgorithm) *Walker[T] {
	return &Walker[T]{
		alg:     alg,
		graph:   g,
		visited: sets.New[T](),
	}
}

// Visited checks if a vertex has been visited during traversal.
//
// complexity:
//   - time : O(1) average
//   - space: O(1)
func (w *Walker[T]) Visited(node T) bool {
	return w.visited.Exists(node)
}

// Explored returns true if all vertices have been visited.
//
// complexity:
//   - time : O(V)
//   - space: O(1)
func (w *Walker[T]) Explored() bool {
	for n := range w.graph.Vertex {
		if !w.Visited(n) {
			return false
		}
	}
	return true
}

// WalkAll traverses the entire graph, including disconnected components.
//
//	Disconnected graph:
//
//	    A ──► B       C ──► D
//	   (component 1)  (component 2)
//
//	Walk(A, visit) only visits A, B
//	WalkAll(visit) visits A, B, C, D
//
// complexity:
//   - time : O(V + E)
//   - space: O(V) for visited set
func (w *Walker[T]) WalkAll(visit func(T)) {
	for n := range w.graph.Vertex {
		if !w.Visited(n) {
			w.Walk(n, visit)
		}
	}
}

// Walk traverses the graph starting from a specific vertex.
//
// BFS example:
//
//	    A
//	   / \
//	  B   C
//	 / \
//	D   E
//
//	walker := NewWalker(g, BFS)
//	walker.Walk("A", print)
//	// Output: A, B, C, D, E (level by level)
//
// DFS Pre-Order example:
//
//	walker := NewWalker(g, DFSPreOrder)
//	walker.Walk("A", print)
//	// Output: A, B, D, E, C (deep first)
//
// complexity:
//   - time : O(V + E) for reachable vertices
//   - space: O(V) for BFS queue or DFS recursion stack
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

// dfs performs depth-first search using recursion.
//
// Pre-Order: visit BEFORE recursing (root → children)
// Post-Order: visit AFTER recursing (children → root)
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

// bfs performs breadth-first search using a queue.
//
//	Queue-based BFS:
//
//	1. Enqueue start, mark visited
//	2. While queue not empty:
//	   a. Dequeue vertex
//	   b. Visit vertex
//	   c. Enqueue unvisited neighbors
//
//	    A
//	   / \
//	  B   C
//	 /
//	D
//
//	Step 1: queue=[A], visited={A}
//	Step 2: visit A, queue=[B,C], visited={A,B,C}
//	Step 3: visit B, queue=[C,D], visited={A,B,C,D}
//	Step 4: visit C, queue=[D]
//	Step 5: visit D, queue=[]
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

// HasCycle detects if the graph contains a cycle.
//
// A cycle exists when you can start at a vertex and follow edges
// back to the same vertex.
//
//	No cycle:              Has cycle:
//	A ──► B ──► C          A ──► B
//	                        ↑     ↓
//	                        └──── C
//
// For directed graphs: uses recursion stack to detect back edges.
// For undirected graphs: tracks parent to avoid false positives from bidirectional edges.
//
// complexity:
//   - time : O(V + E)
//   - space: O(V)
func (w *Walker[T]) HasCycle() bool {
	if w.graph.Directed() {
		return w.hasCycleDirected()
	}
	return w.hasCycleUndirected()
}

func (w *Walker[T]) hasCycleDirected() bool {
	stack := sets.New[T]()
	visited := sets.New[T]()

	var visit func(T) bool
	visit = func(n T) bool {
		if stack.Exists(n) {
			return true
		}
		if visited.Exists(n) {
			return false
		}

		visited.Add(n)
		stack.Add(n)

		for neighbor := range w.graph.Neighbors(n) {
			if visit(neighbor) {
				return true
			}
		}

		stack.Del(n)
		return false
	}

	for n := range w.graph.Vertex {
		if !visited.Exists(n) {
			if visit(n) {
				return true
			}
		}
	}

	return false
}

func (w *Walker[T]) hasCycleUndirected() bool {
	visited := sets.New[T]()

	var visit func(n T, parent *T) bool
	visit = func(n T, parent *T) bool {
		visited.Add(n)

		for neighbor := range w.graph.Neighbors(n) {
			if !visited.Exists(neighbor) {
				if visit(neighbor, &n) {
					return true
				}
			} else if parent == nil || neighbor != *parent {
				return true
			}
		}

		return false
	}

	for n := range w.graph.Vertex {
		if !visited.Exists(n) {
			if visit(n, nil) {
				return true
			}
		}
	}

	return false
}
