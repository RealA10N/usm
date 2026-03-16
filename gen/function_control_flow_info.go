package gen

import (
	"alon.kr/x/graph"
)

// FunctionControlFlowInfo holds the flattened basic-block list, the
// block-to-index mapping, and the control-flow graph for a function. It is
// derived purely from FunctionInfo and can be used by any pass that needs
// CFG traversal (e.g. constant propagation, liveness analysis) or block-index
// lookups (e.g. SSA construction).
type FunctionControlFlowInfo struct {
	// A linear representation of all basic blocks in the function.
	BasicBlocks []*BasicBlockInfo

	// Maps each basic block to its index in BasicBlocks.
	BasicBlocksToIndex map[*BasicBlockInfo]uint

	ControlFlowGraph *graph.Graph
}

// NewFunctionControlFlowInfo builds a FunctionControlFlowInfo from the given
// function by linearising its basic blocks and constructing the CFG from their
// forward edges.
func NewFunctionControlFlowInfo(function *FunctionInfo) FunctionControlFlowInfo {
	basicBlocks := function.CollectBasicBlocks()

	blockToIndex := make(map[*BasicBlockInfo]uint, len(basicBlocks))
	for i, b := range basicBlocks {
		blockToIndex[b] = uint(i)
	}

	forwardEdges := make([][]uint, len(basicBlocks))
	for i, b := range basicBlocks {
		edges := make([]uint, 0, len(b.ForwardEdges))
		for _, target := range b.ForwardEdges {
			edges = append(edges, blockToIndex[target])
		}
		forwardEdges[i] = edges
	}

	cfg := graph.NewGraph(forwardEdges)
	return FunctionControlFlowInfo{
		BasicBlocks:        basicBlocks,
		BasicBlocksToIndex: blockToIndex,
		ControlFlowGraph:   &cfg,
	}
}
