package graph

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 0 → 1 → 2 → 0 (SCC: {0,1,2})
// 1 → 3 → 4 → 5 → 3 (SCC: {3,4,5})
// 4 → 6 (SCC: {6})
func buildSCCGraph() *Graph {
	g := NewGraph()
	g.AddEdge(0, 1)
	g.AddEdge(1, 2)
	g.AddEdge(2, 0)
	g.AddEdge(1, 3)
	g.AddEdge(3, 4)
	g.AddEdge(4, 5)
	g.AddEdge(5, 3)
	g.AddEdge(4, 6)
	return g
}

func TestGraph_Kosaraju(t *testing.T) {
	g := buildSCCGraph()
	scc := g.Kosaraju()

	// SCC 3개
	assert.Equal(t, 3, len(scc.Components))

	// {0, 1, 2}는 같은 SCC
	assert.Equal(t, scc.SCCOf[0], scc.SCCOf[1])
	assert.Equal(t, scc.SCCOf[1], scc.SCCOf[2])

	// {3, 4, 5}는 같은 SCC
	assert.Equal(t, scc.SCCOf[3], scc.SCCOf[4])
	assert.Equal(t, scc.SCCOf[4], scc.SCCOf[5])

	// {6}은 별도 SCC
	assert.NotEqual(t, scc.SCCOf[0], scc.SCCOf[3])
	assert.NotEqual(t, scc.SCCOf[3], scc.SCCOf[6])
	assert.NotEqual(t, scc.SCCOf[0], scc.SCCOf[6])
}

func TestGraph_Kosaraju_SingleNodes(t *testing.T) {
	g := NewGraph()
	g.AddEdge(0, 1)
	g.AddEdge(2, 3)

	scc := g.Kosaraju()

	// 순환이 없으므로 각 정점이 독립 SCC
	assert.Equal(t, 4, len(scc.Components))
}

func TestGraph_Condense(t *testing.T) {
	g := buildSCCGraph()
	scc := g.Kosaraju()
	dag := g.Condense(scc)

	// 축약 그래프의 정점 수 = SCC 수
	assert.Equal(t, 3, len(dag.adj))

	// 축약 그래프는 DAG
	assert.True(t, dag.IsDAG())

	// 위상정렬이 가능해야 한다
	topo := dag.KahnTopologicalSort()
	assert.NotNil(t, topo)
	assert.Equal(t, 3, len(topo))
}

func TestReachability(t *testing.T) {
	g := buildSCCGraph()
	r := NewReachability(g)

	t.Run("same SCC - reachable both ways", func(t *testing.T) {
		// {0, 1, 2} 내부
		assert.True(t, r.IsReachable(0, 1))
		assert.True(t, r.IsReachable(1, 0))
		assert.True(t, r.IsReachable(2, 0))

		// {3, 4, 5} 내부
		assert.True(t, r.IsReachable(3, 5))
		assert.True(t, r.IsReachable(5, 3))
	})

	t.Run("cross SCC - forward reachable", func(t *testing.T) {
		assert.True(t, r.IsReachable(0, 3))
		assert.True(t, r.IsReachable(0, 6))
		assert.True(t, r.IsReachable(1, 6))
		assert.True(t, r.IsReachable(3, 6))
	})

	t.Run("cross SCC - backward not reachable", func(t *testing.T) {
		assert.False(t, r.IsReachable(3, 0))
		assert.False(t, r.IsReachable(6, 0))
		assert.False(t, r.IsReachable(6, 3))
	})
}

func TestReachability_DAG(t *testing.T) {
	// 순환이 없는 단순 DAG
	g := NewGraph()
	g.AddEdge(0, 1)
	g.AddEdge(1, 2)
	g.AddEdge(0, 2)

	r := NewReachability(g)

	assert.True(t, r.IsReachable(0, 1))
	assert.True(t, r.IsReachable(0, 2))
	assert.True(t, r.IsReachable(1, 2))
	assert.False(t, r.IsReachable(1, 0))
	assert.False(t, r.IsReachable(2, 0))
}
