// This file contains an implementation of the "PiggyBank" data structure,
// as described in section 4 of the Sreedhar & Gao Paper about DJ-Graphs
// and construction of Iterative Dominator Frontier in linear time:
// https://doi.org/10.1145/199448.199464

package graph

type piggyBank struct {
	// We only use the 'Depth' field from the Dfs embedded struct.
	Depth []uint

	// This is the essence of the data structure.
	// DepthToNodes[l] contains a set of dominator frontiers that have
	// been discovered at depth l (by join edges from <= depths), and that are
	// have been yet to be processed. You can think of it as a stack of next
	// to process nodes in each level.
	DepthToNodes [][]uint
}

func getMaxDepth(depth []uint, nodes []uint) uint {
	max := uint(0)
	for _, node := range nodes {
		if depth[node] > max {
			max = depth[node]
		}
	}
	return max
}

func newPiggyBank(dfs *Dfs, nodes []uint) piggyBank {
	maxDepth := getMaxDepth(dfs.Depth, nodes)

	piggyBank := piggyBank{
		Depth:        dfs.Depth,
		DepthToNodes: make([][]uint, maxDepth+1),
	}

	for _, node := range nodes {
		piggyBank.Push(node)
	}

	return piggyBank
}

// Its OK to NOT use a pointer receiver here, as we are not modifying the
// piggyBank struct itself, but rather the underlying slices in all methods.

// MARK: Queries

func (pb piggyBank) IsEmptyAtDepth(depth uint) bool {
	return len(pb.DepthToNodes[depth]) == 0
}

func (pb piggyBank) MaxDepth() uint {
	return uint(len(pb.DepthToNodes) - 1)
}

// MARK: Operations

func (pb piggyBank) Push(node uint) {
	depth := pb.Depth[node]
	pb.DepthToNodes[depth] = append(pb.DepthToNodes[depth], node)
}

func (pb piggyBank) Pop(depth uint) uint {
	len := len(pb.DepthToNodes[depth])
	node := pb.DepthToNodes[depth][len-1]
	pb.DepthToNodes[depth] = pb.DepthToNodes[depth][:len-1]
	return node
}
