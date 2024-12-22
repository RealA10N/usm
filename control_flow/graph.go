package control_flow

type Node struct {
	ForwardEdges  []uint
	BackwardEdges []uint
}

type Graph struct {
	Nodes []Node
}

func NewEmptyGraph(size uint) Graph {
	return Graph{
		Nodes: make([]Node, size),
	}
}

func NewGraph(forwardEdges [][]uint) Graph {
	n := uint(len(forwardEdges))
	nodes := make([]Node, n)
	for from, edges := range forwardEdges {
		nodes[from].ForwardEdges = edges
		for _, to := range edges {
			nodes[to].BackwardEdges = append(nodes[to].BackwardEdges, uint(from))
		}
	}

	return Graph{Nodes: nodes}
}

// MARK: Operations

func (g *Graph) AddEdge(from, to uint) {
	g.Nodes[from].ForwardEdges = append(g.Nodes[from].ForwardEdges, to)
	g.Nodes[to].BackwardEdges = append(g.Nodes[to].BackwardEdges, from)
}

// MARK: Queries

// Returns the number of nodes in the graph.
func (g *Graph) Size() uint {
	return uint(len(g.Nodes))
}

// MARK: Algorithms

// Returns the 'Dfs' type that contains information about the graph that have
// been collected in a linear-time depth-first traversal of the graph from
// the provided node as the initial location.
func (g *Graph) Dfs(root uint) Dfs {
	builder := newDfsBuilder(g)
	builder.dfs(root, root, 0)
	return builder.toDfs()
}

// Returns the 'DominatorTree' type that encapsulates the Dominator Tree data
// structure, and provides efficient queries of the dominators and immediate
// dominators of nodes.
//
// Construction of the data structure is based on the Lengauer-Tarjan algorithm:
// https://doi.org/10.1145/357062.357071
func (g *Graph) DominatorTree(entry uint) DominatorTree {
	lengauerTarjan := newLengauerTarjanContext(g, entry)
	immDom := lengauerTarjan.LengauerTarjan()
	return DominatorTree{
		ImmDom: immDom,
		Dfs:    lengauerTarjan.Dfs,
	}
}

func (g *Graph) ControlFlowGraph(entry uint) ControlFlowGraph {
	builder := newControlFlowGraphBuilder(g)
	builder.exploreBasicBlock(entry)
	return builder.ControlFlowGraph
}

func (g *Graph) DominatorJoinGraph(entry uint) DominatorJoinGraph {
	dominatorTree := g.DominatorTree(entry)
	joinGraph := newJoinGraph(g, &dominatorTree)
	return DominatorJoinGraph{
		DominatorTree: dominatorTree,
		JoinGraph:     joinGraph,
	}
}
