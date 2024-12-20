package control_flow

type DfsBuilder[InstT SupportsControlFlow] struct {
	*ControlFlowGraph[InstT]
	Visited  []bool
	Order    []uint
	NextTime uint
}

func (g *DfsBuilder[InstT]) preOrderDfs(node uint) {
	if g.Visited[node] {
		return
	}

	g.Visited[node] = true
	g.Order[node] = g.NextTime
	g.NextTime++

	for _, next := range g.BasicBlocks[node].ForwardEdges {
		g.preOrderDfs(next)
	}
}
