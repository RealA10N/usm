package graph

import (
	"slices"
)

type Node struct {
	ForwardEdges  []uint
	BackwardEdges []uint
}

type Graph struct {
	Nodes []Node
}

func NewEmptyGraph(size uint) Graph {
	return Graph{Nodes: make([]Node, size)}
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

// Converts a parent array that implicitly represents a rooted tree (by
// representing back edges only), into a full explicit graph representation,
// with explicit front edges.
//
// Assumes that the root is identified by the only node v for which
// parent[v] = v.
func NewGraphFromRootedTree(parents []uint) Graph {
	n := uint(len(parents))
	g := Graph{Nodes: make([]Node, n)}
	for u := uint(0); u < n; u++ {
		if parents[u] != u {
			g.AddEdge(parents[u], u)
		}
	}
	return g
}

// MARK: Operations

func (g *Graph) AddEdge(from, to uint) {
	g.Nodes[from].ForwardEdges = append(g.Nodes[from].ForwardEdges, to)
	g.Nodes[to].BackwardEdges = append(g.Nodes[to].BackwardEdges, from)
}

func (g *Graph) Equal(gt *Graph) bool {
	n := g.Size()
	if n != gt.Size() {
		return false
	}

	for u := uint(0); u < n; u++ {
		// Edges are not guaranteed to be sorted.
		// We do take a performance hit here, since we need to sort every edge
		// slice before we compare for equality. However the decision was made
		// that keeping the insertion of nodes O(1) is better.
		// Also, a comparison of whole graphs is not used as much in code, and
		// mainly used in testing.
		// Also, we we DO sort edges in place, and thus sorting them again
		// in a subsequent graph comparison cost only O(|E|), assuming that
		// the sorting algorithm performs in linear time given a sorted input,
		// (which is true for any state-of-the-art sorting algorithm).
		slices.Sort(g.Nodes[u].ForwardEdges)
		slices.Sort(gt.Nodes[u].ForwardEdges)
		if !slices.Equal(g.Nodes[u].ForwardEdges, gt.Nodes[u].ForwardEdges) {
			return false
		}
	}

	return true
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
	immediateDominators := lengauerTarjan.LengauerTarjan()

	// We do not need to store the explicit full graph representation of
	// the dominator tree, however we do need to perform a DFS traversal on
	// the dominator tree for efficient queries, and thus we do need to
	// create a (temporary) explicit representation of the tree.
	dominatorTreeGraph := NewGraphFromRootedTree(immediateDominators)

	return DominatorTree{
		Graph:  g,
		Dfs:    dominatorTreeGraph.Dfs(entry),
		ImmDom: immediateDominators,
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
