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
	"strings"

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

// SCORE: 10
func (g *Graph[V]) ensureNode(v V) *linkedlist.SinglyLinkedList[V] {
	// hint: 1) get neighbors list from g.adjacency.Get(v)
	//       2) if not found, create new list and g.adjacency.Put(v, list)
	//       3) return the neighbors list
	if list, found := g.adjacency.Get(v); found {
		return list
	}
	newList := linkedlist.NewSinglyLinkedList[V]()
	g.adjacency.Put(v, newList)
	return newList
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
//
// SCORE: 20
func (g *Graph[V]) AddEdge(from, to V) {
	// hint: 1) list := ensureNode(from)
	//       2) check if edge exists (iterate list, if v == to, return)
	//       3) list.Append(to)
	//       4) if undirected (!g.directed), also add reverse edge
	list := g.ensureNode(from)
	for l := range list.Iter {
		if l == to {
			return
		}
	}
	list.Append(to)
	if !g.directed {
		list := g.ensureNode(to)
		for l := range list.Iter {
			if l == from {
				return
			}
		}
		list.Append(from)
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
//
// SCORE: 15
func (g *Graph[V]) DelEdge(from, to V) {
	// hint: 1) get list from g.adjacency.Get(from)
	//       2) iterate with sequence.Enum, find 'to', call list.Remove(i)
	//       3) if undirected, repeat for reverse direction
	if list, found := g.adjacency.Get(from); found {
		for i, v := range sequence.Enum(list.Iter) {
			if v == to {
				list.Remove(i)
			}
		}
		if !g.Directed() {
			if list, found = g.adjacency.Get(to); found {
				for i, v := range sequence.Enum(list.Iter) {
					if v == from {
						list.Remove(i)
					}
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
//
// SCORE: 15
func (g *Graph[V]) HasEdge(from, to V) bool {
	// hint: 1) get list from g.adjacency.Get(from)
	//       2) iterate list, if v == to, return true
	//       3) return false
	if list, found := g.adjacency.Get(from); found {
		for l := range list.Iter {
			if l == to {
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
//
// SCORE: 5
func (g *Graph[V]) HasVertex(v V) bool {
	// hint: return g.adjacency.Exists(v)
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
//
// SCORE: 10
func (g *Graph[V]) Vertex(yield func(V) bool) {
	// hint: iterate g.adjacency.Iter, yield entry.Key()
	for v := range g.adjacency.Iter {
		if !yield(v.Key()) {
			return
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
//
// SCORE: 15
func (g *Graph[V]) Neighbors(v V) iter.Seq[V] {
	// hint: return a function that gets list from adjacency
	//       and iterates list.Iter, yielding each neighbor
	return func(yield func(V) bool) {
		if list, found := g.adjacency.Get(v); found {
			for v := range list.Iter {
				if !yield(v) {
					return
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
//
// SCORE: 10
func (g *Graph[V]) String() string {
	// hint: iterate g.adjacency.Iter, format as "Graph{V: [neighbors], ...}"
	//       use sequence.String(entry.Value().Iter) for neighbors
	var sb strings.Builder
	sb.WriteString("Graph{")
	for v := range g.adjacency.Iter {
		n := sequence.String(v.Value().Iter)
		sb.WriteString(fmt.Sprintf("%v: [%s]", v, n))
	}
	sb.WriteByte('}')
	return sb.String()
}
