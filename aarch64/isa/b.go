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

func NewBranch() gen.InstructionDefinition {
	return Branch{}
}

func (i Branch) Operator(*gen.InstructionInfo) string {
	return "b"
}

func (i Branch) Target(
	info *gen.InstructionInfo,
) (*gen.LabelInfo, core.ResultList) {
	results := gen.AssertArgumentsExactly(info, 1)
	if !results.IsEmpty() {
		return nil, results
	}

	label, results := aarch64translation.ArgumentToLabelInfo(info.Arguments[0])
	if !results.IsEmpty() {
		return nil, results
	}

	return label, core.ResultList{}
}

func (i Branch) PossibleNextSteps(
	info *gen.InstructionInfo,
) (gen.StepInfo, core.ResultList) {
	target, results := i.Target(info)
	return gen.StepInfo{
		PossibleBranches: []*gen.LabelInfo{target},
	}, results
}

type branchValidationArtifacts struct {
	Target *gen.LabelInfo
}

func (i Branch) internalValidate(
	info *gen.InstructionInfo,
) (*branchValidationArtifacts, core.ResultList) {
	results := core.ResultList{}

	target, curResults := i.Target(info)
	results.Extend(&curResults)

	curResults = gen.AssertTargetsExactly(info, 0)
	results.Extend(&curResults)

	if !results.IsEmpty() {
		return nil, results
	}

	artifacts := &branchValidationArtifacts{
		Target: target,
	}

	return artifacts, core.ResultList{}
}

func (i Branch) Codegen(
	ctx *aarch64codegen.InstructionCodegenContext,
) (instructions.Instruction, core.ResultList) {
	info := ctx.InstructionInfo

	artifacts, results := i.internalValidate(info)
	if !results.IsEmpty() {
		return nil, results
	}

	target := artifacts.Target
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

func (i Branch) Validate(
	info *gen.InstructionInfo,
) core.ResultList {
	_, results := i.internalValidate(info)
	return results
}
