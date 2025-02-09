package aarch64isa

import (
	"alon.kr/x/aarch64codegen/instructions"
	aarch64translation "alon.kr/x/usm/aarch64/translation"
	"alon.kr/x/usm/core"
	"alon.kr/x/usm/gen"
)

type Movz struct {
	instructions.Movz
}

func (Movz) String() string {
	return "MOVZ"
}

func (Movz) PossibleNextSteps() (gen.StepInfo, core.ResultList) {
	return gen.StepInfo{PossibleContinue: true}, core.ResultList{}
}

type MovzDefinition struct{}

func (MovzDefinition) BuildInstruction(
	info *gen.InstructionInfo,
) (gen.BaseInstruction, core.ResultList) {
	results := aarch64translation.AssertArgumentsBetween(info, 2, 2)
	if !results.IsEmpty() {
		return nil, results
	}

	Xd, curResults := aarch64translation.ArgumentToAarch64GPRegister(
		info.Arguments[0],
	)
	results.Extend(&curResults)

	imm, curResults := aarch64translation.ArgumentToAarch64Immediate16(
		info.Arguments[1],
	)
	results.Extend(&curResults)

	shift := instructions.MovShift0
	if len(info.Arguments) > 2 {
		shift, curResults = aarch64translation.ArgumentToAarch64MovShift(
			info.Arguments[2],
		)
		results.Extend(&curResults)
	}

	if !results.IsEmpty() {
		return nil, results
	}

	return Movz{
		instructions.MOVZ(Xd, imm, shift),
	}, core.ResultList{}
}
