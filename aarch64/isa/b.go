package aarch64isa

import (
	"alon.kr/x/aarch64codegen/instructions"
	"alon.kr/x/list"
	aarch64codegen "alon.kr/x/usm/aarch64/codegen"
	aarch64translation "alon.kr/x/usm/aarch64/translation"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
)

type Branch struct {
	Target *gen.LabelInfo
}

func (b Branch) Operator() string {
	return "b"
}

func (b Branch) PossibleNextSteps() (gen.StepInfo, core.ResultList) {
	return gen.StepInfo{
		PossibleBranches: []*gen.LabelInfo{b.Target},
	}, core.ResultList{}
}

func (b Branch) Generate(
	ctx *aarch64codegen.InstructionCodegenContext,
) (instructions.Instruction, core.ResultList) {
	targetBasicBlock := b.Target.BasicBlock
	targetOffset := ctx.BasicBlockOffsets[targetBasicBlock]
	currentOffset := ctx.InstructionOffsetInFunction
	offset, err := aarch64translation.Uint64DiffToOffset26Align4(targetOffset, currentOffset)

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

	instruction := instructions.B(offset)
	return instruction, core.ResultList{}
}

type BranchDefinition struct{}

func (BranchDefinition) BuildInstruction(
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

	return Branch{Target: target}, core.ResultList{}
}

func NewBranchInstructionDefinition() gen.InstructionDefinition {
	return BranchDefinition{}
}
