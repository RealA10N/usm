package gen

import (
	"alon.kr/x/list"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/parse"
)

type InstructionGenerator struct {
	ArgumentGenerator InstructionContextGenerator[parse.ArgumentNode, ArgumentInfo]
	TargetGenerator   FunctionContextGenerator[parse.TargetNode, *TargetInfo]
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
	ctx *InstructionGenerationContext,
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

func (g *InstructionGenerator) generateTargets(
	ctx *InstructionGenerationContext,
	node parse.InstructionNode,
) ([]*TargetInfo, core.ResultList) {
	targets := make([]*TargetInfo, len(node.Targets))
	results := core.ResultList{}

	for i, target := range node.Targets {
		targetInfo, curResults := g.TargetGenerator.Generate(
			ctx.FunctionGenerationContext,
			target,
		)

		if !curResults.IsEmpty() {
			results.Extend(&curResults)
			continue
		}

		if targetInfo.Register == nil {
			curResults := UndefinedRegisterResult(target.Register)
			results.Extend(&curResults)
		}

		targets[i] = targetInfo
	}

	if !results.IsEmpty() {
		return nil, results
	}

	return targets, core.ResultList{}
}

func (g *InstructionGenerator) generateLabels(
	ctx *InstructionGenerationContext,
	node parse.InstructionNode,
) ([]*LabelInfo, core.ResultList) {
	labels := make([]*LabelInfo, 0, len(node.Labels))
	for _, node := range node.Labels {
		name := NodeToSourceString(ctx.FileGenerationContext, node)
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

	v := node.View()
	instCtx := InstructionGenerationContext{
		FunctionGenerationContext: ctx,
		InstructionInfo:           NewEmptyInstructionInfo(&v),
	}

	instName := ViewToSourceString(ctx.FileGenerationContext, node.Operator)
	instDef, results := ctx.Instructions.GetInstructionDefinition(instName, node)

	arguments, curResults := g.generateArguments(&instCtx, node)
	results.Extend(&curResults)

	targets, curResults := g.generateTargets(&instCtx, node)
	results.Extend(&curResults)

	// Now it's time to check if we have any errors so far.
	if !results.IsEmpty() {
		return nil, results
	}

	instCtx.InstructionInfo.AppendTarget(targets...)
	instCtx.InstructionInfo.AppendArgument(arguments...)

	instruction, results := instDef.BuildInstruction(instCtx.InstructionInfo)
	if !results.IsEmpty() {
		return nil, results
	}

	instCtx.InstructionInfo.Instruction = instruction
	return instCtx.InstructionInfo, core.ResultList{}
}
