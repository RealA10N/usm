package control_flow

type ControlFlowGraph struct {
	Graph

	// BasicBlockToNodes[i] contains the set of nodes in the original graph that
	// the i-th basic block contains.
	BasicBlockToNodes [][]uint

	// NodeToBasicBlock[i] contains the index of the basic block that the i-th
	// node in the original graph belongs to.
	NodeToBasicBlock []uint
}
