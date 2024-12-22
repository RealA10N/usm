// This file contains the implementation of the dominator join graph (DJ-Graph)
// data structure, first introduced in the paper by Sreedhar & Gao:
// https://doi.org/10.1145/199448.199464

package control_flow

type DominatorJoinGraph struct {
	DominatorTree
	JoinGraph Graph
}

// Provided a graph and it's dominator tree, we construct a new graph containing
// only join edges from the original graph.
func newJoinGraph(g *Graph, d *DominatorTree) Graph {
	n := g.Size()
	joinGraph := NewEmptyGraph(n)

	for from := uint(0); from < n; from++ {
		for _, to := range g.Nodes[from].ForwardEdges {
			if !d.IsStrictDominatorOf(from, to) {
				joinGraph.AddEdge(from, to)
			}
		}
	}

	return joinGraph
}
