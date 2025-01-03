package ssa

import (
	"alon.kr/x/set"
	"alon.kr/x/usm/gen"
	"alon.kr/x/usm/graph"
)

type PhiInstructionDescriptor struct {
	PhiInstruction
	base *gen.RegisterInfo
}
type FunctionSsaInfo struct {
	*gen.FunctionInfo

	SsaConstructionScheme SsaConstructionScheme

	// A linear representation of all basic blocks in the function.
	BasicBlocks []*gen.BasicBlockInfo

	// A mapping between all basic blocks in the function and their index in the
	// Blocks slice.
	BasicBlocksToIndex map[*gen.BasicBlockInfo]uint

	// A mapping between all basic blocks and the phi instructions that they
	// define in their entry.
	//
	// Initially, this slice is empty for each block, and each new phi instruction
	// which we create is inserted to to corresponding block's slice.
	PhiInstructionsPerBlock [][]PhiInstructionDescriptor

	// A mapping from (base) registers to their index in the registers slice.
	RegistersToIndex map[*gen.RegisterInfo]uint

	ControlFlowGraph   *graph.Graph
	DominatorJoinGraph *graph.DominatorJoinGraph
}

func NewFunctionSsaInfo(function *gen.FunctionInfo) FunctionSsaInfo {
	basicBlocks := collectBasicBlocks(function.EntryBlock)
	basicBlockToIndex := createMappingToIndex(basicBlocks)
	forwardEdges := getBasicBlocksForwardEdges(basicBlocks, basicBlockToIndex)
	graph := graph.NewGraph(forwardEdges)
	dominatorJoinGraph := graph.DominatorJoinGraph(0)
	registersToIndex := createMappingToIndex(function.Registers)

	return FunctionSsaInfo{
		FunctionInfo:            function,
		BasicBlocks:             basicBlocks,
		BasicBlocksToIndex:      basicBlockToIndex,
		PhiInstructionsPerBlock: make([][]PhiInstructionDescriptor, len(basicBlocks)),
		RegistersToIndex:        registersToIndex,
		ControlFlowGraph:        &graph,
		DominatorJoinGraph:      &dominatorJoinGraph,
	}
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

func createMappingToIndex[T comparable](
	slice []T,
) map[T]uint {
	mapping := make(map[T]uint)
	for i, element := range slice {
		mapping[element] = uint(i)
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

// Returns all the basic blocks in which the provided register is defined.
func (i *FunctionSsaInfo) getDefinitions(
	register *gen.RegisterInfo,
) set.Set[*gen.BasicBlockInfo] {
	blocks := set.New[*gen.BasicBlockInfo]()
	for _, instruction := range register.Definitions {
		blocks.Add(instruction.BasicBlockInfo)
	}

	return blocks
}

func (i *FunctionSsaInfo) blockInfosToIndices(
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
	definitions := i.getDefinitions(register)
	definitionsIndices := i.blockInfosToIndices(definitions)
	phiBlocksIndices := i.DominatorJoinGraph.IteratedDominatorFrontier(definitionsIndices)
	return i.blockIndicesToBlockInfos(phiBlocksIndices)
}

func (i *FunctionSsaInfo) InsertPhiInstructions() {
	for _, register := range i.Registers {
		phiBlocks := i.getRegisterPhiInsertionPoints(register)
		for _, block := range phiBlocks {
			blockIndex := i.BasicBlocksToIndex[block]
			phi := i.SsaConstructionScheme.NewPhiInstruction(block, register)
			descriptor := PhiInstructionDescriptor{
				PhiInstruction: phi,
				base:           register,
			}
			i.PhiInstructionsPerBlock[blockIndex] = append(
				i.PhiInstructionsPerBlock[blockIndex],
				descriptor)
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

			basicBlockDominatorTreeNode := i.DominatorJoinGraph.DominatorTree.Nodes[basicBlockIndex]
			for _, childIndex := range basicBlockDominatorTreeNode.ForwardEdges {
				childBlock := i.BasicBlocks[childIndex]
				for _, phiDescriptor := range i.PhiInstructionsPerBlock[childIndex] {
					renamed := reachingSet.GetReachingDefinition(phiDescriptor.base)
					phiDescriptor.AddForwardingRegister(childBlock, renamed)
				}
			}
		}
	}
}