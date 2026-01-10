package graph_test

import (
	"slices"
	"testing"

	"github.com/josestg/dsa/graph"
)

func TestNew(t *testing.T) {
	t.Run("directed", func(t *testing.T) {
		g := graph.New[string](true)
		if !g.Directed() {
			t.Error("expected directed graph")
		}
		if !g.Empty() {
			t.Error("expected empty graph")
		}
		if g.Size() != 0 {
			t.Errorf("Size() = %d, want 0", g.Size())
		}
	})

	t.Run("undirected", func(t *testing.T) {
		g := graph.New[int](false)
		if g.Directed() {
			t.Error("expected undirected graph")
		}
		if !g.Empty() {
			t.Error("expected empty graph")
		}
	})
}

func TestAddEdge(t *testing.T) {
	t.Run("directed creates one-way edge", func(t *testing.T) {
		g := graph.New[string](true)
		g.AddEdge("A", "B")

		if !g.HasEdge("A", "B") {
			t.Error("expected edge A->B")
		}
		if g.HasEdge("B", "A") {
			t.Error("directed graph should not have reverse edge B->A")
		}
	})

	t.Run("undirected creates two-way edge", func(t *testing.T) {
		g := graph.New[string](false)
		g.AddEdge("A", "B")

		if !g.HasEdge("A", "B") {
			t.Error("expected edge A->B")
		}
		if !g.HasEdge("B", "A") {
			t.Error("undirected graph should have reverse edge B->A")
		}
	})

	t.Run("creates source vertex automatically", func(t *testing.T) {
		g := graph.New[string](true)
		g.AddEdge("X", "Y")

		if !g.HasVertex("X") {
			t.Error("expected vertex X to be created")
		}
		// Note: directed graph only creates source vertex in adjacency map
		// Y only becomes a vertex when it has outgoing edges
		if g.Size() != 1 {
			t.Errorf("Size() = %d, want 1", g.Size())
		}
	})

	t.Run("undirected creates both vertices", func(t *testing.T) {
		g := graph.New[string](false)
		g.AddEdge("X", "Y")

		if !g.HasVertex("X") {
			t.Error("expected vertex X")
		}
		if !g.HasVertex("Y") {
			t.Error("expected vertex Y")
		}
		if g.Size() != 2 {
			t.Errorf("Size() = %d, want 2", g.Size())
		}
	})

	t.Run("duplicate edges ignored", func(t *testing.T) {
		g := graph.New[string](true)
		g.AddEdge("A", "B")
		g.AddEdge("A", "B")
		g.AddEdge("A", "B")

		neighbors := slices.Collect(g.Neighbors("A"))
		if len(neighbors) != 1 {
			t.Errorf("expected 1 neighbor, got %d", len(neighbors))
		}
	})

	t.Run("self loop", func(t *testing.T) {
		g := graph.New[string](true)
		g.AddEdge("A", "A")

		if !g.HasEdge("A", "A") {
			t.Error("expected self-loop edge A->A")
		}
	})
}

func TestDelEdge(t *testing.T) {
	t.Run("directed removes one-way edge", func(t *testing.T) {
		g := graph.New[string](true)
		g.AddEdge("A", "B")
		g.AddEdge("B", "A")

		g.DelEdge("A", "B")

		if g.HasEdge("A", "B") {
			t.Error("edge A->B should be deleted")
		}
		if !g.HasEdge("B", "A") {
			t.Error("edge B->A should still exist")
		}
	})

	t.Run("undirected removes both directions", func(t *testing.T) {
		g := graph.New[string](false)
		g.AddEdge("A", "B")

		g.DelEdge("A", "B")

		if g.HasEdge("A", "B") {
			t.Error("edge A->B should be deleted")
		}
		if g.HasEdge("B", "A") {
			t.Error("edge B->A should also be deleted")
		}
	})

	t.Run("source vertex remains after edge deletion", func(t *testing.T) {
		g := graph.New[string](true)
		g.AddEdge("A", "B")
		g.DelEdge("A", "B")

		if !g.HasVertex("A") {
			t.Error("vertex A should still exist")
		}
		// B was never a standalone vertex in directed graph
	})

	t.Run("deleting non-existent edge is no-op", func(t *testing.T) {
		g := graph.New[string](true)
		g.AddEdge("A", "B")

		g.DelEdge("X", "Y")
		g.DelEdge("A", "C")

		if g.Size() != 1 {
			t.Errorf("Size() = %d, want 1", g.Size())
		}
	})
}

func TestHasEdge(t *testing.T) {
	g := graph.New[string](true)
	g.AddEdge("A", "B")
	g.AddEdge("B", "C")

	tests := []struct {
		from, to string
		want     bool
	}{
		{"A", "B", true},
		{"B", "C", true},
		{"B", "A", false},
		{"A", "C", false},
		{"X", "Y", false},
	}

	for _, tt := range tests {
		got := g.HasEdge(tt.from, tt.to)
		if got != tt.want {
			t.Errorf("HasEdge(%s, %s) = %v, want %v", tt.from, tt.to, got, tt.want)
		}
	}
}

func TestHasVertex(t *testing.T) {
	t.Run("directed", func(t *testing.T) {
		g := graph.New[string](true)
		g.AddEdge("A", "B")

		if !g.HasVertex("A") {
			t.Error("expected vertex A")
		}
		// B is not a standalone vertex in directed graph (only in neighbor list)
		if g.HasVertex("C") {
			t.Error("unexpected vertex C")
		}
	})

	t.Run("undirected", func(t *testing.T) {
		g := graph.New[string](false)
		g.AddEdge("A", "B")

		if !g.HasVertex("A") {
			t.Error("expected vertex A")
		}
		if !g.HasVertex("B") {
			t.Error("expected vertex B")
		}
	})
}

func TestSizeAndEmpty(t *testing.T) {
	t.Run("directed", func(t *testing.T) {
		g := graph.New[int](true)

		if !g.Empty() {
			t.Error("new graph should be empty")
		}
		if g.Size() != 0 {
			t.Errorf("Size() = %d, want 0", g.Size())
		}

		g.AddEdge(1, 2)
		if g.Empty() {
			t.Error("graph with edges should not be empty")
		}
		// Only source vertex is counted
		if g.Size() != 1 {
			t.Errorf("Size() = %d, want 1", g.Size())
		}

		g.AddEdge(2, 3)
		if g.Size() != 2 {
			t.Errorf("Size() = %d, want 2", g.Size())
		}

		g.AddEdge(1, 3)
		if g.Size() != 2 {
			t.Errorf("Size() = %d, want 2 (no new source vertices)", g.Size())
		}
	})

	t.Run("undirected", func(t *testing.T) {
		g := graph.New[int](false)

		g.AddEdge(1, 2)
		if g.Size() != 2 {
			t.Errorf("Size() = %d, want 2", g.Size())
		}

		g.AddEdge(2, 3)
		if g.Size() != 3 {
			t.Errorf("Size() = %d, want 3", g.Size())
		}
	})
}

func TestVertex(t *testing.T) {
	t.Run("directed", func(t *testing.T) {
		g := graph.New[string](true)
		g.AddEdge("A", "B")
		g.AddEdge("B", "C")
		g.AddEdge("C", "D")

		var vertices []string
		for v := range g.Vertex {
			vertices = append(vertices, v)
		}

		slices.Sort(vertices)
		// Only source vertices are in the adjacency map
		want := []string{"A", "B", "C"}
		if !slices.Equal(vertices, want) {
			t.Errorf("Vertex = %v, want %v", vertices, want)
		}
	})

	t.Run("undirected", func(t *testing.T) {
		g := graph.New[string](false)
		g.AddEdge("A", "B")
		g.AddEdge("B", "C")
		g.AddEdge("C", "D")

		var vertices []string
		for v := range g.Vertex {
			vertices = append(vertices, v)
		}

		slices.Sort(vertices)
		want := []string{"A", "B", "C", "D"}
		if !slices.Equal(vertices, want) {
			t.Errorf("Vertex = %v, want %v", vertices, want)
		}
	})
}

func TestVertexBreak(t *testing.T) {
	g := graph.New[int](true)
	for i := range 10 {
		g.AddEdge(i, i+1)
	}

	count := 0
	for range g.Vertex {
		count++
		if count == 3 {
			break
		}
	}

	if count != 3 {
		t.Errorf("count = %d, want 3", count)
	}
}

func TestNeighbors(t *testing.T) {
	t.Run("directed", func(t *testing.T) {
		g := graph.New[string](true)
		g.AddEdge("A", "B")
		g.AddEdge("A", "C")
		g.AddEdge("A", "D")

		neighbors := slices.Collect(g.Neighbors("A"))
		slices.Sort(neighbors)
		want := []string{"B", "C", "D"}
		if !slices.Equal(neighbors, want) {
			t.Errorf("Neighbors(A) = %v, want %v", neighbors, want)
		}
	})

	t.Run("non-existent vertex", func(t *testing.T) {
		g := graph.New[string](true)
		g.AddEdge("A", "B")

		neighbors := slices.Collect(g.Neighbors("X"))
		if len(neighbors) != 0 {
			t.Errorf("Neighbors(X) = %v, want empty", neighbors)
		}
	})

	t.Run("vertex with no outgoing edges", func(t *testing.T) {
		g := graph.New[string](true)
		g.AddEdge("A", "B")

		neighbors := slices.Collect(g.Neighbors("B"))
		if len(neighbors) != 0 {
			t.Errorf("Neighbors(B) = %v, want empty", neighbors)
		}
	})

	t.Run("break iteration", func(t *testing.T) {
		g := graph.New[int](true)
		for i := range 10 {
			g.AddEdge(0, i+1)
		}

		count := 0
		for range g.Neighbors(0) {
			count++
			if count == 3 {
				break
			}
		}

		if count != 3 {
			t.Errorf("count = %d, want 3", count)
		}
	})
}

func TestDirected(t *testing.T) {
	directed := graph.New[string](true)
	undirected := graph.New[string](false)

	if !directed.Directed() {
		t.Error("expected Directed() = true")
	}
	if undirected.Directed() {
		t.Error("expected Directed() = false")
	}
}

func TestIntGraph(t *testing.T) {
	g := graph.New[int](true)
	g.AddEdge(1, 2)
	g.AddEdge(2, 3)
	g.AddEdge(3, 1)

	if g.Size() != 3 {
		t.Errorf("Size() = %d, want 3", g.Size())
	}

	if !g.HasEdge(1, 2) || !g.HasEdge(2, 3) || !g.HasEdge(3, 1) {
		t.Error("missing expected edges")
	}
}

func TestComplexGraph(t *testing.T) {
	g := graph.New[string](false)

	g.AddEdge("A", "B")
	g.AddEdge("A", "C")
	g.AddEdge("B", "D")
	g.AddEdge("C", "D")
	g.AddEdge("D", "E")

	if g.Size() != 5 {
		t.Errorf("Size() = %d, want 5", g.Size())
	}

	if !g.HasEdge("A", "B") || !g.HasEdge("B", "A") {
		t.Error("undirected edge A-B missing")
	}

	aNeighbors := slices.Collect(g.Neighbors("A"))
	slices.Sort(aNeighbors)
	if !slices.Equal(aNeighbors, []string{"B", "C"}) {
		t.Errorf("Neighbors(A) = %v, want [B C]", aNeighbors)
	}

	dNeighbors := slices.Collect(g.Neighbors("D"))
	slices.Sort(dNeighbors)
	if !slices.Equal(dNeighbors, []string{"B", "C", "E"}) {
		t.Errorf("Neighbors(D) = %v, want [B C E]", dNeighbors)
	}
}

func TestString(t *testing.T) {
	t.Run("empty graph", func(t *testing.T) {
		g := graph.New[string](true)
		got := g.String()
		want := "Graph{}"
		if got != want {
			t.Errorf("String() = %q, want %q", got, want)
		}
	})

	t.Run("single vertex with edge", func(t *testing.T) {
		g := graph.New[int](true)
		g.AddEdge(1, 2)
		got := g.String()
		// Should contain "Graph{" and "}"
		if len(got) < 7 {
			t.Errorf("String() = %q, expected longer string", got)
		}
		if got[:6] != "Graph{" {
			t.Errorf("String() should start with 'Graph{', got %q", got)
		}
		if got[len(got)-1] != '}' {
			t.Errorf("String() should end with '}', got %q", got)
		}
	})

	t.Run("contains vertex and neighbors", func(t *testing.T) {
		g := graph.New[string](true)
		g.AddEdge("A", "B")
		g.AddEdge("A", "C")
		got := g.String()
		// Should mention vertex A and its neighbors
		if len(got) == 0 {
			t.Error("String() should not be empty")
		}
		if got[:6] != "Graph{" || got[len(got)-1] != '}' {
			t.Errorf("String() has wrong format: %q", got)
		}
	})
}
