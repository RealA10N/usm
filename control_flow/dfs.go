package control_flow

type dfsBuilder struct {
	*ControlFlowGraph

	Visited []bool
	Parent  []uint

	PreOrder  []uint
	PostOrder []uint

	NextPreTime  uint
	NextPostTime uint
}

type DfsResult struct {
	// PreOrder[v] is the index of v in the DFS pre-order traversal.
	PreOrder []uint

	// PostOrder[v] is the index of v in the DFS post-order traversal.
	PostOrder []uint

	// Parent[v] is the parent of v in the DFS spanning tree.
	// The parent of the root of the tree is itself.
	Parent []uint
}

func newDfsBuilder(cfg *ControlFlowGraph) dfsBuilder {
	n := cfg.Size()
	return dfsBuilder{
		ControlFlowGraph: cfg,
		Visited:          make([]bool, n),
		Parent:           make([]uint, n),
		PreOrder:         make([]uint, n),
		PostOrder:        make([]uint, n),
		NextPreTime:      0,
		NextPostTime:     0,
	}
}

func (g *dfsBuilder) dfs(node uint, from uint) {
	if g.Visited[node] {
		return
	}

	g.Visited[node] = true
	g.Parent[node] = from
	g.PreOrder[node] = g.NextPreTime
	g.NextPreTime++

	for _, next := range g.BasicBlocks[node].ForwardEdges {
		g.dfs(next, node)
	}

	g.PostOrder[node] = g.NextPostTime
	g.NextPostTime++
}

func (g *dfsBuilder) toDfsResult() DfsResult {
	return DfsResult{
		PreOrder:  g.PreOrder,
		PostOrder: g.PostOrder,
		Parent:    g.Parent,
	}
}
