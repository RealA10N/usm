package aarch64isa

import (
	"alon.kr/x/aarch64codegen/immediates"
	"alon.kr/x/aarch64codegen/instructions"
	"alon.kr/x/list"
	aarch64codegen "alon.kr/x/usm/aarch64/codegen"
	aarch64translation "alon.kr/x/usm/aarch64/translation"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
)

type Bcond struct {
	Condition immediates.Condition
	Target    *gen.LabelInfo
}

func (b Bcond) Operator() string {
	return "b." + b.Condition.String()
}

func (b Bcond) PossibleNextSteps() (gen.StepInfo, core.ResultList) {
	return gen.StepInfo{
		PossibleBranches: []*gen.LabelInfo{b.Target},
		PossibleContinue: true,
	}, core.ResultList{}
}

func (b Bcond) Generate(
	ctx *aarch64codegen.InstructionCodegenContext,
) (instructions.Instruction, core.ResultList) {
	targetBasicBlock := b.Target.BasicBlock
	targetOffset := ctx.BasicBlockOffsets[targetBasicBlock]
	currentOffset := ctx.InstructionOffsetInFunction
	offset, err := aarch64translation.Uint64DiffToOffset19Align4(targetOffset, currentOffset)

	if err != nil {
		return nil, list.FromSingle(core.Result{
			{
				Type:     core.ErrorResult,
				Message:  "Invalid branch offset (arbitrary large offsets are not yet supported)",
				Location: ctx.Declaration,
			},
			{
				Type:    core.DebugResult,
				Message: err.Error(),
			},
		})
	}

	instruction := instructions.BCOND(b.Condition, offset)
	return instruction, core.ResultList{}
}

type BcondDefinition struct {
	Condition immediates.Condition
}

func (d BcondDefinition) BuildInstruction(
	info *gen.InstructionInfo,
) (gen.BaseInstruction, core.ResultList) {
	results := core.ResultList{}

	curResults := aarch64translation.AssertArgumentsExactly(info, 1)
	results.Extend(&curResults)

	curResults = aarch64translation.AssertTargetsExactly(info, 0)
	results.Extend(&curResults)

	if !results.IsEmpty() {
		return nil, results
	}

	target, curResults := aarch64translation.ArgumentToLabelInfo(info.Arguments[0])
	results.Extend(&curResults)

	if !results.IsEmpty() {
		return nil, results
	}

	return Bcond{Condition: d.Condition, Target: target}, core.ResultList{}
}

func NewBcondInstructionDefinition(condition immediates.Condition) gen.InstructionDefinition {
	return BcondDefinition{Condition: condition}
}
