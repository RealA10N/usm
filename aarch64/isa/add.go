package aarch64isa

import (
	"alon.kr/x/aarch64codegen/instructions"
	aarch64translation "alon.kr/x/usm/aarch64/translation"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
)

type Add struct {
	instructions.AddImm
}

func (Add) String() string {
	return "ADD"
}

func (Add) PossibleNextSteps() (gen.StepInfo, core.ResultList) {
	return gen.StepInfo{PossibleContinue: true}, core.ResultList{}
}

type AddDefinition struct{}

func (AddDefinition) BuildInstruction(
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

	Xd, curResults := aarch64translation.TargetToAarch64GPorSPRegister(info.Targets[0])
	results.Extend(&curResults)

	Xn, curRecurResults := aarch64translation.ArgumentToAarch64GPorSPRegister(info.Arguments[0])
	results.Extend(&curRecurResults)

	imm, curResults := aarch64translation.ArgumentToAarch64Immediate12(info.Arguments[1])
	results.Extend(&curResults)

	// TODO: Add LSL #12 support.
	// TODO: Add support for other ADD variants.

	if !results.IsEmpty() {
		return nil, results
	}

	return Add{
		instructions.ADDI(Xd, Xn, imm),
	}, core.ResultList{}
}

func NewAddInstructionDefinition() gen.InstructionDefinition {
	return AddDefinition{}
}
