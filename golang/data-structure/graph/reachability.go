package graph

// Reachability는 O(1) 도달 가능성 판정을 위한 전처리 구조체다.
// 교육 목적으로 map을 사용했다. 프로덕션에서는 []uint64 기반 bitset을 사용하면
// OR 연산 한 번으로 도달 가능 집합을 합산할 수 있어 성능이 크게 향상된다.
type Reachability struct {
	scc   *SCCResult
	reach map[int]map[int]bool // SCC ID → 도달 가능한 SCC ID 집합
}

// NewReachability는 그래프에 대해 전처리를 수행한다.
func NewReachability(g *Graph) *Reachability {
	// 1. SCC 분해
	scc := g.Kosaraju()

	// 2. DAG 축약
	dag := g.Condense(scc)

	// 3. 위상정렬 (Kahn)
	topoOrder := dag.KahnTopologicalSort()

	// 4. Bitset 역순 전파
	reach := make(map[int]map[int]bool)
	for i := range len(scc.Components) {
		reach[i] = map[int]bool{i: true}
	}

	for i := len(topoOrder) - 1; i >= 0; i-- {
		u := topoOrder[i]
		for _, v := range dag.adj[u] {
			for sccID := range reach[v] {
				reach[u][sccID] = true
			}
		}
	}

	return &Reachability{scc: scc, reach: reach}
}

// IsReachable은 정점 from에서 정점 to로 도달 가능한지 O(1)으로 판정한다.
func (r *Reachability) IsReachable(from, to int) bool {
	sccFrom := r.scc.SCCOf[from]
	sccTo := r.scc.SCCOf[to]
	return r.reach[sccFrom][sccTo]
}
