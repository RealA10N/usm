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

type BaseAdd struct{}

func (BaseAdd) Operator() string {
	return "ADD"
}

func (BaseAdd) PossibleNextSteps() (gen.StepInfo, core.ResultList) {
	return gen.StepInfo{PossibleContinue: true}, core.ResultList{}
}

type Add struct {
	BaseAdd
	instructions.Add
}

func (i Add) Generate(
	*aarch64codegen.InstructionCodegenContext,
) (instructions.Instruction, core.ResultList) {
	return i, core.ResultList{}
}

type AddImm struct {
	BaseAdd
	instructions.AddImm
}

func (i AddImm) Generate(
	*aarch64codegen.InstructionCodegenContext,
) (instructions.Instruction, core.ResultList) {
	return i, core.ResultList{}
}

type AddDefinition struct{}

func (AddDefinition) buildRegisterVariant(
	info *gen.InstructionInfo,
) (gen.BaseInstruction, core.ResultList) {
	results := core.ResultList{}

	Xd, curResults := aarch64translation.TargetToAarch64GPRegister(info.Targets[0])
	results.Extend(&results)

	Xn, curResults := aarch64translation.ArgumentToAarch64GPRegister(info.Arguments[0])
	results.Extend(&curResults)

	// TODO: Add shifted register support
	Xm, results := aarch64translation.ArgumentToAarch64GPRegister(info.Arguments[1])
	results.Extend(&curResults)

	if !results.IsEmpty() {
		return nil, results
	}

	return Add{
		Add: instructions.ADD(Xd, Xn, Xm, immediates.DoNotSetFlags),
	}, core.ResultList{}
}

func (AddDefinition) buildImmediateVariant(
	info *gen.InstructionInfo,
) (gen.BaseInstruction, core.ResultList) {

	results := core.ResultList{}

	Xd, curResults := aarch64translation.TargetToAarch64GPorSPRegister(info.Targets[0])
	results.Extend(&results)

	Xn, curResults := aarch64translation.ArgumentToAarch64GPorSPRegister(info.Arguments[0])
	results.Extend(&curResults)

	// TODO: Add shifted immediate support
	imm, curResults := aarch64translation.ArgumentToAarch64Immediate12(info.Arguments[1])
	results.Extend(&curResults)

	if !results.IsEmpty() {
		return nil, results
	}

	return AddImm{
		AddImm: instructions.ADDI(Xd, Xn, imm, immediates.DoNotSetFlags),
	}, core.ResultList{}
}

func (d AddDefinition) BuildInstruction(
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
				Message:  "Second ADD argument must be a register or immediate",
				Location: info.Arguments[1].Declaration(),
			},
		})
	}
}

func NewAddInstructionDefinition() gen.InstructionDefinition {
	return AddDefinition{}
}
