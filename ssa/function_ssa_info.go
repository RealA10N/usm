package ssa

import (
	"alon.kr/x/set"
	"alon.kr/x/usm/gen"
	"alon.kr/x/usm/graph"
)

type FunctionSsaInfo struct {
	*gen.FunctionInfo

	SsaConstructionScheme SsaConstructionScheme

	// A linear representation of all basic blocks in the function.
	BasicBlocks []*gen.BasicBlockInfo

	// A mapping between all basic blocks in the function and their index in the
	// Blocks slice.
	BasicBlocksToIndex map[*gen.BasicBlockInfo]uint

	// A mapping from (base) registers to their index in the registers slice.
	RegistersToIndex map[*gen.RegisterInfo]uint

	ControlFlowGraph   *graph.Graph
	DominatorJoinGraph *graph.DominatorJoinGraph
}

func collectBasicBlocks(block *gen.BasicBlockInfo) []*gen.BasicBlockInfo {
	// TODO: this slice actually exists in the previous step in the compilation,
	// in the `gen.FunctionGenerator`. The current implementation creates the
	// array again instead of just passing it through so the implementation is
	// more complete and independent. However, if it is still the case we should
	// find a way to pass the array through as an optimization.

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
	basicBlocks := collectBasicBlocks(function.EntryBlock)
	basicBlockToIndex := createBasicBlockToIndexMapping(basicBlocks)
	forwardEdges := getBasicBlocksForwardEdges(basicBlocks, basicBlockToIndex)
	graph := graph.NewGraph(forwardEdges)
	dominatorJoinGraph := graph.DominatorJoinGraph(0)

	return FunctionSsaInfo{
		FunctionInfo:       function,
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
		blocks = append(blocks, i.BasicBlocks[index])
	}
	return blocks
}

func (i *FunctionSsaInfo) getRegisterPhiInsertionPoints(
	register *gen.RegisterInfo,
) []*gen.BasicBlockInfo {
	definitions := i.GetDefinitions(register)
	definitionsIndices := i.BlockInfosToIndices(definitions)
	phiBlocksIndices := i.DominatorJoinGraph.IteratedDominatorFrontier(definitionsIndices)
	return i.blockIndicesToBlockInfos(phiBlocksIndices)
}

func (i *FunctionSsaInfo) InsertPhiInstructions() {
	for _, register := range i.Registers {
		phiBlocks := i.getRegisterPhiInsertionPoints(register)
		for _, block := range phiBlocks {
			i.SsaConstructionScheme.NewPhiInstruction(block, register)
		}
	}
}

func (i *FunctionSsaInfo) RenameRegisters() {
	reachingSet := NewReachingDefinitionsSet(i)
	n := uint(len(i.BasicBlocks))

	for _, event := range i.DominatorJoinGraph.Dfs.Timeline {
		isPop := event >= n
		if isPop {
			reachingSet.popBlock()
		} else {
			reachingSet.pushBlock()
			basicBlockIndex := event
			basicBlock := i.BasicBlocks[basicBlockIndex]
			i.SsaConstructionScheme.RenameBasicBlock(basicBlock, reachingSet)
		}
	}
}
