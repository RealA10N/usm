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
}

func NewBcond(condition immediates.Condition) gen.InstructionDefinition {
	return Bcond{
		Condition: condition,
	}
}

func (b Bcond) Target(
	info *gen.InstructionInfo,
) (*gen.LabelInfo, core.ResultList) {
	results := gen.AssertArgumentsExactly(info, 1)
	if !results.IsEmpty() {
		return nil, results
	}

	target, results := aarch64translation.ArgumentToLabelInfo(info.Arguments[0])
	if !results.IsEmpty() {
		return nil, results
	}

	return target, core.ResultList{}
}

func (b Bcond) Operator(*gen.InstructionInfo) string {
	return "b." + b.Condition.String()
}

func (b Bcond) PossibleNextSteps(info *gen.InstructionInfo) (gen.StepInfo, core.ResultList) {
	target, results := b.Target(info)
	if !results.IsEmpty() {
		return gen.StepInfo{}, results
	}

	return gen.StepInfo{
		PossibleBranches: []*gen.LabelInfo{target},
		PossibleContinue: true,
	}, core.ResultList{}
}

func (b Bcond) Codegen(
	ctx *aarch64codegen.InstructionCodegenContext,
) (instructions.Instruction, core.ResultList) {
	target, results := b.Target(ctx.InstructionInfo)
	if !results.IsEmpty() {
		return nil, results
	}

	targetBasicBlock := target.BasicBlock
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

func (b Bcond) Validate(
	info *gen.InstructionInfo,
) core.ResultList {
	results := core.ResultList{}

	_, curResults := b.Target(info)
	results.Extend(&curResults)

	curResults = gen.AssertTargetsExactly(info, 0)
	results.Extend(&curResults)

	if !results.IsEmpty() {
		return results
	}

	return core.ResultList{}
}
