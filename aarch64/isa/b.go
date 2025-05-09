package aarch64isa

import (
	"alon.kr/x/aarch64codegen/instructions"
	"alon.kr/x/list"
	aarch64codegen "alon.kr/x/usm/aarch64/codegen"
	aarch64translation "alon.kr/x/usm/aarch64/translation"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
)

type Branch struct{}

func NewBranch() Branch {
	return Branch{}
}

func (b Branch) Operator() string {
	return "b"
}

func (b Branch) Target(
	info *gen.InstructionInfo,
) (*gen.LabelInfo, core.ResultList) {
	results := aarch64translation.AssertArgumentsExactly(info, 1)
	if !results.IsEmpty() {
		return nil, results
	}

	label, results := aarch64translation.ArgumentToLabelInfo(info.Arguments[0])
	if !results.IsEmpty() {
		return nil, results
	}

	return label, core.ResultList{}
}

func (b Branch) PossibleNextSteps(
	info *gen.InstructionInfo,
) (gen.StepInfo, core.ResultList) {
	target, results := b.Target(info)
	return gen.StepInfo{
		PossibleBranches: []*gen.LabelInfo{target},
	}, results
}

func (b Branch) Codegen(
	ctx *aarch64codegen.InstructionCodegenContext,
) (instructions.Instruction, core.ResultList) {
	info := ctx.InstructionInfo

	results := core.ResultList{}
	curResults := aarch64translation.AssertArgumentsExactly(info, 1)
	results.Extend(&curResults)

	curResults = aarch64translation.AssertTargetsExactly(info, 0)
	results.Extend(&curResults)

	if !results.IsEmpty() {
		return nil, results
	}

	target, results := b.Target(info)
	if !results.IsEmpty() {
		return nil, results
	}

	targetBasicBlock := target.BasicBlock
	targetOffset := ctx.BasicBlockOffsets[targetBasicBlock]
	currentOffset := ctx.InstructionOffsetInFunction
	offset, err := aarch64translation.Uint64DiffToOffset26Align4(targetOffset, currentOffset)

	if err != nil {
		return nil, list.FromSingle(core.Result{
			{
				Type:     core.ErrorResult,
				Message:  "Branch offset too large",
				Location: ctx.Declaration,
			},
			{
				Type:    core.DebugResult,
				Message: err.Error(),
			},
		})
	}

	inst := instructions.NewBranch(offset)
	return inst, core.ResultList{}
}
