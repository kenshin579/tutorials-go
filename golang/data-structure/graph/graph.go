package graph

import "sort"

const (
	White = iota // 미방문
	Gray         // 진행 중
	Black        // 완료
)

// Graph는 방향 그래프를 인접 리스트로 표현한다.
type Graph struct {
	adj map[int][]int
}

func NewGraph() *Graph {
	return &Graph{adj: make(map[int][]int)}
}

// AddEdge는 u → v 간선을 추가한다.
func (g *Graph) AddEdge(u, v int) {
	g.adj[u] = append(g.adj[u], v)
	if _, ok := g.adj[v]; !ok {
		g.adj[v] = []int{}
	}
}

// Vertices는 모든 정점을 정렬된 순서로 반환한다.
func (g *Graph) Vertices() []int {
	verts := make([]int, 0, len(g.adj))
	for v := range g.adj {
		verts = append(verts, v)
	}
	sort.Ints(verts)
	return verts
}

// Transpose는 모든 간선의 방향을 뒤집은 전치 그래프를 반환한다.
func (g *Graph) Transpose() *Graph {
	gt := NewGraph()
	for u, neighbors := range g.adj {
		if _, ok := gt.adj[u]; !ok {
			gt.adj[u] = []int{}
		}
		for _, v := range neighbors {
			gt.adj[v] = append(gt.adj[v], u)
		}
	}
	return gt
}

// DFSResult는 DFS 수행 결과를 담는다.
type DFSResult struct {
	Discovery   map[int]int       // 발견 시점
	Finish      map[int]int       // 완료 시점
	FinishOrder []int             // 완료 순서
	EdgeType    map[[2]int]string // 간선 분류
}

// DFS는 타임스탬프와 간선 분류를 포함한 DFS를 수행한다.
func (g *Graph) DFS() *DFSResult {
	result := &DFSResult{
		Discovery: make(map[int]int),
		Finish:    make(map[int]int),
		EdgeType:  make(map[[2]int]string),
	}
	color := make(map[int]int)
	clock := 0

	var visit func(u int)
	visit = func(u int) {
		clock++
		result.Discovery[u] = clock
		color[u] = Gray

		for _, v := range g.adj[u] {
			edge := [2]int{u, v}
			switch color[v] {
			case White:
				result.EdgeType[edge] = "tree"
				visit(v)
			case Gray:
				result.EdgeType[edge] = "back"
			case Black:
				if result.Discovery[u] < result.Discovery[v] {
					result.EdgeType[edge] = "forward"
				} else {
					result.EdgeType[edge] = "cross"
				}
			}
		}

		color[u] = Black
		clock++
		result.Finish[u] = clock
		result.FinishOrder = append(result.FinishOrder, u)
	}

	for _, v := range g.Vertices() {
		if color[v] == White {
			visit(v)
		}
	}
	return result
}

// IsDAG는 그래프에 순환이 없는지(DAG인지) 확인한다.
func (g *Graph) IsDAG() bool {
	color := make(map[int]int)
	for _, v := range g.Vertices() {
		if color[v] == White {
			if g.hasCycle(v, color) {
				return false
			}
		}
	}
	return true
}

func (g *Graph) hasCycle(u int, color map[int]int) bool {
	color[u] = Gray
	for _, v := range g.adj[u] {
		if color[v] == Gray {
			return true
		}
		if color[v] == White && g.hasCycle(v, color) {
			return true
		}
	}
	color[u] = Black
	return false
}

// KahnTopologicalSort는 BFS 기반 위상정렬을 수행한다.
// 순환이 있으면 nil을 반환한다.
func (g *Graph) KahnTopologicalSort() []int {
	inDegree := make(map[int]int)
	for _, v := range g.Vertices() {
		inDegree[v] = 0
	}
	for _, neighbors := range g.adj {
		for _, v := range neighbors {
			inDegree[v]++
		}
	}

	queue := []int{}
	for _, v := range g.Vertices() {
		if inDegree[v] == 0 {
			queue = append(queue, v)
		}
	}

	result := []int{}
	for len(queue) > 0 {
		u := queue[0]
		queue = queue[1:]
		result = append(result, u)

		for _, v := range g.adj[u] {
			inDegree[v]--
			if inDegree[v] == 0 {
				queue = append(queue, v)
			}
		}
	}

	if len(result) != len(g.adj) {
		return nil
	}
	return result
}

// TopologicalSortDFS는 DFS Finish Order의 역순으로 위상정렬을 수행한다.
func (g *Graph) TopologicalSortDFS() []int {
	color := make(map[int]int)
	stack := []int{}

	var dfs func(u int)
	dfs = func(u int) {
		color[u] = Gray
		for _, v := range g.adj[u] {
			if color[v] == White {
				dfs(v)
			}
		}
		color[u] = Black
		stack = append(stack, u)
	}

	for _, v := range g.Vertices() {
		if color[v] == White {
			dfs(v)
		}
	}

	result := make([]int, len(stack))
	for i, v := range stack {
		result[len(stack)-1-i] = v
	}
	return result
}
