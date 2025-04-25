package aarch64isa

import (
	"alon.kr/x/aarch64codegen/instructions"
	"alon.kr/x/list"
	aarch64codegen "alon.kr/x/usm/aarch64/codegen"
	aarch64translation "alon.kr/x/usm/aarch64/translation"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
)

type BaseAdd struct{}

func (BaseAdd) Operator() string {
	return "add"
}

func (BaseAdd) PossibleNextSteps() (gen.StepInfo, core.ResultList) {
	return gen.StepInfo{PossibleContinue: true}, core.ResultList{}
}

type AddReg struct {
	BaseAdd
	instructions.AddShiftedRegister
}

func (i AddReg) Generate(
	*aarch64codegen.InstructionCodegenContext,
) (instructions.Instruction, core.ResultList) {
	return i, core.ResultList{}
}

type AddImm struct {
	BaseAdd
	instructions.AddImmediate
}

func (i AddImm) Generate(
	*aarch64codegen.InstructionCodegenContext,
) (instructions.Instruction, core.ResultList) {
	return i, core.ResultList{}
}

type AddDefinition struct{}

func (d AddDefinition) buildRegisterVariant(
	info *gen.InstructionInfo,
) (gen.BaseInstruction, core.ResultList) {
	Xd, Xn, Xm, results := aarch64translation.BinaryInstructionToAarch64(info)
	if !results.IsEmpty() {
		return nil, results
	}

	return AddReg{
		AddShiftedRegister: instructions.NewAddShiftedRegister(Xd, Xn, Xm),
	}, core.ResultList{}
}

func (AddDefinition) buildImmediateVariant(
	info *gen.InstructionInfo,
) (gen.BaseInstruction, core.ResultList) {
	Xd, Xn, imm, results := aarch64translation.Immediate12InstructionToAarch64(info)
	if !results.IsEmpty() {
		return nil, results
	}

	return AddImm{
		AddImmediate: instructions.NewAddImmediate(Xd, Xn, imm),
	}, core.ResultList{}
}

func (d AddDefinition) BuildInstruction(
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
				Message:  "Second \"add\" argument must be a register or immediate",
				Location: info.Arguments[1].Declaration(),
			},
		})
	}
}

func NewAddInstructionDefinition() gen.InstructionDefinition {
	return AddDefinition{}
}
