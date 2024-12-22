package graph

const Unreachable = ^uint(0)

type ControlFlowGraph struct {
	Graph

	// BasicBlockToNodes[i] contains the set of nodes in the original graph that
	// the i-th basic block contains.
	BasicBlockToNodes [][]uint

	// NodeToBasicBlock[i] contains the index of the basic block that the i-th
	// node in the original graph belongs to.
	// NodeToBasicBlock[i] = Unreachable if the original node i is not included
	// in the control flow graph, i.e., it is unreachable from the entry node.
	NodeToBasicBlock []uint
}
