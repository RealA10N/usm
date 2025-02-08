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
	LabelDefinitionGenerator LabelContextGenerator[parse.LabelNode, *LabelInfo]
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

func (g *FunctionGenerator) createFunctionContext(
	ctx *FileGenerationContext,
) *FunctionGenerationContext {
	return &FunctionGenerationContext{
		FileGenerationContext: ctx,
		Registers:             ctx.RegisterManagerCreator(),
		Labels:                ctx.LabelManagerCreator(),
	}
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

func (g *FunctionGenerator) collectLabelDefinitions(
	ctx *FunctionGenerationContext,
	instructions []parse.InstructionNode,
) (map[*LabelInfo]uint, core.ResultList) {
	results := core.ResultList{}
	labelToInstructionIndex := make(map[*LabelInfo]uint)

	labelCtx := LabelGenerationContext{
		FunctionGenerationContext: ctx,
		CurrentInstructionIndex:   0,
	}

	for i, instruction := range instructions {
		for _, label := range instruction.Labels {
			info, curResults := g.LabelDefinitionGenerator.Generate(&labelCtx, label)
			labelToInstructionIndex[info] = uint(i)
			results.Extend(&curResults)
		}

		labelCtx.CurrentInstructionIndex++
	}

	return labelToInstructionIndex, results
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
	labelToInstructionIndex map[*LabelInfo]uint,
) ([]uint, core.ResultList) {
	steps, results := info.Instruction.PossibleNextSteps()
	if !results.IsEmpty() {
		return nil, results
	}

	destinationIndices := []uint{}
	for _, label := range steps.PossibleBranches {
		destinationIndex := labelToInstructionIndex[label]
		destinationIndices = append(destinationIndices, destinationIndex)
	}

	return destinationIndices, core.ResultList{}
}

func (g *FunctionGenerator) getInstructionBranchingEdges(
	instructions []*InstructionInfo,
	labelToInstructionIndex map[*LabelInfo]uint,
) ([]bool, []bool, core.ResultList) {
	instructionCount := uint(len(instructions))

	forwardEdges := make([]bool, instructionCount)
	backwardEdges := make([]bool, instructionCount)
	for i := uint(0); i < instructionCount; i++ {
		destinationIndices, results := g.getInstructionBranchingDestinations(
			instructions[i],
			labelToInstructionIndex,
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

func (g *FunctionGenerator) generateBasicBlocks(
	instructions []*InstructionInfo,
	function *FunctionInfo,
	labelToInstructionIndex map[*LabelInfo]uint,
) (*BasicBlockInfo, core.ResultList) {
	instructionCount := uint(len(instructions))

	forwardBranchingEdges, backwardBranchingEdges, results := g.getInstructionBranchingEdges(
		instructions,
		labelToInstructionIndex,
	)

	if !results.IsEmpty() {
		return nil, results
	}

	entryBasicBlock := NewEmptyBasicBlockInfo(function)
	currentBasicBlock := entryBasicBlock
	currentBasicBlock.AppendInstruction(instructions[0])

	for i := uint(1); i < instructionCount; i++ {
		currentInstruction := instructions[i]
		previousInstructionBranches := forwardBranchingEdges[i-1]
		branchingToCurrentInstruction := backwardBranchingEdges[i]
		shouldStartNewBlock := previousInstructionBranches || branchingToCurrentInstruction

		if shouldStartNewBlock {
			newBasicBlock := NewEmptyBasicBlockInfo(function)
			currentBasicBlock.AppendBasicBlock(newBasicBlock)
			currentBasicBlock = newBasicBlock
		}

		currentBasicBlock.AppendInstruction(currentInstruction)
	}

	currentBasicBlock = entryBasicBlock
	for currentBasicBlock = entryBasicBlock; currentBasicBlock != nil; currentBasicBlock = currentBasicBlock.NextBlock {
		basicBlockLength := len(currentBasicBlock.Instructions)
		lastInstruction := currentBasicBlock.Instructions[basicBlockLength-1]

		steps, results := lastInstruction.Instruction.PossibleNextSteps()
		if !results.IsEmpty() {
			return nil, results
		}

		if steps.PossibleContinue {
			nextBasicBlock := currentBasicBlock.NextBlock
			if nextBasicBlock == nil {
				return nil, list.FromSingle(core.Result{
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
			branchToInstructionIndex := labelToInstructionIndex[label]
			branchToInstruction := instructions[branchToInstructionIndex]
			branchToBlock := branchToInstruction.BasicBlockInfo
			currentBasicBlock.AppendForwardEdge(branchToBlock)
		}
	}

	return entryBasicBlock, core.ResultList{}
}

func (g *FunctionGenerator) Generate(
	ctx *FileGenerationContext,
	node parse.FunctionNode,
) (*FunctionInfo, core.ResultList) {
	var results core.ResultList
	funcCtx := g.createFunctionContext(ctx)

	parameters, paramResults := g.createParameterRegisters(funcCtx, node.Signature.Parameters)
	results.Extend(&paramResults)

	targets, targetResults := g.generateTargets(funcCtx, node.Signature.Returns)
	results.Extend(&targetResults)

	labelToInstructionIndex, labelResults := g.collectLabelDefinitions(funcCtx, node.Instructions.Nodes)
	results.Extend(&labelResults)

	if !results.IsEmpty() {
		return nil, results
	}

	name := viewToSourceString(funcCtx.FileGenerationContext, node.Signature.Identifier)

	function := &FunctionInfo{
		Name:       name,
		EntryBlock: nil, // will be defined later.
		Registers:  funcCtx.Registers,
		Labels:     funcCtx.Labels,
		Parameters: parameters,
		Targets:    targets,
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

	entryBlock, results := g.generateBasicBlocks(
		instructions,
		function,
		labelToInstructionIndex,
	)

	if !results.IsEmpty() {
		return nil, results
	}

	function.EntryBlock = entryBlock
	return function, core.ResultList{}
}
