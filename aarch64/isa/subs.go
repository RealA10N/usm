package aarch64isa

import (
	"alon.kr/x/aarch64codegen/instructions"
	"alon.kr/x/list"
	aarch64codegen "alon.kr/x/usm/aarch64/codegen"
	aarch64translation "alon.kr/x/usm/aarch64/translation"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
)

type BaseSubs struct {
	NonBranchingInstruction
}

func (BaseSubs) Operator() string {
	return "subs"
}

type SubsReg struct {
	BaseSubs
	instructions.SubShiftedRegister
}

func (i SubsReg) Generate(
	*aarch64codegen.InstructionCodegenContext,
) (instructions.Instruction, core.ResultList) {
	return i, core.ResultList{}
}

type SubsImm struct {
	BaseSubs
	instructions.SubImmediate
}

func (i SubsImm) Generate(
	*aarch64codegen.InstructionCodegenContext,
) (instructions.Instruction, core.ResultList) {
	return i, core.ResultList{}
}

type SubsDefinition struct{}

func (d SubsDefinition) buildRegisterVariant(
	info *gen.InstructionInfo,
) (gen.BaseInstruction, core.ResultList) {
	Xd, Xn, Xm, results := aarch64translation.BinaryInstructionToAarch64(info)
	if !results.IsEmpty() {
		return nil, results
	}

	return SubsReg{
		SubShiftedRegister: instructions.NewSubsShiftedRegister(Xd, Xn, Xm),
	}, core.ResultList{}
}

func (SubsDefinition) buildImmediateVariant(
	info *gen.InstructionInfo,
) (gen.BaseInstruction, core.ResultList) {
	Xd, Xn, imm, results := aarch64translation.Immediate12GPRegisterTargetInstructionToAarch64(info)
	if !results.IsEmpty() {
		return nil, results
	}

	return SubsImm{
		SubImmediate: instructions.NewSubsImmediate(Xd, Xn, imm),
	}, core.ResultList{}
}

func (d SubsDefinition) BuildInstruction(
	info *gen.InstructionInfo,
) (gen.BaseInstruction, core.ResultList) {
	results := aarch64translation.ValidateBinaryInstruction(info)
	if !results.IsEmpty() {
		return nil, results
	}

	switch info.Arguments[1].(type) {
	case *gen.RegisterArgumentInfo:
		return d.buildRegisterVariant(info)

	case *gen.ImmediateInfo:
		return d.buildImmediateVariant(info)

	default:
		return nil, list.FromSingle(core.Result{
			{
				Type:     core.ErrorResult,
				Message:  "Second \"subs\" argument must be a register or immediate",
				Location: info.Arguments[1].Declaration(),
			},
		})
	}
}

func NewSubsInstructionDefinition() gen.InstructionDefinition {
	return SubsDefinition{}
}
