package control_flow

type dfsBuilder struct {
	*ControlFlowGraph
	Visited  []bool
	Preorder []uint
	Parent   []uint
	NextTime uint
}

type DfsResult struct {
	// Preorder[v] is the index of v in the DFS pre-order traversal.
	Preorder []uint

	// Parent[v] is the parent of v in the DFS spanning tree.
	// The parent of the root of the tree is itself.
	Parent []uint
}

func newDfsBuilder(cfg *ControlFlowGraph) dfsBuilder {
	size := cfg.Size()
	return dfsBuilder{
		ControlFlowGraph: cfg,
		Visited:          make([]bool, size),
		Preorder:         make([]uint, size),
		Parent:           make([]uint, size),
		NextTime:         0,
	}
}

func (g *dfsBuilder) dfs(node uint, from uint) {
	if g.Visited[node] {
		return
	}

	g.Visited[node] = true
	g.Parent[node] = from
	g.Preorder[node] = g.NextTime
	g.NextTime++

	for _, next := range g.BasicBlocks[node].ForwardEdges {
		g.dfs(next, node)
	}
}

func (g *dfsBuilder) toDfsResult() DfsResult {
	return DfsResult{
		Preorder: g.Preorder,
		Parent:   g.Parent,
	}
}
