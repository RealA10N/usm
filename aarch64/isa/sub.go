package aarch64isa

import (
	"alon.kr/x/aarch64codegen/instructions"
	"alon.kr/x/list"
	aarch64codegen "alon.kr/x/usm/aarch64/codegen"
	aarch64translation "alon.kr/x/usm/aarch64/translation"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
)

type BaseSub struct {
	NonBranchingInstruction
}

func (BaseSub) Operator() string {
	return "sub"
}

type SubReg struct {
	BaseSub
	instructions.SubShiftedRegister
}

func (i SubReg) Generate(
	*aarch64codegen.InstructionCodegenContext,
) (instructions.Instruction, core.ResultList) {
	return i, core.ResultList{}
}

type SubImm struct {
	BaseSub
	instructions.SubImmediate
}

func (i SubImm) Generate(
	*aarch64codegen.InstructionCodegenContext,
) (instructions.Instruction, core.ResultList) {
	return i, core.ResultList{}
}

type SubDefinition struct{}

func (SubDefinition) buildRegisterVariant(
	info *gen.InstructionInfo,
) (gen.BaseInstruction, core.ResultList) {
	Xd, Xn, Xm, results := aarch64translation.BinaryInstructionToAarch64(info)
	if !results.IsEmpty() {
		return nil, results
	}

	return SubReg{
		SubShiftedRegister: instructions.NewSubShiftedRegister(Xd, Xn, Xm),
	}, core.ResultList{}
}

func (SubDefinition) buildImmediateVariant(
	info *gen.InstructionInfo,
) (gen.BaseInstruction, core.ResultList) {
	Xd, Xn, imm, results := aarch64translation.Immediate12GPRegisterTargetInstructionToAarch64(info)
	if !results.IsEmpty() {
		return nil, results
	}

	return SubImm{
		SubImmediate: instructions.NewSubImmediate(Xd, Xn, imm),
	}, core.ResultList{}
}

func (d SubDefinition) BuildInstruction(
	info *gen.InstructionInfo,
) (gen.BaseInstruction, core.ResultList) {
	results := core.ResultList{}

	curResults := aarch64translation.AssertTargetsExactly(info, 1)
	results.Extend(&curResults)

	curResults = aarch64translation.AssertArgumentsExactly(info, 2)
	results.Extend(&curResults)

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
				Message:  "Second \"sub\" argument must be a register or immediate",
				Location: info.Arguments[1].Declaration(),
			},
		})
	}
}

func NewSubInstructionDefinition() gen.InstructionDefinition {
	return SubDefinition{}
}
