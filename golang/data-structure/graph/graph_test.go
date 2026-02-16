package graph

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func buildDAG() *Graph {
	// 0 → 1 → 3
	// 0 → 2 → 3 → 4
	g := NewGraph()
	g.AddEdge(0, 1)
	g.AddEdge(0, 2)
	g.AddEdge(1, 3)
	g.AddEdge(2, 3)
	g.AddEdge(3, 4)
	return g
}

func TestGraph_AddEdge(t *testing.T) {
	g := NewGraph()
	g.AddEdge(0, 1)
	g.AddEdge(0, 2)

	assert.Equal(t, []int{1, 2}, g.adj[0])
	assert.Equal(t, []int{}, g.adj[1])
	assert.Equal(t, []int{}, g.adj[2])
}

func TestGraph_Vertices(t *testing.T) {
	g := buildDAG()
	assert.Equal(t, []int{0, 1, 2, 3, 4}, g.Vertices())
}

func TestGraph_Transpose(t *testing.T) {
	g := NewGraph()
	g.AddEdge(0, 1)
	g.AddEdge(1, 2)
	g.AddEdge(2, 0)

	gt := g.Transpose()

	// 전치 그래프에서 1→0, 2→1, 0→2 간선이 있어야 한다
	assert.Contains(t, gt.adj[1], 0)
	assert.Contains(t, gt.adj[2], 1)
	assert.Contains(t, gt.adj[0], 2)
}

func TestGraph_DFS(t *testing.T) {
	g := buildDAG()
	result := g.DFS()

	// 모든 정점에 Discovery/Finish가 기록되었는지 확인
	for _, v := range g.Vertices() {
		assert.Contains(t, result.Discovery, v)
		assert.Contains(t, result.Finish, v)
		assert.Greater(t, result.Finish[v], result.Discovery[v])
	}

	// Tree Edge 확인
	assert.Equal(t, "tree", result.EdgeType[[2]int{0, 1}])

	// FinishOrder 길이 확인
	assert.Equal(t, 5, len(result.FinishOrder))
}

func TestGraph_DFS_EdgeClassification(t *testing.T) {
	// Back Edge가 있는 그래프
	g := NewGraph()
	g.AddEdge(0, 1)
	g.AddEdge(1, 2)
	g.AddEdge(2, 0) // Back Edge

	result := g.DFS()

	assert.Equal(t, "tree", result.EdgeType[[2]int{0, 1}])
	assert.Equal(t, "tree", result.EdgeType[[2]int{1, 2}])
	assert.Equal(t, "back", result.EdgeType[[2]int{2, 0}])
}

func TestGraph_IsDAG(t *testing.T) {
	t.Run("DAG", func(t *testing.T) {
		dag := buildDAG()
		assert.True(t, dag.IsDAG())
	})

	t.Run("cyclic graph", func(t *testing.T) {
		cyclic := NewGraph()
		cyclic.AddEdge(0, 1)
		cyclic.AddEdge(1, 2)
		cyclic.AddEdge(2, 0)
		assert.False(t, cyclic.IsDAG())
	})
}

func TestGraph_KahnTopologicalSort(t *testing.T) {
	t.Run("DAG", func(t *testing.T) {
		g := buildDAG()
		result := g.KahnTopologicalSort()

		assert.NotNil(t, result)
		assert.Equal(t, 5, len(result))

		// 위상 순서 검증: 모든 간선 u→v에 대해 u가 v보다 앞에 있어야 한다
		pos := make(map[int]int)
		for i, v := range result {
			pos[v] = i
		}
		edges := [][2]int{{0, 1}, {0, 2}, {1, 3}, {2, 3}, {3, 4}}
		for _, e := range edges {
			assert.Less(t, pos[e[0]], pos[e[1]], "edge %d→%d violated", e[0], e[1])
		}
	})

	t.Run("cyclic graph returns nil", func(t *testing.T) {
		g := NewGraph()
		g.AddEdge(0, 1)
		g.AddEdge(1, 2)
		g.AddEdge(2, 0)

		assert.Nil(t, g.KahnTopologicalSort())
	})
}

func TestGraph_TopologicalSortDFS(t *testing.T) {
	g := buildDAG()
	result := g.TopologicalSortDFS()

	assert.Equal(t, 5, len(result))

	pos := make(map[int]int)
	for i, v := range result {
		pos[v] = i
	}
	edges := [][2]int{{0, 1}, {0, 2}, {1, 3}, {2, 3}, {3, 4}}
	for _, e := range edges {
		assert.Less(t, pos[e[0]], pos[e[1]], "edge %d→%d violated", e[0], e[1])
	}
}
