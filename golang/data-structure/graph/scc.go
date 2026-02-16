package graph

import "sort"

// SCCResult는 Kosaraju SCC 결과를 담는다.
type SCCResult struct {
	SCCOf      map[int]int // 정점 → SCC ID
	Components [][]int     // SCC ID → 정점 리스트
}

// Kosaraju는 2-pass DFS로 SCC를 구한다.
func (g *Graph) Kosaraju() *SCCResult {
	// Step 1: 1차 DFS — Finish Order 기록
	finishOrder := []int{}
	visited := make(map[int]bool)

	var dfs1 func(u int)
	dfs1 = func(u int) {
		visited[u] = true
		for _, v := range g.adj[u] {
			if !visited[v] {
				dfs1(v)
			}
		}
		finishOrder = append(finishOrder, u)
	}

	for _, v := range g.Vertices() {
		if !visited[v] {
			dfs1(v)
		}
	}

	// Step 2: 그래프 전치
	gt := g.Transpose()

	// Step 3: 2차 DFS — Finish Order 역순으로 전치 그래프 탐색
	visited = make(map[int]bool)
	result := &SCCResult{
		SCCOf: make(map[int]int),
	}
	sccID := 0

	var dfs2 func(u int, component *[]int)
	dfs2 = func(u int, component *[]int) {
		visited[u] = true
		*component = append(*component, u)
		result.SCCOf[u] = sccID
		for _, v := range gt.adj[u] {
			if !visited[v] {
				dfs2(v, component)
			}
		}
	}

	for i := len(finishOrder) - 1; i >= 0; i-- {
		u := finishOrder[i]
		if !visited[u] {
			component := []int{}
			dfs2(u, &component)
			sort.Ints(component)
			result.Components = append(result.Components, component)
			sccID++
		}
	}

	return result
}

// Condense는 SCC를 하나의 노드로 축약한 DAG를 생성한다.
func (g *Graph) Condense(scc *SCCResult) *Graph {
	dag := NewGraph()

	for i := range len(scc.Components) {
		if _, ok := dag.adj[i]; !ok {
			dag.adj[i] = []int{}
		}
	}

	seen := make(map[[2]int]bool)
	for u, neighbors := range g.adj {
		for _, v := range neighbors {
			su, sv := scc.SCCOf[u], scc.SCCOf[v]
			if su != sv {
				edge := [2]int{su, sv}
				if !seen[edge] {
					dag.AddEdge(su, sv)
					seen[edge] = true
				}
			}
		}
	}

	return dag
}
