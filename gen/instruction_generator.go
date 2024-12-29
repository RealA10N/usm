package gen

import (
	"alon.kr/x/list"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/parse"
)

type InstructionGenerator struct {
	ArgumentGenerator FunctionContextGenerator[parse.ArgumentNode, ArgumentInfo]
	TargetGenerator   FunctionContextGenerator[parse.TargetNode, registerPartialInfo]
}

func NewInstructionGenerator() FunctionContextGenerator[
	parse.InstructionNode,
	*InstructionInfo,
] {
	return FunctionContextGenerator[
		parse.InstructionNode,
		*InstructionInfo,
	](
		&InstructionGenerator{
			ArgumentGenerator: NewArgumentGenerator(),
			TargetGenerator:   NewTargetGenerator(),
		},
	)
}

func (g *InstructionGenerator) generateArguments(
	ctx *FunctionGenerationContext,
	node parse.InstructionNode,
) ([]ArgumentInfo, core.ResultList) {
	arguments := make([]ArgumentInfo, len(node.Arguments))
	results := core.ResultList{}

	// Different arguments should not effect one another.
	// Thus, we just collect all of the errors along the way, and return
	// them in one chunk.
	for i, argument := range node.Arguments {
		argInfo, curResults := g.ArgumentGenerator.Generate(ctx, argument)
		results.Extend(&curResults)
		arguments[i] = argInfo
	}

	return arguments, results
}

func (g *InstructionGenerator) generatePartialTargetsInfo(
	ctx *FunctionGenerationContext,
	node parse.InstructionNode,
) ([]registerPartialInfo, core.ResultList) {
	targets := make([]registerPartialInfo, len(node.Targets))
	results := core.ResultList{}

	// Different targets should not effect one another.
	// Thus, we just collect all of the errors along the way, and return
	// them in one chunk.
	for i, target := range node.Targets {
		typeInfo, curResults := g.TargetGenerator.Generate(ctx, target)
		results.Extend(&curResults)
		targets[i] = typeInfo
	}

	return targets, results
}

func (g *InstructionGenerator) generateLabels(
	ctx *FunctionGenerationContext,
	node parse.InstructionNode,
) ([]*LabelInfo, core.ResultList) {
	labels := make([]*LabelInfo, 0, len(node.Labels))
	for _, node := range node.Labels {
		name := nodeToSourceString(ctx.FileGenerationContext, node)
		label := ctx.Labels.GetLabel(name)
		if label == nil {
			v := node.View()
			return nil, list.FromSingle(core.Result{{
				Type:     core.ErrorResult,
				Message:  "Label does not exist",
				Location: &v,
			}})
		}

		labels = append(labels, label)
	}

	return labels, core.ResultList{}
}

func partialTargetsToTypes(targets []registerPartialInfo) []*ReferencedTypeInfo {
	types := make([]*ReferencedTypeInfo, len(targets))
	for i, target := range targets {
		types[i] = target.Type
	}
	return types
}

func argumentsToTypes(arguments []ArgumentInfo) []*ReferencedTypeInfo {
	types := make([]*ReferencedTypeInfo, len(arguments))
	for i, arg := range arguments {
		types[i] = arg.GetType()
	}
	return types
}

func (g *InstructionGenerator) getTargetRegister(
	ctx *FunctionGenerationContext,
	node parse.TargetNode,
	targetType ReferencedTypeInfo,
) (*RegisterInfo, core.Result) {
	registerName := nodeToSourceString(ctx.FileGenerationContext, node.Register)
	registerInfo := ctx.Registers.GetRegister(registerName)
	nodeView := node.View()

	if registerInfo == nil {
		// register is defined here; we should create the register and define
		// it's type.
		newRegisterInfo := &RegisterInfo{
			Name:        registerName,
			Type:        targetType,
			Declaration: nodeView,
		}

		return newRegisterInfo, ctx.Registers.NewRegister(newRegisterInfo)
	}

	// register is already defined
	if !registerInfo.Type.Equal(targetType) {
		// notest: sanity check only
		return nil, core.Result{{
			Type:     core.InternalErrorResult,
			Message:  "Internal register type mismatch",
			Location: &nodeView,
		}}
	}

	return registerInfo, nil
}

// Registers can be defined by being a target of an instruction.
// After we have determined 100% of the instruction targets types (either
// if they were explicitly declared or not), we call this function with the
// target types, and here we iterate over all target types and define missing
// registers.
//
// This also returns the full list of register targets for the provided
// instruction.
func (g *InstructionGenerator) defineAndGetTargetRegisters(
	ctx *FunctionGenerationContext,
	node parse.InstructionNode,
	targetTypes []ReferencedTypeInfo,
	instructionInfo *InstructionInfo,
) ([]*RegisterArgumentInfo, core.ResultList) {
	if len(node.Targets) != len(targetTypes) {
		// notest: sanity check: ensure lengths match.
		v := node.View()
		return nil, list.FromSingle(core.Result{{
			Type:     core.InternalErrorResult,
			Message:  "Targets length mismatch",
			Location: &v,
		}})
	}

	registers := make([]*RegisterArgumentInfo, len(node.Targets))
	results := core.ResultList{}
	for i, target := range node.Targets {
		// register errors should not effect one another, so we collect them.
		registerInfo, result := g.getTargetRegister(
			ctx,
			target,
			targetTypes[i],
		)

		if result != nil {
			results.Append(result)
		}

		if registerInfo == nil {
			// notest: sanity check, should not happen.
			v := target.View()
			results.Append(core.Result{{
				Type:     core.InternalErrorResult,
				Message:  "Unexpected nil register",
				Location: &v,
			}})
		}

		registerArgument := &RegisterArgumentInfo{
			Register:    registerInfo,
			declaration: node.View(),
		}

		registerInfo.AddDefinition(instructionInfo)
		registers[i] = registerArgument
	}

	return registers, results
}

// Convert an instruction parsed node into an instruction that is in the
// instruction set.
// If new registers are defined in the instruction (by assigning values to
// instruction targets), the register is created and added to the generation
// context.
func (g *InstructionGenerator) Generate(
	ctx *FunctionGenerationContext,
	node parse.InstructionNode,
) (*InstructionInfo, core.ResultList) {
	// We start generating the instruction, by getting the definition interface,
	// and processing the targets and arguments. We accumulate the results,
	// since those processes do not effect each other.

	instName := viewToSourceString(ctx.FileGenerationContext, node.Operator)
	instDef, results := ctx.Instructions.GetInstructionDefinition(instName, node)

	arguments, curResults := g.generateArguments(ctx, node)
	results.Extend(&curResults)

	partialTargets, curResults := g.generatePartialTargetsInfo(ctx, node)
	results.Extend(&curResults)

	labels, curResults := g.generateLabels(ctx, node)

	// Now it's time to check if we have any errors so far.
	if !results.IsEmpty() {
		return nil, results
	}

	explicitTargetTypes := partialTargetsToTypes(partialTargets)
	argumentTypes := argumentsToTypes(arguments)
	targetTypes, results := instDef.InferTargetTypes(ctx, explicitTargetTypes, argumentTypes)
	// TODO: validate that the returned target types matches expected constraints.

	if !results.IsEmpty() {
		return nil, results
	}

	v := node.View()
	info := &InstructionInfo{
		Instruction: nil, // will after BuildInstruction is called.
		Arguments:   arguments,
		Labels:      labels,
		Declaration: &v,
	}

	targets, results := g.defineAndGetTargetRegisters(
		ctx,
		node,
		targetTypes,
		info,
	)

	if !results.IsEmpty() {
		return nil, results
	}
	info.Targets = targets

	instruction, results := instDef.BuildInstruction(info)
	if !results.IsEmpty() {
		return nil, results
	}

	info.Instruction = instruction
	return info, core.ResultList{}
}
