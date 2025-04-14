package aarch64isa

import (
	"alon.kr/x/aarch64codegen/instructions"
	"alon.kr/x/list"
	aarch64translation "alon.kr/x/usm/aarch64/translation"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
)

type Add struct {
	instructions.Add
}

func (Add) String() string {
	return "ADD"
}

func (Add) PossibleNextSteps() (gen.StepInfo, core.ResultList) {
	return gen.StepInfo{PossibleContinue: true}, core.ResultList{}
}

type AddImm struct {
	instructions.AddImm
}

func (AddImm) String() string {
	return "ADD"
}

func (AddImm) PossibleNextSteps() (gen.StepInfo, core.ResultList) {
	return gen.StepInfo{PossibleContinue: true}, core.ResultList{}
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

	return Add{instructions.ADD(Xd, Xn, Xm)}, core.ResultList{}
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

	return AddImm{instructions.ADDI(Xd, Xn, imm)}, core.ResultList{}
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
