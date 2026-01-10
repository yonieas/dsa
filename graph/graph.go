// Package graph provides graph data structure and traversal algorithms.
//
// # What is a Graph?
//
// A graph models relationships between entities. It consists of vertices
// (nodes) connected by edges. Social networks, road maps, web page links,
// and dependency trees are all graphs.
//
// # Terminology
//
//	Vertex (Node): An entity in the graph
//	Edge:          A connection between two vertices
//	Directed:      Edges have direction (A to B is different from B to A)
//	Undirected:    Edges are bidirectional
//	Neighbor:      A vertex directly connected by an edge
//	Degree:        Number of edges connected to a vertex
//	Path:          Sequence of vertices connected by edges
//	Cycle:         A path that starts and ends at the same vertex
//
// # Representation
//
// This implementation uses an adjacency list, where each vertex maps to a
// list of its neighbors:
//
//	A: [B, C]
//	B: [A, D]
//	C: [A]
//	D: [B]
//
// Adjacency lists use O(V + E) space, efficient for sparse graphs (few edges
// relative to vertices). Dense graphs might prefer an adjacency matrix with
// O(V^2) space but O(1) edge lookup.
//
// # Traversals
//
// The package provides BFS (level-by-level) and DFS (go deep first) traversals.
// BFS finds shortest paths in unweighted graphs. DFS is useful for cycle
// detection, topological sorting, and finding connected components.
//
// # Complexity
//
//	AddEdge:    O(1)
//	HasEdge:    O(degree) where degree is the number of neighbors
//	BFS/DFS:    O(V + E)
//	Space:      O(V + E)
//
// # Further Reading
//
// CLRS "Introduction to Algorithms", Chapters 22-24.
// Sedgewick "Algorithms", Part 5 (Graphs).
// https://en.wikipedia.org/wiki/Graph_(abstract_data_type)
package graph

import (
	"fmt"
	"iter"

	"github.com/josestg/dsa/hashmap"
	"github.com/josestg/dsa/linkedlist"
	"github.com/josestg/dsa/sequence"
)

// Graph represents a graph using an adjacency list.
//
// Undirected graph example:
//
//	A ——— B
//	|     |
//	C ——— D
//
// Adjacency list:
//
//	A → [B, C]
//	B → [A, D]
//	C → [A, D]
//	D → [B, C]
//
// Directed graph example:
//
//	A ──► B
//	│     │
//	▼     ▼
//	C ◄── D
//
// Adjacency list:
//
//	A → [B, C]
//	B → [D]
//	C → []
//	D → [C]
type Graph[V comparable] struct {
	directed  bool
	adjacency *hashmap.HashMap[V, *linkedlist.SinglyLinkedList[V]]
}

// New creates an empty graph.
//
// Parameters:
//   - directed: true for directed graph, false for undirected
//
// Directed graph: Edge A→B does NOT create edge B→A
// Undirected graph: Edge A—B creates both A→B and B→A
//
// complexity:
//   - time : O(1)
//   - space: O(1)
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

// Size returns the number of vertices.
//
//	    A ──► B
//	    │     │
//	    ▼     ▼
//	    C ◄── D
//
//	Size() → 4
//
// complexity:
//   - time : O(1)
//   - space: O(1)
func (g *Graph[V]) Size() int {
	return g.adjacency.Size()
}

// Empty returns true if the graph has no vertices.
//
// complexity:
//   - time : O(1)
//   - space: O(1)
func (g *Graph[V]) Empty() bool {
	return g.adjacency.Empty()
}

// Directed returns true if this is a directed graph.
//
// complexity:
//   - time : O(1)
//   - space: O(1)
func (g *Graph[V]) Directed() bool {
	return g.directed
}

// AddEdge adds an edge between two vertices.
//
// Directed graph - AddEdge(A, B):
//
//	Before:              After:
//	A    B               A ──► B
//
// Undirected graph - AddEdge(A, B):
//
//	Before:              After:
//	A    B               A ◄──► B
//
// Both vertices are created if they don't exist.
// Duplicate edges are ignored.
//
// complexity:
//   - time : O(degree) to check for duplicates
//   - space: O(1)
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

// DelEdge removes an edge between two vertices.
//
// Directed graph - DelEdge(A, B):
//
//	Before:              After:
//	A ──► B              A    B
//
// Undirected graph - DelEdge(A, B):
//
//	Before:              After:
//	A ◄──► B             A    B
//
// Vertices are NOT removed even if they have no edges.
// If edge doesn't exist, no action is taken.
//
// complexity:
//   - time : O(degree)
//   - space: O(1)
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

// HasEdge checks if an edge exists from one vertex to another.
//
//	A ──► B ──► C
//
//	HasEdge(A, B) → true
//	HasEdge(B, A) → false (directed graph)
//	HasEdge(A, C) → false (no direct edge)
//
// complexity:
//   - time : O(degree)
//   - space: O(1)
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

// HasVertex checks if a vertex exists in the graph.
//
//	Graph: A ──► B
//
//	HasVertex(A) → true
//	HasVertex(C) → false
//
// complexity:
//   - time : O(1) average
//   - space: O(1)
func (g *Graph[V]) HasVertex(v V) bool {
	return g.adjacency.Exists(v)
}

// Vertex iterates over all vertices in the graph.
//
//	    A ──► B
//	    │     │
//	    ▼     ▼
//	    C ◄── D
//
//	for v := range g.Vertex {
//	    fmt.Println(v)  // A, B, C, D (order not guaranteed)
//	}
//
// complexity:
//   - time : O(V)
//   - space: O(1)
func (g *Graph[V]) Vertex(yield func(V) bool) {
	for v := range g.adjacency.Iter {
		if !yield(v.Key()) {
			break
		}
	}
}

// Neighbors returns an iterator over a vertex's neighbors.
//
//	    A ──► B
//	    │     │
//	    ▼     ▼
//	    C ◄── D
//
//	for n := range g.Neighbors(A) {
//	    fmt.Println(n)  // B, C
//	}
//
// If vertex doesn't exist, returns empty iterator.
//
// complexity:
//   - time : O(degree)
//   - space: O(1)
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

// String returns a string representation of the graph.
//
//	    A ──► B
//	    │
//	    ▼
//	    C
//
//	String() → "Graph{A: [B C], B: [], C: []}" (order may vary)
//
// complexity:
//   - time : O(V + E)
//   - space: O(V + E)
func (g *Graph[V]) String() string {
	result := "Graph{"
	first := true
	for entry := range g.adjacency.Iter {
		if !first {
			result += ", "
		}
		first = false
		result += fmt.Sprintf("%v: %s", entry.Key(), sequence.String(entry.Value().Iter))
	}
	result += "}"
	return result
}
