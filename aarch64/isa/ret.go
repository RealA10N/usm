package aarch64isa

import (
	"alon.kr/x/aarch64codegen/instructions"
	"alon.kr/x/aarch64codegen/registers"
	aarch64codegen "alon.kr/x/usm/aarch64/codegen"
	aarch64translation "alon.kr/x/usm/aarch64/translation"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
)

type Ret struct {
	instructions.Ret
}

func (Ret) Operator() string {
	return "ret"
}

func (Ret) PossibleNextSteps() (gen.StepInfo, core.ResultList) {
	return gen.StepInfo{PossibleReturn: true}, core.ResultList{}
}

func (i Ret) Generate(
	*aarch64codegen.InstructionCodegenContext,
) (instructions.Instruction, core.ResultList) {
	return i, core.ResultList{}
}

type RetDefinition struct{}

func (RetDefinition) BuildInstruction(
	info *gen.InstructionInfo,
) (gen.BaseInstruction, core.ResultList) {
	results := core.ResultList{}

	curResults := aarch64translation.AssertArgumentsBetween(info, 0, 1)
	results.Extend(&curResults)

	curResults = aarch64translation.AssertTargetsExactly(info, 0)
	results.Extend(&curResults)

	if !results.IsEmpty() {
		return nil, results
	}

	Xn := registers.GPRegisterX30
	if len(info.Arguments) > 0 {
		Xn, curResults = aarch64translation.ArgumentToAarch64GPRegister(info.Arguments[0])
		results.Extend(&curResults)
	}

	if !results.IsEmpty() {
		return nil, results
	}

	return Ret{
		instructions.RET(Xn),
	}, core.ResultList{}
}

func NewRetInstructionDefinition() gen.InstructionDefinition {
	return RetDefinition{}
}
