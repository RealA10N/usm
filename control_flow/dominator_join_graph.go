// This file contains the implementation of the dominator join graph (DJ-Graph)
// data structure, first introduced in the paper by Sreedhar & Gao:
// https://doi.org/10.1145/199448.199464

package control_flow

import (
	"slices"
	"sort"

	"golang.org/x/exp/constraints"
)

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

func removeDuplicates[T constraints.Integer](slice []T) []T {
	sort.Slice(slice, func(i, j int) bool { return slice[i] < slice[j] })
	return slices.Compact(slice)
}

// MARK: Queries

// Computes the dominator frontier of the provided node in linear time.
func (g *DominatorJoinGraph) DominatorFrontier(node uint) []uint {
	frontier := []uint{}
	for _, subtreeNode := range g.Subtree(node) {
		for _, joinNode := range g.JoinGraph.Nodes[subtreeNode].ForwardEdges {
			if g.IsDeeper(node, joinNode) {
				frontier = append(frontier, joinNode)
			}
		}
	}

	return removeDuplicates(frontier)
}
