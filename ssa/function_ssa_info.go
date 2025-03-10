package ssa

import (
	"alon.kr/x/set"
	"alon.kr/x/usm/core"
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

	BaseRegisters []*gen.RegisterInfo

	// A mapping from (base) registers to their index in the registers slice.
	RegistersToIndex map[*gen.RegisterInfo]uint

	ControlFlowGraph   *graph.Graph
	DominatorJoinGraph *graph.DominatorJoinGraph
}

func NewFunctionSsaInfo(
	function *gen.FunctionInfo,
	ssaConstructionScheme SsaConstructionScheme,
) FunctionSsaInfo {
	basicBlocks := function.CollectBasicBlocks()
	basicBlockToIndex := createMappingToIndex(basicBlocks)
	forwardEdges := getBasicBlocksForwardEdges(basicBlocks, basicBlockToIndex)
	graph := graph.NewGraph(forwardEdges)
	dominatorJoinGraph := graph.DominatorJoinGraph(0)

	baseRegisters := function.Registers.GetAllRegisters()
	registersToIndex := createMappingToIndex(baseRegisters)

	return FunctionSsaInfo{
		FunctionInfo:            function,
		SsaConstructionScheme:   ssaConstructionScheme,
		BasicBlocks:             basicBlocks,
		BasicBlocksToIndex:      basicBlockToIndex,
		PhiInstructionsPerBlock: make([][]PhiInstructionDescriptor, len(basicBlocks)),
		BaseRegisters:           baseRegisters,
		RegistersToIndex:        registersToIndex,
		ControlFlowGraph:        &graph,
		DominatorJoinGraph:      &dominatorJoinGraph,
	}
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
	indices := make([]uint, 0, len(block.ForwardEdges))
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
	indices := make([]uint, 0, len(blocks))
	for block := range blocks {
		indices = append(indices, i.BasicBlocksToIndex[block])
	}
	return indices
}

func (i *FunctionSsaInfo) blockIndicesToBlockInfos(
	indices []uint,
) []*gen.BasicBlockInfo {
	blocks := make([]*gen.BasicBlockInfo, 0, len(indices))
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

func (i *FunctionSsaInfo) InsertPhiInstructions() core.ResultList {
	for _, register := range i.BaseRegisters {
		phiBlocks := i.getRegisterPhiInsertionPoints(register)
		for _, block := range phiBlocks {
			blockIndex := i.BasicBlocksToIndex[block]
			phi, results := i.SsaConstructionScheme.NewPhiInstruction(block, register)
			if !results.IsEmpty() {
				return results
			}
			descriptor := PhiInstructionDescriptor{
				PhiInstruction: phi,
				base:           register,
			}
			i.PhiInstructionsPerBlock[blockIndex] = append(
				i.PhiInstructionsPerBlock[blockIndex],
				descriptor)
		}
	}

	return core.ResultList{}
}

func (i *FunctionSsaInfo) deleteBaseRegisters() core.ResultList {
	results := core.ResultList{}
	for _, register := range i.BaseRegisters {
		curResults := i.Registers.DeleteRegister(register)
		results.Extend(&curResults)
	}
	return results
}

func (i *FunctionSsaInfo) RenameRegisters() core.ResultList {
	reachingSet := NewReachingDefinitionsSet(i)
	n := uint(len(i.BasicBlocks))

	for _, event := range i.DominatorJoinGraph.Dfs.Timeline {
		isPop := event >= n
		if isPop {
			reachingSet.popBlock()
		} else {
			// We have currently entered a new basic block in the dominator tree
			// traversal.
			reachingSet.pushBlock()

			basicBlockIndex := event
			basicBlock := i.BasicBlocks[basicBlockIndex]

			// Now, we let the specific implementation to handle the renaming
			// of the basic block registers (arguments and targets). We pass
			// the reaching definition set that we have built so far, and
			// the implementation should use it to query what is the live
			// definition of each register in the current basic block.
			results := i.SsaConstructionScheme.RenameBasicBlock(basicBlock, reachingSet)
			if !results.IsEmpty() {
				return results
			}

			// Now that the basic block has been renamed, we update all phi
			// instructions that are directly dominated by the current basic
			// block about the register that is live in the current basic block,
			// (which is possibly defined in the current basic block).
			basicBlockDominatorTreeNode := i.DominatorJoinGraph.DominatorTree.Nodes[basicBlockIndex]
			for _, childIndex := range basicBlockDominatorTreeNode.ForwardEdges {
				for _, phiDescriptor := range i.PhiInstructionsPerBlock[childIndex] {
					renamed := reachingSet.GetReachingDefinition(phiDescriptor.base)

					// If renamed == nil, it means that definition of the register
					// is undefined if reached from the current basic block.
					// Since we assume that the original representation is well
					// formed (no usage of undefined registers), we assume that
					// this means we can't reach the current basic block from
					// this child. (since otherwise on this path this register
					// value is undefined). So we just do not add the forwarding
					// register to the phi instruction.

					if renamed != nil {
						results := phiDescriptor.AddForwardingRegister(basicBlock, renamed)
						if !results.IsEmpty() {
							return results
						}
					}
				}
			}
		}
	}

	results := i.deleteBaseRegisters()
	if !results.IsEmpty() {
		return results
	}

	return core.ResultList{}
}
