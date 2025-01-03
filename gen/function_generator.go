package gen

import (
	"alon.kr/x/list"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/graph"
	"alon.kr/x/usm/parse"
)

type FunctionGenerator struct {
	InstructionGenerator     FunctionContextGenerator[parse.InstructionNode, *InstructionInfo]
	ParameterGenerator       FunctionContextGenerator[parse.ParameterNode, *RegisterInfo]
	LabelDefinitionGenerator LabelContextGenerator[parse.LabelNode, *LabelInfo]
	ReferencedTypeGenerator  FileContextGenerator[parse.TypeNode, ReferencedTypeInfo]
}

func NewFunctionGenerator() FileContextGenerator[parse.FunctionNode, *FunctionInfo] {
	return FileContextGenerator[parse.FunctionNode, *FunctionInfo](
		&FunctionGenerator{
			InstructionGenerator:     NewInstructionGenerator(),
			ParameterGenerator:       NewParameterGenerator(),
			LabelDefinitionGenerator: NewLabelDefinitionGenerator(),
			ReferencedTypeGenerator:  NewReferencedTypeGenerator(),
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

func (g *FunctionGenerator) generateInstructions(
	ctx *FunctionGenerationContext,
	instNodes []parse.InstructionNode,
) ([]*InstructionInfo, core.ResultList) {
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

func (g *FunctionGenerator) generateInstructionsGraph(
	instructions []*InstructionInfo,
	labelToInstructionIndex map[*LabelInfo]uint,
) (graph.Graph, core.ResultList) {
	instructionCount := uint(len(instructions))
	instructionsGraph := graph.NewEmptyGraph(instructionCount)
	results := core.ResultList{}

	for i := uint(0); i < instructionCount; i++ {
		info := instructions[i]
		possibleNextSteps, curResults := info.Instruction.PossibleNextSteps()
		if !results.IsEmpty() {
			results.Extend(&curResults)
			continue
		}

		for _, nextStep := range possibleNextSteps {
			switch typedNextStep := nextStep.(type) {
			case ContinueToNextInstruction:
				if i+1 >= instructionCount {
					results.Append(core.Result{
						{
							Type:     core.ErrorResult,
							Message:  "Unexpected instruction to end a function",
							Location: info.Declaration,
						},
						{
							Type:    core.HintResult,
							Message: "Perhaps you forgot a return instruction?",
						},
					})
					continue
				}
				instructionsGraph.AddEdge(i, i+1)

			case JumpToLabel:
				j := labelToInstructionIndex[typedNextStep.Label]
				instructionsGraph.AddEdge(i, j)

			case ReturnFromFunction:
				// Don't add an edge.

			default:
				// notest
				results.Append(core.Result{{
					Type:     core.InternalErrorResult,
					Message:  "Unknown next step type",
					Location: info.Declaration,
				}})
			}
		}
	}

	if !results.IsEmpty() {
		return graph.Graph{}, results
	}

	return instructionsGraph, core.ResultList{}
}

func (g *FunctionGenerator) generateBasicBlocks(
	cfg graph.ControlFlowGraph,
	instructions []*InstructionInfo,
	function *FunctionInfo,
) (blocks []*BasicBlockInfo, results core.ResultList) {
	blocksCount := cfg.Size()
	blocks = make([]*BasicBlockInfo, blocksCount)

	// first, initialize ("malloc") blocks so we can take references to them.
	// on the way, also compute and fill any trivial fields that do not require
	// references to other blocks.
	for i := uint(0); i < blocksCount; i++ {
		blockInstructionIndices := cfg.BasicBlockToNodes[i]
		blocks[i] = NewEmptyBasicBlockInfo(function)

		for _, instructionIndex := range blockInstructionIndices {
			blocks[i].AppendInstruction(instructions[instructionIndex])
		}
	}

	// now fill in the missing edges fields.
	for i := uint(0); i < blocksCount; i++ {
		node := cfg.Nodes[i]
		for _, j := range node.ForwardEdges {
			blocks[i].ForwardEdges = append(blocks[i].ForwardEdges, blocks[j])
		}

		for _, j := range node.BackwardEdges {
			blocks[i].BackwardEdges = append(blocks[i].BackwardEdges, blocks[j])
		}
	}

	// finally, fill in the NextBlock field.
	// We do it this way, to guarantee that the order of the blocks in the
	// `NextBlock` list matches the order of them in the source code.
	for i := uint(0); i < blocksCount; i++ {
		lastBlockInstructionIndex := cfg.BasicBlockToNodes[i][len(cfg.BasicBlockToNodes[i])-1]
		firstNextBlockInstructionIndex := lastBlockInstructionIndex + 1
		if firstNextBlockInstructionIndex >= uint(len(instructions)) {
			continue // This is the last block in the function.
		}

		nextBlockIndex := cfg.NodeToBasicBlock[firstNextBlockInstructionIndex]
		blocks[i].NextBlock = blocks[nextBlockIndex]
	}

	return blocks, core.ResultList{}
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
		Parameters: parameters,
		Targets:    targets,
	}

	instructions, results := g.generateInstructions(funcCtx, node.Instructions.Nodes)
	if !results.IsEmpty() {
		return nil, results
	}

	graph, results := g.generateInstructionsGraph(instructions, labelToInstructionIndex)
	if !results.IsEmpty() {
		return nil, results
	}

	if graph.Size() == 0 {
		v := node.View()
		return nil, list.FromSingle(core.Result{{
			Type:     core.ErrorResult,
			Message:  "Function must contain at least one instruction",
			Location: &v,
		}})
	}

	cfg := graph.ControlFlowGraph(0)

	blocks, results := g.generateBasicBlocks(cfg, instructions, function)
	if !results.IsEmpty() {
		return nil, results
	}

	function.EntryBlock = blocks[0]
	return function, core.ResultList{}
}
