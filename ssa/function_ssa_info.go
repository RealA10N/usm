package ssa

import (
	"alon.kr/x/set"
	"alon.kr/x/usm/gen"
	"alon.kr/x/usm/graph"
)

type FunctionSsaInfo struct {
	*gen.FunctionInfo

	SsaConstructionScheme

	// A slice that contains all blocks that are defined in the function.
	// It can be assumed that the length of the slice is >= 1, and that the
	// first block in the slice is the entry block.
	Blocks []*gen.BasicBlockInfo

	// A mapping between all basic blocks in the function and their index in the
	// Blocks slice.
	BasicBlocksToIndex map[*gen.BasicBlockInfo]uint

	ControlFlowGraph   *graph.Graph
	DominatorJoinGraph *graph.DominatorJoinGraph
}

func collectAllFunctionBasicBlocks(
	block *gen.BasicBlockInfo,
) []*gen.BasicBlockInfo {
	blocks := make([]*gen.BasicBlockInfo, 0)
	for block != nil {
		blocks = append(blocks, block)
		block = block.NextBlock
	}
	return blocks
}

func createBasicBlockToIndexMapping(
	blocks []*gen.BasicBlockInfo,
) map[*gen.BasicBlockInfo]uint {
	mapping := make(map[*gen.BasicBlockInfo]uint)
	for i, block := range blocks {
		mapping[block] = uint(i)
	}
	return mapping
}

func getSingleBasicBlockForwardEdges(
	block *gen.BasicBlockInfo,
	basicBlockToIndex map[*gen.BasicBlockInfo]uint,
) []uint {
	indices := make([]uint, len(block.ForwardEdges))
	for _, targetBlock := range block.ForwardEdges {
		indices = append(indices, basicBlockToIndex[targetBlock])
	}
	return indices
}

func getBasicBlocksForwardEdges(
	blocks []*gen.BasicBlockInfo,
	basicBlockToIndex map[*gen.BasicBlockInfo]uint,
) [][]uint {
	edges := make([][]uint, len(blocks))
	for i, block := range blocks {
		edges[i] = getSingleBasicBlockForwardEdges(block, basicBlockToIndex)
	}
	return edges
}

func NewFunctionSsaInfo(function *gen.FunctionInfo) FunctionSsaInfo {
	basicBlocks := collectAllFunctionBasicBlocks(function.EntryBlock)
	basicBlockToIndex := createBasicBlockToIndexMapping(basicBlocks)
	forwardEdges := getBasicBlocksForwardEdges(basicBlocks, basicBlockToIndex)
	graph := graph.NewGraph(forwardEdges)
	dominatorJoinGraph := graph.DominatorJoinGraph(0)

	return FunctionSsaInfo{
		FunctionInfo:       function,
		Blocks:             basicBlocks,
		BasicBlocksToIndex: basicBlockToIndex,
		ControlFlowGraph:   &graph,
		DominatorJoinGraph: &dominatorJoinGraph,
	}
}

// Returns all the basic blocks in which the provided register is defined.
func (i *FunctionSsaInfo) GetDefinitions(
	register *gen.RegisterInfo,
) set.Set[*gen.BasicBlockInfo] {
	blocks := set.New[*gen.BasicBlockInfo]()
	for _, instruction := range register.Definitions {
		blocks.Add(instruction.BasicBlockInfo)
	}

	return blocks
}

func (i *FunctionSsaInfo) BlockInfosToIndices(
	blocks set.Set[*gen.BasicBlockInfo],
) []uint {
	indices := make([]uint, len(blocks))
	for block := range blocks {
		indices = append(indices, i.BasicBlocksToIndex[block])
	}
	return indices
}

func (i *FunctionSsaInfo) blockIndicesToBlockInfos(
	indices []uint,
) []*gen.BasicBlockInfo {
	blocks := make([]*gen.BasicBlockInfo, len(indices))
	for _, index := range indices {
		blocks = append(blocks, i.Blocks[index])
	}
	return blocks
}

func (i *FunctionSsaInfo) GetPhiInsertionPoints(
	register *gen.RegisterInfo,
) []*gen.BasicBlockInfo {
	definitions := i.GetDefinitions(register)
	definitionsIndices := i.BlockInfosToIndices(definitions)
	phiBlocksIndices := i.DominatorJoinGraph.IteratedDominatorFrontier(definitionsIndices)
	return i.blockIndicesToBlockInfos(phiBlocksIndices)
}
