package graph_test

import (
	"testing"

	"github.com/josestg/dsa/graph"
)

func TestCycleDetection(t *testing.T) {
	t.Run("directed no cycle", func(t *testing.T) {
		g := graph.New[string](true)
		g.AddEdge("A", "B")
		g.AddEdge("B", "C")
		w := graph.NewWalker(g, graph.BFS)
		if w.HasCycle() {
			t.Error("Directed A->B->C should NOT have cycle")
		}
	})

	t.Run("directed with cycle", func(t *testing.T) {
		g := graph.New[string](true)
		g.AddEdge("A", "B")
		g.AddEdge("B", "C")
		g.AddEdge("C", "A")
		w := graph.NewWalker(g, graph.BFS)
		if !w.HasCycle() {
			t.Error("Directed A->B->C->A SHOULD have cycle")
		}
	})

	t.Run("undirected simple edge", func(t *testing.T) {
		g := graph.New[string](false)
		g.AddEdge("A", "B")
		w := graph.NewWalker(g, graph.BFS)
		if w.HasCycle() {
			t.Error("Undirected A-B should NOT have cycle")
		}
	})

	t.Run("undirected triangle", func(t *testing.T) {
		g := graph.New[string](false)
		g.AddEdge("A", "B")
		g.AddEdge("B", "C")
		g.AddEdge("C", "A")
		w := graph.NewWalker(g, graph.BFS)
		if !w.HasCycle() {
			t.Error("Undirected triangle SHOULD have cycle")
		}
	})
}
