package gen

import (
	"alon.kr/x/list"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/parse"
)

type FunctionGenerationScheme interface {
	NewJumpInstruction(label *LabelInfo) *InstructionInfo
}

type FunctionGenerator struct {
	FunctionGenerationScheme
	InstructionGenerator     FunctionContextGenerator[parse.InstructionNode, *InstructionInfo]
	ParameterGenerator       FunctionContextGenerator[parse.ParameterNode, *RegisterInfo]
	LabelDefinitionGenerator FunctionContextGenerator[parse.LabelNode, *LabelInfo]
	ReferencedTypeGenerator  FileContextGenerator[parse.TypeNode, ReferencedTypeInfo]
	TargetGenerator          FunctionContextGenerator[parse.TargetNode, *TargetInfo]
}

func NewFunctionGenerator() FileContextGenerator[parse.FunctionNode, *FunctionInfo] {
	return FileContextGenerator[parse.FunctionNode, *FunctionInfo](
		&FunctionGenerator{
			InstructionGenerator:     NewInstructionGenerator(),
			ParameterGenerator:       NewParameterGenerator(),
			LabelDefinitionGenerator: NewLabelDefinitionGenerator(),
			ReferencedTypeGenerator:  NewReferencedTypeGenerator(),
			TargetGenerator:          NewTargetGenerator(),
		},
	)
}

func (g *FunctionGenerator) createParameterRegisters(
	ctx *FunctionGenerationContext,
	parameters []parse.ParameterNode,
) (registers []*RegisterInfo, results core.ResultList) {
	registers = make([]*RegisterInfo, 0, len(parameters))

	for _, parameter := range parameters {
		register, curResults := g.ParameterGenerator.Generate(ctx, parameter)
		results.Extend(&curResults)
		registers = append(registers, register)
	}

	return registers, results
}

func (g *FunctionGenerator) generateTargets(
	ctx *FunctionGenerationContext,
	nodes []parse.TypeNode,
) ([]ReferencedTypeInfo, core.ResultList) {
	types := make([]ReferencedTypeInfo, 0, len(nodes))
	results := core.ResultList{}
	for _, node := range nodes {
		typeInfo, curResults := g.ReferencedTypeGenerator.Generate(ctx.FileGenerationContext, node)
		results.Extend(&curResults)
		types = append(types, typeInfo)
	}

	return types, results
}

// A helper struct to store information about labels in a function.
//
// This is used to store the mapping from instruction index to labels, and
// from label to instruction index, which are used when generating the basic
// blocks of the function.
type functionLabelData struct {
	InstructionIndexToLabels map[int][]*LabelInfo
	LabelToInstructionIndex  map[*LabelInfo]int
}

func (g *FunctionGenerator) collectLabelDefinitions(
	ctx *FunctionGenerationContext,
	instructions []parse.InstructionNode,
) (functionLabelData, core.ResultList) {
	results := core.ResultList{}
	indexToLabels := make(map[int][]*LabelInfo)
	labelToIndex := make(map[*LabelInfo]int)

	for i, instruction := range instructions {
		for _, label := range instruction.Labels {
			info, curResults := g.LabelDefinitionGenerator.Generate(ctx, label)
			results.Extend(&curResults)

			if curResults.IsEmpty() {
				labelToIndex[info] = i
				indexToLabels[i] = append(indexToLabels[i], info)
			}
		}
	}

	return functionLabelData{
		InstructionIndexToLabels: indexToLabels,
		LabelToInstructionIndex:  labelToIndex,
	}, results
}

// Before actually generating the instructions, we iterate over instruction and
// only collect information about target registers.
// Since USM requires all registers to be defined at least once with an explicit
// type, after collecting all register definitions the register manager should
// contain all registers with the required information.
func (g *FunctionGenerator) collectRegisterDefinitions(
	ctx *FunctionGenerationContext,
	instructions []parse.InstructionNode,
) (results core.ResultList) {
	for _, instruction := range instructions {
		for _, target := range instruction.Targets {
			_, curResults := g.TargetGenerator.Generate(ctx, target)
			results.Extend(&curResults)
		}
	}

	return results
}

func (g *FunctionGenerator) generateInstructions(
	ctx *FunctionGenerationContext,
	instNodes []parse.InstructionNode,
) ([]*InstructionInfo, core.ResultList) {
	g.collectRegisterDefinitions(ctx, instNodes)

	instructions := make([]*InstructionInfo, 0, len(instNodes))

	for _, instNode := range instNodes {
		inst, results := g.InstructionGenerator.Generate(ctx, instNode)
		if !results.IsEmpty() {
			// If encountered an error in the middle of the function, it might
			// effect the rest of the function (for example, a register might
			// not be defined correctly, which can cause other errors further
			// down the function). Thus, we return early.
			return nil, results
		}

		instructions = append(instructions, inst)
	}

	return instructions, core.ResultList{}
}

func (g *FunctionGenerator) getInstructionBranchingDestinations(
	info *InstructionInfo,
	labels functionLabelData,
) ([]int, core.ResultList) {
	steps, results := info.Instruction.PossibleNextSteps()
	if !results.IsEmpty() {
		return nil, results
	}

	destinationIndices := []int{}
	for _, label := range steps.PossibleBranches {
		destinationIndex := labels.LabelToInstructionIndex[label]
		destinationIndices = append(destinationIndices, destinationIndex)
	}

	return destinationIndices, core.ResultList{}
}

func (g *FunctionGenerator) getInstructionBranchingEdges(
	instructions []*InstructionInfo,
	labels functionLabelData,
) ([]bool, []bool, core.ResultList) {
	instructionCount := len(instructions)
	forwardEdges := make([]bool, instructionCount)
	backwardEdges := make([]bool, instructionCount)

	for i := 0; i < instructionCount; i++ {
		destinationIndices, results := g.getInstructionBranchingDestinations(
			instructions[i],
			labels,
		)

		if !results.IsEmpty() {
			return nil, nil, results
		}

		for _, j := range destinationIndices {
			forwardEdges[i] = true
			backwardEdges[j] = true
		}
	}

	return forwardEdges, backwardEdges, core.ResultList{}
}

// Creates new basic blocks from an instruction, and returns the *last* basic
// block in the created chain.
//
// If more than one basic block is created, all except the last one should not
// contain any instructions. Notice that this function only creates the basic
// blocks, but does not propagate them with instructions (even not with the
// provided instruction).
func createBasicBlocksFromLabels(
	ctx *FunctionGenerationContext,
	function *FunctionInfo,
	previousBlock *BasicBlockInfo,
	labels []*LabelInfo,
) (*BasicBlockInfo, core.ResultList) {

	if len(labels) == 0 {
		// Generate a label for the new basic block.
		label := ctx.Labels.GenerateLabel()
		results := ctx.Labels.NewLabel(label)
		if !results.IsEmpty() {
			return nil, results
		}
		labels = append(labels, label)
	}

	for _, label := range labels {
		newBasicBlock := NewEmptyBasicBlockInfo(function)
		newBasicBlock.SetLabel(label)

		if previousBlock != nil {
			previousBlock.AppendBasicBlock(newBasicBlock)
		}

		if function.EntryBlock == nil {
			function.EntryBlock = newBasicBlock
		}

		previousBlock = newBasicBlock
	}

	return previousBlock, core.ResultList{}
}

func (g *FunctionGenerator) generateBasicBlocks(
	ctx *FunctionGenerationContext,
	instructions []*InstructionInfo,
	function *FunctionInfo,
	labels functionLabelData,
) core.ResultList {
	instructionCount := len(instructions)

	forwardBranchingEdges, backwardBranchingEdges, results := g.getInstructionBranchingEdges(
		instructions,
		labels,
	)

	if !results.IsEmpty() {
		return results
	}

	entryBasicBlock, results := createBasicBlocksFromLabels(
		ctx,
		function,
		nil,
		labels.InstructionIndexToLabels[0],
	)

	if !results.IsEmpty() {
		return results
	}

	entryBasicBlock.AppendInstruction(instructions[0])

	currentBasicBlock := entryBasicBlock
	for i := 1; i < instructionCount; i++ {
		currentInstruction := instructions[i]
		previousInstructionBranches := forwardBranchingEdges[i-1]
		branchingToCurrentInstruction := backwardBranchingEdges[i]
		instructionLabels := labels.InstructionIndexToLabels[i]
		hasLabels := len(instructionLabels) > 0
		shouldStartNewBlock := previousInstructionBranches || branchingToCurrentInstruction || hasLabels

		if shouldStartNewBlock {
			currentBasicBlock, results = createBasicBlocksFromLabels(
				ctx,
				function,
				currentBasicBlock,
				instructionLabels,
			)
			if !results.IsEmpty() {
				return results
			}
		}

		currentBasicBlock.AppendInstruction(currentInstruction)
	}

	currentBasicBlock = entryBasicBlock
	for currentBasicBlock = entryBasicBlock; currentBasicBlock != nil; currentBasicBlock = currentBasicBlock.NextBlock {
		basicBlockLength := len(currentBasicBlock.Instructions)
		lastInstruction := currentBasicBlock.Instructions[basicBlockLength-1]

		steps, results := lastInstruction.Instruction.PossibleNextSteps()
		if !results.IsEmpty() {
			return results
		}

		if steps.PossibleContinue {
			nextBasicBlock := currentBasicBlock.NextBlock
			if nextBasicBlock == nil {
				return list.FromSingle(core.Result{
					{
						Type:     core.ErrorResult,
						Message:  "Unexpected instruction to end a function",
						Location: lastInstruction.Declaration,
					}, {
						Type:    core.HintResult,
						Message: "Perhaps you forgot a return instruction?",
					},
				})
			}

			currentBasicBlock.AppendForwardEdge(nextBasicBlock)
		}

		for _, label := range steps.PossibleBranches {
			branchToInstructionIndex := labels.LabelToInstructionIndex[label]
			branchToInstruction := instructions[branchToInstructionIndex]
			branchToBlock := branchToInstruction.BasicBlockInfo
			currentBasicBlock.AppendForwardEdge(branchToBlock)
		}
	}

	return core.ResultList{}
}

func (g *FunctionGenerator) Generate(
	ctx *FileGenerationContext,
	node parse.FunctionNode,
) (*FunctionInfo, core.ResultList) {
	var results core.ResultList
	funcCtx := ctx.NewFunctionGenerationContext()

	parameters, paramResults := g.createParameterRegisters(funcCtx, node.Signature.Parameters)
	results.Extend(&paramResults)

	targets, targetResults := g.generateTargets(funcCtx, node.Signature.Returns)
	results.Extend(&targetResults)

	labels, labelResults := g.collectLabelDefinitions(funcCtx, node.Instructions.Nodes)
	results.Extend(&labelResults)

	if !results.IsEmpty() {
		return nil, results
	}

	name := ViewToSourceString(funcCtx.FileGenerationContext, node.Signature.Identifier)

	function := &FunctionInfo{
		Name:        name,
		Declaration: &node.UnmanagedSourceView,
		EntryBlock:  nil, // will be defined later.
		Registers:   funcCtx.Registers,
		Labels:      funcCtx.Labels,
		Parameters:  parameters,
		Targets:     targets,
	}

	instructions, results := g.generateInstructions(funcCtx, node.Instructions.Nodes)
	if !results.IsEmpty() {
		return nil, results
	}

	if len(instructions) == 0 {
		v := node.View()
		return nil, list.FromSingle(core.Result{{
			Type:     core.ErrorResult,
			Message:  "Function must contain at least one instruction",
			Location: &v,
		}})
	}

	results = g.generateBasicBlocks(
		funcCtx,
		instructions,
		function,
		labels,
	)

	if !results.IsEmpty() {
		return nil, results
	}

	return function, core.ResultList{}
}
