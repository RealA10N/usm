// This file contains the implementation of the dominator join graph (DJ-Graph)
// data structure, first introduced in the paper by Sreedhar & Gao:
// https://doi.org/10.1145/199448.199464

package graph

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

// Computes the dominator frontier of the provided node in linear time.
//
// Note if your purpose is to compute the dominator frontier of a set of multiple
// nodes, or the iterated dominator frontier, there are better methods then
// calling this method multiple times, resulting in a quadratic time.
// See the IteratedDominatorFrontier method for more information.
//
// For more information, see https://doi.org/10.1145/199448.199464
func (g *DominatorJoinGraph) DominatorFrontier(node uint) []uint {
	frontier := []uint{}

	// TODO: Although the algorithm is linear anyways, perhaps in practice,
	// creating an array in the size of the graph can be expensive?
	visited := make([]bool, g.JoinGraph.Size())

	for _, subtreeNode := range g.Subtree(node) {
		for _, joinNode := range g.JoinGraph.Nodes[subtreeNode].ForwardEdges {
			if !visited[joinNode] && g.IsDeeper(node, joinNode) {
				visited[joinNode] = true
				frontier = append(frontier, joinNode)
			}
		}
	}

	return frontier
}

func (g *DominatorJoinGraph) IteratedDominatorFrontier(nodes []uint) []uint {
	n := g.JoinGraph.Size()
	frontier := []uint{}
	isNodeInFrontier := make([]bool, n)
	isProcessedSubtree := make([]bool, n)
	piggyBank := newPiggyBank(&g.Dfs, nodes)

	for depth := piggyBank.MaxDepth(); depth != ^uint(0); depth-- {
		for !piggyBank.IsEmptyAtDepth(depth) {
			node := piggyBank.Pop(depth)

			subtreeCurrentIndex := g.PreOrder[node]
			subtreeEndIndex := g.PreOrder[node] + g.SubtreeSize[node]
			for subtreeCurrentIndex < subtreeEndIndex {
				subtreeNode := g.PreOrderReversed[subtreeCurrentIndex]

				if isProcessedSubtree[subtreeNode] {
					subtreeCurrentIndex += g.SubtreeSize[subtreeNode]
					continue
				}

				isProcessedSubtree[subtreeNode] = true

				for _, joinNode := range g.JoinGraph.Nodes[subtreeNode].ForwardEdges {
					if !isNodeInFrontier[joinNode] && g.IsDeeper(node, joinNode) {
						isNodeInFrontier[joinNode] = true
						frontier = append(frontier, joinNode)
					}
				}

				subtreeCurrentIndex++
			}
		}
	}

	return frontier
}
