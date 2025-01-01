package graph

type dfsBuilder struct {
	*Graph
	Dfs

	Visited      []bool
	NextPreTime  uint
	NextPostTime uint
}

func newDfsBuilder(g *Graph) dfsBuilder {
	n := g.Size()
	return dfsBuilder{
		Graph: g,
		Dfs: Dfs{
			PreOrder:    make([]uint, n),
			PostOrder:   make([]uint, n),
			Timeline:    make([]uint, 0, 2*n),
			Parent:      make([]uint, n),
			Depth:       make([]uint, n),
			SubtreeSize: make([]uint, n),
		},
		Visited:      make([]bool, n),
		NextPreTime:  0,
		NextPostTime: 0,
	}
}

func (g *dfsBuilder) recordPre(node uint) {
	g.Timeline = append(g.Timeline, node)
	g.PreOrder[node] = g.NextPreTime
	g.NextPreTime++
}

func (g *dfsBuilder) recordPost(node uint) {
	n := uint(len(g.Nodes))
	g.Timeline = append(g.Timeline, n+node)
	g.PostOrder[node] = g.NextPostTime
	g.NextPostTime++
}

func (g *dfsBuilder) dfs(node, from, depth uint) bool {
	if g.Visited[node] {
		return false
	}

	g.Visited[node] = true
	g.Parent[node] = from
	g.Depth[node] = depth
	g.SubtreeSize[node] = 1

	g.recordPre(node)

	for _, next := range g.Nodes[node].ForwardEdges {
		if g.dfs(next, node, depth+1) {
			g.SubtreeSize[node] += g.SubtreeSize[next]
		}
	}

	g.recordPost(node)
	return true
}

func (g *dfsBuilder) toDfs() Dfs {
	g.Dfs.PreOrderReversed = reversePermutation(g.Dfs.PreOrder)
	return g.Dfs
}
