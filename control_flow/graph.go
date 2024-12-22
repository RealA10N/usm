package control_flow

type Node struct {
	ForwardEdges  []uint
	BackwardEdges []uint
}

type Graph struct {
	Nodes []Node
}

func NewEmptyGraph() Graph {
	return Graph{
		Nodes: []Node{},
	}
}

func NewGraph(n uint, forwardEdges [][]uint) Graph {
	nodes := make([]Node, n)
	for from, edges := range forwardEdges {
		nodes[from].ForwardEdges = edges
		for _, to := range edges {
			nodes[to].BackwardEdges = append(nodes[to].BackwardEdges, uint(from))
		}
	}

	return Graph{Nodes: nodes}
}

// Returns the number of nodes in the graph.
func (g *Graph) Size() uint {
	return uint(len(g.Nodes))
}

// Returns the 'Dfs' type that contains information about the graph that have
// been collected in a linear-time depth-first traversal of the graph from
// the provided node as the initial location.
func (g *Graph) Dfs(root uint) Dfs {
	builder := newDfsBuilder(g)
	builder.dfs(root, root)
	return builder.toDfs()
}

// Returns the 'DominatorTree' type that encapsulates the Dominator Tree data
// structure, and provides efficient queries of the dominators and immediate
// dominators of nodes.
//
// Construction of the data structure is based on the Lengauer-Tarjan algorithm:
// https://doi.org/10.1145/357062.357071
func (g *Graph) DominatorTree() DominatorTree {
	lengauerTarjan := newLengauerTarjanContext(g)
	immDom := lengauerTarjan.LengauerTarjan()
	return DominatorTree{
		ImmDom: immDom,
	}
}
