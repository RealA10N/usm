package aarch64isa

import (
	"alon.kr/x/aarch64codegen/instructions"
	"alon.kr/x/list"
	aarch64codegen "alon.kr/x/usm/aarch64/codegen"
	aarch64translation "alon.kr/x/usm/aarch64/translation"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
)

type BaseAdds struct {
	NonBranchingInstruction
}

func (BaseAdds) Operator() string {
	return "adds"
}

type AddsReg struct {
	BaseAdd
	instructions.AddShiftedRegister
}

func (i AddsReg) Generate(
	*aarch64codegen.InstructionCodegenContext,
) (instructions.Instruction, core.ResultList) {
	return i, core.ResultList{}
}

type AddsImm struct {
	BaseAdd
	instructions.AddsImmediate
}

func (i AddsImm) Generate(
	*aarch64codegen.InstructionCodegenContext,
) (instructions.Instruction, core.ResultList) {
	return i, core.ResultList{}
}

type AddsDefinition struct{}

func (d AddsDefinition) buildRegisterVariant(
	info *gen.InstructionInfo,
) (gen.BaseInstruction, core.ResultList) {
	Xd, Xn, Xm, results := aarch64translation.BinaryInstructionToAarch64(info)
	if !results.IsEmpty() {
		return nil, results
	}

	return AddReg{
		AddShiftedRegister: instructions.NewAddsShiftedRegister(Xd, Xn, Xm),
	}, core.ResultList{}
}

func (AddsDefinition) buildImmediateVariant(
	info *gen.InstructionInfo,
) (gen.BaseInstruction, core.ResultList) {
	Xd, Xn, imm, results := aarch64translation.Immediate12GPRegisterTargetInstructionToAarch64(info)
	if !results.IsEmpty() {
		return nil, results
	}

	return AddsImm{
		AddsImmediate: instructions.NewAddsImmediate(Xd, Xn, imm),
	}, core.ResultList{}
}

func (d AddsDefinition) BuildInstruction(
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
				Message:  "Second \"adds\" argument must be a register or immediate",
				Location: info.Arguments[1].Declaration(),
			},
		})
	}
}

func NewAddsInstructionDefinition() gen.InstructionDefinition {
	return AddsDefinition{}
}
