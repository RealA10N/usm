package gen

import (
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/parse"
)

type FunctionInfo[InstT BaseInstruction] struct {
	Instructions []InstT
	Parameters   []*RegisterInfo
	// TODO: add targets
}

type FunctionGenerator[InstT BaseInstruction] struct {
	InstructionGenerator     Generator[InstT, parse.InstructionNode, InstT]
	ParameterGenerator       Generator[InstT, parse.ParameterNode, *RegisterInfo]
	LabelDefinitionGenerator LabelGenerator[InstT, parse.LabelNode, LabelInfo]
}

func NewFunctionGenerator[InstT BaseInstruction]() Generator[InstT, parse.FunctionNode, *FunctionInfo[InstT]] {
	return Generator[InstT, parse.FunctionNode, *FunctionInfo[InstT]](
		&FunctionGenerator[InstT]{
			InstructionGenerator:     NewInstructionGenerator[InstT](),
			ParameterGenerator:       NewParameterGenerator[InstT](),
			LabelDefinitionGenerator: NewLabelDefinitionGenerator[InstT](),
		},
	)
}

func (g *FunctionGenerator[InstT]) createParameterRegisters(
	ctx *GenerationContext[InstT],
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
	ctx *GenerationContext[InstT],
	instructions []parse.InstructionNode,
) (results core.ResultList) {

	labelCtx := LabelGenerationContext[InstT]{
		GenerationContext:       ctx,
		CurrentInstructionIndex: 0,
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
	ctx *GenerationContext[InstT],
	instNodes []parse.InstructionNode,
) ([]InstT, core.ResultList) {
	instructions := make([]InstT, 0, len(instNodes))

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
	ctx *GenerationContext[InstT],
	node parse.FunctionNode,
) (*FunctionInfo[InstT], core.ResultList) {
	var results core.ResultList
	parameters, paramResults := g.createParameterRegisters(ctx, node.Signature.Parameters)
	results.Extend(&paramResults)

	labelResults := g.collectLabelDefinitions(ctx, node.Instructions.Nodes)
	results.Extend(&labelResults)

	if !results.IsEmpty() {
		return nil, results
	}

	instructions, results := g.generateFunctionBody(ctx, node.Instructions.Nodes)
	if !results.IsEmpty() {
		return nil, results
	}

	return &FunctionInfo[InstT]{
		Instructions: instructions,
		Parameters:   parameters,
	}, core.ResultList{}
}
