package ssa

import (
	"fmt"

	"alon.kr/x/graph"
	"alon.kr/x/list"
	"alon.kr/x/set"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
)

type forwardingRegisterDescriptor struct {
	block   *gen.BasicBlockInfo
	renamed *gen.RegisterInfo
}

type phiInstructionDescriptor struct {
	// The new phi instruction.
	*gen.InstructionInfo

	// The original base register that this phi instruction is a definition for.
	base *gen.RegisterInfo

	// Pairs of (block, renamed register) that this phi instruction receives.
	// There should be one pair for each block that has a forward edge to the
	// block in which this phi instruction is inserted.
	forwards []forwardingRegisterDescriptor
}

func (p *phiInstructionDescriptor) AddForwardingRegister(
	block *gen.BasicBlockInfo,
	renamed *gen.RegisterInfo,
) {
	p.forwards = append(p.forwards, forwardingRegisterDescriptor{block, renamed})
}

// Move the forwarding registers to the phi instruction.
//
// We use this separation, and not directly modify the phi instruction, to not
// add arguments to the phi instruction before it was processed itself, which
// will result in some weird behavior.
func (p *phiInstructionDescriptor) CommitForwardingRegisters() core.ResultList {
	results := core.ResultList{}
	phiDefinition, ok := p.Definition.(PhiInstructionDefinition)

	if !ok {
		defName := p.Definition.Operator(p.InstructionInfo)
		return list.FromSingle(core.Result{
			{
				Type:     core.InternalErrorResult,
				Message:  "Expected phi instruction definition",
				Location: p.Declaration,
			},
			{
				Type:    core.HintResult,
				Message: fmt.Sprintf("Got \"%s\" instruction definition", defName),
			},
		})
	}

	for _, forward := range p.forwards {
		curResults := phiDefinition.AddForwardingRegister(
			p.InstructionInfo,
			forward.block,
			forward.renamed,
		)
		results.Extend(&curResults)
	}

	p.forwards = nil
	return results
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
	PhiInstructionsPerBlock [][]phiInstructionDescriptor

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
		PhiInstructionsPerBlock: make([][]phiInstructionDescriptor, len(basicBlocks)),
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
			descriptor := phiInstructionDescriptor{
				InstructionInfo: phi,
				base:            register,
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

func (i *FunctionSsaInfo) commitPhiInstructions() core.ResultList {
	results := core.ResultList{}

	for _, perBasicBlockPhiInstructions := range i.PhiInstructionsPerBlock {
		for _, phiDescriptor := range perBasicBlockPhiInstructions {
			curResults := phiDescriptor.CommitForwardingRegisters()
			results.Extend(&curResults)
		}
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
			// instructions in blocks that can be reached from the current block,
			// to forward them with the definition of the register if coming=
			// from the current block.
			forwardEdgeIndices := i.ControlFlowGraph.Nodes[basicBlockIndex].ForwardEdges
			for _, childIndex := range forwardEdgeIndices {
				for idx := range i.PhiInstructionsPerBlock[childIndex] {
					phiDescriptor := &i.PhiInstructionsPerBlock[childIndex][idx]
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
						phiDescriptor.AddForwardingRegister(basicBlock, renamed)
					}
				}
			}
		}
	}

	results := i.deleteBaseRegisters()
	if !results.IsEmpty() {
		return results
	}

	results = i.commitPhiInstructions()
	if !results.IsEmpty() {
		return results
	}

	return core.ResultList{}
}
