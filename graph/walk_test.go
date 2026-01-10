package graph_test

import (
	"slices"
	"testing"

	"github.com/josestg/dsa/graph"
)

func createTestGraph() *graph.Graph[string] {
	g := graph.New[string](true)
	g.AddEdge("A", "B")
	g.AddEdge("A", "C")
	g.AddEdge("B", "D")
	g.AddEdge("C", "D")
	return g
}

func TestWalker_BFS(t *testing.T) {
	g := createTestGraph()
	w := graph.NewWalker(g, graph.BFS)

	var visited []string
	w.Walk("A", func(v string) {
		visited = append(visited, v)
	})

	expected := []string{"A", "B", "C", "D"}
	if !slices.Equal(visited, expected) {
		t.Errorf("BFS got %v, want %v", visited, expected)
	}
	if !w.Explored() {
		t.Error("expected Explored() to be true")
	}
}

func TestWalker_DFSPreOrder(t *testing.T) {
	g := createTestGraph()
	w := graph.NewWalker(g, graph.DFSPreOrder)

	var visited []string
	w.Walk("A", func(v string) {
		visited = append(visited, v)
	})

	expected := []string{"A", "B", "D", "C"}
	if !slices.Equal(visited, expected) {
		t.Errorf("DFSPreOrder got %v, want %v", visited, expected)
	}
	if !w.Explored() {
		t.Error("expected Explored() to be true")
	}
}

func TestWalker_DFSPostOrder(t *testing.T) {
	g := createTestGraph()
	w := graph.NewWalker(g, graph.DFSPostOrder)

	var visited []string
	w.Walk("A", func(v string) {
		visited = append(visited, v)
	})

	expected := []string{"D", "B", "C", "A"}
	if !slices.Equal(visited, expected) {
		t.Errorf("DFSPostOrder got %v, want %v", visited, expected)
	}
	if !w.Explored() {
		t.Error("expected Explored() to be true")
	}
}

func TestWalker_WalkAll(t *testing.T) {
	g := graph.New[string](true)
	g.AddEdge("A", "B")
	g.AddEdge("C", "D")

	w := graph.NewWalker(g, graph.BFS)

	var visited []string
	w.WalkAll(func(v string) {
		visited = append(visited, v)
	})

	if len(visited) != 4 {
		t.Errorf("WalkAll got %d elements, want 4", len(visited))
	}
	slices.Sort(visited)
	expected := []string{"A", "B", "C", "D"}
	if !slices.Equal(visited, expected) {
		t.Errorf("WalkAll got %v, want %v (sorted)", visited, expected)
	}
	if !w.Explored() {
		t.Error("expected Explored() to be true")
	}
}

func TestWalker_VisitedAndExplored(t *testing.T) {
	g := createTestGraph()
	w := graph.NewWalker(g, graph.DFSPreOrder)

	if w.Visited("A") {
		t.Error("Visited(A) should be false before Walk")
	}
	if w.Explored() {
		t.Error("Explored() should be false before Walk")
	}

	w.Walk("A", func(_ string) {})

	if !w.Visited("A") {
		t.Error("Visited(A) should be true after Walk")
	}
	if !w.Visited("D") {
		t.Error("Visited(D) should be true after Walk")
	}
	if !w.Explored() {
		t.Error("Explored() should be true after Walk")
	}
}
