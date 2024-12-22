package control_flow

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
			PreOrder:  make([]uint, n),
			PostOrder: make([]uint, n),
			Parent:    make([]uint, n),
			Depth:     make([]uint, n),
		},
		Visited:      make([]bool, n),
		NextPreTime:  0,
		NextPostTime: 0,
	}
}

func (g *dfsBuilder) dfs(node, from, depth uint) {
	if g.Visited[node] {
		return
	}

	g.Visited[node] = true
	g.Depth[node] = depth
	g.Parent[node] = from
	g.PreOrder[node] = g.NextPreTime
	g.NextPreTime++

	for _, next := range g.Nodes[node].ForwardEdges {
		g.dfs(next, node, depth+1)
	}

	g.PostOrder[node] = g.NextPostTime
	g.NextPostTime++
}

func (g *dfsBuilder) toDfs() Dfs {
	return g.Dfs
}
