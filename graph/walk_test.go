package graph_test

import (
	"testing"

	"github.com/josestg/dsa/graph"
	"github.com/stretchr/testify/assert"
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
	assert.Equal(t, expected, visited)
	assert.True(t, w.Explored())
}

func TestWalker_DFSPreOrder(t *testing.T) {
	g := createTestGraph()
	w := graph.NewWalker(g, graph.DFSPreOrder)

	var visited []string
	w.Walk("A", func(v string) {
		visited = append(visited, v)
	})

	// PreOrder may vary depending on neighbor order, but expected consistent here.
	expected := []string{"A", "B", "D", "C"}
	assert.Equal(t, expected, visited)
	assert.True(t, w.Explored())
}

func TestWalker_DFSPostOrder(t *testing.T) {
	g := createTestGraph()
	w := graph.NewWalker(g, graph.DFSPostOrder)

	var visited []string
	w.Walk("A", func(v string) {
		visited = append(visited, v)
	})

	expected := []string{"D", "B", "C", "A"}
	assert.Equal(t, expected, visited)
	assert.True(t, w.Explored())
}

func TestWalker_WalkAll(t *testing.T) {
	g := graph.New[string](true)
	g.AddEdge("A", "B")
	g.AddEdge("C", "D") // disconnected component.

	w := graph.NewWalker(g, graph.BFS)

	var visited []string
	w.WalkAll(func(v string) {
		visited = append(visited, v)
	})

	assert.ElementsMatch(t, []string{"A", "B", "C", "D"}, visited)
	assert.True(t, w.Explored())
}

func TestWalker_VisitedAndExplored(t *testing.T) {
	g := createTestGraph()
	w := graph.NewWalker(g, graph.DFSPreOrder)

	assert.False(t, w.Visited("A"))
	assert.False(t, w.Explored())

	w.Walk("A", func(_ string) {})

	assert.True(t, w.Visited("A"))
	assert.True(t, w.Visited("D"))
	assert.True(t, w.Explored())
}
