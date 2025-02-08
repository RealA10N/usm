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
		v := target.View()
		targetInfo, curResults := g.TargetGenerator.Generate(
			ctx.FunctionGenerationContext,
			target,
		)
		results.Extend(&curResults)

		if targetInfo.Register == nil {
			// TODO: improve error message
			results.Append(core.Result{
				{
					Type:     core.ErrorResult,
					Message:  "Undefined or untyped register",
					Location: &v,
				},
				{
					Type:    core.HintResult,
					Message: "A register must be defined with an explicit type at least once",
				},
			})
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

	instName := viewToSourceString(ctx.FileGenerationContext, node.Operator)
	instDef, results := ctx.Instructions.GetInstructionDefinition(instName, node)

	arguments, curResults := g.generateArguments(&instCtx, node)
	results.Extend(&curResults)

	targets, curResults := g.generateTargets(&instCtx, node)
	results.Extend(&curResults)

	labels, curResults := g.generateLabels(&instCtx, node)
	results.Extend(&curResults)

	// Now it's time to check if we have any errors so far.
	if !results.IsEmpty() {
		return nil, results
	}

	instCtx.InstructionInfo.AppendTarget(targets...)
	instCtx.InstructionInfo.AppendArgument(arguments...)
	instCtx.InstructionInfo.AppendLabels(labels...)

	instruction, results := instDef.BuildInstruction(instCtx.InstructionInfo)
	if !results.IsEmpty() {
		return nil, results
	}

	instCtx.InstructionInfo.Instruction = instruction
	return instCtx.InstructionInfo, core.ResultList{}
}
