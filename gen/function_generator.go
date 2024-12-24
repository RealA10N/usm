package gen

import (
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/parse"
)

type FunctionGenerator[InstT BaseInstruction] struct {
	InstructionGenerator     FunctionContextGenerator[InstT, parse.InstructionNode, *InstructionInfo[InstT]]
	ParameterGenerator       FunctionContextGenerator[InstT, parse.ParameterNode, *RegisterInfo]
	LabelDefinitionGenerator LabelContextGenerator[InstT, parse.LabelNode, LabelInfo]
}

func NewFunctionGenerator[InstT BaseInstruction]() FileContextGenerator[InstT, parse.FunctionNode, *FunctionInfo[InstT]] {
	return FileContextGenerator[InstT, parse.FunctionNode, *FunctionInfo[InstT]](
		&FunctionGenerator[InstT]{
			InstructionGenerator:     NewInstructionGenerator[InstT](),
			ParameterGenerator:       NewParameterGenerator[InstT](),
			LabelDefinitionGenerator: NewLabelDefinitionGenerator[InstT](),
		},
	)
}

func (g *FunctionGenerator[InstT]) createFunctionContext(
	ctx *FileGenerationContext[InstT],
) *FunctionGenerationContext[InstT] {
	return &FunctionGenerationContext[InstT]{
		FileGenerationContext: ctx,
		Registers:             ctx.RegisterManagerCreator(),
		Labels:                ctx.LabelManagerCreator(),
	}
}

func (g *FunctionGenerator[InstT]) createParameterRegisters(
	ctx *FunctionGenerationContext[InstT],
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

func (g *FunctionGenerator[InstT]) collectLabelDefinitions(
	ctx *FunctionGenerationContext[InstT],
	instructions []parse.InstructionNode,
) (results core.ResultList) {

	labelCtx := LabelGenerationContext[InstT]{
		FunctionGenerationContext: ctx,
		CurrentInstructionIndex:   0,
	}

	for _, instruction := range instructions {
		for _, label := range instruction.Labels {
			_, curResults := g.LabelDefinitionGenerator.Generate(&labelCtx, label)
			results.Extend(&curResults)
		}

		labelCtx.CurrentInstructionIndex++
	}

	return results
}

func (g *FunctionGenerator[InstT]) generateFunctionBody(
	ctx *FunctionGenerationContext[InstT],
	instNodes []parse.InstructionNode,
) ([]*InstructionInfo[InstT], core.ResultList) {
	instructions := make([]*InstructionInfo[InstT], 0, len(instNodes))

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

func (g *FunctionGenerator[InstT]) Generate(
	ctx *FileGenerationContext[InstT],
	node parse.FunctionNode,
) (*FunctionInfo[InstT], core.ResultList) {
	var results core.ResultList
	funcCtx := g.createFunctionContext(ctx)

	parameters, paramResults := g.createParameterRegisters(funcCtx, node.Signature.Parameters)
	results.Extend(&paramResults)

	labelResults := g.collectLabelDefinitions(funcCtx, node.Instructions.Nodes)
	results.Extend(&labelResults)

	if !results.IsEmpty() {
		return nil, results
	}

	instructions, results := g.generateFunctionBody(funcCtx, node.Instructions.Nodes)
	if !results.IsEmpty() {
		return nil, results
	}

	return &FunctionInfo[InstT]{
		Instructions: instructions,
		Parameters:   parameters,
	}, core.ResultList{}
}
